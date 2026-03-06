package tong_quan_master

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangTongQuanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")

	nhanVien, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	tongSoShop := 0
	tongSoGoiDaBan := 0

	core.KhoaHeThong.RLock()
	danhSachKhach := core.CacheKhachHang[shopID]
	core.KhoaHeThong.RUnlock()

	for _, kh := range danhSachKhach {
		hasActivePlan := false
		for _, goi := range kh.GoiDichVu {
			if goi.TrangThai == "active" {
				hasActivePlan = true
				tongSoGoiDaBan++
			}
		}
		if hasActivePlan {
			tongSoShop++
		}
	}

	// Gọi đúng tên Define HTML
	c.HTML(http.StatusOK, "tong_quan_master", gin.H{
		"TieuDe":         "Tổng Quan Master",
		"NhanVien":       nhanVien,
		"TongSoShop":     tongSoShop,
		"TongSoGoiDaBan": tongSoGoiDaBan,
	})
}
