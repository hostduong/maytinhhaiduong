package chuc_nang_master

import (
	"net/http"
	"strings"

	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangTinNhanMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	meCopy := *me 
	meCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, userID, vaiTro)

	soTinChuaDoc := 0
	for _, msg := range meCopy.Inbox {
		if !msg.DaDoc {
			soTinChuaDoc++
		}
	}

	c.HTML(http.StatusOK, "master_tin_nhan", gin.H{
		"TieuDe":       "Hộp thư hệ thống",
		"NhanVien":     &meCopy,
		"SoTinChuaDoc": soTinChuaDoc, 
	})
}

func API_DanhDauDaDocMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	msgID := strings.TrimSpace(c.PostForm("msg_id"))

	// Kích hoạt hàm lõi để chuyển trạng thái "Đã đọc"
	core.DanhDauDocTinNhan(masterShopID, userID, msgID)

	c.JSON(200, gin.H{"status": "ok"})
}
