package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type CauHinhHeThong struct {
	IdFileSheetMaster string
	IdFileSheetAdmin  string
	CongChayWeb       string
	GoogleAuthJson    string
	
	// CẤU HÌNH QUẢN TRỊ BỘ NHỚ (RAM)
	MaxRamMB         int // Tổng RAM của máy chủ (VD: 1024 MB)
	HighWatermarkPct int // Mức báo động đỏ (VD: 75%)
	LowWatermarkPct  int // Mức an toàn sau khi xả lũ (VD: 60%)
}

var BienCauHinh CauHinhHeThong

const (
	ThoiGianHetHanCookie = 30 * time.Minute
	ThoiGianAnHan        = 5 * time.Minute
)

var MapDomainShop = map[string]string{
	"localhost": "ID_SHEET_CUA_MOT_SHOP_TEST_NAO_DO", 
}

// Hàm phụ trợ lấy số từ ENV
func getEnvInt(key string, defaultVal int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil { return val }
	}
	return defaultVal
}

func KhoiTaoCauHinh() {
	idMaster := os.Getenv("SHEET_ID_MASTER")
	if idMaster == "" { idMaster = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8" }

	idAdmin := os.Getenv("SHEET_ID_ADMIN")
	if idAdmin == "" { idAdmin = "1afUZ6nmk9VeHkDJ-qZz5wJesqE2JRCRqn8MYzahoFBk" }

	congWeb := os.Getenv("PORT")
	if congWeb == "" { congWeb = "8080" }

	jsonKey := os.Getenv("GOOGLE_JSON_KEY")
	
	BienCauHinh = CauHinhHeThong{
		IdFileSheetMaster: idMaster,
		IdFileSheetAdmin:  idAdmin,
		CongChayWeb:       congWeb,
		GoogleAuthJson:    jsonKey,
		
		// Setup RAM linh hoạt
		MaxRamMB:         getEnvInt("MAX_RAM_MB", 1024), // Mặc định 1GB (1024MB)
		HighWatermarkPct: getEnvInt("HIGH_WATERMARK_PCT", 75), // Báo động ở 75%
		LowWatermarkPct:  getEnvInt("LOW_WATERMARK_PCT", 60),  // Xóa đến khi còn 60% thì dừng
	}
	log.Printf("--- [CẤU HÌNH 4.1] Khởi động. RAM Lõi: %dMB (Báo động: %d%%) ---", BienCauHinh.MaxRamMB, BienCauHinh.HighWatermarkPct)
}
