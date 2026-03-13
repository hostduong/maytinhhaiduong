package auth_master

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
		c.HTML(http.StatusOK, "dang_nhap_master", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	if ghiNho { maxAge = 30 * 24 * 3600 }
	// Cookie đặt cho tên miền sss.99k.vn để tăng tính bảo mật
	c.SetCookie("session_token", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", "", false, true)
	
	c.Redirect(http.StatusFound, "/master/tong-quan")
}

func API_Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// Sử dụng lại sức mạnh của Trạm xác thực trung tâm (auth_verify)
func API_ResetByPin(c *gin.Context) {
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	// Ép cứng APP_MODE = "MASTER_CORE" khi gọi sang auth_verify
	if err := auth_verify.GlobalService.VerifyPin("MASTER_CORE", config.BienCauHinh.IdFileSheetMaster, dinhDanh, pin); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}

	if err := service.ResetPassword(dinhDanh, passMoi); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đổi mật khẩu Master thành công!"})
}

func API_ResetByOtp(c *gin.Context) {
	dinhDanh := strings.TrimSpace(c.PostForm("dinh_danh"))
	otp := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	if err := auth_verify.GlobalService.VerifyOtp("MASTER_CORE", config.BienCauHinh.IdFileSheetMaster, dinhDanh, otp); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}

	if err := service.ResetPassword(dinhDanh, passMoi); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": err.Error()}); return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Đổi mật khẩu Master thành công!"})
}
