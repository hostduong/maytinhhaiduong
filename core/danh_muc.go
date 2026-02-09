package core

import (
	"fmt" // [THÊM]
	"app/cau_hinh"
)

// ... (Giữ nguyên Struct và NapDanhMuc) ...

// =============================================================
// 1. CẤU HÌNH CỘT (DANH_MUC)
// =============================================================
const (
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0
	CotDM_TenDanhMuc = 1
	CotDM_HinhAnh    = 2
	CotDM_MoTa       = 3
	CotDM_TrangThai  = 4
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaDanhMuc  string `json:"ma_danh_muc"`
	TenDanhMuc string `json:"ten_danh_muc"`
	HinhAnh    string `json:"hinh_anh"`
	MoTa       string `json:"mo_ta"`
	TrangThai  int    `json:"trang_thai"`
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
		if _, daTonTai := _Map_DanhMuc[key]; daTonTai { continue }

		dm := &DanhMuc{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaDanhMuc:  maDM,
			TenDanhMuc: layString(r, CotDM_TenDanhMuc),
			HinhAnh:    layString(r, CotDM_HinhAnh),
			MoTa:       layString(r, CotDM_MoTa),
			TrangThai:  layInt(r, CotDM_TrangThai),
		}
		_DS_DanhMuc = append(_DS_DanhMuc, dm)
		_Map_DanhMuc[key] = dm
	}
}

func LayDanhSachDanhMuc() []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	var kq []*DanhMuc
	for _, dm := range _DS_DanhMuc {
		if dm.SpreadsheetID == currentSheetID { kq = append(kq, dm) }
	}
	return kq
}

func LayChiTietDanhMuc(maDM string) (*DanhMuc, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	return dm, ok
}

func TaoMaDanhMucMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	for {
		id := fmt.Sprintf("DM%s", LayChuoiSoNgauNhien(3))
		key := TaoCompositeKey(currentSheetID, id)
		if _, tonTai := _Map_DanhMuc[key]; !tonTai { return id }
	}
}

func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if dm.SpreadsheetID == "" { dm.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_DanhMuc = append(_DS_DanhMuc, dm)
	key := TaoCompositeKey(dm.SpreadsheetID, dm.MaDanhMuc)
	_Map_DanhMuc[key] = dm
}
