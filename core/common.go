package core

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"app/cau_hinh"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// =============================================================
// 1. CẤU TRÚC HẠ TẦNG (INFRASTRUCTURE)
// =============================================================
var (
	// Khóa an toàn (Mutex) bảo vệ toàn bộ dữ liệu RAM
	KhoaHeThong sync.RWMutex

	// Dịch vụ Google Sheets API
	DichVuSheet *sheets.Service

	// Cờ báo hiệu hệ thống đang bận (khi Reload toàn bộ)
	HeThongDangBan bool
)

// Struct phục vụ cho Hàng Chờ Ghi (Write Queue)
type YeuCauGhi struct {
	SpreadsheetID string      // ID file Google Sheet
	SheetName     string      // Tên Sheet
	RowIndex      int         // Dòng cần ghi
	ColIndex      int         // Cột cần ghi
	Value         interface{} // Giá trị cần ghi
}

// Callback để main.go đăng ký hàm xử lý ghi
var CallbackThemVaoHangCho func(req YeuCauGhi)

// =============================================================
// 2. KHỞI TẠO KẾT NỐI
// =============================================================
func KhoiTaoNenTang() {
	log.Println("🔌 [CORE] Đang kết nối Google Sheets (Chế độ Đa Nhiệm)...")

	ctx := context.Background()
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson

	// Nếu không có JSON Key thì báo lỗi hoặc tìm cách xử lý
	if jsonKey == "" {
		log.Println("⚠️ CẢNH BÁO: Chưa cấu hình GoogleAuthJson (Env: GOOGLE_JSON_KEY).")
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		log.Fatalf("❌ LỖI KẾT NỐI GOOGLE SHEETS: %v", err)
	}

	DichVuSheet = srv
	log.Println("✅ [CORE] Kết nối thành công!")
}

// =============================================================
// 3. HÀM TIỆN ÍCH CỐT LÕI (HELPER)
// =============================================================

// Tạo khóa duy nhất trong RAM: "SheetID__EntityID"
func TaoCompositeKey(sheetID, entityID string) string {
	return fmt.Sprintf("%s__%s", sheetID, entityID)
}

// Hàm đọc dữ liệu hỗ trợ chỉ định ID File (QUAN TRỌNG: 2 THAM SỐ)
func loadSheetData(spreadsheetID string, tenSheet string) ([][]interface{}, error) {
	// Nếu không truyền ID, lấy ID mặc định trong Config
	if spreadsheetID == "" {
		spreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	readRange := tenSheet + "!A:AZ" // Đọc rộng
	resp, err := DichVuSheet.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Printf("⚠️ Lỗi đọc sheet %s (ID: %s...): %v", tenSheet, spreadsheetID[:5], err)
		return nil, err
	}
	return resp.Values, nil
}

// --- CÁC HÀM PARSE DỮ LIỆU ---

func layString(row []interface{}, index int) string {
	if index >= len(row) || row[index] == nil { return "" }
	return strings.TrimSpace(fmt.Sprintf("%v", row[index]))
}

func layInt(row []interface{}, index int) int {
	s := layString(row, index)
	if s == "" { return 0 }
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	val, _ := strconv.Atoi(s)
	return val
}

func layFloat(row []interface{}, index int) float64 {
	s := layString(row, index)
	if s == "" { return 0 }
	s = strings.ReplaceAll(s, "đ", "")
	s = strings.ReplaceAll(s, "USD", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
