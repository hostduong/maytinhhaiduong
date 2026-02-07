package core

import (
	"sort"
	"app/cau_hinh"
)

// 1. CẤU HÌNH CỘT
const (
	CotTH_MaThuongHieu  = 0
	CotTH_TenThuongHieu = 1
	CotTH_LogoUrl       = 2
)

// 2. STRUCT
type ThuongHieu struct {
	SpreadsheetID string `json:"-"`

	MaThuongHieu  string `json:"ma_thuong_hieu"`
	TenThuongHieu string `json:"ten_thuong_hieu"`
	LogoUrl       string `json:"logo_url"`
}

// 3. KHO LƯU TRỮ
var (
	_DS_ThuongHieu  []*ThuongHieu
	_Map_ThuongHieu map[string]*ThuongHieu
)

// 4. NẠP DỮ LIỆU
func NapThuongHieu(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" {
		targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	raw, err := loadSheetData(targetSpreadsheetID, "THUONG_HIEU")
	if err != nil { return }

	if _Map_ThuongHieu == nil {
		_Map_ThuongHieu = make(map[string]*ThuongHieu)
		_DS_ThuongHieu = []*ThuongHieu{}
	}
	_DS_ThuongHieu = []*ThuongHieu{}

	for i, r := range raw {
		if i < DongBatDauDuLieu-1 { continue }
		
		ma := layString(r, CotTH_MaThuongHieu)
		if ma == "" { continue }

		item := &ThuongHieu{
			SpreadsheetID: targetSpreadsheetID,
			MaThuongHieu:  ma,
			TenThuongHieu: layString(r, CotTH_TenThuongHieu),
			LogoUrl:       layString(r, CotTH_LogoUrl),
		}

		_DS_ThuongHieu = append(_DS_ThuongHieu, item)
		key := TaoCompositeKey(targetSpreadsheetID, ma)
		_Map_ThuongHieu[key] = item
	}
}

// 5. TRUY VẤN
func LayDanhSachThuongHieu() []*ThuongHieu {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()

	currentID := cau_hinh.BienCauHinh.IdFileSheet
	var kq []*ThuongHieu
	
	for _, th := range _DS_ThuongHieu {
		if th.SpreadsheetID == currentID {
			kq = append(kq, th)
		}
	}

	// Sắp xếp A-Z
	sort.Slice(kq, func(i, j int) bool {
		return kq[i].TenThuongHieu < kq[j].TenThuongHieu
	})
	
	return kq
}

func LayThuongHieu(maTH string) (*ThuongHieu, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maTH)
	th, ok := _Map_ThuongHieu[key]
	return th, ok
}
