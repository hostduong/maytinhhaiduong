package cau_hinh

import (
	"log"
	"os"
	"time"
)

type CauHinhHeThong struct {
	IdFileSheet    string
	CongChayWeb    string
	GoogleAuthJson string // [MỚI] Thêm trường này
}

var BienCauHinh CauHinhHeThong

const (
	ChuKyGhiSheet = 5 * time.Second
)

func KhoiTaoCauHinh() {
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

	congWeb := os.Getenv("PORT")
	if congWeb == "" {
		congWeb = "8080"
	}

	// Lấy JSON Key từ biến môi trường
	jsonKey := os.Getenv("GOOGLE_JSON_KEY")
	// Nếu không có, bạn có thể hardcode chuỗi JSON vào đây để test (nhưng không khuyến khích)
	
	BienCauHinh = CauHinhHeThong{
		IdFileSheet:    idSheet,
		CongChayWeb:    congWeb,
		GoogleAuthJson: jsonKey,
	}

	log.Println("--- [CẤU HÌNH] Đã tải xong ---")
}
