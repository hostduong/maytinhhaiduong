package goi_dich_vu_master

import (
	"app/core"
	"encoding/json"
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

	// Chuyển List thành chuỗi JSON để truyền ra Javascript cực gọn
	jsonBytes, _ := json.Marshal(listGoi)

	c.HTML(http.StatusOK, "goi_dich_vu_master", gin.H{
		"TieuDe":       "Quản Lý Gói Dịch Vụ SaaS",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"ListGoi":      listGoi, // Để build HTML Card
		"ListGoiJson":  string(jsonBytes), // Để fill Form Edit JS
	})
}
