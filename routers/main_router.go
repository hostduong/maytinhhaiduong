package routers

import (
	"net/http"
	"strings"

	"app/middlewares"
	
	// IMPORT HỆ THỐNG AUTH MỚI TỪ THƯ MỤC /auth/
	"app/auth/auth_admin"
	"app/auth/auth_master"
	"app/auth/auth_store"
	"app/auth/auth_verify"

	// CÁC MODULE NGHIỆP VỤ BÌNH THƯỜNG
	"app/modules/bang_gia"
	"app/modules/bang_gia_admin"
	"app/modules/cau_hinh"
	"app/modules/database_admin"
	"app/modules/dong_bo_sheets"
	"app/modules/goi_dich_vu_master"
	"app/modules/hien_thi_web"
	"app/modules/ho_so"
	"app/modules/nhap_hang"
	"app/modules/thanh_toan"
	"app/modules/thanh_vien_master" 
	"app/modules/tin_nhan_master"
	"app/modules/tong_quan_admin"
	"app/modules/tong_quan_master"
	"app/modules/cua_hang_master"
	"app/modules/product_master"
	"app/modules/setup"
	"app/modules/phan_quyen_master"

	"github.com/gin-gonic/gin"
)

// =======================================================
// MIDDLEWARE: KHÓA VŨ TRỤ (CHỐNG LEO THANG ĐẶC QUYỀN)
// =======================================================
func RequireAppMode(allowedMode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("APP_MODE") != allowedMode {
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
	// TRẠM KHỞI THỦY HỆ THỐNG (CHỈ DÀNH CHO SSS.99K.VN)
	// =======================================================
	router.GET("/setup", func(c *gin.Context) {
		if c.GetString("APP_MODE") == "MASTER_CORE" {
			setup.TrangSetup(c)
		} else {
			c.HTML(http.StatusNotFound, "404", nil) 
		}
	})
	router.POST("/setup", func(c *gin.Context) {
		if c.GetString("APP_MODE") == "MASTER_CORE" {
			setup.API_Setup(c)
		} else {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Truy cập bị từ chối!"})
		}
	})

	// =======================================================
	// 1. VÙNG MẶC ĐỊNH (PUBLIC)
	// =======================================================
	router.GET("/", func(c *gin.Context) {
		mode := c.GetString("APP_MODE")
		host := c.Request.Host
		
		if mode == "MASTER_CORE" {
			c.Redirect(http.StatusFound, "/master/tong-quan") 
		} else if mode == "TENANT_ADMIN" {
			if strings.HasPrefix(host, "shop.") {
				c.Redirect(http.StatusFound, "/tong-quan") 
			} else {
				hien_thi_web.TrangChu(c) // www.99k.vn
			}
		} else {
			hien_thi_web.TrangChu(c) // [cuahang].99k.vn
		}
	})

	router.GET("/san-pham/:id", hien_thi_web.ChiTietSanPham)

	// =======================================================
	// BỘ ĐIỀU PHỐI XÁC THỰC ĐA VŨ TRỤ (MULTIVERSE SWITCHER)
	// =======================================================
	
	// ĐĂNG NHẬP (Render View & Xử lý POST)
	router.GET("/login", func(c *gin.Context) {
		switch c.GetString("APP_MODE") {
		case "MASTER_CORE": auth_master.TrangDangNhap(c)
		case "TENANT_ADMIN": auth_admin.TrangDangNhap(c)
		default: auth_store.TrangDangNhap(c)
		}
	})
	router.POST("/login", func(c *gin.Context) {
		switch c.GetString("APP_MODE") {
		case "MASTER_CORE": auth_master.API_Login(c)
		case "TENANT_ADMIN": auth_admin.API_Login(c)
		default: auth_store.API_Login(c)
		}
	})

	// ĐĂNG KÝ (Master không có tính năng Đăng ký)
	router.GET("/register", func(c *gin.Context) {
		switch c.GetString("APP_MODE") {
		case "MASTER_CORE": c.HTML(http.StatusNotFound, "404", nil)
		case "TENANT_ADMIN": auth_admin.TrangDangKy(c)
		default: auth_store.TrangDangKy(c)
		}
	})
	router.POST("/register", func(c *gin.Context) {
		switch c.GetString("APP_MODE") {
		case "MASTER_CORE": c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Master không hỗ trợ đăng ký từ bên ngoài!"})
		case "TENANT_ADMIN": auth_admin.API_Register(c)
		default: auth_store.API_Register(c)
		}
	})

	// QUÊN MẬT KHẨU
	router.GET("/forgot-password", func(c *gin.Context) {
		switch c.GetString("APP_MODE") {
		case "MASTER_CORE": auth_master.TrangQuenMatKhau(c)
		case "TENANT_ADMIN": auth_admin.TrangQuenMatKhau(c)
		default: auth_store.TrangQuenMatKhau(c)
		}
	})

	// ĐĂNG XUẤT
	router.GET("/logout", func(c *gin.Context) {
		switch c.GetString("APP_MODE") {
		case "MASTER_CORE": auth_master.API_Logout(c)
		case "TENANT_ADMIN": auth_admin.API_Logout(c)
		default: auth_store.API_Logout(c)
		}
	})

	// API BẢO MẬT & XÁC THỰC (OTP/PIN)
	apiAuth := router.Group("/api/auth")
	{
		// Cụm gọi trực tiếp về Trạm kiểm soát chung (Verify)
		apiAuth.POST("/send-otp", auth_verify.API_SendOtp)
		apiAuth.POST("/check-otp", auth_verify.API_CheckOtp) 
		apiAuth.POST("/check-pin", auth_verify.API_CheckPin) 

		// Cụm đổi mật khẩu: Sau khi Verify xong, phân luồng về Module để đổi pass
		apiAuth.POST("/reset-by-pin", func(c *gin.Context) {
			switch c.GetString("APP_MODE") {
			case "MASTER_CORE": auth_master.API_ResetByPin(c)
			case "TENANT_ADMIN": auth_admin.API_ResetByPin(c)
			default: auth_store.API_ResetByPin(c)
			}
		})
		apiAuth.POST("/reset-by-otp", func(c *gin.Context) {
			switch c.GetString("APP_MODE") {
			case "MASTER_CORE": auth_master.API_ResetByOtp(c)
			case "TENANT_ADMIN": auth_admin.API_ResetByOtp(c)
			default: auth_store.API_ResetByOtp(c)
			}
		})
	}

	// =======================================================
	// 2. VŨ TRỤ CHỦ SHOP (TENANT ADMIN)
	// =======================================================
	admin := router.Group("")
	admin.Use(RequireAppMode("TENANT_ADMIN")) 
	admin.Use(middlewares.CheckAuth())
	{
		admin.GET("/tong-quan", tong_quan_admin.TrangTongQuanAdmin)
		admin.GET("/database", database_admin.TrangThietLapDatabaseAdmin)
		admin.GET("/bang-gia", bang_gia_admin.TrangBangGiaAdmin)
		
		apiAdmin := admin.Group("/api")
		{
			apiAdmin.POST("/database/setup", database_admin.API_ThietLapDatabase)
			apiAdmin.POST("/thanh-toan/check-price", thanh_toan.API_CheckPrice)
			apiAdmin.POST("/thanh-toan/mua-goi", thanh_toan.API_MuaGoi)
		}
	}

	// =======================================================
	// 3. VŨ TRỤ SẾP (MASTER CORE)
	// =======================================================
	master := router.Group("/master")
	master.Use(RequireAppMode("MASTER_CORE")) 
	master.Use(middlewares.CheckAuth())
	master.Use(middlewares.RequireLevel(2))
	{
		master.GET("/tong-quan", tong_quan_master.TrangTongQuanMaster)
		master.GET("/goi-dich-vu", goi_dich_vu_master.TrangGoiDichVuMaster)
		master.GET("/ho-so", ho_so.TrangHoSoMaster)
		master.GET("/nhap-hang", nhap_hang.TrangNhapHangMaster)
		master.GET("/quan-ly-san-pham", product_master.TrangQuanLySanPhamMaster)
		master.GET("/tin-nhan", tin_nhan_master.TrangTinNhanMaster)
		master.GET("/quan-ly-cua-hang", cua_hang_master.TrangQuanLyCuaHangMaster)
		master.GET("/cau-hinh", cau_hinh.TrangCaiDatCauHinhMaster)
		master.GET("/thanh-vien", thanh_vien_master.TrangQuanLyThanhVienMaster) 
		master.GET("/dong-bo-sheets", dong_bo_sheets.TrangDongBoSheetsMaster)
		master.GET("/phan-quyen", phan_quyen_master.TrangPhanQuyenMaster)

		apiMaster := master.Group("/api")
		{
			apiMaster.POST("/ho-so", ho_so.API_LuuHoSoMaster)
			apiMaster.POST("/change-pass", ho_so.API_DoiMatKhauMaster)
			apiMaster.POST("/change-pin", ho_so.API_DoiMaPinMaster)
			apiMaster.POST("/doc-tin-nhan", tin_nhan_master.API_DanhDauDaDocMaster)
            apiMaster.POST("/tin-nhan/send-chat", tin_nhan_master.API_GuiTinNhanChat)

			apiMaster.POST("/product/save", middlewares.RequirePermission("product.edit"), product_master.API_LuuSanPhamMaster)
			apiMaster.POST("/cai-dat-cau-hinh/nha-cung-cap/save", middlewares.RequirePermission("config.edit"), cau_hinh.API_LuuNhaCungCap)
			apiMaster.POST("/nhap-hang/save", middlewares.RequirePermission("stock.import"), nhap_hang.API_LuuPhieuNhap)
			apiMaster.POST("/nhap-hang/status", middlewares.RequirePermission("stock.import"), nhap_hang.API_DoiTrangThaiPhieu)

			apiMaster.POST("/phan-quyen/save", phan_quyen_master.API_LuuPhanQuyenMaster)
			apiMaster.POST("/dong-bo-sheets", dong_bo_sheets.API_NapLaiDuLieuMasterCoPIN)
			apiMaster.POST("/thanh-vien/save", thanh_vien_master.API_LuuThanhVienMaster) 
			apiMaster.POST("/thanh-vien/send-msg", thanh_vien_master.API_GuiTinNhanMaster) 
			apiMaster.POST("/goi-dich-vu/save", goi_dich_vu_master.API_LuuGoiDichVuMaster)
			apiMaster.POST("/cua-hang/save", cua_hang_master.API_LuuCuaHangMaster)
		}
	}

	// =======================================================
	// 4. GIỮ LẠI VIEW BẢNG GIÁ DỰ PHÒNG CHO WWW.99K.VN 
	// =======================================================
	portal := router.Group("/bang-gia")
	portal.Use(middlewares.CheckAuth())
	{
		portal.GET("/", bang_gia.TrangCongPortalBangGia)
	}

	return router
}
