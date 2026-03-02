package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var ThemVaoHangCho = PushUpdate

// =======================================================
// KHÁCH HÀNG & NHÂN SỰ (Dùng khóa TenSheetKhachHang)
// =======================================================
func LayKhachHang(shopID, userID string) (*KhachHang, bool) {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	kh, ok := CacheMapKhachHang[TaoCompositeKey(shopID, userID)]
	return kh, ok
}

func LayDanhSachKhachHang(shopID string) []*KhachHang {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	return CacheKhachHang[shopID]
}

func TimKhachHangTheoCookie(shopID, cookie string) (*KhachHang, bool) {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	for _, kh := range CacheKhachHang[shopID] {
		if info, ok := kh.RefreshTokens[cookie]; ok {
			if time.Now().Unix() <= info.ExpiresAt { return kh, true }
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(shopID, input string) (*KhachHang, bool) {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range CacheKhachHang[shopID] {
		if strings.ToLower(kh.TenDangNhap) == input || (kh.Email != "" && strings.ToLower(kh.Email) == input) {
			if kh.MaKhachHang == "0000000000000000000" { return nil, false } 
			return kh, true
		}
	}
	return nil, false
}

func TaoMaKhachHangMoi(shopID string) string {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	for {
		id := LayChuoiSoNgauNhien(19)
		if _, exist := CacheMapKhachHang[TaoCompositeKey(shopID, id)]; !exist { return id }
	}
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	sID := kh.SpreadsheetID
	if sID == "" { sID = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8" }
	lock := GetSheetLock(sID, TenSheetKhachHang)
	lock.Lock(); defer lock.Unlock()
	
	CacheKhachHang[sID] = append(CacheKhachHang[sID], kh)
	CacheMapKhachHang[TaoCompositeKey(sID, kh.MaKhachHang)] = kh
}

// =======================================================
// PHÂN QUYỀN & HỆ THỐNG (Dùng khóa TenSheetPhanQuyen)
// =======================================================
func LayCapBacVaiTro(shopID, userID, role string) int {
	if userID == "0000000000000000001" || role == "quan_tri_he_thong" { return 0 }
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.RLock(); defer lock.RUnlock()
	for _, v := range CacheDanhSachVaiTro[shopID] {
		if v.MaVaiTro == role { return v.StyleLevel }
	}
	return 9
}

func KiemTraQuyen(shopID, role, maChucNang string) bool {
	if role == "quan_tri_he_thong" { return true }
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.RLock(); defer lock.RUnlock()
	if shopMap, ok := CachePhanQuyen[shopID]; ok {
		if listQuyen, exists := shopMap[role]; exists {
			if allowed, has := listQuyen[maChucNang]; has && allowed { return true }
		}
	}
	return false
}

// =======================================================
// SẢN PHẨM & MASTER DATA (Khóa phân mảnh theo bảng)
// =======================================================
func TaoMaSPMayTinhMoi(shopID, prefix string) string {
	lock := GetSheetLock(shopID, TenSheetMayTinh)
	lock.RLock(); defer lock.RUnlock()
	if prefix == "" { prefix = "SP" }
	max := 0
	for _, sp := range CacheSanPhamMayTinh[shopID] {
		if strings.HasPrefix(sp.MaSanPham, prefix) {
			var num int; fmt.Sscanf(strings.TrimPrefix(sp.MaSanPham, prefix), "%d", &num)
			if num > max { max = num }
		}
	}
	return fmt.Sprintf("%s%04d", prefix, max+1)
}

func CapNhatSlotThuCong(shopID, dmMa string, slotMoi int) {
	lock := GetSheetLock(shopID, TenSheetDanhMuc)
	lock.Lock(); defer lock.Unlock()
	for _, dm := range CacheDanhMuc[shopID] {
		if dm.MaDanhMuc == dmMa {
			if slotMoi > dm.Slot { dm.Slot = slotMoi; PushUpdate(shopID, TenSheetDanhMuc, dm.DongTrongSheet, CotDM_Slot, slotMoi) }
			break
		}
	}
}

func LayDanhSachSanPhamMayTinh(shopID string) []*SanPhamMayTinh { lock := GetSheetLock(shopID, TenSheetMayTinh); lock.RLock(); defer lock.RUnlock(); return CacheSanPhamMayTinh[shopID] }
func LayDanhSachDanhMuc(shopID string) []*DanhMuc { lock := GetSheetLock(shopID, TenSheetDanhMuc); lock.RLock(); defer lock.RUnlock(); return CacheDanhMuc[shopID] }
func LayDanhSachThuongHieu(shopID string) []*ThuongHieu { lock := GetSheetLock(shopID, TenSheetThuongHieu); lock.RLock(); defer lock.RUnlock(); return CacheThuongHieu[shopID] }
func LayDanhSachBienLoiNhuan(shopID string) []*BienLoiNhuan { lock := GetSheetLock(shopID, TenSheetBienLoiNhuan); lock.RLock(); defer lock.RUnlock(); return CacheBienLoiNhuan[shopID] }
func LayDanhSachNhaCungCap(shopID string) []*NhaCungCap { lock := GetSheetLock(shopID, TenSheetNhaCungCap); lock.RLock(); defer lock.RUnlock(); return CacheNhaCungCap[shopID] }

func LayChiTietSKUMayTinh(shopID, id string) (*SanPhamMayTinh, bool) {
	lock := GetSheetLock(shopID, TenSheetMayTinh)
	lock.RLock(); defer lock.RUnlock()
	sp, ok := CacheMapSKUMayTinh[TaoCompositeKey(shopID, id)]
	if !ok { for _, s := range CacheSanPhamMayTinh[shopID] { if s.Slug == id && s.TrangThai == 1 { return s, true } } }
	return sp, ok
}

// =======================================================
// TIN NHẮN (Dùng khóa TenSheetTinNhan)
// =======================================================
func LayHopThuNguoiDung(shopID, userID, role string) []*TinNhan {
	lock := GetSheetLock(shopID, TenSheetTinNhan)
	lock.RLock(); defer lock.RUnlock()
	var rs []*TinNhan
	for _, tn := range CacheTinNhan[shopID] {
		if tn.NguoiNhanID == userID || tn.NguoiGuiID == userID { rs = append(rs, tn); continue }
		if tn.LoaiTinNhan == "ALL" { rs = append(rs, tn); continue }
		if tn.LoaiTinNhan == "ROLE" && strings.Contains(tn.NguoiNhanID, role) { rs = append(rs, tn); continue }
	}
	return rs
}

func ThemMoiTinNhan(shopID string, tn *TinNhan) {
	lock := GetSheetLock(shopID, TenSheetTinNhan)
	lock.Lock()
	tn.DongTrongSheet = DongBatDau_TinNhan + len(CacheTinNhan[shopID])
	CacheTinNhan[shopID] = append(CacheTinNhan[shopID], tn)
	lock.Unlock()
	
	PushAppend(shopID, TenSheetTinNhan, []interface{}{ tn.MaTinNhan, tn.LoaiTinNhan, tn.NguoiGuiID, tn.NguoiNhanID, tn.TieuDe, tn.NoiDung, "", tn.ThamChieuID, tn.ReplyChoID, tn.NgayTao, "[]", "" })
}

func DanhDauDocTinNhan(shopID, userID, msgID string) {
	lock := GetSheetLock(shopID, TenSheetTinNhan)
	lock.Lock(); defer lock.Unlock()
	for _, tn := range CacheTinNhan[shopID] {
		if tn.MaTinNhan == msgID {
			daDoc := false; for _, u := range tn.NguoiDoc { if u == userID { daDoc = true; break } }
			if !daDoc { tn.NguoiDoc = append(tn.NguoiDoc, userID); b, _ := json.Marshal(tn.NguoiDoc); PushUpdate(shopID, TenSheetTinNhan, tn.DongTrongSheet, CotTN_NguoiDocJson, string(b)) }
			break
		}
	}
}

func ToJSON(v interface{}) string { b, _ := json.Marshal(v); return string(b) }
