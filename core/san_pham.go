package core

import (
	"fmt"
	"strings"

	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT
// =============================================================
const (
	// [CẬP NHẬT] Dữ liệu bắt đầu từ dòng 11
	DongBatDauDuLieuSP = 11

	CotSP_MaSanPham    = 0
	CotSP_TenSanPham   = 1
	CotSP_TenRutGon    = 2
	CotSP_Sku          = 3
	CotSP_MaDanhMuc    = 4
	CotSP_MaThuongHieu = 5
	CotSP_DonVi        = 6
	CotSP_MauSac       = 7
	CotSP_UrlHinhAnh   = 8
	CotSP_ThongSo      = 9
	CotSP_MoTaChiTiet  = 10
	CotSP_BaoHanhThang = 11
	CotSP_TinhTrang    = 12
	CotSP_TrangThai    = 13
	CotSP_GiaBanLe     = 14
	CotSP_GhiChu       = 15
	CotSP_NguoiTao     = 16
	CotSP_NgayTao      = 17
	CotSP_NgayCapNhat  = 18
)

type SanPham struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
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

var (
	_DS_SanPham  []*SanPham
	_Map_SanPham map[string]*SanPham
)

func NapSanPham(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" {
		targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	raw, err := loadSheetData(targetSpreadsheetID, "SAN_PHAM")
	if err != nil { return }

	_Map_SanPham = make(map[string]*SanPham)
	_DS_SanPham = []*SanPham{}

	for i, r := range raw {
		// [CẬP NHẬT] Bỏ qua 10 dòng đầu
		if i < DongBatDauDuLieuSP-1 { continue }
		
		maSP := layString(r, CotSP_MaSanPham)
		
		// [LOGIC MỚI] Chỉ cần có mã là lấy (Bỏ check SP_)
		if maSP == "" { continue }

		key := TaoCompositeKey(targetSpreadsheetID, maSP)

		// Chống trùng lặp (Ưu tiên dòng đầu)
		if _, daTonTai := _Map_SanPham[key]; daTonTai {
			continue
		}

		sp := &SanPham{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			
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

		_DS_SanPham = append(_DS_SanPham, sp)
		_Map_SanPham[key] = sp
	}
}

// ... (Các hàm truy vấn giữ nguyên) ...
func LayDanhSachSanPham() []*SanPham {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	var kq []*SanPham
	for _, sp := range _DS_SanPham {
		if sp.SpreadsheetID == currentSheetID { kq = append(kq, sp) }
	}
	return kq
}

func LayChiTietSanPham(maSP string) (*SanPham, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	key := TaoCompositeKey(currentSheetID, maSP)
	sp, ok := _Map_SanPham[key]
	return sp, ok
}

func TaoMaSPMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	maxID := 0
	for _, sp := range _DS_SanPham {
		if sp.SpreadsheetID != currentSheetID { continue }
		parts := strings.Split(sp.MaSanPham, "_")
		if len(parts) == 2 {
			var id int
			fmt.Sscanf(parts[1], "%d", &id)
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("SP_%04d", maxID+1)
}

func ThemSanPhamVaoRam(sp *SanPham) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if sp.SpreadsheetID == "" { sp.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_SanPham = append(_DS_SanPham, sp)
	key := TaoCompositeKey(sp.SpreadsheetID, sp.MaSanPham)
	_Map_SanPham[key] = sp
}
