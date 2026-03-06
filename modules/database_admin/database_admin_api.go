package database_admin

import (
	"app/config"
	"app/core"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TrangThietLapDatabaseAdmin: Hàm hiển thị giao diện nhập ID / Tạo tự động
func TrangThietLapDatabaseAdmin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// Lấy thông tin khách để hiển thị Email lên View
	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "database_admin", gin.H{
		"TieuDe":   "Kết Nối Dữ Liệu", // Đổi tên cho hợp với Header
		"NhanVien": kh,              // [QUAN TRỌNG]: Đổi từ KhachHang -> NhanVien
	})
}

// API_ThietLapDatabase: Xử lý logic Clone file hoặc lưu ID thủ công
func API_ThietLapDatabase(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // ID Master
	userID := c.GetString("USER_ID") // ID của Chủ shop đang thao tác

	loaiThietLap := c.PostForm("loai_thiet_lap") // "auto" hoặc "manual"
	sheetIDInput := c.PostForm("spreadsheet_id")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy tài khoản!"})
		return
	}

	var newSpreadsheetID string

	if loaiThietLap == "auto" {
		// Gọi hàm clone file mẫu và share quyền cho Email của khách
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

	// Cập nhật ID mới vào RAM Master
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.DataSheets.SpreadsheetID = newSpreadsheetID
	jsonBytes, _ := json.Marshal(kh.DataSheets)
	strJson := string(jsonBytes)
	row := kh.DongTrongSheet
	tenSubdomain := kh.TenDangNhap // Lấy tên đăng nhập để làm Subdomain
	lock.Unlock()

	// Đẩy lệnh xuống Queue để lưu vào Google Sheet Master
	core.PushUpdate(shopID, core.TenSheetKhachHang, row, core.CotKH_DataSheetsJson, strJson)

	// [SỬA Ở ĐÂY]: Không văng ra subdomain của khách nữa, giữ họ lại Tổng Hành Dinh
	redirectURL := "https://admin.99k.vn/tong-quan"

	c.JSON(200, gin.H{
		"status":       "ok",
		"msg":          "Khởi tạo Database thành công!",
		"redirect_url": redirectURL,
	})
}
