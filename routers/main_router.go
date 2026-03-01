package routers

import (
	"app/middlewares"
	"app/modules/cau_hinh" // Trỏ vào thư mục ngắn gọn
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")

	workspace := router.Group("/master")
	workspace.Use(middlewares.CheckAuth())
	{
		cauHinh := workspace.Group("/cau-hinh")
		cauHinh.Use(middlewares.RequireLevel(2))
		{
			cauHinh.GET("/", cau_hinh.TrangCauHinhView)

			apiCauHinh := cauHinh.Group("/api")
			apiCauHinh.Use(middlewares.CheckSaaSLimit("cau_hinh"))
			apiCauHinh.POST("/nha-cung-cap/save", 
				middlewares.RequirePermission("system.setting.edit"), 
				cau_hinh.API_LuuNhaCungCap,
			)
		}
	}
	return router
}
