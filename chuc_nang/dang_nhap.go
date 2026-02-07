package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func TrangDangNhap(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, "dang_nhap", gin.H{})
}

func XuLyDangNhap(c *gin.Context) {
	// Nhận input đa năng (Mã KH / User / Email)
	inputDinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))
	ghiNho        := c.PostForm("ghi_nho") // Checkbox: "on" hoặc ""

	// 1. Tìm user
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(inputDinhDanh)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	// 2. Kiểm tra mật khẩu
	if !bao_mat.KiemTraMatKhau(pass, kh.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Mật khẩu không đúng!"})
		return
	}

	// 3. Kiểm tra trạng thái
	if kh.TrangThai == 0 {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản đã bị khóa vĩnh viễn!"})
		return
	}
	if kh.TrangThai == 2 {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản đang bị tạm khóa!"})
		return
	}

	// 4. Xử lý "Ghi nhớ đăng nhập"
	var thoiGianSong time.Duration
	if ghiNho == "on" {
		thoiGianSong = 30 * 24 * time.Hour // 30 ngày
	} else {
		thoiGianSong = cau_hinh.ThoiGianHetHanCookie // Mặc định (30 phút)
	}

	// 5. Tạo Session & Chữ ký
	sessionID := bao_mat.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent() // Lấy thông tin trình duyệt khách
	signature := bao_mat.TaoChuKyBaoMat(sessionID, userAgent)
	
	expTime := time.Now().Add(thoiGianSong).Unix()
	maxAge  := int(thoiGianSong.Seconds())

	// 6. Cập nhật vào Struct trong RAM & Ghi Sheet
	kh.Cookie = sessionID
	kh.CookieExpired = expTime
	nghiep_vu.CapNhatPhienDangNhapKH(kh)

	// 7. Set 2 Cookie xuống trình duyệt
	// Cookie 1: Session ID (Chìa khóa)
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	
	// Cookie 2: Signature (Ổ khóa bảo vệ)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func DangXuat(c *gin.Context) {
	// Xóa sạch cả 2 cookie
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
