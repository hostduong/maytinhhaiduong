package core

import (
	"sort"
	"app/cau_hinh"
)

// 1. CẤU HÌNH CỘT
const (
	CotDM_MaDanhMuc    = 0
	CotDM_ThuTuHienThi = 1
	CotDM_TenDanhMuc   = 2
	CotDM_Slug         = 3
	CotDM_MaDanhMucCha = 4
)

// 2. STRUCT
type DanhMuc struct {
	SpreadsheetID string `json:"-"`
	
	MaDanhMuc    string `json:"ma_danh_muc"`
	ThuTuHienThi int    `json:"thu_tu_hien_thi"`
	TenDanhMuc   string `json:"ten_danh_muc"`
	Slug         string `json:"slug"`
	MaDanhMucCha string `json:"ma_danh_muc_cha"`
}

// 3. KHO LƯU TRỮ
var (
	_DS_DanhMuc  []*DanhMuc
	_Map_DanhMuc map[string]*DanhMuc
)

// 4. NẠP DỮ LIỆU
func NapDanhMuc(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" {
		targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	raw, err := loadSheetData(targetSpreadsheetID, "DANH_MUC")
	if err != nil { return }

	if _Map_DanhMuc == nil {
		_Map_DanhMuc = make(map[string]*DanhMuc)
		_DS_DanhMuc = []*DanhMuc{}
	}
	_DS_DanhMuc = []*DanhMuc{} // Reset tạm thời

	for i, r := range raw {
		if i < DongBatDauDuLieu-1 { continue }
		
		ma := layString(r, CotDM_MaDanhMuc)
		if ma == "" { continue }

		item := &DanhMuc{
			SpreadsheetID: targetSpreadsheetID,
			MaDanhMuc:     ma,
			ThuTuHienThi:  layInt(r, CotDM_ThuTuHienThi),
			TenDanhMuc:    layString(r, CotDM_TenDanhMuc),
			Slug:          layString(r, CotDM_Slug),
			MaDanhMucCha:  layString(r, CotDM_MaDanhMucCha),
		}

		_DS_DanhMuc = append(_DS_DanhMuc, item)
		key := TaoCompositeKey(targetSpreadsheetID, ma)
		_Map_DanhMuc[key] = item
	}
}

// 5. TRUY VẤN
func LayDanhSachDanhMuc() []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()

	currentID := cau_hinh.BienCauHinh.IdFileSheet
	var kq []*DanhMuc
	
	for _, dm := range _DS_DanhMuc {
		if dm.SpreadsheetID == currentID {
			kq = append(kq, dm)
		}
	}

	// Sắp xếp theo Thứ tự hiển thị
	sort.Slice(kq, func(i, j int) bool {
		return kq[i].ThuTuHienThi < kq[j].ThuTuHienThi
	})
	
	return kq
}

func LayDanhMuc(maDM string) (*DanhMuc, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	return dm, ok
}
