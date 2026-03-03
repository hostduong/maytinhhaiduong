package ho_so

import (
	"app/config"
	"app/core"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func API_ThietLapDatabase(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // ID Master
	userID := c.GetString("USER_ID") // ID của Chủ shop đang thao tác

	loaiThietLap := c.PostForm("loai_thiet_lap") // "auto" hoặc "manual"
	sheetIDInput := c.PostForm("spreadsheet_id")

	core.KhoaHeThong.RLock()
	kh, ok := core.CacheMapKhachHang[shopID+"__"+userID]
	core.KhoaHeThong.RUnlock()

	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy tài khoản!"})
		return
	}

	var newSpreadsheetID string

	if loaiThietLap == "auto" {
		// Gọi hàm clone file từ core
		// Lưu ý: config.BienCauHinh.GoogleAuthJson là chuỗi JSON của Service Account
		idMoi, err := core.HamCloneVaCapQuyenSheet(kh.Email, kh.TenDangNhap, config.BienCauHinh.GoogleAuthJson)
		if err != nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Lỗi tạo tự động: " + err.Error()})
			return
		}
		newSpreadsheetID = idMoi
	} else {
		// Nhập thủ công
		if sheetIDInput == "" {
			c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập Spreadsheet ID!"})
			return
		}
		newSpreadsheetID = sheetIDInput
	}

	// Cập nhật vào RAM Master
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.DataSheets.SpreadsheetID = newSpreadsheetID
	jsonBytes, _ := json.Marshal(kh.DataSheets)
	strJson := string(jsonBytes)
	row := kh.DongTrongSheet
	lock.Unlock()

	// Đẩy lệnh xuống Queue để lưu vào Google Sheet Master
	core.PushUpdate(shopID, core.TenSheetKhachHang, row, core.CotKH_DataSheetsJson, strJson)

	// CHỖ NÀY SẼ GỌI API CLOUD RUN ĐỂ TẠO SUBDOMAIN (Ta sẽ hoàn thiện sau)
	// core.TaoCloudRunDomainMapping(kh.TenDangNhap + ".99k.vn")

	c.JSON(200, gin.H{"status": "ok", "msg": "Khởi tạo Database thành công!"})
}
