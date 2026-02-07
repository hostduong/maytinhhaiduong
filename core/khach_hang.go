package core

import (
	"fmt"
	"strings"
	"time"
)

// =============================================================
// 1. CẤU HÌNH CỘT (KHACH_HANG: A -> AA)
// =============================================================
const (
	CotKH_MaKhachHang      = 0  // A
	CotKH_TenDangNhap      = 1  // B
	CotKH_MatKhauHash      = 2  // C
	CotKH_Cookie           = 3  // D
	CotKH_CookieExpired    = 4  // E
	CotKH_MaPinHash        = 5  // F
	CotKH_LoaiKhachHang    = 6  // G
	CotKH_TenKhachHang     = 7  // H
	CotKH_DienThoai        = 8  // I
	CotKH_Email            = 9  // J
	CotKH_UrlFb            = 10 // K
	CotKH_Zalo             = 11 // L
	CotKH_UrlTele          = 12 // M
	CotKH_UrlTiktok        = 13 // N
	CotKH_DiaChi           = 14 // O
	CotKH_NgaySinh         = 15 // P
	CotKH_GioiTinh         = 16 // Q (Nam/Nữ)
	CotKH_MaSoThue         = 17 // R
	CotKH_DangNo           = 18 // S
	CotKH_TongMua          = 19 // T
	CotKH_ChucVu           = 20 // U
	CotKH_VaiTroQuyenHan   = 21 // V
	CotKH_TrangThai        = 22 // W
	CotKH_GhiChu           = 23 // X
	CotKH_NguoiTao         = 24 // Y
	CotKH_NgayTao          = 25 // Z
	CotKH_NgayCapNhat      = 26 // AA
)

// =============================================================
// 2. STRUCT DỮ LIỆU
// =============================================================
type KhachHang struct {
	// Dùng con trỏ để update trạng thái (Cookie) dễ dàng hơn
	DongTrongSheet int `json:"-"` 

	MaKhachHang      string  `json:"ma_khach_hang"`
	TenDangNhap      string  `json:"ten_dang_nhap"`
	MatKhauHash      string  `json:"-"` // Ẩn khi trả về JSON
	Cookie           string  `json:"-"` 
	CookieExpired    int64   `json:"cookie_expired"`
	MaPinHash        string  `json:"-"`
	LoaiKhachHang    string  `json:"loai_khach_hang"`
	TenKhachHang     string  `json:"ten_khach_hang"`
	DienThoai        string  `json:"dien_thoai"`
	Email            string  `json:"email"`
	UrlFb            string  `json:"url_fb"`
	Zalo             string  `json:"zalo"`
	UrlTele          string  `json:"url_tele"`
	UrlTiktok        string  `json:"url_tiktok"`
	DiaChi           string  `json:"dia_chi"`
	NgaySinh         string  `json:"ngay_sinh"`
	GioiTinh         string  `json:"gioi_tinh"`
	MaSoThue         string  `json:"ma_so_thue"`
	DangNo           float64 `json:"dang_no"`
	TongMua          float64 `json:"tong_mua"`
	ChucVu           string  `json:"chuc_vu"`
	VaiTroQuyenHan   string  `json:"vai_tro_quyen_han"`
	TrangThai        int     `json:"trang_thai"`
	GhiChu           string  `json:"ghi_chu"`
	NguoiTao         string  `json:"nguoi_tao"`
	NgayTao          string  `json:"ngay_tao"`
	NgayCapNhat      string  `json:"ngay_cap_nhat"`
}

// =============================================================
// 3. KHO LƯU TRỮ (In-Memory)
// =============================================================
var (
	_DS_KhachHang  []*KhachHang          // Slice chứa con trỏ
	_Map_KhachHang map[string]*KhachHang // Map tìm nhanh theo ID
)

// =============================================================
// 4. LOGIC NẠP DỮ LIỆU
// =============================================================
func NapKhachHang() {
	raw, err := loadSheetData("KHACH_HANG")
	if err != nil { return }

	tempList := []*KhachHang{}
	tempMap := make(map[string]*KhachHang)

	for i, r := range raw {
		if i < DongBatDauDuLieu-1 { continue }
		
		maKH := layString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }

		// Tạo struct
		kh := &KhachHang{
			DongTrongSheet: i + 1, // Google Sheet index từ 1
			MaKhachHang:    maKH,
			TenDangNhap:    layString(r, CotKH_TenDangNhap),
			MatKhauHash:    layString(r, CotKH_MatKhauHash),
			Cookie:         layString(r, CotKH_Cookie),
			CookieExpired:  int64(layFloat(r, CotKH_CookieExpired)),
			MaPinHash:      layString(r, CotKH_MaPinHash),
			LoaiKhachHang:  layString(r, CotKH_LoaiKhachHang),
			TenKhachHang:   layString(r, CotKH_TenKhachHang),
			DienThoai:      layString(r, CotKH_DienThoai),
			Email:          layString(r, CotKH_Email),
			UrlFb:          layString(r, CotKH_UrlFb),
			Zalo:           layString(r, CotKH_Zalo),
			UrlTele:        layString(r, CotKH_UrlTele),
			UrlTiktok:      layString(r, CotKH_UrlTiktok),
			DiaChi:         layString(r, CotKH_DiaChi),
			NgaySinh:       layString(r, CotKH_NgaySinh),
			GioiTinh:       layString(r, CotKH_GioiTinh),
			MaSoThue:       layString(r, CotKH_MaSoThue),
			DangNo:         layFloat(r, CotKH_DangNo),
			TongMua:        layFloat(r, CotKH_TongMua),
			ChucVu:         layString(r, CotKH_ChucVu),
			VaiTroQuyenHan: layString(r, CotKH_VaiTroQuyenHan),
			TrangThai:      layInt(r, CotKH_TrangThai),
			GhiChu:         layString(r, CotKH_GhiChu),
			NguoiTao:       layString(r, CotKH_NguoiTao),
			NgayTao:        layString(r, CotKH_NgayTao),
			NgayCapNhat:    layString(r, CotKH_NgayCapNhat),
		}

		tempList = append(tempList, kh)
		tempMap[maKH] = kh
	}

	// Cập nhật biến toàn cục
	_DS_KhachHang = tempList
	_Map_KhachHang = tempMap
}

// =============================================================
// 5. CÁC HÀM TRUY VẤN & NGHIỆP VỤ (Logic Nghiệp Vụ)
// =============================================================

// Lấy danh sách khách hàng (Trả về copy để an toàn)
func LayDanhSachKhachHang() []*KhachHang {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	// Copy slice of pointers
	kq := make([]*KhachHang, len(_DS_KhachHang))
	copy(kq, _DS_KhachHang)
	return kq
}

// Tìm khách hàng theo ID
func LayKhachHang(maKH string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	kh, ok := _Map_KhachHang[maKH]
	return kh, ok
}

// Tìm khách hàng theo Cookie (Dùng cho Auth)
func TimKhachHangTheoCookie(cookie string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()

	for _, kh := range _DS_KhachHang {
		if kh.Cookie == cookie && cookie != "" {
			// Kiểm tra hạn cookie
			if time.Now().Unix() > kh.CookieExpired {
				return nil, false
			}
			return kh, true
		}
	}
	return nil, false
}

// Tìm theo Tên Đăng Nhập hoặc Email (Dùng cho Login)
func TimKhachHangTheoUserOrEmail(input string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()

	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range _DS_KhachHang {
		if strings.ToLower(kh.TenDangNhap) == input { return kh, true }
		if strings.ToLower(kh.Email) == input && input != "" { return kh, true }
	}
	return nil, false
}

// Tạo mã khách hàng mới (KH_0001)
func TaoMaKhachHangMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()

	maxID := 0
	for _, kh := range _DS_KhachHang {
		parts := strings.Split(kh.MaKhachHang, "_")
		if len(parts) == 2 {
			var id int
			fmt.Sscanf(parts[1], "%d", &id)
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}

// Kiểm tra trùng User/Email
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

// Hàm thêm khách hàng vào RAM (Logic Ghi Sheet sẽ gọi ở layer khác hoặc qua Callback)
func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()

	_DS_KhachHang = append(_DS_KhachHang, kh)
	_Map_KhachHang[kh.MaKhachHang] = kh
}
