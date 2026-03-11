package product_master

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/core"
)

func Service_LayChiTietSanPham(shopID string, maSP string) (*core.ProductJSON, error) {
	core.KhoaHeThong.RLock()
	sp, ok := core.CacheMapSanPham[core.TaoCompositeKey(shopID, maSP)]
	core.KhoaHeThong.RUnlock()

	if !ok || sp == nil {
		return nil, errors.New("Không tìm thấy sản phẩm!")
	}
	return sp, nil
}

func Service_LuuSanPham(masterShopID, adminShopID, vaiTro, userID, maNganh, maSP, dataJSON string) error {
	// 1. Kiểm tra Quyền
	if maSP == "" {
		if vaiTro != "quan_tri_he_thong" && !core.KiemTraQuyen(masterShopID, vaiTro, "product.create") {
			return errors.New("Bạn không có quyền thêm sản phẩm!")
		}
	} else {
		if vaiTro != "quan_tri_he_thong" && !core.KiemTraQuyen(masterShopID, vaiTro, "product.edit") {
			return errors.New("Bạn không có quyền sửa sản phẩm!")
		}
	}

	// 2. Tra cứu định tuyến Sheet vật lý từ RAM Cấu Hình
	core.KhoaHeThong.RLock()
	cfgNganh, ok := core.CacheMapNganh[maNganh]
	core.KhoaHeThong.RUnlock()

	if !ok || cfgNganh.TenSheet == "" {
		return errors.New("Ngành hàng không hợp lệ hoặc chưa khai báo Sheet trong Cấu hình!")
	}

	// 3. Phân rã JSON để đóng dấu
	var inputSP core.ProductJSON
	if err := json.Unmarshal([]byte(dataJSON), &inputSP); err != nil {
		return errors.New("Dữ liệu JSON truyền lên không hợp lệ!")
	}

	now := time.Now().Unix()
	nowStr := time.Unix(now, 0).Format("2006-01-02 15:04:05")

	core.KhoaHeThong.Lock()
	defer core.KhoaHeThong.Unlock()

	spCu, isUpdate := core.CacheMapSanPham[core.TaoCompositeKey(adminShopID, maSP)]

	// 4. Nhào nặn dữ liệu chuẩn hóa (Auto-fill Metadata)
	if !isUpdate || maSP == "" {
		// TẠO MỚI (Cấp phát Mã SP mới)
		if maSP == "" {
			maSP = fmt.Sprintf("SP%d", now) // TODO: Hàm sinh mã quy chuẩn
		}
		inputSP.MaSanPham = maSP
		inputSP.CreatedAt = now
		inputSP.UpdatedAt = now
		inputSP.Version = 1
		inputSP.QuanLy.NguoiTao = userID
		inputSP.QuanLy.NgayTao = nowStr
	} else {
		// CẬP NHẬT
		inputSP.MaSanPham = maSP
		inputSP.CreatedAt = spCu.CreatedAt
		inputSP.Version = spCu.Version + 1 // Tăng version migrate
		inputSP.UpdatedAt = now
		inputSP.QuanLy.NguoiTao = spCu.QuanLy.NguoiTao
		inputSP.QuanLy.NgayTao = spCu.QuanLy.NgayTao
	}
	
	inputSP.QuanLy.NguoiCapNhat = userID
	inputSP.QuanLy.NgayCapNhat = nowStr
	inputSP.Slug = Repo_TaoSlugChuan(inputSP.TenSanPham)
	inputSP.SearchText = Repo_BuildSearchText(&inputSP)

	// Xử lý cấp Mã SKU tự động nếu trống
	for i := range inputSP.SKU {
		if inputSP.SKU[i].MaSKU == "" {
			inputSP.SKU[i].MaSKU = fmt.Sprintf("%s-%02d", maSP, i+1)
		}
	}

	// Đóng gói lại thành chuỗi JSON cuối cùng để lưu
	finalJSONBytes, _ := json.Marshal(inputSP)
	finalJSONString := string(finalJSONBytes)
	spPtr := &inputSP

	// 5. Thao tác RAM Siêu Tốc (Ghi đè O(1))
	if isUpdate {
		// Tìm vị trí trong Array để chèn đè
		listSP := core.CacheSanPham[adminShopID][maNganh]
		for i, v := range listSP {
			if v.MaSanPham == maSP {
				core.CacheSanPham[adminShopID][maNganh][i] = spPtr
				break
			}
		}
		// Dọn rác SKU cũ
		for _, oldSKU := range spCu.SKU {
			delete(core.CacheMapSKU, core.TaoCompositeKey(adminShopID, oldSKU.MaSKU))
		}
	} else {
		// Append mới tinh
		core.CacheSanPham[adminShopID][maNganh] = append(core.CacheSanPham[adminShopID][maNganh], spPtr)
	}

	// Cập nhật Map O(1)
	core.CacheMapSanPham[core.TaoCompositeKey(adminShopID, maSP)] = spPtr
	for i := range inputSP.SKU {
		core.CacheMapSKU[core.TaoCompositeKey(adminShopID, inputSP.SKU[i].MaSKU)] = &inputSP.SKU[i]
	}

	// 6. Gửi Tín hiệu "Note" vào Hàng chờ thông minh (Không gửi JSON)
	if isUpdate {
		core.GhiChuDongBo(adminShopID, cfgNganh.TenSheet, "UPDATE", spPtr.MaSanPham)
	} else {
		// Tính số dòng dự kiến cho SP Mới
		spPtr.DongTrongSheet = core.DongBatDau_Product + len(core.CacheSanPham[adminShopID][maNganh]) - 1
		core.GhiChuDongBo(adminShopID, cfgNganh.TenSheet, "APPEND", spPtr.MaSanPham)
	}

	return nil
}

	core.ThemVaoHangCho(adminShopID, cfgNganh.TenSheet, dongSheet, core.CotProd_MaSanPham, maSP)
	core.ThemVaoHangCho(adminShopID, cfgNganh.TenSheet, dongSheet, core.CotProd_DataJSON, finalJSONString)

	return nil
}
