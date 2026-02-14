package core

import (
	"strings"
	"app/cau_hinh"
)

const (
	DongBatDau_DanhMuc = 11

	CotDM_MaDanhMuc  = 0
	CotDM_TenDanhMuc = 1
	CotDM_DanhMucMe  = 2
	CotDM_ThueVAT    = 3
	CotDM_LoiNhuan   = 4
	CotDM_Slot       = 5
	CotDM_TrangThai  = 6
)

type DanhMuc struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	
	MaDanhMuc  string  `json:"ma_danh_muc"`
	TenDanhMuc string  `json:"ten_danh_muc"`
	DanhMucMe  string  `json:"danh_muc_me"`
	ThueVAT    float64 `json:"thue_vat"`
	LoiNhuan   float64 `json:"bien_loi_nhuan"`
	Slot       int     `json:"slot"`
	TrangThai  int     `json:"trang_thai"` 
}

// BỘ NHỚ ĐA SHOP
var (
	CacheDanhMuc    = make(map[string][]*DanhMuc)
	// Sửa lại thành map phẳng để dùng với TaoCompositeKey
	CacheMapDanhMuc = make(map[string]*DanhMuc) 
)

func NapDanhMuc(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(shopID, "DANH_MUC")
	if err != nil { return }

	list := []*DanhMuc{}

	for i, r := range raw {
		if i < DongBatDau_DanhMuc-1 { continue }
		maDM := layString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }

		dm := &DanhMuc{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaDanhMuc:      maDM,
			TenDanhMuc:     layString(r, CotDM_TenDanhMuc),
			DanhMucMe:      layString(r, CotDM_DanhMucMe),
			ThueVAT:        layFloat(r, CotDM_ThueVAT),
			LoiNhuan:       layFloat(r, CotDM_LoiNhuan),
			Slot:           layInt(r, CotDM_Slot),
			TrangThai:      layInt(r, CotDM_TrangThai), 
		}
		
		list = append(list, dm)
		
		// Map lookup: Dùng key kết hợp ShopID__MaDM
		key := TaoCompositeKey(shopID, maDM)
		CacheMapDanhMuc[key] = dm
	}

	KhoaHeThong.Lock()
	CacheDanhMuc[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachDanhMuc(shopID string) []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	if list, ok := CacheDanhMuc[shopID]; ok {
		return list
	}
	return []*DanhMuc{}
}

func LayChiTietDanhMuc(shopID, maDM string) (*DanhMuc, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	key := TaoCompositeKey(shopID, maDM)
	dm, ok := CacheMapDanhMuc[key]
	return dm, ok
}

func TimMaDanhMucTheoTen(shopID, tenDM string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	
	list := CacheDanhMuc[shopID]
	for _, dm := range list {
		if strings.EqualFold(dm.TenDanhMuc, tenDM) { return dm.MaDanhMuc }
	}
	return "" 
}

func LaySlotTiepTheo(shopID, maDM string) int {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()

	key := TaoCompositeKey(shopID, maDM)
	dm, ok := CacheMapDanhMuc[key]
	
	if !ok { return 1 }

	dm.Slot++ 
	newSlot := dm.Slot
	
	ThemVaoHangCho(dm.SpreadsheetID, "DANH_MUC", dm.DongTrongSheet, CotDM_Slot, newSlot)
	return newSlot
}

// --- THÊM HÀM MỚI: CẬP NHẬT SLOT CƯỠNG BỨC (Dùng cho logic đồng bộ ngược) ---
func CapNhatSlotThuCong(shopID, maDM string, slotMoi int) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()

	key := TaoCompositeKey(shopID, maDM)
	dm, ok := CacheMapDanhMuc[key]
	
	// Chỉ cập nhật nếu tìm thấy danh mục VÀ số mới lớn hơn số cũ
	if ok && slotMoi > dm.Slot {
		dm.Slot = slotMoi
		// Ghi đè số lớn nhất này vào Sheet (Vẫn dùng hàng đợi an toàn)
		ThemVaoHangCho(dm.SpreadsheetID, "DANH_MUC", dm.DongTrongSheet, CotDM_Slot, slotMoi)
	}
}

func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := dm.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
	
	CacheDanhMuc[sID] = append(CacheDanhMuc[sID], dm)
	
	key := TaoCompositeKey(sID, dm.MaDanhMuc)
	CacheMapDanhMuc[key] = dm
}
