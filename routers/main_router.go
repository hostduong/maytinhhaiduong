package routers

import (
	"app/middlewares"
	"app/modules/cau_hinh_he_thong"

	"github.com/gin-gonic/gin"
)

// SetupRouter: Hàm khởi tạo toàn bộ mạng lưới đường dẫn (URL) của hệ thống
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 1. Mở cửa cho phép tải file tĩnh (CSS, JS, Hình ảnh)
	router.Static("/static", "./static")

	// ==============================================================================
	// KHU VỰC PUBLIC (Trang chủ, Đăng nhập - Không cần bảo vệ)
	// ==============================================================================
	// public := router.Group("/")
	// {
	// 	public.GET("/login", auth_module.TrangDangNhap)
	// }

	// ==============================================================================
	// KHU VỰC WORKSPACE (Bên trong ứng dụng - Gộp chung Admin & Master)
	// ==============================================================================
	workspace := router.Group("/master")

	// [MÀNG LỌC 1 & 2]: Chặn toàn bộ người lạ. Chỉ User hợp lệ & Đang hoạt động mới lọt qua.
	workspace.Use(middlewares.CheckAuth())
	{
		// --------------------------------------------------------------------------
		// MODULE 1: CẤU HÌNH HỆ THỐNG
		// --------------------------------------------------------------------------
		cauHinh := workspace.Group("/cau-hinh-he-thong")

		// [MÀNG LỌC 5]: Tường lửa cấp bậc. Chỉ có Level 0, 1, 2 mới được vào phân khu này.
		cauHinh.Use(middlewares.RequireLevel(2))
		{
			// 1. API Render Giao diện Web HTML
			cauHinh.GET("/", cau_hinh_he_thong.TrangCauHinhHeThongView)

			// 2. Nhóm API Xử lý dữ liệu (Cần bảo vệ khắt khe hơn)
			apiCauHinh := cauHinh.Group("/api")

			// [MÀNG LỌC 3]: Chặn theo dung lượng Gói cước SaaS
			apiCauHinh.Use(middlewares.CheckSaaSLimit("cau_hinh"))

			// [MÀNG LỌC 4]: Chặn theo Mã chức năng (RBAC) trên từng Nút bấm cụ thể
			apiCauHinh.POST("/nha-cung-cap/save", 
				middlewares.RequirePermission("system.setting.edit"), // Mã quyền dựa theo PDF
				cau_hinh_he_thong.API_LuuNhaCungCap,
			)

			// Khai báo sẵn chỗ cho các API sau này
			// apiCauHinh.POST("/danh-muc/save", middlewares.RequirePermission("system.setting.edit"), cau_hinh_he_thong.API_LuuDanhMuc)
		}

		// --------------------------------------------------------------------------
		// MODULE 2: QUẢN LÝ SẢN PHẨM (Ví dụ)
		// --------------------------------------------------------------------------
		// sanPham := workspace.Group("/san-pham")
		// {
		//     sanPham.GET("/", quan_ly_san_pham.TrangQuanLySanPhamView)
		//     // ...
		// }
	}

	return router
}
