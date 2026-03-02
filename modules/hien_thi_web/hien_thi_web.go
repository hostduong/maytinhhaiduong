package hien_thi_web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Giữ nguyên bộ hàm Format HTML
func LayBoHamHTML() template.FuncMap {
	return template.FuncMap{
		"firstImg": func(s string) string {
			if s == "" { return "" }
			parts := strings.Split(s, "|")
			return strings.TrimSpace(parts[0])
		},
		"format_money": func(n float64) string {
			p := message.NewPrinter(language.Vietnamese)
			return p.Sprintf("%.0f", n)
		},
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"split": strings.Split,
	}
}

// Render Giao diện Mặt tiền
func TrangChu(c *gin.Context) {
	c.HTML(http.StatusOK, "trang_chu", gin.H{"TieuDe": "Nền tảng vận hành 99K"})
}

func TrangDangNhap(c *gin.Context) {
	c.HTML(http.StatusOK, "dang_nhap", gin.H{"TieuDe": "Đăng Nhập"})
}

func TrangDangKy(c *gin.Context) {
	c.HTML(http.StatusOK, "dang_ky", gin.H{"TieuDe": "Đăng Ký Tài Khoản Mới"})
}

func TrangQuenMatKhau(c *gin.Context) {
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{"TieuDe": "Khôi phục Mật Khẩu"})
}

func TrangXacThucOTP(c *gin.Context) {
	c.HTML(http.StatusOK, "xac_thuc_otp", gin.H{"TieuDe": "Xác thực OTP"})
}
