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

	"app/config"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/drive/v3"
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
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

// KiemTraVaKhoiTaoSheetNganh: Kiểm tra quyền truy cập và tự động tạo Tab theo chuyên ngành
func KiemTraVaKhoiTaoSheetNganh(shopID, spreadsheetID, authJson, chuyenNganh string) error {
	// 1. Nạp API Riêng nếu có (Để đảm bảo dùng đúng thông tin vừa nhập)
	if authJson != "" && spreadsheetID != "" {
		KetNoiGoogleSheetRieng(shopID, authJson)
	}

	srv := LayDichVuSheet(shopID)
	if srv == nil {
		return fmt.Errorf("Không thể khởi tạo kết nối Google API. Vui lòng kiểm tra lại JSON Auth.")
	}

	// 2. Chọc thử vào Google Sheet để lấy MetaData (Kiểm tra quyền)
	resp, err := srv.Spreadsheets.Get(spreadsheetID).Fields("sheets(properties(title))").Do()
	if err != nil {
		return fmt.Errorf("Không thể truy cập Spreadsheet. Sai ID hoặc chưa cấp quyền Editor cho www.99k.vn@gmail.com.")
	}

	// 3. Quy chuẩn tên Sheet (VD: may_tinh -> MAY_TINH)
	tenTabCanTao := strings.ToUpper(chuyenNganh)
	tabDaTonTai := false

	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == tenTabCanTao {
			tabDaTonTai = true
			break
		}
	}

	// 4. Nếu chưa có Tab -> Ra lệnh tạo ngay lập tức (Sync)
	if !tabDaTonTai {
		req := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					AddSheet: &sheets.AddSheetRequest{
						Properties: &sheets.SheetProperties{
							Title: tenTabCanTao,
						},
					},
				},
			},
		}

		_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, req).Do()
		if err != nil {
			return fmt.Errorf("Đã kết nối thành công nhưng hệ thống không thể tự tạo Tab '%s'. Lỗi: %v", tenTabCanTao, err)
		}
		log.Printf("✨ [AUTO-SETUP] Đã tự động tạo Tab '%s' cho Shop [%s]", tenTabCanTao, shopID)
	}

	return nil // Mọi thứ ĐỀU XANH!
}

// KiemTraFolderDrive: Kiểm tra quyền truy cập và tính hợp lệ của Folder ID
func KiemTraFolderDrive(folderID string, jsonKey string) error {
	if folderID == "" {
		return nil // Không nhập thì không kiểm tra
	}

	ctx := context.Background()
	var srv *drive.Service
	var err error

	// Khởi tạo kết nối Drive
	if jsonKey != "" {
		srv, err = drive.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	} else if cau_hinh.BienCauHinh.GoogleAuthJson != "" {
		srv, err = drive.NewService(ctx, option.WithCredentialsJSON([]byte(cau_hinh.BienCauHinh.GoogleAuthJson)))
	} else {
		srv, err = drive.NewService(ctx, option.WithScopes(drive.DriveReadonlyScope))
	}

	if err != nil {
		return fmt.Errorf("Lỗi cấu hình Google API, không thể kiểm tra Drive.")
	}

	// Chọc thử vào Google Drive để lấy thông tin
	f, err := srv.Files.Get(folderID).Fields("id, mimeType").Do()
	if err != nil {
		return fmt.Errorf("Không thể truy cập Thư mục Drive. Vui lòng kiểm tra lại ID hoặc đảm bảo đã Share quyền Editor cho www.99k.vn@gmail.com.")
	}

	// Đảm bảo ID cung cấp là một Thư mục chứ không phải ID của một File ảnh/File doc
	if f.MimeType != "application/vnd.google-apps.folder" {
		return fmt.Errorf("ID bạn nhập không phải là một Thư mục (Folder). Vui lòng copy đúng ID của Thư mục gốc.")
	}

	return nil // Xanh mượt!
}
