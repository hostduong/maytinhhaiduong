package core

// ==============================================================================
// 1. ĐỊNH NGHĨA TÊN SHEET CHUẨN CHO 3 VŨ TRỤ (3-TIER ARCHITECTURE)
// ==============================================================================

// --- VŨ TRỤ 1: MASTER (sss.99k.vn) ---
const (
	TenSheetKhachHangMaster  = "KHACH_HANG_MASTER"
	TenSheetPhanQuyenMaster  = "PHAN_QUYEN_MASTER"
	TenSheetGoiDichVuMaster  = "GOI_DICH_VU_MASTER"
	TenSheetTinNhanMaster    = "TIN_NHAN_MASTER"
	TenSheetCauHinhThuocTinh = "CAU_HINH_THUOC_TINH" // Sheet cấu hình Meta-data (Dynamic UI)
)

// --- VŨ TRỤ 2: ADMIN (admin.99k.vn) ---
const (
	TenSheetKhachHangAdmin   = "KHACH_HANG_ADMIN"
	TenSheetPhanQuyenAdmin   = "PHAN_QUYEN_ADMIN"
)

// --- VŨ TRỤ 3: CỬA HÀNG ([cuahang].99k.vn) ---
const (
	TenSheetKhachHang        = "KHACH_HANG"
	TenSheetPhanQuyen        = "PHAN_QUYEN"
	
	// Dữ liệu Vận hành kinh doanh (Chỉ có ở Shop)
	TenSheetDanhMuc          = "DANH_MUC"
	TenSheetThuongHieu       = "THUONG_HIEU"
	TenSheetBienLoiNhuan     = "BIEN_LOI_NHUAN"
	TenSheetNhaCungCap       = "NHA_CUNG_CAP"
	TenSheetSerial           = "SERIAL_SAN_PHAM"
	
	TenSheetPhieuNhap        = "PHIEU_NHAP"
	TenSheetChiTietPhieuNhap = "CHI_TIET_PHIEU_NHAP"
	TenSheetPhieuXuat        = "PHIEU_XUAT"
	TenSheetChiTietPhieuXuat = "CHI_TIET_PHIEU_XUAT"

	// Các Sheet mới cập nhật
	TenSheetHoaDon           = "HOA_DON"
	TenSheetHoaDonChiTiet    = "HOA_DON_CHI_TIET"
	TenSheetPhieuThuChi      = "PHIEU_THU_CHI"
	TenSheetPhieuBaoHanh     = "PHIEU_BAO_HANH"
)

// ==============================================================================
// 2. CẤU TRÚC PHÂN QUYỀN (Dùng chung cho cả 3 vũ trụ, khác nhau ở Level Check)
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
	StyleLevel int    `json:"style_level"` // Cấp bậc: Master(0-9), Admin(3-9), Shop(4-9)
	StyleTheme int    `json:"style_theme"` 
}

// ==============================================================================
// 3. CẤU TRÚC KHÁCH HÀNG: MASTER & ADMIN (FULL 26 CỘT)
// ==============================================================================
const (
	DongBatDau_KhachHang = 11
	CotKH_MaKhachHang      = 0  // A
	CotKH_TenDangNhap      = 1  // B
	CotKH_Email            = 2  // C
	CotKH_MatKhauHash      = 3  // D
	CotKH_MaPinHash        = 4  // E
	CotKH_RefreshTokenJson = 5  // F
	CotKH_VaiTroQuyenHan   = 6  // G
	CotKH_ChucVu           = 7  // H
	CotKH_TrangThai        = 8  // I
	CotKH_DataSheetsJson   = 9  // J
	CotKH_GoiDichVuJson    = 10 // K
	CotKH_CauHinhJson      = 11 // L
	CotKH_NguonKhachHang   = 12 // M
	CotKH_TenKhachHang     = 13 // N
	CotKH_DienThoai        = 14 // O
	CotKH_AnhDaiDien       = 15 // P
	CotKH_MangXaHoiJson    = 16 // Q
	CotKH_DiaChi           = 17 // R
	CotKH_NgaySinh         = 18 // S
	CotKH_GioiTinh         = 19 // T
	CotKH_MaSoThue         = 20 // U
	CotKH_ViTienJson       = 21 // V
	CotKH_GhiChu           = 22 // W
	CotKH_NgayTao          = 23 // X
	CotKH_NguoiCapNhat     = 24 // Y
	CotKH_NgayCapNhat      = 25 // Z
)

type TokenInfo struct { DeviceName string `json:"dev"`; ExpiresAt int64 `json:"exp"` }
type DataSheetInfo struct { SpreadsheetID string `json:"sheet_id"`; GoogleAuthJson string `json:"google_auth_json"`; FolderDriveID string `json:"folder_drive_id"` }
type PlanInfo struct {
	MaGoi          string `json:"ma_goi"`
	TenGoi         string `json:"ten_goi"`
	LoaiGoi        string `json:"loai_goi"`       
	NgayHetHan     string `json:"ngay_het_han"`
	TrangThai      string `json:"trang_thai"`     
	MaxSanPham     int    `json:"max_san_pham"`   
	MaxNhanVien    int    `json:"max_nhan_vien"`  
	
	// CÁC QUYỀN LỢI MỚI MỞ RỘNG CỦA GÓI CƯỚC
	TinNhan        bool   `json:"tin_nhan"`       // true: Bật module Chat CSKH
	TenMienFree    bool   `json:"ten_mien_free"`  // true: Được dùng Subdomain hệ thống
	TenMienRieng   bool   `json:"ten_mien_rieng"` // true: Được trỏ Custom Domain
}

type UserConfig struct { 
	Theme         string `json:"theme"` 
	ChuyenNganh   string `json:"chuyen_nganh"` 
	Language      string `json:"lang"` 
	DarkMode      bool   `json:"dark_mode"` 
	
	// CẤU HÌNH TÊN MIỀN VÀ KHÔNG GIAN
	CustomDomain  string `json:"custom_domain"`   // Ưu tiên 1 (VD: shopcuatoi.com)
	Subdomain     string `json:"subdomain"`       // Ưu tiên 2 (VD: shopcuatoi.99k.vn)
	Website       bool   `json:"website"`         // true: Có Web bán hàng / false: Chỉ dùng Admin Kế toán
}
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
// 4. CẤU TRÚC KHÁCH HÀNG: DÀNH RIÊNG CHO SHOP (17 CỘT - BẢN RÚT GỌN SIÊU NHẸ)
// ==============================================================================
const (
	DongBatDau_KhachHangShop = 11
	CotKHS_MaKhachHang    = 0  // A
	CotKHS_TenDangNhap    = 1  // B
	CotKHS_Email          = 2  // C
	CotKHS_MatKhauHash    = 3  // D
	CotKHS_MaPinHash      = 4  // E
	CotKHS_RefreshToken   = 5  // F
	CotKHS_VaiTro         = 6  // G
	CotKHS_ChucVu         = 7  // H
	CotKHS_TrangThai      = 8  // I
	CotKHS_TenKhachHang   = 9  // J
	CotKHS_DienThoai      = 10 // K
	CotKHS_AnhDaiDien     = 11 // L
	CotKHS_DiaChi         = 12 // M
	CotKHS_NgaySinh       = 13 // N
	CotKHS_GioiTinh       = 14 // O
	CotKHS_GhiChu         = 15 // P
	CotKHS_NgayTao        = 16 // Q
)

type KhachHangShop struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaKhachHang   string `json:"ma_khach_hang"`
	TenDangNhap   string `json:"ten_dang_nhap"`
	Email         string `json:"email"`
	MatKhauHash   string `json:"-"`
	MaPinHash     string `json:"-"`
	RefreshTokens map[string]TokenInfo `json:"-"`
	
	VaiTroQuyenHan string `json:"vai_tro_quyen_han"`
	ChucVu         string `json:"chuc_vu"`
	TrangThai      int    `json:"trang_thai"`
	
	TenKhachHang  string `json:"ten_khach_hang"`
	DienThoai     string `json:"dien_thoai"`
	AnhDaiDien    string `json:"anh_dai_dien"`
	DiaChi        string `json:"dia_chi"`
	NgaySinh      string `json:"ngay_sinh"`
	GioiTinh      int    `json:"gioi_tinh"`
	GhiChu        string `json:"ghi_chu"`
	NgayTao       string `json:"ngay_tao"`
}

// ==============================================================================
// CẤU TRÚC ĐỘNG: CẤU HÌNH THUỘC TÍNH (EAV META-DATA)
// ==============================================================================
const (
	DongBatDau_CauHinhThuocTinh = 11
	CotCHTT_MaThuocTinh  = 0 // A
	CotCHTT_TenThuocTinh = 1 // B
	CotCHTT_KieuNhap     = 2 // C
	CotCHTT_DonVi        = 3 // D
	CotCHTT_StartNganh   = 4 // E trở đi (Chứa JSON cấu hình)
)

// ConfigNganhHang dùng để hứng JSON trên dòng 1 của Cột Ngành
type ConfigNganhHang struct {
	MaNganh    string `json:"ma"`
	TenSheet   string `json:"sheet"`
	TenHienThi string `json:"ten"`
	Icon       string `json:"icon,omitempty"`
	Sort       int    `json:"sort,omitempty"`
}

type ThuocTinhNganh struct {
	MaThuocTinh  string `json:"ma_thuoc_tinh"`
	TenThuocTinh string `json:"ten_thuoc_tinh"`
	KieuNhap     string `json:"kieu_nhap"`
	DonVi        string `json:"don_vi"`
}

// ==============================================================================
// CẤU TRÚC CHUẨN: SẢN PHẨM NOSQL (DÙNG CHUNG MỌI NGÀNH HÀNG)
// ==============================================================================
const (
	DongBatDau_Product = 2
	CotProd_MaSanPham  = 0 // Cột A
	CotProd_DataJSON   = 1 // Cột B
)

// ProductRow đại diện cho 1 dòng vật lý trên Google Sheets
type ProductRow struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaSanPham      string `json:"ma_san_pham"`
	DataJSON       string `json:"data_json"`
}

// ProductJSON định nghĩa lõi dữ liệu sẽ được Serialize/Deserialize từ DataJSON (Cột B)
type ProductJSON struct {
	MaSanPham    string        `json:"ma_san_pham"`
	Version      int           `json:"version"`      // Quản lý migrate schema
	CreatedAt    int64         `json:"created_at"`   // Unix timestamp
	UpdatedAt    int64         `json:"updated_at"`   // Unix timestamp
	TrangThai    int           `json:"trang_thai"`   // 1=hoạt động, 0=ẩn, -1=xóa mềm, 2=hết hàng, v.v.
	MaNganh      string        `json:"ma_nganh"`
	MaDanhMuc    []string      `json:"ma_danh_muc"`
	MaThuongHieu string        `json:"ma_thuong_hieu"`
	TenSanPham   string        `json:"ten_san_pham"`
	TenRutGon    string        `json:"ten_rut_gon"`
	Slug         string        `json:"slug"`
	NenTang      []string      `json:"nen_tang"`     // web, shopee, pos...
	LoaiSanPham  string        `json:"loai_san_pham"`
	SearchText   string        `json:"search_text"`  // Text không dấu tối ưu query RAM
	SearchBoost  int           `json:"search_boost"`
	Sort         int           `json:"sort"`
	Views        int           `json:"views"`
	DaBan        int           `json:"da_ban"`
	SKUChinh     string        `json:"sku_chinh"`
	Tags         []string      `json:"tags"`
	SKU          []ProductSKU  `json:"sku"`
	SEO          ProductSEO    `json:"seo"`
	QuanLy       ProductQuanLy `json:"quan_ly"`
}

type ProductSKU struct {
	MaSanPham    string                 `json:"ma_san_pham"`
	MaSKU        string                 `json:"ma_sku"`
	TrangThai    int                    `json:"trang_thai"`
	TenSKU       string                 `json:"ten_sku"`
	Barcode      string                 `json:"barcode"`
	TinhTrang    string                 `json:"tinh_trang"`
	BaoHanh      string                 `json:"bao_hanh"`
	MaNhaCungCap string                 `json:"ma_nha_cung_cap"`
	XuatXu       string                 `json:"xuat_xu"`
	DonVi        string                 `json:"don_vi"`
	GhiChu       string                 `json:"ghi_chu"`
	AnhDaiDien   string                 `json:"anh_dai_dien"`
	HinhAnh      []string               `json:"hinh_anh"`
	MoTaHTML     string                 `json:"mo_ta_html"`
	ThongSoHTML  string                 `json:"thong_so_html"`
	Gia          ProductGia             `json:"gia"`
	TonKho       int                    `json:"ton_kho"`
	DaBan        int                    `json:"da_ban"`
	DatTruoc     int                    `json:"dat_truoc"`
	ThuocTinh    map[string]interface{} `json:"thuoc_tinh"` // Ma trận thuộc tính động (Dynamic EAV)
}

type ProductGia struct {
	GiaNhap     float64 `json:"gia_nhap"`
	PhanTramLai float64 `json:"phan_tram_lai"`
	PhanTramVAT float64 `json:"phan_tram_vat"`
	ChiPhiNhap  float64 `json:"chi_phi_nhap"`
	GiaNiemYet  float64 `json:"gia_niem_yet"`
	GiaBan      float64 `json:"gia_ban"`
}

type ProductSEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	OGImage     string `json:"og_image"`
}

type ProductQuanLy struct {
	NguoiTao     string `json:"nguoi_tao"`
	NgayTao      string `json:"ngay_tao"` // Có thể dùng timestamp hoặc chuỗi ISO
	NguoiCapNhat string `json:"nguoi_cap_nhat"`
	NgayCapNhat  string `json:"ngay_cap_nhat"`
}


// ==============================================================================
// 5. CÁC CẤU TRÚC VẬN HÀNH: TÀI CHÍNH & HÓA ĐƠN
// ==============================================================================

// --- HÓA ĐƠN TÀI CHÍNH (20 Cột) ---
const (
	DongBatDau_HoaDon = 11
	CotHD_MaHoaDon           = 0  // A
	CotHD_MaTraCuu           = 1  // B
	CotHD_XmlUrl             = 2  // C
	CotHD_LoaiHoaDon         = 3  // D
	CotHD_MaPhieuXuat        = 4  // E
	CotHD_MaKhachHang        = 5  // F
	CotHD_NgayHoaDon         = 6  // G
	CotHD_KyHieu             = 7  // H
	CotHD_SoHoaDon           = 8  // I
	CotHD_MauSo              = 9  // J
	CotHD_LinkChungTu        = 10 // K
	CotHD_TongTienPhieu      = 11 // L
	CotHD_TongVAT            = 12 // M
	CotHD_TongTienSauVAT     = 13 // N
	CotHD_TrangThai          = 14 // O
	CotHD_TrangThaiThanhToan = 15 // P
	CotHD_GhiChu             = 16 // Q
	CotHD_NguoiTao           = 17 // R
	CotHD_NgayTao            = 18 // S
	CotHD_NgayCapNhat        = 19 // T
)

type HoaDon struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaHoaDon           string  `json:"ma_hoa_don"`
	MaTraCuu           string  `json:"ma_tra_cuu"`
	XmlUrl             string  `json:"xml_url"`
	LoaiHoaDon         string  `json:"loai_hoa_don"`
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`
	MaKhachHang        string  `json:"ma_khach_hang"`
	NgayHoaDon         string  `json:"ngay_hoa_don"`
	KyHieu             string  `json:"ky_hieu"`
	SoHoaDon           string  `json:"so_hoa_don"`
	MauSo              string  `json:"mau_so"`
	LinkChungTu        string  `json:"link_chung_tu"`
	TongTienPhieu      float64 `json:"tong_tien_phieu"`
	TongVAT            float64 `json:"tong_vat"`
	TongTienSauVAT     float64 `json:"tong_tien_sau_vat"`
	TrangThai          int     `json:"trang_thai"`
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`

	ChiTiet            []*HoaDonChiTiet `json:"chi_tiet"`
}

// --- HÓA ĐƠN CHI TIẾT (11 Cột) ---
const (
	CotHDCT_MaHoaDon    = 0  // A
	CotHDCT_MaSanPham   = 1  // B
	CotHDCT_MaSKU       = 2  // C
	CotHDCT_MaNganhHang = 3  // D
	CotHDCT_TenSanPham  = 4  // E
	CotHDCT_DonVi       = 5  // F
	CotHDCT_SoLuong     = 6  // G
	CotHDCT_DonGiaBan   = 7  // H
	CotHDCT_VATPercent  = 8  // I
	CotHDCT_TienVAT     = 9  // J
	CotHDCT_ThanhTien   = 10 // K
)

type HoaDonChiTiet struct {
	MaHoaDon    string  `json:"ma_hoa_don"`
	MaSanPham   string  `json:"ma_san_pham"`
	MaSKU       string  `json:"ma_sku"`
	MaNganhHang string  `json:"ma_nganh_hang"`
	TenSanPham  string  `json:"ten_san_pham"`
	DonVi       string  `json:"don_vi"`
	SoLuong     int     `json:"so_luong"`
	DonGiaBan   float64 `json:"don_gia_ban"`
	VATPercent  float64 `json:"vat_percent"`
	TienVAT     float64 `json:"tien_vat"`
	ThanhTien   float64 `json:"thanh_tien"`
}

// --- PHIẾU THU CHI QUỸ KẾ TOÁN (16 Cột) ---
const (
	DongBatDau_PhieuThuChi = 11
	CotPTC_MaPhieu             = 0  // A
	CotPTC_NgayTaoPhieu        = 1  // B
	CotPTC_LoaiPhieu           = 2  // C
	CotPTC_DoiTuongLoai        = 3  // D
	CotPTC_DoiTuongID          = 4  // E
	CotPTC_HangMucThuChi       = 5  // F
	CotPTC_CoHoaDonDo          = 6  // G
	CotPTC_MaChungTuThamChieu  = 7  // H
	CotPTC_SoTien              = 8  // I
	CotPTC_PhuongThucThanhToan = 9  // J
	CotPTC_TrangThaiDuyet      = 10 // K
	CotPTC_NguoiDuyet          = 11 // L
	CotPTC_GhiChu              = 12 // M
	CotPTC_NguoiTao            = 13 // N
	CotPTC_NgayTao             = 14 // O
	CotPTC_NgayCapNhat         = 15 // P
)

type PhieuThuChi struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuThuChi        string  `json:"ma_phieu_thu_chi"`
	NgayTaoPhieu         string  `json:"ngay_tao_phieu"`
	LoaiPhieu            string  `json:"loai_phieu"`
	DoiTuongLoai         string  `json:"doi_tuong_loai"`
	DoiTuongID           string  `json:"doi_tuong_id"`
	HangMucThuChi        string  `json:"hang_muc_thu_chi"`
	CoHoaDonDo           int     `json:"co_hoa_don_do"`
	MaChungTuThamChieu   string  `json:"ma_chung_tu_tham_chieu"`
	SoTien               float64 `json:"so_tien"`
	PhuongThucThanhToan  string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiDuyet       int     `json:"trang_thai_duyet"`
	NguoiDuyet           string  `json:"nguoi_duyet"`
	GhiChu               string  `json:"ghi_chu"`
	NguoiTao             string  `json:"nguoi_tao"`
	NgayTao              string  `json:"ngay_tao"`
	NgayCapNhat          string  `json:"ngay_cap_nhat"`
}

// --- PHIẾU BẢO HÀNH KỸ THUẬT (22 Cột) ---
const (
	DongBatDau_PhieuBaoHanh = 11
	CotPBH_MaPhieu           = 0  // A
	CotPBH_LoaiPhieu         = 1  // B
	CotPBH_SerialIMEI        = 2  // C
	CotPBH_MaSanPham         = 3  // D
	CotPBH_MaKhachHang       = 4  // E
	CotPBH_TenNguoiGui       = 5  // F
	CotPBH_SDTNguoiGui       = 6  // G
	CotPBH_NgayNhan          = 7  // H
	CotPBH_TinhTrangLoi      = 8  // I
	CotPBH_HinhThuc          = 9  // J
	CotPBH_TrangThai         = 10 // K
	CotPBH_NgayTraDuKien     = 11 // L
	CotPBH_NgayTraThucTe     = 12 // M
	CotPBH_ChiPhiSua         = 13 // N
	CotPBH_PhiThuKhach       = 14 // O
	CotPBH_KetQuaSuaChua     = 15 // P
	CotPBH_LinhKienThayThe   = 16 // Q
	CotPBH_MaNhanVienKyThuat = 17 // R
	CotPBH_GhiChu            = 18 // S
	CotPBH_NguoiTao          = 19 // T
	CotPBH_NgayTao           = 20 // U
	CotPBH_NgayCapNhat       = 21 // V
)

type PhieuBaoHanh struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuBaoHanh    string  `json:"ma_phieu_bao_hanh"`
	LoaiPhieu         string  `json:"loai_phieu"`
	SerialIMEI        string  `json:"serial_imei"`
	MaSanPham         string  `json:"ma_san_pham"`
	MaKhachHang       string  `json:"ma_khach_hang"`
	TenNguoiGui       string  `json:"ten_nguoi_gui"`
	SDTNguoiGui       string  `json:"sdt_nguoi_gui"`
	NgayNhan          string  `json:"ngay_nhan"`
	TinhTrangLoi      string  `json:"tinh_trang_loi"`
	HinhThuc          string  `json:"hinh_thuc"`
	TrangThai         int     `json:"trang_thai"`
	NgayTraDuKien     string  `json:"ngay_tra_du_kien"`
	NgayTraThucTe     string  `json:"ngay_tra_thuc_te"`
	ChiPhiSua         float64 `json:"chi_phi_sua"`
	PhiThuKhach       float64 `json:"phi_thu_khach"`
	KetQuaSuaChua     string  `json:"ket_qua_sua_chua"`
	LinhKienThayThe   string  `json:"linh_kien_thay_the"`
	MaNhanVienKyThuat string  `json:"ma_nhan_vien_ky_thuat"`
	GhiChu            string  `json:"ghi_chu"`
	NguoiTao          string  `json:"nguoi_tao"`
	NgayTao           string  `json:"ngay_tao"`
	NgayCapNhat       string  `json:"ngay_cap_nhat"`
}


// ==============================================================================
// CẤU TRÚC: NHÀ CUNG CẤP & MASTER DATA
// ==============================================================================
const (
	DongBatDau_NhaCungCap = 11
	CotNCC_MaNhaCungCap       = 0  // A
	CotNCC_TenNhaCungCap      = 1  // B
	CotNCC_MaSoThue           = 2  // C
	CotNCC_DienThoai          = 3  // D
	CotNCC_Email              = 4  // E
	CotNCC_KhuVuc             = 5  // F
	CotNCC_DiaChi             = 6  // G
	CotNCC_NguoiLienHe        = 7  // H
	CotNCC_AnhDaiDien         = 8  // I
	CotNCC_NganHang           = 9  // J
	CotNCC_NhomNhaCungCap     = 10 // K
	CotNCC_LoaiNhaCungCap     = 11 // L
	CotNCC_DieuKhoanThanhToan = 12 // M
	CotNCC_ChietKhauMacDinh   = 13 // N
	CotNCC_HanMucCongNo       = 14 // O
	CotNCC_CongNoDauKy        = 15 // P
	CotNCC_TongMua            = 16 // Q
	CotNCC_NoCanTra           = 17 // R
	CotNCC_ThongTinThemJson   = 18 // S
	CotNCC_TrangThai          = 19 // T
	CotNCC_GhiChu             = 20 // U
	CotNCC_NguoiTao           = 21 // V
	CotNCC_NgayTao            = 22 // W
	CotNCC_NgayCapNhat        = 23 // X
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
	AnhDaiDien         string  `json:"anh_dai_dien"` 
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

// ==============================================================================
// CẤU TRÚC: DANH MỤC
// ==============================================================================
const (
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0
	CotDM_TenDanhMuc = 1
	CotDM_DanhMucMe  = 2
	CotDM_ThueVAT    = 3
	CotDM_LoiNhuan   = 4
	CotDM_Slot       = 5
	CotDM_TrangThai  = 6
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaDanhMuc      string `json:"ma_danh_muc"`
	TenDanhMuc     string `json:"ten_danh_muc"`
	DanhMucMe      string `json:"danh_muc_me"`
	ThueVAT        float64 `json:"thue_vat"`
	LoiNhuan       float64 `json:"bien_loi_nhuan"`
	Slot           int    `json:"slot"`
	TrangThai      int    `json:"trang_thai"`
}

// ==============================================================================
// CẤU TRÚC: THƯƠNG HIỆU
// ==============================================================================
const (
	DongBatDau_ThuongHieu = 11 
	CotTH_MaThuongHieu  = 0
	CotTH_TenThuongHieu = 1
	CotTH_LogoUrl       = 2
	CotTH_MoTa          = 3
	CotTH_TrangThai     = 4
)

type ThuongHieu struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaThuongHieu   string `json:"ma_thuong_hieu"`
	TenThuongHieu  string `json:"ten_thuong_hieu"`
	LogoUrl        string `json:"logo_url"`
	MoTa           string `json:"mo_ta"`
	TrangThai      int    `json:"trang_thai"`
}

// ==============================================================================
// CẤU TRÚC: BIÊN LỢI NHUẬN
// ==============================================================================
const (
	DongBatDau_BienLoiNhuan = 11
	CotBLN_KhungGiaNhap = 0
	CotBLN_BienLoiNhuan = 1
	CotBLN_TrangThai    = 2
)

type BienLoiNhuan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	GiaTu          float64 `json:"gia_tu"`
	KhungGiaNhap   float64 `json:"khung_gia_nhap"`
	BienLoiNhuan   float64 `json:"bien_loi_nhuan"`
	TrangThai      int     `json:"trang_thai"`
}

// ==============================================================================
// CẤU TRÚC: TIN NHẮN
// ==============================================================================
const (
	DongBatDau_TinNhan = 11

	CotTN_MaTinNhan    = 0  // A
	CotTN_LoaiTinNhan  = 1  // B
	CotTN_NguoiGuiID   = 2  // C
	CotTN_NguoiNhanID  = 3  // D
	CotTN_TieuDe       = 4  // E
	CotTN_NoiDung      = 5  // F
	CotTN_DinhKemJson  = 6  // G
	CotTN_ThamChieuID  = 7  // H
	CotTN_ReplyChoID   = 8  // I
	CotTN_NgayTao      = 9  // J
	CotTN_NguoiDocJson = 10 // K
	CotTN_TrangThaiXoa = 11 // L
)

type FileDinhKem struct {
	TenFile string `json:"name"`
	URL     string `json:"url"`
	Loai    string `json:"type"` 
}

type TinNhan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaTinNhan      string        `json:"ma_tin_nhan"`
	LoaiTinNhan    string        `json:"loai_tin_nhan"`
	NguoiGuiID     string        `json:"nguoi_gui_id"`
	NguoiNhanID    string        `json:"nguoi_nhan_id"`
	TieuDe         string        `json:"tieu_de"`
	NoiDung        string        `json:"noi_dung"`
	DinhKem        []FileDinhKem `json:"dinh_kem"`
	ThamChieuID    string        `json:"tham_chieu_id"`
	ReplyChoID     string        `json:"reply_cho_id"`
	NgayTao        string        `json:"ngay_tao"`
	NguoiDoc       []string      `json:"nguoi_doc"`
	TrangThaiXoa   []string      `json:"trang_thai_xoa"`
	
	DaDoc          bool          `json:"da_doc"` 
}

// ==============================================================================
// CẤU TRÚC: PHIẾU NHẬP KHÔNG GIAN (LƯU NHÁP JSON)
// ==============================================================================
const (
	DongBatDau_PhieuNhap       = 11
	CotPN_MaPhieuNhap          = 0  // A
	CotPN_MaNhaCungCap         = 1  // B
	CotPN_MaKho                = 2  // C
	CotPN_NgayNhap             = 3  // D
	CotPN_ChiTietJson          = 4  // E (Mới: Lưu mảng JSON chi tiết)
	CotPN_TrangThai            = 5  // F (-1: Xóa, 0: Nháp, 1: Hoàn thành, 2: Chờ duyệt)
	CotPN_SoHoaDon             = 6  // G
	CotPN_NgayHoaDon           = 7  // H
	CotPN_UrlChungTu           = 8  // I
	CotPN_TongTienPhieu        = 9  // J
	CotPN_GiamGiaPhieu         = 10 // K
	CotPN_ChiPhiNhap           = 11 // L (Mới: Phí vận chuyển, bốc vác...)
	CotPN_DaThanhToan          = 12 // M
	CotPN_ConNo                = 13 // N
	CotPN_PhuongThucThanhToan  = 14 // O
	CotPN_TrangThaiThanhToan   = 15 // P
	CotPN_GhiChu               = 16 // Q
	CotPN_NguoiTao             = 17 // R
	CotPN_NgayTao              = 18 // S
	CotPN_NguoiDuyet           = 19 // T (Mới)
	CotPN_NgayDuyet            = 20 // U (Mới)
	CotPN_NguoiCapNhat         = 21 // V
	CotPN_NgayCapNhat          = 22 // W
)

// --- CHI TIẾT PHIẾU NHẬP (15 Cột) ---
const (
	CotCTPN_MaPhieuNhap     = 0  // A
	CotCTPN_MaSanPham       = 1  // B
	CotCTPN_MaSKU           = 2  // C
	CotCTPN_MaNganhHang     = 3  // D
	CotCTPN_TenSanPham      = 4  // E
	CotCTPN_DonVi           = 5  // F
	CotCTPN_SoLuong         = 6  // G
	CotCTPN_DonGiaNhap      = 7  // H
	CotCTPN_VATPercent      = 8  // I
	CotCTPN_GiaSauVAT       = 9  // J
	CotCTPN_ChietKhauDong   = 10 // K
	CotCTPN_ThanhTienDong   = 11 // L
	CotCTPN_GiaVonThucTe    = 12 // M
	CotCTPN_BaoHanhThang    = 13 // N
	CotCTPN_GhiChuDong      = 14 // O
)

// --- SERIAL SẢN PHẨM (19 Cột) ---
const (
	CotSR_SerialIMEI             = 0  // A
	CotSR_MaSanPham              = 1  // B
	CotSR_MaSKU                  = 2  // C
	CotSR_MaNganhHang            = 3  // D
	CotSR_MaNhaCungCap           = 4  // E
	CotSR_MaPhieuNhap            = 5  // F
	CotSR_MaPhieuXuat            = 6  // G
	CotSR_TrangThai              = 7  // H
	CotSR_BaoHanhNhaCungCap      = 8  // I
	CotSR_HanBaoHanhNhaCungCap   = 9  // J
	CotSR_MaKhachHangHienTai     = 10 // K
	CotSR_NgayNhapKho            = 11 // L
	CotSR_NgayXuatKho            = 12 // M
	CotSR_GiaVonNhap             = 13 // N
	CotSR_KichHoatBaoHanhKhach   = 14 // O
	CotSR_HanBaoHanhKhach        = 15 // P
	CotSR_MaKho                  = 16 // Q
	CotSR_GhiChu                 = 17 // R
	CotSR_NgayCapNhat            = 18 // S
)

// ====================================================================
// 3. KHAI BÁO STRUCT GIAO TIẾP JSON
// ====================================================================

type PhieuNhap struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuNhap          string  `json:"ma_phieu_nhap"`
	MaNhaCungCap         string  `json:"ma_nha_cung_cap"`
	MaKho                string  `json:"ma_kho"`
	NgayNhap             string  `json:"ngay_nhap"`
	ChiTietJson          string  `json:"chi_tiet_json"` // Thêm mới
	TrangThai            int     `json:"trang_thai"`
	SoHoaDon             string  `json:"so_hoa_don"`
	NgayHoaDon           string  `json:"ngay_hoa_don"`
	UrlChungTu           string  `json:"url_chung_tu"`
	TongTienPhieu        float64 `json:"tong_tien_phieu"`
	GiamGiaPhieu         float64 `json:"giam_gia_phieu"`
	ChiPhiNhap           float64 `json:"chi_phi_nhap"` // Thêm mới
	DaThanhToan          float64 `json:"da_thanh_toan"`
	ConNo                float64 `json:"con_no"`
	PhuongThucThanhToan  string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiThanhToan   string  `json:"trang_thai_thanh_toan"`
	GhiChu               string  `json:"ghi_chu"`
	NguoiTao             string  `json:"nguoi_tao"`
	NgayTao              string  `json:"ngay_tao"`
	NguoiDuyet           string  `json:"nguoi_duyet"` // Thêm mới
	NgayDuyet            string  `json:"ngay_duyet"`  // Thêm mới
	NguoiCapNhat         string  `json:"nguoi_cap_nhat"`
	NgayCapNhat          string  `json:"ngay_cap_nhat"`
	
	// Thuộc tính này để chứa dữ liệu dạng Struct khi RAM bung JSON ra (Không lưu trực tiếp biến này xuống sheet Phiếu Nhập)
	ChiTiet              []*ChiTietPhieuNhap `json:"chi_tiet"`
}

type ChiTietPhieuNhap struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuNhap    string  `json:"ma_phieu_nhap"`
	MaSanPham      string  `json:"ma_san_pham"`
	MaSKU          string  `json:"ma_sku"`
	MaNganhHang    string  `json:"ma_nganh_hang"`
	TenSanPham     string  `json:"ten_san_pham"`
	DonVi          string  `json:"don_vi"`
	SoLuong        int     `json:"so_luong"`
	DonGiaNhap     float64 `json:"don_gia_nhap"`
	VATPercent     float64 `json:"vat_percent"`
	GiaSauVAT      float64 `json:"gia_sau_vat"`
	ChietKhauDong  float64 `json:"chiet_khau_dong"`
	ThanhTienDong  float64 `json:"thanh_tien_dong"`
	GiaVonThucTe   float64 `json:"gia_von_thuc_te"`
	BaoHanhThang   int     `json:"bao_hanh_thang"`
	GhiChuDong     string  `json:"ghi_chu_dong"`
}

type SerialSanPham struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	SerialIMEI               string  `json:"serial_imei"`
	MaSanPham                string  `json:"ma_san_pham"`
	MaSKU                    string  `json:"ma_sku"`
	MaNganhHang              string  `json:"ma_nganh_hang"`
	MaNhaCungCap             string  `json:"ma_nha_cung_cap"`
	MaPhieuNhap              string  `json:"ma_phieu_nhap"`
	MaPhieuXuat              string  `json:"ma_phieu_xuat"`
	TrangThai                int     `json:"trang_thai"` 
	BaoHanhNhaCungCap        int     `json:"bao_hanh_nha_cung_cap"` 
	HanBaoHanhNhaCungCap     string  `json:"han_bao_hanh_nha_cung_cap"`
	MaKhachHangHienTai       string  `json:"ma_khach_hang_hien_tai"`
	NgayNhapKho              string  `json:"ngay_nhap_kho"`
	NgayXuatKho              string  `json:"ngay_xuat_kho"`
	GiaVonNhap               float64 `json:"gia_von_nhap"`
	KichHoatBaoHanhKhach     string  `json:"kich_hoat_bao_hanh_khach"`
	HanBaoHanhKhach          string  `json:"han_bao_hanh_khach"`
	MaKho                    string  `json:"ma_kho"`
	GhiChu                   string  `json:"ghi_chu"`
	NgayCapNhat              string  `json:"ngay_cap_nhat"`
}

// ==============================================================================
// CẤU TRÚC: PHIẾU XUẤT BÁN HÀNG
// ==============================================================================
const (
	CotPX_MaPhieuXuat          = 0  
	CotPX_LoaiXuat             = 1  
	CotPX_NgayXuat             = 2  
	CotPX_MaKho                = 3  
	CotPX_MaKhachHang          = 4  
	CotPX_TrangThai            = 5  
	CotPX_MaVoucher            = 6  
	CotPX_TienGiamVoucher      = 7  
	CotPX_TongTienPhieu        = 8  
	CotPX_LinkChungTu          = 9  
	CotPX_DaThu                = 10 
	CotPX_ConNo                = 11 
	CotPX_PhuongThucThanhToan  = 12 
	CotPX_TrangThaiThanhToan   = 13 
	CotPX_PhiVanChuyen         = 14 
	CotPX_NguonDonHang         = 15 
	CotPX_ThongTinGiaoHang     = 16 
	CotPX_GhiChu               = 17 
	CotPX_NguoiTao             = 18 
	CotPX_NgayTao              = 19 
	CotPX_NgayCapNhat          = 20 

	CotCTPX_MaPhieuXuat        = 0  
	CotCTPX_MaSanPham          = 1  
	CotCTPX_MaSKU              = 2  
	CotCTPX_MaNganhHang        = 3  
	CotCTPX_TenSanPham         = 4  
	CotCTPX_DonVi              = 5  
	CotCTPX_SoLuong            = 6  
	CotCTPX_DonGiaBan          = 7  
	CotCTPX_VATPercent         = 8  
	CotCTPX_GiaSauVAT          = 9  
	CotCTPX_ChietKhauDong      = 10 
	CotCTPX_ThanhTienDong      = 11 
	CotCTPX_GiaVon             = 12 
	CotCTPX_BaoHanhThang       = 13 
	CotCTPX_GhiChuDong         = 14 
)

type PhieuXuat struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuXuat          string  `json:"ma_phieu_xuat"`
	LoaiXuat             string  `json:"loai_xuat"`
	NgayXuat             string  `json:"ngay_xuat"`
	MaKho                string  `json:"ma_kho"`
	MaKhachHang          string  `json:"ma_khach_hang"`
	TrangThai            int     `json:"trang_thai"`
	MaVoucher            string  `json:"ma_voucher"`
	TienGiamVoucher      float64 `json:"tien_giam_voucher"`
	TongTienPhieu        float64 `json:"tong_tien_phieu"`
	LinkChungTu          string  `json:"link_chung_tu"`
	DaThu                float64 `json:"da_thu"`
	ConNo                float64 `json:"con_no"`
	PhuongThucThanhToan  string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiThanhToan   string  `json:"trang_thai_thanh_toan"`
	PhiVanChuyen         float64 `json:"phi_van_chuyen"`
	NguonDonHang         string  `json:"nguon_don_hang"`
	ThongTinGiaoHang     string  `json:"thong_tin_giao_hang"`
	GhiChu               string  `json:"ghi_chu"`
	NguoiTao             string  `json:"nguoi_tao"`
	NgayTao              string  `json:"ngay_tao"`
	NgayCapNhat          string  `json:"ngay_cap_nhat"`

	ChiTiet              []*ChiTietPhieuXuat `json:"chi_tiet"`
}

type ChiTietPhieuXuat struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuXuat    string  `json:"ma_phieu_xuat"`
	MaSanPham      string  `json:"ma_san_pham"`
	MaSKU          string  `json:"ma_sku"`
	MaNganhHang    string  `json:"ma_nganh_hang"`
	TenSanPham     string  `json:"ten_san_pham"`
	DonVi          string  `json:"don_vi"`
	SoLuong        int     `json:"so_luong"`
	DonGiaBan      float64 `json:"don_gia_ban"`
	VATPercent     float64 `json:"vat_percent"`
	GiaSauVAT      float64 `json:"gia_sau_vat"`
	ChietKhauDong  float64 `json:"chiet_khau_dong"`
	ThanhTienDong  float64 `json:"thanh_tien_dong"`
	GiaVon         float64 `json:"gia_von"`
	BaoHanhThang   int     `json:"bao_hanh_thang"`
	GhiChuDong     string  `json:"ghi_chu_dong"`
}

// ==============================================================================
// CẤU TRÚC: GÓI DỊCH VỤ SAAS
// ==============================================================================
const (
	DongBatDau_GoiDichVu = 11

	CotGDV_MaGoi              = 0  // A
	CotGDV_TenGoi             = 1  // B
	CotGDV_LoaiGoi            = 2  // C
	CotGDV_ThoiHanNgay        = 3  // D
	CotGDV_ThoiHanHienThi     = 4  // E 
	CotGDV_NhanHienThi        = 5  // F
	CotGDV_GiaNiemYet         = 6  // G
	CotGDV_GiaBan             = 7  // H
	CotGDV_MaCodeKichHoatJson = 8  // I
	CotGDV_GioiHanJson        = 9  // J
	CotGDV_MoTa               = 10 // K
	CotGDV_NgayBatDau         = 11 // L
	CotGDV_NgayKetThuc        = 12 // M
	CotGDV_SoLuongConLai      = 13 // N
	CotGDV_TrangThai          = 14 // O
)

// Cấu trúc để bóc tách JSON Cột I (Mã Code Khuyến Mãi)
type CodeKichHoat struct {
	Code     string  `json:"code"`
	GiamTien float64 `json:"giam_tien"`
	SoLuong  int     `json:"so_luong"` // -1 là không giới hạn
}

// Cấu trúc dữ liệu chính của Gói Dịch Vụ
type GoiDichVu struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaGoi              string          `json:"ma_goi"`
	TenGoi             string          `json:"ten_goi"`
	LoaiGoi            string          `json:"loai_goi"`
	ThoiHanNgay        int             `json:"thoi_han_ngay"`
	ThoiHanHienThi     string          `json:"thoi_han_hien_thi"` 
	NhanHienThi        string          `json:"nhan_hien_thi"`
	GiaNiemYet         float64         `json:"gia_niem_yet"`
	GiaBan             float64         `json:"gia_ban"`
	MaCodeKichHoatJson string          `json:"-"` 
	DanhSachCode       []CodeKichHoat  `json:"danh_sach_code"` 
	GioiHanJson        string          `json:"gioi_han_json"`
	MoTa               string          `json:"mo_ta"`
	NgayBatDau         string          `json:"ngay_bat_dau"`
	NgayKetThuc        string          `json:"ngay_ket_thuc"`
	SoLuongConLai      int             `json:"so_luong_con_lai"`
	TrangThai          int             `json:"trang_thai"`
}

// ==============================================================================
// CẤU TRÚC TẦNG 3: NHÂN VIÊN & KHÁCH LẺ CỦA CỬA HÀNG (BẢN RÚT GỌN)
// Tên Sheet: THANH_VIEN_SHOP
// ==============================================================================
const (
	DongBatDau_ThanhVienShop = 11

	CotTVS_MaThanhVien    = 0  // A
	CotTVS_TenDangNhap    = 1  // B
	CotTVS_MatKhauHash    = 2  // C
	CotTVS_MaPinHash      = 3  // D
	CotTVS_HoTen          = 4  // E
	CotTVS_DienThoai      = 5  // F
	CotTVS_Email          = 6  // G
	CotTVS_VaiTro         = 7  // H (VD: ban_hang, ke_toan, khach_le)
	CotTVS_TrangThai      = 8  // I (1: Hoạt động, 0: Khóa)
	CotTVS_DiemThuong     = 9  // J
	CotTVS_GhiChu         = 10 // K
	CotTVS_NgayTao        = 11 // L
)

type ThanhVienShop struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaThanhVien string `json:"ma_thanh_vien"`
	TenDangNhap string `json:"ten_dang_nhap"`
	MatKhauHash string `json:"-"`
	MaPinHash   string `json:"-"`
	HoTen       string `json:"ho_ten"`
	DienThoai   string `json:"dien_thoai"`
	Email       string `json:"email"`
	VaiTro      string `json:"vai_tro"`
	TrangThai   int    `json:"trang_thai"`
	DiemThuong  int    `json:"diem_thuong"`
	GhiChu      string `json:"ghi_chu"`
	NgayTao     string `json:"ngay_tao"`
}
