package middlewares

import (
	"net/http"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// CheckAuth: Xác thực danh tính và trạng thái tài khoản
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Lấy ShopID từ Gateway (Giả định đã có middleware Gateway bắt Domain set vào đây)
		shopID := c.GetString("SHOP_ID")
		if shopID == "" {
			// Fallback an toàn nếu chưa cấu hình Domain
			shopID = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
			c.Set("SHOP_ID", shopID)
		}

		cookie, err := c.Cookie("session_token")
		if err != nil || cookie == "" {
			TuChoiTruyCap(c, http.StatusUnauthorized, "Vui lòng đăng nhập để tiếp tục!")
			return
		}

		// 2. Lấy Khóa an toàn của Sheet Khách Hàng và tra cứu
		lockKH := core.GetSheetLock(shopID, core.TenSheetKhachHang)
		lockKH.RLock()
		
		// Tìm user dựa trên Cookie (Logic này bạn tự map với cấu trúc cũ)
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

		// LỚP 2: KIỂM TRA TRẠNG THÁI HOẠT ĐỘNG
		// 1 = Hoạt động, 0 = Tạm khóa, -1 = Đợi xóa
		if user.TrangThai != 1 && user.MaKhachHang != "0000000000000000001" {
			TuChoiTruyCap(c, http.StatusForbidden, "Tài khoản của bạn đã bị khóa hoặc ngừng giao dịch!")
			return
		}

		// Tính toán Level hiện tại của User
		lockPQ := core.GetSheetLock(shopID, core.TenSheetPhanQuyen)
		lockPQ.RLock()
		userLevel := 9 // Mặc định thấp nhất
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

		// 3. Đóng dấu Thân phận vào Request (Context) để truyền vào trong
		c.Set("USER_ID", user.MaKhachHang)
		c.Set("USER_ROLE", user.VaiTroQuyenHan)
		c.Set("USER_LEVEL", userLevel)

		c.Next() // Cấp phép đi qua cửa
	}
}

// Hàm tiện ích: Tự động rẽ nhánh trả về JSON hoặc chuyển hướng HTML
func TuChoiTruyCap(c *gin.Context, code int, msg string) {
	if strings.Contains(c.Request.URL.Path, "/api/") {
		c.AbortWithStatusJSON(code, gin.H{"status": "error", "msg": msg})
	} else {
		if code == http.StatusUnauthorized {
			c.Redirect(http.StatusFound, "/login")
		} else {
			c.AbortWithStatusHTML(code, "<h3>⛔ "+msg+"</h3><a href='/'>Về trang chủ</a>")
		}
	}
}
