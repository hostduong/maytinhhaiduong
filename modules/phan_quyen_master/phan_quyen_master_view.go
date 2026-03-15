package phan_quyen_master

import (
	"app/core"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrangPhanQuyenMaster(c *gin.Context) {
	masterID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// [LUẬT THÉP 4]: CHỈ DUY NHẤT ID 001 MỚI ĐƯỢC VÀO TRANG NÀY
	if userID != "0000000000000000001" {
		c.Redirect(http.StatusFound, "/master/tong-quan")
		return
	}

	core.KhoaHeThong.RLock()
	kh := core.CacheMapKhachHang[core.TaoCompositeKey(masterID, userID)]
	listPQ := core.CachePhanQuyen[masterID]
	core.KhoaHeThong.RUnlock()

	if kh == nil { c.Redirect(http.StatusFound, "/login"); return }

	jsonBytes, _ := json.Marshal(listPQ)

	c.HTML(http.StatusOK, "phan_quyen_master", gin.H{
		"TieuDe":       "Quản Lý Phân Quyền (RBAC)",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"ListPQ":       listPQ,
		"ListPQJson":   string(jsonBytes), 
	})
}
