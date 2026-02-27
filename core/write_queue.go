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
	DataUpdate map[string]map[string]map[int]map[int]interface{}
	DataAppend map[string]map[string][][]interface{} 
}

var BoNhoGhi = &SmartQueue{
	DataUpdate: make(map[string]map[string]map[int]map[int]interface{}),
	DataAppend: make(map[string]map[string][][]interface{}),
}

var KenhBaoThuc = make(chan struct{}, 1)
const ChuKyGhiSheet = 5 * time.Second

func ThemVaoHangCho(spreadId string, sheetName string, row int, col int, value interface{}) {
	BoNhoGhi.Lock()
	defer BoNhoGhi.Unlock()

	if BoNhoGhi.DataUpdate[spreadId] == nil {
		BoNhoGhi.DataUpdate[spreadId] = make(map[string]map[int]map[int]interface{})
	}
	if BoNhoGhi.DataUpdate[spreadId][sheetName] == nil {
		BoNhoGhi.DataUpdate[spreadId][sheetName] = make(map[int]map[int]interface{})
	}
	if BoNhoGhi.DataUpdate[spreadId][sheetName][row] == nil {
		BoNhoGhi.DataUpdate[spreadId][sheetName][row] = make(map[int]interface{})
	}

	BoNhoGhi.DataUpdate[spreadId][sheetName][row][col] = value

	select {
	case KenhBaoThuc <- struct{}{}:
	default:
	}
}

func ThemDongVaoHangCho(spreadId string, sheetName string, rowData []interface{}) {
	BoNhoGhi.Lock()
	defer BoNhoGhi.Unlock()

	if BoNhoGhi.DataAppend[spreadId] == nil {
		BoNhoGhi.DataAppend[spreadId] = make(map[string][][]interface{})
	}
	
	BoNhoGhi.DataAppend[spreadId][sheetName] = append(BoNhoGhi.DataAppend[spreadId][sheetName], rowData)

	select {
	case KenhBaoThuc <- struct{}{}:
	default:
	}
}

func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Printf("🚀 [CORE WORKER] Đã khởi động. Chế độ Kép: Update & Append (%v).", ChuKyGhiSheet)
		for {
			<-KenhBaoThuc
			time.Sleep(ChuKyGhiSheet)
			ThucHienGhiSheet()
		}
	}()
}

func ThucHienGhiSheet() {
	BoNhoGhi.Lock()
	if len(BoNhoGhi.DataUpdate) == 0 && len(BoNhoGhi.DataAppend) == 0 {
		BoNhoGhi.Unlock()
		return
	}

	snapshotUpdate := BoNhoGhi.DataUpdate
	snapshotAppend := BoNhoGhi.DataAppend
	
	BoNhoGhi.DataUpdate = make(map[string]map[string]map[int]map[int]interface{})
	BoNhoGhi.DataAppend = make(map[string]map[string][][]interface{})
	BoNhoGhi.Unlock()

	log.Println("⚡ [SMART BATCH] Đang xử lý đồng bộ dữ liệu kép xuống Google Sheets...")

	allSpreadIDs := make(map[string]bool)
	for id := range snapshotUpdate { allSpreadIDs[id] = true }
	for id := range snapshotAppend { allSpreadIDs[id] = true }

	for spreadId := range allSpreadIDs {
		srv := LayDichVuSheet(spreadId)
		if srv == nil {
			log.Printf("❌ LỖI GHI %s: Không tìm thấy kết nối Google API", spreadId[:5])
			continue
		}

		// LUỒNG 1: XỬ LÝ UPDATE
		if sheetsMap, ok := snapshotUpdate[spreadId]; ok && len(sheetsMap) > 0 {
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
				_, err := srv.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
				if err != nil {
					log.Printf("❌ LỖI GHI UPDATE %s: %v", spreadId[:5], err)
				} else {
					log.Printf("✅ Đã ghi UPDATE %d dải dữ liệu vào Sheet %s", len(requests), spreadId[:5])
				}
			}
		}

		// LUỒNG 2: XỬ LÝ APPEND
		if appendSheets, ok := snapshotAppend[spreadId]; ok && len(appendSheets) > 0 {
			for sheetName, rowsData := range appendSheets {
				
				vr := &sheets.ValueRange{
					Values: rowsData,
				}

				// [ĐÃ FIX LỖI "LỢN QUÈ" Ở ĐÂY]: Đổi USER_ENTERED thành RAW
				_, err := srv.Spreadsheets.Values.Append(spreadId, sheetName, vr).
					ValueInputOption("RAW"). 
					InsertDataOption("INSERT_ROWS").
					Do()

				if err != nil {
					log.Printf("❌ LỖI GHI APPEND %s (Tab: %s): %v", spreadId[:5], sheetName, err)
				} else {
					log.Printf("✅ Đã APPEND %d dòng mới vào Tab %s (Sheet %s)", len(rowsData), sheetName, spreadId[:5])
				}
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
