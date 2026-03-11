package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/sheets/v4"
)

// ==============================================================================
// 1. CẤU TRÚC HÀNG CHỜ
// ==============================================================================
const (
	ActionUpdate     = "UPDATE"
	ActionAppend     = "APPEND"
	ActionSmartSync  = "SMART_SYNC" // Hành động mới: Đồng bộ nguyên dòng JSON 2 cột
)

type QueueJob struct {
	ShopID      string
	SheetName   string
	Action      string
	
	// Dùng cho Cổ điển (Nhiều cột)
	Row         int
	Col         int
	Value       interface{}
	RowData     []interface{}
	
	// Dùng cho Smart Queue (NoSQL 2 cột)
	ObjectID    string // ID Sản phẩm, Khách hàng...
}

type WriteQueueManager struct {
	mu    sync.Mutex
	Jobs  []QueueJob
}

var (
	Queue       = &WriteQueueManager{Jobs: make([]QueueJob, 0)}
	WakeUpQueue = make(chan struct{}, 1)
)

// ==============================================================================
// 2. CÁC HÀM BƠM DATA VÀO QUEUE (PRODUCER)
// ==============================================================================

// Cổ điển: Dùng để sửa 1 Ô dữ liệu (Các bảng chưa chuyển JSON)
func PushUpdate(shopID, sheetName string, row, col int, value interface{}) {
	Queue.mu.Lock()
	Queue.Jobs = append(Queue.Jobs, QueueJob{
		ShopID: shopID, SheetName: sheetName, Action: ActionUpdate,
		Row: row, Col: col, Value: value,
	})
	Queue.mu.Unlock()
	TriggerWorker()
}

// Cổ điển: Dùng để chèn 1 Dòng mới (Tuyệt đối không đè định dạng)
func PushAppend(shopID, sheetName string, rowData []interface{}) {
	Queue.mu.Lock()
	Queue.Jobs = append(Queue.Jobs, QueueJob{
		ShopID: shopID, SheetName: sheetName, Action: ActionAppend,
		RowData: rowData,
	})
	Queue.mu.Unlock()
	TriggerWorker()
}

// [MỚI] Tương lai: Bơm một tín hiệu "Báo mộng" (Smart Sync)
func GhiChuDongBo(shopID, tableName, action, objectID string) {
	// Khi gọi hàm này, chúng ta khóa RAM lại để tránh bị xóa
	TangTaskQueue(shopID)
	
	Queue.mu.Lock()
	Queue.Jobs = append(Queue.Jobs, QueueJob{
		ShopID: shopID, SheetName: tableName, Action: ActionSmartSync,
		ObjectID: objectID,
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

// ==============================================================================
// 3. ĐỘNG CƠ XỬ LÝ NGẦM (CONSUMER)
// ==============================================================================
func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Println("🚀 [CORE WORKER] Đã khởi động Hàng Chờ Hỗn Hợp (Classic + Smart)...")
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
	Queue.Jobs = make([]QueueJob, 0) // Làm rỗng Queue cũ để nạp đợt mới
	Queue.mu.Unlock()

	log.Printf("⚡ [QUEUE] Đang xử lý mẻ %d tác vụ hỗn hợp...", len(currentJobs))

	jobsByShop := make(map[string][]QueueJob)
	
	// [MỚI] Bộ lọc Debounce cho Smart Queue (Xóa bỏ các lệnh SmartSync bị trùng lặp)
	smartSyncMap := make(map[string]QueueJob) 
	
	for _, job := range currentJobs {
		if job.Action == ActionSmartSync {
			// Ghi đè: Chỉ giữ lại hành động mới nhất của Object đó
			key := job.ShopID + "_" + job.SheetName + "_" + job.ObjectID
			smartSyncMap[key] = job
		} else {
			jobsByShop[job.ShopID] = append(jobsByShop[job.ShopID], job)
		}
	}
	
	// Đẩy ngược lại Smart Sync vào list theo Shop
	for _, job := range smartSyncMap {
		jobsByShop[job.ShopID] = append(jobsByShop[job.ShopID], job)
	}

	// Xử lý nã súng lên Google
	for shopID, jobs := range jobsByShop {
		srv := LayDichVuSheet(shopID)
		if srv == nil {
			log.Printf("❌ [QUEUE] Lỗi: Không lấy được API Google cho Shop %s", shopID[:5])
			// Nhét lại toàn bộ jobs của Shop này vào Queue để lần sau Retry
			Queue.mu.Lock()
			Queue.Jobs = append(Queue.Jobs, jobs...)
			Queue.mu.Unlock()
			continue
		}

		var updateRequests []*sheets.ValueRange
		var retryJobs []QueueJob // Những thằng cần chạy lại

		for _, job := range jobs {
			switch job.Action {
			
			case ActionAppend:
				rangeToAppend := fmt.Sprintf("%s!A:Z", job.SheetName)
				vr := &sheets.ValueRange{Values: [][]interface{}{job.RowData}}
				_, err := srv.Spreadsheets.Values.Append(shopID, rangeToAppend, vr).
					ValueInputOption("RAW").InsertDataOption("OVERWRITE").Do()
				
				if err != nil {
					log.Printf("❌ [APPEND ERROR] Shop %s - %s: %v", shopID[:5], job.SheetName, err)
					retryJobs = append(retryJobs, job)
				}
				
			case ActionUpdate:
				rangeStr := fmt.Sprintf("%s!%s%d", job.SheetName, layTenCot(job.Col), job.Row)
				updateRequests = append(updateRequests, &sheets.ValueRange{
					Range:  rangeStr,
					Values: [][]interface{}{{job.Value}},
				})
				
			case ActionSmartSync:
				// [MỚI] Cơ chế NoSQL 2 cột
				success := thucThiSmartSync(srv, shopID, job)
				if !success {
					retryJobs = append(retryJobs, job) // Lỗi thì Retry
				} else {
					GiamTaskQueue(shopID) // Thành công thì Mở khóa RAM
				}
			}
		}

		// Đẩy toàn bộ Cổ điển Update lên 1 lần (BatchUpdate)
		if len(updateRequests) > 0 {
			req := &sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             updateRequests,
			}
			_, err := srv.Spreadsheets.Values.BatchUpdate(shopID, req).Do()
			if err != nil {
				log.Printf("❌ [UPDATE BATCH ERROR] Shop %s: %v", shopID[:5], err)
				// Khôi phục các job update lại
				for _, j := range jobs {
					if j.Action == ActionUpdate { retryJobs = append(retryJobs, j) }
				}
			}
		}

		// Đưa các lệnh xịt vào lại Queue chờ 5 giây sau nã tiếp
		if len(retryJobs) > 0 {
			Queue.mu.Lock()
			Queue.Jobs = append(Queue.Jobs, retryJobs...)
			Queue.mu.Unlock()
		}
	}
}

// [MỚI] Thực thi lệnh của NoSQL
func thucThiSmartSync(srv *sheets.Service, shopID string, job QueueJob) bool {
	// Lấy danh sách Tên Sheet Sản Phẩm để xác định nó là SP hay KH
	KhoaHeThong.RLock()
	isSanPham := false
	for _, nganh := range CacheDanhSachNganh {
		if nganh.TenSheet == job.SheetName { isSanPham = true; break }
	}
	KhoaHeThong.RUnlock()

	if isSanPham {
		KhoaHeThong.RLock()
		sp, ok := CacheMapSanPham[TaoCompositeKey(shopID, job.ObjectID)]
		KhoaHeThong.RUnlock()

		// Nếu RAM đã xóa, nghĩa là object không tồn tại
		if !ok || sp == nil { return true } 

		b, _ := json.Marshal(sp)
		jsonStr := string(b)
		
		// Luôn luôn UPDATE (Dù thêm mới hay sửa)
		// (Hệ thống Service đã có trách nhiệm tự Append 1 dòng trống mới trước khi gọi SmartSync)
		rangeID := fmt.Sprintf("%s!A%d", job.SheetName, sp.DongTrongSheet)
		rangeJSON := fmt.Sprintf("%s!B%d", job.SheetName, sp.DongTrongSheet)
		
		req := &sheets.BatchUpdateValuesRequest{
			ValueInputOption: "RAW",
			Data: []*sheets.ValueRange{
				{Range: rangeID, Values: [][]interface{}{{sp.MaSanPham}}},
				{Range: rangeJSON, Values: [][]interface{}{{jsonStr}}},
			},
		}
		
		_, err := srv.Spreadsheets.Values.BatchUpdate(shopID, req).Do()
		if err != nil {
			log.Printf("❌ [SMART_SYNC ERROR] %s: %v", job.ObjectID, err)
			return false
		}
		return true
	}

	// Mở rộng sau này cho Khách Hàng, Đơn Hàng...
	return true 
}

func layTenCot(i int) string {
	const text = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 0 { return "A" }
	if i < 26 { return string(text[i]) }
	return string(text[i/26-1]) + string(text[i%26])
}
