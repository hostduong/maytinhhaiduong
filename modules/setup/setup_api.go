package setup

import (
	"strings"

	"app/config"
	"github.com/gin-gonic/gin"
)

func API_Setup(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	
	// Ép kiểu và gọt khoảng trắng thừa ngay từ cửa
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full"))
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, err := Service_ThucThiSetup(shopID, hoTen, user, email, pass, maPin, dienThoai, c.PostForm("ngay_sinh"), c.PostForm("gioi_tinh"), c.Request.UserAgent())
	
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)

	c.JSON(200, gin.H{
		"status": "ok", 
		"msg": "Thần Sáng Lập đã giáng trần! Dữ liệu Lõi đã sẵn sàng.",
		"redirect": "/master/tong-quan",
	})
}
