package middlewares

import (
	"net/http"
	"app/core"

	"github.com/gin-gonic/gin"
)

// CheckSaaSLimit: Trạm kiểm soát tài nguyên gói cước
func CheckSaaSLimit(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID := c.GetString("SHOP_ID")
		
		if resourceType == "san_pham" {
			maxSP := core.LayGioiHanSanPhamCuaShop(shopID)
			if maxSP != -1 {
				// Khóa bộ nhớ siêu tốc O(1) mới để đếm tổng sản phẩm
				lock := core.GetSheetLock(shopID, "PRODUCTS_CACHE")
				lock.RLock()
				currentCount := 0
				if nganhMap, exists := core.CacheSanPham[shopID]; exists {
					for _, ds := range nganhMap {
						currentCount += len(ds)
					}
				}
				lock.RUnlock()

				if currentCount >= maxSP {
					TuChoiTruyCap(c, http.StatusPaymentRequired, "Gói dịch vụ của bạn đã đạt giới hạn tối đa. Vui lòng nâng cấp gói cước!")
					return
				}
			}
		}
		
		c.Next()
	}
}
