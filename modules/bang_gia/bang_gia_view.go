package bang_gia

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

// GoiDichVuView: Struct bọc thêm dữ liệu hiển thị động
type GoiDichVuView struct {
	*core.GoiDichVu
	TextNutBam string // Chỉ giữ lại Text nút bấm động
}

// ĐÃ XÓA HÀM formatThoiGian VÌ KHÔNG CẦN THIẾT NỮA

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
		if p.TrangThai != 1 { continue }

		if !hasStarter {
			if p.LoaiGoi != "STARTER" { continue }
		} else {
			if p.LoaiGoi == "STARTER" { continue }
		}

		textNut := "ĐĂNG KÝ NGAY"
		if hasStarter { textNut = "MUA THÊM" }

		viewItem := GoiDichVuView{
			GoiDichVu:  p,
			TextNutBam: textNut,
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
