package auth

import (
	"net/http"
	"strings"

	"app/config"
	"github.com/gin-gonic/gin"
)

var service = Service{repo: Repo{}}

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
	c.SetCookie("session_token", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}

func API_Register(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME")
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")); if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, vaiTro, err := service.Register(shopID, theme, hoTen, user, email, pass, maPin, dienThoai, c.PostForm("ngay_sinh"), c.PostForm("gioi_tinh"), c.Request.UserAgent())
	
	if err != nil {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_token", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", "", false, true)

	if vaiTro == "quan_tri_he_thong" || vaiTro == "quan_tri_vien_he_thong" || vaiTro == "quan_tri_vien" {
		c.Redirect(http.StatusFound, "/master/tong-quan")
	} else if theme == "theme_master" { c.Redirect(http.StatusFound, "/verify")
	} else { c.Redirect(http.StatusFound, "/") }
}

func API_Verify(c *gin.Context) {
	if err := service.VerifyOTPAndActivate(c.GetString("SHOP_ID"), strings.TrimSpace(c.PostForm("dinh_danh")), strings.TrimSpace(c.PostForm("otp"))); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Xác thực thành công! Hệ thống đang khởi tạo Tên miền."})
}

func API_SendOtp(c *gin.Context) {
	if err := service.SendOtp(c.GetString("SHOP_ID"), strings.TrimSpace(c.PostForm("dinh_danh"))); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã OTP đến Email đăng ký của bạn!"})
}

func API_ResetByOtp(c *gin.Context) {
	if err := service.ResetByOtp(c.GetString("SHOP_ID"), strings.TrimSpace(c.PostForm("dinh_danh")), strings.TrimSpace(c.PostForm("otp")), strings.TrimSpace(c.PostForm("pass_moi"))); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

func API_ResetByPin(c *gin.Context) {
	if err := service.ResetByPin(c.GetString("SHOP_ID"), strings.TrimSpace(c.PostForm("dinh_danh")), strings.TrimSpace(c.PostForm("pin")), strings.TrimSpace(c.PostForm("pass_moi"))); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

func API_Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
