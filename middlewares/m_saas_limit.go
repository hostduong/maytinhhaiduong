package middlewares

import (
	"github.com/gin-gonic/gin"
)

// CheckSaaSLimit: Trạm kiểm soát tài nguyên gói cước
func CheckSaaSLimit(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// shopID := c.GetString("SHOP_ID")
		
		// TODO: Tương lai sẽ móc vào bảng cấu hình gói cước của Master để check.
		// Ví dụ: KiemTraGioiHan(shopID, "so_luong_san_pham")
		// Nếu vượt quá:
		// TuChoiTruyCap(c, http.StatusPaymentRequired, "Gói dịch vụ của bạn đã hết dung lượng. Vui lòng nâng cấp!")
		// return
		
		// Hiện tại pass qua cửa này
		c.Next()
	}
}
