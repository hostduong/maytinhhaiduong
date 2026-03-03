package goi_dich_vu

import (
	"app/core"
	"net/http"
	"github.com/gin-gonic/gin"
)

func TrangGoiDichVuMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	core.KhoaHeThong.RLock()
	kh := core.CacheMapKhachHang[shopID+"__"+userID]
	listGoi := core.CacheGoiDichVu[shopID]
	core.KhoaHeThong.RUnlock()

	if kh == nil { c.Redirect(http.StatusFound, "/login"); return }

	c.HTML(http.StatusOK, "master_goi_dich_vu", gin.H{
		"TieuDe":       "Quản Lý Gói Dịch Vụ SaaS",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"ListGoi":      listGoi,
	})
}
