package core

import (
	"fmt"
	"strings"
	"app/cau_hinh"
)

// =============================================================
// CẤU HÌNH CỘT (GIỮ NGUYÊN)
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
	CotSP_BaoHanh        = 12 
	CotSP_TinhTrang      = 13 
	CotSP_TrangThai      = 14 
	CotSP_GiaNhap        = 15 
	CotSP_GiaBanLe       = 16 
	CotSP_GiamGia        = 17 
	CotSP_GiaBanThuc     = 18 
	CotSP_GhiChu         = 19 
	CotSP_NguoiTao       = 20 
	CotSP_NgayTao        = 21 
	CotSP_NgayCapNhat    = 22 
)

// STRUCT SẢN PHẨM (GIỮ NGUYÊN)
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

// BỘ NHỚ ĐA SHOP (THAY THẾ BIẾN CŨ)
// Key của map ngoài cùng là ShopID (SpreadsheetID)
var (
	CacheSanPham    = make(map[string][]*SanPham)
	CacheMapSanPham = make(map[string]map[string]*SanPham) // Key: ShopID__MaSP
)

// HÀM NẠP (THÊM THAM SỐ shopID)
func NapSanPham(shopID string) {
	// Nếu shopID rỗng (trường hợp init), lấy mặc định
	if shopID == "" { 
		shopID = cau_hinh.BienCauHinh.IdFileSheet 
	}

	raw, err := loadSheetData(shopID, "SAN_PHAM")
	if err != nil { return }

	// Tạo list tạm cho Shop này
	list := []*SanPham{}

	for i, r := range raw {
		if i < DongBatDau_SanPham-1 { continue }
		maSP := layString(r, CotSP_MaSanPham)
		if maSP == "" { continue }

		sp := &SanPham{
			SpreadsheetID:  shopID,
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
			BaoHanh:        layString(r, CotSP_BaoHanh),
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
		
		list = append(list, sp)
		
		// Lưu vào Map tra cứu nhanh
		key := TaoCompositeKey(shopID, maSP)
		CacheMapSanPham[key] = sp
	}

	// Cập nhật vào Cache Tổng
	KhoaHeThong.Lock()
	CacheSanPham[shopID] = list
	KhoaHeThong.Unlock()
}

// HÀM LẤY DANH SÁCH (THEO SHOP ID)
func LayDanhSachSanPham(shopID string) []*SanPham {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	// Trả về đúng list của Shop đó
	if list, ok := CacheSanPham[shopID]; ok {
		return list
	}
	return []*SanPham{}
}

// HÀM LẤY CHI TIẾT (DÙNG COMPOSITE KEY)
func LayChiTietSanPham(shopID, maSP string) (*SanPham, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	key := TaoCompositeKey(shopID, maSP)
	sp, ok := CacheMapSanPham[key]
	return sp, ok
}

// HÀM THÊM MỚI (UPDATE CACHE CỦA SHOP)
func ThemSanPhamVaoRam(sp *SanPham) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := sp.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
	
	// Append vào list của Shop đó
	CacheSanPham[sID] = append(CacheSanPham[sID], sp)
	
	// Update Map
	key := TaoCompositeKey(sID, sp.MaSanPham)
	CacheMapSanPham[key] = sp
}

// HÀM SINH MÃ MỚI (CẦN SHOP ID ĐỂ ĐẾM SLOT ĐÚNG SHOP)
func TaoMaSPMoi(shopID, maDanhMuc string) string {
	maDanhMuc = strings.ToUpper(strings.TrimSpace(maDanhMuc))
	if maDanhMuc == "" { maDanhMuc = "SP" }
	
	// Gọi hàm LaySlotTiepTheo cần truyền shopID vào
	slot := LaySlotTiepTheo(shopID, maDanhMuc) 
	return fmt.Sprintf("%s%04d", maDanhMuc, slot)
}
