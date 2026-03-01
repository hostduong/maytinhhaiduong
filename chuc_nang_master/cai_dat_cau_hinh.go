package chuc_nang_master

import (
	"net/http"
	"strconv"
	"strings"

	"app/core"
	"github.com/gin-gonic/gin"
)

// =============================================================
// 1. TRANG QUẢN TRỊ CÀI ĐẶT CẤU HÌNH MASTER
// =============================================================
func TrangCaiDatCauHinhMaster(c *gin.Context) {
	// [SAAS] Lấy thông tin Shop & User
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	kh, _ := core.LayKhachHang(shopID, userID)

	c.HTML(http.StatusOK, "master_cai_dat_cau_hinh", gin.H{
		"TieuDe":         "Cài đặt hệ thống",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		
		// [SAAS] Load dữ liệu theo Shop
		"ListDanhMuc":    core.LayDanhSachDanhMuc(shopID),
		"ListThuongHieu": core.LayDanhSachThuongHieu(shopID),
		"ListBLN":        core.LayDanhSachBienLoiNhuan(shopID),
		"ListNCC":        core.LayDanhSachNhaCungCap(shopID),
	})
}

// =============================================================
// 2. API DANH MỤC
// =============================================================
func API_LuuDanhMucMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maDM := strings.TrimSpace(c.PostForm("ma_danh_muc"))
	tenDM := strings.TrimSpace(c.PostForm("ten_danh_muc"))
	dmMe := strings.TrimSpace(c.PostForm("danh_muc_me"))
	thueVAT, _ := strconv.ParseFloat(c.PostForm("thue_vat"), 64)
	loiNhuan, _ := strconv.ParseFloat(c.PostForm("loi_nhuan"), 64)
	trangThai := 0; if c.PostForm("trang_thai") == "on" { trangThai = 1 }
	isNew := c.PostForm("is_new") == "true"

	if maDM == "" || tenDM == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên danh mục không được để trống!"})
		return
	}
	if strings.ToUpper(maDM) == strings.ToUpper(dmMe) {
		c.JSON(200, gin.H{"status": "error", "msg": "Danh mục mẹ không được trùng với chính nó!"})
		return
	}

	var targetRow int

	if isNew {
		if _, ok := core.LayChiTietDanhMuc(shopID, maDM); ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã danh mục này đã tồn tại!"})
			return
		}
		
		targetRow = core.DongBatDau_DanhMuc + len(core.LayDanhSachDanhMuc(shopID))
		
		newDM := &core.DanhMuc{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			MaDanhMuc:      strings.ToUpper(maDM),
			TenDanhMuc:     tenDM,
			DanhMucMe:      dmMe,
			ThueVAT:        thueVAT,
			LoiNhuan:       loiNhuan,
			Slot:           0,
			TrangThai:      trangThai,
		}
		core.ThemDanhMucVaoRam(newDM) 
	} else {
		found, ok := core.LayChiTietDanhMuc(shopID, maDM)
		if !ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy danh mục để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet
		
		func() {
			core.KhoaHeThong.Lock()
			defer core.KhoaHeThong.Unlock()
			found.TenDanhMuc = tenDM
			found.DanhMucMe = dmMe
			found.ThueVAT = thueVAT
			found.LoiNhuan = loiNhuan
			found.TrangThai = trangThai
		}()
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_MaDanhMuc, strings.ToUpper(maDM))
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_TenDanhMuc, tenDM)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_DanhMucMe, dmMe)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_ThueVAT, thueVAT)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_LoiNhuan, loiNhuan)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_TrangThai, trangThai)
	
	if isNew { ghi(shopID, "DANH_MUC", targetRow, core.CotDM_Slot, 0) }

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Danh mục thành công!"})
}

func API_DongBoSlotDanhMucMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	listSP := core.LayDanhSachSanPhamMayTinh(shopID)
	mapMaxSlot := make(map[string]int) 
	
	for _, sp := range listSP {
		maSP := strings.TrimSpace(sp.MaSanPham)
		ketThucSo := len(maSP)
		batDauSo := ketThucSo
		
		for i := len(maSP) - 1; i >= 0; i-- {
			char := maSP[i]
			if char >= '0' && char <= '9' {
				batDauSo = i
			} else {
				break
			}
		}

		if batDauSo < ketThucSo {
			prefix := strings.ToUpper(maSP[0:batDauSo])
			soStr := maSP[batDauSo:ketThucSo]
			so, err := strconv.Atoi(soStr)
			if err == nil {
				if so > mapMaxSlot[prefix] {
					mapMaxSlot[prefix] = so
				}
			}
		}
	}

	listDM := core.LayDanhSachDanhMuc(shopID)

	core.KhoaHeThong.Lock()
	countUpdate := 0
	
	for _, dm := range listDM {
		maxThucTe, coDuLieu := mapMaxSlot[dm.MaDanhMuc]
		if coDuLieu && maxThucTe > dm.Slot {
			dm.Slot = maxThucTe
			core.ThemVaoHangCho(shopID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_Slot, dm.Slot)
			countUpdate++
		}
	}
	core.KhoaHeThong.Unlock()

	msg := "Đã đồng bộ xong. Các bộ đếm đều chính xác."
	if countUpdate > 0 {
		msg = "Đã cập nhật lại Slot cho " + strconv.Itoa(countUpdate) + " danh mục."
	}

	c.JSON(200, gin.H{"status": "ok", "msg": msg})
}

// =============================================================
// 3. API THƯƠNG HIỆU
// =============================================================
func API_LuuThuongHieuMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maTH := strings.TrimSpace(c.PostForm("ma_thuong_hieu"))
	tenTH := strings.TrimSpace(c.PostForm("ten_thuong_hieu"))
	logoUrl := strings.TrimSpace(c.PostForm("logo_url"))
	moTa := strings.TrimSpace(c.PostForm("mo_ta"))
	trangThai := 0; if c.PostForm("trang_thai") == "on" { trangThai = 1 }
	isNew := c.PostForm("is_new") == "true"

	if maTH == "" || tenTH == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên thương hiệu không được trống!"})
		return
	}

	var targetRow int

	if isNew {
		targetRow = core.DongBatDau_ThuongHieu + len(core.LayDanhSachThuongHieu(shopID))
		newTH := &core.ThuongHieu{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			MaThuongHieu:   strings.ToUpper(maTH),
			TenThuongHieu:  tenTH,
			LogoUrl:        logoUrl,
			MoTa:           moTa,
			TrangThai:      trangThai,
		}
		core.ThemThuongHieuVaoRam(newTH) 
	} else {
		var found *core.ThuongHieu
		for _, item := range core.LayDanhSachThuongHieu(shopID) {
			if item.MaThuongHieu == maTH { found = item; break }
		}
		if found == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thương hiệu để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet
		func() {
			core.KhoaHeThong.Lock()
			defer core.KhoaHeThong.Unlock()
			found.TenThuongHieu = tenTH
			found.LogoUrl = logoUrl
			found.MoTa = moTa
			found.TrangThai = trangThai
		}()
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_MaThuongHieu, strings.ToUpper(maTH))
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_TenThuongHieu, tenTH)
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_LogoUrl, logoUrl)
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_MoTa, moTa)
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_TrangThai, trangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Thương hiệu thành công!"})
}

// =============================================================
// 4. API BIÊN LỢI NHUẬN
// =============================================================
func API_LuuBienLoiNhuanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	khungGia, err1 := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("khung_gia_nhap"), ".", ""), 64)
	loiNhuan, err2 := strconv.ParseFloat(c.PostForm("bien_loi_nhuan"), 64)
	trangThai := 0; if c.PostForm("trang_thai") == "on" { trangThai = 1 }
	isNew := c.PostForm("is_new") == "true"
	dongCu, _ := strconv.Atoi(c.PostForm("dong_cu")) 

	if err1 != nil || err2 != nil || khungGia <= 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Khung giá và Biên lợi nhuận phải là số hợp lệ > 0!"})
		return
	}

	var targetRow int

	if isNew {
		targetRow = core.DongBatDau_BienLoiNhuan + len(core.LayDanhSachBienLoiNhuan(shopID))
		newBLN := &core.BienLoiNhuan{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			KhungGiaNhap:   khungGia,
			BienLoiNhuan:   loiNhuan,
			TrangThai:      trangThai,
		}
		core.ThemBienLoiNhuanVaoRam(newBLN) 
	} else {
		targetRow = dongCu
		core.SuaBienLoiNhuanTrongRam(shopID, dongCu, khungGia, loiNhuan, trangThai) 
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_KhungGiaNhap, khungGia)
	ghi(shopID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_BienLoiNhuan, loiNhuan)
	ghi(shopID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_TrangThai, trangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu cấu hình Khung giá thành công!"})
}

// =============================================================
// 5. API NHÀ CUNG CẤP
// =============================================================
func API_LuuNhaCungCapMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maNCC := strings.TrimSpace(c.PostForm("ma_nha_cung_cap"))
	tenNCC := strings.TrimSpace(c.PostForm("ten_nha_cung_cap"))
	sdt := strings.TrimSpace(c.PostForm("dien_thoai"))
	email := strings.TrimSpace(c.PostForm("email"))
	diaChi := strings.TrimSpace(c.PostForm("dia_chi"))
	mst := strings.TrimSpace(c.PostForm("ma_so_thue"))
	nguoiLh := strings.TrimSpace(c.PostForm("nguoi_lien_he"))
	nganHang := strings.TrimSpace(c.PostForm("ngan_hang"))
	ghiChu := strings.TrimSpace(c.PostForm("ghi_chu"))
	trangThai := 0
	if c.PostForm("trang_thai") == "on" {
		trangThai = 1
	}
	isNew := c.PostForm("is_new") == "true"

	if tenNCC == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên nhà cung cấp không được để trống!"})
		return
	}

	var targetRow int

	if isNew {
		maMoi := maNCC
		if maMoi == "" { maMoi = core.TaoMaNhaCungCapMoi(shopID) }
		targetRow = core.DongBatDau_NhaCungCap + len(core.LayDanhSachNhaCungCap(shopID))

		newNCC := &core.NhaCungCap{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			MaNhaCungCap:   maMoi,
			TenNhaCungCap:  tenNCC,
			DienThoai:      sdt,
			Email:          email,
			DiaChi:         diaChi,
			MaSoThue:       mst,
			NguoiLienHe:    nguoiLh,
			NganHang:       nganHang,
			TrangThai:      trangThai,
			GhiChu:         ghiChu,
		}
		
		core.KhoaHeThong.Lock()
		core.CacheNhaCungCap[shopID] = append(core.CacheNhaCungCap[shopID], newNCC)
		core.KhoaHeThong.Unlock()
		
		maNCC = maMoi
	} else {
		list := core.LayDanhSachNhaCungCap(shopID)
		var found *core.NhaCungCap
		for _, item := range list {
			if item.MaNhaCungCap == maNCC { found = item; break }
		}

		if found == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy Nhà cung cấp để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet

		core.KhoaHeThong.Lock()
		found.TenNhaCungCap = tenNCC
		found.DienThoai = sdt
		found.Email = email
		found.DiaChi = diaChi
		found.MaSoThue = mst
		found.NguoiLienHe = nguoiLh
		found.NganHang = nganHang
		found.TrangThai = trangThai
		found.GhiChu = ghiChu
		core.KhoaHeThong.Unlock()
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_MaNhaCungCap, maNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TenNhaCungCap, tenNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DienThoai, sdt)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_Email, email)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DiaChi, diaChi)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_MaSoThue, mst)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NguoiLienHe, nguoiLh)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NganHang, nganHang)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TrangThai, trangThai)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_GhiChu, ghiChu)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Nhà cung cấp thành công!"})
}
