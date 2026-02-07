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
	"app/core" // [QUAN TRỌNG] Sử dụng Core mới

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG HỆ THỐNG MỚI (CORE)...")

	// 1. Tải cấu hình
	cau_hinh.KhoiTaoCauHinh()

	// 2. Khởi tạo Core (Kết nối Google Sheet & Worker)
	core.KhoiTaoNenTang()
	core.KhoiTaoWorkerGhiSheet()

	// 3. Nạp dữ liệu vào RAM (Master Data)
	log.Println("📦 [BOOT] Đang nạp dữ liệu từ Google Sheet...")
	
	// Nạp dữ liệu của Shop hiện tại (ID trong Config)
	// Hàm Nap...("") nghĩa là lấy ID mặc định
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapSanPham("")
	core.NapKhachHang("")

	// 4. Cấu hình Router
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER ---

	// Public
	router.GET("/", chuc_nang.TrangChu)

	// Auth (Tạm thời vẫn dùng logic cũ, sẽ refactor sau)
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)

	// Admin Group
	admin := router.Group("/admin")
	{
		// Dashboard (Tạm dùng logic cũ)
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)

		// [ĐÃ NÂNG CẤP] Quản lý sản phẩm dùng app/core
		admin.GET("/san-pham", chuc_nang.TrangQuanLySanPham)
		admin.POST("/api/product/save", chuc_nang.API_LuuSanPham)
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
