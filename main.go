package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings" // [1] BẮT BUỘC CÓ để dùng hàm split
	"syscall"

	"app/cau_hinh"
	"app/chuc_nang"
	"app/core"

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG HỆ THỐNG...")

	// 1. Core & Config
	cau_hinh.KhoiTaoCauHinh()
	core.KhoiTaoNenTang()
	core.KhoiTaoWorkerGhiSheet()

	// 2. Nạp dữ liệu
	log.Println("📦 [BOOT] Đang nạp dữ liệu Master Data...")
	core.NapPhanQuyen("") 
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapSanPham("")
	core.NapKhachHang("")

	// 3. Router & Template
	router := gin.Default()

	// [2] Đăng ký hàm "split" để file HTML dùng (Khắc phục lỗi trắng trang)
	funcMap := template.FuncMap{
		"split": strings.Split,
	}

	// Nạp template với FuncMap
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER ---

	// A. Public
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	// Auth
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	router.GET("/logout", chuc_nang.DangXuat)
	
	// User Profile
	router.GET("/tai-khoan", chuc_nang.TrangHoSo)
	router.GET("/forgot-password", chuc_nang.TrangQuenMatKhau)

	// B. API Public
	api := router.Group("/api")
	{
		api.GET("/san-pham", chuc_nang.API_LayDanhSachSanPham)
		api.GET("/cau-hinh", chuc_nang.API_LayMenu)
		api.GET("/san-pham/:id", chuc_nang.API_ChiTietSanPham)
		api.POST("/auth/send-otp", chuc_nang.XuLyGuiOTPEmail)
		api.POST("/auth/reset-by-pin", chuc_nang.XuLyQuenPassBangPIN)
		api.POST("/auth/reset-by-otp", chuc_nang.XuLyQuenPassBangOTP)
	}

	// C. API User (Cần đăng nhập)
	userApi := router.Group("/api/user")
	userApi.Use(chuc_nang.KiemTraDangNhap)
	{
		userApi.POST("/update-info", chuc_nang.API_DoiThongTin)
		userApi.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userApi.POST("/change-pin", chuc_nang.API_DoiMaPin)
		userApi.POST("/send-otp-pin", chuc_nang.API_GuiOTPPin)
	}

// D. Admin Group
	admin := router.Group("/admin")
	
	// [SỬA LẠI DÒNG NÀY] 
	// Phải gọi cả 2 theo thứ tự: Đăng Nhập (để lấy User) -> Quyền Hạn (để check Role)
	admin.Use(chuc_nang.KiemTraDangNhap, chuc_nang.KiemTraQuyenHan) 
	{
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
		
		// Quản lý sản phẩm
		admin.GET("/san-pham", chuc_nang.TrangQuanLySanPham)
		admin.POST("/api/product/save", chuc_nang.API_LuuSanPham)
		
		// Quản lý thành viên
		admin.GET("/thanh-vien", chuc_nang.TrangQuanLyThanhVien)
		admin.POST("/api/member/save", chuc_nang.API_Admin_LuuThanhVien) 
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
