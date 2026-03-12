package setup

import (
	"net/http"
	"strings"

	"app/config"
	"github.com/gin-gonic/gin"
)

func API_Setup(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full"))
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, err := Service_ThucThiSetup(shopID, hoTen, user, email, pass, maPin, dienThoai, c.PostForm("ngay_sinh"), c.PostForm("gioi_tinh"), c.Request.UserAgent())
	
	if err != nil {
		c.HTML(http.StatusOK, "setup", gin.H{
			"Loi": err.Error(), 
			"TieuDe": "Khởi Tạo Lõi Hệ Thống",
		})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)

	// Đăng ký xong vút thẳng vào trang Tổng Quan
	c.Redirect(http.StatusFound, "/master/tong-quan")
}
