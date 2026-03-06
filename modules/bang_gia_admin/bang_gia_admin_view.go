package bang_gia_admin // Đổi tên package

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

type GoiDichVuView struct {
	*core.GoiDichVu
	TextNutBam string 
}

// Đổi tên hàm thành TrangBangGiaAdmin
func TrangBangGiaAdmin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	

	// 1. Kiểm tra khách hàng đã có gói nền tảng (STARTER) chưa
	hasStarter := false
	for _, p := range kh.GoiDichVu {
		if p.LoaiGoi == "STARTER" && p.TrangThai == "active" {
			hasStarter = true
			break
		}
	}

	// 2. Kéo danh sách Gói Cước từ RAM của Master
	core.KhoaHeThong.RLock()
	allPackages := core.CacheGoiDichVu[shopID]
	core.KhoaHeThong.RUnlock()

	var listGoiView []GoiDichVuView

	// 3. Phân luồng: Khách mới xem gói Khởi tạo, Khách cũ xem gói Nâng cấp
	for _, p := range allPackages {
		if p.TrangThai != 1 {
			continue // Bỏ qua gói bị ẩn
		}

		if !hasStarter {
			if p.LoaiGoi != "STARTER" { continue }
		} else {
			if p.LoaiGoi == "STARTER" { continue }
		}

		textNut := "ĐĂNG KÝ NGAY"
		if hasStarter {
			textNut = "MUA THÊM"
		}

		viewItem := GoiDichVuView{
			GoiDichVu:  p,
			TextNutBam: textNut,
		}
		listGoiView = append(listGoiView, viewItem)
	}

	// 4. Bơm dữ liệu ra View (Kết hợp Layout Admin)
    c.HTML(http.StatusOK, "bang_gia_admin", gin.H{
		"TieuDe":     "Bảng Giá Dịch Vụ", 
		"NhanVien":   kh,                   
		"ListGoi":    listGoiView,
		"HasStarter": hasStarter,
	})
}
