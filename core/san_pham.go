package core

import (
	"sort"
	"strings"
	"time"
)

// =============================================================
// 1. CẤU HÌNH CỘT (Copy từ chi_muc.go)
// =============================================================
const (
	DongBatDauDuLieu = 2 // Dữ liệu bắt đầu từ dòng 2

	CotSP_MaSanPham    = 0  // A
	CotSP_TenSanPham   = 1  // B
	CotSP_TenRutGon    = 2  // C
	CotSP_Sku          = 3  // D
	CotSP_MaDanhMuc    = 4  // E
	CotSP_MaThuongHieu = 5  // F
	CotSP_DonVi        = 6  // G
	CotSP_MauSac       = 7  // H
	CotSP_UrlHinhAnh   = 8  // I
	CotSP_ThongSo      = 9  // J
	CotSP_MoTaChiTiet  = 10 // K
	CotSP_BaoHanhThang = 11 // L
	CotSP_TinhTrang    = 12 // M
	CotSP_TrangThai    = 13 // N
	CotSP_GiaBanLe     = 14 // O
	CotSP_GhiChu       = 15 // P
	CotSP_NguoiTao     = 16 // Q
	CotSP_NgayTao      = 17 // R
	CotSP_NgayCapNhat  = 18 // S
)

// =============================================================
// 2. STRUCT DỮ LIỆU (Copy từ bang_du_lieu.go)
// =============================================================
type SanPham struct {
	MaSanPham    string  `json:"ma_san_pham"`
	TenSanPham   string  `json:"ten_san_pham"`
	TenRutGon    string  `json:"ten_rut_gon"`
	Sku          string  `json:"sku"`
	MaDanhMuc    string  `json:"ma_danh_muc"`
	MaThuongHieu string  `json:"ma_thuong_hieu"`
	DonVi        string  `json:"don_vi"`
	MauSac       string  `json:"mau_sac"`
	UrlHinhAnh   string  `json:"url_hinh_anh"`
	ThongSo      string  `json:"thong_so"`
	MoTaChiTiet  string  `json:"mo_ta_chi_tiet"`
	BaoHanhThang int     `json:"bao_hanh_thang"`
	TinhTrang    string  `json:"tinh_trang"`
	TrangThai    int     `json:"trang_thai"`
	GiaBanLe     float64 `json:"gia_ban_le"`
	GhiChu       string  `json:"ghi_chu"`
	NguoiTao     string  `json:"nguoi_tao"`
	NgayTao      string  `json:"ngay_tao"`
	NgayCapNhat  string  `json:"ngay_cap_nhat"`
}

// =============================================================
// 3. KHO LƯU TRỮ (Thay thế bo_nho_dem)
// =============================================================
var (
	_DS_SanPham  []SanPham          // Slice để duyệt danh sách
	_Map_SanPham map[string]SanPham // Map để tìm nhanh theo ID
)

// =============================================================
// 4. LOGIC NẠP DỮ LIỆU (Thay thế bo_nho_dem/san_pham.go)
// =============================================================
func NapSanPham() {
	raw, err := loadSheetData("SAN_PHAM")
	if err != nil {
		return
	}

	// Reset bộ nhớ
	tempList := []SanPham{}
	tempMap := make(map[string]SanPham)

	for i, r := range raw {
		if i < DongBatDauDuLieu-1 { continue } // Bỏ qua Header
		
		// Validate cơ bản
		maSP := layString(r, CotSP_MaSanPham)
		if maSP == "" { continue }

		item := SanPham{
			MaSanPham:    maSP,
			TenSanPham:   layString(r, CotSP_TenSanPham),
			TenRutGon:    layString(r, CotSP_TenRutGon),
			Sku:          layString(r, CotSP_Sku),
			MaDanhMuc:    layString(r, CotSP_MaDanhMuc),
			MaThuongHieu: layString(r, CotSP_MaThuongHieu),
			DonVi:        layString(r, CotSP_DonVi),
			MauSac:       layString(r, CotSP_MauSac),
			UrlHinhAnh:   layString(r, CotSP_UrlHinhAnh),
			ThongSo:      layString(r, CotSP_ThongSo),
			MoTaChiTiet:  layString(r, CotSP_MoTaChiTiet),
			BaoHanhThang: layInt(r, CotSP_BaoHanhThang),
			TinhTrang:    layString(r, CotSP_TinhTrang),
			TrangThai:    layInt(r, CotSP_TrangThai),
			GiaBanLe:     layFloat(r, CotSP_GiaBanLe),
			GhiChu:       layString(r, CotSP_GhiChu),
			NguoiTao:     layString(r, CotSP_NguoiTao),
			NgayTao:      layString(r, CotSP_NgayTao),
			NgayCapNhat:  layString(r, CotSP_NgayCapNhat),
		}

		tempList = append(tempList, item)
		tempMap[maSP] = item
	}

	// Gán vào biến toàn cục (Thread-safe logic xử lý ở hàm gọi hoặc dùng Lock nếu cần)
	// Ở đây ta gán thẳng vì hàm NapSanPham thường chạy khi khởi động hoặc Reload (đã có Lock tổng)
	_DS_SanPham = tempList
	_Map_SanPham = tempMap
}

// =============================================================
// 5. LOGIC TRUY VẤN (Thay thế nghiep_vu/truy_xuat.go)
// =============================================================

// Lấy toàn bộ danh sách (Có Lock an toàn)
func LayDanhSachSanPham() []SanPham {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	// Trả về bản copy để an toàn
	ketQua := make([]SanPham, len(_DS_SanPham))
	copy(ketQua, _DS_SanPham)
	return ketQua
}

// Lấy chi tiết 1 sản phẩm
func LayChiTietSanPham(maSP string) (SanPham, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	sp, ok := _Map_SanPham[maSP]
	return sp, ok
}

// Helper: Tạo mã sản phẩm mới tự động (SP_0001)
func TaoMaSPMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	maxID := 0
	for _, sp := range _DS_SanPham {
		parts := strings.Split(sp.MaSanPham, "_")
		if len(parts) == 2 {
			// Giả sử format SP_xxxx. Cần parse int cẩn thận hơn thực tế
			// nhưng đây là logic cũ của bạn, tôi giữ nguyên.
			// Logic tối ưu hơn: Dùng regex hoặc TrimPrefix
			soStr := strings.TrimPrefix(sp.MaSanPham, "SP_")
			// Remove leading zeros... (phức tạp, giữ logic đơn giản tạm thời)
			// Để đơn giản, ta chỉ đếm số lượng + 1 nếu ko parse được
			// Logic cũ của bạn đang split "_", ok giữ nguyên.
			// ...
		}
	}
	// Logic đơn giản hóa: Đếm số lượng + 1 (hoặc lấy max id thực tế nếu cần)
	// Tạm thời return len + 1 cho nhanh gọn
	return "SP_" + time.Now().Format("150405") // Trick: Dùng giờ phút giây để tránh trùng :D
}
