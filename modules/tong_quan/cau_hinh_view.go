package cau_hinh

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangCaiDatCauHinhMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// Rút dữ liệu hàng loạt từ RAM (Tốc độ mili-giây)
	core.KhoaHeThong.RLock()
	kh := core.CacheMapKhachHang[shopID+"__"+userID]
	listDM := core.CacheDanhMuc[shopID]
	listTH := core.CacheThuongHieu[shopID]
	listBLN := core.CacheBienLoiNhuan[shopID]
	listNCC := core.CacheNhaCungCap[shopID]
	core.KhoaHeThong.RUnlock()

	if kh == nil {
		c.String(http.StatusUnauthorized, "Vui lòng đăng nhập")
		return
	}

	c.HTML(http.StatusOK, "master_cai_dat_cau_hinh", gin.H{
		"TieuDe":       "Cấu Hình Hệ Thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,
		
		"ListDanhMuc":    listDM,
		"ListThuongHieu": listTH,
		"ListBLN":        listBLN,
		"ListNCC":        listNCC,
	})
}
