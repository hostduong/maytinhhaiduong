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

// BỘ NHỚ ĐA SHOP
var (
	CacheThuongHieu    = make(map[string][]*ThuongHieu)
	CacheMapThuongHieu = make(map[string]map[string]*ThuongHieu)
)

func NapThuongHieu(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(shopID, "THUONG_HIEU")
	if err != nil { return }

	list := []*ThuongHieu{}

	for i, r := range raw {
		if i < DongBatDau_ThuongHieu-1 { continue }
		maTH := layString(r, CotTH_MaThuongHieu)
		if maTH == "" { continue }

		th := &ThuongHieu{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaThuongHieu:  maTH,
			TenThuongHieu: layString(r, CotTH_TenThuongHieu),
			LogoUrl:       layString(r, CotTH_LogoUrl),
			MoTa:          layString(r, CotTH_MoTa),
			TrangThai:     layInt(r, CotTH_TrangThai),
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
	
	if list, ok := CacheThuongHieu[shopID]; ok {
		return list
	}
	return []*ThuongHieu{}
}

func ThemThuongHieuVaoRam(th *ThuongHieu) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := th.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
	
	CacheThuongHieu[sID] = append(CacheThuongHieu[sID], th)
	
	key := TaoCompositeKey(sID, th.MaThuongHieu)
	CacheMapThuongHieu[key] = th
}
