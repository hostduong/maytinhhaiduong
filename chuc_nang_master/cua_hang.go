package chuc_nang_master

import (
	"net/http"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// TrangQuanLyCuaHang: Render giao diện Dashboard quản lý hạ tầng
func TrangQuanLyCuaHang(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") // ID của File Master 99k.vn
	userID := c.GetString("USER_ID")

	// Lấy thông tin Chủ shop từ DB Master
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

// API_CapNhatHaTang: Lưu SpreadsheetID và Domain riêng
func API_CapNhatHaTang(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	sheetID := strings.TrimSpace(c.PostForm("spreadsheet_id"))
	domain := strings.TrimSpace(c.PostForm("custom_domain"))

	chuShop, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiên đăng nhập không hợp lệ!"})
		return
	}

	// Khóa RAM để cập nhật
	core.KhoaHeThong.Lock()
	chuShop.DataSheets.SpreadsheetID = sheetID
	chuShop.CauHinh.CustomDomain = domain
	chuShop.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

	// Ghi đè JSON xuống Sheet MASTER
	ghi := core.ThemVaoHangCho
	r := chuShop.DongTrongSheet
	sh := "KHACH_HANG"

	jsonDS := core.ToJSON(chuShop.DataSheets)
	jsonCH := core.ToJSON(chuShop.CauHinh)

	ghi(masterShopID, sh, r, core.CotKH_DataSheetsJson, jsonDS)
	ghi(masterShopID, sh, r, core.CotKH_CauHinhJson, jsonCH)
	ghi(masterShopID, sh, r, core.CotKH_NgayCapNhat, chuShop.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã cập nhật hệ thống máy chủ Cửa hàng thành công!"})
}
