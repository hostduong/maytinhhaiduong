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

	// Cờ báo hiệu hệ thống đang bận
	HeThongDangBan bool
)

// Struct phục vụ cho Hàng Chờ Ghi
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
	log.Println("🔌 [CORE] Đang kết nối Google Sheets...")

	ctx := context.Background()
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson

	var srv *sheets.Service
	var err error

	if jsonKey != "" {
		log.Println("🔑 [AUTH] Phát hiện JSON Key, sử dụng chế độ Service Account Key.")
		srv, err = sheets.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	} else {
		log.Println("☁️ [AUTH] Không có JSON Key, chuyển sang chế độ Cloud Run (ADC).")
		srv, err = sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsScope))
	}

	if err != nil {
		log.Printf("❌ LỖI KẾT NỐI GOOGLE SHEETS: %v", err)
		log.Println("⚠️ Hệ thống sẽ chạy ở chế độ Offline (Chỉ xem giao diện, không có dữ liệu).")
		return
	}

	DichVuSheet = srv
	log.Println("✅ [CORE] Kết nối thành công!")
}

// =============================================================
// 3. HÀM TIỆN ÍCH CỐT LÕI (HELPER)
// =============================================================

func TaoCompositeKey(sheetID, entityID string) string {
	return fmt.Sprintf("%s__%s", sheetID, entityID)
}

func loadSheetData(spreadsheetID string, tenSheet string) ([][]interface{}, error) {
	if DichVuSheet == nil {
		return nil, fmt.Errorf("chưa kết nối được Google Sheets")
	}

	if spreadsheetID == "" {
		spreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	readRange := tenSheet + "!A:AZ"
	resp, err := DichVuSheet.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Printf("⚠️ Lỗi đọc sheet %s: %v", tenSheet, err)
		return nil, err
	}
	return resp.Values, nil
}

func ThemVaoHangCho(sheetID, sheetName string, row, col int, val interface{}) {
	if CallbackThemVaoHangCho != nil {
		CallbackThemVaoHangCho(YeuCauGhi{
			SpreadsheetID: sheetID,
			SheetName:     sheetName,
			RowIndex:      row,
			ColIndex:      col,
			Value:         val,
		})
	}
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
