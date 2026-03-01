package middlewares

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// RequireLevel: Yêu cầu người gọi phải có cấp bậc <= minLevel
// VD: RequireLevel(2) -> Chỉ Level 0, 1, 2 mới được vào. Level 3 trở lên bị đá văng.
func RequireLevel(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userLevel := c.GetInt("USER_LEVEL")

		// Đặc quyền của Level 0 (ID 001) luôn luôn qua cửa
		if userLevel == 0 {
			c.Next()
			return
		}

		if userLevel > minLevel {
			TuChoiTruyCap(c, http.StatusForbidden, "Truy cập bị từ chối! Yêu cầu cấp bậc đặc quyền hệ thống.")
			return
		}

		c.Next()
	}
}
