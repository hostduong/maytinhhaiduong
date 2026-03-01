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

	"app/config" // ĐÃ CHUẨN HÓA

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// =============================================================
// 1. CẤU TRÚC HẠ TẦNG (INFRASTRUCTURE)
// =============================================================

var (
	KhoaHeThong sync.RWMutex
	HeThongDangBan bool

	MapDichVuSheet = make(map[string]*sheets.Service)
	MutexDichVu    sync.RWMutex 
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
	jsonKey := config.BienCauHinh.GoogleAuthJson // ĐÃ SỬA

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

	MutexDichVu.Lock()
	MapDichVuSheet["default"] = srv
	MapDichVuSheet[config.BienCauHinh.IdFileSheet] = srv // ĐÃ SỬA
	MutexDichVu.Unlock()
	
	log.Println("✅ [CORE] Khởi tạo API mặc định thành công!")
}

// =============================================================
// 3. QUẢN LÝ POOL KẾT NỐI (MULTITENANT API)
// =============================================================

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

func LayDichVuSheet(shopID string) *sheets.Service {
	MutexDichVu.RLock()
	srv, ok := MapDichVuSheet[shopID]
	MutexDichVu.RUnlock()

	if ok && srv != nil { return srv }

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

func LoadSheetData(spreadsheetID string, tenSheet string) ([][]interface{}, error) {
	if spreadsheetID == "" {
		spreadsheetID = config.BienCauHinh.IdFileSheet // ĐÃ SỬA
	}

	srv := LayDichVuSheet(spreadsheetID)
	if srv == nil { return nil, fmt.Errorf("chưa kết nối được Google Sheets API") }

	readRange := tenSheet + "!A:AZ"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Printf("⚠️ Lỗi đọc sheet %s (ID: %s): %v", tenSheet, spreadsheetID[:5], err)
		return nil, err
	}
	return resp.Values, nil
}

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
	s := LayString(row, index)
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
	for i := range b { b[i] = digits[rand.Intn(len(digits))] }
	return string(b)
}

func KiemTraVaKhoiTaoSheetNganh(shopID, spreadsheetID, authJson, chuyenNganh string) error {
	if authJson != "" && spreadsheetID != "" { KetNoiGoogleSheetRieng(shopID, authJson) }
	srv := LayDichVuSheet(shopID)
	if srv == nil { return fmt.Errorf("Không thể khởi tạo kết nối Google API.") }

	resp, err := srv.Spreadsheets.Get(spreadsheetID).Fields("sheets(properties(title))").Do()
	if err != nil { return fmt.Errorf("Không thể truy cập Spreadsheet.") }

	tenTabCanTao := strings.ToUpper(chuyenNganh)
	tabDaTonTai := false
	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == tenTabCanTao { tabDaTonTai = true; break }
	}

	if !tabDaTonTai {
		req := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{{AddSheet: &sheets.AddSheetRequest{Properties: &sheets.SheetProperties{Title: tenTabCanTao}}}},
		}
		_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, req).Do()
		if err != nil { return fmt.Errorf("Lỗi tạo Tab: %v", err) }
	}
	return nil
}

func KiemTraFolderDrive(folderID string, jsonKey string) error {
	if folderID == "" { return nil }

	ctx := context.Background()
	var srv *drive.Service
	var err error

	if jsonKey != "" {
		srv, err = drive.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	} else if config.BienCauHinh.GoogleAuthJson != "" { // ĐÃ SỬA
		srv, err = drive.NewService(ctx, option.WithCredentialsJSON([]byte(config.BienCauHinh.GoogleAuthJson))) // ĐÃ SỬA
	} else {
		srv, err = drive.NewService(ctx, option.WithScopes(drive.DriveReadonlyScope))
	}

	if err != nil { return fmt.Errorf("Lỗi cấu hình Google API.") }

	f, err := srv.Files.Get(folderID).Fields("id, mimeType").Do()
	if err != nil { return fmt.Errorf("Không thể truy cập Thư mục Drive.") }

	if f.MimeType != "application/vnd.google-apps.folder" { return fmt.Errorf("ID không phải là Folder.") }
	return nil
}
