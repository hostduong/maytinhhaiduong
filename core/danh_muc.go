package core

import (
	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT (DANH_MUC)
// =============================================================
const (
	// [CHUẨN HÓA]
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0
	CotDM_TenDanhMuc = 1
	CotDM_HinhAnh    = 2
	CotDM_MoTa       = 3
	CotDM_TrangThai  = 4
)

// =============================================================
// 2. STRUCT DỮ LIỆU
// =============================================================
type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	MaDanhMuc  string `json:"ma_danh_muc"`
	TenDanhMuc string `json:"ten_danh_muc"`
	HinhAnh    string `json:"hinh_anh"`
	MoTa       string `json:"mo_ta"`
	TrangThai  int    `json:"trang_thai"`
}

// =============================================================
// 3. KHO LƯU TRỮ
// =============================================================
var (
	_DS_DanhMuc  []*DanhMuc
	_Map_DanhMuc map[string]*DanhMuc
)

// =============================================================
// 4. LOGIC NẠP DỮ LIỆU
// =============================================================
func NapDanhMuc(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" {
		targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	raw, err := loadSheetData(targetSpreadsheetID, "DANH_MUC")
	if err != nil { return }

	_Map_DanhMuc = make(map[string]*DanhMuc)
	_DS_DanhMuc = []*DanhMuc{}

	for i, r := range raw {
		// [CHUẨN HÓA] Dùng biến DongBatDau_DanhMuc
		if i < DongBatDau_DanhMuc-1 { continue }
		
		maDM := layString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }

		key := TaoCompositeKey(targetSpreadsheetID, maDM)

		// [AN TOÀN] Chống trùng lặp
		if _, daTonTai := _Map_DanhMuc[key]; daTonTai {
			continue
		}

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

// =============================================================
// 5. TRUY VẤN
// =============================================================
func LayDanhSachDanhMuc() []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	currentSheetID := cau_hinh.BienCauHinh.IdFileSheet
	var kq []*DanhMuc

	for _, dm := range _DS_DanhMuc {
		if dm.SpreadsheetID == currentSheetID {
			kq = append(kq, dm)
		}
	}
	return kq
}
