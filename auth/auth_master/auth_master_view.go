package auth_master

import (
	"net/http"
	"app/core"
	"app/config"
	"github.com/gin-gonic/gin"
)

func checkLogin(c *gin.Context) bool {
	cookie, _ := c.Cookie("session_token")
	if cookie != "" {
		masterID := config.BienCauHinh.IdFileSheetMaster
		_ = core.EnsureKhachHangLoaded(masterID)
		
		lock := core.GetSheetLock(masterID, core.TenSheetKhachHangMaster)
		lock.RLock()
		defer lock.RUnlock()
		
		for _, kh := range core.CacheKhachHang[masterID] {
			if _, ok := kh.RefreshTokens[cookie]; ok { return true }
		}
	}
	return false
}

func TrangDangNhap(c *gin.Context) {
	if checkLogin(c) { c.Redirect(http.StatusFound, "/master/tong-quan"); return }
	// Gọi file: themes/template_master/dang_nhap.html
	c.HTML(http.StatusOK, "dang_nhap_master", gin.H{"TieuDe": "Đăng Nhập Quản Trị Lõi (Master)"})
}

func TrangQuenMatKhau(c *gin.Context) {
	// Gọi file: themes/template_master/quen_mat_khau.html
	c.HTML(http.StatusOK, "quen_mat_khau_master", gin.H{"TieuDe": "Khôi phục Mật Khẩu Lõi"})
}
