package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// =================================================================
// 1. TRẠM KIỂM SOÁT LÃNH THỔ (DOMAIN BOUNDARY ENFORCER)
// Chạy tự động để kiểm tra Cấp bậc (Level) so với Tên miền truy cập
// =================================================================
func EnforceDomainBoundary() gin.HandlerFunc {
	return func(c *gin.Context) {
		appMode := c.GetString("APP_MODE")
		userLevel := c.GetInt("USER_LEVEL") // Cấp độ từ 0 đến 9 

		// 1. Tầng Tuyệt Mật (sss.99k.vn) -> Chỉ cho phép Level 0, 1, 2 
		if appMode == "MASTER_CORE" && userLevel > 2 {
			TuChoiTruyCap(c, http.StatusForbidden, "Cảnh báo an ninh: Khu vực tuyệt mật, chỉ dành cho Ban Quản Trị Lõi!")
			return
		}

		// 2. Tầng Quản Trị Shop (admin.99k.vn) -> Chỉ cho phép Level <= 3 (Chủ shop và Sếp) 
		if appMode == "TENANT_ADMIN" && userLevel > 3 {
			TuChoiTruyCap(c, http.StatusForbidden, "Truy cập bị từ chối! Bạn không có quyền vào Không gian Quản trị của Cửa hàng.")
			return
		}

		// 3. Tầng Cửa hàng ([cuahang].99k.vn) -> Level từ 4-9 và Sếp đều được vào 
		// Trạm này mở cửa cho qua
		c.Next()
	}
}

// =================================================================
// 2. TRẠM YÊU CẦU CẤP BẬC CỤ THỂ (GRANULAR LEVEL CHECK)
// Dùng để bọc các API nhạy cảm bên trong (VD: API Xóa cửa hàng cần Level 0)
// =================================================================
func RequireLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel := c.GetInt("USER_LEVEL")

		// Đặc quyền tuyệt đối của Chúa tể (Level 0 / ID 001) luôn luôn qua cửa 
		if userLevel == 0 {
			c.Next()
			return
		}

		if userLevel > minLevel {
			TuChoiTruyCap(c, http.StatusForbidden, "Thao tác bị từ chối! Yêu cầu cấp bậc đặc quyền hệ thống cao hơn.")
			return
		}

		c.Next()
	}
}
