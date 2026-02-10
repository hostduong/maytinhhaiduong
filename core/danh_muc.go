package core

import (
	"fmt"
	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT
// =============================================================
const (
	DongBatDau_DanhMuc = 11

	// A=0, B=1, C=2, D=3, E=4
	CotDM_MaDanhMuc      = 0 // A
	CotDM_TenDanhMuc     = 1 // B
	CotDM_ThuTuHienThi   = 2 // C
	CotDM_Slug           = 3 // D
	CotDM_MaDanhMucCha   = 4 // E
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaDanhMuc    string `json:"ma_danh_muc"`
	TenDanhMuc   string `json:"ten_danh_muc"`
	Slug         string `json:"slug"` 
	MaDanhMucCha string `json:"ma_danh_muc_cha"`
	TrangThai    int    `json:"trang_thai"` 
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
		// 1. Bỏ qua các dòng tiêu đề
		if i < DongBatDau_DanhMuc-1 { continue }
		
		// 2. Lấy Mã (Cột A)
		maDM := layString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }

		// 3. Kiểm tra trùng lặp
		key := TaoCompositeKey(targetSpreadsheetID, maDM)
		if _, daTonTai := _Map_DanhMuc[key]; daTonTai { continue }

		// 4. Map dữ liệu chuẩn xác theo cấu hình mới
		dm := &DanhMuc{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			
			MaDanhMuc:    maDM,
			TenDanhMuc:   layString(r, CotDM_TenDanhMuc),   // Cột B
			
			// Slug và Danh mục cha
			Slug:         layString(r, CotDM_Slug),         // Cột D
			MaDanhMucCha: layString(r, CotDM_MaDanhMucCha), // Cột E
			
			// Mặc định cho = 1 (Hiện) vì sheet không có cột trạng thái
			TrangThai:    1, 
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

// Hàm hỗ trợ tạo mã mới (giữ nguyên để tương thích)
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

// Hàm hỗ trợ thêm vào RAM (giữ nguyên để tương thích)
func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if dm.SpreadsheetID == "" { dm.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_DanhMuc = append(_DS_DanhMuc, dm)
	key := TaoCompositeKey(dm.SpreadsheetID, dm.MaDanhMuc)
	_Map_DanhMuc[key] = dm
}
