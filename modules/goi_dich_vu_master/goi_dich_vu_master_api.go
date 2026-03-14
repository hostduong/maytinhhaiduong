package goi_dich_vu_master

import (
	"app/config"
	"app/core"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

func API_LuuGoiDichVuMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	me, _ := core.LayKhachHang(shopID, userID)
	
	if me == nil || me.BaoMat.MaPinHash == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn chưa thiết lập Mã PIN bảo mật!"})
		return
	}
	if !config.KiemTraMatKhau(pinXacNhan, me.BaoMat.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}
	
	isNew := c.PostForm("is_new") == "true"
	payloadJson := c.PostForm("payload_json")
	
	var input core.GoiDichVu
	if err := json.Unmarshal([]byte(payloadJson), &input); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Dữ liệu cấu hình không hợp lệ: " + err.Error()})
		return
	}

	input.MaGoi = strings.ToUpper(strings.TrimSpace(input.MaGoi))

	if err := Service_XuLyLuu(shopID, isNew, &input); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu cấu hình Gói dịch vụ thành công!"})
}
