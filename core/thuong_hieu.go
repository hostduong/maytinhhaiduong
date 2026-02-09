package core

import (
	"fmt" // [THÊM]
	"app/cau_hinh"
)

// ... (Giữ nguyên Struct và NapThuongHieu) ...

const (
	DongBatDau_ThuongHieu = 11 
	CotTH_MaThuongHieu  = 0
	CotTH_TenThuongHieu = 1
	CotTH_HinhAnh       = 2
	CotTH_MoTa          = 3
	CotTH_TrangThai     = 4
)

type ThuongHieu struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaThuongHieu  string `json:"ma_thuong_hieu"`
	TenThuongHieu string `json:"ten_thuong_hieu"`
	HinhAnh       string `json:"hinh_anh"`
	MoTa          string `json:"mo_ta"`
	TrangThai     int    `json:"trang_thai"`
}

var (
	_DS_ThuongHieu  []*ThuongHieu
	_Map_ThuongHieu map[string]*ThuongHieu
)

func NapThuongHieu(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" { targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(targetSpreadsheetID, "THUONG_HIEU")
	if err != nil { return }

	_Map_ThuongHieu = make(map[string]*ThuongHieu)
	_DS_ThuongHieu = []*ThuongHieu{}

	for i, r := range raw {
		if i < DongBatDau_ThuongHieu-1 { continue }
		maTH := layString(r, CotTH_MaThuongHieu)
		if maTH == "" { continue }

		key := TaoCompositeKey(targetSpreadsheetID, maTH)
		if _, daTonTai := _Map_ThuongHieu[key]; daTonTai { continue }

		th := &ThuongHieu{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaThuongHieu:  maTH,
			TenThuongHieu: layString(r, CotTH_TenThuongHieu),
			HinhAnh:       layString(r, CotTH_HinhAnh),
			MoTa:          layString(r, CotTH_MoTa),
			TrangThai:     layInt(r, CotTH_TrangThai),
		}
		_DS_ThuongHieu = append(_DS_ThuongHieu, th)
		_Map_ThuongHieu[key] = th
	}
}

func LayDanhSachThuongHieu() []*ThuongHieu {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	var kq []*ThuongHieu
	for _, th := range _DS_ThuongHieu {
		if th.SpreadsheetID == currentSheetID { kq = append(kq, th) }
	}
	return kq
}

func LayChiTietThuongHieu(maTH string) (*ThuongHieu, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maTH)
	th, ok := _Map_ThuongHieu[key]
	return th, ok
}

func TaoMaThuongHieuMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	for {
		id := fmt.Sprintf("TH%s", LayChuoiSoNgauNhien(3))
		key := TaoCompositeKey(currentSheetID, id)
		if _, tonTai := _Map_ThuongHieu[key]; !tonTai { return id }
	}
}

func ThemThuongHieuVaoRam(th *ThuongHieu) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if th.SpreadsheetID == "" { th.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_ThuongHieu = append(_DS_ThuongHieu, th)
	key := TaoCompositeKey(th.SpreadsheetID, th.MaThuongHieu)
	_Map_ThuongHieu[key] = th
}
