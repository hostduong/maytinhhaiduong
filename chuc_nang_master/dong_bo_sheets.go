package chuc_nang_master

import (
	"net/http"
	"strings"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangDongBoSheetsMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	me, _ := core.LayKhachHang(masterShopID, userID)

	c.HTML(http.StatusOK, "master_dong_bo_sheets", gin.H{
		"TieuDe":   "Đồng Bộ Master",
		"NhanVien": me,
		"QuyenHan": vaiTro,
	})
}

func API_NapLaiDuLieuMasterCoPIN(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	me, _ := core.LayKhachHang(masterShopID, userID)

	if me.MaPinHash == "" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Bạn chưa thiết lập mã PIN bảo mật trong phần Hồ sơ!"})
		return
	}

	if !cau_hinh.KiemTraMatKhau(pinXacNhan, me.MaPinHash) {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}

	go func() {
		core.HeThongDangBan = true
		// Nạp lại dữ liệu Lõi
		core.NapPhanQuyen(masterShopID) 
		core.NapKhachHang(masterShopID)
		core.NapTinNhan(masterShopID)
		core.HeThongDangBan = false
	}()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok", 
		"msg": "Xác thực PIN thành công. Đang tải lại dữ liệu ngầm...",
	})
}
