package may_tinh_master

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app/core"
)

func Service_LayChiTietMayTinh(adminShopID string, maSP string) ([]*core.SanPhamMayTinh, error) {
	core.KhoaHeThong.RLock()
	listSKU := core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(adminShopID, maSP)]
	core.KhoaHeThong.RUnlock()

	if len(listSKU) == 0 {
		return nil, errors.New("Không tìm thấy sản phẩm!")
	}
	return listSKU, nil
}

func Service_LuuMayTinh(masterShopID, adminShopID, vaiTro, userID, maSP, dataJSON string) error {
	if maSP == "" {
		if vaiTro != "quan_tri_he_thong" && !core.KiemTraQuyen(masterShopID, vaiTro, "product.create") {
			return errors.New("Bạn không có quyền thêm sản phẩm!")
		}
	} else {
		if vaiTro != "quan_tri_he_thong" && !core.KiemTraQuyen(masterShopID, vaiTro, "product.edit") {
			return errors.New("Bạn không có quyền sửa sản phẩm!")
		}
	}

	var inputSKUs []InputSKUMayTinh
	if err := json.Unmarshal([]byte(dataJSON), &inputSKUs); err != nil || len(inputSKUs) == 0 {
		return errors.New("Dữ liệu không hợp lệ!")
	}

	hasMain := false
	for _, sku := range inputSKUs { if sku.SKUChinh == 1 { hasMain = true; break } }
	if !hasMain { inputSKUs[0].SKUChinh = 1 }

	loc := time.FixedZone("ICT", 7*3600)
	nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")

	core.KhoaHeThong.RLock()
	existingSKUs := core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(adminShopID, maSP)] 
	core.KhoaHeThong.RUnlock()

	if maSP == "" {
		firstCodeDM := ""
		if inputSKUs[0].MaDanhMuc != "" { 
			parsedDM := Repo_XuLyTags(inputSKUs[0].MaDanhMuc)
			if parsedDM != "" { firstCodeDM = strings.Split(parsedDM, "|")[0] }
		}
		maSP = core.TaoMaSPMayTinhMoi(adminShopID, firstCodeDM) 
	} else {
		if len(existingSKUs) == 0 {
			firstCodeDM := ""
			if inputSKUs[0].MaDanhMuc != "" { 
				parsedDM := Repo_XuLyTags(inputSKUs[0].MaDanhMuc)
				if parsedDM != "" { firstCodeDM = strings.Split(parsedDM, "|")[0] }
			}
			re := regexp.MustCompile(`[0-9]+`)
			nums := re.FindAllString(maSP, -1)
			if len(nums) > 0 {
				lastNumStr := nums[len(nums)-1] 
				if slotMoi, err := strconv.Atoi(lastNumStr); err == nil {
					if firstCodeDM != "" { core.CapNhatSlotThuCong(adminShopID, firstCodeDM, slotMoi) }
				}
			}
		}
	}

	existingMap := make(map[string]*core.SanPhamMayTinh)
	for _, sp := range existingSKUs { existingMap[sp.LayIDDuyNhat()] = sp }
	processedSKUs := make(map[string]bool) 

	core.KhoaHeThong.Lock()
	defer core.KhoaHeThong.Unlock()

	for i, in := range inputSKUs {
		skuID := in.MaSKU
		if skuID == "" { skuID = fmt.Sprintf("%s-%02d", maSP, i+1) }
		
		var sp *core.SanPhamMayTinh
		isNewSKU := false
		
		if exist, ok := existingMap[skuID]; ok {
			sp = exist; processedSKUs[skuID] = true
		} else {
			isNewSKU = true
			currentList := core.CacheSanPhamMayTinh[adminShopID] 
			sp = &core.SanPhamMayTinh{
				SpreadsheetID:  adminShopID, 
				DongTrongSheet: core.DongBatDau_SanPhamMayTinh + len(currentList),
				MaSanPham:      maSP,
				MaSKU:          skuID,
			}
		}

		newTenSanPham   := strings.TrimSpace(in.TenSanPham)
		newTenRutGon    := strings.TrimSpace(in.TenRutGon)
		newSlug         := Repo_TaoSlugChuan(newTenSanPham)
		newTenSKU       := strings.TrimSpace(in.TenSKU)
		newMaDanhMuc    := Repo_XuLyTags(in.MaDanhMuc)
		newMaThuongHieu := Repo_XuLyTags(in.MaThuongHieu)
		newDonVi        := Repo_XuLyTags(in.DonVi)
		newMauSac       := Repo_XuLyTags(in.MauSac)
		newUrlHinhAnh   := strings.TrimSpace(in.UrlHinhAnh)
		newTinhTrang    := Repo_XuLyTags(in.TinhTrang)

		isChanged := false

		if isNewSKU {
			isChanged = true
			sp.NgayTao = nowStr; sp.NguoiTao = userID
			sp.NgayCapNhat = nowStr; sp.NguoiCapNhat = userID
		} else {
			if sp.TenSanPham != newTenSanPham || sp.TenRutGon != newTenRutGon || sp.Slug != newSlug || sp.TenSKU != newTenSKU ||
				sp.SKUChinh != in.SKUChinh || sp.TrangThai != in.TrangThai || sp.MaDanhMuc != newMaDanhMuc || sp.MaThuongHieu != newMaThuongHieu ||
				sp.DonVi != newDonVi || sp.MauSac != newMauSac || sp.KhoiLuong != in.KhoiLuong || sp.KichThuoc != in.KichThuoc ||
				sp.UrlHinhAnh != newUrlHinhAnh || sp.ThongSoHTML != in.ThongSoHTML || sp.MoTaHTML != in.MoTaHTML || sp.BaoHanh != in.BaoHanh ||
				sp.TinhTrang != newTinhTrang || sp.GiaNhap != in.GiaNhap || sp.PhanTramLai != in.PhanTramLai || sp.GiaNiemYet != in.GiaNiemYet ||
				sp.PhanTramGiam != in.PhanTramGiam || sp.SoTienGiam != in.SoTienGiam || sp.GiaBan != in.GiaBan || sp.GhiChu != in.GhiChu {
				
				isChanged = true
				sp.NgayCapNhat = nowStr; sp.NguoiCapNhat = userID
			}
		}

		if isChanged {
			sp.TenSanPham = newTenSanPham; sp.TenRutGon = newTenRutGon; sp.Slug = newSlug; sp.TenSKU = newTenSKU
			sp.SKUChinh = in.SKUChinh; sp.TrangThai = in.TrangThai; sp.MaDanhMuc = newMaDanhMuc; sp.MaThuongHieu = newMaThuongHieu
			sp.DonVi = newDonVi; sp.MauSac = newMauSac; sp.KhoiLuong = in.KhoiLuong; sp.KichThuoc = in.KichThuoc
			sp.UrlHinhAnh = newUrlHinhAnh; sp.ThongSoHTML = in.ThongSoHTML; sp.MoTaHTML = in.MoTaHTML; sp.BaoHanh = in.BaoHanh
			sp.TinhTrang = newTinhTrang; sp.GiaNhap = in.GiaNhap; sp.PhanTramLai = in.PhanTramLai; sp.GiaNiemYet = in.GiaNiemYet
			sp.PhanTramGiam = in.PhanTramGiam; sp.SoTienGiam = in.SoTienGiam; sp.GiaBan = in.GiaBan; sp.GhiChu = in.GhiChu

			if isNewSKU {
				core.CacheSanPhamMayTinh[adminShopID] = append(core.CacheSanPhamMayTinh[adminShopID], sp)
				core.CacheMapSKUMayTinh[core.TaoCompositeKey(adminShopID, sp.LayIDDuyNhat())] = sp
				core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(adminShopID, sp.MaSanPham)] = append(core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(adminShopID, sp.MaSanPham)], sp)
			}

			ghi := core.ThemVaoHangCho
			sheet := core.TenSheetMayTinh
			r := sp.DongTrongSheet
			
			ghi(adminShopID, sheet, r, core.CotPC_MaSanPham, sp.MaSanPham); ghi(adminShopID, sheet, r, core.CotPC_TenSanPham, sp.TenSanPham)
			ghi(adminShopID, sheet, r, core.CotPC_TenRutGon, sp.TenRutGon); ghi(adminShopID, sheet, r, core.CotPC_Slug, sp.Slug)
			ghi(adminShopID, sheet, r, core.CotPC_MaSKU, sp.MaSKU); ghi(adminShopID, sheet, r, core.CotPC_TenSKU, sp.TenSKU)
			ghi(adminShopID, sheet, r, core.CotPC_SKUChinh, sp.SKUChinh); ghi(adminShopID, sheet, r, core.CotPC_TrangThai, sp.TrangThai)
			ghi(adminShopID, sheet, r, core.CotPC_MaDanhMuc, sp.MaDanhMuc); ghi(adminShopID, sheet, r, core.CotPC_MaThuongHieu, sp.MaThuongHieu)
			ghi(adminShopID, sheet, r, core.CotPC_DonVi, sp.DonVi); ghi(adminShopID, sheet, r, core.CotPC_MauSac, sp.MauSac)
			ghi(adminShopID, sheet, r, core.CotPC_KhoiLuong, sp.KhoiLuong); ghi(adminShopID, sheet, r, core.CotPC_KichThuoc, sp.KichThuoc)
			ghi(adminShopID, sheet, r, core.CotPC_UrlHinhAnh, sp.UrlHinhAnh); ghi(adminShopID, sheet, r, core.CotPC_ThongSoHTML, sp.ThongSoHTML)
			ghi(adminShopID, sheet, r, core.CotPC_MoTaHTML, sp.MoTaHTML); ghi(adminShopID, sheet, r, core.CotPC_BaoHanh, sp.BaoHanh)
			ghi(adminShopID, sheet, r, core.CotPC_TinhTrang, sp.TinhTrang); ghi(adminShopID, sheet, r, core.CotPC_GiaNhap, sp.GiaNhap)
			ghi(adminShopID, sheet, r, core.CotPC_PhanTramLai, sp.PhanTramLai); ghi(adminShopID, sheet, r, core.CotPC_GiaNiemYet, sp.GiaNiemYet)
			ghi(adminShopID, sheet, r, core.CotPC_PhanTramGiam, sp.PhanTramGiam); ghi(adminShopID, sheet, r, core.CotPC_SoTienGiam, sp.SoTienGiam)
			ghi(adminShopID, sheet, r, core.CotPC_GiaBan, sp.GiaBan); ghi(adminShopID, sheet, r, core.CotPC_GhiChu, sp.GhiChu)
			
			if isNewSKU {
				ghi(adminShopID, sheet, r, core.CotPC_NguoiTao, sp.NguoiTao); ghi(adminShopID, sheet, r, core.CotPC_NgayTao, sp.NgayTao)
			}
			ghi(adminShopID, sheet, r, core.CotPC_NguoiCapNhat, sp.NguoiCapNhat); ghi(adminShopID, sheet, r, core.CotPC_NgayCapNhat, sp.NgayCapNhat)
		}
	}

	for skuID, sp := range existingMap {
		if !processedSKUs[skuID] {
			if sp.TrangThai != -1 {
				sp.TrangThai = -1; sp.SKUChinh = 0; sp.NgayCapNhat = nowStr; sp.NguoiCapNhat = userID
				core.ThemVaoHangCho(adminShopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_TrangThai, -1)
				core.ThemVaoHangCho(adminShopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_SKUChinh, 0)
				core.ThemVaoHangCho(adminShopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_NgayCapNhat, nowStr)
				core.ThemVaoHangCho(adminShopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_NguoiCapNhat, userID)
			}
		}
	}
	return nil
}
