package core

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"app/cau_hinh"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// =============================================================
// 1. CẤU TRÚC HẠ TẦNG (INFRASTRUCTURE)
// =============================================================

var (
	// Khóa an toàn bảo vệ dữ liệu RAM
	KhoaHeThong sync.RWMutex
	HeThongDangBan bool

	// --- [MỚI] BỂ CHỨA KẾT NỐI API (CONNECTION POOL) ---
	MapDichVuSheet = make(map[string]*sheets.Service)
	MutexDichVu    sync.RWMutex // Khóa riêng cho Pool
)

type YeuCauGhi struct {
	SpreadsheetID string      
	SheetName     string      
	RowIndex      int         
	ColIndex      int         
	Value         interface{} 
}

var CallbackThemVaoHangCho func(req YeuCauGhi)

// =============================================================
// 2. KHỞI TẠO KẾT NỐI (SERVER DEFAULT)
// =============================================================
func KhoiTaoNenTang() {
	log.Println("🔌 [CORE] Đang kết nối Google Sheets (API Mặc định)...")

	ctx := context.Background()
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson

	var srv *sheets.Service
	var err error

	if jsonKey != "" {
		log.Println("🔑 [AUTH] Phát hiện JSON Key hệ thống, sử dụng Service Account.")
		srv, err = sheets.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	} else {
		log.Println("☁️ [AUTH] Không có JSON Key, dùng chế độ Cloud Run (ADC).")
		srv, err = sheets.NewService(ctx, option.WithScopes(sheets.SpreadsheetsScope))
	}

	if err != nil {
		log.Printf("❌ LỖI KẾT NỐI MẶC ĐỊNH: %v", err)
		return
	}

	// Lưu API mặc định vào Pool
	MutexDichVu.Lock()
	MapDichVuSheet["default"] = srv
	MapDichVuSheet[cau_hinh.BienCauHinh.IdFileSheet] = srv // Lưu cho Master Shop
	MutexDichVu.Unlock()
	
	log.Println("✅ [CORE] Khởi tạo API mặc định thành công!")
}

// =============================================================
// 3. QUẢN LÝ POOL KẾT NỐI (MULTITENANT API)
// =============================================================

// KetNoiGoogleSheetRieng: Tạo đường truyền API riêng cho Shop VIP
func KetNoiGoogleSheetRieng(shopID string, jsonKey string) {
	if jsonKey == "" || shopID == "" { return }

	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		log.Printf("⚠️ [AUTH] Shop [%s] sai định dạng JSON API: %v", shopID, err)
		return
	}

	MutexDichVu.Lock()
	MapDichVuSheet[shopID] = srv
	MutexDichVu.Unlock()
	log.Printf("🚀 [AUTH] Kích hoạt đường truyền API Riêng (VIP) cho Shop [%s]", shopID)
}

// LayDichVuSheet: Lấy API của shop, nếu ko có thì lấy mặc định
func LayDichVuSheet(shopID string) *sheets.Service {
	MutexDichVu.RLock()
	srv, ok := MapDichVuSheet[shopID]
	MutexDichVu.RUnlock()

	if ok && srv != nil {
		return srv
	}

	// Fallback
	MutexDichVu.RLock()
	defaultSrv := MapDichVuSheet["default"]
	MutexDichVu.RUnlock()
	return defaultSrv
}

// =============================================================
// 4. HÀM TIỆN ÍCH CỐT LÕI (HELPER)
// =============================================================

func TaoCompositeKey(sheetID, entityID string) string {
	return fmt.Sprintf("%s__%s", sheetID, entityID)
}

// Lấy dữ liệu thông minh (Tự tìm đúng API của Shop)
func LoadSheetData(spreadsheetID string, tenSheet string) ([][]interface{}, error) {
	if spreadsheetID == "" {
		spreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	srv := LayDichVuSheet(spreadsheetID)
	if srv == nil {
		return nil, fmt.Errorf("chưa kết nối được Google Sheets API")
	}

	readRange := tenSheet + "!A:AZ"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Printf("⚠️ Lỗi đọc sheet %s (ID: %s): %v", tenSheet, spreadsheetID[:5], err)
		return nil, err
	}
	return resp.Values, nil
}

// --- CÁC HÀM PARSE DỮ LIỆU ---

func LayString(row []interface{}, index int) string {
	if index >= len(row) || row[index] == nil { return "" }
	return strings.TrimSpace(fmt.Sprintf("%v", row[index]))
}

func LayInt(row []interface{}, index int) int {
	s := LayString(row, index)
	if s == "" { return 0 }
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, " ", "")
	val, _ := strconv.Atoi(s)
	return val
}

func LayFloat(row []interface{}, index int) float64 {
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

func LayChuoiSoNgauNhien(doDai int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	b := make([]byte, doDai)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
