package core

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"google.golang.org/api/sheets/v4"
)

// KHÔNG CÓ TỪ "HOẶC" - CHỈ CÓ 2 ACTION DUY NHẤT
const (
	ActionUpdate = "UPDATE"
	ActionAppend = "APPEND"
)

// QueueJob: Định dạng lệnh duy nhất được đẩy xuống Google Sheets
type QueueJob struct {
	ShopID      string
	SheetName   string
	Action      string
	
	// Dùng cho UPDATE
	Row         int
	Col         int
	Value       interface{}
	
	// Dùng cho APPEND
	RowData     []interface{}
}

type WriteQueueManager struct {
	mu    sync.Mutex
	Jobs  []QueueJob
}

var (
	Queue       = &WriteQueueManager{Jobs: make([]QueueJob, 0)}
	WakeUpQueue = make(chan struct{}, 1)
)

// PushUpdate: Dùng để sửa 1 Ô dữ liệu
func PushUpdate(shopID, sheetName string, row, col int, value interface{}) {
	Queue.mu.Lock()
	Queue.Jobs = append(Queue.Jobs, QueueJob{
		ShopID: shopID, SheetName: sheetName, Action: ActionUpdate,
		Row: row, Col: col, Value: value,
	})
	Queue.mu.Unlock()
	TriggerWorker()
}

// PushAppend: Dùng để chèn 1 Dòng dữ liệu mới (Tuyệt đối không đè định dạng)
func PushAppend(shopID, sheetName string, rowData []interface{}) {
	Queue.mu.Lock()
	Queue.Jobs = append(Queue.Jobs, QueueJob{
		ShopID: shopID, SheetName: sheetName, Action: ActionAppend,
		RowData: rowData,
	})
	Queue.mu.Unlock()
	TriggerWorker()
}

func TriggerWorker() {
	select {
	case WakeUpQueue <- struct{}{}:
	default:
	}
}

func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Println("🚀 [CORE WORKER] Đã khởi động đường ống Ghi dữ liệu (Append & Update)...")
		for {
			<-WakeUpQueue
			time.Sleep(5 * time.Second) // Gom mẻ (Batching window)
			ProcessQueue()
		}
	}()
}

func ProcessQueue() {
	Queue.mu.Lock()
	if len(Queue.Jobs) == 0 {
		Queue.mu.Unlock()
		return
	}
	currentJobs := Queue.Jobs
	Queue.Jobs = make([]QueueJob, 0)
	Queue.mu.Unlock()

	log.Printf("⚡ [QUEUE] Đang xử lý mẻ %d tác vụ...", len(currentJobs))

	// Gom nhóm theo ShopID
	jobsByShop := make(map[string][]QueueJob)
	for _, job := range currentJobs {
		jobsByShop[job.ShopID] = append(jobsByShop[job.ShopID], job)
	}

	for shopID, jobs := range jobsByShop {
		srv := LayDichVuSheet(shopID) // Hàm này phải nằm trong utils.go hoặc sheet_driver.go
		if srv == nil {
			log.Printf("❌ [QUEUE] Lỗi: Không lấy được API Google cho Shop %s", shopID[:5])
			continue
		}

		var updateRequests []*sheets.ValueRange

		for _, job := range jobs {
			if job.Action == ActionAppend {
				// Xử lý APPEND (Ghi chèn dòng cuối)
				rangeToAppend := fmt.Sprintf("%s!A:Z", job.SheetName)
				vr := &sheets.ValueRange{Values: [][]interface{}{job.RowData}}
				
				_, err := srv.Spreadsheets.Values.Append(shopID, rangeToAppend, vr).
					ValueInputOption("RAW").InsertDataOption("OVERWRITE").Do()
				
				if err != nil {
					log.Printf("❌ [APPEND ERROR] Shop %s - Sheet %s: %v", shopID[:5], job.SheetName, err)
				}
			} else if job.Action == ActionUpdate {
				// Gom Update (Ghi đè ô)
				rangeStr := fmt.Sprintf("%s!%s%d", job.SheetName, layTenCot(job.Col), job.Row)
				updateRequests = append(updateRequests, &sheets.ValueRange{
					Range:  rangeStr,
					Values: [][]interface{}{{job.Value}},
				})
			}
		}

		// Đẩy toàn bộ Update lên 1 lần (BatchUpdate)
		if len(updateRequests) > 0 {
			req := &sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             updateRequests,
			}
			_, err := srv.Spreadsheets.Values.BatchUpdate(shopID, req).Do()
			if err != nil {
				log.Printf("❌ [UPDATE ERROR] Shop %s: %v", shopID[:5], err)
			}
		}
	}
	log.Println("✅ [QUEUE] Đã giải quyết xong mẻ tác vụ.")
}

func layTenCot(i int) string {
	const text = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 0 { return "A" }
	if i < 26 { return string(text[i]) }
	return string(text[i/26-1]) + string(text[i%26])
}
