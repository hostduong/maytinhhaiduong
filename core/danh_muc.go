package core

import (
	"fmt"
	"strings"
	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT (DANH_MUC)
// =============================================================
const (
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0 // A: Mã (VD: MAIN)
	CotDM_TenDanhMuc = 1 // B: Tên (VD: Mainboard)
	CotDM_ThueVAT    = 2 // C
	CotDM_LoiNhuan   = 3 // D
	CotDM_STT        = 4 // E: Đếm số (VD: 105)
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaDanhMuc  string  `json:"ma_danh_muc"`
	TenDanhMuc string  `json:"ten_danh_muc"`
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

// [MỚI] Tìm Mã (MAIN) dựa trên Tên (Mainboard)
// Vì giao diện gửi tên lên, nên ta phải tìm ngược lại mã
func TimMaDanhMucTheoTen(tenDM string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	tenDM = strings.ToLower(strings.TrimSpace(tenDM))
	for _, dm := range _DS_DanhMuc {
		if strings.ToLower(dm.TenDanhMuc) == tenDM {
			return dm.MaDanhMuc
		}
	}
	return "" // Không tìm thấy
}

// [MỚI] Lấy STT tiếp theo và cập nhật luôn vào Sheet
func LaySTTtiepTheo(maDM string) int {
	KhoaHeThong.Lock() // Lock ghi
	defer KhoaHeThong.Unlock()

	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	
	// Nếu danh mục chưa có, hoặc lỗi, trả về 1
	if !ok { return 1 }

	// Tăng số đếm
	dm.STT++ 
	newSTT := dm.STT

	// Ghi ngay xuống Sheet cột E (Index 4) để lưu lại trạng thái
	ThemVaoHangCho(dm.SpreadsheetID, "DANH_MUC", dm.DongTrongSheet, CotDM_STT, newSTT)
	
	return newSTT
}
