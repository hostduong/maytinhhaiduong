package core

// ==============================================================================
// 1. ĐỊNH NGHĨA TÊN SHEET CHUẨN VÀ PREFIX (KEY-VALUE NOSQL)
// ==============================================================================

const (
	TenSheetKhachHangMaster  = "KHACH_HANG_MASTER"
	TenSheetPhanQuyenMaster  = "PHAN_QUYEN_MASTER"
	TenSheetCauHinhMaster    = "CAU_HINH_MASTER"
	TenSheetTinNhanMaster    = "TIN_NHAN_MASTER"
	TenSheetCauHinhThuocTinh = "CAU_HINH_THUOC_TINH" 
)

// BỘ PREFIX ĐỊNH TUYẾN DỮ LIỆU
const (
	PreGoiDichVu = "GoiDichVu__"
	PrePhanQuyen = "PhanQuyen__"
	PreHeThong   = "HeThong__"
)

const (
	TenSheetKhachHangAdmin   = "KHACH_HANG_ADMIN"
	TenSheetPhanQuyenAdmin   = "PHAN_QUYEN_ADMIN"
)

const (
	TenSheetKhachHang        = "KHACH_HANG"
	TenSheetPhanQuyen        = "PHAN_QUYEN"
	TenSheetDanhMuc          = "DANH_MUC"
	TenSheetThuongHieu       = "THUONG_HIEU"
	TenSheetBienLoiNhuan     = "BIEN_LOI_NHUAN"
	TenSheetNhaCungCap       = "NHA_CUNG_CAP"
	TenSheetSerial           = "SERIAL_SAN_PHAM"
	TenSheetPhieuNhap        = "PHIEU_NHAP"
	TenSheetChiTietPhieuNhap = "CHI_TIET_PHIEU_NHAP"
	TenSheetPhieuXuat        = "PHIEU_XUAT"
	TenSheetChiTietPhieuXuat = "CHI_TIET_PHIEU_XUAT"
	TenSheetHoaDon           = "HOA_DON"
	TenSheetHoaDonChiTiet    = "HOA_DON_CHI_TIET"
	TenSheetPhieuThuChi      = "PHIEU_THU_CHI"
	TenSheetPhieuBaoHanh     = "PHIEU_BAO_HANH"
)

// ==============================================================================
// 2. CẤU TRÚC PHÂN QUYỀN
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
	StyleLevel int    `json:"style_level"` 
	StyleTheme int    `json:"style_theme"` 
}

// ==============================================================================
// 3. CẤU TRÚC KHÁCH HÀNG SAAS (MASTER & ADMIN) - NOSQL 2 CỘT
// ==============================================================================
const (
	DongBatDau_KhachHang = 11
	CotKH_MaKhachHang    = 0 // Cột A
	CotKH_DataJSON       = 1 // Cột B
)

type TenantBaoMat struct {
	MatKhauHash string `json:"mat_khau_hash"`
	MaPinHash   string `json:"ma_pin_hash"`
}

type TenantDeviceToken struct {
	DeviceID string `json:"device_id"`
	Dev      string `json:"dev"`
	Exp      int64  `json:"exp"`
	Created  int64  `json:"created"`
}

type TenantGoiDichVu struct {
	MaGoi       string   `json:"ma_goi"`
	TenGoi      string   `json:"ten_goi"`
	LoaiGoi     string   `json:"loai_goi"`
	TrangThai   string   `json:"trang_thai"`
	NgayHetHan  int64    `json:"ngay_het_han"`
	MaxSanPham  int      `json:"max_san_pham"`
	MaxNhanVien int      `json:"max_nhan_vien"`
	Modules     []string `json:"modules"`
}

type TenantDomain struct {
	CapTenMien   bool   `json:"cap_ten_mien"`
	CustomDomain string `json:"custom_domain"`
	Subdomain    string `json:"subdomain"`
}

type TenantCauHinh struct {
	Theme       string `json:"theme"`
	ChuyenNganh string `json:"chuyen_nganh"`
	Lang        string `json:"lang"`
	DarkMode    bool   `json:"dark_mode"`
}

type TenantSystem struct {
	SheetID        string `json:"sheet_id"`
	GoogleAuthJson string `json:"google_auth_json"`
	FolderDriveID  string `json:"folder_drive_id"`
}

type TenantHoaDonConfig struct {
	Enabled  bool   `json:"enabled"`
	Provider string `json:"provider"`
	Serial   string `json:"serial"`
	MauSo    string `json:"mau_so"`
	ChuKySo  bool   `json:"chu_ky_so"`
	TokenAPI string `json:"token_api"`
}

type TenantThongTin struct {
	NguonKhachHang string `json:"nguon_khach_hang"`
	TenCuaHang     string `json:"ten_cua_hang"`
	TenKhachHang   string `json:"ten_khach_hang"`
	DienThoai      string `json:"dien_thoai"`
	Zalo           string `json:"zalo"`
	AnhDaiDien     string `json:"anh_dai_dien"`
	DiaChi         string `json:"dia_chi"`
	NgaySinh       string `json:"ngay_sinh"`
	GioiTinh       int    `json:"gioi_tinh"`
	DiemTichLuy    int    `json:"diem_tich_luy"`
	MaSoThue       string `json:"ma_so_thue"`
}

type TenantNganHang struct {
	TenNganHang string `json:"ten_ngan_hang"`
	SoTaiKhoan  string `json:"so_tai_khoan"`
	ChuTaiKhoan string `json:"chu_tai_khoan"`
}

type TenantViTien struct {
	SoDu    float64 `json:"so_du"`
	DaTieu  float64 `json:"da_tieu"`
	TongNap float64 `json:"tong_nap"`
}

type KhachHang struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	Version        int    `json:"version"`
	MaKhachHang    string `json:"ma_khach_hang"`
	TenDangNhap    string `json:"ten_dang_nhap"`
	Email          string `json:"email"`

	BaoMat         TenantBaoMat                 `json:"bao_mat"`
	RefreshTokens  map[string]TenantDeviceToken `json:"refresh_tokens"`

	VaiTroQuyenHan string `json:"vai_tro_quyen_han"`
	ChucVu         string `json:"chuc_vu"`
	TrangThai      int    `json:"trang_thai"`

	GoiDichVu      []TenantGoiDichVu            `json:"goi_dich_vu"`
	Modules        map[string]bool              `json:"modules"`
	Domain         TenantDomain                 `json:"domain"`
	CauHinh        TenantCauHinh                `json:"cau_hinh"`
	System         TenantSystem                 `json:"system"`
	HoaDonConfig   TenantHoaDonConfig           `json:"hoa_don_config"`
	ThongTin       TenantThongTin               `json:"thong_tin"`
	NganHang       TenantNganHang               `json:"ngan_hang"`
	MangXaHoi      map[string]string            `json:"mang_xa_hoi"`
	ViTien         TenantViTien                 `json:"vi_tien"`
	KetNoiDoiTac   map[string]interface{}       `json:"ket_noi_doi_tac"`

	GhiChu         string `json:"ghi_chu"`
	NguoiTao       string `json:"nguoi_tao"`
	NguoiCapNhat   string `json:"nguoi_cap_nhat"`
	NgayTao        int64  `json:"ngay_tao"`
	NgayCapNhat    int64  `json:"ngay_cap_nhat"`

	// Trường phụ chạy RAM
	Inbox      []*TinNhan `json:"-"`
	StyleLevel int        `json:"-"`
	StyleTheme int        `json:"-"`
}

// ==============================================================================
// 4. CẤU TRÚC KHÁCH HÀNG TẦNG CỬA HÀNG (NOSQL 2 CỘT)
// ==============================================================================
const (
	DongBatDau_KhachHangShop = 11
	CotKHS_MaKhachHang    = 0 
	CotKHS_DataJSON       = 1 
)

type KhachHangShop struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	Version        int    `json:"version"`
	MaKhachHang    string `json:"ma_khach_hang"`
	TenDangNhap    string `json:"ten_dang_nhap"`
	Email          string `json:"email"`

	BaoMat         TenantBaoMat                 `json:"bao_mat"`
	RefreshTokens  map[string]TenantDeviceToken `json:"refresh_tokens"`

	VaiTroQuyenHan string `json:"vai_tro_quyen_han"`
	ChucVu         string `json:"chuc_vu"`
	TrangThai      int    `json:"trang_thai"`

	ThongTin       TenantThongTin               `json:"thong_tin"`
	NganHang       TenantNganHang               `json:"ngan_hang"`
	MangXaHoi      map[string]string            `json:"mang_xa_hoi"`
	ViTien         TenantViTien                 `json:"vi_tien"`

	GhiChu         string `json:"ghi_chu"`
	NguoiTao       string `json:"nguoi_tao"`
	NguoiCapNhat   string `json:"nguoi_cap_nhat"`
	NgayTao        int64  `json:"ngay_tao"`
	NgayCapNhat    int64  `json:"ngay_cap_nhat"`
}

// ==============================================================================
// 5. CẤU TRÚC TIN NHẮN (NOSQL 2 CỘT)
// ==============================================================================
const (
	DongBatDau_TinNhan = 11
	CotTN_MaTinNhan    = 0 // Cột A
	CotTN_DataJSON     = 1 // Cột B
)

type FileDinhKem struct {
	TenFile string `json:"name"`
	URL     string `json:"url"`
	Loai    string `json:"type"` 
}

type TinNhan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaTinNhan    string        `json:"ma_tin_nhan"`
	LoaiTinNhan  string        `json:"loai_tin_nhan"`
	NguoiGuiID   string        `json:"nguoi_gui_id"`
	NguoiNhanID  []string      `json:"nguoi_nhan_id"` // Mảng ["ALL"] hoặc cụ thể
	TieuDe       string        `json:"tieu_de"`
	NoiDung      string        `json:"noi_dung"`
	DinhKem      []FileDinhKem `json:"dinh_kem"`
	ThamChieuID  []string      `json:"tham_chieu_id"`
	ReplyChoID   string        `json:"reply_cho_id"`
	NguoiDoc     []string      `json:"nguoi_doc"`
	TrangThaiXoa []string      `json:"trang_thai_xoa"`
	NgayTao      int64         `json:"ngay_tao"`

	DaDoc bool `json:"-"` // Cờ logic trên RAM
}

// ==============================================================================
// CÁC CẤU TRÚC VẬN HÀNH KINH DOANH CHUNG
// ==============================================================================

const (
	DongBatDau_CauHinhThuocTinh = 11
	CotCHTT_MaThuocTinh  = 0 
	CotCHTT_TenThuocTinh = 1 
	CotCHTT_KieuNhap     = 2 
	CotCHTT_DonVi        = 3 
	CotCHTT_StartNganh   = 4 
)

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

const (
	DongBatDau_Product = 11
	CotProd_MaSanPham  = 0
	CotProd_DataJSON   = 1
)

type ProductRow struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaSanPham      string `json:"ma_san_pham"`
	DataJSON       string `json:"data_json"`
}

type ProductJSON struct {
	SpreadsheetID  string        `json:"-"`
	DongTrongSheet int           `json:"-"`          
	MaSanPham    string        `json:"ma_san_pham"`
	Version      int           `json:"version"`      
	CreatedAt    int64         `json:"created_at"`   
	UpdatedAt    int64         `json:"updated_at"`   
	TrangThai    int           `json:"trang_thai"`   
	MaNganh      string        `json:"ma_nganh"`
	MaDanhMuc    []string      `json:"ma_danh_muc"`
	MaThuongHieu string        `json:"ma_thuong_hieu"`
	TenSanPham   string        `json:"ten_san_pham"`
	TenRutGon    string        `json:"ten_rut_gon"`
	Slug         string        `json:"slug"`
	NenTang      []string      `json:"nen_tang"`     
	LoaiSanPham  string        `json:"loai_san_pham"`
	SearchText   string        `json:"search_text"`  
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
	TenSanPham   string                 `json:"ten_san_pham"` 
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
	ThuocTinh    map[string]interface{} `json:"thuoc_tinh"` 
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
	NgayTao      string `json:"ngay_tao"` 
	NguoiCapNhat string `json:"nguoi_cap_nhat"`
	NgayCapNhat  string `json:"ngay_cap_nhat"`
}

const (
	DongBatDau_HoaDon = 11
	CotHD_MaHoaDon           = 0  
	CotHD_MaTraCuu           = 1  
	CotHD_XmlUrl             = 2  
	CotHD_LoaiHoaDon         = 3  
	CotHD_MaPhieuXuat        = 4  
	CotHD_MaKhachHang        = 5  
	CotHD_NgayHoaDon         = 6  
	CotHD_KyHieu             = 7  
	CotHD_SoHoaDon           = 8  
	CotHD_MauSo              = 9  
	CotHD_LinkChungTu        = 10 
	CotHD_TongTienPhieu      = 11 
	CotHD_TongVAT            = 12 
	CotHD_TongTienSauVAT     = 13 
	CotHD_TrangThai          = 14 
	CotHD_TrangThaiThanhToan = 15 
	CotHD_GhiChu             = 16 
	CotHD_NguoiTao           = 17 
	CotHD_NgayTao            = 18 
	CotHD_NgayCapNhat        = 19 
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

const (
	CotHDCT_MaHoaDon    = 0  
	CotHDCT_MaSanPham   = 1  
	CotHDCT_MaSKU       = 2  
	CotHDCT_MaNganhHang = 3  
	CotHDCT_TenSanPham  = 4  
	CotHDCT_DonVi       = 5  
	CotHDCT_SoLuong     = 6  
	CotHDCT_DonGiaBan   = 7  
	CotHDCT_VATPercent  = 8  
	CotHDCT_TienVAT     = 9  
	CotHDCT_ThanhTien   = 10 
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

const (
	DongBatDau_PhieuThuChi = 11
	CotPTC_MaPhieu             = 0  
	CotPTC_NgayTaoPhieu        = 1  
	CotPTC_LoaiPhieu           = 2  
	CotPTC_DoiTuongLoai        = 3  
	CotPTC_DoiTuongID          = 4  
	CotPTC_HangMucThuChi       = 5  
	CotPTC_CoHoaDonDo          = 6  
	CotPTC_MaChungTuThamChieu  = 7  
	CotPTC_SoTien              = 8  
	CotPTC_PhuongThucThanhToan = 9  
	CotPTC_TrangThaiDuyet      = 10 
	CotPTC_NguoiDuyet          = 11 
	CotPTC_GhiChu              = 12 
	CotPTC_NguoiTao            = 13 
	CotPTC_NgayTao             = 14 
	CotPTC_NgayCapNhat         = 15 
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

const (
	DongBatDau_PhieuBaoHanh = 11
	CotPBH_MaPhieu           = 0  
	CotPBH_LoaiPhieu         = 1  
	CotPBH_SerialIMEI        = 2  
	CotPBH_MaSanPham         = 3  
	CotPBH_MaKhachHang       = 4  
	CotPBH_TenNguoiGui       = 5  
	CotPBH_SDTNguoiGui       = 6  
	CotPBH_NgayNhan          = 7  
	CotPBH_TinhTrangLoi      = 8  
	CotPBH_HinhThuc          = 9  
	CotPBH_TrangThai         = 10 
	CotPBH_NgayTraDuKien     = 11 
	CotPBH_NgayTraThucTe     = 12 
	CotPBH_ChiPhiSua         = 13 
	CotPBH_PhiThuKhach       = 14 
	CotPBH_KetQuaSuaChua     = 15 
	CotPBH_LinhKienThayThe   = 16 
	CotPBH_MaNhanVienKyThuat = 17 
	CotPBH_GhiChu            = 18 
	CotPBH_NguoiTao          = 19 
	CotPBH_NgayTao           = 20 
	CotPBH_NgayCapNhat       = 21 
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
	CotNCC_AnhDaiDien         = 8  
	CotNCC_NganHang           = 9  
	CotNCC_NhomNhaCungCap     = 10 
	CotNCC_LoaiNhaCungCap     = 11 
	CotNCC_DieuKhoanThanhToan = 12 
	CotNCC_ChietKhauMacDinh   = 13 
	CotNCC_HanMucCongNo       = 14 
	CotNCC_CongNoDauKy        = 15 
	CotNCC_TongMua            = 16 
	CotNCC_NoCanTra           = 17 
	CotNCC_ThongTinThemJson   = 18 
	CotNCC_TrangThai          = 19 
	CotNCC_GhiChu             = 20 
	CotNCC_NguoiTao           = 21 
	CotNCC_NgayTao            = 22 
	CotNCC_NgayCapNhat        = 23 
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

const (
	DongBatDau_PhieuNhap       = 11
	CotPN_MaPhieuNhap          = 0  
	CotPN_MaNhaCungCap         = 1  
	CotPN_MaKho                = 2  
	CotPN_NgayNhap             = 3  
	CotPN_ChiTietJson          = 4  
	CotPN_TrangThai            = 5  
	CotPN_SoHoaDon             = 6  
	CotPN_NgayHoaDon           = 7  
	CotPN_UrlChungTu           = 8  
	CotPN_TongTienPhieu        = 9  
	CotPN_GiamGiaPhieu         = 10 
	CotPN_ChiPhiNhap           = 11 
	CotPN_DaThanhToan          = 12 
	CotPN_ConNo                = 13 
	CotPN_PhuongThucThanhToan  = 14 
	CotPN_TrangThaiThanhToan   = 15 
	CotPN_GhiChu               = 16 
	CotPN_NguoiTao             = 17 
	CotPN_NgayTao              = 18 
	CotPN_NguoiDuyet           = 19 
	CotPN_NgayDuyet            = 20 
	CotPN_NguoiCapNhat         = 21 
	CotPN_NgayCapNhat          = 22 
)

const (
	CotCTPN_MaPhieuNhap     = 0  
	CotCTPN_MaSanPham       = 1  
	CotCTPN_MaSKU           = 2  
	CotCTPN_MaNganhHang     = 3  
	CotCTPN_TenSanPham      = 4  
	CotCTPN_DonVi           = 5  
	CotCTPN_SoLuong         = 6  
	CotCTPN_DonGiaNhap      = 7  
	CotCTPN_VATPercent      = 8  
	CotCTPN_GiaSauVAT       = 9  
	CotCTPN_ChietKhauDong   = 10 
	CotCTPN_ThanhTienDong   = 11 
	CotCTPN_GiaVonThucTe    = 12 
	CotCTPN_BaoHanhThang    = 13 
	CotCTPN_GhiChuDong      = 14 
)

type PhieuNhap struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuNhap          string  `json:"ma_phieu_nhap"`
	MaNhaCungCap         string  `json:"ma_nha_cung_cap"`
	MaKho                string  `json:"ma_kho"`
	NgayNhap             string  `json:"ngay_nhap"`
	ChiTietJson          string  `json:"chi_tiet_json"` 
	TrangThai            int     `json:"trang_thai"`
	SoHoaDon             string  `json:"so_hoa_don"`
	NgayHoaDon           string  `json:"ngay_hoa_don"`
	UrlChungTu           string  `json:"url_chung_tu"`
	TongTienPhieu        float64 `json:"tong_tien_phieu"`
	GiamGiaPhieu         float64 `json:"giam_gia_phieu"`
	ChiPhiNhap           float64 `json:"chi_phi_nhap"` 
	DaThanhToan          float64 `json:"da_thanh_toan"`
	ConNo                float64 `json:"con_no"`
	PhuongThucThanhToan  string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiThanhToan   string  `json:"trang_thai_thanh_toan"`
	GhiChu               string  `json:"ghi_chu"`
	NguoiTao             string  `json:"nguoi_tao"`
	NgayTao              string  `json:"ngay_tao"`
	NguoiDuyet           string  `json:"nguoi_duyet"` 
	NgayDuyet            string  `json:"ngay_duyet"`  
	NguoiCapNhat         string  `json:"nguoi_cap_nhat"`
	NgayCapNhat          string  `json:"ngay_cap_nhat"`
	
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

const (
	CotSR_SerialIMEI             = 0  
	CotSR_MaSanPham              = 1  
	CotSR_MaSKU                  = 2  
	CotSR_MaNganhHang            = 3  
	CotSR_MaNhaCungCap           = 4  
	CotSR_MaPhieuNhap            = 5  
	CotSR_MaPhieuXuat            = 6  
	CotSR_TrangThai              = 7  
	CotSR_BaoHanhNhaCungCap      = 8  
	CotSR_HanBaoHanhNhaCungCap   = 9  
	CotSR_MaKhachHangHienTai     = 10 
	CotSR_NgayNhapKho            = 11 
	CotSR_NgayXuatKho            = 12 
	CotSR_GiaVonNhap             = 13 
	CotSR_KichHoatBaoHanhKhach   = 14 
	CotSR_HanBaoHanhKhach        = 15 
	CotSR_MaKho                  = 16 
	CotSR_GhiChu                 = 17 
	CotSR_NgayCapNhat            = 18 
)

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

const (
	DongBatDau_ThanhVienShop = 11

	CotTVS_MaThanhVien    = 0  
	CotTVS_TenDangNhap    = 1  
	CotTVS_MatKhauHash    = 2  
	CotTVS_MaPinHash      = 3  
	CotTVS_HoTen          = 4  
	CotTVS_DienThoai      = 5  
	CotTVS_Email          = 6  
	CotTVS_VaiTro         = 7  
	CotTVS_TrangThai      = 8  
	CotTVS_DiemThuong     = 9  
	CotTVS_GhiChu         = 10 
	CotTVS_NgayTao        = 11 
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


// ==============================================================================
// ĐỊNH NGHĨA BẢNG CẤU HÌNH (NOSQL 2 CỘT)
// ==============================================================================
const (
	DongBatDau_CauHinh = 11

	CotCH_MaCauHinh = 0  // Cột A (GoiDichVu__PRO, PhanQuyen__ADMIN)
	CotCH_DataJSON  = 1  // Cột B (Toàn bộ Object JSON)
)

type CodeKichHoat struct {
	Code           string  `json:"code"`
	Loai           string  `json:"loai"`
	GiamTien       float64 `json:"giam_tien"`
	GiamPhanTram   float64 `json:"giam_phan_tram"`
	SoLuong        int     `json:"so_luong"`
	DaSuDung       int     `json:"da_su_dung"`
	NguoiGioiThieu string  `json:"nguoi_gioi_thieu"`
	NgayHetHan     int64   `json:"ngay_het_han"`
	TrangThai      int     `json:"trang_thai"`
}

type GoiDichVuGia struct {
	GiaNiemYet   float64 `json:"gia_niem_yet"`
	PhanTramGiam float64 `json:"phan_tram_giam"`
	GiaBan       float64 `json:"gia_ban"`
	TienTe       string  `json:"tien_te"`
}

type GoiDichVuGioiHan struct {
	MaxSanPham        int `json:"max_san_pham"`
	MaxNhanVien       int `json:"max_nhan_vien"`
	MaxHoaDonThang    int `json:"max_hoa_don_thang"`
	MaxKhachHang      int `json:"max_khach_hang"`
	MaxApiRequestNgay int `json:"max_api_request_ngay"`
}

type GoiDichVu struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	Version          int              `json:"version"`
	CreatedAt        int64            `json:"created_at"`
	UpdatedAt        int64            `json:"updated_at"`

	MaGoi            string           `json:"ma_goi"`
	SlugGoi          string           `json:"slug_goi"`
	MaGoiHienThi     string           `json:"ma_goi_hien_thi"`
	TenGoi           string           `json:"ten_goi"`
	LoaiGoi          string           `json:"loai_goi"`
	XepHang          int              `json:"xep_hang"`
	PhuDe            string           `json:"phu_de"`
	Badge            string           `json:"badge"`
	NhanHienThi      string           `json:"nhan_hien_thi"`
	MoTa             string           `json:"mo_ta"`
	MoTaGoi          []string         `json:"mo_ta_goi"`
	
	ThoiHanNgay      int              `json:"thoi_han_ngay"`
	ThoiHanHienThi   string           `json:"thoi_han_hien_thi"`
	KhuVucHienThi    []string         `json:"khu_vuc_hien_thi"`
	
	Gia              GoiDichVuGia     `json:"gia"`
	CodeKichHoat     []CodeKichHoat   `json:"code_kich_hoat"`
	GioiHan          GoiDichVuGioiHan `json:"gioi_han"`
	Modules          map[string]bool  `json:"modules"`
	
	NgayBatDau       int64            `json:"ngay_bat_dau"`
	NgayKetThuc      int64            `json:"ngay_ket_thuc"`
	SoLuongConLai    int              `json:"so_luong_con_lai"`
	TrangThai        int              `json:"trang_thai"`
}

// ==============================================================================
// CẤU TRÚC PHÂN QUYỀN (NOSQL 2 CỘT - LƯU CHUNG BẢNG CẤU HÌNH)
// ==============================================================================
type PhanQuyen struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	Version   int      `json:"version"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`

	MaVaiTro  string   `json:"ma_vai_tro"`
	TenVaiTro string   `json:"ten_vai_tro"`
	Level     int      `json:"level"`
	MoTa      string   `json:"mo_ta"`
	QuyenHan  []string `json:"quyen_han"`  // Mảng chứa: ["nhap_kho.tao", "hoa_don.doi_gia"]
	TrangThai int      `json:"trang_thai"` // 1: Hoạt động, 0: Tạm khóa
	IsLocked  bool     `json:"is_locked"`  // [MỚI] Cờ khóa cứng từ Master
}

// BẢNG MASTER QUYỀN HẠN (White-list: Dùng để filter data từ trình duyệt gửi lên)
var DanhSachQuyenHanChuan = map[string]bool{
	// Bán Hàng & Hóa Đơn
	"ban_hang_pos.xem": true, "ban_hang_pos.tao": true, "ban_hang_pos.sua": true, "ban_hang_pos.xoa": true, "ban_hang_pos.in": true, "ban_hang_pos.huy": true, "ban_hang_pos.doi_gia": true, "ban_hang_pos.void": true,
	"hoa_don.xem": true, "hoa_don.tao": true, "hoa_don.sua": true, "hoa_don.xoa": true, "hoa_don.duyet": true, "hoa_don.huy": true, "hoa_don.in": true, "hoa_don.export": true, "hoa_don.copy": true, "hoa_don.doi_gia": true, "hoa_don.khoa": true, "hoa_don.chuyen_chi_nhanh": true, "hoa_don.xac_nhan_thanh_toan": true,
	"tra_hang.xem": true, "tra_hang.tao": true, "tra_hang.duyet": true, "tra_hang.huy": true, "tra_hang.in": true,
	"phong_ban.xem": true, "phong_ban.tao": true, "phong_ban.sua": true, "phong_ban.xoa": true,
	"bao_hanh.xem": true, "bao_hanh.tao": true, "bao_hanh.sua": true, "bao_hanh.xoa": true, "bao_hanh.huy": true, "bao_hanh.in": true,
	
	// Sửa Chữa
	"phieu_sua_chua.xem": true, "phieu_sua_chua.tao": true, "phieu_sua_chua.sua": true, "phieu_sua_chua.xoa": true, "phieu_sua_chua.duyet": true, "phieu_sua_chua.huy": true, "phieu_sua_chua.in": true, "phieu_sua_chua.phan_cong": true,
	"bao_gia_sua_chua.xem": true, "bao_gia_sua_chua.tao": true, "bao_gia_sua_chua.sua": true, "bao_gia_sua_chua.duyet": true, "bao_gia_sua_chua.huy": true, "bao_gia_sua_chua.in": true,
	"tinh_trang_sua.xem": true, "tinh_trang_sua.sua": true,
	"linh_kien_thay_the.xem": true, "linh_kien_thay_the.tao": true, "linh_kien_thay_the.sua": true, "linh_kien_thay_the.xoa": true,
	"lich_su_sua.xem": true, "lich_su_sua.export": true,
	"hoa_hong_sua_chua.xem": true, "hoa_hong_sua_chua.sua": true,

	// Sản Phẩm (Đã bổ sung: Biên Lợi Nhuận & Thuộc Tính)
	"san_pham.xem": true, "san_pham.tao": true, "san_pham.sua": true, "san_pham.xoa": true, "san_pham.import": true, "san_pham.export": true, "san_pham.copy": true, "san_pham.khoa": true, "san_pham.xem_gia_von": true, "san_pham.xem_gia_ban": true, "san_pham.doi_gia": true, "san_pham.in": true, "san_pham.in_tem": true,
	"danh_muc.xem": true, "danh_muc.tao": true, "danh_muc.sua": true, "danh_muc.xoa": true,
	"thuong_hieu.xem": true, "thuong_hieu.tao": true, "thuong_hieu.sua": true, "thuong_hieu.xoa": true,
	"khuyen_mai.xem": true, "khuyen_mai.tao": true, "khuyen_mai.sua": true, "khuyen_mai.xoa": true, "khuyen_mai.huy": true, "khuyen_mai.copy": true, "khuyen_mai.duyet": true,
	"serial_number.xem": true, "serial_number.tao": true, "serial_number.sua": true, "serial_number.xoa": true, "serial_number.import": true, "serial_number.in_tem": true,
	"bien_loi_nhuan.xem": true, "bien_loi_nhuan.tao": true, "bien_loi_nhuan.sua": true, "bien_loi_nhuan.xoa": true,
	"thuoc_tinh_san_pham.xem": true, "thuoc_tinh_san_pham.tao": true, "thuoc_tinh_san_pham.sua": true, "thuoc_tinh_san_pham.xoa": true,

	// Kho Hàng
	"kho_hang.xem": true, "kho_hang.tao": true, "kho_hang.sua": true, "kho_hang.xoa": true,
	"nhap_kho.xem": true, "nhap_kho.tao": true, "nhap_kho.sua": true, "nhap_kho.duyet": true, "nhap_kho.duyet_lai": true, "nhap_kho.huy": true, "nhap_kho.in": true, "nhap_kho.export": true, "nhap_kho.xem_gia_nhap": true, "nhap_kho.sua_gia_nhap": true,
	"xuat_kho.xem": true, "xuat_kho.tao": true, "xuat_kho.sua": true, "xuat_kho.duyet": true, "xuat_kho.duyet_lai": true, "xuat_kho.huy": true, "xuat_kho.in": true, "xuat_kho.export": true,
	"kiem_kho.xem": true, "kiem_kho.tao": true, "kiem_kho.sua": true, "kiem_kho.duyet": true, "kiem_kho.huy": true, "kiem_kho.export": true, "kiem_kho.import": true,
	"dieu_chuyen_kho.xem": true, "dieu_chuyen_kho.tao": true, "dieu_chuyen_kho.duyet": true, "dieu_chuyen_kho.nhan_hang": true, "dieu_chuyen_kho.huy": true, "dieu_chuyen_kho.export": true,
	"tra_hang_nha_cung_cap.xem": true, "tra_hang_nha_cung_cap.tao": true, "tra_hang_nha_cung_cap.duyet": true, "tra_hang_nha_cung_cap.huy": true, "tra_hang_nha_cung_cap.export": true,
	"canh_bao_ton_kho.xem": true, "canh_bao_ton_kho.cau_hinh": true,

	// Khách Hàng & Đối Tác
	"khach_hang.xem": true, "khach_hang.tao": true, "khach_hang.sua": true, "khach_hang.xoa": true, "khach_hang.export": true, "khach_hang.import": true, "khach_hang.phan_cong": true,
	"tich_diem.xem": true, "tich_diem.sua": true,
	"hang_thanh_vien.xem": true, "hang_thanh_vien.tao": true, "hang_thanh_vien.sua": true, "hang_thanh_vien.xoa": true,
	"cham_soc_khach_hang.xem": true, "cham_soc_khach_hang.tao": true,
	"lich_su_mua_hang.xem": true, "lich_su_mua_hang.export": true,
	"nha_cung_cap.xem": true, "nha_cung_cap.tao": true, "nha_cung_cap.sua": true, "nha_cung_cap.xoa": true, "nha_cung_cap.import": true, "nha_cung_cap.export": true,

	// Tài Chính
	"thu_chi.xem": true, "thu_chi.tao": true, "thu_chi.sua": true, "thu_chi.xoa": true, "thu_chi.duyet": true, "thu_chi.duyet_lai": true, "thu_chi.huy": true, "thu_chi.in": true, "thu_chi.export": true,
	"so_quy.xem": true, "so_quy.chot_ca": true, "so_quy.export": true,
	"cong_no_khach_hang.xem": true, "cong_no_khach_hang.sua": true, "cong_no_khach_hang.thanh_toan_no": true, "cong_no_khach_hang.export": true,
	"cong_no_nha_cung_cap.xem": true, "cong_no_nha_cung_cap.sua": true, "cong_no_nha_cung_cap.thanh_toan_no": true, "cong_no_nha_cung_cap.export": true,
	"quyet_toan_ca.xem": true, "quyet_toan_ca.tao": true, "quyet_toan_ca.duyet": true, "quyet_toan_ca.chot_ca": true, "quyet_toan_ca.export": true,
	"luong_nhan_vien.xem": true, "luong_nhan_vien.tao": true, "luong_nhan_vien.sua": true, "luong_nhan_vien.duyet": true, "luong_nhan_vien.chot_luong": true, "luong_nhan_vien.export": true,
	"chi_phi_hoat_dong.xem": true, "chi_phi_hoat_dong.tao": true, "chi_phi_hoat_dong.sua": true, "chi_phi_hoat_dong.xoa": true,

	// Website
	"don_hang_online.xem": true, "don_hang_online.sua": true, "don_hang_online.duyet": true, "don_hang_online.huy": true, "don_hang_online.phan_cong": true, "don_hang_online.in": true, "don_hang_online.export": true,
	"giao_dien.xem": true, "giao_dien.sua": true, "giao_dien.cau_hinh": true,
	"trang_noi_dung.xem": true, "trang_noi_dung.tao": true, "trang_noi_dung.sua": true, "trang_noi_dung.xoa": true,
	"menu_website.xem": true, "menu_website.tao": true, "menu_website.sua": true, "menu_website.xoa": true,
	"banner.xem": true, "banner.tao": true, "banner.sua": true, "banner.xoa": true,
	"bai_viet.xem": true, "bai_viet.tao": true, "bai_viet.sua": true, "bai_viet.xoa": true, "bai_viet.copy": true,
	"danh_muc_bai_viet.xem": true, "danh_muc_bai_viet.tao": true, "danh_muc_bai_viet.sua": true, "danh_muc_bai_viet.xoa": true,
	"binh_luan.xem": true, "binh_luan.tra_loi": true, "binh_luan.duyet": true, "binh_luan.xoa": true,
	"cau_hinh_seo.xem": true, "cau_hinh_seo.sua": true,
	"domain_website.xem": true, "domain_website.cau_hinh": true,
	"form_lien_he.xem": true, "form_lien_he.xoa": true, "form_lien_he.export": true,
	"live_chat.xem": true, "live_chat.tra_loi": true, "live_chat.cau_hinh": true,
	"landing_page.xem": true, "landing_page.tao": true, "landing_page.sua": true, "landing_page.xoa": true, "landing_page.copy": true,

	// Vận Chuyển & Đa Kênh
	"van_chuyen.xem": true, "van_chuyen.tao": true, "van_chuyen.sua": true, "van_chuyen.huy": true, "van_chuyen.in": true,
	"don_vi_van_chuyen.xem": true, "don_vi_van_chuyen.cau_hinh": true,
	"so_sanh_gia_cuoc.xem": true,
	"tao_van_don.tao": true,
	"tracking_don_hang.xem": true,
	"chi_nhanh.xem": true, "chi_nhanh.tao": true, "chi_nhanh.sua": true, "chi_nhanh.xoa": true,
	"bao_cao_chi_nhanh.xem": true, "bao_cao_chi_nhanh.export": true,
	"ton_kho_chi_nhanh.xem": true,
	"gia_chi_nhanh.xem": true, "gia_chi_nhanh.sua": true,
	"nhan_vien_chi_nhanh.xem": true, "nhan_vien_chi_nhanh.sua": true, "nhan_vien_chi_nhanh.phan_cong": true,
	"ket_noi_san_tmdt.xem": true, "ket_noi_san_tmdt.cau_hinh": true,
	"chat_da_kenh.xem": true, "chat_da_kenh.tra_loi": true, "chat_da_kenh.tao": true,
	"dong_bo_ton_kho.xem": true, "dong_bo_ton_kho.sua": true,
	"gia_theo_kenh.xem": true, "gia_theo_kenh.sua": true,

	// Thanh Toán, Nhân Sự & Báo Cáo
	"qr_code.tao": true,
	"cong_thanh_toan.xem": true, "cong_thanh_toan.cau_hinh": true,
	"hoa_don_dien_tu.xem": true, "hoa_don_dien_tu.tao": true, "hoa_don_dien_tu.phat_hanh": true, "hoa_don_dien_tu.huy": true, "hoa_don_dien_tu.in": true,
	"nhan_vien.xem": true, "nhan_vien.tao": true, "nhan_vien.sua": true, "nhan_vien.xoa": true, "nhan_vien.reset_mat_khau": true, "nhan_vien.khoa": true, "nhan_vien.phan_cong": true, "nhan_vien.xem_lich_su": true,
	"phan_quyen.xem": true, "phan_quyen.tao": true, "phan_quyen.sua": true, "phan_quyen.xoa": true,
	"bao_cao_ban_hang.xem": true, "bao_cao_ban_hang.export": true,
	"bao_cao_nhap_hang.xem": true, "bao_cao_nhap_hang.export": true,
	"bao_cao_lai_lo.xem": true, "bao_cao_lai_lo.export": true,
	"bao_cao_cong_no.xem": true, "bao_cao_cong_no.export": true,
	"bao_cao_ton_kho.xem": true, "bao_cao_ton_kho.export": true,
	"bao_cao_san_pham_ban_chay.xem": true, "bao_cao_san_pham_ban_chay.export": true,
	"bao_cao_nhan_vien_ban_hang.xem": true, "bao_cao_nhan_vien_ban_hang.export": true,
	"bao_cao_khach_hang.xem": true, "bao_cao_khach_hang.export": true,

	// Marketing & AI
	"sms_marketing.xem": true, "sms_marketing.tao": true, "sms_marketing.gui": true,
	"email_marketing.xem": true, "email_marketing.tao": true, "email_marketing.gui": true,
	"khuyen_mai_tu_dong.xem": true, "khuyen_mai_tu_dong.tao": true, "khuyen_mai_tu_dong.sua": true, "khuyen_mai_tu_dong.xoa": true,
	"workflow.xem": true, "workflow.tao": true, "workflow.sua": true, "workflow.xoa": true,
	"webhook.xem": true, "webhook.tao": true, "webhook.xoa": true,
	"ai_bao_cao.xem": true,
	"ai_du_doan_ton_kho.xem": true,

	// Hệ Thống (Đã bổ sung: Tin nhắn Hệ thống)
	"goi_dich_vu.xem": true, "goi_dich_vu.tao": true, "goi_dich_vu.sua": true, "goi_dich_vu.xoa": true,
	"thanh_toan_goi.xem": true, "thanh_toan_goi.tao": true,
	"lich_su_thanh_toan.xem": true,
	"import_du_lieu.thuc_hien": true, "export_du_lieu.thuc_hien": true, "sao_luu_du_lieu.thuc_hien": true, "khoi_phuc_du_lieu.thuc_hien": true,
	"nhat_ky_hoat_dong.xem": true, "nhat_ky_truy_cap.xem": true, "nhat_ky_truy_cap.export": true,
	"cau_hinh_cua_hang.xem": true, "cau_hinh_cua_hang.sua": true,
	"in_an.cau_hinh": true,
	"cau_hinh_vat.xem": true, "cau_hinh_vat.sua": true,
	"ho_so.xem": true, "ho_so.sua": true,
	"api_key.xem": true, "api_key.tao": true, "api_key.sua": true, "api_key.xoa": true,
	"tep_tin.xem": true, "tep_tin.tao": true, "tep_tin.xoa": true,
	"tin_nhan_he_thong.xem": true, "tin_nhan_he_thong.gui": true, "tin_nhan_he_thong.xoa": true,
}
