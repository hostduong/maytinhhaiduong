package khach_hang

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangCongPortalKhachHang(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // Trạm IdentifyTenant đã cấp thẻ Master ID

	// 1. Đọc RAM Master an toàn
	core.KhoaHeThong.RLock()
	allPackages := core.CacheGoiDichVu[shopID]
	core.KhoaHeThong.RUnlock()

	// 2. Bộ lọc thông minh: Chỉ lấy gói STARTER đang Active
	var starterPackages []*core.GoiDichVu
	for _, p := range allPackages {
		if p.LoaiGoi == "STARTER" && p.TrangThai == 1 {
			starterPackages = append(starterPackages, p)
		}
	}

	// 3. Đẩy ra Giao diện
	c.HTML(http.StatusOK, "portal_khach_hang", gin.H{
		"TieuDe":   "Chào mừng bạn đến với 99K.VN",
		"ListGoi":  starterPackages,
		"IsPortal": true,
	})
}
