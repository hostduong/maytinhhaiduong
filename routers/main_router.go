package routers

import (
	"app/middlewares"
	"app/modules/auth"
	"app/modules/cau_hinh"
	"app/modules/dong_bo_sheets"
	"app/modules/hien_thi_web" 
	"app/modules/ho_so"
	"app/modules/nhap_hang"
	"app/modules/san_pham"
	"app/modules/thanh_vien"
	"app/modules/tin_nhan"
	"app/modules/tong_quan"
	"app/modules/goi_dich_vu"
	"app/modules/bang_gia"
	"app/modules/database_admin"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")

	// Trạm IdentifyTenant: Nhận diện Shop/Master từ Domain
	router.Use(middlewares.IdentifyTenant())

	// =======================================================
	// 1. KHU VỰC PUBLIC (Không cần đăng nhập)
	// =======================================================
	router.GET("/", hien_thi_web.TrangChu) 
	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham) 

	// =======================================================
	// 2. KHU VỰC ĐỊNH DANH (AUTH)
	// =======================================================
	router.GET("/login", auth.TrangDangNhap)
	router.GET("/register", auth.TrangDangKy)
	router.GET("/forgot-password", auth.TrangQuenMatKhau)
	router.GET("/verify", auth.TrangXacThucOTP)
	router.GET("/logout", auth.API_Logout)

	router.POST("/login", auth.API_Login)
	router.POST("/register", auth.API_Register)

	apiAuth := router.Group("/api/auth")
	{
		// Đã xóa API_Verify (Xác thực đăng ký) vì luồng đăng ký mới cho qua thẳng
		apiAuth.POST("/send-otp", auth.API_SendOtp)
		apiAuth.POST("/reset-by-pin", auth.API_ResetByPin)
		apiAuth.POST("/reset-by-otp", auth.API_ResetByOtp)
	}

	// =======================================================
	// 3. KHU VỰC CHỌN GÓI (BẢNG GIÁ) - Bắt buộc Login
	// =======================================================
	portal := router.Group("/bang-gia")
	portal.Use(middlewares.CheckAuth())
	{
		portal.GET("/", bang_gia.TrangCongPortalBangGia)
		portal.POST("/api/check-price", bang_gia.API_CheckGia)
		portal.POST("/api/mua-goi", bang_gia.API_MuaGoi)
	}

	// =======================================================
	// KHU VỰC THIẾT LẬP DATABASE CHUNG
	// =======================================================
	admin := router.Group("/admin")
	admin.Use(middlewares.CheckAuth())
	admin.Use(middlewares.EnforceDomainBoundary())
	{
		// Cổng vào thiết lập Database sau khi mua gói
		admin.GET("/database", database_admin.TrangThietLapDatabaseAdmin) 
		admin.POST("/api/database/setup", database_admin.API_ThietLapDatabase)
	}

	// =======================================================
	// 5. KHU VỰC QUẢN TRỊ HỆ THỐNG (MASTER)
	// =======================================================
	workspace := router.Group("/master")
	workspace.Use(middlewares.CheckAuth())
	workspace.Use(middlewares.RequireLevel(2))
	{
		// Giao diện Workspace
		workspace.GET("/tong-quan", tong_quan.TrangTongQuanMaster)
		workspace.GET("/goi-dich-vu", goi_dich_vu.TrangGoiDichVuMaster)
		workspace.GET("/ho-so", ho_so.TrangHoSoMaster)
		workspace.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		workspace.GET("/quan-ly-may-tinh", san_pham.TrangQuanLyMayTinhMaster)
		workspace.GET("/tin-nhan", tin_nhan.TrangTinNhanMaster)

		// Cấu hình nâng cao (Yêu cầu Level 2)
		cauHinhUI := workspace.Group("/cau-hinh")
		cauHinhUI.Use(middlewares.RequireLevel(2))
		cauHinhUI.GET("/", cau_hinh.TrangCaiDatCauHinhMaster) 

		thanhVienUI := workspace.Group("/thanh-vien")
		thanhVienUI.Use(middlewares.RequireLevel(2))
		thanhVienUI.GET("/", thanh_vien.TrangQuanLyThanhVienMaster)

		dongBoUI := workspace.Group("/dong-bo-sheets")
		dongBoUI.Use(middlewares.RequireLevel(2))
		dongBoUI.GET("/", dong_bo_sheets.TrangDongBoSheetsMaster)

		// API Xử lý Master
		api := workspace.Group("/api")
		{
			api.POST("/ho-so", ho_so.API_LuuHoSoMaster)
			api.POST("/change-pass", ho_so.API_DoiMatKhauMaster)
			api.POST("/change-pin", ho_so.API_DoiMaPinMaster)
			api.POST("/doc-tin-nhan", tin_nhan.API_DanhDauDaDocMaster)
			api.POST("/tin-nhan/send-chat", tin_nhan.API_GuiTinNhanChat)

			// RBAC Product & Stock
			api.POST("/may-tinh/save", middlewares.RequirePermission("product.edit"), san_pham.API_LuuMayTinhMaster)
			api.POST("/cai-dat-cau-hinh/nha-cung-cap/save", middlewares.RequirePermission("config.edit"), cau_hinh.API_LuuNhaCungCap)
			api.POST("/nhap-hang/save", middlewares.RequirePermission("stock.import"), nhap_hang.API_LuuPhieuNhap)
			api.POST("/nhap-hang/status", middlewares.RequirePermission("stock.import"), nhap_hang.API_DoiTrangThaiPhieu)

			// API Admin Master (Level 2)
			apiAdmin := api.Group("")
			apiAdmin.Use(middlewares.RequireLevel(2))
			{
				apiAdmin.POST("/dong-bo-sheets", dong_bo_sheets.API_NapLaiDuLieuMasterCoPIN)
				apiAdmin.POST("/thanh-vien/save", thanh_vien.API_LuuThanhVienMaster)
				apiAdmin.POST("/thanh-vien/send-msg", thanh_vien.API_GuiTinNhanMaster)
				apiAdmin.POST("/goi-dich-vu/save", goi_dich_vu.API_LuuGoiDichVu)
			}
		}
	}

	return router
}
