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
	"app/core" // [MỚI] Sử dụng package Core

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG HỆ THỐNG MỚI (CORE)...")

	// 1. Tải cấu hình môi trường (.env)
	cau_hinh.KhoiTaoCauHinh()

	// 2. Khởi tạo Core (Kết nối Google Sheet)
	core.KhoiTaoNenTang()

	// 3. Khởi chạy Worker ghi sheet (Chạy ngầm)
	core.KhoiTaoWorkerGhiSheet()

	// 4. Nạp dữ liệu vào RAM (Ưu tiên)
	// Lưu ý: ID rỗng "" nghĩa là lấy ID mặc định trong Config
	log.Println("📦 [BOOT] Đang nạp dữ liệu Master Data...")
	
	// Sử dụng WaitGroup nếu muốn nạp song song (Tạm thời nạp tuần tự cho an toàn)
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapSanPham("")
	core.NapKhachHang("") 
	// core.NapCauHinhWeb("") ... (Nếu có)

	// 5. Cấu hình Router
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER ---
	
	// Public
	router.GET("/", chuc_nang.TrangChu)
	// router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham) // Tạm đóng để sửa sau
	
	// Auth
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	// router.POST("/register", chuc_nang.XuLyDangKy) // Tạm đóng
	
	// Admin Group
	admin := router.Group("/admin")
	// admin.Use(chuc_nang.KiemTraQuyenHan) // Tạm đóng Middleware cũ để test
	{
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)
		
		// Quản lý sản phẩm (Đã nâng cấp View, chờ nâng cấp Controller)
		admin.GET("/san-pham", chuc_nang.TrangQuanLySanPham)
		admin.POST("/api/product/save", chuc_nang.API_LuuSanPham)
	}

	// --- KHỞI CHẠY SERVER ---
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server chạy tại 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server...")
	// Có thể thêm logic chờ worker ghi hết dữ liệu còn dư
	log.Println("✅ Server tắt an toàn.")
}
