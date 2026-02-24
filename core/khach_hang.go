package core

import (
	"encoding/json"
	"strings"
	"time"
	"app/cau_hinh" 
)

// =============================================================
// 1. ĐỊNH NGHĨA VỊ TRÍ CỘT (27 Cột từ A -> AA)
// =============================================================
const (
	DongBatDau_KhachHang = 11

	CotKH_MaKhachHang        = 0  // A
	CotKH_TenDangNhap        = 1  // B (Đóng vai trò là Subdomain)
	CotKH_Email              = 2  // C
	CotKH_MatKhauHash        = 3  // D
	CotKH_MaPinHash          = 4  // E
	CotKH_RefreshTokenJson   = 5  // F
	CotKH_VaiTroQuyenHan     = 6  // G
	CotKH_ChucVu             = 7  // H
	CotKH_TrangThai          = 8  // I
	
	// [CỤM LÕI SAAS]
	CotKH_DataSheetsJson     = 9  // J
	CotKH_GoiDichVuJson      = 10 // K
	CotKH_CauHinhJson        = 11 // L
	
	// [CỤM THÔNG TIN LIÊN HỆ / HỒ SƠ]
	CotKH_NguonKhachHang     = 12 // M
	CotKH_TenKhachHang       = 13 // N
	CotKH_DienThoai          = 14 // O
	CotKH_AnhDaiDien         = 15 // P
	CotKH_MangXaHoiJson      = 16 // Q
	CotKH_DiaChi             = 17 // R
	CotKH_NgaySinh           = 18 // S
	CotKH_GioiTinh           = 19 // T
	CotKH_MaSoThue           = 20 // U
	CotKH_ViTienJson         = 21 // V
	CotKH_InboxJson          = 22 // W
	
	// [LƯU VẾT]
	CotKH_GhiChu             = 23 // X
	CotKH_NgayTao            = 24 // Y
	CotKH_NguoiCapNhat       = 25 // Z (Mới sửa)
	CotKH_NgayCapNhat        = 26 // AA
)

// =============================================================
// 2. CÁC STRUCT THÀNH PHẦN (JSON)
// =============================================================

// Cột F
type TokenInfo struct { 
	DeviceName string `json:"dev"`
	ExpiresAt  int64  `json:"exp"` 
}

// Cột J: Cấu trúc kết nối DB của riêng Shop
type DataSheetInfo struct { 
	SpreadsheetID  string `json:"sheet_id"`
	GoogleAuthJson string `json:"google_auth_json"`
	FolderDriveID  string `json:"folder_drive_id"` 
}

// Định nghĩa cấu trúc Gói Dịch Vụ
type PlanInfo struct {
	MaGoi      string `json:"ma_goi"`       // Mã gói: TRIAL_3DAYS, BASIC_1M, VIP_12M...
	TenGoi     string `json:"ten_goi"`      // Tên hiển thị: Dùng thử 3 ngày, Gói 1 tháng...
	NgayHetHan string `json:"ngay_het_han"` // Thời gian hết hạn
	TrangThai  string `json:"trang_thai"`   // Trạng thái: "active", "expired"
}

// Cột L: Cấu hình giao diện shop / người dùng
type UserConfig struct { 
	Theme        string `json:"theme"`
	ChuyenNganh  string `json:"chuyen_nganh"`
	CustomDomain string `json:"custom_domain"`
	DarkMode     bool   `json:"dark_mode"`
	Language     string `json:"lang"` 
}

// Các cột Q, V, W
type SocialInfo struct { Zalo string `json:"zalo"`; Facebook string `json:"fb"`; Tiktok string `json:"tiktok"` }
type WalletInfo struct { SoDuHienTai float64 `json:"so_du"` }
type MessageInfo struct { 
	ID             string `json:"id"`
	TieuDe         string `json:"title"`
	NoiDung        string `json:"content"`
	DaDoc          bool   `json:"is_read"`
	NgayTao        string `json:"date"`
	NguoiGuiID     string `json:"sender_id"`
	NguoiGuiTen    string `json:"sender_name"`
	NguoiGuiChucVu string `json:"sender_title"`
}
// =============================================================
// 3. STRUCT CHÍNH (ĐỐI TƯỢNG KHÁCH HÀNG / CHỦ SHOP)
// =============================================================
type KhachHang struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaKhachHang      string `json:"ma_khach_hang"`
	TenDangNhap      string `json:"ten_dang_nhap"` // == Subdomain
	Email            string `json:"email"`
	MatKhauHash      string `json:"-"`
	MaPinHash        string `json:"-"`
	RefreshTokens    map[string]TokenInfo `json:"-"`

	VaiTroQuyenHan   string `json:"vai_tro_quyen_han"`
	ChucVu           string `json:"chuc_vu"`
	TrangThai        int    `json:"trang_thai"`

	// --- CỤM LÕI SAAS ---
	DataSheets       DataSheetInfo `json:"data_sheets"`
	GoiDichVu        []PlanInfo    `json:"goi_dich_vu"` // Đã chuyển thành Mảng
	CauHinh          UserConfig    `json:"cau_hinh"`    // Kéo lên trên theo thiết kế

	// --- CỤM HỒ SƠ ---
	NguonKhachHang   string     `json:"nguon_khach_hang"`
	TenKhachHang     string     `json:"ten_khach_hang"`
	DienThoai        string     `json:"dien_thoai"`
	AnhDaiDien       string     `json:"anh_dai_dien"`
	MangXaHoi        SocialInfo `json:"mang_xa_hoi"`
	DiaChi           string     `json:"dia_chi"`
	NgaySinh         string     `json:"ngay_sinh"`
	GioiTinh         int        `json:"gioi_tinh"`
	MaSoThue         string     `json:"ma_so_thue"`
	ViTien           WalletInfo `json:"vi_tien"`
	Inbox            []MessageInfo `json:"inbox"`

	// --- LƯU VẾT ---
	GhiChu           string     `json:"ghi_chu"`
	NgayTao          string     `json:"ngay_tao"`
	NguoiCapNhat     string     `json:"nguoi_cap_nhat"`
	NgayCapNhat      string     `json:"ngay_cap_nhat"`
}

// BỘ NHỚ ĐA SHOP
var (
	CacheKhachHang    = make(map[string][]*KhachHang)
	CacheMapKhachHang = make(map[string]*KhachHang) // Map Phẳng
)

// =============================================================
// 4. LOGIC ĐỌC VÀ LƯU VÀO RAM
// =============================================================
func NapKhachHang(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }

	raw, err := LoadSheetData(shopID, "KHACH_HANG")
	if err != nil { return }

	list := []*KhachHang{}

	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := LayString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }

		kh := &KhachHang{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaKhachHang:    maKH,
			TenDangNhap:    LayString(r, CotKH_TenDangNhap),
			Email:          LayString(r, CotKH_Email),
			MatKhauHash:    LayString(r, CotKH_MatKhauHash),
			MaPinHash:      LayString(r, CotKH_MaPinHash),
			VaiTroQuyenHan: strings.TrimSpace(LayString(r, CotKH_VaiTroQuyenHan)),
			ChucVu:         strings.TrimSpace(LayString(r, CotKH_ChucVu)),
			TrangThai:      LayInt(r, CotKH_TrangThai),
			NguonKhachHang: LayString(r, CotKH_NguonKhachHang),
			TenKhachHang:   LayString(r, CotKH_TenKhachHang),
			DienThoai:      LayString(r, CotKH_DienThoai),
			AnhDaiDien:     LayString(r, CotKH_AnhDaiDien),
			DiaChi:         LayString(r, CotKH_DiaChi),
			NgaySinh:       LayString(r, CotKH_NgaySinh),
			GioiTinh:       LayInt(r, CotKH_GioiTinh),
			MaSoThue:       LayString(r, CotKH_MaSoThue),
			GhiChu:         LayString(r, CotKH_GhiChu),
			NgayTao:        LayString(r, CotKH_NgayTao),
			NguoiCapNhat:   LayString(r, CotKH_NguoiCapNhat),
			NgayCapNhat:    LayString(r, CotKH_NgayCapNhat),
		}

		// --- PARSE JSON (KÈM BẢO VỆ CHỐNG LỖI) ---
		
		// Token Map
		_ = json.Unmarshal([]byte(LayString(r, CotKH_RefreshTokenJson)), &kh.RefreshTokens)
		if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]TokenInfo) }
		
		// Lõi SaaS (Cột J, K, L)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_DataSheetsJson)), &kh.DataSheets)
		// Kích hoạt API riêng nếu chủ shop có cung cấp
        if kh.DataSheets.GoogleAuthJson != "" && kh.DataSheets.SpreadsheetID != "" {
        KetNoiGoogleSheetRieng(kh.DataSheets.SpreadsheetID, kh.DataSheets.GoogleAuthJson)
        }
		_ = json.Unmarshal([]byte(LayString(r, CotKH_GoiDichVuJson)), &kh.GoiDichVu)
		if kh.GoiDichVu == nil { kh.GoiDichVu = make([]PlanInfo, 0) } // Ép mảng rỗng nếu chưa có
		_ = json.Unmarshal([]byte(LayString(r, CotKH_CauHinhJson)), &kh.CauHinh)

		// Thông tin râu ria
		_ = json.Unmarshal([]byte(LayString(r, CotKH_MangXaHoiJson)), &kh.MangXaHoi)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_ViTienJson)), &kh.ViTien)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_InboxJson)), &kh.Inbox)
		if kh.Inbox == nil { kh.Inbox = make([]MessageInfo, 0) }

		list = append(list, kh)

		// Lưu vào Map Phẳng tra cứu nhanh
		key := TaoCompositeKey(shopID, maKH)
		CacheMapKhachHang[key] = kh
	}

	KhoaHeThong.Lock()
	CacheKhachHang[shopID] = list 
	KhoaHeThong.Unlock()
}

// =============================================================
// 5. CÁC HÀM HELPER ĐỂ CONTROLLER GỌI
// =============================================================

func LayDanhSachKhachHang(shopID string) []*KhachHang {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheKhachHang[shopID]
}

func TimKhachHangTheoCookie(shopID, cookie string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	list := CacheKhachHang[shopID]
	for _, kh := range list {
		if info, ok := kh.RefreshTokens[cookie]; ok {
			if time.Now().Unix() > info.ExpiresAt { return nil, false }
			return kh, true
		}
	}
	return nil, false
}

// Do Cột B đóng vai trò là Subdomain, hàm này vừa tìm Tên đăng nhập vừa là tìm Cửa hàng
func TimKhachHangTheoUserOrEmail(shopID, input string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	input = strings.ToLower(strings.TrimSpace(input))
	list := CacheKhachHang[shopID]
	
	for _, kh := range list {
		if strings.ToLower(kh.TenDangNhap) == input { return kh, true }
		if kh.Email != "" && strings.ToLower(kh.Email) == input { return kh, true }
	}
	return nil, false
}

func LayKhachHang(shopID, maKH string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(shopID, maKH)
	kh, ok := CacheMapKhachHang[key]
	return kh, ok
}

func TaoMaKhachHangMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	for {
		id := LayChuoiSoNgauNhien(19)
		key := TaoCompositeKey(shopID, id)
		if _, exist := CacheMapKhachHang[key]; !exist { return id }
	}
}

func ToJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := kh.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet } 
	
	CacheKhachHang[sID] = append(CacheKhachHang[sID], kh)
	
	key := TaoCompositeKey(sID, kh.MaKhachHang)
	CacheMapKhachHang[key] = kh
}
