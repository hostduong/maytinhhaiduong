package routers

import (
	"app/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRouter: Cấu hình toàn bộ tuyến đường của hệ thống
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 1. Phục vụ tài nguyên tĩnh (CSS, JS, Ảnh)
	router.Static("/static", "./static")

	// 2. KHU VỰC PUBLIC (Không cần đăng nhập)
	public := router.Group("/")
	{
		// Sẽ gọi đến module Auth ở Giai đoạn 4
		// public.GET("/login", auth_module.TrangDangNhapView)
		// public.POST("/api/login", auth_module.API_XuLyDangNhap)
	}

	// 3. KHU VỰC WORKSPACE (Gộp chung Admin & Master như đã chốt)
	// Bất cứ ai đăng nhập thành công đều vào /master. Giao diện & Data hiển thị sẽ tự co giãn theo Level.
	workspace := router.Group("/master")
	workspace.Use(middlewares.CheckAuth()) // [BẢO MẬT LỚP 1 & 2]: Cookie & Status
	{
		// =======================================================
		// MODULE 1: CẤU HÌNH HỆ THỐNG LÕI (Chỉ dành cho Master)
		// =======================================================
		cauHinh := workspace.Group("/cau-hinh-he-thong")
		cauHinh.Use(middlewares.RequireLevel(2)) // [BẢO MẬT LỚP 5]: Quét Cấp Bậc (Chỉ Level 0,1,2 lọt qua)
		{
			// Render HTML
			// cauHinh.GET("/", cau_hinh_he_thong.TrangCauHinhHeThongView)

			// Nhóm API xử lý (Gắn thêm check gói cước)
			apiCauHinh := cauHinh.Group("/api")
			apiCauHinh.Use(middlewares.CheckSaaSLimit("cau_hinh")) // [BẢO MẬT LỚP 3]: Gói cước

			// Khai báo chính xác ma_chuc_nang từ file PDF
			// apiCauHinh.POST("/danh-muc/save", middlewares.RequirePermission("system.setting.edit"), cau_hinh_he_thong.API_LuuDanhMuc)
			// apiCauHinh.POST("/nha-cung-cap/save", middlewares.RequirePermission("system.setting.edit"), cau_hinh_he_thong.API_LuuNhaCungCap)
		}

		// =======================================================
		// MODULE 2: QUẢN LÝ SẢN PHẨM (Dành cho toàn mạng lưới)
		// =======================================================
		sanPham := workspace.Group("/san-pham")
		// (Không chặn Level ở đây vì Chủ Shop (Lv3) hay Thủ Kho (Lv6) đều cần vào trang này)
		{
			// Render HTML
			// sanPham.GET("/", quan_ly_san_pham.TrangSanPhamView)

			// Gắn quyền chi tiết cho từng nút bấm
			// sanPham.POST("/api/create", middlewares.RequirePermission("product.create"), quan_ly_san_pham.API_TaoSanPham)
			// sanPham.POST("/api/edit", middlewares.RequirePermission("product.edit"), quan_ly_san_pham.API_SuaSanPham)
			// sanPham.POST("/api/delete", middlewares.RequirePermission("product.delete"), quan_ly_san_pham.API_XoaSanPham)
		}
	}

	return router
}
