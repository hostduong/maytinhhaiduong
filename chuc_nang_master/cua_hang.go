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

	// [MỚI] Clone object ra để nhét Inbox vào, tránh đụng RAM của nhau
	chuShopCopy := *chuShop 
	chuShopCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, userID, chuShopCopy.VaiTroQuyenHan)

	// Đếm số tin chưa đọc
	soTinChuaDoc := 0
	for _, msg := range chuShopCopy.Inbox {
		if !msg.DaDoc {
			soTinChuaDoc++
		}
	}

	c.HTML(http.StatusOK, "cua_hang_master", gin.H{
		"TieuDe":       "Quản trị Hạ tầng Cửa hàng",
		"ChuShop":      &chuShopCopy, // Đẩy bản copy ra giao diện
		"SoTinChuaDoc": soTinChuaDoc, 
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

	if sheetID == "" || chuyenNganh == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập Spreadsheet ID và chọn Chuyên ngành kinh doanh!"})
		return
	}

	chuShop, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiên đăng nhập không hợp lệ!"})
		return
	}

	err := core.KiemTraVaKhoiTaoSheetNganh(masterShopID, sheetID, authJson, chuyenNganh)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return 
	}

	if folderDrive != "" {
		errDrive := core.KiemTraFolderDrive(folderDrive, authJson)
		if errDrive != nil {
			c.JSON(200, gin.H{"status": "error", "msg": errDrive.Error()})
			return
		}
	}

	core.KhoaHeThong.Lock()
	chuShop.DataSheets.SpreadsheetID = sheetID
	chuShop.DataSheets.FolderDriveID = folderDrive
	chuShop.DataSheets.GoogleAuthJson = authJson
	chuShop.CauHinh.CustomDomain = domain
	chuShop.CauHinh.ChuyenNganh = chuyenNganh
	chuShop.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

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

// API: Đánh dấu tin nhắn đã đọc
func API_DanhDauDaDoc(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	msgID := strings.TrimSpace(c.PostForm("msg_id"))

	// [MỚI]: Gọi thẳng lệnh đánh dấu đọc từ lõi TIN_NHAN mới (Không còn đụng vào JSON cũ)
	core.DanhDauDocTinNhan(masterShopID, userID, msgID)

	c.JSON(200, gin.H{"status": "ok"})
}
