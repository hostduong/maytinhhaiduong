package routers

import (
	"app/middlewares" // Sử dụng thật sự màng lọc bảo mật
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

	// =======================================================
	// KHU VỰC PUBLIC (Mặt tiền: Khách hàng truy cập tự do)
	// =======================================================
	router.GET("/", hien_thi_web.TrangChu)
	router.GET("/login", hien_thi_web.TrangDangNhap)
	router.GET("/register", hien_thi_web.TrangDangKy)
	router.GET("/forgot-password", hien_thi_web.TrangQuenMatKhau)
	router.GET("/verify", hien_thi_web.TrangXacThucOTP)
	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham)

	// =======================================================
	// KHU VỰC WORKSPACE (Bên trong ứng dụng - BẢO MẬT 5 LỚP)
	// =======================================================
	workspace := router.Group("/master")
	
	// SỬ DỤNG AUTH THẬT: Kiểm tra Cookie -> RAM Cache (Sheet KHACH_HANG)
	workspace.Use(middlewares.CheckAuth()) 
	{
		// 1. Nhóm Render Giao diện HTML
		workspace.GET("/tong-quan", tong_quan.TrangTongQuanMaster)
		workspace.GET("/cau-hinh", cau_hinh.TrangCaiDatCauHinhMaster)
		workspace.GET("/dong-bo-sheets", dong_bo_sheets.TrangDongBoSheetsMaster)
		workspace.GET("/ho-so", ho_so.TrangHoSoMaster)
		workspace.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		workspace.GET("/quan-ly-may-tinh", san_pham.TrangQuanLyMayTinhMaster)
		workspace.GET("/thanh-vien", thanh_vien.TrangQuanLyThanhVienMaster)
		workspace.GET("/tin-nhan", tin_nhan.TrangTinNhanMaster)

		// 2. Nhóm Nhận Request AJAX (APIs)
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
