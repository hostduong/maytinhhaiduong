package goi_dich_vu_master

import (
	"app/core"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrangGoiDichVuMaster(c *gin.Context) {
	masterID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	core.KhoaHeThong.RLock()
	kh := core.CacheMapKhachHang[core.TaoCompositeKey(masterID, userID)]
	listGoi := core.CacheGoiDichVu[masterID]
	core.KhoaHeThong.RUnlock()

	if kh == nil { c.Redirect(http.StatusFound, "/login"); return }

	// Chuyển List thành chuỗi JSON để truyền ra Javascript
	jsonBytes, _ := json.Marshal(listGoi)

	c.HTML(http.StatusOK, "goi_dich_vu_master", gin.H{
		"TieuDe":       "Quản Lý Gói Dịch Vụ SaaS",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"ListGoi":      listGoi,
		"ListGoiJson":  string(jsonBytes), 
	})
}
