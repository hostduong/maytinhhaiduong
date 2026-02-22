package chuc_nang

import (
	"net/http"
	"strings"
	"app/cau_hinh"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangQuenMatKhau(c *gin.Context) { c.HTML(http.StatusOK, "quen_mat_khau", gin.H{}) }

func XuLyQuenPassBangPIN(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	pinInput := strings.TrimSpace(c.PostForm("pin"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))
	
	if !cau_hinh.KiemTraMaPin(pinInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN phải đúng 8 chữ số!"})
		return
	}
	if !cau_hinh.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, dinhDanh)
	
	if !ok || !cau_hinh.KiemTraMatKhau(pinInput, kh.MaPinHash) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản hoặc mã PIN không chính xác!"})
		return 
	}
	
	hash, _ := cau_hinh.HashMatKhau(passMoi)
	core.KhoaHeThong.Lock()
	kh.MatKhauHash = hash
	core.KhoaHeThong.Unlock()
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

func XuLyGuiOTPEmail(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	
	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, dinhDanh)
	if !ok { 
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu tài khoản tồn tại, mã OTP sẽ được gửi đến Email đăng ký."})
		return 
	}

	if kh.Email == "" || !strings.Contains(kh.Email, "@") {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản này chưa cập nhật Email, vui lòng dùng PIN."})
		return
	}

	okLimit, msg := core.KiemTraRateLimit(kh.Email)
	if !okLimit { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	code := core.TaoMaOTP6So()
	
	// 👉 Gọi thẳng API gửi mail thật của bạn!
	if err := core.GuiMailXacMinhAPI(kh.Email, code); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi hệ thống gửi mail: " + err.Error()})
		return
	}
	
	cacheKey := shopID + "_" + kh.TenDangNhap
	core.LuuOTP(cacheKey, code)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã OTP đến Email đăng ký của bạn!"})
}

func XuLyQuenPassBangOTP(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	otp      := strings.TrimSpace(c.PostForm("otp"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))

	if !cau_hinh.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, dinhDanh)
	cacheKey := shopID + "_" + kh.TenDangNhap
	
	// Kiểm tra Cache RAM
	if !ok || !core.KiemTraOTP(cacheKey, otp) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP không đúng hoặc đã hết hạn!"})
		return 
	}

	hash, _ := cau_hinh.HashMatKhau(passMoi)
	core.KhoaHeThong.Lock()
	kh.MatKhauHash = hash
	core.KhoaHeThong.Unlock()
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
