package bang_gia

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

// TrangCongPortalBangGia: Hiển thị trang bảng giá cho khách hàng
func TrangCongPortalBangGia(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // Thường là ID file Master

	core.KhoaHeThong.RLock()
	allPackages := core.CacheGoiDichVu[shopID]
	core.KhoaHeThong.RUnlock()

	// Lọc các gói STARTER đang ở trạng thái hoạt động (TrangThai = 1)
	var starterPackages []*core.GoiDichVu
	for _, p := range allPackages {
		if p.LoaiGoi == "STARTER" && p.TrangThai == 1 {
			starterPackages = append(starterPackages, p)
		}
	}

	c.HTML(http.StatusOK, "bang_gia", gin.H{
		"TieuDe":   "Bảng Giá Dịch Vụ - 99K.VN",
		"ListGoi":  starterPackages,
		"IsPortal": true,
	})
}
