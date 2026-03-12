package dang_ky_master

import (
	"net/http"
	"strings"

	"app/config"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangDangKyMaster(c *gin.Context) {
	// Kiểm tra xem đã đăng nhập chưa
	cookie, _ := c.Cookie("session_token")
	if cookie != "" {
		shopID := c.GetString("SHOP_ID")
		_ = core.EnsureKhachHangLoaded(shopID)
		
		lock := core.GetSheetLock(shopID, core.TenSheetKhachHangMaster)
		lock.RLock()
		isLogged := false
		for _, kh := range core.CacheKhachHang[shopID] {
			if _, ok := kh.RefreshTokens[cookie]; ok { isLogged = true; break }
		}
		lock.RUnlock()
		
		if isLogged {
			c.Redirect(http.StatusFound, "/master/tong-quan")
			return
		}
	}

	// Render giao diện riêng biệt của Master
	c.HTML(http.StatusOK, "dang_ky_master", gin.H{
		"TieuDe": "Khởi Tạo Lõi Hệ Thống",
		"Loi":    c.Query("loi"),
	})
}

func API_RegisterMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full"))
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, err := Service_DangKyMaster(shopID, hoTen, user, email, pass, maPin, dienThoai, c.PostForm("ngay_sinh"), c.PostForm("gioi_tinh"), c.Request.UserAgent())
	
	if err != nil {
		c.HTML(http.StatusOK, "dang_ky_master", gin.H{"Loi": err.Error(), "TieuDe": "Khởi Tạo Lõi Hệ Thống"})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)

	c.Redirect(http.StatusFound, "/master/tong-quan")
}
