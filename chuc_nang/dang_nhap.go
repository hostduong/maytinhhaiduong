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

	kh, ok := core.TimKhachHangTheoUserOrEmail(inputDinhDanh)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	if !bao_mat.KiemTraMatKhau(pass, kh.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Mật khẩu không đúng!"})
		return
	}

	if kh.TrangThai == 0 {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản bị khóa!"})
		return
	}

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

	core.KhoaHeThong.Lock()
	kh.Cookie = sessionID
	kh.CookieExpired = expTime
	core.KhoaHeThong.Unlock()
	
	// [QUAN TRỌNG] Ghi Sheet bằng Goroutine để không block người dùng
	go func() {
		sID := kh.SpreadsheetID
		if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_Cookie, sessionID)
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_CookieExpired, expTime)
	}()

	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func DangXuat(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
