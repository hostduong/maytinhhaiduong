package may_tinh

import (

	"app/cau_hinh"
	"app/core"
)

// ĐỔI TÊN SHEET THÀNH "MAY_TINH" NHƯ BẠN THIẾT KẾ
const TenSheet = "MAY_TINH" 

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
	CotSP_PhanTramLai    = 20 // U 
	CotSP_GiaNiemYet     = 21 // V
	CotSP_PhanTramGiam   = 22 // W
	CotSP_SoTienGiam     = 23 // X
	CotSP_GiaBan         = 24 // Y
	CotSP_GhiChu         = 25 // Z
	CotSP_NguoiTao       = 26 // AA
	CotSP_NgayTao        = 27 // AB
	CotSP_NguoiCapNhat   = 28 // AC
	CotSP_NgayCapNhat    = 29 // AD
)

// Struct của Ngành Máy Tính (Giữ nguyên 30 cột)
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

func (sp *SanPham) LayIDDuyNhat() string {
	if sp.MaSKU != "" { return sp.MaSKU }
	return sp.MaSanPham
}

var (
	CacheSanPham      = make(map[string][]*SanPham)
	CacheMapSKU       = make(map[string]*SanPham) 
	CacheGroupSanPham = make(map[string][]*SanPham)
)

func NapDuLieu(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	
	// Gọi hàm LoadSheetData từ package core (Viết Hoa chữ L)
	raw, err := core.LoadSheetData(shopID, TenSheet) 
	if err != nil { return }

	list := []*SanPham{}
	groupMap := make(map[string][]*SanPham)

	for i, r := range raw {
		if i < DongBatDau_SanPham-1 { continue }
		
		// Gọi hàm LayString từ package core
		maSP := core.LayString(r, CotSP_MaSanPham)
		if maSP == "" { continue }

		sp := &SanPham{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaSanPham:      maSP,
			TenSanPham:     core.LayString(r, CotSP_TenSanPham),
			TenRutGon:      core.LayString(r, CotSP_TenRutGon),
			Slug:           core.LayString(r, CotSP_Slug),
			MaSKU:          core.LayString(r, CotSP_MaSKU),
			TenSKU:         core.LayString(r, CotSP_TenSKU),
			SKUChinh:       core.LayInt(r, CotSP_SKUChinh),
			TrangThai:      core.LayInt(r, CotSP_TrangThai),
			MaDanhMuc:      core.LayString(r, CotSP_MaDanhMuc),
			MaThuongHieu:   core.LayString(r, CotSP_MaThuongHieu),
			DonVi:          core.LayString(r, CotSP_DonVi),
			MauSac:         core.LayString(r, CotSP_MauSac),
			KhoiLuong:      core.LayFloat(r, CotSP_KhoiLuong),
			KichThuoc:      core.LayString(r, CotSP_KichThuoc),
			UrlHinhAnh:     core.LayString(r, CotSP_UrlHinhAnh),
			ThongSoHTML:    core.LayString(r, CotSP_ThongSoHTML),
			MoTaHTML:       core.LayString(r, CotSP_MoTaHTML),
			BaoHanh:        core.LayString(r, CotSP_BaoHanh),
			TinhTrang:      core.LayString(r, CotSP_TinhTrang),
			GiaNhap:        core.LayFloat(r, CotSP_GiaNhap),
			PhanTramLai:    core.LayFloat(r, CotSP_PhanTramLai),
			GiaNiemYet:     core.LayFloat(r, CotSP_GiaNiemYet),
			PhanTramGiam:   core.LayFloat(r, CotSP_PhanTramGiam),
			SoTienGiam:     core.LayFloat(r, CotSP_SoTienGiam),
			GiaBan:         core.LayFloat(r, CotSP_GiaBan),
			GhiChu:         core.LayString(r, CotSP_GhiChu),
			NguoiTao:       core.LayString(r, CotSP_NguoiTao),
			NgayTao:        core.LayString(r, CotSP_NgayTao),
			NguoiCapNhat:   core.LayString(r, CotSP_NguoiCapNhat),
			NgayCapNhat:    core.LayString(r, CotSP_NgayCapNhat),
		}
		
		list = append(list, sp)
		
		keySingle := core.TaoCompositeKey(shopID, sp.LayIDDuyNhat())
		CacheMapSKU[keySingle] = sp

		keyGroup := core.TaoCompositeKey(shopID, maSP)
		groupMap[keyGroup] = append(groupMap[keyGroup], sp)
	}

	core.KhoaHeThong.Lock()
	CacheSanPham[shopID] = list
	for k, v := range groupMap {
		CacheGroupSanPham[k] = v
	}
	core.KhoaHeThong.Unlock()
}

func LayDanhSachSanPham(shopID string) []*SanPham {
	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()
	return CacheSanPham[shopID]
}

func LayChiTietSKU(shopID, idDuyNhat string) (*SanPham, bool) {
	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()
	key := core.TaoCompositeKey(shopID, idDuyNhat)
	sp, ok := CacheMapSKU[key]
	return sp, ok
}

func TaoMaSPMoi(shopID string, prefix string) string {
	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()
	if prefix == "" { prefix = "SP" }
	
	maxNum := 0
	list := CacheSanPham[shopID]
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
