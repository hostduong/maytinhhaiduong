package chuc_nang_admin

import (
	"net/http"

	"app/core"
	"github.com/gin-gonic/gin"
)

// TrangNhapHang hiển thị giao diện tạo phiếu nhập (Master-Detail)
func TrangNhapHang(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// 1. Kiểm tra đăng nhập
	kh, found := core.LayKhachHang(shopID, userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 2. Lấy danh sách Dữ liệu Lõi từ RAM Cache để mớm xuống giao diện (JS)
	listSP := core.LayDanhSachSanPhamMayTinh(shopID)
	listNCC := core.LayDanhSachNhaCungCap(shopID)

	// Lọc danh sách SP: Chỉ lấy những phiên bản đang được Bán (TrangThai == 1)
	var dsSanPhamGoiY []*core.SanPhamMayTinh
	for _, sp := range listSP {
		if sp.TrangThai == 1 {
			dsSanPhamGoiY = append(dsSanPhamGoiY, sp)
		}
	}

	// 3. Đẩy ra View
	c.HTML(http.StatusOK, "admin_nhap_hang", gin.H{
		"TieuDe":       "Nhập Hàng",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,
		
		// Đổ data xuống Data Island cho JS xử lý
		"DanhSachSP":   dsSanPhamGoiY,
		"DanhSachNCC":  listNCC,
	})
}
