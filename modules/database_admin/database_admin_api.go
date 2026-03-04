package ho_so

import (
	"app/config"
	"app/core"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

// API_ThietLapDatabase: Xử lý cài đặt ID Google Sheet cho shop
func API_ThietLapDatabase(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") // Thường là ID file Master
	userID := c.GetString("USER_ID")       // ID của chủ shop

	loaiThietLap := c.PostForm("loai_thiet_lap")
	sheetIDInput := c.PostForm("spreadsheet_id")

	// 1. TRẠM KIỂM SOÁT ZERO-TRUST: Lấy thông tin khách hàng từ RAM Master
	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiên đăng nhập không hợp lệ!"})
		return
	}

	// 2. KIỂM TRA GÓI CƯỚC: Phải có ít nhất 1 gói STARTER mới được cài Database
	hasStarter := false
	for _, p := range kh.GoiDichVu {
		if p.LoaiGoi == "STARTER" && p.TrangThai == "active" {
			hasStarter = true
			break
		}
	}
	if !hasStarter {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn chưa kích hoạt gói dịch vụ STARTER!"})
		return
	}

	var newSpreadsheetID string

	// 3. THỰC THI KHỞI TẠO
	if loaiThietLap == "auto" {
		// Gọi hàm Clone Sheet mẫu từ core (Sử dụng Service Account của hệ thống)
		idMoi, err := core.HamCloneVaCapQuyenSheet(kh.Email, kh.TenDangNhap, config.BienCauHinh.GoogleAuthJson)
		if err != nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Lỗi hạ tầng Google: " + err.Error()})
			return
		}
		newSpreadsheetID = idMoi
	} else {
		if sheetIDInput == "" {
			c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập Spreadsheet ID!"})
			return
		}
		newSpreadsheetID = sheetIDInput
	}

	// 4. CẬP NHẬT RAM VÀ ĐẨY QUEUE (Sheet KHACH_HANG)
	lock := core.GetSheetLock(masterShopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.DataSheets.SpreadsheetID = newSpreadsheetID
	jsonBytes, _ := json.Marshal(kh.DataSheets)
	strJson := string(jsonBytes)
	row, tenSubdomain := kh.DongTrongSheet, kh.TenDangNhap
	lock.Unlock()

	// Đẩy lệnh ghi Sheet Master qua Queue
	core.PushUpdate(masterShopID, core.TenSheetKhachHang, row, core.CotKH_DataSheetsJson, strJson)

	// 5. PHẢN HỒI VÀ BẺ LÁI TUYỆT ĐỐI
	// Trả về URL dẫn về Subdomain riêng của họ để bắt đầu làm việc
	redirectURL := fmt.Sprintf("https://%s.99k.vn/admin/tong-quan", tenSubdomain)

	c.JSON(200, gin.H{
		"status":       "ok", 
		"msg":          "Database đã được thiết lập!",
		"redirect_url": redirectURL,
	})
}
