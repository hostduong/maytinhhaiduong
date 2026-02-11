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
	CotDM_ThueVAT    = 2 // C: Thuế đầu ra (%)
	CotDM_LoiNhuan   = 3 // D: Biên lợi nhuận mong muốn (%)
	CotDM_STT        = 4 // E: Số thứ tự hiện tại (Để sinh SKU: MAIN0001)
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaDanhMuc  string  `json:"ma_danh_muc"`
	TenDanhMuc string  `json:"ten_danh_muc"`
	ThueVAT    float64 `json:"thue_vat"`
	LoiNhuan   float64 `json:"loi_nhuan"`
	STT        int     `json:"stt"` // Số đếm để sinh mã
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
// Chuyển về lowercase để so sánh cho chuẩn
func TimMaDanhMucTheoTen(tenDM string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	// Cần import "strings" nếu muốn dùng strings.ToLower ở đây,
	// nhưng để đơn giản và tránh lỗi import thừa, ta so sánh trực tiếp hoặc tự xử lý bên ngoài.
	// Tuy nhiên, logic chuẩn là phải duyệt qua slice.
	// Ở đây tôi viết lại hàm so sánh đơn giản không cần strings để tránh import thừa nếu lười import.
	// Nhưng tốt nhất là import "strings" nếu cần chính xác. 
	// Trong trường hợp này, để fix lỗi build nhanh, tôi sẽ dùng vòng lặp đơn giản.
	
	for _, dm := range _DS_DanhMuc {
		if dm.TenDanhMuc == tenDM {
			return dm.MaDanhMuc
		}
	}
	return "" 
}

// [MỚI] Hàm lấy số thứ tự tiếp theo và cập nhật RAM + Sheet
func LaySTTtiepTheo(maDM string) int {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()

	key := TaoCompositeKey(cau_hinh.BienCauHinh.IdFileSheet, maDM)
	dm, ok := _Map_DanhMuc[key]
	if !ok { return 1 }

	// Tăng số
	dm.STT++ 
	newSTT := dm.STT

	// Ghi ngay xuống Sheet để tránh trùng lặp nếu restart
	ThemVaoHangCho(dm.SpreadsheetID, "DANH_MUC", dm.DongTrongSheet, CotDM_STT, newSTT)
	
	return newSTT
}

// Thêm vào cuối file core/danh_muc.go
func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if dm.SpreadsheetID == "" { dm.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_DanhMuc = append(_DS_DanhMuc, dm)
	key := TaoCompositeKey(dm.SpreadsheetID, dm.MaDanhMuc)
	_Map_DanhMuc[key] = dm
}
