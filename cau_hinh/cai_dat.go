package cau_hinh

import (
	"log"
	"os"
	"time"
)

type CauHinhHeThong struct {
	IdFileSheet string
	CongChayWeb string
}

var BienCauHinh CauHinhHeThong

// [CẤU HÌNH] Chu kỳ ghi dữ liệu xuống Sheet (5 Giây)
const (
	ChuKyGhiSheet = 2 * time.Second
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

	BienCauHinh = CauHinhHeThong{
		IdFileSheet: idSheet,
		CongChayWeb: congWeb,
	}

	log.Println("--- [CẤU HÌNH] Đã tải xong (Mode: Public + Batch 5s) ---")
}
