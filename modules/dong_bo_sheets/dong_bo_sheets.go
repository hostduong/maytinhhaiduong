package dong_bo_sheets

import (
	"net/http"
	"strings"

	"app/config"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangDongBoSheetsMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// [FIX LỖI TRẮNG TRANG]: Bọc an toàn, nếu mất Session thì đá ra Login
	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok || me == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Tạo bản sao (Deep Copy) để Render HTML không bao giờ bị Crash
	meCopy := *me
	if meCopy.MangXaHoi == nil { meCopy.MangXaHoi = make(map[string]string) }

	// Bơm chỉ số giao diện (Level/Theme) để Layout render Sidebar chuẩn màu
	if meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" {
		meCopy.StyleLevel, meCopy.StyleTheme = 0, 9
	} else {
		meCopy.StyleLevel, meCopy.StyleTheme = 1, 4
	}

	c.HTML(http.StatusOK, "master_dong_bo_sheets", gin.H{
		"TieuDe":   "Đồng Bộ Sheets",
		"NhanVien": &meCopy, // Truyền bản Copy an toàn xuống Giao diện
		"QuyenHan": vaiTro,
	})
}

func API_NapLaiDuLieuMasterCoPIN(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	
	// [FIX LỖI]: Bọc an toàn cả luồng gọi API
	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok || me == nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Phiên đăng nhập lỗi!"})
		return
	}

	if me.BaoMat.MaPinHash == "" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Bạn chưa thiết lập mã PIN bảo mật trong phần Hồ sơ!"})
		return
	}

	if !config.KiemTraMatKhau(pinXacNhan, me.BaoMat.MaPinHash) {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}

	go func() {
		core.KhoiDongHeThongNapDuLieu()
	}()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok", 
		"msg": "Xác thực PIN thành công. Đang tải lại dữ liệu ngầm...",
	})
}
