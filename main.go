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
	"app/chuc_nang" // Giữ lại để lấy funcMap (Format tiền, số...)
	"app/core"
	"app/routers"

	"github.com/gin-gonic/gin"
)

// BẮT BUỘC: Quét thư mục giao_dien_he_thong (bao gồm cả file nằm trực tiếp và file trong thư mục con)
//go:embed giao_dien_he_thong/*.html giao_dien_he_thong/*/*.html
var f embed.FS

func main() {
	log.Println(">>> [99K.VN SAAS] KHỞI ĐỘNG HỆ THỐNG KIẾN TRÚC LÕI V1.0...")

	cau_hinh.KhoiTaoCauHinh()
	core.KhoiTaoNenTang() 
	core.KhoiTaoWorkerGhiSheet()

	log.Println("📦 [BOOT] Đang nạp toàn bộ Master Data lên RAM...")
	core.NapPhanQuyen("")
	core.NapKhachHang("")
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapBienLoiNhuan("")
	core.NapNhaCungCap("")
	core.NapMayTinh("")
	core.NapTinNhan("")

	// Khởi tạo phòng Điều phối Router
	router := routers.SetupRouter()
	
	// Nạp hàm tiện ích cho HTML và Build UI từ Embed
	funcMap := chuc_nang.LayBoHamHTML()
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, "giao_dien/*.html", "giao_dien/*/*.html"))
	router.SetHTMLTemplate(templ)

	// Mở cổng mạng
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: router}

	go func() {
		log.Printf("✅ [RUNNING] Cỗ máy SaaS đang lắng nghe tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SẬP MẠNG: %v", err)
		}
	}()

	// Đóng băng hệ thống an toàn khi tắt Server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("\n⚠️ [SHUTDOWN] Đang tiến hành đóng băng hệ thống...")
	core.ProcessQueue() 
	log.Println("✅ [SHUTDOWN] Đóng băng thành công! Tạm biệt.")
}
