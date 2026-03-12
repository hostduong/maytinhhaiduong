package setup

import (
	"net/http"

	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangSetup(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")

	// Đảm bảo dữ liệu đã được nạp hoàn chỉnh từ Sheet
	_ = core.EnsureKhachHangLoaded(shopID)

	// [CHỐT CHẶN]: Nếu đã tồn tại ID 001, tuyệt đối không hiển thị Form
	_, hasGod := core.LayKhachHang(shopID, "0000000000000000001")
	if hasGod {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "setup", gin.H{
		"TieuDe": "Khởi Tạo Lõi Hệ Thống",
		"Loi":    c.Query("loi"),
	})
}
