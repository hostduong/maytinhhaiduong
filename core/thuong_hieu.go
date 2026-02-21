package core

import "app/cau_hinh"

const (
	DongBatDau_ThuongHieu = 11 
	CotTH_MaThuongHieu  = 0
	CotTH_TenThuongHieu = 1
	CotTH_LogoUrl       = 2
	CotTH_MoTa          = 3
	CotTH_TrangThai     = 4
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
	CacheThuongHieu    = make(map[string][]*ThuongHieu)
	CacheMapThuongHieu = make(map[string]*ThuongHieu) // Đã sửa thành Map phẳng
)

func NapThuongHieu(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "THUONG_HIEU")
	if err != nil { return }

	list := []*ThuongHieu{}
	for i, r := range raw {
		if i < DongBatDau_ThuongHieu-1 { continue }
		maTH := LayString(r, CotTH_MaThuongHieu)
		if maTH == "" { continue }

		th := &ThuongHieu{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaThuongHieu:  maTH,
			TenThuongHieu: LayString(r, CotTH_TenThuongHieu),
			LogoUrl:       LayString(r, CotTH_LogoUrl),
			MoTa:          LayString(r, CotTH_MoTa),
			TrangThai:     LayInt(r, CotTH_TrangThai),
		}
		list = append(list, th)
		key := TaoCompositeKey(shopID, maTH)
		CacheMapThuongHieu[key] = th
	}
	KhoaHeThong.Lock()
	CacheThuongHieu[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachThuongHieu(shopID string) []*ThuongHieu {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheThuongHieu[shopID]
}

func ThemThuongHieuVaoRam(th *ThuongHieu) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	sID := th.SpreadsheetID
	CacheThuongHieu[sID] = append(CacheThuongHieu[sID], th)
	key := TaoCompositeKey(sID, th.MaThuongHieu)
	CacheMapThuongHieu[key] = th
}
