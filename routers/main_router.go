package routers

import (
	"net/http"
	"app/modules/cau_hinh"
	"app/modules/tong_quan"
	"github.com/gin-gonic/gin"
)

// FakeAuth: Màng lọc giả lập đăng nhập để test Giao diện
// Nó tự động cấp thẻ Founder (Level 0) cho bạn
func FakeAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("SHOP_ID", "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8")
		c.Set("USER_ID", "0000000000000000001") // ID của Sáng lập viên
		c.Set("USER_ROLE", "quan_tri_he_thong")
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	
	// Phục vụ CSS, JS, Ảnh
	router.Static("/static", "./static")

	// Vào trang chủ sẽ tự động nảy sang Dashboard
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/master/tong-quan")
	})

	// Khu vực làm việc
	workspace := router.Group("/master")
	workspace.Use(FakeAuth()) // Quẹt thẻ VIP ở đây
	{
		workspace.GET("/tong-quan", tong_quan.TrangTongQuanMaster)
		workspace.GET("/cau-hinh", cau_hinh.TrangCaiDatCauHinhMaster)
	}

	return router
}
