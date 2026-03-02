package middlewares

import (
	"strings"

	"app/config"
	"github.com/gin-gonic/gin"
)

// IdentifyTenant: Trạm kiểm soát phân luồng SaaS dựa vào tên miền (Domain)
func IdentifyTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Đọc tên miền từ Request (VD: abc.99k.vn, maytinhhaiduong.run.app)
		host := c.Request.Host
		domain := strings.Split(host, ":")[0]

		// Giá trị dự phòng (Fallback an toàn cho môi trường Dev/Cloud Run)
		shopID := config.BienCauHinh.IdFileSheet 
		theme := "theme_master"

		// Ánh xạ Domain tĩnh từ Config (Cho Localhost)
		if mappedID, ok := config.MapDomainShop[domain]; ok {
			shopID = mappedID
		} else if strings.HasSuffix(domain, ".99k.vn") && domain != "www.99k.vn" && domain != "99k.vn" {
			// Tương lai: Logic tra cứu Domain động trong Cache Tenant
			// subdomain := strings.TrimSuffix(domain, ".99k.vn")
			// shopID = CacheTenant[subdomain].ShopID
			// theme = CacheTenant[subdomain].Theme
		}

		// Đóng dấu thân phận Tenant vào Request
		c.Set("SHOP_ID", shopID)
		c.Set("THEME", theme)
		
		c.Next()
	}
}
