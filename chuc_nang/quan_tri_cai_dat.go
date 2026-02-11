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
	dmMe := strings.TrimSpace(c.PostForm("danh_muc_me")) // [MỚI] Lấy từ form
	thueVAT, _ := strconv.ParseFloat(c.PostForm("thue_vat"), 64)
	loiNhuan, _ := strconv.ParseFloat(c.PostForm("loi_nhuan"), 64)
	isNew := c.PostForm("is_new") == "true"

	if maDM == "" || tenDM == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên danh mục không được để trống!"})
		return
	}

	// Chặn lỗi: Danh mục mẹ không được là chính nó
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
			DanhMucMe:      dmMe, // [MỚI]
			ThueVAT:        thueVAT,
			LoiNhuan:       loiNhuan,
			STT:            0,
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
			found.DanhMucMe = dmMe // [MỚI]
			found.ThueVAT = thueVAT
			found.LoiNhuan = loiNhuan
		}()
	}

	ghi := core.ThemVaoHangCho
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_MaDanhMuc, strings.ToUpper(maDM))
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_TenDanhMuc, tenDM)
	ghi(sheetID, "DANH_MUC", targetRow, core.CotDM_DanhMucMe, dmMe) // [MỚI] Ghi cột C
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
		core.ThemThuongHieuVaoRam(newTH) 
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
        // Tự động giải phóng Lock an toàn bằng defer
		func() {
			core.KhoaHeThong.Lock()
			defer core.KhoaHeThong.Unlock()
			found.TenThuongHieu = tenTH
			found.Logo = logo
			found.MoTa = moTa
		}()
	}

	ghi := core.ThemVaoHangCho
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_MaThuongHieu, strings.ToUpper(maTH))
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_TenThuongHieu, tenTH)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_Logo, logo)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_MoTa, moTa)
	ghi(sheetID, "THUONG_HIEU", targetRow, core.CotTH_TrangThai, 1)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Thương hiệu thành công!"})
}
