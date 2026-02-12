package core

import (
	"encoding/json"
	"strings"
	"time"

	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT (MAPPING A -> Z)
// =============================================================
const (
	DongBatDau_KhachHang = 2 // Dòng 1 là tiêu đề

	CotKH_MaKhachHang        = 0  // A
	CotKH_TenDangNhap        = 1  // B
	CotKH_Email              = 2  // C
	CotKH_MatKhauHash        = 3  // D
	CotKH_MaPinHash          = 4  // E
	CotKH_RefreshTokenJson   = 5  // F [JSON]
	CotKH_VaiTroQuyenHan     = 6  // G
	CotKH_ChucVu             = 7  // H
	CotKH_TrangThai          = 8  // I
	
	CotKH_DataSheetsJson     = 9  // J [JSON]
	CotKH_GoiDichVuJson      = 10 // K [JSON]
	
	CotKH_NguonKhachHang     = 11 // L
	CotKH_TenKhachHang       = 12 // M
	CotKH_DienThoai          = 13 // N
	CotKH_AnhDaiDien         = 14 // O
	CotKH_MangXaHoiJson      = 15 // P [JSON]
	CotKH_DiaChi             = 16 // Q
	CotKH_NgaySinh           = 17 // R
	CotKH_GioiTinh           = 18 // S
	CotKH_MaSoThue           = 19 // T
	
	CotKH_ViTienJson         = 20 // U [JSON] (Gộp tiền & nợ & điểm)
	CotKH_CauHinhJson        = 21 // V [JSON]
	CotKH_InboxJson          = 22 // W [JSON] (Tin nhắn Admin)
	
	CotKH_GhiChu             = 23 // X
	CotKH_NgayTao            = 24 // Y
	CotKH_NgayCapNhat        = 25 // Z
)

// =============================================================
// 2. STRUCT GOLANG (Object trong RAM)
// =============================================================

type KhachHang struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	// --- ĐỊNH DANH & BẢO MẬT ---
	MaKhachHang      string `json:"ma_khach_hang"`
	TenDangNhap      string `json:"ten_dang_nhap"`
	Email            string `json:"email"`
	MatKhauHash      string `json:"-"` 
	MaPinHash        string `json:"-"` 
	RefreshTokens    map[string]TokenInfo `json:"-"` // Cột F (Map để lưu nhiều thiết bị)

	// --- PHÂN QUYỀN ---
	VaiTroQuyenHan   string `json:"vai_tro_quyen_han"`
	ChucVu           string `json:"chuc_vu"`
	TrangThai        int    `json:"trang_thai"` 

	// --- SAAS ---
	DataSheets       DataSheetInfo `json:"data_sheets"` // Cột J
	GoiDichVu        PlanInfo      `json:"goi_dich_vu"` // Cột K

	// --- THÔNG TIN CÁ NHÂN ---
	NguonKhachHang   string     `json:"nguon_khach_hang"`
	TenKhachHang     string     `json:"ten_khach_hang"`
	DienThoai        string     `json:"dien_thoai"`
	AnhDaiDien       string     `json:"anh_dai_dien"`
	MangXaHoi        SocialInfo `json:"mang_xa_hoi"` // Cột P
	DiaChi           string     `json:"dia_chi"`
	NgaySinh         string     `json:"ngay_sinh"`
	GioiTinh         int        `json:"gioi_tinh"` 
	MaSoThue         string     `json:"ma_so_thue"`

	// --- TÀI CHÍNH & TƯƠNG TÁC (NEW) ---
	ViTien           WalletInfo    `json:"vi_tien"`  // Cột U (Tổng hợp tiền)
	CauHinh          UserConfig    `json:"cau_hinh"` // Cột V
	Inbox            []MessageInfo `json:"inbox"`    // Cột W (List tin nhắn)

	// --- META ---
	GhiChu           string     `json:"ghi_chu"`
	NgayTao          string     `json:"ngay_tao"`
	NgayCapNhat      string     `json:"ngay_cap_nhat"`
}

// =============================================================
// 3. CÁC STRUCT CON (Để Parse JSON)
// =============================================================

// Cột F: Token
type TokenInfo struct {
	HashedToken string `json:"t"`   
	ExpiresAt   int64  `json:"exp"` 
	DeviceName  string `json:"dev"` 
	CreatedAt   string `json:"at"`  
}

// Cột J: Data Sheet
type DataSheetInfo struct {
	SpreadsheetID string `json:"sheet_id"` 
	DriveFolderID string `json:"drive_id"` 
	CustomAuth    string `json:"auth_json,omitempty"` // Nếu shop dùng Auth riêng
}

// Cột K: Gói cước
type PlanInfo struct {
	PlanName   string `json:"name"`  
	ExpiredAt  string `json:"exp"`    
	MaxProduct int    `json:"limit"` 
}

// Cột P: Mạng xã hội
type SocialInfo struct {
	Zalo     string `json:"zalo"`
	Facebook string `json:"fb"`
	Telegram string `json:"tele"`
	Tiktok   string `json:"tiktok"`
}

// Cột U: Ví tiền & Điểm (GOM NHÓM)
type WalletInfo struct {
	SoDuHienTai  float64 `json:"so_du"`    // Tiền mặt đang có
	CongNo       float64 `json:"cong_no"`  // Đang nợ shop
	TongTienMua  float64 `json:"tong_mua"` // Tổng tiền tích lũy
	DiemThuong   int     `json:"diem"`     // Điểm Loyalty
}

// Cột V: Cấu hình
type UserConfig struct {
	Theme    string `json:"theme"`      
	Language string `json:"lang"`       
}

// Cột W: Inbox/Thông báo (MỚI)
type MessageInfo struct {
	ID       string `json:"id"`
	TieuDe   string `json:"title"`
	NoiDung  string `json:"body"`
	Link     string `json:"link,omitempty"` // Bấm vào nhảy đi đâu
	DaXem    int    `json:"read"`           // 1: Rồi, 0: Chưa
	ThoiGian string `json:"time"`
}

// ... (Giữ nguyên các biến _DS_KhachHang, _Map_KhachHang) ...

// =============================================================
// 4. LOGIC NẠP (CẬP NHẬT PARSE JSON)
// =============================================================
func NapKhachHang(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" { targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(targetSpreadsheetID, "KHACH_HANG")
	if err != nil { return }

	_Map_KhachHang = make(map[string]*KhachHang)
	_DS_KhachHang = []*KhachHang{}

	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := layString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }

		key := TaoCompositeKey(targetSpreadsheetID, maKH)
		if _, daTonTai := _Map_KhachHang[key]; daTonTai { continue }

		kh := &KhachHang{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaKhachHang:    maKH,
			TenDangNhap:    layString(r, CotKH_TenDangNhap),
			Email:          layString(r, CotKH_Email),
			MatKhauHash:    layString(r, CotKH_MatKhauHash),
			MaPinHash:      layString(r, CotKH_MaPinHash),
			VaiTroQuyenHan: layString(r, CotKH_VaiTroQuyenHan),
			ChucVu:         layString(r, CotKH_ChucVu),
			TrangThai:      layInt(r, CotKH_TrangThai),
			NguonKhachHang: layString(r, CotKH_NguonKhachHang),
			TenKhachHang:   layString(r, CotKH_TenKhachHang),
			DienThoai:      layString(r, CotKH_DienThoai),
			AnhDaiDien:     layString(r, CotKH_AnhDaiDien),
			DiaChi:         layString(r, CotKH_DiaChi),
			NgaySinh:       layString(r, CotKH_NgaySinh),
			GioiTinh:       layInt(r, CotKH_GioiTinh),
			MaSoThue:       layString(r, CotKH_MaSoThue),
			GhiChu:         layString(r, CotKH_GhiChu),
			NgayTao:        layString(r, CotKH_NgayTao),
			NgayCapNhat:    layString(r, CotKH_NgayCapNhat),
		}

		// --- PARSE CÁC CỘT JSON ---
		parseJSON(layString(r, CotKH_RefreshTokenJson), &kh.RefreshTokens)
		parseJSON(layString(r, CotKH_DataSheetsJson), &kh.DataSheets)
		parseJSON(layString(r, CotKH_GoiDichVuJson), &kh.GoiDichVu)
		parseJSON(layString(r, CotKH_MangXaHoiJson), &kh.MangXaHoi)
		parseJSON(layString(r, CotKH_ViTienJson), &kh.ViTien)
		parseJSON(layString(r, CotKH_CauHinhJson), &kh.CauHinh)
		parseJSON(layString(r, CotKH_InboxJson), &kh.Inbox)

		_DS_KhachHang = append(_DS_KhachHang, kh)
		_Map_KhachHang[key] = kh
	}
}

// Helper nhỏ để parse
func parseJSON(jsonStr string, target interface{}) {
	if jsonStr != "" && jsonStr != "{}" && jsonStr != "[]" {
		_ = json.Unmarshal([]byte(jsonStr), target)
	}
}

// Helper: Convert Struct -> JSON String (Dùng khi ghi Sheet)
func ToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil { return "" }
	return string(b)
}

func LayDanhSachKhachHang() []*KhachHang {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	kq := make([]*KhachHang, len(_DS_KhachHang))
	copy(kq, _DS_KhachHang)
	return kq
}

func LayKhachHang(maKH string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	sheetID := cau_hinh.BienCauHinh.IdFileSheet
	key := TaoCompositeKey(sheetID, maKH)
	kh, ok := _Map_KhachHang[key]
	return kh, ok
}

func TimKhachHangTheoCookie(cookie string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	for _, kh := range _DS_KhachHang {
		if kh.Cookie == cookie && cookie != "" {
			if time.Now().Unix() > kh.CookieExpired { return nil, false }
			return kh, true
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(input string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range _DS_KhachHang {
		if strings.ToLower(kh.TenDangNhap) == input { return kh, true }
		if kh.Email != "" && strings.ToLower(kh.Email) == input { return kh, true }
	}
	return nil, false
}

func KiemTraTonTaiUserEmail(user, email string) bool {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	user = strings.ToLower(strings.TrimSpace(user))
	email = strings.ToLower(strings.TrimSpace(email))
	for _, kh := range _DS_KhachHang {
		if strings.ToLower(kh.TenDangNhap) == user { return true }
		if email != "" && strings.ToLower(kh.Email) == email { return true }
	}
	return false
}

func TaoMaKhachHangMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	for {
		id := LayChuoiSoNgauNhien(19)
		key := TaoCompositeKey(currentSheetID, id)
		if _, tonTai := _Map_KhachHang[key]; !tonTai { return id }
	}
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if kh.SpreadsheetID == "" { kh.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_KhachHang = append(_DS_KhachHang, kh)
	key := TaoCompositeKey(kh.SpreadsheetID, kh.MaKhachHang)
	_Map_KhachHang[key] = kh
}
