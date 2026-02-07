package chuc_nang

import (
	"net/http"
	"app/nghiep_vu"
	"github.com/gin-gonic/gin"
)

// Hàm hỗ trợ lấy thông tin User từ Cookie (Dùng nội bộ file này)
func layThongTinNguoiDung(c *gin.Context) (bool, string, string) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		// Tìm Khách Hàng thay vì Nhân Viên
		if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			// Trả về: Đã đăng nhập, Tên hiển thị (Họ tên), Quyền hạn
			return true, kh.TenKhachHang, kh.VaiTroQuyenHan
		}
	}
	return false, "", ""
}

// TrangChu : Hiển thị trang chủ
func TrangChu(c *gin.Context) {
	// 1. Lấy dữ liệu sản phẩm
	danhSachSP := nghiep_vu.LayDanhSachSanPham()
	
	// 2. Kiểm tra đăng nhập
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)

	// 3. Trả về HTML kèm thông tin User
	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":          "Trang Chủ",
		"DanhSachSanPham": danhSachSP,
		"DaDangNhap":      daLogin,   // Biến cờ để giao diện biết
		"TenNguoiDung":    tenUser,   // Tên để hiển thị "Chào A"
		"QuyenHan":        quyen,     // Để hiện nút Admin nếu cần
	})
}

// ChiTietSanPham : Hiển thị trang chi tiết
func ChiTietSanPham(c *gin.Context) {
	id := c.Param("id")
	sp, tonTai := nghiep_vu.LayChiTietSanPham(id)

	if !tonTai {
		c.String(http.StatusNotFound, "Không tìm thấy sản phẩm này!")
		return
	}

	// Kiểm tra đăng nhập
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)

	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":       sp.TenSanPham,
		"SanPham":      sp,
		"DaDangNhap":   daLogin,
		"TenNguoiDung": tenUser,
		"QuyenHan":     quyen,
	})
}
