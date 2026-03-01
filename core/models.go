package core

import "sync"

// ==============================================================================
// ĐỊNH NGHĨA TÊN SHEET CHUẨN
// ==============================================================================
const (
	TenSheetPhanQuyen        = "PHAN_QUYEN"
	TenSheetDanhMuc          = "DANH_MUC"
	TenSheetThuongHieu       = "THUONG_HIEU"
	TenSheetBienLoiNhuan     = "BIEN_LOI_NHUAN"
	TenSheetNhaCungCap       = "NHA_CUNG_CAP"
	TenSheetKhachHang        = "KHACH_HANG"
	TenSheetMayTinh          = "MAY_TINH"
	TenSheetTinNhan          = "TIN_NHAN"
	TenSheetPhieuNhap        = "PHIEU_NHAP"
	TenSheetChiTietPhieuNhap = "CHI_TIET_PHIEU_NHAP"
	TenSheetSerial           = "SERIAL_SAN_PHAM"
)

// ==============================================================================
// CẤU TRÚC: CÀI ĐẶT CẤU HÌNH & BẢO MẬT
// ==============================================================================
const (
	DongBatDau_PhanQuyen = 11
	CotPQ_MaChucNang     = 0
	CotPQ_Nhom           = 1
	CotPQ_MoTa           = 2
	CotPQ_StartRole      = 3
)

type VaiTroInfo struct {
	MaVaiTro   string `json:"ma_vai_tro"`
	TenVaiTro  string `json:"ten_vai_tro"`
	StyleLevel int    `json:"style_level"` // Cấp bậc từ 0-9
	StyleTheme int    `json:"style_theme"` // Đuôi giao diện
}

// ==============================================================================
// CẤU TRÚC: KHÁCH HÀNG & NHÂN VIÊN
// ==============================================================================
const (
	DongBatDau_KhachHang = 11
	CotKH_MaKhachHang      = 0
	CotKH_TenDangNhap      = 1
	CotKH_Email            = 2
	CotKH_MatKhauHash      = 3
	CotKH_MaPinHash        = 4
	CotKH_RefreshTokenJson = 5
	CotKH_VaiTroQuyenHan   = 6
	CotKH_ChucVu           = 7
	CotKH_TrangThai        = 8
	CotKH_DataSheetsJson   = 9
	CotKH_GoiDichVuJson    = 10
	CotKH_CauHinhJson      = 11
	CotKH_NguonKhachHang   = 12
	CotKH_TenKhachHang     = 13
	CotKH_DienThoai        = 14
	CotKH_AnhDaiDien       = 15
	CotKH_MangXaHoiJson    = 16
	CotKH_DiaChi           = 17
	CotKH_NgaySinh         = 18
	CotKH_GioiTinh         = 19
	CotKH_MaSoThue         = 20
	CotKH_ViTienJson       = 21
	CotKH_GhiChu           = 22
	CotKH_NgayTao          = 23
	CotKH_NguoiCapNhat     = 24
	CotKH_NgayCapNhat      = 25
)

// Các struct phụ trợ cho KhachHang
type TokenInfo struct { DeviceName string `json:"dev"`; ExpiresAt int64 `json:"exp"` }
type DataSheetInfo struct { SpreadsheetID string `json:"sheet_id"`; GoogleAuthJson string `json:"google_auth_json"`; FolderDriveID string `json:"folder_drive_id"` }
type PlanInfo struct { MaGoi string `json:"ma_goi"`; TenGoi string `json:"ten_goi"`; NgayHetHan string `json:"ngay_het_han"`; TrangThai string `json:"trang_thai"` }
type UserConfig struct { Theme string `json:"theme"`; ChuyenNganh string `json:"chuyen_nganh"`; CustomDomain string `json:"custom_domain"`; DarkMode bool `json:"dark_mode"`; Language string `json:"lang"` }
type SocialInfo struct { Zalo string `json:"zalo"`; Facebook string `json:"fb"`; Tiktok string `json:"tiktok"` }
type WalletInfo struct { SoDuHienTai float64 `json:"so_du"` }

type KhachHang struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaKhachHang    string `json:"ma_khach_hang"`
	TenDangNhap    string `json:"ten_dang_nhap"`
	Email          string `json:"email"`
	MatKhauHash    string `json:"-"`
	MaPinHash      string `json:"-"`
	RefreshTokens  map[string]TokenInfo `json:"-"`

	VaiTroQuyenHan string `json:"vai_tro_quyen_han"`
	ChucVu         string `json:"chuc_vu"`
	TrangThai      int    `json:"trang_thai"`

	DataSheets     DataSheetInfo `json:"data_sheets"`
	GoiDichVu      []PlanInfo    `json:"goi_dich_vu"`
	CauHinh        UserConfig    `json:"cau_hinh"`

	NguonKhachHang string     `json:"nguon_khach_hang"`
	TenKhachHang   string     `json:"ten_khach_hang"`
	DienThoai      string     `json:"dien_thoai"`
	AnhDaiDien     string     `json:"anh_dai_dien"`
	MangXaHoi      SocialInfo `json:"mang_xa_hoi"`
	DiaChi         string     `json:"dia_chi"`
	NgaySinh       string     `json:"ngay_sinh"`
	GioiTinh       int        `json:"gioi_tinh"`
	MaSoThue       string     `json:"ma_so_thue"`
	ViTien         WalletInfo `json:"vi_tien"`
	
	Inbox          []*TinNhan `json:"-"`
	StyleLevel     int        `json:"-"`
	StyleTheme     int        `json:"-"`

	GhiChu         string `json:"ghi_chu"`
	NgayTao        string `json:"ngay_tao"`
	NguoiCapNhat   string `json:"nguoi_cap_nhat"`
	NgayCapNhat    string `json:"ngay_cap_nhat"`
}

// ==============================================================================
// CẤU TRÚC: NHÀ CUNG CẤP & MASTER DATA
// ==============================================================================
const (
	DongBatDau_NhaCungCap = 11
	CotNCC_MaNhaCungCap       = 0
	CotNCC_TenNhaCungCap      = 1
	CotNCC_MaSoThue           = 2
	CotNCC_DienThoai          = 3
	CotNCC_Email              = 4
	CotNCC_KhuVuc             = 5
	CotNCC_DiaChi             = 6
	CotNCC_NguoiLienHe        = 7
	CotNCC_NganHang           = 8
	CotNCC_NhomNhaCungCap     = 9
	CotNCC_LoaiNhaCungCap     = 10
	CotNCC_DieuKhoanThanhToan = 11
	CotNCC_ChietKhauMacDinh   = 12
	CotNCC_HanMucCongNo       = 13
	CotNCC_CongNoDauKy        = 14
	CotNCC_TongMua            = 15
	CotNCC_NoCanTra           = 16
	CotNCC_ThongTinThemJson   = 17
	CotNCC_TrangThai          = 18
	CotNCC_GhiChu             = 19
	CotNCC_NguoiTao           = 20
	CotNCC_NgayTao            = 21
	CotNCC_NgayCapNhat        = 22
)

type NhaCungCap struct {
	SpreadsheetID      string  `json:"-"`
	DongTrongSheet     int     `json:"-"`
	MaNhaCungCap       string  `json:"ma_nha_cung_cap"`
	TenNhaCungCap      string  `json:"ten_nha_cung_cap"`
	MaSoThue           string  `json:"ma_so_thue"`
	DienThoai          string  `json:"dien_thoai"`
	Email              string  `json:"email"`
	KhuVuc             string  `json:"khu_vuc"`
	DiaChi             string  `json:"dia_chi"`
	NguoiLienHe        string  `json:"nguoi_lien_he"`
	NganHang           string  `json:"ngan_hang"`
	NhomNhaCungCap     string  `json:"nhom_nha_cung_cap"`
	LoaiNhaCungCap     string  `json:"loai_nha_cung_cap"`
	DieuKhoanThanhToan string  `json:"dieu_khoan_thanh_toan"`
	ChietKhauMacDinh   float64 `json:"chiet_khau_mac_dinh"`
	HanMucCongNo       float64 `json:"han_muc_cong_no"`
	CongNoDauKy        float64 `json:"cong_no_dau_ky"`
	TongMua            float64 `json:"tong_mua"`
	NoCanTra           float64 `json:"no_can_tra"`
	ThongTinThemJson   string  `json:"thong_tin_them_json"`
	TrangThai          int     `json:"trang_thai"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`
}

// Bổ sung các Struct còn lại: DanhMuc, ThuongHieu, SanPhamMayTinh, TinNhan, PhieuNhap...
// (Vì khuôn khổ có hạn, bạn gộp nốt Struct từ các file cũ vào đây. Nguyên tắc là Cột + Struct nằm chung).
