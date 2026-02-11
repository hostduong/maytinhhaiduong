package chuc_nang

import (
	"net/http"
	"strconv"
	"strings"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// Hiển thị trang
func TrangQuanLyDanhMuc(c *gin.Context) {
	userID := c.GetString("USER_ID")
	kh, _ := core.LayKhachHang(userID)

	c.HTML(http.StatusOK, "quan_tri_danh_muc", gin.H{
		"TieuDe":         "Cấu hình Danh mục & Thương hiệu",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"ListDanhMuc":    core.LayDanhSachDanhMuc(),
		"ListThuongHieu": core.LayDanhSachThuongHieu(),
	})
}

// API_LuuDanhMuc: Thêm hoặc Sửa Danh Mục
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

	var dm *core.DanhMuc
	sheetID := cau_hinh.BienCauHinh.IdFileSheet

	if isNew {
		if _, ok := core.LayChiTietDanhMuc(maDM); ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã danh mục này đã tồn tại!"})
			return
		}
		dm = &core.DanhMuc{
			SpreadsheetID:  sheetID,
			DongTrongSheet: core.DongBatDau_DanhMuc + len(core.LayDanhSachDanhMuc()),
			MaDanhMuc:      strings.ToUpper(maDM), // Mã luôn in hoa (VD: MAIN)
			STT:            0,
		}
	} else {
		found, ok := core.LayChiTietDanhMuc(maDM)
		if !ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy danh mục để sửa!"})
			return
		}
		dm = found
	}

	// Lưu vào RAM
	core.KhoaHeThong.Lock()
	dm.TenDanhMuc = tenDM
	dm.ThueVAT = thueVAT
	dm.LoiNhuan = loiNhuan
	core.KhoaHeThong.Unlock()

	// Ghi xuống Sheet
	ghi := core.ThemVaoHangCho
	ghi(sheetID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_MaDanhMuc, dm.MaDanhMuc)
	ghi(sheetID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_TenDanhMuc, dm.TenDanhMuc)
	ghi(sheetID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_ThueVAT, dm.ThueVAT)
	ghi(sheetID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_LoiNhuan, dm.LoiNhuan)
	if isNew {
		ghi(sheetID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_STT, dm.STT)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu Danh mục thành công! (Nên bấm Đồng Bộ để cập nhật RAM)"})
}

// API_LuuThuongHieu: Thêm hoặc Sửa Thương Hiệu
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
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên thương hiệu không được để trống!"})
		return
	}

	var th *core.ThuongHieu
	sheetID := cau_hinh.BienCauHinh.IdFileSheet

	if isNew {
		th = &core.ThuongHieu{
			SpreadsheetID:  sheetID,
			DongTrongSheet: core.DongBatDau_ThuongHieu + len(core.LayDanhSachThuongHieu()),
			MaThuongHieu:   strings.ToUpper(maTH),
		}
	} else {
		// Tìm thương hiệu trong list
		for _, item := range core.LayDanhSachThuongHieu() {
			if item.MaThuongHieu == maTH {
				th = item
				break
			}
		}
		if th == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thương hiệu để sửa!"})
			return
		}
	}

	// Lưu vào RAM
	core.KhoaHeThong.Lock()
	th.TenThuongHieu = tenTH
	th.Logo = logo
	th.MoTa = moTa
	th.TrangThai = 1
	core.KhoaHeThong.Unlock()

	// Ghi xuống Sheet
	ghi := core.ThemVaoHangCho
	ghi(sheetID, "THUONG_HIEU", th.DongTrongSheet, core.CotTH_MaThuongHieu, th.MaThuongHieu)
	ghi(sheetID, "THUONG_HIEU", th.DongTrongSheet, core.CotTH_TenThuongHieu, th.TenThuongHieu)
	ghi(sheetID, "THUONG_HIEU", th.DongTrongSheet, core.CotTH_Logo, th.Logo)
	ghi(sheetID, "THUONG_HIEU", th.DongTrongSheet, core.CotTH_MoTa, th.MoTa)
	ghi(sheetID, "THUONG_HIEU", th.DongTrongSheet, core.CotTH_TrangThai, th.TrangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu Thương hiệu thành công! (Nên bấm Đồng Bộ để cập nhật RAM)"})
}
