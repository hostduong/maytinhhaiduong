package cau_hinh_he_thong

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangCauHinhHeThongView(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// Lấy thông tin user đăng nhập (Để hiện Avatar)
	kh, _ := core.CacheMapKhachHang[core.TaoCompositeKey(shopID, userID)]

	// Gọi Repo lấy dữ liệu các bảng
	listNCC := repo.LayDanhSachNCC(shopID)
	// listDM := repo.LayDanhSachDM(shopID)... (bạn tự bổ sung)

	c.HTML(http.StatusOK, "master_cai_dat_cau_hinh", gin.H{
		"TieuDe":       "Cấu Hình Hệ Thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"ListNCC":      listNCC,
		// "ListDanhMuc": listDM,
	})
}
