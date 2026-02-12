package core

import (
	"encoding/json"
	"strings"
	"time"
	"app/cau_hinh" // Import thêm để lấy ID mặc định
)

// ... (Giữ nguyên phần Const và Struct KhachHang như cũ) ...
// (Lưu ý: Giữ nguyên Struct KhachHang và các Struct con TokenInfo, UserConfig... đã sửa ở bước trước)

const (
	DongBatDau_KhachHang = 11

	CotKH_MaKhachHang        = 0
	CotKH_TenDangNhap        = 1
	CotKH_Email              = 2
	CotKH_MatKhauHash        = 3
	CotKH_MaPinHash          = 4
	CotKH_RefreshTokenJson   = 5
	CotKH_VaiTroQuyenHan     = 6
	CotKH_ChucVu             = 7
	CotKH_TrangThai          = 8
	CotKH_DataSheetsJson     = 9
	CotKH_GoiDichVuJson      = 10
	CotKH_NguonKhachHang     = 11
	CotKH_TenKhachHang       = 12
	CotKH_DienThoai          = 13
	CotKH_AnhDaiDien         = 14
	CotKH_MangXaHoiJson      = 15
	CotKH_DiaChi             = 16
	CotKH_NgaySinh           = 17
	CotKH_GioiTinh           = 18
	CotKH_MaSoThue           = 19
	CotKH_ViTienJson         = 20
	CotKH_CauHinhJson        = 21
	CotKH_InboxJson          = 22
	CotKH_GhiChu             = 23
	CotKH_NgayTao            = 24
	CotKH_NgayCapNhat        = 25
)

// STRUCT CHÍNH
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
	GoiDichVu        PlanInfo      `json:"goi_dich_vu"`

	NguonKhachHang   string     `json:"nguon_khach_hang"`
	TenKhachHang     string     `json:"ten_khach_hang"`
	DienThoai        string     `json:"dien_thoai"`
	AnhDaiDien       string     `json:"anh_dai_dien"`
	MangXaHoi        SocialInfo `json:"mang_xa_hoi"`
	DiaChi           string     `json:"dia_chi"`
	NgaySinh         string     `json:"ngay_sinh"`
	GioiTinh         int        `json:"gioi_tinh"`
	MaSoThue         string     `json:"ma_so_thue"`

	ViTien           WalletInfo    `json:"vi_tien"`
	CauHinh          UserConfig    `json:"cau_hinh"`
	Inbox            []MessageInfo `json:"inbox"`

	GhiChu           string     `json:"ghi_chu"`
	NgayTao          string     `json:"ngay_tao"`
	NgayCapNhat      string     `json:"ngay_cap_nhat"`
}

type TokenInfo struct { DeviceName string `json:"dev"`; ExpiresAt int64 `json:"exp"` }
type DataSheetInfo struct { SpreadsheetID string `json:"sheet_id"` }
type PlanInfo struct { PlanName string `json:"name"` }
type SocialInfo struct { Zalo string `json:"zalo"`; Facebook string `json:"fb"`; Tiktok string `json:"tiktok"` }
type WalletInfo struct { SoDuHienTai float64 `json:"so_du"` }
type UserConfig struct { Theme string `json:"theme"`; Language string `json:"lang"` }
type MessageInfo struct { TieuDe string `json:"title"` }

// BỘ NHỚ ĐA SHOP
var (
	CacheKhachHang    = make(map[string][]*KhachHang)
	CacheMapKhachHang = make(map[string]*KhachHang) // Map Phẳng
)

func NapKhachHang(shopID string) {
	// [QUAN TRỌNG] Chuẩn hóa ID Shop trước khi dùng
	if shopID == "" { 
		shopID = cau_hinh.BienCauHinh.IdFileSheet 
	}

	raw, err := loadSheetData(shopID, "KHACH_HANG")
	if err != nil { return }

	list := []*KhachHang{}

	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := layString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }

		kh := &KhachHang{
			SpreadsheetID:  shopID, // Lưu ID đã chuẩn hóa
			DongTrongSheet: i + 1,
			MaKhachHang:    maKH,
			TenDangNhap:    layString(r, CotKH_TenDangNhap),
			Email:          layString(r, CotKH_Email),
			MatKhauHash:    layString(r, CotKH_MatKhauHash),
			MaPinHash:      layString(r, CotKH_MaPinHash),
			VaiTroQuyenHan: layString(r, CotKH_VaiTroQuyenHan),
			ChucVu:         layString(r, CotKH_ChucVu),
			TrangThai:      layInt(r, CotKH_TrangThai),
			TenKhachHang:   layString(r, CotKH_TenKhachHang),
			NguonKhachHang: layString(r, CotKH_NguonKhachHang),
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

		_ = json.Unmarshal([]byte(layString(r, CotKH_RefreshTokenJson)), &kh.RefreshTokens)
		if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]TokenInfo) }
		
		_ = json.Unmarshal([]byte(layString(r, CotKH_DataSheetsJson)), &kh.DataSheets)
		_ = json.Unmarshal([]byte(layString(r, CotKH_GoiDichVuJson)), &kh.GoiDichVu)
		_ = json.Unmarshal([]byte(layString(r, CotKH_MangXaHoiJson)), &kh.MangXaHoi)
		_ = json.Unmarshal([]byte(layString(r, CotKH_ViTienJson)), &kh.ViTien)
		_ = json.Unmarshal([]byte(layString(r, CotKH_CauHinhJson)), &kh.CauHinh)
		_ = json.Unmarshal([]byte(layString(r, CotKH_InboxJson)), &kh.Inbox)

		list = append(list, kh)

		// Lưu vào Map Phẳng
		key := TaoCompositeKey(shopID, maKH)
		CacheMapKhachHang[key] = kh
	}

	KhoaHeThong.Lock()
	CacheKhachHang[shopID] = list // Lưu vào Key ID thật (không phải rỗng)
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

func LayKhachHang(shopID, maKH string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(shopID, maKH)
	kh, ok := CacheMapKhachHang[key]
	return kh, ok
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := kh.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet } // Fallback an toàn
	
	CacheKhachHang[sID] = append(CacheKhachHang[sID], kh)
	
	key := TaoCompositeKey(sID, kh.MaKhachHang)
	CacheMapKhachHang[key] = kh
}
