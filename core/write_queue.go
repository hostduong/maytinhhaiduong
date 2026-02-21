package core

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"google.golang.org/api/sheets/v4"
)

type SmartQueue struct {
	sync.Mutex
	Data map[string]map[string]map[int]map[int]interface{}
}

var BoNhoGhi = &SmartQueue{
	Data: make(map[string]map[string]map[int]map[int]interface{}),
}

var KenhBaoThuc = make(chan struct{}, 1)
const ChuKyGhiSheet = 5 * time.Second

func ThemVaoHangCho(spreadId string, sheetName string, row int, col int, value interface{}) {
	BoNhoGhi.Lock()
	defer BoNhoGhi.Unlock()

	if BoNhoGhi.Data[spreadId] == nil {
		BoNhoGhi.Data[spreadId] = make(map[string]map[int]map[int]interface{})
	}
	if BoNhoGhi.Data[spreadId][sheetName] == nil {
		BoNhoGhi.Data[spreadId][sheetName] = make(map[int]map[int]interface{})
	}
	if BoNhoGhi.Data[spreadId][sheetName][row] == nil {
		BoNhoGhi.Data[spreadId][sheetName][row] = make(map[int]interface{})
	}

	BoNhoGhi.Data[spreadId][sheetName][row][col] = value

	select {
	case KenhBaoThuc <- struct{}{}:
	default:
	}
}

func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Printf("🚀 [CORE WORKER] Đã khởi động. Chế độ Smart Batch (%v).", ChuKyGhiSheet)
		for {
			<-KenhBaoThuc
			time.Sleep(ChuKyGhiSheet)
			ThucHienGhiSheet()
		}
	}()
}

func ThucHienGhiSheet() {
	BoNhoGhi.Lock()
	if len(BoNhoGhi.Data) == 0 {
		BoNhoGhi.Unlock()
		return
	}

	snapshotData := BoNhoGhi.Data
	BoNhoGhi.Data = make(map[string]map[string]map[int]map[int]interface{})
	BoNhoGhi.Unlock()

	log.Println("⚡ [SMART BATCH] Đang đẩy dữ liệu xuống Google Sheets...")

	for spreadId, sheetsMap := range snapshotData {
		var requests []*sheets.ValueRange
		
		for sheetName, rows := range sheetsMap {
			for r, cols := range rows {
				var colIndexes []int
				for c := range cols { colIndexes = append(colIndexes, c) }
				sort.Ints(colIndexes)

				if len(colIndexes) == 0 { continue }
				
				startCol := colIndexes[0]
				prevCol := colIndexes[0]
				currentValues := []interface{}{cols[startCol]}

				for i := 1; i < len(colIndexes); i++ {
					currCol := colIndexes[i]
					if currCol == prevCol+1 { 
						currentValues = append(currentValues, cols[currCol])
						prevCol = currCol
					} else { 
						rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
						vr := &sheets.ValueRange{ Range: rangeStr, Values: [][]interface{}{currentValues} }
						requests = append(requests, vr)

						startCol = currCol
						prevCol = currCol
						currentValues = []interface{}{cols[currCol]}
					}
				}
				if len(currentValues) > 0 {
					rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
					vr := &sheets.ValueRange{ Range: rangeStr, Values: [][]interface{}{currentValues} }
					requests = append(requests, vr)
				}
			}
		}

		if len(requests) > 0 {
			req := &sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             requests,
			}
			
			// --- [MỚI] TỰ ĐỘNG CHỌN ĐƯỜNG TRUYỀN API THEO SHOP ---
			srv := LayDichVuSheet(spreadId)
			if srv == nil {
				log.Printf("❌ LỖI GHI %s: Không tìm thấy kết nối Google API", spreadId[:5])
				continue
			}

			_, err := srv.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
			
			if err != nil {
				log.Printf("❌ LỖI GHI %s: %v", spreadId[:5], err)
			} else {
				log.Printf("✅ Đã ghi %d dải dữ liệu vào Sheet %s", len(requests), spreadId[:5])
			}
		}
	}
}

func layTenCot(i int) string {
	const text = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 0 { return "A" }
	if i < 26 {
		return string(text[i])
	}
	return string(text[i/26-1]) + string(text[i%26])
}
