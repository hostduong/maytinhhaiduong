package core

import (
	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT (DANH_MUC)
// =============================================================
const (
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0 // A: Mã (VD: MAIN, CPU, MON)
	CotDM_TenDanhMuc = 1 // B: Tên hiển thị
	CotDM_DanhMucMe  = 2 // C: [MỚI] Mã danh mục mẹ (Rỗng nếu là cấp 1)
	CotDM_ThueVAT    = 3 // D: Thuế đầu ra (%)
	CotDM_LoiNhuan   = 4 // E: Biên lợi nhuận mong muốn (%)
	CotDM_STT        = 5 // F: Số thứ tự hiện tại (Để sinh SKU: MAIN0001)
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaDanhMuc  string  `json:"ma_danh_muc"`
	TenDanhMuc string  `json:"ten_danh_muc"`
	DanhMucMe  string  `json:"danh_muc_me"` // [MỚI]
	ThueVAT    float64 `json:"thue_vat"`
	LoiNhuan   float64 `json:"loi_nhuan"`
	STT        int     `json:"stt"` 
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
			MaDanhMuc:  maDM,
			TenDanhMuc: layString(r, CotDM_TenDanhMuc),
			DanhMucMe:  layString(r, CotDM_DanhMucMe), // [MỚI]
			ThueVAT:    layFloat(r, CotDM_ThueVAT),
			LoiNhuan:   layFloat(r, CotDM_LoiNhuan),
			STT:        layInt(r, CotDM_STT),
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
		if dm.TenDanhMuc == tenDM {
			return dm.MaDanhMuc
		}
	}
	return "" 
}

func LaySTTtiepTheo(maDM string) int {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()

	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	if !ok { return 1 }

	dm.STT++ 
	newSTT := dm.STT

	// Ghi ngay xuống Sheet 
	ThemVaoHangCho(dm.SpreadsheetID, "DANH_MUC", dm.DongTrongSheet, CotDM_STT, newSTT)
	
	return newSTT
}

func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if dm.SpreadsheetID == "" { dm.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_DanhMuc = append(_DS_DanhMuc, dm)
	key := TaoCompositeKey(dm.SpreadsheetID, dm.MaDanhMuc)
	_Map_DanhMuc[key] = dm
}
