package auth

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func checkLogin(c *gin.Context) bool {
	cookie, _ := c.Cookie("session_token")
	if cookie != "" {
		shopID := c.GetString("SHOP_ID")
		// Bức tường lửa: Chặn việc RAM trống làm văng user
		_ = core.EnsureKhachHangLoaded(shopID)
		
		// [LOCK CHUẨN]: Tìm user phải RLock
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

	if mode == "MASTER_CORE" {
		c.Redirect(http.StatusFound, "/login")
		return
	} else if mode == "TENANT_ADMIN" {
		host := c.Request.Host
		if strings.HasPrefix(host, "admin.") {
			c.Redirect(http.StatusFound, "https://www.99k.vn/register")
			return
		}
	}

	// CHẺ LUỒNG RENDER GIAO DIỆN Ở ĐÂY
	if mode == "TENANT_STORE" {
		// Dành cho Khách mua lẻ (Không có trường Tên miền)
		c.HTML(http.StatusOK, "dang_ky_khach_hang", gin.H{
			"TieuDe": "Đăng Ký Thành Viên",
			"Loi":    c.Query("loi"),
		})
	} else {
		// Dành cho Chủ Shop tại www.99k.vn (Có trường chọn Tên miền .99k.vn)
		c.HTML(http.StatusOK, "dang_ky", gin.H{
			"TieuDe": "Tạo Cửa Hàng Mới",
			"Loi":    c.Query("loi"),
		})
	}
}

func TrangQuenMatKhau(c *gin.Context) {
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{"TieuDe": "Khôi phục Mật Khẩu"})
}

func TrangXacThucOTP(c *gin.Context) {
	userName := ""
	cookie, _ := c.Cookie("session_token")
	shopID := c.GetString("SHOP_ID")
	
	_ = core.EnsureKhachHangLoaded(shopID)
	
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.RLock()
	for _, kh := range core.CacheKhachHang[shopID] {
		if _, ok := kh.RefreshTokens[cookie]; ok { 
			userName = kh.TenDangNhap
			break
		}
	}
	lock.RUnlock()

	c.HTML(http.StatusOK, "xac_thuc_otp", gin.H{"TieuDe": "Xác thực OTP", "User": userName})
}
