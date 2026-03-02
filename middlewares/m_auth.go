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
// 1. TRẠM LỄ TÂN: PHÂN LUỒNG TÊN MIỀN (Chạy cho TẤT CẢ request)
// =================================================================
func IdentifyTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		domain := strings.Split(host, ":")[0]

		shopID := config.BienCauHinh.IdFileSheet 
		theme := "theme_master"

		// Ánh xạ Domain lấy ShopID
		if mappedID, ok := config.MapDomainShop[domain]; ok {
			shopID = mappedID
		} else if strings.HasSuffix(domain, ".99k.vn") && domain != "www.99k.vn" && domain != "99k.vn" {
			// Tương lai: Logic tra cứu Domain động trong Cache Tenant
		}

		c.Set("SHOP_ID", shopID)
		c.Set("THEME", theme)
		c.Next()
	}
}

// =================================================================
// 2. TRẠM BẢO VỆ: KIỂM TRA ĐĂNG NHẬP (Chỉ chạy cho Khu vực Quản trị)
// =================================================================
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID := c.GetString("SHOP_ID") // Đã được Lễ tân cấp ở trên

		cookie, err := c.Cookie("session_token") // Đảm bảo đồng bộ tên cookie
		if err != nil || cookie == "" {
			TuChoiTruyCap(c, http.StatusUnauthorized, "Vui lòng đăng nhập để tiếp tục!")
			return
		}

		lockKH := core.GetSheetLock(shopID, core.TenSheetKhachHang)
		lockKH.RLock()
		
		var user *core.KhachHang
		danhSach := core.CacheKhachHang[shopID]
		for _, kh := range danhSach {
			if info, ok := kh.RefreshTokens[cookie]; ok {
				if time.Now().Unix() <= info.ExpiresAt {
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

		lockPQ := core.GetSheetLock(shopID, core.TenSheetPhanQuyen)
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
