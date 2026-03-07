package may_tinh_master

import (
	"strings"

	"app/config"
	"github.com/gin-gonic/gin"
)

func API_LayChiTietMayTinhMaster(c *gin.Context) {
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	maSP := c.Param("ma_sp")
	if maSP == "" { c.JSON(200, gin.H{"status": "error", "msg": "Thiếu mã sản phẩm!"}); return }

	data, err := Service_LayChiTietMayTinh(adminShopID, maSP)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "data": data})
}

func API_LuuMayTinhMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") 
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	vaiTro := c.GetString("USER_ROLE")
	userID := c.GetString("USER_ID")
	
	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	dataJSON := c.PostForm("data_skus")

	err := Service_LuuMayTinh(masterShopID, adminShopID, vaiTro, userID, maSP, dataJSON)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu thành công!"})
}
