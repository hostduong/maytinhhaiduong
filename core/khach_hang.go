package core

import (
	"encoding/json"
	"strings"
	"time"
	"app/cau_hinh" 
)

// =============================================================
// 1. ĐỊNH NGHĨA VỊ TRÍ CỘT (ĐÃ DỊCH CHUYỂN BỎ INBOX)
// =============================================================
const (
	DongBatDau_KhachHang = 11

	CotKH_MaKhachHang        = 0  // A
	CotKH_TenDangNhap        = 1  // B (Subdomain)
	CotKH_Email              = 2  // C
	CotKH_MatKhauHash        = 3  // D
	CotKH_MaPinHash          = 4  // E
	CotKH_RefreshTokenJson   = 5  // F
	CotKH_VaiTroQuyenHan     = 6  // G
	CotKH_ChucVu             = 7  // H
	CotKH_TrangThai          = 8  // I
	
	CotKH_DataSheetsJson     = 9  // J
	CotKH_GoiDichVuJson      = 10 // K
	CotKH_CauHinhJson        = 11 // L
	
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
	
	// [ĐÃ XÓA CotKH_InboxJson Ở ĐÂY - CHUYỂN SANG SHEET MỚI]
	
	CotKH_GhiChu             = 22 // W (Đã kéo lên cột 22)
	CotKH_NgayTao            = 23 // X
	CotKH_NguoiCapNhat       = 24 // Y 
	CotKH_NgayCapNhat        = 25 // Z
)

// =============================================================
// 2. CÁC STRUCT THÀNH PHẦN (JSON)
// =============================================================
type TokenInfo struct { 
	DeviceName string `json:"dev"`
	ExpiresAt  int64  `json:"exp"` 
}
type DataSheetInfo struct { 
	SpreadsheetID  string `json:"sheet_id"`
	GoogleAuthJson string `json:"google_auth_json"`
	FolderDriveID  string `json:"folder_drive_id"` 
}
type PlanInfo struct {
	MaGoi      string `json:"ma_goi"`       
	TenGoi     string `json:"ten_goi"`      
	NgayHetHan string `json:"ngay_het_han"` 
	TrangThai  string `json:"trang_thai"`   
}
type UserConfig struct { 
	Theme        string `json:"theme"`
	ChuyenNganh  string `json:"chuyen_nganh"`
	CustomDomain string `json:"custom_domain"`
	DarkMode     bool   `json:"dark_mode"`
	Language     string `json:"lang"` 
}
type SocialInfo struct { Zalo string `json:"zalo"`; Facebook string `json:"fb"`; Tiktok string `json:"tiktok"` }
type WalletInfo struct { SoDuHienTai float64 `json:"so_du"` }

// =============================================================
// 3. STRUCT CHÍNH (ĐỐI TƯỢNG KHÁCH HÀNG / CHỦ SHOP)
// =============================================================
type KhachHang struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaKhachHang      string `json:"ma_khach_hang"`
	TenDangNhap      string `json:"ten_dang_nhap"`
	Email            string `json:"email"`
	MatKhauHash      string `json:"-"`
	MaPinHash        string `json:"-"`
	RefreshTokens    map[string]TokenInfo `json:"-"`

	VaiTroQuyenHan   string `json:"vai_tro_quyen_han"`
	ChucVu           string `json:"chuc_vu"`
	TrangThai        int    `json:"trang_thai"`

	DataSheets       DataSheetInfo `json:"data_sheets"`
	GoiDichVu        []PlanInfo    `json:"goi_dich_vu"` 
	CauHinh          UserConfig    `json:"cau_hinh"`    

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
	
	// [MỚI] Trỏ thẳng sang Lõi Tin nhắn mới. Đánh dấu json:"-" để KHÔNG ghi xuống sheet Khách Hàng.
	Inbox            []*TinNhan `json:"-"` 

	GhiChu           string     `json:"ghi_chu"`
	NgayTao          string     `json:"ngay_tao"`
	NguoiCapNhat     string     `json:"nguoi_cap_nhat"`
	NgayCapNhat      string     `json:"ngay_cap_nhat"`
}

var (
	CacheKhachHang    = make(map[string][]*KhachHang)
	CacheMapKhachHang = make(map[string]*KhachHang) 
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
			Inbox:          make([]*TinNhan, 0), // Tránh lỗi Nil Pointer
		}

		_ = json.Unmarshal([]byte(LayString(r, CotKH_RefreshTokenJson)), &kh.RefreshTokens)
		if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]TokenInfo) }
		
		_ = json.Unmarshal([]byte(LayString(r, CotKH_DataSheetsJson)), &kh.DataSheets)
		if kh.DataSheets.GoogleAuthJson != "" && kh.DataSheets.SpreadsheetID != "" {
			KetNoiGoogleSheetRieng(kh.DataSheets.SpreadsheetID, kh.DataSheets.GoogleAuthJson)
		}
		_ = json.Unmarshal([]byte(LayString(r, CotKH_GoiDichVuJson)), &kh.GoiDichVu)
		if kh.GoiDichVu == nil { kh.GoiDichVu = make([]PlanInfo, 0) } 
		
		_ = json.Unmarshal([]byte(LayString(r, CotKH_CauHinhJson)), &kh.CauHinh)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_MangXaHoiJson)), &kh.MangXaHoi)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_ViTienJson)), &kh.ViTien)

		list = append(list, kh)
		key := TaoCompositeKey(shopID, maKH)
		CacheMapKhachHang[key] = kh
	}

	KhoaHeThong.Lock()
	CacheKhachHang[shopID] = list 
	KhoaHeThong.Unlock()
}

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
	for { id := LayChuoiSoNgauNhien(19); key := TaoCompositeKey(shopID, id); if _, exist := CacheMapKhachHang[key]; !exist { return id } }
}
func ToJSON(v interface{}) string { b, _ := json.Marshal(v); return string(b) }
func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	sID := kh.SpreadsheetID; if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet } 
	CacheKhachHang[sID] = append(CacheKhachHang[sID], kh); key := TaoCompositeKey(sID, kh.MaKhachHang); CacheMapKhachHang[key] = kh
}
