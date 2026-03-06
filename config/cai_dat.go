package config

import (
	"log"
	"os"
	"time"
)

type CauHinhHeThong struct {
	IdFileSheetMaster string // Của Sếp (Chứa Core Team, Gói Cước)
	IdFileSheetAdmin  string // Của Hệ thống (Chứa danh sách Chủ Shop đã mua gói)
	CongChayWeb       string
	GoogleAuthJson    string
}

var BienCauHinh CauHinhHeThong

const (
	ThoiGianHetHanCookie = 30 * time.Minute
	ThoiGianAnHan        = 5 * time.Minute
)

// Map cứng dùng để test Localhost trên máy tính
var MapDomainShop = map[string]string{
	"localhost": "ID_SHEET_CUA_MOT_SHOP_TEST_NAO_DO", 
}

func KhoiTaoCauHinh() {
	// 1. ID Của Sếp (Vùng Xám)
	idMaster := os.Getenv("SHEET_ID_MASTER")
	if idMaster == "" { idMaster = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8" }

	// 2. ID Của Các Chủ Shop (Vùng Đệm)
	idAdmin := os.Getenv("SHEET_ID_ADMIN")
	if idAdmin == "" { idAdmin = "1afUZ6nmk9VeHkDJ-qZz5wJesqE2JRCRqn8MYzahoFBk" }

	congWeb := os.Getenv("PORT")
	if congWeb == "" { congWeb = "8080" }

	jsonKey := os.Getenv("GOOGLE_JSON_KEY") // Dành cho Local test, Cloud Run sẽ bỏ qua
	
	BienCauHinh = CauHinhHeThong{
		IdFileSheetMaster: idMaster,
		IdFileSheetAdmin:  idAdmin,
		CongChayWeb:       congWeb,
		GoogleAuthJson:    jsonKey,
	}
	log.Println("--- [CẤU HÌNH 4.0] Đã tải xong Kiến trúc Phân mảnh - Port:", congWeb, "---")
}
