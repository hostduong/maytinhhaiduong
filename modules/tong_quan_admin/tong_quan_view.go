package tong_quan

import (
	"net/http"

	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangTongQuanCuaHang(c *gin.Context) {
	// Ở không gian Tenant Admin, SHOP_ID chính là ID của file Master chứa dữ liệu Khách hàng
	masterShopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")

	// Lấy thông tin Chủ shop
	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Tính toán hạn mức tài nguyên (Cộng dồn từ các gói cước đang Active)
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

	// Lấy số lượng sản phẩm thực tế đang có trong kho của Cửa hàng (Đọc từ RAM riêng của shop)
	soSanPhamHienTai := 0
	if kh.DataSheets.SpreadsheetID != "" {
		core.KhoaHeThong.RLock()
		if ds, exists := core.CacheSanPhamMayTinh[kh.DataSheets.SpreadsheetID]; exists {
			soSanPhamHienTai = len(ds)
		}
		core.KhoaHeThong.RUnlock()
	}

	// Tính % sử dụng tài nguyên để vẽ thanh Progress Bar
	phanTramSP := 0
	if maxSanPham > 0 {
		phanTramSP = (soSanPhamHienTai * 100) / maxSanPham
		if phanTramSP > 100 { phanTramSP = 100 }
	}

	c.HTML(http.StatusOK, "tong_quan_cua_hang", gin.H{
		"TieuDe":           "Dashboard Cửa Hàng",
		"NhanVien":         kh,
		"TenGoiHienTai":    tenGoiHienTai,
		"MaxSanPham":       maxSanPham,
		"MaxNhanVien":      maxNhanVien,
		"SoSanPhamHienTai": soSanPhamHienTai,
		"PhanTramSP":       phanTramSP,
	})
}
