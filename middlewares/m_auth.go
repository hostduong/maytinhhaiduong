package middlewares

import (
	"net/http"
	"strings"
	"time"

	"app/config"
	"app/core"

	"github.com/gin-gonic/gin"
)

// =================================================================
// 1. TRẠM LỄ TÂN: NGÃ BA ĐỊNH TUYẾN 3 TẦNG (Chạy tốc độ O(1))
// =================================================================
func IdentifyTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		domain := strings.Split(host, ":")[0]

		var appMode, theme, shopID string

		// TẦNG 1: VÙNG TUYỆT MẬT (Master)
		if domain == "sss.99k.vn" {
			appMode = "MASTER_CORE"
			theme = "template_master"
			shopID = config.BienCauHinh.IdFileSheetMaster

		// TẦNG 2: VÙNG QUẢN TRỊ SHOP (Admin & Marketing)
		} else if domain == "admin.99k.vn" || domain == "www.99k.vn" || domain == "localhost" {
			appMode = "TENANT_ADMIN"
			theme = "template_admin"
			if domain == "www.99k.vn" { theme = "default" } // www dùng giao diện public
			shopID = config.BienCauHinh.IdFileSheetAdmin

		// TẦNG 3: VÙNG TIỀN TUYẾN (Cửa hàng Subdomain)
		} else {
			appMode = "STOREFRONT"
			theme = "default"
			
			core.KhoaHeThong.RLock() 
			id, exists := core.CacheDomainToSheetID[domain]
			core.KhoaHeThong.RUnlock()

			if exists {
				shopID = id
			} else {
				TuChoiTruyCap(c, http.StatusNotFound, "Không tìm thấy Cửa hàng này trên hệ thống!")
				return
			}
		}

		c.Set("APP_MODE", appMode)
		c.Set("THEME", theme)
		c.Set("SHOP_ID", shopID)

		if appMode == "STOREFRONT" && shopID != "" {
			core.DanhDauTruyCapShop(shopID)
		}

		c.Next()
	}
}

// =================================================================
// 2. TRẠM BẢO VỆ: KIỂM TRA ĐĂNG NHẬP (Chỉ chạy cho Khu vực Quản trị)
// =================================================================
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID := c.GetString("SHOP_ID")
		appMode := c.GetString("APP_MODE") 

		cookie, err := c.Cookie("session_token")
		if err != nil || cookie == "" {
			TuChoiTruyCap(c, http.StatusUnauthorized, "Vui lòng đăng nhập để tiếp tục!")
			return
		}

		sheetKH := core.TenSheetKhachHang
		sheetPQ := core.TenSheetPhanQuyen
		if appMode == "MASTER_CORE" {
			sheetKH = core.TenSheetKhachHangMaster
			sheetPQ = core.TenSheetPhanQuyenMaster
		} else if appMode == "TENANT_ADMIN" {
			sheetKH = core.TenSheetKhachHangAdmin
			sheetPQ = core.TenSheetPhanQuyenAdmin
		}

		lockKH := core.GetSheetLock(shopID, sheetKH)
		lockKH.RLock()
		
		var user *core.KhachHang
		danhSach := core.CacheKhachHang[shopID]
		for _, kh := range danhSach {
			if info, ok := kh.RefreshTokens[cookie]; ok {
				// [ĐÃ FIX]: Sử dụng .Exp thay cho .ExpiresAt theo chuẩn JSON NoSQL
				if time.Now().Unix() <= info.Exp {
					user = kh
					break
				}
			}
		}
		lockKH.RUnlock()

		if user == nil {
			TuChoiTruyCap(c, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ hoặc đã hết hạn!")
			return
		}

		if user.TrangThai != 1 && user.MaKhachHang != "0000000000000000001" {
			TuChoiTruyCap(c, http.StatusForbidden, "Tài khoản của bạn đã bị khóa hoặc ngừng giao dịch!")
			return
		}

		lockPQ := core.GetSheetLock(shopID, sheetPQ)
		lockPQ.RLock()
		userLevel := 9
		if user.MaKhachHang == "0000000000000000001" || user.VaiTroQuyenHan == "quan_tri_he_thong" {
			userLevel = 0
		} else {
			for _, v := range core.CacheDanhSachVaiTro[shopID] {
				if v.MaVaiTro == user.VaiTroQuyenHan {
					userLevel = v.StyleLevel
					break
				}
			}
		}
		lockPQ.RUnlock()

		c.Set("USER_ID", user.MaKhachHang)
		c.Set("USER_ROLE", user.VaiTroQuyenHan)
		c.Set("USER_LEVEL", userLevel)

		c.Next()
	}
}

func TuChoiTruyCap(c *gin.Context, code int, msg string) {
	if strings.Contains(c.Request.URL.Path, "/api/") {
		c.AbortWithStatusJSON(code, gin.H{"status": "error", "msg": msg})
	} else {
		if code == http.StatusUnauthorized {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(code, "<h3>⛔ "+msg+"</h3><a href='/'>Về trang chủ</a>")
			c.Abort()
		}
	}
}
