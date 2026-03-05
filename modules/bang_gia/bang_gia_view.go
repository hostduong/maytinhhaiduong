package bang_gia

import (
	"fmt"
	"net/http"

	"app/core"
	"github.com/gin-gonic/gin"
)

// GoiDichVuView: Struct bọc thêm dữ liệu hiển thị động
type GoiDichVuView struct {
	*core.GoiDichVu
	ThoiGianHienThi string
	TextNutBam      string
	PhanTramGiam    int // MỚI: Thêm % giảm giá
}

// Hàm quy đổi thời gian thông minh
func formatThoiGian(ngay int) string {
	if ngay == 9999 {
		return "vĩnh viễn"
	}
	if ngay%365 == 0 {
		return fmt.Sprintf("%d năm", ngay/365)
	} else if ngay%30 == 0 {
		return fmt.Sprintf("%d tháng", ngay/30)
	}
	return fmt.Sprintf("%d ngày", ngay)
}

func TrangCongPortalBangGia(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// 1. Kiểm tra khách hàng đã có gói nền tảng (STARTER) chưa
	hasStarter := false
	if kh, ok := core.LayKhachHang(shopID, userID); ok {
		for _, p := range kh.GoiDichVu {
			if p.LoaiGoi == "STARTER" && p.TrangThai == "active" {
				hasStarter = true
				break
			}
		}
	}

	core.KhoaHeThong.RLock()
	allPackages := core.CacheGoiDichVu[shopID]
	core.KhoaHeThong.RUnlock()

	var listGoiView []GoiDichVuView

	// 2. Phân loại và gắp gói cước đẩy ra View
	for _, p := range allPackages {
		if p.TrangThai != 1 {
			continue // Bỏ qua các gói đang bị ẩn
		}

		if !hasStarter {
			// KHÁCH MỚI
			if p.LoaiGoi != "STARTER" {
				continue
			}
		} else {
			// KHÁCH CŨ (MUA THÊM)
			if p.LoaiGoi == "STARTER" {
				continue
			}
		}

		textNut := "ĐĂNG KÝ NGAY"
		if hasStarter {
			textNut = "MUA THÊM"
		}

		// Tự tính % giảm giá (Đoạn này phải nằm TRONG vòng lặp)
		phanTram := 0
		if p.GiaNiemYet > p.GiaBan && p.GiaNiemYet > 0 {
			phanTram = int((p.GiaNiemYet - p.GiaBan) / p.GiaNiemYet * 100)
		}

		// Đóng gói dữ liệu
		viewItem := GoiDichVuView{
			GoiDichVu:       p,
			ThoiGianHienThi: formatThoiGian(p.ThoiHanNgay),
			TextNutBam:      textNut,
			PhanTramGiam:    phanTram,
		}
		listGoiView = append(listGoiView, viewItem)
	}

	// Tiêu đề trang động
	tieuDeTrang := "Bảng Giá Dịch Vụ - 99K.VN"
	moTaTrang := "Chọn gói cước phù hợp với quy mô cửa hàng"
	if hasStarter {
		tieuDeTrang = "Nâng Cấp Hệ Thống - 99K.VN"
		moTaTrang = "Mua thêm dung lượng và tài khoản nhân viên"
	}

	c.HTML(http.StatusOK, "bang_gia", gin.H{
		"TieuDe":     tieuDeTrang,
		"MoTaTrang":  moTaTrang,
		"ListGoi":    listGoiView,
		"HasStarter": hasStarter,
	})
}
