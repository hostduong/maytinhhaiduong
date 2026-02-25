package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/cau_hinh" 
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangDangNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 
	cookie, _ := c.Cookie("session_id")
	
	if cookie != "" {
		if _, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, "dang_nhap", gin.H{"TieuDe": "Đăng Nhập"})
}

func XuLyDangNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 

	inputDinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))
	ghiNho        := c.PostForm("ghi_nho")

	// [MỚI] LỚP TỪ CHỐI 1: Chặn cứng tài khoản Bot Hệ Thống từ vòng gửi xe
	if inputDinhDanh == "admin" {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, inputDinhDanh)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	if !cau_hinh.KiemTraMatKhau(pass, kh.MatKhauHash) {
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

	sessionID := cau_hinh.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	signature := cau_hinh.TaoChuKyBaoMat(sessionID, userAgent)
	
	expTime := time.Now().Add(thoiGianSong).Unix()
	maxAge  := int(thoiGianSong.Seconds())

	core.KhoaHeThong.Lock()
	if kh.RefreshTokens == nil {
		kh.RefreshTokens = make(map[string]core.TokenInfo)
	}
	
	kh.RefreshTokens[sessionID] = core.TokenInfo{
		DeviceName: userAgent,
		ExpiresAt:  expTime,
	}
	core.KhoaHeThong.Unlock()
	
	go func() {
		sID := kh.SpreadsheetID
		if sID == "" { sID = shopID } 
		jsonStr := core.ToJSON(kh.RefreshTokens)
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_RefreshTokenJson, jsonStr)
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
