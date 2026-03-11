package product_master

import (
	"strings"
	"app/config"
	"github.com/gin-gonic/gin"
)

// Dùng chung cho mọi ngành hàng (Lấy 1 SP)
func API_LayChiTietSanPhamMaster(c *gin.Context) {
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không có quyền truy cập!"})
		return
	}

	maSP := c.Param("ma_sp")
	if maSP == "" { c.JSON(200, gin.H{"status": "error", "msg": "Thiếu mã sản phẩm!"}); return }

	data, err := Service_LayChiTietSanPham(adminShopID, maSP)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "data": data})
}

// Dùng chung cho mọi ngành hàng (Lưu/Update SP)
func API_LuuSanPhamMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") 
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	vaiTro := c.GetString("USER_ROLE")
	userID := c.GetString("USER_ID")
	
	maNganh := strings.TrimSpace(c.PostForm("ma_nganh"))
	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	dataJSON := c.PostForm("data_json") // Hứng toàn bộ JSON từ UI vẽ ra

	if maNganh == "" || dataJSON == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Thiếu mã ngành hoặc dữ liệu JSON!"})
		return
	}

	err := Service_LuuSanPham(masterShopID, adminShopID, vaiTro, userID, maNganh, maSP, dataJSON)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}
