package auth_verify

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Khởi tạo một lần duy nhất để giữ RAM cho OTP không bị mất
var repo = NewRepo()
var service = NewService(repo)

func API_SendOtp(c *gin.Context) {
	appMode := c.GetString("APP_MODE")
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))

	if err := service.SendOtp(appMode, shopID, dinhDanh); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đã gửi mã OTP bảo mật đến Email của bạn!"})
}

func API_CheckOtp(c *gin.Context) {
	appMode := c.GetString("APP_MODE")
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	otp := strings.TrimSpace(c.PostForm("otp"))

	if err := service.VerifyOtp(appMode, shopID, dinhDanh, otp); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Xác thực OTP thành công!"})
}

func API_CheckPin(c *gin.Context) {
	appMode := c.GetString("APP_MODE")
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	pin := strings.TrimSpace(c.PostForm("pin"))

	if err := service.VerifyPin(appMode, shopID, dinhDanh, pin); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Xác thực mã PIN thành công!"})
}
