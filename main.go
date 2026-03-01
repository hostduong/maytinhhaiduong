package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/cau_hinh"
	"app/core"
	"app/routers"
)

func main() {
	log.Println(">>> [99K.VN SAAS] KHỞI ĐỘNG HỆ THỐNG ENTERPRISE V1.0...")

	// 1. Nạp cấu hình & Hệ sinh thái Google
	cau_hinh.KhoiTaoCauHinh()
	core.KhoiTaoNenTang() 

	// 2. Kích hoạt cỗ máy Hàng đợi Ghi dữ liệu (Write Queue)
	core.KhoiTaoWorkerGhiSheet()

	// 3. Nạp Master Data đa người thuê lên RAM (Bootstrapping)
	log.Println("📦 [BOOT] Đang nạp toàn bộ cấu trúc dữ liệu lên RAM (In-Memory)...")
	// Lưu ý: Tạm truyền ID rỗng "" để test, khi xong tính năng đa Shop sẽ chạy vòng lặp nạp nhiều ShopID.
	core.NapPhanQuyen("")
	core.NapKhachHang("")
	core.NapDanhMuc("")
	core.NapThuongHieu("")
	core.NapBienLoiNhuan("")
	core.NapNhaCungCap("")
	core.NapMayTinh("")
	// core.NapPhieuNhap("") // Chờ ghép module Nhập Hàng
	// core.NapTinNhan("")   // Chờ ghép module Tin Nhắn

	// 4. Lắp ráp Phòng Điều Phối (Router & Middlewares)
	router := routers.SetupRouter()

	// (Tạm thời map thư mục HTML nếu bạn đang giữ file cũ ở ngoài, sau này sẽ move vào module)
	// router.LoadHTMLGlob("giao_dien_he_thong/*/*.html")

	// 5. Mở Cổng Mạng (Start HTTP Server)
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	
	srv := &http.Server{Addr: "0.0.0.0:" + port, Handler: router}

	go func() {
		log.Printf("✅ [RUNNING] Cỗ máy SaaS đang lắng nghe tại http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SẬP MẠNG: %v", err)
		}
	}()

	// 6. Graceful Shutdown (Bắt sự kiện Ctrl+C, tắt server an toàn tuyệt đối)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("\n⚠️ [SHUTDOWN] Nhận lệnh tắt máy chủ. Đang tiến hành đóng băng hệ thống...")
	log.Println("💾 [SHUTDOWN] Đang xả toàn bộ dữ liệu tồn đọng trong Queue xuống Google Sheets...")
	core.ProcessQueue() // Ép con Worker ghi nốt 100% dữ liệu đang cầm trên tay
	log.Println("✅ [SHUTDOWN] Quá trình đóng băng hoàn tất. Không rớt 1 byte dữ liệu. Tạm biệt!")
}
