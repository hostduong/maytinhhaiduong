package tong_quan_admin

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangTongQuanAdmin(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")

	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	maxSanPham := 0
	maxNhanVien := 0
	tenGoiHienTai := "Chưa kích hoạt"

	for _, p := range kh.GoiDichVu {
		if p.TrangThai == "active" {
			maxSanPham += p.MaxSanPham
			maxNhanVien += p.MaxNhanVien
			if p.LoaiGoi == "STARTER" {
				tenGoiHienTai = p.TenGoi
			}
		}
	}

	soSanPhamHienTai := 0
	if kh.DataSheets.SpreadsheetID != "" {
		core.KhoaHeThong.RLock()
		if nganhMap, exists := core.CacheSanPham[kh.DataSheets.SpreadsheetID]; exists {
			for _, ds := range nganhMap {
				soSanPhamHienTai += len(ds)
			}
		}
		core.KhoaHeThong.RUnlock()
	}

	phanTramSP := 0
	if maxSanPham > 0 {
		phanTramSP = (soSanPhamHienTai * 100) / maxSanPham
		if phanTramSP > 100 { phanTramSP = 100 }
	}

	c.HTML(http.StatusOK, "tong_quan_admin", gin.H{
		"TieuDe":           "Dashboard Cửa Hàng",
		"NhanVien":         kh,
		"TenGoiHienTai":    tenGoiHienTai,
		"MaxSanPham":       maxSanPham,
		"MaxNhanVien":      maxNhanVien,
		"SoSanPhamHienTai": soSanPhamHienTai,
		"PhanTramSP":       phanTramSP,
	})
}
