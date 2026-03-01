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
	"app/chuc_nang" // Tạm giữ để lấy hàm FuncMap (Format tiền...)
	"app/core"
	"app/routers"
)

// Khai báo nhúng toàn bộ thư mục giao diện thống nhất mới
//go:embed giao_dien_he_thong/*.html giao_dien_he_thong/*/*.html
var f embed.FS

func main() {
	log.Println(">>> [99K.VN SAAS] KHỞI ĐỘNG HỆ THỐNG ENTERPRISE V1.0...")

	// 1. Nạp cấu hình & Hệ sinh thái Google
	cau_hinh.KhoiTaoCauHinh()
	core.KhoiTaoNenTang() 

	// 2. Kích hoạt cỗ máy Hàng đợi Ghi dữ liệu kép (Write Queue)
	core.KhoiTaoWorkerGhiSheet()

	// 3. Nạp Master Data đa người thuê lên RAM (Bootstrapping)
	log.Println("📦 [BOOT] Đang nạp toàn bộ cấu trúc dữ liệu lên RAM (In-Memory)...")
	core.NapPhanQuyen("")
	core.NapKhachHang("")
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapBienLoiNhuan("")
	core.NapNhaCungCap("")
	core.NapMayTinh("")
	core.NapTinNhan("")
	// core.NapPhieuNhap("") // Chờ sửa xong module Nhập Hàng sẽ mở ra

	// 4. Lắp ráp Phòng Điều Phối & Load Giao diện
	router := routers.SetupRouter()
	
	// Nạp FuncMap (Format số, tiền...) từ code cũ của bạn
	funcMap := chuc_nang.LayBoHamHTML()
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "giao_dien_he_thong/*.html", "giao_dien_he_thong/*/*.html"))
	router.SetHTMLTemplate(templ)

	// 5. Mở Cổng Mạng
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: router}

	go func() {
		log.Printf("✅ [RUNNING] Cỗ máy SaaS đang lắng nghe tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SẬP MẠNG: %v", err)
		}
	}()

	// 6. Graceful Shutdown (Bắt tín hiệu tắt Server)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("\n⚠️ [SHUTDOWN] Đang tiến hành đóng băng hệ thống...")
	core.ProcessQueue() // Ép Worker ghi nốt 100% dữ liệu đang cầm trên tay
	log.Println("✅ [SHUTDOWN] Đóng băng thành công. Không rớt 1 byte. Tạm biệt!")
}
