package core

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"app/cau_hinh" // Vẫn giữ cấu hình cũ để lấy ID Sheet & JSON Key

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// =============================================================
// 1. BIẾN TOÀN CỤC (INFRASTRUCTURE)
// =============================================================
var (
	// Khóa an toàn cho toàn bộ hệ thống (Thay thế cho QuanLyKhoa cũ)
	KhoaHeThong sync.RWMutex

	// Dịch vụ Google Sheets
	DichVuSheet *sheets.Service

	// Cờ báo hiệu hệ thống đang bận (khi Reload)
	HeThongDangBan bool
	
	// Callback để gọi ngược ra ngoài (nếu cần ghi log hoặc xử lý async)
	CallbackGhiSheet func(bool)
)

// =============================================================
// 2. KHỞI TẠO KẾT NỐI (Gọi 1 lần ở main.go)
// =============================================================
func KhoiTaoNenTang() {
	log.Println("🔌 [CORE] Đang kết nối Google Sheets...")

	ctx := context.Background()
	// Lấy Credentials từ package cau_hinh (giữ nguyên logic cũ)
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		log.Fatalf("❌ LỖI KẾT NỐI GOOGLE SHEETS: %v", err)
	}

	DichVuSheet = srv
	log.Println("✅ [CORE] Kết nối thành công!")
}

// =============================================================
// 3. HÀM ĐỌC DỮ LIỆU (Helper dùng chung cho các file con)
// =============================================================
func loadSheetData(tenSheet string) ([][]interface{}, error) {
	spreadsheetId := cau_hinh.BienCauHinh.IdFileSheet
	readRange := tenSheet + "!A:AZ" // Đọc rộng ra đến cột AZ cho chắc

	resp, err := DichVuSheet.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Printf("⚠️ Lỗi đọc sheet %s: %v", tenSheet, err)
		return nil, err
	}
	return resp.Values, nil
}

// =============================================================
// 4. CÁC HÀM CHUYỂN ĐỔI DỮ LIỆU (Parser)
// =============================================================
func layString(row []interface{}, index int) string {
	if index >= len(row) || row[index] == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", row[index]))
}

func layInt(row []interface{}, index int) int {
	s := layString(row, index)
	if s == "" { return 0 }
	// Xóa dấu chấm/phẩy ngăn cách nghìn (VD: 1.000 -> 1000)
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
	// Xử lý dấu chấm/phẩy tùy theo locale, ở đây ta assume format 1.000.000
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func layBool(row []interface{}, index int) bool {
	s := strings.ToLower(layString(row, index))
	return s == "1" || s == "true" || s == "yes" || s == "co"
}
