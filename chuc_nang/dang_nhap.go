package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangDangNhap(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := core.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, "dang_nhap", gin.H{"TieuDe": "Đăng Nhập"})
}

func XuLyDangNhap(c *gin.Context) {
	inputDinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))
	ghiNho        := c.PostForm("ghi_nho")

	// 1. Tìm user
	kh, ok := core.TimKhachHangTheoUserOrEmail(inputDinhDanh)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	// 2. Check Pass
	if !bao_mat.KiemTraMatKhau(pass, kh.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Mật khẩu không đúng!"})
		return
	}

	if kh.TrangThai == 0 {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản bị khóa!"})
		return
	}

	// 3. Tạo Session
	var thoiGianSong time.Duration
	if ghiNho == "on" {
		thoiGianSong = 30 * 24 * time.Hour
	} else {
		thoiGianSong = cau_hinh.ThoiGianHetHanCookie
	}

	sessionID := bao_mat.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	signature := bao_mat.TaoChuKyBaoMat(sessionID, userAgent)
	
	expTime := time.Now().Add(thoiGianSong).Unix()
	maxAge  := int(thoiGianSong.Seconds())

	// 4. Cập nhật RAM (Nhanh)
	core.KhoaHeThong.Lock()
	kh.Cookie = sessionID
	kh.CookieExpired = expTime
	core.KhoaHeThong.Unlock()
	
	// 5. Ghi Sheet (Bất đồng bộ - Không block người dùng)
	// Việc này sẽ được Worker xử lý sau, người dùng không cần chờ
	go func() {
		sID := kh.SpreadsheetID
		if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_Cookie, sessionID)
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_CookieExpired, expTime)
	}()

	// 6. Set Cookie & Redirect ngay lập tức
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func DangXuat(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
