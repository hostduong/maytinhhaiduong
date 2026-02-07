package chuc_nang

import (
	"net/http"
	"strings"
	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"
	"github.com/gin-gonic/gin"
)

func TrangQuenMatKhau(c *gin.Context) { c.HTML(http.StatusOK, "quen_mat_khau", gin.H{}) }

// [CÁCH 1]: Đổi mật khẩu bằng Mã PIN
func XuLyQuenPassBangPIN(c *gin.Context) {
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	pinInput := strings.TrimSpace(c.PostForm("pin"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))
	
	// Validate
	if !bao_mat.KiemTraMaPin(pinInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN phải đúng 8 chữ số!"})
		return
	}
	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(dinhDanh)
	
	// Kiểm tra PIN (Dùng hàm so sánh Hash an toàn)
	if !ok || !bao_mat.KiemTraMatKhau(pinInput, kh.MaPinHash) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản hoặc mã PIN không chính xác!"})
		return 
	}
	
	// Hash mật khẩu mới và lưu
	hash, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hash
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hash)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

// [CÁCH 2 - BƯỚC 1]: Gửi OTP (Người dùng nhập User -> Hệ thống tìm Email -> Gửi)
func XuLyGuiOTPEmail(c *gin.Context) {
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(dinhDanh)
	if !ok { 
		// Fake thành công để tránh dò User
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu tài khoản tồn tại, mã OTP sẽ được gửi đến Email đăng ký."})
		return 
	}

	// Kiểm tra xem tài khoản có Email không
	if kh.Email == "" || !strings.Contains(kh.Email, "@") {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản này chưa cập nhật Email, vui lòng dùng PIN."})
		return
	}

	// Kiểm tra Rate Limit trên Email thực tế
	okLimit, msg := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !okLimit { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	code := nghiep_vu.TaoMaOTP6So()
	
	// Gửi mail
	if err := nghiep_vu.GuiMailXacMinhAPI(kh.Email, code); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi hệ thống gửi mail: " + err.Error()})
		return
	}
	
	// Lưu OTP vào Cache (Key là Tên Đăng Nhập để bước sau đối chiếu)
	nghiep_vu.LuuOTP(kh.TenDangNhap, code)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã OTP đến Email đăng ký của bạn!"})
}

// [CÁCH 2 - BƯỚC 2]: Xác nhận OTP và Đổi Pass
func XuLyQuenPassBangOTP(c *gin.Context) {
	// Người dùng gửi lại định danh để tìm lại User
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	otp      := strings.TrimSpace(c.PostForm("otp"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))

	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(dinhDanh)
	// Key OTP lưu theo TenDangNhap, nên phải check đúng key đó
	if !ok || !nghiep_vu.KiemTraOTP(kh.TenDangNhap, otp) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP không đúng hoặc đã hết hạn!"})
		return 
	}

	hash, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hash
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
