package core

import "app/cau_hinh"

const (
	DongBatDau_ThuongHieu = 11 
	CotTH_MaThuongHieu  = 0 // A
	CotTH_TenThuongHieu = 1 // B
	CotTH_LogoUrl       = 2 // C
	CotTH_MoTa          = 3 // D
	CotTH_TrangThai     = 4 // E
)

type ThuongHieu struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaThuongHieu  string `json:"ma_thuong_hieu"`
	TenThuongHieu string `json:"ten_thuong_hieu"`
	LogoUrl       string `json:"logo_url"`
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
		th := &ThuongHieu{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaThuongHieu:  maTH,
			TenThuongHieu: layString(r, CotTH_TenThuongHieu),
			LogoUrl:       layString(r, CotTH_LogoUrl),
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
	return _DS_ThuongHieu
}

func ThemThuongHieuVaoRam(th *ThuongHieu) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if th.SpreadsheetID == "" { th.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_ThuongHieu = append(_DS_ThuongHieu, th)
	key := TaoCompositeKey(th.SpreadsheetID, th.MaThuongHieu)
	_Map_ThuongHieu[key] = th
}
