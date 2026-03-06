package routers

import (
	"net/http"

	"app/middlewares"
	"app/modules/auth"
	"app/modules/bang_gia"
	"app/modules/cau_hinh"
	"app/modules/database_admin"
	"app/modules/dong_bo_sheets"
	"app/modules/goi_dich_vu"
	"app/modules/hien_thi_web"
	"app/modules/ho_so"
	"app/modules/nhap_hang"
	"app/modules/san_pham"
	"app/modules/thanh_toan"
	"app/modules/thanh_vien"
	"app/modules/tin_nhan"
	"app/modules/tong_quan_admin"
    "app/modules/tong_quan_master"
	"app/modules/trang_chu_admin"

	"github.com/gin-gonic/gin"
)

// =======================================================
// MIDDLEWARE: KHÓA VŨ TRỤ (CHỐNG LEO THANG ĐẶC QUYỀN)
// =======================================================
func RequireAppMode(allowedMode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentMode := c.GetString("APP_MODE")
		if currentMode != allowedMode {
			// Báo 404 để giấu luôn sự tồn tại của URL này với người ngoài
			c.HTML(http.StatusNotFound, "404", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")

	// Trạm IdentifyTenant: Đóng dấu APP_MODE từ Domain (Do sếp đã viết trước đó)
	router.Use(middlewares.IdentifyTenant())

	// =======================================================
	// 1. NGÃ BA ĐƯỜNG: ĐIỀU PHỐI TRANG CHỦ "/" DỰA VÀO TÊN MIỀN
	// =======================================================
	router.GET("/", func(c *gin.Context) {
		mode := c.GetString("APP_MODE")
		if mode == "MASTER_CORE" {
			// Sếp gõ sss.99k.vn -> Vào thẳng Tổng quan Master
			c.Redirect(http.StatusFound, "/master/tong-quan") 
		} else if mode == "TENANT_ADMIN" {
			// [ĐÃ SỬA] Chủ shop gõ admin.99k.vn -> Vào thẳng Tổng quan Shop
			c.Redirect(http.StatusFound, "/tong-quan") 
		} else {
			hien_thi_web.TrangChu(c) // www.99k.vn 
		}
	})

	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham)

	// =======================================================
	// 2. KHU VỰC ĐỊNH DANH (AUTH) - Hoạt động tự do trên các Domain
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
		apiAuth.POST("/send-otp", auth.API_SendOtp)
		apiAuth.POST("/reset-by-pin", auth.API_ResetByPin)
		apiAuth.POST("/reset-by-otp", auth.API_ResetByOtp)
	}

	// =======================================================
	// 3. VŨ TRỤ TENANT ADMIN (admin.99k.vn) - Dành cho Khách Hàng
	// =======================================================
	tenantAdmin := router.Group("")
	tenantAdmin.Use(RequireAppMode("TENANT_ADMIN")) // KHÓA CHẶT TÊN MIỀN ADMIN
	tenantAdmin.Use(middlewares.CheckAuth())
	{
		// Mới chỉ cấp quyền cho 2 chức năng đã hoàn thiện
		tenantAdmin.GET("/tong-quan", tong_quan_admin.TrangTongQuanAdmin)
		tenantAdmin.GET("/database", database_admin.TrangThietLapDatabaseAdmin)
		tenantAdmin.POST("/api/database/setup", database_admin.API_ThietLapDatabase)
		tenantAdmin.GET("/bang-gia", bang_gia_admin.TrangBangGiaAdmin)

		// Kế toán trưởng: Xử lý Thanh toán
		apiThanhToan := tenantAdmin.Group("/api/thanh-toan")
		{
			apiThanhToan.POST("/check-price", thanh_toan.API_CheckPrice)
			apiThanhToan.POST("/mua-goi", thanh_toan.API_MuaGoi)
		}
	}

	// =======================================================
	// 4. VŨ TRỤ GOD MODE (sss.99k.vn) - Full tính năng cho Sếp Test
	// =======================================================
	workspace := router.Group("/master")
	workspace.Use(RequireAppMode("MASTER_CORE")) // KHÓA CHẶT TÊN MIỀN SẾP
	workspace.Use(middlewares.CheckAuth())
	workspace.Use(middlewares.RequireLevel(2))
	{
		// Giao diện Workspace (Giữ lại tất cả để sếp test)
		workspace.GET("/tong-quan", tong_quan_master.TrangTongQuanMaster)
		workspace.GET("/goi-dich-vu", goi_dich_vu.TrangGoiDichVuMaster)
		workspace.GET("/ho-so", ho_so.TrangHoSoMaster)
		workspace.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		workspace.GET("/quan-ly-may-tinh", san_pham.TrangQuanLyMayTinhMaster)
		workspace.GET("/tin-nhan", tin_nhan.TrangTinNhanMaster)

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

			api.POST("/may-tinh/save", middlewares.RequirePermission("product.edit"), san_pham.API_LuuMayTinhMaster)
			api.POST("/cai-dat-cau-hinh/nha-cung-cap/save", middlewares.RequirePermission("config.edit"), cau_hinh.API_LuuNhaCungCap)
			api.POST("/nhap-hang/save", middlewares.RequirePermission("stock.import"), nhap_hang.API_LuuPhieuNhap)
			api.POST("/nhap-hang/status", middlewares.RequirePermission("stock.import"), nhap_hang.API_DoiTrangThaiPhieu)

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

	// =======================================================
	// 5. GIỮ LẠI VIEW BẢNG GIÁ DỰ PHÒNG CHO WWW.99K.VN (Nếu cần)
	// =======================================================
	portal := router.Group("/bang-gia")
	portal.Use(middlewares.CheckAuth())
	{
		portal.GET("/", bang_gia.TrangCongPortalBangGia)
	}

	return router
}
