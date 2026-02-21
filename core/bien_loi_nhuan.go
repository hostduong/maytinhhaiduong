package core

import (
	"app/cau_hinh"
	"sort"
)

const (
	DongBatDau_BienLoiNhuan = 11
	CotBLN_KhungGiaNhap = 0
	CotBLN_BienLoiNhuan = 1
	CotBLN_TrangThai    = 2
)

type BienLoiNhuan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`

	GiaTu          float64 `json:"gia_tu"`
	KhungGiaNhap   float64 `json:"khung_gia_nhap"`
	BienLoiNhuan   float64 `json:"bien_loi_nhuan"`
	TrangThai      int     `json:"trang_thai"`
}

// BỘ NHỚ ĐA SHOP
var (
	CacheBienLoiNhuan = make(map[string][]*BienLoiNhuan)
)

// Helper nội bộ: Cập nhật khoảng giá cho 1 list cụ thể
func capNhatKhoangGia(list []*BienLoiNhuan) {
	var prev float64 = 0
	for _, b := range list {
		b.GiaTu = prev
		prev = b.KhungGiaNhap + 1
	}
}

func NapBienLoiNhuan(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "BIEN_LOI_NHUAN")
	if err != nil { return }

	list := []*BienLoiNhuan{}

	for i, r := range raw {
		if i < DongBatDau_BienLoiNhuan-1 { continue }
		khungGia := LayFloat(r, CotBLN_KhungGiaNhap)
		if khungGia <= 0 { continue } 

		bln := &BienLoiNhuan{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			KhungGiaNhap:   khungGia,
			BienLoiNhuan:   LayFloat(r, CotBLN_BienLoiNhuan),
			TrangThai:      LayInt(r, CotBLN_TrangThai),
		}
		list = append(list, bln)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].KhungGiaNhap < list[j].KhungGiaNhap
	})

	capNhatKhoangGia(list)

	KhoaHeThong.Lock()
	CacheBienLoiNhuan[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachBienLoiNhuan(shopID string) []*BienLoiNhuan {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	if list, ok := CacheBienLoiNhuan[shopID]; ok {
		return list
	}
	return []*BienLoiNhuan{}
}

func ThemBienLoiNhuanVaoRam(bln *BienLoiNhuan) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	sID := bln.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
	
	list := append(CacheBienLoiNhuan[sID], bln)
	
	sort.Slice(list, func(i, j int) bool {
		return list[i].KhungGiaNhap < list[j].KhungGiaNhap
	})
	
	capNhatKhoangGia(list)
	CacheBienLoiNhuan[sID] = list
}

func SuaBienLoiNhuanTrongRam(shopID string, dong int, khungGia, loiNhuan float64, trangThai int) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	
	list := CacheBienLoiNhuan[shopID]
	for _, item := range list {
		if item.DongTrongSheet == dong {
			item.KhungGiaNhap = khungGia
			item.BienLoiNhuan = loiNhuan
			item.TrangThai = trangThai
			break
		}
	}
	
	sort.Slice(list, func(i, j int) bool {
		return list[i].KhungGiaNhap < list[j].KhungGiaNhap
	})
	capNhatKhoangGia(list)
}
