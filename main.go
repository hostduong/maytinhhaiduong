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
	"app/core" // Sử dụng Core

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG HỆ THỐNG MỚI (CORE)...")

	// 1. Tải cấu hình
	cau_hinh.KhoiTaoCauHinh()

	// 2. Khởi tạo Core (Kết nối & Worker)
	core.KhoiTaoNenTang()
	core.KhoiTaoWorkerGhiSheet()

	// 3. Nạp dữ liệu vào RAM
	log.Println("📦 [BOOT] Đang nạp dữ liệu Master Data...")
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapSanPham("")
	core.NapKhachHang("")

	// 4. Cấu hình Router
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER ---

	// A. Public (Không cần đăng nhập)
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	// Auth
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	router.GET("/logout", chuc_nang.DangXuat) // Hàm này nằm trong dang_nhap.go

	// Chức năng Tài khoản & Quên mật khẩu
	router.GET("/tai-khoan", chuc_nang.TrangHoSo)           // [MỚI]
	router.GET("/quen-mat-khau", chuc_nang.TrangQuenMatKhau) // [MỚI]

	// B. API Public (Ajax gọi)
	api := router.Group("/api")
	{
		// Sản phẩm & Menu
		api.GET("/san-pham", chuc_nang.API_LayDanhSachSanPham)
		api.GET("/cau-hinh", chuc_nang.API_LayMenu)
		api.GET("/san-pham/:id", chuc_nang.API_ChiTietSanPham)

		// Quên mật khẩu (Auth API)
		api.POST("/auth/send-otp", chuc_nang.XuLyGuiOTPEmail)
		api.POST("/auth/reset-by-pin", chuc_nang.XuLyQuenPassBangPIN)
		api.POST("/auth/reset-by-otp", chuc_nang.XuLyQuenPassBangOTP)
	}

	// C. API User (Cần đăng nhập)
	userApi := router.Group("/api/user")
	userApi.Use(chuc_nang.KiemTraDangNhap) // Middleware chặn nếu chưa login
	{
		userApi.POST("/update-info", chuc_nang.API_DoiThongTin)
		userApi.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userApi.POST("/change-pin", chuc_nang.API_DoiMaPin)
		userApi.POST("/send-otp-pin", chuc_nang.API_GuiOTPPin)
	}

	// D. Admin Group
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan) // Middleware check quyền Admin
	{
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
		
		// Quản lý sản phẩm
		admin.GET("/san-pham", chuc_nang.TrangQuanLySanPham)
		admin.POST("/api/product/save", chuc_nang.API_LuuSanPham)
		
		// Quản lý thành viên (API)
		admin.POST("/api/member/update", chuc_nang.API_Admin_SuaThanhVien)
	}

	// --- KHỞI CHẠY SERVER ---
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server chạy tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠️ Đang tắt Server...")
	log.Println("✅ Server tắt an toàn.")
}
