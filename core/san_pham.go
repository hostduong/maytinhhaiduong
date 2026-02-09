package core

import (
	"fmt"      // [THÊM]
	"strings"
	"time"     // [THÊM]

	"app/cau_hinh"
)

// ... (Giữ nguyên phần Const và Struct) ...

// =============================================================
// 1. CẤU HÌNH CỘT
// =============================================================
const (
	DongBatDau_SanPham = 11

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
func TaoMaSPMoi(maThuongHieu string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	
	partBrand := "000"
	maThuongHieu = strings.ToUpper(strings.TrimSpace(maThuongHieu))
	
	if strings.HasPrefix(maThuongHieu, "TH") && len(maThuongHieu) >= 5 {
		partBrand = maThuongHieu[2:] 
	} else if len(maThuongHieu) >= 3 {
		partBrand = maThuongHieu[len(maThuongHieu)-3:]
	}
	if len(partBrand) > 3 { partBrand = partBrand[:3] }

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
