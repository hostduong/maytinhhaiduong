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

// BẮT BUỘC: Quét tất cả các thư mục giao diện
//go:embed themes/*/*.html
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

	// Kích hoạt Động cơ Nạp Lõi & Admin chạy ngầm
	go core.KhoiDongHeThongNapDuLieu()

	router := routers.SetupRouter()
	
	// Chỉ đường cho máy chủ tìm thư mục themes mới
	funcMap := layBoHamHTML()
	templ := template.Must(template.New("").Funcs(funcMap).ParseFS(f, 
		"themes/default/*.html", 
		"themes/template_master/*.html",
		"themes/template_admin/*.html", // <--- BỔ SUNG THÊM DÒNG NÀY LÀ XONG!
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
