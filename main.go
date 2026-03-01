package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/config"
	"app/core"
	"app/routers"
)

//go:embed giao_dien/*.html giao_dien/*/*.html
var f embed.FS

func main() {
	log.Println(">>> [99K.VN SAAS] KHỞI ĐỘNG HỆ THỐNG KIẾN TRÚC LÕI V1.0...")

	config.KhoiTaoCauHinh()
	core.KhoiTaoNenTang() 
	core.KhoiTaoWorkerGhiSheet()

	// [THAY ĐỔI LỚN]: Đẩy quá trình nạp RAM vào một tiến trình chạy nền (Background Goroutine)
	// Để Server không bị block và có thể mở Port báo cáo cho Google Cloud ngay lập tức.
	go func() {
		log.Println("📦 [BOOT BACKGROUND] Đang nạp toàn bộ Master Data lên RAM...")
		core.NapPhanQuyen("")
		core.NapKhachHang("")
		core.NapDanhMuc("")
		core.NapThuongHieu("")
		core.NapBienLoiNhuan("")
		core.NapNhaCungCap("")
		core.NapMayTinh("")
		core.NapTinNhan("")
		log.Println("✅ [BOOT BACKGROUND] Nạp dữ liệu hoàn tất!")
	}()

	router := routers.SetupRouter()
	
	basicFuncMap := template.FuncMap{ "dummy": func() string { return "" } }
	templ := template.Must(template.New("").Funcs(basicFuncMap).ParseFS(f, "giao_dien/*.html", "giao_dien/*/*.html"))
	router.SetHTMLTemplate(templ)

	// MỞ CỔNG MẠNG BÁO CÁO GOOGLE CLOUD NGAY
	port := config.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: router}

	go func() {
		log.Printf("✅ [RUNNING] Cỗ máy SaaS đang lắng nghe tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SẬP MẠNG: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("\n⚠️ [SHUTDOWN] Đang tiến hành đóng băng hệ thống...")
	core.ProcessQueue() 
	log.Println("✅ [SHUTDOWN] Đóng băng thành công! Tạm biệt.")
}
