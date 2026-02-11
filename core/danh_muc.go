package core

import "app/cau_hinh"

const (
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0 // A: Mã (VD: MAIN, CPU, MON)
	CotDM_TenDanhMuc = 1 // B: Tên hiển thị
	CotDM_DanhMucMe  = 2 // C: Mã danh mục mẹ (Rỗng nếu là cấp 1)
	CotDM_ThueVAT    = 3 // D: Thuế đầu ra (%)
	CotDM_LoiNhuan   = 4 // E: Biên lợi nhuận mong muốn (%)
	CotDM_Slot       = 5 // F: Slot hiện tại (Để sinh SKU: MAIN0001)
	CotDM_TrangThai  = 6 // G
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaDanhMuc  string  `json:"ma_danh_muc"`
	TenDanhMuc string  `json:"ten_danh_muc"`
	DanhMucMe  string  `json:"danh_muc_me"`
	ThueVAT    float64 `json:"thue_vat"`
	LoiNhuan   float64 `json:"bien_loi_nhuan"`
	Slot       int     `json:"slot"`       // [ĐÃ SỬA] Đổi từ STT thành Slot
	TrangThai  int     `json:"trang_thai"` 
}

var (
	_DS_DanhMuc  []*DanhMuc
	_Map_DanhMuc map[string]*DanhMuc
)

func NapDanhMuc(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" { targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(targetSpreadsheetID, "DANH_MUC")
	if err != nil { return }

	_Map_DanhMuc = make(map[string]*DanhMuc)
	_DS_DanhMuc = []*DanhMuc{}

	for i, r := range raw {
		if i < DongBatDau_DanhMuc-1 { continue }
		maDM := layString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }

		key := TaoCompositeKey(targetSpreadsheetID, maDM)
		
		dm := &DanhMuc{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaDanhMuc:      maDM,
			TenDanhMuc:     layString(r, CotDM_TenDanhMuc),
			DanhMucMe:      layString(r, CotDM_DanhMucMe),
			ThueVAT:        layFloat(r, CotDM_ThueVAT),
			LoiNhuan:       layFloat(r, CotDM_LoiNhuan),
			Slot:           layInt(r, CotDM_Slot), // [ĐÃ SỬA]
			TrangThai:      layInt(r, CotDM_TrangThai), 
		}
		_DS_DanhMuc = append(_DS_DanhMuc, dm)
		_Map_DanhMuc[key] = dm
	}
}

func LayDanhSachDanhMuc() []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return _DS_DanhMuc
}

func LayChiTietDanhMuc(maDM string) (*DanhMuc, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	return dm, ok
}

func TimMaDanhMucTheoTen(tenDM string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	for _, dm := range _DS_DanhMuc {
		if dm.TenDanhMuc == tenDM { return dm.MaDanhMuc }
	}
	return "" 
}

// [ĐÃ SỬA] Đổi tên hàm cho chuẩn Slot
func LaySlotTiepTheo(maDM string) int {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()

	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	if !ok { return 1 }

	dm.Slot++ 
	newSlot := dm.Slot
	ThemVaoHangCho(dm.SpreadsheetID, "DANH_MUC", dm.DongTrongSheet, CotDM_Slot, newSlot)
	return newSlot
}

func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if dm.SpreadsheetID == "" { dm.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_DanhMuc = append(_DS_DanhMuc, dm)
	key := TaoCompositeKey(dm.SpreadsheetID, dm.MaDanhMuc)
	_Map_DanhMuc[key] = dm
}
