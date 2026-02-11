// API_LuuDanhMuc
func API_LuuDanhMuc(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maDM := strings.TrimSpace(c.PostForm("ma_danh_muc"))
	tenDM := strings.TrimSpace(c.PostForm("ten_danh_muc"))
	thueVAT, _ := strconv.ParseFloat(c.PostForm("thue_vat"), 64)
	loiNhuan, _ := strconv.ParseFloat(c.PostForm("loi_nhuan"), 64)
	isNew := c.PostForm("is_new") == "true"

	if maDM == "" || tenDM == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên danh mục không được để trống!"})
		return
	}

	sheetID := cau_hinh.BienCauHinh.IdFileSheet
	var targetRow int

	// TÁCH BẠCH LUỒNG LƯU RAM ĐỂ TRÁNH TREO LOCK
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
			ThueVAT:        thueVAT,
			LoiNhuan:       loiNhuan,
			STT:            0,
		}
		core.ThemDanhMucVaoRam(newDM) // Đã an toàn 100%
	} else {
		found, ok := core.LayChiTietDanhMuc(maDM)
		if !ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy danh mục để sửa!"})
			return
		}
		
		targetRow = found.DongTrongSheet
		core.KhoaHeThong.Lock()
		found.TenDanhMuc = tenDM
		found.ThueVAT = thueVAT
		found.LoiNhuan = loiNhuan
		core.KhoaHeThong.Unlock()
	}

	// Ghi Hàng chờ (Queue)
	ghi := core.ThemVaoHangCho
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_MaDanhMuc, strings.ToUpper(maDM))
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_TenDanhMuc, tenDM)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_ThueVAT, thueVAT)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_LoiNhuan, loiNhuan)
	if isNew {
		ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_STT, 0)
	}

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
	logo := strings.TrimSpace(c.PostForm("logo"))
	moTa := strings.TrimSpace(c.PostForm("mo_ta"))
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
			Logo:           logo,
			MoTa:           moTa,
			TrangThai:      1,
		}
		core.ThemThuongHieuVaoRam(newTH) // An toàn
	} else {
		var found *core.ThuongHieu
		for _, item := range core.LayDanhSachThuongHieu() {
			if item.MaThuongHieu == maTH {
				found = item
				break
			}
		}
		if found == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thương hiệu để sửa!"})
			return
		}
		
		targetRow = found.DongTrongSheet
		core.KhoaHeThong.Lock()
		found.TenThuongHieu = tenTH
		found.Logo = logo
		found.MoTa = moTa
		core.KhoaHeThong.Unlock()
	}

	// Ghi Hàng chờ (Queue)
	ghi := core.ThemVaoHangCho
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_MaThuongHieu, strings.ToUpper(maTH))
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_TenThuongHieu, tenTH)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_Logo, logo)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_MoTa, moTa)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_TrangThai, 1)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Thương hiệu thành công!"})
}
