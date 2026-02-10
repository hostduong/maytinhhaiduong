package core

import (
	"fmt"
	"strings"
	"time"

	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT (CẬP NHẬT THEO CHIẾN THUẬT MỚI)
// =============================================================
const (
	DongBatDau_SanPham = 11

	// A=0, B=1, ... T=19
	CotSP_MaSanPham      = 0  // A
	CotSP_TenSanPham     = 1  // B
	CotSP_TenRutGon      = 2  // C
	CotSP_Slug           = 3  // D [MỚI]
	CotSP_Sku            = 4  // E
	CotSP_DanhMuc        = 5  // F (Lưu trực tiếp: "Main|Máy tính")
	CotSP_ThuongHieu     = 6  // G (Lưu trực tiếp: "Asus")
	CotSP_DonVi          = 7  // H
	CotSP_MauSac         = 8  // I
	CotSP_UrlHinhAnh     = 9  // J
	CotSP_ThongSo        = 10 // K
	CotSP_MoTaChiTiet    = 11 // L
	CotSP_BaoHanhThang   = 12 // M
	CotSP_TinhTrang      = 13 // N
	CotSP_TrangThai      = 14 // O
	CotSP_GiaBanLe       = 15 // P
	CotSP_GhiChu         = 16 // Q
	CotSP_NguoiTao       = 17 // R
	CotSP_NgayTao        = 18 // S
	CotSP_NgayCapNhat    = 19 // T
)

type SanPham struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaSanPham      string  `json:"ma_san_pham"`
	TenSanPham     string  `json:"ten_san_pham"`
	TenRutGon      string  `json:"ten_rut_gon"`
	Slug           string  `json:"slug"` // [MỚI]
	Sku            string  `json:"sku"`
	DanhMuc        string  `json:"danh_muc"`    // Chuỗi trực tiếp
	ThuongHieu     string  `json:"thuong_hieu"` // Chuỗi trực tiếp
	DonVi          string  `json:"don_vi"`
	MauSac         string  `json:"mau_sac"`
	UrlHinhAnh     string  `json:"url_hinh_anh"`
	ThongSo        string  `json:"thong_so"`
	MoTaChiTiet    string  `json:"mo_ta_chi_tiet"`
	BaoHanhThang   int     `json:"bao_hanh_thang"`
	TinhTrang      string  `json:"tinh_trang"`
	TrangThai      int     `json:"trang_thai"`
	GiaBanLe       float64 `json:"gia_ban_le"`
	GhiChu         string  `json:"ghi_chu"`
	NguoiTao       string  `json:"nguoi_tao"`
	NgayTao        string  `json:"ngay_tao"`
	NgayCapNhat    string  `json:"ngay_cap_nhat"`
}

var (
	_DS_SanPham  []*SanPham
	_Map_SanPham map[string]*SanPham
)

func NapSanPham(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" { targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(targetSpreadsheetID, "SAN_PHAM")
	if err != nil { return }

	_Map_SanPham = make(map[string]*SanPham)
	_DS_SanPham = []*SanPham{}

	for i, r := range raw {
		if i < DongBatDau_SanPham-1 { continue }
		maSP := layString(r, CotSP_MaSanPham)
		if maSP == "" { continue }

		key := TaoCompositeKey(targetSpreadsheetID, maSP)
		if _, daTonTai := _Map_SanPham[key]; daTonTai { continue }

		sp := &SanPham{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaSanPham:      maSP,
			TenSanPham:     layString(r, CotSP_TenSanPham),
			TenRutGon:      layString(r, CotSP_TenRutGon),
			Slug:           layString(r, CotSP_Slug), // Đọc cột D
			Sku:            layString(r, CotSP_Sku),
			DanhMuc:        layString(r, CotSP_DanhMuc),    // Cột F
			ThuongHieu:     layString(r, CotSP_ThuongHieu), // Cột G
			DonVi:          layString(r, CotSP_DonVi),
			MauSac:         layString(r, CotSP_MauSac),
			UrlHinhAnh:     layString(r, CotSP_UrlHinhAnh),
			ThongSo:        layString(r, CotSP_ThongSo),
			MoTaChiTiet:    layString(r, CotSP_MoTaChiTiet),
			BaoHanhThang:   layInt(r, CotSP_BaoHanhThang),
			TinhTrang:      layString(r, CotSP_TinhTrang),
			TrangThai:      layInt(r, CotSP_TrangThai),
			GiaBanLe:       layFloat(r, CotSP_GiaBanLe),
			GhiChu:         layString(r, CotSP_GhiChu),
			NguoiTao:       layString(r, CotSP_NguoiTao),
			NgayTao:        layString(r, CotSP_NgayTao),
			NgayCapNhat:    layString(r, CotSP_NgayCapNhat),
		}
		_DS_SanPham = append(_DS_SanPham, sp)
		_Map_SanPham[key] = sp
	}
}

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
	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maSP)
	sp, ok := _Map_SanPham[key]
	return sp, ok
}

func ThemSanPhamVaoRam(sp *SanPham) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if sp.SpreadsheetID == "" { sp.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_SanPham = append(_DS_SanPham, sp)
	key := TaoCompositeKey(sp.SpreadsheetID, sp.MaSanPham)
	_Map_SanPham[key] = sp
}

// [HÀM SINH MÃ CHUẨN]
func TaoMaSPMoi(prefixThuongHieu string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	
	partBrand := "000"
	prefixThuongHieu = strings.ToUpper(strings.TrimSpace(prefixThuongHieu))
	
	// Logic lấy 3 ký tự đầu của thương hiệu làm mã (VD: ASUS -> ASU)
	if len(prefixThuongHieu) >= 3 {
		partBrand = prefixThuongHieu[:3]
	} else if len(prefixThuongHieu) > 0 {
		partBrand = prefixThuongHieu + "X" // VD: HP -> HPX
	}

	partTime := time.Now().Format("0601") // YYMM

	for {
		rand4 := LayChuoiSoNgauNhien(4)
		id := fmt.Sprintf("HD%s%s%s", partBrand, partTime, rand4)
		key := TaoCompositeKey(currentSheetID, id)
		if _, tonTai := _Map_SanPham[key]; !tonTai {
			return id
		}
	}
}
