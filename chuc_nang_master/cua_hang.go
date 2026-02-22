package chuc_nang_master

import (
	"net/http"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangQuanLyCuaHang(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")

	chuShop, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "cua_hang_master", gin.H{
		"TieuDe":   "Quản trị Hạ tầng Cửa hàng",
		"ChuShop":  chuShop,
	})
}

func API_CapNhatHaTang(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	sheetID := strings.TrimSpace(c.PostForm("spreadsheet_id"))
	chuyenNganh := strings.TrimSpace(c.PostForm("chuyen_nganh")) // <--- Hứng thêm trường mới
	domain := strings.TrimSpace(c.PostForm("custom_domain"))
	folderDrive := strings.TrimSpace(c.PostForm("folder_drive_id"))
	authJson := strings.TrimSpace(c.PostForm("google_auth_json"))

	// Validate Bắt buộc
	if sheetID == "" || chuyenNganh == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập Spreadsheet ID và chọn Chuyên ngành kinh doanh!"})
		return
	}

	chuShop, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiên đăng nhập không hợp lệ!"})
		return
	}

	core.KhoaHeThong.Lock()
	chuShop.DataSheets.SpreadsheetID = sheetID
	chuShop.DataSheets.FolderDriveID = folderDrive
	chuShop.DataSheets.GoogleAuthJson = authJson
	chuShop.CauHinh.CustomDomain = domain
	chuShop.CauHinh.ChuyenNganh = chuyenNganh // <--- Cập nhật vào RAM
	chuShop.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	
	// Nếu có cấp JSON Riêng -> Kích hoạt API riêng ngay lập tức
	if authJson != "" && sheetID != "" {
		core.KetNoiGoogleSheetRieng(sheetID, authJson)
	}
	core.KhoaHeThong.Unlock()

	ghi := core.ThemVaoHangCho
	r := chuShop.DongTrongSheet
	sh := "KHACH_HANG"

	jsonDS := core.ToJSON(chuShop.DataSheets)
	jsonCH := core.ToJSON(chuShop.CauHinh)

	ghi(masterShopID, sh, r, core.CotKH_DataSheetsJson, jsonDS)
	ghi(masterShopID, sh, r, core.CotKH_CauHinhJson, jsonCH)
	ghi(masterShopID, sh, r, core.CotKH_NgayCapNhat, chuShop.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã cập nhật hệ thống Cửa hàng thành công!"})
}
