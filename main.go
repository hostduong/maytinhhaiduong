package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/config" // Trỏ vào thư mục Config mới tạo
	"app/core"
	"app/routers"
)

// BẮT BUỘC: Quét thư mục giao_dien (bao gồm cả file nằm trực tiếp và file trong thư mục con)
//go:embed giao_dien/*.html giao_dien/*/*.html
var f embed.FS

func main() {
	log.Println(">>> [99K.VN SAAS] KHỞI ĐỘNG HỆ THỐNG KIẾN TRÚC LÕI V1.0...")

	// 1. Khởi tạo cấu hình Server
	config.KhoiTaoCauHinh()
	core.KhoiTaoNenTang() 
	core.KhoiTaoWorkerGhiSheet()

	// 2. Nạp toàn bộ dữ liệu lên RAM
	log.Println("📦 [BOOT] Đang nạp toàn bộ Master Data lên RAM...")
	core.NapPhanQuyen("")
	core.NapKhachHang("")
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapBienLoiNhuan("")
	core.NapNhaCungCap("")
	core.NapMayTinh("")
	core.NapTinNhan("")

	// 3. Khởi tạo phòng Điều phối Router
	router := routers.SetupRouter()
	
	// 4. Định nghĩa một FuncMap cơ bản rỗng để Bypass lỗi Compile
	basicFuncMap := template.FuncMap{
		"dummy": func() string { return "" },
	}
	templ := template.Must(template.New("").Funcs(basicFuncMap).ParseFS(f, "giao_dien/*.html", "giao_dien/*/*.html"))
	router.SetHTMLTemplate(templ)

	// 5. Mở cổng mạng
	port := config.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: router}

	go func() {
		log.Printf("✅ [RUNNING] Cỗ máy SaaS đang lắng nghe tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SẬP MẠNG: %v", err)
		}
	}()

	// 6. Đóng băng hệ thống an toàn khi tắt Server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("\n⚠️ [SHUTDOWN] Đang tiến hành đóng băng hệ thống...")
	core.ProcessQueue() 
	log.Println("✅ [SHUTDOWN] Đóng băng thành công! Tạm biệt.")
}
