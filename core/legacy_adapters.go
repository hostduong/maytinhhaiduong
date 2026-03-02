package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Alias cho cơ chế hàng đợi ghi mới
var ThemVaoHangCho = PushUpdate

func LayKhachHang(shopID, userID string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	kh, ok := CacheMapKhachHang[TaoCompositeKey(shopID, userID)]
	return kh, ok
}

func LayDanhSachKhachHang(shopID string) []*KhachHang {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheKhachHang[shopID]
}

func TimKhachHangTheoCookie(shopID, cookie string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	for _, kh := range CacheKhachHang[shopID] {
		if info, ok := kh.RefreshTokens[cookie]; ok {
			if time.Now().Unix() <= info.ExpiresAt { return kh, true }
		}
	}
	return nil, false
}

func KiemTraQuyen(shopID, role, maChucNang string) bool {
	if role == "quan_tri_he_thong" { return true }
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.RLock()
	defer lock.RUnlock()
	if shopMap, ok := CachePhanQuyen[shopID]; ok {
		if listQuyen, exists := shopMap[role]; exists {
			if allowed, has := listQuyen[maChucNang]; has && allowed { return true }
		}
	}
	return false
}

func LayCapBacVaiTro(shopID, userID, role string) int {
	if userID == "0000000000000000001" || role == "quan_tri_he_thong" { return 0 }
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.RLock()
	defer lock.RUnlock()
	for _, v := range CacheDanhSachVaiTro[shopID] {
		if v.MaVaiTro == role { return v.StyleLevel }
	}
	return 9
}

func LayDanhSachSanPhamMayTinh(shopID string) []*SanPhamMayTinh {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheSanPhamMayTinh[shopID]
}

func LayChiTietSKUMayTinh(shopID, id string) (*SanPhamMayTinh, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	sp, ok := CacheMapSKUMayTinh[TaoCompositeKey(shopID, id)]
	if !ok {
		for _, s := range CacheSanPhamMayTinh[shopID] {
			if s.Slug == id && s.TrangThai == 1 { return s, true }
		}
	}
	return sp, ok
}

func TaoMaSPMayTinhMoi(shopID, prefix string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	if prefix == "" { prefix = "SP" }
	max := 0
	for _, sp := range CacheSanPhamMayTinh[shopID] {
		if strings.HasPrefix(sp.MaSanPham, prefix) {
			numStr := strings.TrimPrefix(sp.MaSanPham, prefix)
			if num, err := strconv.Atoi(numStr); err == nil && num > max { max = num }
		}
	}
	return fmt.Sprintf("%s%04d", prefix, max+1)
}

func CapNhatSlotThuCong(shopID, dmMa string, slotMoi int) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	for _, dm := range CacheDanhMuc[shopID] {
		if dm.MaDanhMuc == dmMa {
			if slotMoi > dm.Slot {
				dm.Slot = slotMoi
				ThemVaoHangCho(shopID, TenSheetDanhMuc, dm.DongTrongSheet, CotDM_Slot, slotMoi)
			}
			break
		}
	}
}

func LayDanhSachDanhMuc(shopID string) []*DanhMuc {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheDanhMuc[shopID]
}

func LayDanhSachThuongHieu(shopID string) []*ThuongHieu {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheThuongHieu[shopID]
}

func LayDanhSachBienLoiNhuan(shopID string) []*BienLoiNhuan {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheBienLoiNhuan[shopID]
}

func LayDanhSachNhaCungCap(shopID string) []*NhaCungCap {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	return CacheNhaCungCap[shopID]
}

func LayHopThuNguoiDung(shopID, userID, role string) []*TinNhan {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	var rs []*TinNhan
	for _, tn := range CacheTinNhan[shopID] {
		if tn.NguoiNhanID == userID || tn.NguoiGuiID == userID { rs = append(rs, tn); continue }
		if tn.LoaiTinNhan == "ALL" { rs = append(rs, tn); continue }
		if tn.LoaiTinNhan == "ROLE" && strings.Contains(tn.NguoiNhanID, role) { rs = append(rs, tn); continue }
	}
	return rs
}

func ThemMoiTinNhan(shopID string, tn *TinNhan) {
	KhoaHeThong.Lock()
	tn.DongTrongSheet = DongBatDau_TinNhan + len(CacheTinNhan[shopID])
	CacheTinNhan[shopID] = append(CacheTinNhan[shopID], tn)
	KhoaHeThong.Unlock()
	rowData := []interface{}{
		tn.MaTinNhan, tn.LoaiTinNhan, tn.NguoiGuiID, tn.NguoiNhanID, tn.TieuDe, tn.NoiDung, "", tn.ThamChieuID, tn.ReplyChoID, tn.NgayTao, "[]", "",
	}
	PushAppend(shopID, TenSheetTinNhan, rowData)
}

func DanhDauDocTinNhan(shopID, userID, msgID string) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	for _, tn := range CacheTinNhan[shopID] {
		if tn.MaTinNhan == msgID {
			daDoc := false
			for _, u := range tn.NguoiDoc { if u == userID { daDoc = true; break } }
			if !daDoc {
				tn.NguoiDoc = append(tn.NguoiDoc, userID)
				b, _ := json.Marshal(tn.NguoiDoc)
				ThemVaoHangCho(shopID, TenSheetTinNhan, tn.DongTrongSheet, CotTN_NguoiDocJson, string(b))
			}
			break
		}
	}
}

func ToJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
