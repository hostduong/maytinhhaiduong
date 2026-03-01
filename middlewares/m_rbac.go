package middlewares

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

// RequirePermission: Yêu cầu chính xác 1 mã quyền (VD: "product.create")
func RequirePermission(maChucNang string) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID := c.GetString("SHOP_ID")
		role := c.GetString("USER_ROLE")

		// Đặc quyền tuyệt đối cho Quản trị Lõi
		if role == "quan_tri_he_thong" {
			c.Next()
			return
		}

		lockPQ := core.GetSheetLock(shopID, core.TenSheetPhanQuyen)
		lockPQ.RLock()
		defer lockPQ.RUnlock()

		hasPermission := false
		if shopMap, ok := core.CachePhanQuyen[shopID]; ok {
			if listQuyen, exists := shopMap[role]; exists {
				if allowed, has := listQuyen[maChucNang]; has && allowed {
					hasPermission = true
				}
			}
		}

		if !hasPermission {
			TuChoiTruyCap(c, http.StatusForbidden, "Bạn không có quyền thao tác tính năng này ("+maChucNang+")!")
			return
		}

		c.Next()
	}
}
