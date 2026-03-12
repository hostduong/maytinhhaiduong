package database_admin

import (
	"app/config"
	"app/core"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrangThietLapDatabaseAdmin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "database_admin", gin.H{
		"TieuDe":   "Kết Nối Dữ Liệu", 
		"NhanVien": kh,              
	})
}

func API_ThietLapDatabase(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID") 

	loaiThietLap := c.PostForm("loai_thiet_lap") 
	sheetIDInput := c.PostForm("spreadsheet_id")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy tài khoản!"})
		return
	}

	var newSpreadsheetID string

	if loaiThietLap == "auto" {
		idMoi, err := core.HamCloneVaCapQuyenSheet(kh.Email, kh.TenDangNhap, config.BienCauHinh.GoogleAuthJson)
		if err != nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Lỗi tạo tự động: " + err.Error()})
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

	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	// [ĐÃ FIX]: Sử dụng kh.System thay cho kh.DataSheets
	kh.System.SheetID = newSpreadsheetID
	b, _ := json.Marshal(kh)
	strJson := string(b)
	row := kh.DongTrongSheet
	lock.Unlock()

	// [ĐÃ FIX]: Cập nhật cả cục DataJSON vào cột B
	core.PushUpdate(shopID, core.TenSheetKhachHang, row, core.CotKH_DataJSON, strJson)

	redirectURL := "https://admin.99k.vn/tong-quan"

	c.JSON(200, gin.H{
		"status":       "ok",
		"msg":          "Khởi tạo Database thành công!",
		"redirect_url": redirectURL,
	})
}
