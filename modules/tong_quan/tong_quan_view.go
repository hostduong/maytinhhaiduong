package tong_quan

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangTongQuanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	// Đọc siêu tốc từ RAM Cache (Không cần hàm helper)
	core.KhoaHeThong.RLock()
	me := core.CacheMapKhachHang[shopID+"__"+userID]
	core.KhoaHeThong.RUnlock()

	if me == nil {
		c.String(http.StatusUnauthorized, "Vui lòng đăng nhập")
		return
	}

	// Trỏ đúng vào tên define "master_tong_quan" trong HTML của bạn
	c.HTML(http.StatusOK, "master_tong_quan", gin.H{
		"TieuDe":   "Tổng quan Master",
		"NhanVien": me,
		"QuyenHan": vaiTro,
	})
}
