package auth_admin

import (
	"net/http"
	"strings"

	"app/config"
	"app/auth/auth_verify"
	"github.com/gin-gonic/gin"
)

var service = Service{repo: Repo{}}

func API_Login(c *gin.Context) {
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	ghiNho := c.PostForm("ghi_nho") == "on"

	sessionID, sign, err := service.Login(dinhDanh, pass, c.Request.UserAgent(), ghiNho)
	if err != nil {
		// [ĐÃ VÁ LỖI]: Gọi đúng tên template giao diện Tím Premium
		c.HTML(http.StatusOK, "dang_nhap_admin", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	if ghiNho { maxAge = 30 * 24 * 3600 }
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)
	c.Redirect(http.StatusFound, "https://shop.99k.vn/tong-quan")
}

func API_Register(c *gin.Context) {
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full"))
	if dienThoai == "" { 
		dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) 
	}

	sessionID, sign, err := service.Register(hoTen, user, email, pass, maPin, dienThoai, c.Request.UserAgent())
	
	if err != nil {
		// [ĐÃ VÁ LỖI]: Gọi đúng tên template giao diện Tím Premium
		c.HTML(http.StatusOK, "dang_ky_admin", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)
	c.Redirect(http.StatusFound, "https://shop.99k.vn/tong-quan") 
}

func API_Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", "", -1, "/", ".99k.vn", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// ===============================================================
// QUÊN MẬT KHẨU - TÍCH HỢP TRẠM KIỂM SOÁT CHUNG AUTH_VERIFY
// ===============================================================
func API_ResetByPin(c *gin.Context) {
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	// 1. Gọi auth_verify để kiểm tra mã PIN với quyền TENANT_ADMIN
	if err := auth_verify.GlobalService.VerifyPin("TENANT_ADMIN", config.BienCauHinh.IdFileSheetAdmin, dinhDanh, pin); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	// 2. Cập nhật mật khẩu mới
	if err := service.ResetPassword(dinhDanh, passMoi); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

func API_ResetByOtp(c *gin.Context) {
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	otp := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	// 1. Gọi auth_verify để kiểm tra OTP với quyền TENANT_ADMIN
	if err := auth_verify.GlobalService.VerifyOtp("TENANT_ADMIN", config.BienCauHinh.IdFileSheetAdmin, dinhDanh, otp); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	// 2. Cập nhật mật khẩu mới
	if err := service.ResetPassword(dinhDanh, passMoi); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
