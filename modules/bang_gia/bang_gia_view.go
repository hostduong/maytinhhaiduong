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
}

// Hàm quy đổi thời gian thông minh
func formatThoiGian(ngay int) string {
	if ngay == 9999 {
		return "vĩnh viễn" // 9999 = vĩnh viễn
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

	for _, p := range allPackages {
		if p.TrangThai != 1 {
			continue 
		}

		if !hasStarter {
			if p.LoaiGoi != "STARTER" {
				continue
			}
		} else {
			if p.LoaiGoi == "STARTER" {
				continue
			}
		}

		textNut := "ĐĂNG KÝ NGAY"
		if hasStarter {
			textNut = "MUA THÊM"
		}

		viewItem := GoiDichVuView{
			GoiDichVu:       p,
			ThoiGianHienThi: formatThoiGian(p.ThoiHanNgay),
			TextNutBam:      textNut,
		}
		listGoiView = append(listGoiView, viewItem)
	}

	tieuDeTrang := "Bảng Giá Dịch Vụ"
	moTaTrang := "CHỌN GÓI CƯỚC PHÙ HỢP VỚI QUY MÔ CỬA HÀNG"
	if hasStarter {
		tieuDeTrang = "Nâng Cấp Hệ Thống"
		moTaTrang = "MUA THÊM DUNG LƯỢNG VÀ TÀI KHOẢN NHÂN VIÊN"
	}

	c.HTML(http.StatusOK, "bang_gia", gin.H{
		"TieuDe":     tieuDeTrang,
		"MoTaTrang":  moTaTrang,
		"ListGoi":    listGoiView,
		"HasStarter": hasStarter,
	})
}
