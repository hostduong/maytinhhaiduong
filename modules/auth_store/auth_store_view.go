package auth_store

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
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
			if _, ok := kh.RefreshTokens[cookie]; ok { return true }
		}
	}
	return false
}

func TrangDangNhap(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "/"); return }
	// Gọi file: /themes/default/dang_nhap.html (Dùng header mặc định)
	c.HTML(http.StatusOK, "dang_nhap", gin.H{"TieuDe": "Đăng Nhập Khách Hàng"})
}

func TrangDangKy(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "/"); return }
	// Gọi file: /themes/default/dang_ky_khach_hang.html
	c.HTML(http.StatusOK, "dang_ky_khach_hang", gin.H{"TieuDe": "Đăng Ký Thành Viên"})
}

func TrangQuenMatKhau(c *gin.Context) {
	// Gọi file: /themes/default/quen_mat_khau.html
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{"TieuDe": "Khôi phục Mật Khẩu"})
}
