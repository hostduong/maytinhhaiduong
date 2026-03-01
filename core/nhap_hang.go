package core

import (
	"fmt"
	"strconv"
	"strings"

	"app/cau_hinh"
)

// ====================================================================
// 1. KHAI BÁO TÊN SHEET 
// ====================================================================
const (
	TenSheetPhieuNhap        = "PHIEU_NHAP"
	TenSheetChiTietPhieuNhap = "CHI_TIET_PHIEU_NHAP"
	TenSheetSerial           = "SERIAL_SAN_PHAM"

	DongBatDau_PhieuNhap        = 2
	DongBatDau_ChiTietPhieuNhap = 2
	DongBatDau_Serial           = 2
)

// ====================================================================
// 2. TỌA ĐỘ CỘT
// ====================================================================

// --- PHIẾU NHẬP (18 Cột) ---
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

// ====================================================================
// 4. BỘ NHỚ ĐỆM (RAM CACHE)
// ====================================================================

var (
	CachePhieuNhap      = make(map[string][]*PhieuNhap)
	CacheMapPhieuNhap   = make(map[string]*PhieuNhap)  
)

// ====================================================================
// 5. CÁC HÀM NẠP DỮ LIỆU TỪ GOOGLE SHEETS LÊN RAM
// ====================================================================

func NapPhieuNhap(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	
	rawPN, err := LoadSheetData(shopID, TenSheetPhieuNhap)
	if err != nil { return }

	listPN := []*PhieuNhap{}
	mapPN := make(map[string]*PhieuNhap)

	for i, r := range rawPN {
		if i < DongBatDau_PhieuNhap-1 { continue }
		maPN := LayString(r, CotPN_MaPhieuNhap)
		if maPN == "" { continue }

		pn := &PhieuNhap{
			SpreadsheetID:       shopID,
			DongTrongSheet:      i + 1,
			MaPhieuNhap:         maPN,
			MaNhaCungCap:        LayString(r, CotPN_MaNhaCungCap),
			MaKho:               LayString(r, CotPN_MaKho),
			NgayNhap:            LayString(r, CotPN_NgayNhap),
			TrangThai:           LayInt(r, CotPN_TrangThai),
			SoHoaDon:            LayString(r, CotPN_SoHoaDon),
			NgayHoaDon:          LayString(r, CotPN_NgayHoaDon),
			UrlChungTu:          LayString(r, CotPN_UrlChungTu),
			TongTienPhieu:       LayFloat(r, CotPN_TongTienPhieu),
			GiamGiaPhieu:        LayFloat(r, CotPN_GiamGiaPhieu),
			DaThanhToan:         LayFloat(r, CotPN_DaThanhToan),
			ConNo:               LayFloat(r, CotPN_ConNo),
			PhuongThucThanhToan: LayString(r, CotPN_PhuongThucThanhToan),
			TrangThaiThanhToan:  LayString(r, CotPN_TrangThaiThanhToan),
			GhiChu:              LayString(r, CotPN_GhiChu),
			NguoiTao:            LayString(r, CotPN_NguoiTao),
			NgayTao:             LayString(r, CotPN_NgayTao),
			NgayCapNhat:         LayString(r, CotPN_NgayCapNhat),
			ChiTiet:             []*ChiTietPhieuNhap{}, 
		}
		listPN = append(listPN, pn)
		mapPN[maPN] = pn
		CacheMapPhieuNhap[TaoCompositeKey(shopID, maPN)] = pn
	}

	rawCT, errCT := LoadSheetData(shopID, TenSheetChiTietPhieuNhap)
	if errCT == nil {
		for i, r := range rawCT {
			if i < DongBatDau_ChiTietPhieuNhap-1 { continue }
			maPN := LayString(r, CotCTPN_MaPhieuNhap)
			if maPN == "" { continue }

			ct := &ChiTietPhieuNhap{
				SpreadsheetID:  shopID,
				DongTrongSheet: i + 1,
				MaPhieuNhap:    maPN,
				MaSanPham:      LayString(r, CotCTPN_MaSanPham),
				MaSKU:          LayString(r, CotCTPN_MaSKU),
				MaNganhHang:    LayString(r, CotCTPN_MaNganhHang),
				TenSanPham:     LayString(r, CotCTPN_TenSanPham),
				DonVi:          LayString(r, CotCTPN_DonVi),
				SoLuong:        LayInt(r, CotCTPN_SoLuong),
				DonGiaNhap:     LayFloat(r, CotCTPN_DonGiaNhap),
				VATPercent:     LayFloat(r, CotCTPN_VATPercent),
				GiaSauVAT:      LayFloat(r, CotCTPN_GiaSauVAT),
				ChietKhauDong:  LayFloat(r, CotCTPN_ChietKhauDong),
				ThanhTienDong:  LayFloat(r, CotCTPN_ThanhTienDong),
				GiaVonThucTe:   LayFloat(r, CotCTPN_GiaVonThucTe),
				BaoHanhThang:   LayInt(r, CotCTPN_BaoHanhThang),
				GhiChuDong:     LayString(r, CotCTPN_GhiChuDong),
			}

			if phieu, ok := mapPN[maPN]; ok {
				phieu.ChiTiet = append(phieu.ChiTiet, ct)
			}
		}
	}

	KhoaHeThong.Lock()
	CachePhieuNhap[shopID] = listPN
	KhoaHeThong.Unlock()
}

// ====================================================================
// 6. CÁC HÀM TIỆN ÍCH
// ====================================================================

func LayDanhSachPhieuNhap(shopID string) []*PhieuNhap {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CachePhieuNhap[shopID]
}

func TaoMaPhieuNhapMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	prefix := "PN"
	maxNum := 0
	list := CachePhieuNhap[shopID]
	for _, pn := range list {
		if strings.HasPrefix(pn.MaPhieuNhap, prefix) {
			numStr := strings.TrimPrefix(pn.MaPhieuNhap, prefix)
			if num, err := strconv.Atoi(numStr); err == nil {
				if num > maxNum { maxNum = num }
			}
		}
	}
	return fmt.Sprintf("%s%05d", prefix, maxNum+1) 
}
