package routers

import (
	"app/middlewares"
	"app/modules/auth"
	"app/modules/cau_hinh"
	"app/modules/dong_bo_sheets"
	"app/modules/hien_thi_web" // Đã đổi từ home sang hien_thi_web
	"app/modules/ho_so"
	"app/modules/nhap_hang"
	"app/modules/san_pham"
	"app/modules/thanh_vien"
	"app/modules/tin_nhan"
	"app/modules/tong_quan"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")

	// Lễ tân phân luồng Domain
	router.Use(middlewares.IdentifyTenant())

	// =======================================================
	// KHU VỰC PUBLIC
	// =======================================================
	router.GET("/", hien_thi_web.TrangChu) // Đã cập nhật
	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham) // Đã cập nhật

	// =======================================================
	// KHU VỰC AUTH
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
		apiAuth.POST("/verify-register", auth.API_Verify)
		apiAuth.POST("/send-otp", auth.API_SendOtp)
		apiAuth.POST("/reset-by-pin", auth.API_ResetByPin)
		apiAuth.POST("/reset-by-otp", auth.API_ResetByOtp)
	}

	// =======================================================
	// KHU VỰC WORKSPACE (Bảo vệ 5 lớp)
	// =======================================================
	workspace := router.Group("/master")
	workspace.Use(middlewares.CheckAuth())
	{
		// 1. Nhóm Load Giao diện HTML
		workspace.GET("/tong-quan", tong_quan.TrangTongQuanMaster)
		workspace.GET("/ho-so", ho_so.TrangHoSoMaster)
		workspace.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		workspace.GET("/quan-ly-may-tinh", san_pham.TrangQuanLyMayTinhMaster)
		workspace.GET("/tin-nhan", tin_nhan.TrangTinNhanMaster)

		// Giao diện Yêu cầu Cấp bậc Quản trị (Level 0, 1, 2)
		cauHinhUI := workspace.Group("/cau-hinh")
		cauHinhUI.Use(middlewares.RequireLevel(2))
		cauHinhUI.GET("/", cau_hinh.TrangCauHinhView)

		thanhVienUI := workspace.Group("/thanh-vien")
		thanhVienUI.Use(middlewares.RequireLevel(2))
		thanhVienUI.GET("/", thanh_vien.TrangQuanLyThanhVienMaster)

		dongBoUI := workspace.Group("/dong-bo-sheets")
		dongBoUI.Use(middlewares.RequireLevel(2))
		dongBoUI.GET("/", dong_bo_sheets.TrangDongBoSheetsMaster)

		// 2. Nhóm Gọi API Xử lý Dữ liệu
		api := workspace.Group("/api")
		{
			// API Không giới hạn quyền (Cá nhân tự đổi)
			api.POST("/ho-so", ho_so.API_LuuHoSoMaster)
			api.POST("/change-pass", ho_so.API_DoiMatKhauMaster)
			api.POST("/change-pin", ho_so.API_DoiMaPinMaster)
			api.POST("/doc-tin-nhan", tin_nhan.API_DanhDauDaDocMaster)
			api.POST("/tin-nhan/send-chat", tin_nhan.API_GuiTinNhanChat)

			// [LƯỚI THÉP RBAC CẨM NANG]: Kiểm tra quyền chuẩn xác cho từng nhóm hành động
			api.POST("/may-tinh/save", middlewares.RequirePermission("product.edit"), san_pham.API_LuuMayTinhMaster)
			api.POST("/cai-dat-cau-hinh/nha-cung-cap/save", middlewares.RequirePermission("config.edit"), cau_hinh.API_LuuNhaCungCap)

			// Các API Yêu cầu Cấp bậc Hệ thống
			apiAdmin := api.Group("")
			apiAdmin.Use(middlewares.RequireLevel(2))
			{
				apiAdmin.POST("/dong-bo-sheets", dong_bo_sheets.API_NapLaiDuLieuMasterCoPIN)
				apiAdmin.POST("/thanh-vien/save", thanh_vien.API_LuuThanhVienMaster)
				apiAdmin.POST("/thanh-vien/send-msg", thanh_vien.API_GuiTinNhanMaster)
			}
		}
	}

	return router
}
