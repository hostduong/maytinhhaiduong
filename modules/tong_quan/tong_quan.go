package tong_quan

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

// ==========================================================
// 1. TRANG TỔNG QUAN HỆ THỐNG LÕI (DASHBOARD)
// ==========================================================
func TrangTongQuanMaster(c *gin.Context) {
	// Lấy ID của Shop Mẹ (Master)
	masterShopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	// Chặn quyền ngặt nghèo: Chỉ Cấp 1 và Cấp 2 mới được vào
	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Đẩy ra View (Sử dụng đúng define "master_tong_quan" của HTML)
	c.HTML(http.StatusOK, "master_tong_quan", gin.H{
		"TieuDe":   "Tổng quan Master",
		"NhanVien": me,
		"QuyenHan": vaiTro,
	})
}

// ==========================================================
// 2. API ĐỒNG BỘ DỮ LIỆU LÕI
// ==========================================================
func API_NapLaiDuLieuMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	go func() {
		core.HeThongDangBan = true
		
		// Chỉ nạp lại Sheet Master
		core.NapPhanQuyen(masterShopID) 
		core.NapKhachHang(masterShopID)
		core.NapTinNhan(masterShopID)

		core.NapDanhMuc(masterShopID)
		core.NapThuongHieu(masterShopID)
		core.NapBienLoiNhuan(masterShopID)
		core.NapMayTinh(masterShopID)
		
		core.HeThongDangBan = false
	}()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok", 
		"msg": "Đang tiến hành đồng bộ dữ liệu Core System...",
	})
}
