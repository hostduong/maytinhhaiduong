package cau_hinh

import (
	"log"
	"os"
	"time"
)

// --- CẤU HÌNH HỆ THỐNG ---
type CauHinhHeThong struct {
	IdFileSheet    string // ID mặc định (cho shop gốc hoặc fallback)
	CongChayWeb    string
	GoogleAuthJson string
}

var BienCauHinh CauHinhHeThong

// --- CẤU HÌNH GIỚI HẠN & BẢO MẬT ---
const (
	ThoiGianHetHanCookie = 30 * time.Minute
	ThoiGianAnHan        = 5 * time.Minute
	GioiHanNguoiDung     = 100 // request/giây
)

// --- MAPPING DOMAIN -> SHEET ID (SAAS) ---
// Key: Domain (vd: shop1.com), Value: SpreadsheetID
var MapDomainShop = map[string]string{
	"localhost": "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8", // Mặc định Local
	// Sau này có thể load từ DB hoặc file JSON riêng
}

func KhoiTaoCauHinh() {
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

	congWeb := os.Getenv("PORT")
	if congWeb == "" { congWeb = "8080" }

	jsonKey := os.Getenv("GOOGLE_JSON_KEY")
	
	BienCauHinh = CauHinhHeThong{
		IdFileSheet:    idSheet,
		CongChayWeb:    congWeb,
		GoogleAuthJson: jsonKey,
	}
	log.Println("--- [CẤU HÌNH] Đã tải xong (Chế độ SaaS) ---")
}
