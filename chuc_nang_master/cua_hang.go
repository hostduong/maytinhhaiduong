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
	chuyenNganh := strings.TrimSpace(c.PostForm("chuyen_nganh"))
	domain := strings.TrimSpace(c.PostForm("custom_domain"))
	folderDrive := strings.TrimSpace(c.PostForm("folder_drive_id"))
	authJson := strings.TrimSpace(c.PostForm("google_auth_json"))

	// 1. Validate Bắt buộc
	if sheetID == "" || chuyenNganh == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập Spreadsheet ID và chọn Chuyên ngành kinh doanh!"})
		return
	}

	chuShop, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiên đăng nhập không hợp lệ!"})
		return
	}

	// 2. [KIỂM TRA SPREADSHEET]
	err := core.KiemTraVaKhoiTaoSheetNganh(masterShopID, sheetID, authJson, chuyenNganh)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return 
	}

	// 2.5 [KIỂM TRA THÊM DRIVE NẾU CÓ NHẬP]
	if folderDrive != "" {
		errDrive := core.KiemTraFolderDrive(folderDrive, authJson)
		if errDrive != nil {
			c.JSON(200, gin.H{"status": "error", "msg": errDrive.Error()})
			return
		}
	}

	// 3. NẾU VƯỢ QUA 2 LỚP BẢO VỆ -> CẬP NHẬT RAM
	core.KhoaHeThong.Lock()
	chuShop.DataSheets.SpreadsheetID = sheetID
	chuShop.DataSheets.FolderDriveID = folderDrive
	chuShop.DataSheets.GoogleAuthJson = authJson
	chuShop.CauHinh.CustomDomain = domain
	chuShop.CauHinh.ChuyenNganh = chuyenNganh
	chuShop.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

	// 4. ĐẨY XUỐNG HÀNG ĐỢI GHI SHEET
	ghi := core.ThemVaoHangCho
	r := chuShop.DongTrongSheet
	sh := "KHACH_HANG"

	jsonDS := core.ToJSON(chuShop.DataSheets)
	jsonCH := core.ToJSON(chuShop.CauHinh)

	ghi(masterShopID, sh, r, core.CotKH_DataSheetsJson, jsonDS)
	ghi(masterShopID, sh, r, core.CotKH_CauHinhJson, jsonCH)
	ghi(masterShopID, sh, r, core.CotKH_NgayCapNhat, chuShop.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Tuyệt vời! Kết nối Database và Google Drive thành công."})
}
