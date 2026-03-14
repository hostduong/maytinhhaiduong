package auth_admin

import (
	"net/http"
	"app/core"
	"app/config"
	"github.com/gin-gonic/gin"
)

func checkLogin(c *gin.Context) bool {
	cookie, _ := c.Cookie("session_token")
	if cookie != "" {
		adminID := config.BienCauHinh.IdFileSheetAdmin
		_ = core.EnsureKhachHangLoaded(adminID)
		
		lock := core.GetSheetLock(adminID, core.TenSheetKhachHangAdmin)
		lock.RLock()
		defer lock.RUnlock()
		
		for _, kh := range core.CacheKhachHang[adminID] {
			if _, ok := kh.RefreshTokens[cookie]; ok { return true }
		}
	}
	return false
}

func TrangDangNhap(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "https://shop.99k.vn/tong-quan"); return }
	// Render file /themes/template_admin/dang_nhap_admin.html
	c.HTML(http.StatusOK, "dang_nhap_admin", gin.H{"TieuDe": "Đăng Nhập Hệ Thống"})
}

func TrangDangKy(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "https://shop.99k.vn/tong-quan"); return }
	// Render file /themes/template_admin/dang_ky_admin.html
	c.HTML(http.StatusOK, "dang_ky_admin", gin.H{"TieuDe": "Đăng Ký Tài Khoản"})
}

func TrangQuenMatKhau(c *gin.Context) {
	// Render file /themes/template_admin/quen_mat_khau_admin.html
	c.HTML(http.StatusOK, "quen_mat_khau_admin", gin.H{"TieuDe": "Khôi Phục Mật Khẩu"})
}
