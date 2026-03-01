package core


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
// CẤU TRÚC: SẢN PHẨM MÁY TÍNH
// ==============================================================================
const (
	DongBatDau_SanPhamMayTinh = 11
	
	CotPC_MaSanPham      = 0  
	CotPC_TenSanPham     = 1  
	CotPC_TenRutGon      = 2  
	CotPC_Slug           = 3  
	CotPC_MaSKU          = 4  
	CotPC_TenSKU         = 5  
	CotPC_SKUChinh       = 6  
	CotPC_TrangThai      = 7  
	CotPC_MaDanhMuc      = 8  
	CotPC_MaThuongHieu   = 9  
	CotPC_DonVi          = 10 
	CotPC_MauSac         = 11 
	CotPC_KhoiLuong      = 12 
	CotPC_KichThuoc      = 13 
	CotPC_UrlHinhAnh     = 14 
	CotPC_ThongSoHTML    = 15 
	CotPC_MoTaHTML       = 16 
	CotPC_BaoHanh        = 17 
	CotPC_TinhTrang      = 18 
	CotPC_GiaNhap        = 19 
	CotPC_PhanTramLai    = 20  
	CotPC_GiaNiemYet     = 21 
	CotPC_PhanTramGiam   = 22 
	CotPC_SoTienGiam     = 23 
	CotPC_GiaBan         = 24 
	CotPC_GhiChu         = 25 
	CotPC_NguoiTao       = 26 
	CotPC_NgayTao        = 27 
	CotPC_NguoiCapNhat   = 28 
	CotPC_NgayCapNhat    = 29 
)

type SanPhamMayTinh struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaSanPham      string  `json:"ma_san_pham"`
	TenSanPham     string  `json:"ten_san_pham"`
	TenRutGon      string  `json:"ten_rut_gon"`
	Slug           string  `json:"slug"`
	MaSKU          string  `json:"ma_sku"`
	TenSKU         string  `json:"ten_sku"`
	SKUChinh       int     `json:"sku_chinh"`
	TrangThai      int     `json:"trang_thai"`
	MaDanhMuc      string  `json:"ma_danh_muc"`
	MaThuongHieu   string  `json:"ma_thuong_hieu"`
	DonVi          string  `json:"don_vi"`
	MauSac         string  `json:"mau_sac"`
	KhoiLuong      float64 `json:"khoi_luong"`
	KichThuoc      string  `json:"kich_thuoc"`
	UrlHinhAnh     string  `json:"url_hinh_anh"`
	ThongSoHTML    string  `json:"thong_so_html"`
	MoTaHTML       string  `json:"mo_ta_html"`
	BaoHanh        string  `json:"bao_hanh"`
	TinhTrang      string  `json:"tinh_trang"`
	GiaNhap        float64 `json:"gia_nhap"`
	PhanTramLai    float64 `json:"phan_tram_lai"`
	GiaNiemYet     float64 `json:"gia_niem_yet"`
	PhanTramGiam   float64 `json:"phan_tram_giam"`
	SoTienGiam     float64 `json:"so_tien_giam"`
	GiaBan         float64 `json:"gia_ban"`
	GhiChu         string  `json:"ghi_chu"`
	NguoiTao       string  `json:"nguoi_tao"`
	NgayTao        string  `json:"ngay_tao"`
	NguoiCapNhat   string  `json:"nguoi_cap_nhat"`
	NgayCapNhat    string  `json:"ngay_cap_nhat"`
}
// Hàm vệ tinh tạo ID duy nhất cho Sản phẩm máy tính
func (sp *SanPhamMayTinh) LayIDDuyNhat() string {
	if sp.MaSKU != "" {
		return sp.MaSKU
	}
	return sp.MaSanPham
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
	
	// [ĐÃ FIX]: Thêm biến ảo để Front-end biết tin này đã đọc hay chưa
	DaDoc          bool          `json:"da_doc"` 
}

// ==============================================================================
// CẤU TRÚC: DANH MỤC
// ==============================================================================
const (
	CotPN_MaPhieuNhap          = 0  // A
	CotPN_MaNhaCungCap         = 1  // B
	CotPN_MaKho                = 2  // C
	CotPN_NgayNhap             = 3  // D
	CotPN_TrangThai            = 4  // E
	CotPN_SoHoaDon             = 5  // F
	CotPN_NgayHoaDon           = 6  // G
	CotPN_UrlChungTu           = 7  // H
	CotPN_TongTienPhieu        = 8  // I
	CotPN_GiamGiaPhieu         = 9  // J
	CotPN_DaThanhToan          = 10 // K
	CotPN_ConNo                = 11 // L
	CotPN_PhuongThucThanhToan  = 12 // M
	CotPN_TrangThaiThanhToan   = 13 // N
	CotPN_GhiChu               = 14 // O
	CotPN_NguoiTao             = 15 // P
	CotPN_NgayTao              = 16 // Q
	CotPN_NgayCapNhat          = 17 // R
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
	TrangThai            int     `json:"trang_thai"`
	SoHoaDon             string  `json:"so_hoa_don"`
	NgayHoaDon           string  `json:"ngay_hoa_don"`
	UrlChungTu           string  `json:"url_chung_tu"`
	TongTienPhieu        float64 `json:"tong_tien_phieu"`
	GiamGiaPhieu         float64 `json:"giam_gia_phieu"`
	DaThanhToan          float64 `json:"da_thanh_toan"`
	ConNo                float64 `json:"con_no"`
	PhuongThucThanhToan  string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiThanhToan   string  `json:"trang_thai_thanh_toan"`
	GhiChu               string  `json:"ghi_chu"`
	NguoiTao             string  `json:"nguoi_tao"`
	NgayTao              string  `json:"ngay_tao"`
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
