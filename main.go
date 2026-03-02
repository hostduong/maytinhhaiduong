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

	"app/config"
	"app/core"
	"app/routers"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// BẮT BUỘC: Quét tất cả các thư mục giao diện (Thêm cả thư mục con của master để dự phòng)
//go:embed giao_dien/*.html giao_dien/*/*.html giao_dien_master/*.html giao_dien_master/*/*.html
var f embed.FS

// --- PHỤC HỒI BỘ HÀM HTML THẬT ĐỂ KHÔNG BỊ CRASH GIAO DIỆN ---
func layBoHamHTML() template.FuncMap {
	return template.FuncMap{
		"firstImg": func(s string) string {
			if s == "" { return "" }
			parts := strings.Split(s, "|")
			return strings.TrimSpace(parts[0])
		},
		"format_money": func(n float64) string {
			p := message.NewPrinter(language.Vietnamese)
			return p.Sprintf("%.0f", n)
		},
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"split": strings.Split,
	}
}

func main() {
	log.Println(">>> [99K.VN SAAS] KHỞI ĐỘNG HỆ THỐNG KIẾN TRÚC LÕI V1.0...")

	config.KhoiTaoCauHinh()
	core.KhoiTaoNenTang() 
	core.KhoiTaoWorkerGhiSheet()

	// Đẩy quá trình nạp RAM xuống nền để Server không bị nghẽn
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
	
	// [ĐÃ SỬA TẠI ĐÂY]: Bổ sung việc Parse (Dịch) toàn bộ file trong giao_dien_master
	funcMap := layBoHamHTML()
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, 
		"giao_dien/*.html", 
		"giao_dien/*/*.html",
		"giao_dien_master/*.html",
		"giao_dien_master/*/*.html",
	))
	router.SetHTMLTemplate(templ)

	// Mở cổng mạng
	port := config.BienCauHinh.CongChayWeb
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
