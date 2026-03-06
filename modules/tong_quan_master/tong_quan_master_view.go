package tong_quan

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangTongQuanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // Đây sẽ là ID của Master
	userID := c.GetString("USER_ID")

	// 1. Kiểm tra quyền truy cập (Dù Middleware đã chặn, nhưng cứ check cho chắc)
	nhanVien, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 2. Tính toán các chỉ số thống kê từ RAM
	tongSoShop := 0
	tongSoGoiDaBan := 0
	doanhThuTamTinh := 0.0

	// Lấy danh sách toàn bộ khách hàng (Các Shop đang thuê SaaS)
	core.KhoaHeThong.RLock()
	danhSachKhach := core.CacheKhachHang[shopID]
	core.KhoaHeThong.RUnlock()

	for _, kh := range danhSachKhach {
		// Chỉ đếm những shop đã kích hoạt gói Starter trở lên
		hasActivePlan := false
		for _, goi := range kh.GoiDichVu {
			if goi.TrangThai == "active" {
				hasActivePlan = true
				tongSoGoiDaBan++
				// (Giả lập: Nếu lưu giá bán vào PlanInfo thì có thể cộng dồn doanh thu ở đây)
			}
		}
		if hasActivePlan {
			tongSoShop++
		}
	}

	// 3. Render ra giao diện HTML
	c.HTML(http.StatusOK, "tong_quan_master", gin.H{
		"TieuDe":         "Tổng Quan Master",
		"NhanVien":       nhanVien,
		"TongSoShop":     tongSoShop,
		"TongSoGoiDaBan": tongSoGoiDaBan,
		"DoanhThuTamTinh": doanhThuTamTinh, // Hiện tại để 0, sau này móc với bảng Thanh toán
	})
}
