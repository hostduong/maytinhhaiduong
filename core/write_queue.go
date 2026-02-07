package core

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"google.golang.org/api/sheets/v4"
)

// =============================================================
// 1. CẤU TRÚC DỮ LIỆU THÔNG MINH (SMART QUEUE)
// =============================================================

// Map 4 cấp: [SpreadsheetID] -> [SheetName] -> [Row] -> [Col] -> Value
type SmartQueue struct {
	sync.Mutex
	Data map[string]map[string]map[int]map[int]interface{}
}

// Bộ nhớ đệm RAM
var BoNhoGhi = &SmartQueue{
	Data: make(map[string]map[string]map[int]map[int]interface{}),
}

// Kênh báo thức Worker (Buffer 1 để tránh block)
var KenhBaoThuc = make(chan struct{}, 1)

// Chu kỳ ghi (Hardcode hoặc lấy từ Config)
const ChuKyGhiSheet = 5 * time.Second

// =============================================================
// 2. HÀM PUBLIC: ĐẨY DỮ LIỆU VÀO HÀNG CHỜ
// =============================================================

// Hàm này thay thế cho nghiep_vu.ThemVaoHangCho cũ
func ThemVaoHangCho(spreadId string, sheetName string, row int, col int, value interface{}) {
	BoNhoGhi.Lock()
	defer BoNhoGhi.Unlock()

	// Init Map nếu chưa có
	if BoNhoGhi.Data[spreadId] == nil {
		BoNhoGhi.Data[spreadId] = make(map[string]map[int]map[int]interface{})
	}
	if BoNhoGhi.Data[spreadId][sheetName] == nil {
		BoNhoGhi.Data[spreadId][sheetName] = make(map[int]map[int]interface{})
	}
	if BoNhoGhi.Data[spreadId][sheetName][row] == nil {
		BoNhoGhi.Data[spreadId][sheetName][row] = make(map[int]interface{})
	}

	// Ghi vào RAM
	BoNhoGhi.Data[spreadId][sheetName][row][col] = value

	// Bắn tín hiệu đánh thức Worker (Non-blocking)
	select {
	case KenhBaoThuc <- struct{}{}:
	default:
	}
}

// =============================================================
// 3. WORKER THÔNG MINH (CHẠY NGẦM)
// =============================================================

func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Printf("🚀 [CORE WORKER] Đã khởi động. Chế độ Smart Batch (%v).", ChuKyGhiSheet)
		
		for {
			// A. NGỦ ĐÔNG: Chờ tín hiệu
			<-KenhBaoThuc
			
			// B. GOM HÀNG (Debounce): Chờ 5s để gom thêm lệnh
			time.Sleep(ChuKyGhiSheet)

			// C. THỰC THI
			ThucHienGhiSheet()
		}
	}()
}

// =============================================================
// 4. LOGIC XỬ LÝ GHI (BATCH UPDATE)
// =============================================================

func ThucHienGhiSheet() {
	BoNhoGhi.Lock()
	if len(BoNhoGhi.Data) == 0 {
		BoNhoGhi.Unlock()
		return
	}

	// 1. Snapshot dữ liệu & Reset RAM
	snapshotData := BoNhoGhi.Data
	BoNhoGhi.Data = make(map[string]map[string]map[int]map[int]interface{})
	BoNhoGhi.Unlock()

	log.Println("⚡ [SMART BATCH] Đang đẩy dữ liệu xuống Google Sheet...")

	// 2. Duyệt qua từng File Sheet (SpreadsheetID)
	for spreadId, sheetsMap := range snapshotData {
		var requests []*sheets.ValueRange
		
		for sheetName, rows := range sheetsMap {
			for r, cols := range rows {
				// --- THUẬT TOÁN GOM CỘT LIỀN KỀ (Của bạn) ---
				var colIndexes []int
				for c := range cols { colIndexes = append(colIndexes, c) }
				sort.Ints(colIndexes)

				if len(colIndexes) == 0 { continue }
				
				// Tìm dải liên tục: A, B, C -> A:C
				startCol := colIndexes[0]
				prevCol := colIndexes[0]
				currentValues := []interface{}{cols[startCol]}

				for i := 1; i < len(colIndexes); i++ {
					currCol := colIndexes[i]
					if currCol == prevCol+1 { // Liền kề
						currentValues = append(currentValues, cols[currCol])
						prevCol = currCol
					} else { // Ngắt quãng -> Đóng gói dải cũ
						rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
						vr := &sheets.ValueRange{
							Range: rangeStr, 
							Values: [][]interface{}{currentValues},
						}
						requests = append(requests, vr)

						// Reset dải mới
						startCol = currCol
						prevCol = currCol
						currentValues = []interface{}{cols[currCol]}
					}
				}
				// Đóng gói dải cuối cùng
				if len(currentValues) > 0 {
					rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
					vr := &sheets.ValueRange{
						Range: rangeStr, 
						Values: [][]interface{}{currentValues},
					}
					requests = append(requests, vr)
				}
			}
		}

		// 3. Gửi Batch Update cho từng File
		if len(requests) > 0 {
			req := &sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             requests,
			}
			
			// Gọi dịch vụ Sheet (Biến toàn cục trong common.go)
			_, err := DichVuSheet.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
			
			if err != nil {
				log.Printf("❌ LỖI GHI %s: %v", spreadId[:5], err)
				// TODO: Logic Retry nếu cần thiết
			} else {
				log.Printf("✅ Đã ghi %d dải dữ liệu vào Sheet %s", len(requests), spreadId[:5])
			}
		}
	}
}

// Helper: 0 -> A, 1 -> B...
func layTenCot(i int) string {
	const text = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 0 { return "A" }
	if i < 26 {
		return string(text[i])
	}
	return string(text[i/26-1]) + string(text[i%26])
}
