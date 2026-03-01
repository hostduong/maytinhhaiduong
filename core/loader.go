package core

import (
	"encoding/json"
	"strings"
)

// --- HÀM TRỢ GIÚP NẠP DATA CHUNG ---
func napDataGeneric(shopID, sheetName string, target interface{}) [][]interface{} {
	if shopID == "" { shopID = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8" }
	raw, err := LoadSheetData(shopID, sheetName)
	if err != nil { return nil }
	return raw
}

// 1. NẠP PHÂN QUYỀN
func NapPhanQuyen(shopID string) {
	raw := napDataGeneric(shopID, TenSheetPhanQuyen, nil)
	if raw == nil { return }
	headerIndex := -1
	for i, row := range raw {
		if len(row) > 0 && strings.TrimSpace(strings.ToLower(LayString(row, 0))) == "ma_chuc_nang" {
			headerIndex = i; break
		}
	}
	if headerIndex == -1 { return }
	tempMap := make(map[string]map[string]bool)
	var listVaiTro []VaiTroInfo
	header := raw[headerIndex]
	var listMaVaiTro []string
	for i := CotPQ_StartRole; i < len(header); i++ {
		text := strings.TrimSpace(LayString(header, i))
		if text == "" { continue }
		roleID := strings.ReplaceAll(strings.ToLower(text), " ", "_")
		listMaVaiTro = append(listMaVaiTro, roleID)
		tempMap[roleID] = make(map[string]bool)
		listVaiTro = append(listVaiTro, VaiTroInfo{MaVaiTro: roleID, TenVaiTro: text, StyleLevel: 9})
	}
	for i, row := range raw {
		if i <= headerIndex { continue }
		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" { continue }
		for j, roleID := range listMaVaiTro {
			if LayString(row, CotPQ_StartRole+j) == "1" { tempMap[roleID][maChucNang] = true }
		}
	}
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.Lock(); defer lock.Unlock()
	CachePhanQuyen[shopID] = tempMap
	CacheDanhSachVaiTro[shopID] = listVaiTro
}

// 2. NẠP KHÁCH HÀNG
func NapKhachHang(shopID string) {
	raw := napDataGeneric(shopID, TenSheetKhachHang, nil)
	if raw == nil { return }
	list := []*KhachHang{}
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.Lock(); defer lock.Unlock()
	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := LayString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }
		kh := &KhachHang{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaKhachHang: maKH,
			TenDangNhap: LayString(r, CotKH_TenDangNhap), Email: LayString(r, CotKH_Email),
			MatKhauHash: LayString(r, CotKH_MatKhauHash), MaPinHash: LayString(r, CotKH_MaPinHash),
			VaiTroQuyenHan: LayString(r, CotKH_VaiTroQuyenHan), ChucVu: LayString(r, CotKH_ChucVu),
			TrangThai: LayInt(r, CotKH_TrangThai), TenKhachHang: LayString(r, CotKH_TenKhachHang),
			DienThoai: LayString(r, CotKH_DienThoai), AnhDaiDien: LayString(r, CotKH_AnhDaiDien),
		}
		json.Unmarshal([]byte(LayString(r, CotKH_RefreshTokenJson)), &kh.RefreshTokens)
		json.Unmarshal([]byte(LayString(r, CotKH_MangXaHoiJson)), &kh.MangXaHoi)
		list = append(list, kh)
		CacheMapKhachHang[shopID+"__"+maKH] = kh
	}
	CacheKhachHang[shopID] = list
}

// 3. NẠP NHÀ CUNG CẤP
func NapNhaCungCap(shopID string) {
	raw := napDataGeneric(shopID, TenSheetNhaCungCap, nil)
	if raw == nil { return }
	list := []*NhaCungCap{}
	lock := GetSheetLock(shopID, TenSheetNhaCungCap)
	lock.Lock(); defer lock.Unlock()
	for i, r := range raw {
		if i < DongBatDau_NhaCungCap-1 { continue }
		maNCC := LayString(r, CotNCC_MaNhaCungCap)
		if maNCC == "" { continue }
		ncc := &NhaCungCap{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaNhaCungCap: maNCC,
			TenNhaCungCap: LayString(r, CotNCC_TenNhaCungCap), NoCanTra: LayFloat(r, CotNCC_NoCanTra),
			TrangThai: LayInt(r, CotNCC_TrangThai), CongNoDauKy: LayFloat(r, CotNCC_CongNoDauKy),
		}
		list = append(list, ncc)
		CacheMapNhaCungCap[shopID+"__"+maNCC] = ncc
	}
	CacheNhaCungCap[shopID] = list
}

// 4. NẠP SẢN PHẨM MÁY TÍNH
func NapMayTinh(shopID string) {
	raw := napDataGeneric(shopID, TenSheetMayTinh, nil)
	if raw == nil { return }
	list := []*SanPhamMayTinh{}
	lock := GetSheetLock(shopID, TenSheetMayTinh)
	lock.Lock(); defer lock.Unlock()
	for i, r := range raw {
		if i < 11-1 { continue }
		maSP := LayString(r, 0)
		if maSP == "" { continue }
		sp := &SanPhamMayTinh{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaSanPham: maSP,
			TenSanPham: LayString(r, 1), MaSKU: LayString(r, 4), SKUChinh: LayInt(r, 6),
			TrangThai: LayInt(r, 7), GiaBan: LayFloat(r, 24),
		}
		list = append(list, sp)
		CacheMapSKUMayTinh[shopID+"__"+sp.LayIDDuyNhat()] = sp
	}
	CacheSanPhamMayTinh[shopID] = list
}

// Các hàm stub cho các sheet khác để tránh lỗi compile
func NapDanhMuc(s string) {}
func NapThuongHieu(s string) {}
func NapBienLoiNhuan(s string) {}
func NapTinNhan(s string) {}
