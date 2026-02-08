package chuc_nang

import (
	"net/http"
	"app/core" // [QUAN TRỌNG] Sử dụng Core
	"github.com/gin-gonic/gin"
)

// Hàm hỗ trợ lấy thông tin User từ Cookie (Dùng Core)
func layThongTinNguoiDung(c *gin.Context) (bool, string, string) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		// Tìm Khách Hàng trong Core
		if kh, ok := core.TimKhachHangTheoCookie(cookie); ok {
			// Trả về: Đã đăng nhập, Tên hiển thị (Họ tên), Quyền hạn
			return true, kh.TenKhachHang, kh.VaiTroQuyenHan
		}
	}
	return false, "", ""
}

// TrangChu : Hiển thị trang chủ
func TrangChu(c *gin.Context) {
	// 1. Lấy dữ liệu sản phẩm từ Core
	danhSachSP := core.LayDanhSachSanPham()
	
	// 2. Kiểm tra đăng nhập
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)

	// 3. Trả về HTML kèm thông tin User
	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":          "Trang Chủ",
		"DanhSachSanPham": danhSachSP, // Core trả về []*SanPham
		"DaDangNhap":      daLogin,
		"TenNguoiDung":    tenUser,
		"QuyenHan":        quyen,
	})
}

// ChiTietSanPham : Hiển thị trang chi tiết
func ChiTietSanPham(c *gin.Context) {
	id := c.Param("id")
	
	// Lấy từ Core
	sp, tonTai := core.LayChiTietSanPham(id)

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
