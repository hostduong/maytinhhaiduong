package phan_quyen_master

import (
	"app/config"
	"app/core"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
)

func API_LuuPhanQuyenMaster(c *gin.Context) {
	masterID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")

	// Xác định đây có phải là Sáng Lập Viên tối cao không
	isMasterUser := (userID == "0000000000000000001")

	// [LUẬT THÉP 4]: CHẶN ĐỨNG API NẾU KHÔNG PHẢI 001
	if !isMasterUser {
		c.JSON(200, gin.H{"status": "error", "msg": "Vùng cấm: Bạn không có quyền truy cập khu vực này!"})
		return
	}

	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	me, _ := core.LayKhachHang(masterID, userID)
	
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
	
	var input core.PhanQuyen
	if err := json.Unmarshal([]byte(payloadJson), &input); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Dữ liệu phân quyền lỗi định dạng: " + err.Error()})
		return
	}

	// Gọi thẳng vào Trái tim xử lý (Đẩy theo cờ isMasterUser)
	if err := Service_XuLyLuu(masterID, isNew, &input, isMasterUser); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu Phân quyền Hệ thống thành công!"})
}
