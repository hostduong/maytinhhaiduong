package core

import (
	"fmt"
	"strings"
	"time"

	"app/cau_hinh"
)

// =============================================================
// CẤU HÌNH CỘT (CẬP NHẬT MỚI NHẤT P-Q-R-S)
// =============================================================
const (
	DongBatDau_SanPham = 11

	CotSP_MaSanPham      = 0  
	CotSP_TenSanPham     = 1  
	CotSP_TenRutGon      = 2  
	CotSP_Slug           = 3  
	CotSP_Sku            = 4  
	CotSP_DanhMuc        = 5  
	CotSP_ThuongHieu     = 6  
	CotSP_DonVi          = 7  
	CotSP_MauSac         = 8  
	CotSP_UrlHinhAnh     = 9  
	CotSP_ThongSo        = 10 
	CotSP_MoTaChiTiet    = 11 
	CotSP_BaoHanhThang   = 12 
	CotSP_TinhTrang      = 13 
	CotSP_TrangThai      = 14 
	
	CotSP_GiaNhap        = 15 // P
	CotSP_GiaBanLe       = 16 // Q
	CotSP_GiamGia        = 17 // R
	CotSP_GiaBanThuc     = 18 // S
	
	CotSP_GhiChu         = 19 // T
	CotSP_NguoiTao       = 20 // U
	CotSP_NgayTao        = 21 // V
	CotSP_NgayCapNhat    = 22 // W
)

type SanPham struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaSanPham      string  `json:"ma_san_pham"`
	TenSanPham     string  `json:"ten_san_pham"`
	TenRutGon      string  `json:"ten_rut_gon"`
	Slug           string  `json:"slug"`
	Sku            string  `json:"sku"`
	DanhMuc        string  `json:"danh_muc"`
	ThuongHieu     string  `json:"thuong_hieu"`
	DonVi          string  `json:"don_vi"`
	MauSac         string  `json:"mau_sac"`
	UrlHinhAnh     string  `json:"url_hinh_anh"`
	ThongSo        string  `json:"thong_so"`
	MoTaChiTiet    string  `json:"mo_ta_chi_tiet"`
	BaoHanh        string  `json:"bao_hanh"` 
	TinhTrang      string  `json:"tinh_trang"`
	TrangThai      int     `json:"trang_thai"`
	
	GiaNhap        float64 `json:"gia_nhap"`
	GiaBanLe       float64 `json:"gia_ban_le"`
	GiamGia        float64 `json:"giam_gia"`
	GiaBanThuc     float64 `json:"gia_ban_thuc"`
	
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
			Slug:           layString(r, CotSP_Slug),
			Sku:            layString(r, CotSP_Sku),
			DanhMuc:        layString(r, CotSP_DanhMuc),
			ThuongHieu:     layString(r, CotSP_ThuongHieu),
			DonVi:          layString(r, CotSP_DonVi),
			MauSac:         layString(r, CotSP_MauSac),
			UrlHinhAnh:     layString(r, CotSP_UrlHinhAnh),
			ThongSo:        layString(r, CotSP_ThongSo),
			MoTaChiTiet:    layString(r, CotSP_MoTaChiTiet),
			BaoHanh:        layString(r, CotSP_BaoHanhThang), 		
			TinhTrang:      layString(r, CotSP_TinhTrang),
			TrangThai:      layInt(r, CotSP_TrangThai),
			
			GiaNhap:        layFloat(r, CotSP_GiaNhap),
			GiaBanLe:       layFloat(r, CotSP_GiaBanLe),
			GiamGia:        layFloat(r, CotSP_GiamGia),
			GiaBanThuc:     layFloat(r, CotSP_GiaBanThuc),
			
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

// [HÀM SINH MÃ MỚI: DỰA VÀO DANH MỤC]
// Input: maDanhMuc (VD: "MAIN")
// Output: "MAIN0001", "MAIN0002"...
func TaoMaSPMoi(maDanhMuc string) string {
	// Không cần Lock ở đây vì hàm LaySTTtiepTheo đã tự Lock rồi
	
	maDanhMuc = strings.ToUpper(strings.TrimSpace(maDanhMuc))
	if maDanhMuc == "" {
		maDanhMuc = "SP" // Fallback nếu không có danh mục
	}

	// Lấy số thứ tự tiếp theo từ core/danh_muc.go
	stt := LaySTTtiepTheo(maDanhMuc)

	// Format thành chuỗi 4 số (0001, 0015...)
	return fmt.Sprintf("%s%04d", maDanhMuc, stt)
}
