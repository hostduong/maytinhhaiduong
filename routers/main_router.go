package routers

import (
	"net/http"
	"strings"

	"app/middlewares"
	"app/modules/auth"
	"app/modules/bang_gia"
	"app/modules/bang_gia_admin"
	"app/modules/cau_hinh"
	"app/modules/database_admin"
	"app/modules/dong_bo_sheets"
	"app/modules/goi_dich_vu_master"
	"app/modules/hien_thi_web"
	"app/modules/ho_so"
	"app/modules/nhap_hang"
	"app/modules/san_pham"
	"app/modules/thanh_toan"
	"app/modules/thanh_vien_master" // [S.O.P] Đã trỏ đúng thư mục Master
	"app/modules/tin_nhan_master"
	"app/modules/tong_quan_admin"
	"app/modules/tong_quan_master"
	"app/modules/cua_hang_master"

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

	// Trạm IdentifyTenant: Đóng dấu APP_MODE từ Domain
	router.Use(middlewares.IdentifyTenant())

	// =======================================================
	// 1. VÙNG MẶC ĐỊNH (PUBLIC): router.GET / router.POST
	// =======================================================
	router.GET("/", func(c *gin.Context) {
		mode := c.GetString("APP_MODE")
		host := c.Request.Host
		
		if mode == "MASTER_CORE" {
			c.Redirect(http.StatusFound, "/master/tong-quan") 
		} else if mode == "TENANT_ADMIN" {
			if strings.HasPrefix(host, "admin.") {
				c.Redirect(http.StatusFound, "/tong-quan") 
			} else {
				hien_thi_web.TrangChu(c) // www.99k.vn
			}
		} else {
			hien_thi_web.TrangChu(c) // [cuahang].99k.vn
		}
	})

	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham)

	// Auth (Xác thực)
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
	// 2. VŨ TRỤ CHỦ SHOP (TENANT ADMIN): admin.GET / apiAdmin.POST
	// =======================================================
	admin := router.Group("")
	admin.Use(RequireAppMode("TENANT_ADMIN")) 
	admin.Use(middlewares.CheckAuth())
	{
		// Render Giao Diện
		admin.GET("/tong-quan", tong_quan_admin.TrangTongQuanAdmin)
		admin.GET("/database", database_admin.TrangThietLapDatabaseAdmin)
		admin.GET("/bang-gia", bang_gia_admin.TrangBangGiaAdmin)
		
		// Render Giao Diện (Module dự kiến)
		// admin.GET("/nha-cung-cap", nha_cung_cap_admin.TrangNhaCungCap)

		// Xử lý Dữ liệu ngầm
		apiAdmin := admin.Group("/api")
		{
			apiAdmin.POST("/database/setup", database_admin.API_ThietLapDatabase)
			apiAdmin.POST("/thanh-toan/check-price", thanh_toan.API_CheckPrice)
			apiAdmin.POST("/thanh-toan/mua-goi", thanh_toan.API_MuaGoi)
			
			// API (Module dự kiến)
			// apiAdmin.POST("/nha-cung-cap/save", nha_cung_cap_admin.API_LuuNhaCungCap)
		}
	}

	// =======================================================
	// 3. VŨ TRỤ SẾP (MASTER CORE): master.GET / apiMaster.POST
	// =======================================================
	master := router.Group("/master")
	master.Use(RequireAppMode("MASTER_CORE")) 
	master.Use(middlewares.CheckAuth())
	master.Use(middlewares.RequireLevel(2))
	{
		// Render Giao Diện
		master.GET("/tong-quan", tong_quan_master.TrangTongQuanMaster)
		master.GET("/goi-dich-vu", goi_dich_vu_master.TrangGoiDichVuMaster)
		master.GET("/ho-so", ho_so.TrangHoSoMaster)
		master.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		master.GET("/quan-ly-may-tinh", san_pham.TrangQuanLyMayTinhMaster)
		master.GET("/tin-nhan", tin_nhan_master.TrangTinNhanMaster)
		master.GET("/quan-ly-cua-hang", cua_hang_master.TrangQuanLyCuaHangMaster)
		
		// (Các UI bọc thêm RequireLevel nếu cần, hiện tại master đã bọc sẵn Level 2)
		master.GET("/cau-hinh", cau_hinh.TrangCaiDatCauHinhMaster)
		master.GET("/thanh-vien", thanh_vien_master.TrangQuanLyThanhVienMaster) 
		master.GET("/dong-bo-sheets", dong_bo_sheets.TrangDongBoSheetsMaster)

		// Xử lý Dữ liệu ngầm
		apiMaster := master.Group("/api")
		{
			apiMaster.POST("/ho-so", ho_so.API_LuuHoSoMaster)
			apiMaster.POST("/change-pass", ho_so.API_DoiMatKhauMaster)
			apiMaster.POST("/change-pin", ho_so.API_DoiMaPinMaster)
			apiMaster.POST("/doc-tin-nhan", tin_nhan_master.API_DanhDauDaDocMaster)
            apiMaster.POST("/tin-nhan/send-chat", tin_nhan_master.API_GuiTinNhanChat)

			// Gắn Middleware kiểm tra Permission từng nút bấm
			apiMaster.POST("/may-tinh/save", middlewares.RequirePermission("product.edit"), san_pham.API_LuuMayTinhMaster)
			apiMaster.POST("/cai-dat-cau-hinh/nha-cung-cap/save", middlewares.RequirePermission("config.edit"), cau_hinh.API_LuuNhaCungCap)
			apiMaster.POST("/nhap-hang/save", middlewares.RequirePermission("stock.import"), nhap_hang.API_LuuPhieuNhap)
			apiMaster.POST("/nhap-hang/status", middlewares.RequirePermission("stock.import"), nhap_hang.API_DoiTrangThaiPhieu)

			apiMaster.POST("/dong-bo-sheets", dong_bo_sheets.API_NapLaiDuLieuMasterCoPIN)
			apiMaster.POST("/thanh-vien/save", thanh_vien_master.API_LuuThanhVienMaster) 
			apiMaster.POST("/thanh-vien/send-msg", thanh_vien_master.API_GuiTinNhanMaster) 
			apiMaster.POST("/goi-dich-vu/save", goi_dich_vu_master.API_LuuGoiDichVuMaster)
			apiMaster.POST("/cua-hang/save", cua_hang_master.API_LuuCuaHangMaster)
		}
	}

	// =======================================================
	// 4. GIỮ LẠI VIEW BẢNG GIÁ DỰ PHÒNG CHO WWW.99K.VN (Nếu cần)
	// =======================================================
	portal := router.Group("/bang-gia")
	portal.Use(middlewares.CheckAuth())
	{
		portal.GET("/", bang_gia.TrangCongPortalBangGia)
	}

	return router
}
