package auth_store

import (
	"net/http"
	"strings"

	"app/config"
	"app/modules/auth_verify"
	"github.com/gin-gonic/gin"
)

var service = Service{repo: Repo{}}

// Hàm tiện ích: Tách chính xác Domain/Subdomain hiện tại để Cookie không bị tràn sang Shop khác
func getDomain(c *gin.Context) string {
	host := c.Request.Host
	return strings.Split(host, ":")[0]
}

func API_Login(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	ghiNho := c.PostForm("ghi_nho") == "on"

	sessionID, sign, err := service.Login(shopID, dinhDanh, pass, c.Request.UserAgent(), ghiNho)
	if err != nil {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	if ghiNho { maxAge = 30 * 24 * 3600 }
	domain := getDomain(c)
	
	// Set Cookie gắn chặt với domain của Shop đó
	c.SetCookie("session_token", sessionID, maxAge, "/", domain, false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", domain, false, true)
	c.Redirect(http.StatusFound, "/")
}

func API_Register(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")); if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, err := service.Register(shopID, hoTen, user, email, pass, maPin, dienThoai, c.Request.UserAgent())
	
	if err != nil {
		c.HTML(http.StatusOK, "dang_ky_khach_hang", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	domain := getDomain(c)
	c.SetCookie("session_token", sessionID, maxAge, "/", domain, false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", domain, false, true)
	c.Redirect(http.StatusFound, "/") 
}

func API_Logout(c *gin.Context) {
	domain := getDomain(c)
	c.SetCookie("session_token", "", -1, "/", domain, false, true)
	c.SetCookie("session_sign", "", -1, "/", domain, false, true)
	c.Redirect(http.StatusFound, "/login")
}

func API_ResetByPin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	// Chọc thẳng vào auth_verify với cờ "STOREFRONT"
	if err := auth_verify.GlobalService.VerifyPin("STOREFRONT", shopID, dinhDanh, pin); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}

	if err := service.ResetPassword(shopID, dinhDanh, passMoi); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

func API_ResetByOtp(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	otp := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	if err := auth_verify.GlobalService.VerifyOtp("STOREFRONT", shopID, dinhDanh, otp); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}

	if err := service.ResetPassword(shopID, dinhDanh, passMoi); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
