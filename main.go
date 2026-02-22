package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/cau_hinh"
	"app/chuc_nang"
	"app/chuc_nang_admin"
	admin_may_tinh "app/chuc_nang_admin/may_tinh"
	"app/chuc_nang_master" // [MỚI] Import Controller quản lý Cửa hàng
	"app/core"

	"github.com/gin-gonic/gin"
)

// [SỬA QUAN TRỌNG] Nạp thêm thư mục giao_dien_master
//go:embed giao_dien/*.html giao_dien/*/*.html giao_dien_admin/*.html giao_dien_admin/*/*.html giao_dien_master/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG HỆ THỐNG SAAS...")

	cau_hinh.KhoiTaoCauHinh()
	core.KhoiTaoNenTang()
	core.KhoiTaoWorkerGhiSheet()

	log.Println("📦 [BOOT] Đang nạp dữ liệu Master Data...")
	core.NapPhanQuyen("")
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapBienLoiNhuan("")	
	core.NapKhachHang("")

	router := gin.Default()
	router.Use(chuc_nang.GatewaySaaS, chuc_nang.KiemTraGoiDichVu)

	funcMap := chuc_nang.LayBoHamHTML()

	// [SỬA QUAN TRỌNG] Parse thêm các file HTML của Cổng Master
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "giao_dien/*.html", "giao_dien/*/*.html", "giao_dien_admin/*.html", "giao_dien_admin/*/*.html", "giao_dien_master/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER ---
	// Public & Auth
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	router.GET("/logout", chuc_nang.DangXuat)
	router.GET("/forgot-password", chuc_nang.TrangQuenMatKhau) // <-- Đã thêm dấu đóng ngoặc
	router.GET("/xac-thuc", chuc_nang.TrangXacThuc)
	router.POST("/api/auth/verify-register", chuc_nang.XuLyXacThucOTP)
	
	router.GET("/tai-khoan", chuc_nang.KiemTraDangNhap, chuc_nang.TrangHoSo)

	// --- [MỚI] CỔNG MERCHANT PORTAL (QUẢN LÝ HẠ TẦNG CỬA HÀNG) ---
	merchant := router.Group("/cua-hang")
	merchant.Use(chuc_nang.KiemTraDangNhap)
	{
		merchant.GET("/", chuc_nang_master.TrangQuanLyCuaHang)
		merchant.POST("/api/cap-nhat-id", chuc_nang_master.API_CapNhatHaTang)
	}

	// API Public
	api := router.Group("/api")
	{
		api.GET("/san-pham", chuc_nang.API_LayDanhSachSanPham)
		api.GET("/cau-hinh", chuc_nang.API_LayMenu)
		api.GET("/san-pham/:id", chuc_nang.API_ChiTietSanPham)
		api.POST("/auth/send-otp", chuc_nang.XuLyGuiOTPEmail)
		api.POST("/auth/reset-by-pin", chuc_nang.XuLyQuenPassBangPIN)
		api.POST("/auth/reset-by-otp", chuc_nang.XuLyQuenPassBangOTP)
	}

	// API User
	userApi := router.Group("/api/user")
	userApi.Use(chuc_nang.KiemTraDangNhap)
	{
		userApi.POST("/update-info", chuc_nang.API_DoiThongTin)
		userApi.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userApi.POST("/change-pin", chuc_nang.API_DoiMaPin)
		userApi.POST("/send-otp-pin", chuc_nang.API_GuiOTPPin)
	}

	// Admin Area (Dùng chung)
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraDangNhap, chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", chuc_nang_admin.TrangTongQuan)
		admin.GET("/reload", chuc_nang_admin.API_NapLaiDuLieu)
		
		admin.GET("/thanh-vien", chuc_nang_admin.TrangQuanLyThanhVien)
		admin.POST("/api/member/save", chuc_nang_admin.API_Admin_LuuThanhVien)
		admin.POST("/api/member/send-msg", chuc_nang_admin.API_Admin_GuiTinNhan)

		admin.GET("/cai-dat", chuc_nang_admin.TrangQuanLyCaiDat)
		admin.POST("/api/category/save", chuc_nang_admin.API_LuuDanhMuc)
		admin.POST("/api/brand/save", chuc_nang_admin.API_LuuThuongHieu)
		admin.POST("/api/margin/save", chuc_nang_admin.API_LuuBienLoiNhuan)
		admin.POST("/api/category/sync-slots", chuc_nang_admin.API_DongBoSlotDanhMuc)

		// --- ĐỊNH TUYẾN RIÊNG CHO NGÀNH MÁY TÍNH ---
		pc := admin.Group("/pc")
		{
			pc.GET("/san-pham", admin_may_tinh.TrangQuanLySanPham)
			pc.POST("/api/product/save", admin_may_tinh.API_LuuSanPham)
		}
	}

	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }

	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server đang chạy tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("✅ Server tắt an toàn.")
}
