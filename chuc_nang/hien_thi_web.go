package chuc_nang

import (
	"net/http"
	"app/core" // Sử dụng Core mới

	"github.com/gin-gonic/gin"
)

// Helper: Lấy thông tin User từ Cookie
func layThongTinNguoiDung(c *gin.Context) (bool, string, string) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		// Tìm trong RAM Core
		if kh, ok := core.TimKhachHangTheoCookie(cookie); ok {
			return true, kh.TenKhachHang, kh.VaiTroQuyenHan
		}
	}
	return false, "", ""
}

// 1. Trang Chủ
func TrangChu(c *gin.Context) {
	// Lấy danh sách sản phẩm từ Core
	danhSachSP := core.LayDanhSachSanPham()
	
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)

	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":          "Trang Chủ",
		"DanhSachSanPham": danhSachSP,
		"DaDangNhap":      daLogin,
		"TenNguoiDung":    tenUser,
		"QuyenHan":        quyen,
	})
}

// 2. Chi Tiết Sản Phẩm
func ChiTietSanPham(c *gin.Context) {
	id := c.Param("id")
	sp, tonTai := core.LayChiTietSanPham(id)

	if !tonTai {
		c.String(http.StatusNotFound, "Không tìm thấy sản phẩm này!")
		return
	}

	daLogin, tenUser, quyen := layThongTinNguoiDung(c)

	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":       sp.TenSanPham,
		"SanPham":      sp,
		"DaDangNhap":   daLogin,
		"TenNguoiDung": tenUser,
		"QuyenHan":     quyen,
	})
}

// 3. [MỚI] Trang Hồ Sơ Cá Nhân
func TrangHoSo(c *gin.Context) {
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	// Chưa đăng nhập thì đá về Login
	if !daLogin {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Lấy chi tiết khách hàng để điền vào Form
	cookie, _ := c.Cookie("session_id")
	kh, _ := core.TimKhachHangTheoCookie(cookie)

	c.HTML(http.StatusOK, "ho_so", gin.H{
		"TieuDe":       "Hồ sơ cá nhân",
		"DaDangNhap":   daLogin,
		"TenNguoiDung": tenUser,
		"QuyenHan":     quyen,
		"NhanVien":     kh, // Truyền biến 'NhanVien' khớp với template ho_so.html
	})
}
