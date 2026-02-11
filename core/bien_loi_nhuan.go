package core

import (
	"app/cau_hinh"
	"sort"
)

const (
	DongBatDau_BienLoiNhuan = 11

	CotBLN_KhungGiaNhap = 0 // A: Giá nhập nhỏ hơn hoặc bằng
	CotBLN_BienLoiNhuan = 1 // B: Lợi nhuận (%)
	CotBLN_TrangThai    = 2 // C: 1=Hoạt động, 0=Tắt
)

type BienLoiNhuan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	KhungGiaNhap float64 `json:"khung_gia_nhap"`
	BienLoiNhuan float64 `json:"bien_loi_nhuan"`
	TrangThai    int     `json:"trang_thai"`
}

var (
	_DS_BienLoiNhuan []*BienLoiNhuan
)

func NapBienLoiNhuan(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" { targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(targetSpreadsheetID, "BIEN_LOI_NHUAN")
	if err != nil { return }

	_DS_BienLoiNhuan = []*BienLoiNhuan{}

	for i, r := range raw {
		if i < DongBatDau_BienLoiNhuan-1 { continue }
		khungGia := layFloat(r, CotBLN_KhungGiaNhap)
		if khungGia <= 0 { continue } // Bỏ qua dòng trống

		bln := &BienLoiNhuan{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			KhungGiaNhap:   khungGia,
			BienLoiNhuan:   layFloat(r, CotBLN_BienLoiNhuan),
			TrangThai:      layInt(r, CotBLN_TrangThai),
		}
		_DS_BienLoiNhuan = append(_DS_BienLoiNhuan, bln)
	}

	// Sắp xếp tự động tăng dần theo Khung giá để lúc tính toán dễ dò tìm
	sort.Slice(_DS_BienLoiNhuan, func(i, j int) bool {
		return _DS_BienLoiNhuan[i].KhungGiaNhap < _DS_BienLoiNhuan[j].KhungGiaNhap
	})
}

func LayDanhSachBienLoiNhuan() []*BienLoiNhuan {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return _DS_BienLoiNhuan
}

func ThemBienLoiNhuanVaoRam(bln *BienLoiNhuan) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if bln.SpreadsheetID == "" { bln.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_BienLoiNhuan = append(_DS_BienLoiNhuan, bln)
	
	sort.Slice(_DS_BienLoiNhuan, func(i, j int) bool {
		return _DS_BienLoiNhuan[i].KhungGiaNhap < _DS_BienLoiNhuan[j].KhungGiaNhap
	})
}
