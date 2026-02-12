package chuc_nang

import (
	"net/http"
	"strconv"
	"strings"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangQuanLyCaiDat(c *gin.Context) {
	userID := c.GetString("USER_ID")
	kh, _ := core.LayKhachHang(userID)

	c.HTML(http.StatusOK, "quan_tri_cai_dat", gin.H{
		"TieuDe":         "Cài đặt hệ thống",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"ListDanhMuc":    core.LayDanhSachDanhMuc(),
		"ListThuongHieu": core.LayDanhSachThuongHieu(),
		"ListBLN":        core.LayDanhSachBienLoiNhuan(), // [MỚI] Truyền dữ liệu Khung giá ra View
	})
}

// API_LuuDanhMuc
func API_LuuDanhMuc(c *gin.Context) {
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

	sheetID := cau_hinh.BienCauHinh.IdFileSheet
	var targetRow int

	if isNew {
		if _, ok := core.LayChiTietDanhMuc(maDM); ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã danh mục này đã tồn tại!"})
			return
		}
		targetRow = core.DongBatDau_DanhMuc + len(core.LayDanhSachDanhMuc())
		newDM := &core.DanhMuc{
			SpreadsheetID:  sheetID,
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
		found, ok := core.LayChiTietDanhMuc(maDM)
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
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_MaDanhMuc, strings.ToUpper(maDM))
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_TenDanhMuc, tenDM)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_DanhMucMe, dmMe)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_ThueVAT, thueVAT)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_LoiNhuan, loiNhuan)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_TrangThai, trangThai)
	if isNew { ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_Slot, 0) }

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Danh mục thành công!"})
}

// API_LuuThuongHieu
func API_LuuThuongHieu(c *gin.Context) {
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

	sheetID := cau_hinh.BienCauHinh.IdFileSheet
	var targetRow int

	if isNew {
		targetRow = core.DongBatDau_ThuongHieu + len(core.LayDanhSachThuongHieu())
		newTH := &core.ThuongHieu{
			SpreadsheetID:  sheetID,
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
		for _, item := range core.LayDanhSachThuongHieu() {
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
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_MaThuongHieu, strings.ToUpper(maTH))
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_TenThuongHieu, tenTH)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_LogoUrl, logoUrl)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_MoTa, moTa)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_TrangThai, trangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Thương hiệu thành công!"})
}

// [MỚI] API_LuuBienLoiNhuan
func API_LuuBienLoiNhuan(c *gin.Context) {
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

	sheetID := cau_hinh.BienCauHinh.IdFileSheet
	var targetRow int

	if isNew {
		targetRow = core.DongBatDau_BienLoiNhuan + len(core.LayDanhSachBienLoiNhuan())
		newBLN := &core.BienLoiNhuan{
			SpreadsheetID:  sheetID,
			DongTrongSheet: targetRow,
			KhungGiaNhap:   khungGia,
			BienLoiNhuan:   loiNhuan,
			TrangThai:      trangThai,
		}
		core.ThemBienLoiNhuanVaoRam(newBLN) 
	} else {
		targetRow = dongCu
		core.SuaBienLoiNhuanTrongRam(dongCu, khungGia, loiNhuan, trangThai) // [ĐÃ SỬA] Update RAM an toàn
	}

	ghi := core.ThemVaoHangCho
	ghi(sheetID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_KhungGiaNhap, khungGia)
	ghi(sheetID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_BienLoiNhuan, loiNhuan)
	ghi(sheetID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_TrangThai, trangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu cấu hình Khung giá thành công!"})
}

// [MỚI] API Đồng bộ lại Slot (Bộ đếm) cho Danh mục dựa trên dữ liệu Sản phẩm thực tế
func API_DongBoSlotDanhMuc(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	core.KhoaHeThong.Lock()
	defer core.KhoaHeThong.Unlock()

	// 1. Quét toàn bộ sản phẩm để tìm số Max của từng Prefix
	mapMaxSlot := make(map[string]int) // Key: Prefix (VD: MAIN), Value: Max Number (VD: 45)
	
	listSP := core.LayDanhSachSanPham()
	for _, sp := range listSP {
		maSP := strings.TrimSpace(sp.MaSanPham)
		// Logic tách số: Chạy từ cuối chuỗi về đầu, lấy hết các ký tự số
		ketThucSo := len(maSP)
		batDauSo := ketThucSo
		
		for i := len(maSP) - 1; i >= 0; i-- {
			char := maSP[i]
			if char >= '0' && char <= '9' {
				batDauSo = i
			} else {
				break // Gặp chữ cái thì dừng
			}
		}

		// Nếu tìm thấy đoạn số ở cuối
		if batDauSo < ketThucSo {
			prefix := strings.ToUpper(maSP[0:batDauSo]) // Phần chữ (VD: MAIN)
			soStr := maSP[batDauSo:ketThucSo]           // Phần số (VD: 0045)
			
			so, err := strconv.Atoi(soStr)
			if err == nil {
				if so > mapMaxSlot[prefix] {
					mapMaxSlot[prefix] = so
				}
			}
		}
	}

	// 2. Cập nhật lại Slot cho Danh Mục
	countUpdate := 0
	for _, dm := range core.LayDanhSachDanhMuc() {
		maxThucTe, coDuLieu := mapMaxSlot[dm.MaDanhMuc]
		
		// Chỉ cập nhật nếu số thực tế lớn hơn Slot hiện tại
		// Hoặc nếu Slot hiện tại bị mất (về 0) mà thực tế lại có hàng
		if coDuLieu && maxThucTe > dm.Slot {
			dm.Slot = maxThucTe
			// Ghi đè lại vào Sheet
			core.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "DANH_MUC", dm.DongTrongSheet, core.CotDM_Slot, dm.Slot)
			countUpdate++
		}
	}

	msg := "Hệ thống đã đồng bộ xong. Không có lệch lạc nào."
	if countUpdate > 0 {
		msg = "Đã phát hiện và đồng bộ lại Slot cho " + strconv.Itoa(countUpdate) + " danh mục."
	}

	c.JSON(200, gin.H{"status": "ok", "msg": msg})
}
