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
	"app/chuc_nang_master"
	"app/core"

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html giao_dien/*/*.html giao_dien_admin/*.html giao_dien_admin/*/*.html giao_dien_master/*.html giao_dien_master/*/*.html
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
	core.NapTinNhan("")
	core.NapMayTinh("")

	router := gin.Default()

	// Mở quyền truy cập thư mục static chứa CSS/JS
	router.Static("/static", "./static")

	router.Use(chuc_nang.GatewaySaaS, chuc_nang.KiemTraGoiDichVu)

	funcMap := chuc_nang.LayBoHamHTML()

	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "giao_dien/*.html", "giao_dien/*/*.html", "giao_dien_admin/*.html", "giao_dien_admin/*/*.html", "giao_dien_master/*.html", "giao_dien_master/*/*.html"))
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
	router.GET("/forgot-password", chuc_nang.TrangQuenMatKhau)

	router.GET("/tai-khoan", chuc_nang.KiemTraDangNhap, chuc_nang.TrangHoSo)

	// Master Area (Quản trị hệ thống Lõi)
	master := router.Group("/master")
	master.Use(chuc_nang.KiemTraDangNhap, chuc_nang.KiemTraQuyenHan)
	{
		master.GET("/tong-quan", chuc_nang_master.TrangTongQuanMaster)
		master.GET("/api/reload", chuc_nang_master.API_NapLaiDuLieuMaster)

		master.GET("/ho-so", chuc_nang_master.TrangHoSoMaster)
		master.POST("/api/ho-so", chuc_nang_master.API_LuuHoSoMaster)
		master.POST("/api/change-pass", chuc_nang_master.API_DoiMatKhauMaster)
		master.POST("/api/change-pin", chuc_nang_master.API_DoiMaPinMaster)

		master.GET("/thanh-vien", chuc_nang_master.TrangQuanLyThanhVienMaster)
		master.POST("/api/thanh-vien/save", chuc_nang_master.API_LuuThanhVienMaster)
		master.POST("/api/thanh-vien/send-msg", chuc_nang_master.API_GuiTinNhanMaster)

		master.GET("/tin-nhan", chuc_nang_master.TrangTinNhanMaster)
		master.POST("/api/doc-tin-nhan", chuc_nang_master.API_DanhDauDaDocMaster)
		master.POST("/api/tin-nhan/send-chat", chuc_nang_master.API_GuiTinNhanChat)

		master.GET("/dong-bo-sheets", chuc_nang_master.TrangDongBoSheetsMaster)
		master.POST("/api/dong-bo-sheets", chuc_nang_master.API_NapLaiDuLieuMasterCoPIN)
		
        master.GET("/nhap-hang", chuc_nang_master.TrangNhapHangMaster)
		master.GET("/quan-ly-may-tinh", chuc_nang_master.TrangQuanLyMayTinhMaster)
		master.GET("/api/may-tinh/detail/:ma_sp", chuc_nang_master.API_LayChiTietMayTinhMaster)
		master.POST("/api/may-tinh/save", chuc_nang_master.API_LuuMayTinhMaster)

		master.GET("/cai-dat-cau-hinh", chuc_nang_master.TrangCaiDatCauHinhMaster)
		master.POST("/api/cai-dat-cau-hinh/danh-muc/save", chuc_nang_master.API_LuuDanhMucMaster)
		master.POST("/api/cai-dat-cau-hinh/danh-muc/sync-slots", chuc_nang_master.API_DongBoSlotDanhMucMaster)
		master.POST("/api/cai-dat-cau-hinh/thuong-hieu/save", chuc_nang_master.API_LuuThuongHieuMaster)
		master.POST("/api/cai-dat-cau-hinh/bien-loi-nhuan/save", chuc_nang_master.API_LuuBienLoiNhuanMaster)
		master.POST("/api/cai-dat-cau-hinh/ncc/save", chuc_nang_master.API_LuuNhaCungCapMaster)
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
		userApi.POST("/verify-softgate", chuc_nang.API_XacThucKichHoat)
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

		admin.GET("/nhap-hang", chuc_nang_admin.TrangNhapHang)
	}

	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: router}

	go func() {
		log.Printf("✅ Server đang chạy tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ [SHUTDOWN] Đã nhận lệnh tắt Server...")
	log.Println("💾 [SHUTDOWN] Đang xả toàn bộ dữ liệu trên RAM xuống Google Sheets lần cuối...")
	core.ThucHienGhiSheet()

	log.Println("✅ [SHUTDOWN] Server tắt an toàn. Không mất dữ liệu.")
}
