package thanh_toan

import (
	"strings"
	"github.com/gin-gonic/gin"
)

var Svc = PaymentService{}

func API_CheckPrice(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	maGoi := c.PostForm("ma_goi")
	maCode := strings.ToUpper(strings.TrimSpace(c.PostForm("ma_code")))

	gia, codeDung, _, err := Svc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status":      "ok",
		"final_price": gia,
		"is_valid":    codeDung != "",
	})
}
