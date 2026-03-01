package core

import (
	"app/cau_hinh"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ==============================================================================
// PHẦN 1: DANH MỤC
// ==============================================================================
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
	MaDanhMuc      string `json:"ma_danh_muc"`
	TenDanhMuc     string `json:"ten_danh_muc"`
	DanhMucMe      string `json:"danh_muc_me"`
	ThueVAT        float64 `json:"thue_vat"`
	LoiNhuan       float64 `json:"bien_loi_nhuan"`
	Slot           int    `json:"slot"`
	TrangThai      int    `json:"trang_thai"`
}

var (
	CacheDanhMuc    = make(map[string][]*DanhMuc)
	CacheMapDanhMuc = make(map[string]*DanhMuc)
)

func NapDanhMuc(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "DANH_MUC")
	if err != nil { return }
	list := []*DanhMuc{}
	for i, r := range raw {
		if i < DongBatDau_DanhMuc-1 { continue }
		maDM := LayString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }
		dm := &DanhMuc{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaDanhMuc:      maDM,
			TenDanhMuc:     LayString(r, CotDM_TenDanhMuc),
			DanhMucMe:      LayString(r, CotDM_DanhMucMe),
			ThueVAT:        LayFloat(r, CotDM_ThueVAT),
			LoiNhuan:       LayFloat(r, CotDM_LoiNhuan),
			Slot:           LayInt(r, CotDM_Slot),
			TrangThai:      LayInt(r, CotDM_TrangThai),
		}
		list = append(list, dm)
		CacheMapDanhMuc[TaoCompositeKey(shopID, maDM)] = dm
	}
	KhoaHeThong.Lock()
	CacheDanhMuc[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachDanhMuc(shopID string) []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheDanhMuc[shopID]
}

func LayChiTietDanhMuc(shopID, maDM string) (*DanhMuc, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	dm, ok := CacheMapDanhMuc[TaoCompositeKey(shopID, maDM)]
	return dm, ok
}

func ThemDanhMucVaoRam(dm *DanhMuc) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	sID := dm.SpreadsheetID
	CacheDanhMuc[sID] = append(CacheDanhMuc[sID], dm)
	CacheMapDanhMuc[TaoCompositeKey(sID, dm.MaDanhMuc)] = dm
}

// ==============================================================================
// PHẦN 2: THƯƠNG HIỆU
// ==============================================================================
const (
	DongBatDau_ThuongHieu = 11
	CotTH_MaThuongHieu    = 0
	CotTH_TenThuongHieu   = 1
	CotTH_LogoUrl         = 2
	CotTH_MoTa            = 3
	CotTH_TrangThai       = 4
)

type ThuongHieu struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaThuongHieu   string `json:"ma_thuong_hieu"`
	TenThuongHieu  string `json:"ten_thuong_hieu"`
	LogoUrl        string `json:"logo_url"`
	MoTa           string `json:"mo_ta"`
	TrangThai      int    `json:"trang_thai"`
}

var (
	CacheThuongHieu    = make(map[string][]*ThuongHieu)
	CacheMapThuongHieu = make(map[string]*ThuongHieu)
)

func NapThuongHieu(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "THUONG_HIEU")
	if err != nil { return }
	list := []*ThuongHieu{}
	for i, r := range raw {
		if i < DongBatDau_ThuongHieu-1 { continue }
		maTH := LayString(r, CotTH_MaThuongHieu)
		if maTH == "" { continue }
		th := &ThuongHieu{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaThuongHieu:   maTH,
			TenThuongHieu:  LayString(r, CotTH_TenThuongHieu),
			LogoUrl:        LayString(r, CotTH_LogoUrl),
			MoTa:           LayString(r, CotTH_MoTa),
			TrangThai:      LayInt(r, CotTH_TrangThai),
		}
		list = append(list, th)
		CacheMapThuongHieu[TaoCompositeKey(shopID, maTH)] = th
	}
	KhoaHeThong.Lock()
	CacheThuongHieu[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachThuongHieu(shopID string) []*ThuongHieu {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheThuongHieu[shopID]
}

func ThemThuongHieuVaoRam(th *ThuongHieu) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	sID := th.SpreadsheetID
	CacheThuongHieu[sID] = append(CacheThuongHieu[sID], th)
	CacheMapThuongHieu[TaoCompositeKey(sID, th.MaThuongHieu)] = th
}

// ==============================================================================
// PHẦN 3: BIÊN LỢI NHUẬN
// ==============================================================================
const (
	DongBatDau_BienLoiNhuan = 11
	CotBLN_KhungGiaNhap     = 0
	CotBLN_BienLoiNhuan     = 1
	CotBLN_TrangThai        = 2
)

type BienLoiNhuan struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	GiaTu          float64 `json:"gia_tu"`
	KhungGiaNhap   float64 `json:"khung_gia_nhap"`
	BienLoiNhuan   float64 `json:"bien_loi_nhuan"`
	TrangThai      int     `json:"trang_thai"`
}

var CacheBienLoiNhuan = make(map[string][]*BienLoiNhuan)

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
	sort.Slice(list, func(i, j int) bool { return list[i].KhungGiaNhap < list[j].KhungGiaNhap })
	capNhatKhoangGia(list)
	KhoaHeThong.Lock()
	CacheBienLoiNhuan[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachBienLoiNhuan(shopID string) []*BienLoiNhuan {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheBienLoiNhuan[shopID]
}

func ThemBienLoiNhuanVaoRam(bln *BienLoiNhuan) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	sID := bln.SpreadsheetID
	list := append(CacheBienLoiNhuan[sID], bln)
	sort.Slice(list, func(i, j int) bool { return list[i].KhungGiaNhap < list[j].KhungGiaNhap })
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
	sort.Slice(list, func(i, j int) bool { return list[i].KhungGiaNhap < list[j].KhungGiaNhap })
	capNhatKhoangGia(list)
}

// ==============================================================================
// PHẦN 4: NHÀ CUNG CẤP
// ==============================================================================
const (
	TenSheetNhaCungCap    = "NHA_CUNG_CAP"
	DongBatDau_NhaCungCap = 2

	CotNCC_MaNhaCungCap = 0
	CotNCC_TenNhaCungCap = 1
	CotNCC_DienThoai    = 2
	CotNCC_Email        = 3
	CotNCC_DiaChi       = 4
	CotNCC_MaSoThue     = 5
	CotNCC_NguoiLienHe  = 6
	CotNCC_NganHang     = 7
	CotNCC_NoCanTra     = 8
	CotNCC_TongMua      = 9
	CotNCC_HanMucCongNo = 10
	CotNCC_TrangThai    = 11
	CotNCC_GhiChu       = 12
	CotNCC_NguoiTao     = 13
	CotNCC_NgayTao      = 14
	CotNCC_NgayCapNhat  = 15
)

type NhaCungCap struct {
	SpreadsheetID  string `json:"-"`
	DongTrongSheet int    `json:"-"`
	MaNhaCungCap   string `json:"ma_nha_cung_cap"`
	TenNhaCungCap  string `json:"ten_nha_cung_cap"`
	DienThoai      string `json:"dien_thoai"`
	Email          string `json:"email"`
	DiaChi         string `json:"dia_chi"`
	MaSoThue       string `json:"ma_so_thue"`
	NguoiLienHe    string `json:"nguoi_lien_he"`
	NganHang       string `json:"ngan_hang"`
	NoCanTra       float64 `json:"no_can_tra"`
	TongMua        float64 `json:"tong_mua"`
	HanMucCongNo   float64 `json:"han_muc_cong_no"`
	TrangThai      int    `json:"trang_thai"`
	GhiChu         string `json:"ghi_chu"`
}

var (
	CacheNhaCungCap    = make(map[string][]*NhaCungCap)
	CacheMapNhaCungCap = make(map[string]*NhaCungCap)
)

func NapNhaCungCap(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, TenSheetNhaCungCap)
	if err != nil { return }
	list := []*NhaCungCap{}
	for i, r := range raw {
		if i < DongBatDau_NhaCungCap-1 { continue }
		maNCC := LayString(r, CotNCC_MaNhaCungCap)
		if maNCC == "" { continue }
		ncc := &NhaCungCap{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaNhaCungCap:   maNCC,
			TenNhaCungCap:  LayString(r, CotNCC_TenNhaCungCap),
			DienThoai:      LayString(r, CotNCC_DienThoai),
			Email:          LayString(r, CotNCC_Email),
			DiaChi:         LayString(r, CotNCC_DiaChi),
			MaSoThue:       LayString(r, CotNCC_MaSoThue),
			NguoiLienHe:    LayString(r, CotNCC_NguoiLienHe),
			NganHang:       LayString(r, CotNCC_NganHang),
			NoCanTra:       LayFloat(r, CotNCC_NoCanTra),
			TongMua:        LayFloat(r, CotNCC_TongMua),
			HanMucCongNo:   LayFloat(r, CotNCC_HanMucCongNo),
			TrangThai:      LayInt(r, CotNCC_TrangThai),
			GhiChu:         LayString(r, CotNCC_GhiChu),
		}
		list = append(list, ncc)
		CacheMapNhaCungCap[TaoCompositeKey(shopID, maNCC)] = ncc
	}
	KhoaHeThong.Lock()
	CacheNhaCungCap[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachNhaCungCap(shopID string) []*NhaCungCap {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheNhaCungCap[shopID]
}

func TaoMaNhaCungCapMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	prefix := "NCC"
	maxNum := 0
	for _, ncc := range CacheNhaCungCap[shopID] {
		if strings.HasPrefix(ncc.MaNhaCungCap, prefix) {
			numStr := strings.TrimPrefix(ncc.MaNhaCungCap, prefix)
			if num, err := strconv.Atoi(numStr); err == nil {
				if num > maxNum { maxNum = num }
			}
		}
	}
	return fmt.Sprintf("%s%03d", prefix, maxNum+1)
}

// ==============================================================================
// PHẦN 5: PHÂN QUYỀN (Nếu có, bạn paste thêm code phân quyền vào đây)
// ==============================================================================
