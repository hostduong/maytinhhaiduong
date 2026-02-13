package core

import (
	"fmt"
	"strings"
	"app/cau_hinh"
)

const (
	DongBatDau_SanPham = 11
	
	CotSP_MaSanPham      = 0  // A
	CotSP_TenSanPham     = 1  // B
	CotSP_TenRutGon      = 2  // C
	CotSP_Slug           = 3  // D
	CotSP_MaSKU          = 4  // E
	CotSP_TenSKU         = 5  // F
	CotSP_SKUChinh       = 6  // G
	CotSP_TrangThai      = 7  // H
	CotSP_MaDanhMuc      = 8  // I
	CotSP_MaThuongHieu   = 9  // J
	CotSP_DonVi          = 10 // K
	CotSP_MauSac         = 11 // L
	CotSP_KhoiLuong      = 12 // M
	CotSP_KichThuoc      = 13 // N
	CotSP_UrlHinhAnh     = 14 // O
	CotSP_ThongSoHTML    = 15 // P
	CotSP_MoTaHTML       = 16 // Q
	CotSP_BaoHanh        = 17 // R
	CotSP_TinhTrang      = 18 // S
	CotSP_GiaNhap        = 19 // T
	CotSP_GiaNiemYet     = 20 // U
	CotSP_PhanTramGiam   = 21 // V
	CotSP_SoTienGiam     = 22 // W
	CotSP_GiaBan         = 23 // X
	CotSP_GhiChu         = 24 // Y
	CotSP_NguoiTao       = 25 // Z
	CotSP_NgayTao        = 26 // AA
	CotSP_NguoiCapNhat   = 27 // AB
	CotSP_NgayCapNhat    = 28 // AC
)

type SanPham struct {
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

// Hàm lấy ID Duy Nhất để làm Key tra cứu phẳng
// Nếu SKU rỗng (Sản phẩm không biến thể) -> Dùng MaSanPham
// Nếu có SKU (Sản phẩm nhiều biến thể) -> Dùng MaSKU
func (sp *SanPham) LấyIDDuyNhat() string {
	if sp.MaSKU != "" {
		return sp.MaSKU
	}
	return sp.MaSanPham
}

var (
	// Cache danh sách toàn bộ sản phẩm (mỗi phần tử là 1 dòng/1 SKU)
	CacheSanPham = make(map[string][]*SanPham)
	
	// Cache tra cứu chi tiết 1 Dòng/SKU cụ thể (Key: ShopID__IDDuyNhat)
	CacheMapSKU  = make(map[string]*SanPham) 

	// Cache gom nhóm anh em cùng nhà (Key: ShopID__MaSanPham -> Mảng các SKU con)
	CacheGroupSanPham = make(map[string][]*SanPham)
)

func NapSanPham(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(shopID, "SAN_PHAM")
	if err != nil { return }

	list := []*SanPham{}
	groupMap := make(map[string][]*SanPham)

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
			MaSKU:          layString(r, CotSP_MaSKU),
			TenSKU:         layString(r, CotSP_TenSKU),
			SKUChinh:       layInt(r, CotSP_SKUChinh),
			TrangThai:      layInt(r, CotSP_TrangThai),
			MaDanhMuc:      layString(r, CotSP_MaDanhMuc),
			MaThuongHieu:   layString(r, CotSP_MaThuongHieu),
			DonVi:          layString(r, CotSP_DonVi),
			MauSac:         layString(r, CotSP_MauSac),
			KhoiLuong:      layFloat(r, CotSP_KhoiLuong),
			KichThuoc:      layString(r, CotSP_KichThuoc),
			UrlHinhAnh:     layString(r, CotSP_UrlHinhAnh),
			ThongSoHTML:    layString(r, CotSP_ThongSoHTML),
			MoTaHTML:       layString(r, CotSP_MoTaHTML),
			BaoHanh:        layString(r, CotSP_BaoHanh),
			TinhTrang:      layString(r, CotSP_TinhTrang),
			GiaNhap:        layFloat(r, CotSP_GiaNhap),
			GiaNiemYet:     layFloat(r, CotSP_GiaNiemYet),
			PhanTramGiam:   layFloat(r, CotSP_PhanTramGiam),
			SoTienGiam:     layFloat(r, CotSP_SoTienGiam),
			GiaBan:         layFloat(r, CotSP_GiaBan),
			GhiChu:         layString(r, CotSP_GhiChu),
			NguoiTao:       layString(r, CotSP_NguoiTao),
			NgayTao:        layString(r, CotSP_NgayTao),
			NguoiCapNhat:   layString(r, CotSP_NguoiCapNhat),
			NgayCapNhat:    layString(r, CotSP_NgayCapNhat),
		}
		
		list = append(list, sp)
		
		// Map tra cứu 1 dòng duy nhất
		keySingle := TaoCompositeKey(shopID, sp.LấyIDDuyNhat())
		CacheMapSKU[keySingle] = sp

		// Gom nhóm theo MaSanPham
		keyGroup := TaoCompositeKey(shopID, maSP)
		groupMap[keyGroup] = append(groupMap[keyGroup], sp)
	}

	KhoaHeThong.Lock()
	CacheSanPham[shopID] = list
	
	// Cập nhật Group Cache vào RAM chung
	for k, v := range groupMap {
		CacheGroupSanPham[k] = v
	}
	KhoaHeThong.Unlock()
}

// Lấy Toàn bộ Dòng dữ liệu (Danh sách hỗn hợp cả gốc lẫn SKU)
func LayDanhSachSanPham(shopID string) []*SanPham {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheSanPham[shopID]
}

// Tra cứu Chi tiết 1 Dòng (Truyền vào MaSKU, nếu không có truyền MaSanPham)
func LayChiTietSKU(shopID, idDuyNhat string) (*SanPham, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(shopID, idDuyNhat)
	sp, ok := CacheMapSKU[key]
	return sp, ok
}

// Lấy tất cả các dòng SKU thuộc cùng 1 mã Sản Phẩm gốc
func LayNhomSanPham(shopID, maSP string) []*SanPham {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(shopID, maSP)
	if list, ok := CacheGroupSanPham[key]; ok {
		return list
	}
	return []*SanPham{}
}

// Dùng cho trường hợp thêm dòng SKU mới tinh vào RAM (Giữ để xài sau này)
func ThemSanPhamVaoRam(sp *SanPham) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := sp.SpreadsheetID
	
	// 1. Thêm vào List tổng
	CacheSanPham[sID] = append(CacheSanPham[sID], sp)
	
	// 2. Thêm vào Map tra cứu SKU
	keySingle := TaoCompositeKey(sID, sp.LấyIDDuyNhat())
	CacheMapSKU[keySingle] = sp
	
	// 3. Thêm vào Group
	keyGroup := TaoCompositeKey(sID, sp.MaSanPham)
	CacheGroupSanPham[keyGroup] = append(CacheGroupSanPham[keyGroup], sp)
}

func TaoMaSPMoi(shopID, maDanhMuc string) string {
	maDanhMuc = strings.ToUpper(strings.TrimSpace(maDanhMuc))
	if maDanhMuc == "" { maDanhMuc = "SP" }
	slot := LaySlotTiepTheo(shopID, maDanhMuc) 
	return fmt.Sprintf("%s%04d", maDanhMuc, slot)
}
