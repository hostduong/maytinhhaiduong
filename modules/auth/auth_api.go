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
	appMode := c.GetString("APP_MODE") // Lấy từ Middleware IdentifyTenant
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")); if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, _, err := service.Register(appMode, shopID, hoTen, user, email, pass, maPin, dienThoai, c.PostForm("ngay_sinh"), c.PostForm("gioi_tinh"), c.Request.UserAgent())
	
	if err != nil {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	
	// [QUAN TRỌNG]: Cookie Domain là ".99k.vn" để đăng nhập thông suốt giữa www và admin
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)

	// ĐIỀU HƯỚNG CHUẨN MỰC THEO YÊU CẦU CỦA SẾP
	if appMode == "TENANT_ADMIN" { 
		// Đăng ký ở www.99k.vn -> Khách của sếp -> Bay thẳng vào Admin làm việc
		c.Redirect(http.StatusFound, "https://admin.99k.vn") 
	} else { 
		// Đăng ký ở cuahang.99k.vn -> Khách của cửa hàng -> Trả về trang chủ cửa hàng đó
		c.Redirect(http.StatusFound, "/") 
	}
}

// ================= CÁC API PHỤC VỤ QUÊN MẬT KHẨU =================

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
