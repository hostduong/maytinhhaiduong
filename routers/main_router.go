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

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./static")

	// [QUAN TRỌNG]: Lễ tân Đứng ở cửa lớn, soi Domain và cấp SHOP_ID cho MỌI request
	router.Use(middlewares.IdentifyTenant())

	// =======================================================
	// KHU VỰC PUBLIC (Trang chủ & View SP - Không cần soát vé)
	// =======================================================
	router.GET("/", hien_thi_web.TrangChu)
	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham)

	// =======================================================
	// KHU VỰC AUTH (Đăng nhập, Đăng ký - Không cần soát vé)
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
	// KHU VỰC WORKSPACE (Bảo vệ đứng ở đây soát vé)
	// =======================================================
	workspace := router.Group("/master")
	workspace.Use(middlewares.CheckAuth())
	{
		workspace.GET("/tong-quan", tong_quan.TrangTongQuanMaster)
		workspace.GET("/dong-bo-sheets", dong_bo_sheets.TrangDongBoSheetsMaster)
		workspace.GET("/ho-so", ho_so.TrangHoSoMaster)
		workspace.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		workspace.GET("/quan-ly-may-tinh", san_pham.TrangQuanLyMayTinhMaster)
		workspace.GET("/thanh-vien", thanh_vien.TrangQuanLyThanhVienMaster)
		workspace.GET("/tin-nhan", tin_nhan.TrangTinNhanMaster)

		cauHinhUI := workspace.Group("/cau-hinh")
		cauHinhUI.Use(middlewares.RequireLevel(2))
		cauHinhUI.GET("", cau_hinh.TrangCaiDatCauHinhMaster)

		api := workspace.Group("/api")
		{
			api.POST("/cai-dat-cau-hinh/nha-cung-cap/save", cau_hinh.API_LuuNhaCungCap)
			api.POST("/dong-bo-sheets", dong_bo_sheets.API_NapLaiDuLieuMasterCoPIN)
			api.POST("/ho-so", ho_so.API_LuuHoSoMaster)
			api.POST("/change-pass", ho_so.API_DoiMatKhauMaster)
			api.POST("/change-pin", ho_so.API_DoiMaPinMaster)
			api.POST("/may-tinh/save", san_pham.API_LuuMayTinhMaster)
			api.POST("/thanh-vien/save", thanh_vien.API_LuuThanhVienMaster)
			api.POST("/thanh-vien/send-msg", thanh_vien.API_GuiTinNhanMaster)
			api.POST("/doc-tin-nhan", tin_nhan.API_DanhDauDaDocMaster)
			api.POST("/tin-nhan/send-chat", tin_nhan.API_GuiTinNhanChat)
		}
	}

	return router
}
