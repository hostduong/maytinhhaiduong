package main

import (
	"embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"app/cau_hinh"
	"app/chuc_nang"
	"app/core"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG HỆ THỐNG...")

	cau_hinh.KhoiTaoCauHinh()
	core.KhoiTaoNenTang()
	core.KhoiTaoWorkerGhiSheet()

	log.Println("📦 [BOOT] Đang nạp dữ liệu Master Data...")
	core.NapPhanQuyen("") 
	core.NapSanPham("")
	core.NapKhachHang("")

	router := gin.Default()

	// --- [QUAN TRỌNG] Đăng ký các hàm hỗ trợ cho HTML (FuncMap) ---
	funcMap := template.FuncMap{
		"split": strings.Split,
		
		// 1. Hàm format tiền: 1000000 -> 1.000.000
		"format_money": func(n float64) string {
			p := message.NewPrinter(language.Vietnamese)
			return p.Sprintf("%.0f", n)
		},

		// 2. Hàm chuyển struct sang JSON (Để JS dùng an toàn)
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	}
	// -------------------------------------------------------------

	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER (Giữ nguyên như cũ) ---
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

	// Admin
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraDangNhap, chuc_nang.KiemTraQuyenHan) 
	{
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
		
		admin.GET("/san-pham", chuc_nang.TrangQuanLySanPham)
		admin.POST("/api/product/save", chuc_nang.API_LuuSanPham)
		admin.POST("/api/product/delete", chuc_nang.API_XoaSanPham)
		
		admin.GET("/thanh-vien", chuc_nang.TrangQuanLyThanhVien)
		admin.POST("/api/member/save", chuc_nang.API_Admin_LuuThanhVien)
		admin.GET("/danh-muc", chuc_nang.TrangQuanLyDanhMuc)
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }

	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server chạy tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("✅ Server tắt an toàn.")
}
