package core

import (
	"fmt"
	"strconv"
	"strings"

	"app/cau_hinh"
)

const TenSheetMayTinh = "MAY_TINH"

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

func (sp *SanPhamMayTinh) LayIDDuyNhat() string {
	if sp.MaSKU != "" { return sp.MaSKU }
	return sp.MaSanPham
}

var (
	CacheSanPhamMayTinh      = make(map[string][]*SanPhamMayTinh)
	CacheMapSKUMayTinh       = make(map[string]*SanPhamMayTinh) 
	CacheGroupSanPhamMayTinh = make(map[string][]*SanPhamMayTinh)
)

func NapDuLieuMayTinh(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, TenSheetMayTinh) 
	if err != nil { return }

	list := []*SanPhamMayTinh{}
	groupMap := make(map[string][]*SanPhamMayTinh)

	for i, r := range raw {
		if i < DongBatDau_SanPhamMayTinh-1 { continue }
		maSP := LayString(r, CotPC_MaSanPham)
		if maSP == "" { continue }

		sp := &SanPhamMayTinh{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaSanPham:      maSP,
			TenSanPham:     LayString(r, CotPC_TenSanPham),
			TenRutGon:      LayString(r, CotPC_TenRutGon),
			Slug:           LayString(r, CotPC_Slug),
			MaSKU:          LayString(r, CotPC_MaSKU),
			TenSKU:         LayString(r, CotPC_TenSKU),
			SKUChinh:       LayInt(r, CotPC_SKUChinh),
			TrangThai:      LayInt(r, CotPC_TrangThai),
			MaDanhMuc:      LayString(r, CotPC_MaDanhMuc),
			MaThuongHieu:   LayString(r, CotPC_MaThuongHieu),
			DonVi:          LayString(r, CotPC_DonVi),
			MauSac:         LayString(r, CotPC_MauSac),
			KhoiLuong:      LayFloat(r, CotPC_KhoiLuong),
			KichThuoc:      LayString(r, CotPC_KichThuoc),
			UrlHinhAnh:     LayString(r, CotPC_UrlHinhAnh),
			ThongSoHTML:    LayString(r, CotPC_ThongSoHTML),
			MoTaHTML:       LayString(r, CotPC_MoTaHTML),
			BaoHanh:        LayString(r, CotPC_BaoHanh),
			TinhTrang:      LayString(r, CotPC_TinhTrang),
			GiaNhap:        LayFloat(r, CotPC_GiaNhap),
			PhanTramLai:    LayFloat(r, CotPC_PhanTramLai),
			GiaNiemYet:     LayFloat(r, CotPC_GiaNiemYet),
			PhanTramGiam:   LayFloat(r, CotPC_PhanTramGiam),
			SoTienGiam:     LayFloat(r, CotPC_SoTienGiam),
			GiaBan:         LayFloat(r, CotPC_GiaBan),
			GhiChu:         LayString(r, CotPC_GhiChu),
			NguoiTao:       LayString(r, CotPC_NguoiTao),
			NgayTao:        LayString(r, CotPC_NgayTao),
			NguoiCapNhat:   LayString(r, CotPC_NguoiCapNhat),
			NgayCapNhat:    LayString(r, CotPC_NgayCapNhat),
		}
		
		list = append(list, sp)
		
		keySingle := TaoCompositeKey(shopID, sp.LayIDDuyNhat())
		CacheMapSKUMayTinh[keySingle] = sp

		keyGroup := TaoCompositeKey(shopID, maSP)
		groupMap[keyGroup] = append(groupMap[keyGroup], sp)
	}

	KhoaHeThong.Lock()
	CacheSanPhamMayTinh[shopID] = list
	for k, v := range groupMap { CacheGroupSanPhamMayTinh[k] = v }
	KhoaHeThong.Unlock()
}

func LayDanhSachSanPhamMayTinh(shopID string) []*SanPhamMayTinh {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheSanPhamMayTinh[shopID]
}

func TaoMaSPMayTinhMoi(shopID string, prefix string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	if prefix == "" { prefix = "SP" }
	maxNum := 0
	list := CacheSanPhamMayTinh[shopID]
	for _, sp := range list {
		if strings.HasPrefix(sp.MaSanPham, prefix) {
			numStr := strings.TrimPrefix(sp.MaSanPham, prefix)
			if num, err := strconv.Atoi(numStr); err == nil {
				if num > maxNum { maxNum = num }
			}
		}
	}
	return fmt.Sprintf("%s%04d", prefix, maxNum+1)
}

func LayChiTietSKUMayTinh(shopID, idDuyNhat string) (*SanPhamMayTinh, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(shopID, idDuyNhat)
	sp, ok := CacheMapSKUMayTinh[key]
	return sp, ok
}
