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
		appMode := c.GetString("APP_MODE") // Lấy App Mode

		if role == "quan_tri_he_thong" {
			c.Next()
			return
		}

		// [ĐÃ FIX LỖI LOCK MISMATCH]
		sheetPQ := core.TenSheetPhanQuyen
		if appMode == "MASTER_CORE" {
			sheetPQ = core.TenSheetPhanQuyenMaster
		} else if appMode == "TENANT_ADMIN" {
			sheetPQ = core.TenSheetPhanQuyenAdmin
		}

		// [LOCK CHUẨN]: Khóa riêng đúng Sheet Phân Quyền
		lockPQ := core.GetSheetLock(shopID, sheetPQ)
		lockPQ.RLock()
		defer lockPQ.RUnlock()

		hasPermission := false
		if pq, ok := core.CacheMapPhanQuyen[core.TaoCompositeKey(shopID, role)]; ok {
			// Thuật toán quét mảng quyền hạn Dot Notation
			for _, q := range pq.QuyenHan {
				if q == maChucNang {
					hasPermission = true
					break
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
