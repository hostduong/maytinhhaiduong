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
	if checkLogin(c) { c.Redirect(http.StatusFound, "https://admin.99k.vn/tong-quan"); return }
	// Render file /themes/template_admin/dang_nhap.html
	c.HTML(http.StatusOK, "dang_nhap_admin", gin.H{"TieuDe": "Đăng Nhập Quản Trị Hệ Thống"})
}

func TrangDangKy(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "https://admin.99k.vn/tong-quan"); return }
	// Render file /themes/default/dang_ky.html
	c.HTML(http.StatusOK, "dang_ky", gin.H{"TieuDe": "Mở Cửa Hàng Mới"})
}

func TrangQuenMatKhau(c *gin.Context) {
	// Render file /themes/template_admin/quen_mat_khau.html
	c.HTML(http.StatusOK, "quen_mat_khau_admin", gin.H{"TieuDe": "Khôi phục Mật Khẩu Admin"})
}
