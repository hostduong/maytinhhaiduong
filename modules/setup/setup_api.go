package setup

import (
	"net/http"
	"strings"

	"app/config"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangSetup(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	_ = core.EnsureKhachHangLoaded(shopID)

	_, hasGod := core.LayKhachHang(shopID, "0000000000000000001")
	if hasGod {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "setup", nil) // Xóa Header/Footer, ném thẳng file HTML thuần
}

func API_Setup(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full"))
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }

	sessionID, sign, err := Service_ThucThiSetup(shopID, hoTen, user, email, pass, maPin, dienThoai, c.PostForm("ngay_sinh"), c.PostForm("gioi_tinh"), c.Request.UserAgent())
	
	if err != nil {
		// [ĐÃ FIX]: Trả về JSON để hiển thị Tab lỗi màu cam
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	maxAge := int(config.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_token", sessionID, maxAge, "/", ".99k.vn", false, true)
	c.SetCookie("session_sign", sign, maxAge, "/", ".99k.vn", false, true)

	// [ĐÃ FIX]: Trả về JSON thành công, do đã set Cookie nên click là login luôn
	c.JSON(200, gin.H{
		"status": "ok", 
		"msg": "Thần Sáng Lập đã giáng trần! Dữ liệu Lõi đã sẵn sàng.",
		"redirect": "/master/tong-quan",
	})
}
