package chuc_nang

import (
	"net/http"

	"app/core"
	"github.com/gin-gonic/gin"
)

// TrangNhapHangMaster hiển thị giao diện tạo phiếu nhập cho Super Admin (Master)
func TrangNhapHangMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, found := core.LayKhachHang(shopID, userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Lấy danh sách Dữ liệu Lõi từ RAM Cache để mớm xuống giao diện (JS)
	listSP := core.LayDanhSachSanPhamMayTinh(shopID)
	listNCC := core.LayDanhSachNhaCungCap(shopID)

	// Lọc danh sách SP: Chỉ lấy những phiên bản đang được Bán (TrangThai == 1)
	var dsSanPhamGoiY []*core.SanPhamMayTinh
	for _, sp := range listSP {
		if sp.TrangThai == 1 {
			dsSanPhamGoiY = append(dsSanPhamGoiY, sp)
		}
	}

	// Trỏ đúng vào template master_nhap_hang
	c.HTML(http.StatusOK, "master_nhap_hang", gin.H{
		"TieuDe":       "Nhập Hàng",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,
		
		"DanhSachSP":   dsSanPhamGoiY,
		"DanhSachNCC":  listNCC,
	})
}
