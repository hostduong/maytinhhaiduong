package config

import (
	"log"
	"os"
	"time"
)

// --- CẤU HÌNH HỆ THỐNG ---
type CauHinhHeThong struct {
	IdFileSheet    string
	CongChayWeb    string
	GoogleAuthJson string
}

var BienCauHinh CauHinhHeThong

const (
	ThoiGianHetHanCookie = 30 * time.Minute
	ThoiGianAnHan        = 5 * time.Minute
	GioiHanNguoiDung     = 100 
)

var MapDomainShop = map[string]string{
	"localhost": "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8", 
}

func KhoiTaoCauHinh() {
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

	// [QUAN TRỌNG NHẤT LÀ ĐOẠN NÀY]: Bắt buộc phải nghe theo lời Google Cloud Run
	congWeb := os.Getenv("PORT")
	if congWeb == "" { 
		congWeb = "8080" // Chỉ dùng 8080 nếu chạy Local trên máy tính của bạn
	}

	jsonKey := os.Getenv("GOOGLE_JSON_KEY")
	
	BienCauHinh = CauHinhHeThong{
		IdFileSheet:    idSheet,
		CongChayWeb:    congWeb,
		GoogleAuthJson: jsonKey,
	}
	log.Println("--- [CẤU HÌNH] Đã tải xong (Chế độ SaaS) - Port:", congWeb, "---")
}
