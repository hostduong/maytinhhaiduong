package trang_chu_admin

import (
	"fmt"
	"net/http"

	"app/core"
	"github.com/gin-gonic/gin"
)

// Bọc thêm dữ liệu hiển thị động cho gói cước
type GoiDichVuView struct {
	*core.GoiDichVu
	TextNutBam string 
}

func TrangChuAdmin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // Đây là ID của Master (Nơi chứa Bảng Giá)
	userID := c.GetString("USER_ID") // ID của khách hàng đang đăng nhập

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
	c.HTML(http.StatusOK, "trang_chu_admin", gin.H{
		"TieuDe":     "Tổng Quan Hệ Thống", // Tên hiển thị trên Header PC
		"NhanVien":   kh,                   // Truyền xuống để hiện Avatar trên Sidebar
		"ListGoi":    listGoiView,
		"HasStarter": hasStarter,
	})
}
