package core

import (
	"fmt"
	"strconv"
	"strings"

	"app/cau_hinh"
)

// ====================================================================
// 1. KHAI BÁO TÊN SHEET (ĐÃ TÁCH RIÊNG CHO NGÀNH MÁY TÍNH)
// ====================================================================
const (
	TenSheetPhieuNhapMayTinh        = "PHIEU_NHAP_MAY_TINH"
	TenSheetCTPhieuNhapMayTinh      = "CT_PHIEU_NHAP_MAY_TINH"
	TenSheetNhaCungCapMayTinh       = "NHA_CUNG_CAP_MAY_TINH"
	TenSheetSerialMayTinh           = "SERIAL_MAY_TINH"

	DongBatDau_PhieuNhap        = 2
	DongBatDau_ChiTietPhieuNhap = 2
	DongBatDau_NhaCungCap       = 2
	DongBatDau_Serial           = 2
)

// ====================================================================
// 2. KHAI BÁO TỌA ĐỘ CỘT
// ====================================================================

// --- PHIẾU NHẬP ---
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
	CotPN_DaThanhToan          = 9  // J
	CotPN_ConNo                = 10 // K
	CotPN_PhuongThucThanhToan  = 11 // L
	CotPN_TrangThaiThanhToan   = 12 // M
	CotPN_GhiChu               = 13 // N
	CotPN_NguoiTao             = 14 // O
	CotPN_NgayTao              = 15 // P
	CotPN_NgayCapNhat          = 16 // Q
)

// --- CHI TIẾT PHIẾU NHẬP ---
const (
	CotCTPN_MaPhieuNhap     = 0  // A
	CotCTPN_MaSanPham       = 1  // B
	CotCTPN_TenSanPham      = 2  // C
	CotCTPN_DonVi           = 3  // D
	CotCTPN_SoLuong         = 4  // E
	CotCTPN_DonGiaNhap      = 5  // F
	CotCTPN_VATPercent      = 6  // G
	CotCTPN_GiaSauVAT       = 7  // H
	CotCTPN_ChietKhauDong   = 8  // I
	CotCTPN_ThanhTienDong   = 9  // J
	CotCTPN_GiaVonThucTe    = 10 // K
	CotCTPN_BaoHanhThang    = 11 // L
	CotCTPN_GhiChuDong      = 12 // M
)

// --- NHÀ CUNG CẤP ---
const (
	CotNCC_MaNhaCungCap     = 0  // A
	CotNCC_TenNhaCungCap    = 1  // B
	CotNCC_DienThoai        = 2  // C
	CotNCC_Email            = 3  // D
	CotNCC_DiaChi           = 4  // E
	CotNCC_MaSoThue         = 5  // F
	CotNCC_NguoiLienHe      = 6  // G
	CotNCC_NganHang         = 7  // H
	CotNCC_NoCanTra         = 8  // I
	CotNCC_TongMua          = 9  // J
	CotNCC_HanMucCongNo     = 10 // K
	CotNCC_TrangThai        = 11 // L
	CotNCC_GhiChu           = 12 // M
	CotNCC_NguoiTao         = 13 // N
	CotNCC_NgayTao          = 14 // O
	CotNCC_NgayCapNhat      = 15 // P
)

// ====================================================================
// 3. KHAI BÁO STRUCT (CÓ HẬU TỐ MAYTINH)
// ====================================================================

type NhaCungCapMayTinh struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaNhaCungCap   string  `json:"ma_nha_cung_cap"`
	TenNhaCungCap  string  `json:"ten_nha_cung_cap"`
	DienThoai      string  `json:"dien_thoai"`
	Email          string  `json:"email"`
	DiaChi         string  `json:"dia_chi"`
	MaSoThue       string  `json:"ma_so_thue"`
	NguoiLienHe    string  `json:"nguoi_lien_he"`
	NganHang       string  `json:"ngan_hang"`
	NoCanTra       float64 `json:"no_can_tra"`
	TongMua        float64 `json:"tong_mua"`
	HanMucCongNo   float64 `json:"han_muc_cong_no"`
	TrangThai      int     `json:"trang_thai"`
	GhiChu         string  `json:"ghi_chu"`
}

type PhieuNhapMayTinh struct {
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
	DaThanhToan          float64 `json:"da_thanh_toan"`
	ConNo                float64 `json:"con_no"`
	PhuongThucThanhToan  string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiThanhToan   string  `json:"trang_thai_thanh_toan"`
	GhiChu               string  `json:"ghi_chu"`
	
	// Slice chứa chi tiết để Front-end dễ render
	ChiTiet              []*ChiTietPhieuNhapMayTinh `json:"chi_tiet"`
}

type ChiTietPhieuNhapMayTinh struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaPhieuNhap    string  `json:"ma_phieu_nhap"`
	MaSanPham      string  `json:"ma_san_pham"`
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

// ====================================================================
// 4. BỘ NHỚ ĐỆM (RAM CACHE)
// ====================================================================

var (
	CacheNhaCungCapMayTinh     = make(map[string][]*NhaCungCapMayTinh)
	CacheMapNhaCungCapMayTinh  = make(map[string]*NhaCungCapMayTinh) 

	CachePhieuNhapMayTinh      = make(map[string][]*PhieuNhapMayTinh)
	CacheMapPhieuNhapMayTinh   = make(map[string]*PhieuNhapMayTinh)  
)

// ====================================================================
// 5. CÁC HÀM NẠP DỮ LIỆU TỪ GOOGLE SHEETS LÊN RAM
// ====================================================================

func NapNhaCungCapMayTinh(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, TenSheetNhaCungCapMayTinh)
	if err != nil { return }

	list := []*NhaCungCapMayTinh{}
	
	for i, r := range raw {
		if i < DongBatDau_NhaCungCap-1 { continue }
		maNCC := LayString(r, CotNCC_MaNhaCungCap)
		if maNCC == "" { continue }

		ncc := &NhaCungCapMayTinh{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaNhaCungCap:   maNCC,
			TenNhaCungCap:  LayString(r, CotNCC_TenNhaCungCap),
			DienThoai:      LayString(r, CotNCC_DienThoai),
			Email:          LayString(r, CotNCC_Email),
			DiaChi:         LayString(r, CotNCC_DiaChi),
			MaSoThue:       LayString(r, CotNCC_MaSoThue),
			NguoiLienHe:    LayString(r, CotNCC_NguoiLienHe),
			NganHang:       LayString(r, CotNCC_NganHang),
			NoCanTra:       LayFloat(r, CotNCC_NoCanTra),
			TongMua:        LayFloat(r, CotNCC_TongMua),
			HanMucCongNo:   LayFloat(r, CotNCC_HanMucCongNo),
			TrangThai:      LayInt(r, CotNCC_TrangThai),
			GhiChu:         LayString(r, CotNCC_GhiChu),
		}
		list = append(list, ncc)
		CacheMapNhaCungCapMayTinh[TaoCompositeKey(shopID, maNCC)] = ncc
	}

	KhoaHeThong.Lock()
	CacheNhaCungCapMayTinh[shopID] = list
	KhoaHeThong.Unlock()
}

func NapPhieuNhapMayTinh(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	
	rawPN, err := LoadSheetData(shopID, TenSheetPhieuNhapMayTinh)
	if err != nil { return }

	listPN := []*PhieuNhapMayTinh{}
	mapPN := make(map[string]*PhieuNhapMayTinh)

	for i, r := range rawPN {
		if i < DongBatDau_PhieuNhap-1 { continue }
		maPN := LayString(r, CotPN_MaPhieuNhap)
		if maPN == "" { continue }

		pn := &PhieuNhapMayTinh{
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
			DaThanhToan:         LayFloat(r, CotPN_DaThanhToan),
			ConNo:               LayFloat(r, CotPN_ConNo),
			PhuongThucThanhToan: LayString(r, CotPN_PhuongThucThanhToan),
			TrangThaiThanhToan:  LayString(r, CotPN_TrangThaiThanhToan),
			GhiChu:              LayString(r, CotPN_GhiChu),
			ChiTiet:             []*ChiTietPhieuNhapMayTinh{}, 
		}
		listPN = append(listPN, pn)
		mapPN[maPN] = pn
		CacheMapPhieuNhapMayTinh[TaoCompositeKey(shopID, maPN)] = pn
	}

	rawCT, errCT := LoadSheetData(shopID, TenSheetCTPhieuNhapMayTinh)
	if errCT == nil {
		for i, r := range rawCT {
			if i < DongBatDau_ChiTietPhieuNhap-1 { continue }
			maPN := LayString(r, CotCTPN_MaPhieuNhap)
			if maPN == "" { continue }

			ct := &ChiTietPhieuNhapMayTinh{
				SpreadsheetID:  shopID,
				DongTrongSheet: i + 1,
				MaPhieuNhap:    maPN,
				MaSanPham:      LayString(r, CotCTPN_MaSanPham),
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
	CachePhieuNhapMayTinh[shopID] = listPN
	KhoaHeThong.Unlock()
}

// ====================================================================
// 6. CÁC HÀM TIỆN ÍCH (GETTERS & GENERATORS)
// ====================================================================

func LayDanhSachNhaCungCapMayTinh(shopID string) []*NhaCungCapMayTinh {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheNhaCungCapMayTinh[shopID]
}

func LayDanhSachPhieuNhapMayTinh(shopID string) []*PhieuNhapMayTinh {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CachePhieuNhapMayTinh[shopID]
}

func TaoMaPhieuNhapMayTinhMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	prefix := "PN"
	maxNum := 0
	list := CachePhieuNhapMayTinh[shopID]
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

func TaoMaNhaCungCapMayTinhMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	prefix := "NCC"
	maxNum := 0
	list := CacheNhaCungCapMayTinh[shopID]
	for _, ncc := range list {
		if strings.HasPrefix(ncc.MaNhaCungCap, prefix) {
			numStr := strings.TrimPrefix(ncc.MaNhaCungCap, prefix)
			if num, err := strconv.Atoi(numStr); err == nil {
				if num > maxNum { maxNum = num }
			}
		}
	}
	return fmt.Sprintf("%s%03d", prefix, maxNum+1) 
}
