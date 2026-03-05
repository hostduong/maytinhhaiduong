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
				// Cần khóa RLock để đếm số lượng hiện tại một cách an toàn
				lock := core.GetSheetLock(shopID, core.TenSheetMayTinh)
				lock.RLock()
				currentCount := len(core.CacheSanPhamMayTinh[shopID])
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
