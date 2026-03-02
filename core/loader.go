package core

import (
	"encoding/json"
	"sort"
	"strings"

	"app/config"
)

// --- HÀM TRỢ GIÚP NẠP DATA CHUNG ---
func napDataGeneric(shopID, sheetName string, target interface{}) [][]interface{} {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, sheetName)
	if err != nil { return nil }
	return raw
}

// 1. NẠP PHÂN QUYỀN
func NapPhanQuyen(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetPhanQuyen, nil)
	if raw == nil { return }

	headerIndex, styleIndex := -1, -1
	for i, row := range raw {
		if len(row) > 0 {
			firstCell := strings.TrimSpace(strings.ToLower(LayString(row, 0)))
			if firstCell == "ma_chuc_nang" { headerIndex = i } else if firstCell == "style" { styleIndex = i }
		}
	}
	if headerIndex == -1 { return }

	tempMap := make(map[string]map[string]bool)
	var danhSachVaiTroCuaShop []VaiTroInfo
	header := raw[headerIndex]
	var listMaVaiTro []string

	for i := CotPQ_StartRole; i < len(header); i++ {
		headerText := strings.TrimSpace(LayString(header, i))
		if headerText == "" { continue }
		parts := strings.Split(headerText, "|")
		roleID := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(parts[0])), " ", "_")
		roleName := roleID
		if len(parts) > 1 { roleName = strings.TrimSpace(parts[1]) }

		if roleID != "" {
			listMaVaiTro = append(listMaVaiTro, roleID)
			tempMap[roleID] = make(map[string]bool)
			
			styleCode := 90
			if styleIndex != -1 {
				val := LayInt(raw[styleIndex], i)
				if val >= 0 { styleCode = val }
			}
			
			var lvl, thm int
			if styleCode >= 10 {
				lvl = styleCode / 10
				thm = styleCode % 10
			} else {
				lvl = styleCode
				switch lvl {
				case 0: thm = 9; case 1: thm = 4; case 2: thm = 7; case 3: thm = 5
				case 4: thm = 4; case 5: thm = 6; case 6: thm = 2; case 7: thm = 1
				default: thm = 0
				}
			}

			danhSachVaiTroCuaShop = append(danhSachVaiTroCuaShop, VaiTroInfo{
				MaVaiTro: roleID, TenVaiTro: roleName, StyleLevel: lvl, StyleTheme: thm,
			})
		}
	}

	for i, row := range raw {
		if i <= headerIndex { continue }
		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" || maChucNang == "style" { continue }

		for j, roleID := range listMaVaiTro {
			val := LayString(row, CotPQ_StartRole+j)
			if val == "1" || strings.ToLower(val) == "true" { tempMap[roleID][maChucNang] = true }
		}
	}

	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.Lock()
	defer lock.Unlock()
	CachePhanQuyen[shopID] = tempMap
	CacheDanhSachVaiTro[shopID] = danhSachVaiTroCuaShop
}

// 2. NẠP KHÁCH HÀNG (FULL CỘT TỪ FILE CŨ)
func NapKhachHang(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet } // FIX LỖI CACHE RỖNG TẠI ĐÂY
	raw := napDataGeneric(shopID, TenSheetKhachHang, nil)
	if raw == nil { return }
	list := []*KhachHang{}
	
	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := LayString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }

		kh := &KhachHang{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaKhachHang:    maKH,
			TenDangNhap:    LayString(r, CotKH_TenDangNhap),
			Email:          LayString(r, CotKH_Email),
			MatKhauHash:    LayString(r, CotKH_MatKhauHash),
			MaPinHash:      LayString(r, CotKH_MaPinHash),
			VaiTroQuyenHan: strings.TrimSpace(LayString(r, CotKH_VaiTroQuyenHan)),
			ChucVu:         strings.TrimSpace(LayString(r, CotKH_ChucVu)),
			TrangThai:      LayInt(r, CotKH_TrangThai),
			NguonKhachHang: LayString(r, CotKH_NguonKhachHang),
			TenKhachHang:   LayString(r, CotKH_TenKhachHang),
			DienThoai:      LayString(r, CotKH_DienThoai),
			AnhDaiDien:     LayString(r, CotKH_AnhDaiDien),
			DiaChi:         LayString(r, CotKH_DiaChi),
			NgaySinh:       LayString(r, CotKH_NgaySinh),
			GioiTinh:       LayInt(r, CotKH_GioiTinh),
			MaSoThue:       LayString(r, CotKH_MaSoThue),
			GhiChu:         LayString(r, CotKH_GhiChu),
			NgayTao:        LayString(r, CotKH_NgayTao),
			NguoiCapNhat:   LayString(r, CotKH_NguoiCapNhat),
			NgayCapNhat:    LayString(r, CotKH_NgayCapNhat),
			Inbox:          make([]*TinNhan, 0),
		}

		_ = json.Unmarshal([]byte(LayString(r, CotKH_RefreshTokenJson)), &kh.RefreshTokens)
		if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]TokenInfo) }
		
		_ = json.Unmarshal([]byte(LayString(r, CotKH_DataSheetsJson)), &kh.DataSheets)
		if kh.DataSheets.GoogleAuthJson != "" && kh.DataSheets.SpreadsheetID != "" {
			KetNoiGoogleSheetRieng(kh.DataSheets.SpreadsheetID, kh.DataSheets.GoogleAuthJson)
		}
		_ = json.Unmarshal([]byte(LayString(r, CotKH_GoiDichVuJson)), &kh.GoiDichVu)
		if kh.GoiDichVu == nil { kh.GoiDichVu = make([]PlanInfo, 0) } 
		
		_ = json.Unmarshal([]byte(LayString(r, CotKH_CauHinhJson)), &kh.CauHinh)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_MangXaHoiJson)), &kh.MangXaHoi)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_ViTienJson)), &kh.ViTien)

		list = append(list, kh)
	}

	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.Lock()
	defer lock.Unlock()
	CacheKhachHang[shopID] = list
	for _, kh := range list {
		CacheMapKhachHang[TaoCompositeKey(shopID, kh.MaKhachHang)] = kh
	}
}

// 3. NẠP NHÀ CUNG CẤP (FULL CỘT)
func NapNhaCungCap(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetNhaCungCap, nil)
	if raw == nil { return }
	list := []*NhaCungCap{}
	for i, r := range raw {
		if i < DongBatDau_NhaCungCap-1 { continue }
		maNCC := LayString(r, CotNCC_MaNhaCungCap)
		if maNCC == "" { continue }
		ncc := &NhaCungCap{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaNhaCungCap: maNCC,
			TenNhaCungCap: LayString(r, CotNCC_TenNhaCungCap), MaSoThue: LayString(r, CotNCC_MaSoThue),
			DienThoai: LayString(r, CotNCC_DienThoai), Email: LayString(r, CotNCC_Email),
			KhuVuc: LayString(r, CotNCC_KhuVuc), DiaChi: LayString(r, CotNCC_DiaChi),
			NguoiLienHe: LayString(r, CotNCC_NguoiLienHe), NganHang: LayString(r, CotNCC_NganHang),
			NhomNhaCungCap: LayString(r, CotNCC_NhomNhaCungCap), LoaiNhaCungCap: LayString(r, CotNCC_LoaiNhaCungCap),
			DieuKhoanThanhToan: LayString(r, CotNCC_DieuKhoanThanhToan), ChietKhauMacDinh: LayFloat(r, CotNCC_ChietKhauMacDinh),
			HanMucCongNo: LayFloat(r, CotNCC_HanMucCongNo), CongNoDauKy: LayFloat(r, CotNCC_CongNoDauKy),
			TongMua: LayFloat(r, CotNCC_TongMua), NoCanTra: LayFloat(r, CotNCC_NoCanTra),
			ThongTinThemJson: LayString(r, CotNCC_ThongTinThemJson), TrangThai: LayInt(r, CotNCC_TrangThai),
			GhiChu: LayString(r, CotNCC_GhiChu), NguoiTao: LayString(r, CotNCC_NguoiTao),
			NgayTao: LayString(r, CotNCC_NgayTao), NgayCapNhat: LayString(r, CotNCC_NgayCapNhat),
		}
		list = append(list, ncc)
	}
	lock := GetSheetLock(shopID, TenSheetNhaCungCap)
	lock.Lock()
	defer lock.Unlock()
	CacheNhaCungCap[shopID] = list
	for _, ncc := range list {
		CacheMapNhaCungCap[TaoCompositeKey(shopID, ncc.MaNhaCungCap)] = ncc
	}
}

// 4. NẠP MÁY TÍNH (FULL TỪNG SKUS)
func NapMayTinh(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetMayTinh, nil)
	if raw == nil { return }
	list := []*SanPhamMayTinh{}
	for i, r := range raw {
		if i < DongBatDau_SanPhamMayTinh-1 { continue }
		maSP := LayString(r, CotPC_MaSanPham)
		if maSP == "" { continue }
		sp := &SanPhamMayTinh{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaSanPham: maSP,
			TenSanPham: LayString(r, CotPC_TenSanPham), TenRutGon: LayString(r, CotPC_TenRutGon),
			Slug: LayString(r, CotPC_Slug), MaSKU: LayString(r, CotPC_MaSKU),
			TenSKU: LayString(r, CotPC_TenSKU), SKUChinh: LayInt(r, CotPC_SKUChinh),
			TrangThai: LayInt(r, CotPC_TrangThai), MaDanhMuc: LayString(r, CotPC_MaDanhMuc),
			MaThuongHieu: LayString(r, CotPC_MaThuongHieu), DonVi: LayString(r, CotPC_DonVi),
			MauSac: LayString(r, CotPC_MauSac), KhoiLuong: LayFloat(r, CotPC_KhoiLuong),
			KichThuoc: LayString(r, CotPC_KichThuoc), UrlHinhAnh: LayString(r, CotPC_UrlHinhAnh),
			ThongSoHTML: LayString(r, CotPC_ThongSoHTML), MoTaHTML: LayString(r, CotPC_MoTaHTML),
			BaoHanh: LayString(r, CotPC_BaoHanh), TinhTrang: LayString(r, CotPC_TinhTrang),
			GiaNhap: LayFloat(r, CotPC_GiaNhap), PhanTramLai: LayFloat(r, CotPC_PhanTramLai),
			GiaNiemYet: LayFloat(r, CotPC_GiaNiemYet), PhanTramGiam: LayFloat(r, CotPC_PhanTramGiam),
			SoTienGiam: LayFloat(r, CotPC_SoTienGiam), GiaBan: LayFloat(r, CotPC_GiaBan),
			GhiChu: LayString(r, CotPC_GhiChu), NguoiTao: LayString(r, CotPC_NguoiTao),
			NgayTao: LayString(r, CotPC_NgayTao), NguoiCapNhat: LayString(r, CotPC_NguoiCapNhat),
			NgayCapNhat: LayString(r, CotPC_NgayCapNhat),
		}
		list = append(list, sp)
	}

	lock := GetSheetLock(shopID, TenSheetMayTinh)
	lock.Lock()
	defer lock.Unlock()
	CacheSanPhamMayTinh[shopID] = list
	
	for k := range CacheGroupSanPhamMayTinh {
		if strings.HasPrefix(k, shopID+"__") { delete(CacheGroupSanPhamMayTinh, k) }
	}
	for _, sp := range list {
		CacheMapSKUMayTinh[TaoCompositeKey(shopID, sp.LayIDDuyNhat())] = sp
		kGroup := TaoCompositeKey(shopID, sp.MaSanPham)
		if CacheGroupSanPhamMayTinh[kGroup] == nil {
			CacheGroupSanPhamMayTinh[kGroup] = []*SanPhamMayTinh{}
		}
		CacheGroupSanPhamMayTinh[kGroup] = append(CacheGroupSanPhamMayTinh[kGroup], sp)
	}
}

// 5. NẠP DANH MỤC
func NapDanhMuc(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetDanhMuc, nil)
	if raw == nil { return }
	list := []*DanhMuc{}
	for i, r := range raw {
		if i < DongBatDau_DanhMuc-1 { continue }
		maDM := LayString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }
		dm := &DanhMuc{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaDanhMuc: maDM,
			TenDanhMuc: LayString(r, CotDM_TenDanhMuc), DanhMucMe: LayString(r, CotDM_DanhMucMe),
			ThueVAT: LayFloat(r, CotDM_ThueVAT), LoiNhuan: LayFloat(r, CotDM_LoiNhuan),
			Slot: LayInt(r, CotDM_Slot), TrangThai: LayInt(r, CotDM_TrangThai),
		}
		list = append(list, dm)
	}
	lock := GetSheetLock(shopID, TenSheetDanhMuc)
	lock.Lock(); defer lock.Unlock()
	CacheDanhMuc[shopID] = list
	for _, dm := range list { CacheMapDanhMuc[TaoCompositeKey(shopID, dm.MaDanhMuc)] = dm }
}

// 6. NẠP THƯƠNG HIỆU
func NapThuongHieu(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetThuongHieu, nil)
	if raw == nil { return }
	list := []*ThuongHieu{}
	for i, r := range raw {
		if i < DongBatDau_ThuongHieu-1 { continue }
		maTH := LayString(r, CotTH_MaThuongHieu)
		if maTH == "" { continue }
		th := &ThuongHieu{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaThuongHieu: maTH,
			TenThuongHieu: LayString(r, CotTH_TenThuongHieu), LogoUrl: LayString(r, CotTH_LogoUrl),
			MoTa: LayString(r, CotTH_MoTa), TrangThai: LayInt(r, CotTH_TrangThai),
		}
		list = append(list, th)
	}
	lock := GetSheetLock(shopID, TenSheetThuongHieu)
	lock.Lock(); defer lock.Unlock()
	CacheThuongHieu[shopID] = list
	for _, th := range list { CacheMapThuongHieu[TaoCompositeKey(shopID, th.MaThuongHieu)] = th }
}

// 7. NẠP BIÊN LỢI NHUẬN
func NapBienLoiNhuan(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetBienLoiNhuan, nil)
	if raw == nil { return }
	list := []*BienLoiNhuan{}
	for i, r := range raw {
		if i < DongBatDau_BienLoiNhuan-1 { continue }
		khungGia := LayFloat(r, CotBLN_KhungGiaNhap)
		if khungGia <= 0 { continue }
		bln := &BienLoiNhuan{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, KhungGiaNhap: khungGia,
			BienLoiNhuan: LayFloat(r, CotBLN_BienLoiNhuan), TrangThai: LayInt(r, CotBLN_TrangThai),
		}
		list = append(list, bln)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].KhungGiaNhap < list[j].KhungGiaNhap })
	var prev float64 = 0
	for _, b := range list { b.GiaTu = prev; prev = b.KhungGiaNhap + 1 }
	
	lock := GetSheetLock(shopID, TenSheetBienLoiNhuan)
	lock.Lock(); defer lock.Unlock()
	CacheBienLoiNhuan[shopID] = list
}

// 8. NẠP TIN NHẮN
func NapTinNhan(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheet }
	raw := napDataGeneric(shopID, TenSheetTinNhan, nil)
	if raw == nil { return }
	list := []*TinNhan{}
	for i, r := range raw {
		if i < DongBatDau_TinNhan-1 { continue }
		maTN := LayString(r, CotTN_MaTinNhan)
		if maTN == "" { continue }
		tn := &TinNhan{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaTinNhan: maTN,
			LoaiTinNhan: LayString(r, CotTN_LoaiTinNhan), NguoiGuiID: LayString(r, CotTN_NguoiGuiID),
			NguoiNhanID: LayString(r, CotTN_NguoiNhanID), TieuDe: LayString(r, CotTN_TieuDe),
			NoiDung: LayString(r, CotTN_NoiDung), ThamChieuID: LayString(r, CotTN_ThamChieuID),
			ReplyChoID: LayString(r, CotTN_ReplyChoID), NgayTao: LayString(r, CotTN_NgayTao),
		}
		json.Unmarshal([]byte(LayString(r, CotTN_DinhKemJson)), &tn.DinhKem)
		json.Unmarshal([]byte(LayString(r, CotTN_NguoiDocJson)), &tn.NguoiDoc)
		json.Unmarshal([]byte(LayString(r, CotTN_TrangThaiXoa)), &tn.TrangThaiXoa)
		list = append(list, tn)
	}
	lock := GetSheetLock(shopID, TenSheetTinNhan)
	lock.Lock(); defer lock.Unlock()
	CacheTinNhan[shopID] = list
}
