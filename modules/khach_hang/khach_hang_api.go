package khach_hang

import (
	"strings"
	"github.com/gin-gonic/gin"

	"app/modules/thanh_toan" // Gọi sang module thanh toán tập trung
)

// Khai báo service để sử dụng nội bộ trong module
var svc = CustomerService{}

// API_CheckPrice: Khách hàng hỏi "Nếu tôi nhập mã này thì giá bao nhiêu?"
func API_CheckPrice(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") // [cite: 17]
	maGoi := c.PostForm("ma_goi")
	maCode := strings.ToUpper(strings.TrimSpace(c.PostForm("ma_code")))

	// Chuyển việc tính toán cho module thanh toán chuyên dụng
	finalPrice, codeHopLe, _, err := thanh_toan.Svc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status":      "ok",
		"final_price": finalPrice,
		"is_valid":    codeHopLe != "",
	})
}

// API_MuaGoiKhachHang: Khách hàng chốt "Tôi mua gói này với mã này"
func API_MuaGoiKhachHang(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	maGoi := c.PostForm("ma_goi")
	maCode := c.PostForm("ma_code")

	// Gọi Service thực thi
	redirectURL, err := svc.BuyStarterPackage(masterShopID, userID, maGoi, maCode)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	// Trả về URL: https://www.99k.vn/admin/database
	c.JSON(200, gin.H{
		"status":       "ok",
		"redirect_url": redirectURL,
	})
}
