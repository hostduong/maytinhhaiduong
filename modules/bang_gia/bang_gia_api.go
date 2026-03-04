package bang_gia

import (
	"strings"
	"github.com/gin-gonic/gin"
	"app/modules/thanh_toan"
)

var svc = BangGiaService{}

// API_CheckGia: Kiểm tra giá cuối cùng (Zero-Trust)
func API_CheckGia(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	maGoi := c.PostForm("ma_goi")
	maCode := strings.ToUpper(strings.TrimSpace(c.PostForm("ma_code")))

	gia, codeHopLe, _, err := thanh_toan.Svc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status":      "ok",
		"final_price": gia,
		"is_valid":    codeHopLe != "",
	})
}

// API_MuaGoi: Kích hoạt mua gói dịch vụ
func API_MuaGoi(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	maGoi := c.PostForm("ma_goi")
	maCode := strings.ToUpper(strings.TrimSpace(c.PostForm("ma_code")))

	url, err := svc.BuyStarterPackage(masterShopID, userID, maGoi, maCode)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "redirect_url": url})
}
