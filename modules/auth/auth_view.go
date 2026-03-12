package auth

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
	"strings"
)

func checkLogin(c *gin.Context) bool {
	cookie, _ := c.Cookie("session_token")
	if cookie != "" {
		shopID := c.GetString("SHOP_ID")
		_ = core.EnsureKhachHangLoaded(shopID)
		
		lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
		lock.RLock()
		defer lock.RUnlock()
		
		for _, kh := range core.CacheKhachHang[shopID] {
			if _, ok := kh.RefreshTokens[cookie]; ok { 
				return true 
			}
		}
	}
	return false
}

func TrangDangNhap(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "/"); return }
	c.HTML(http.StatusOK, "dang_nhap", gin.H{"TieuDe": "Đăng Nhập"})
}

func TrangDangKy(c *gin.Context) {
	if checkLogin(c) {
		return 
	}

	mode := c.GetString("APP_MODE")

	// [ĐÃ MỞ KHÓA]: Bỏ chặn if mode == "MASTER_CORE" ở đây

	if mode == "TENANT_ADMIN" {
		host := c.Request.Host
		if strings.HasPrefix(host, "admin.") {
			c.Redirect(http.StatusFound, "https://www.99k.vn/register")
			return
		}
	}

	if mode == "TENANT_STORE" {
		c.HTML(http.StatusOK, "dang_ky_khach_hang", gin.H{
			"TieuDe": "Đăng Ký Thành Viên",
			"Loi":    c.Query("loi"),
		})
	} else {
		// Form này giờ sẽ dùng chung cho cả Khách hàng (www) và Chủ Tịch (sss)
		tieuDe := "Tạo Cửa Hàng Mới"
		if mode == "MASTER_CORE" { tieuDe = "Đăng Ký Sáng Lập Viên (Master)" }

		c.HTML(http.StatusOK, "dang_ky_cua_hang", gin.H{
			"TieuDe": tieuDe,
			"Loi":    c.Query("loi"),
		})
	}
}
