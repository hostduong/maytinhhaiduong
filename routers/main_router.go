package routers

import (
	"app/middlewares"
	"app/modules/cau_hinh_he_thong"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Phục vụ tài nguyên tĩnh
	router.Static("/static", "./static")

	// KHU VỰC WORKSPACE (Gộp chung Admin & Master)
	workspace := router.Group("/master")
	workspace.Use(middlewares.CheckAuth()) // [BẢO MẬT LỚP 1 & 2]
	{
		// =======================================================
		// MODULE 1: CẤU HÌNH HỆ THỐNG
		// =======================================================
		cauHinh := workspace.Group("/cau-hinh-he-thong")
		cauHinh.Use(middlewares.RequireLevel(2)) // [BẢO MẬT LỚP 5]
		{
			// Render View (Sẽ tạo ở bước sau)
			cauHinh.GET("/", cau_hinh_he_thong.TrangCauHinhHeThongView)

			apiCauHinh := cauHinh.Group("/api")
			apiCauHinh.Use(middlewares.CheckSaaSLimit("cau_hinh")) // [BẢO MẬT LỚP 3]

			// API xử lý dữ liệu
			apiCauHinh.POST("/nha-cung-cap/save", 
				middlewares.RequirePermission("system.setting.edit"), // [BẢO MẬT LỚP 4]
				cau_hinh_he_thong.API_LuuNhaCungCap,
			)
		}
	}

	return router
}
