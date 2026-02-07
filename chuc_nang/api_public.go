package chuc_nang

import (
	"net/http"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API_LayDanhSachSanPham : Trả về JSON list sản phẩm
// Method: GET /api/san-pham
func API_LayDanhSachSanPham(c *gin.Context) {
	// Gọi lớp nghiệp vụ lấy dữ liệu an toàn
	danhSach := nghiep_vu.LayDanhSachSanPham()

	// Trả về JSON thành công
	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"so_luong":   len(danhSach),
		"du_lieu":    danhSach,
	})
}

// API_LayMenu : Trả về danh mục và cấu hình web
// Method: GET /api/cau-hinh
func API_LayMenu(c *gin.Context) {
	menu := nghiep_vu.LayDanhSachDanhMuc()
	banner := nghiep_vu.LayCauHinhWeb()

	c.JSON(http.StatusOK, gin.H{
		"danh_muc": menu,
		"cau_hinh": banner,
	})
}

// API_ChiTietSanPham : Lấy 1 SP cụ thể
// Method: GET /api/san-pham/:id
func API_ChiTietSanPham(c *gin.Context) {
	id := c.Param("id")
	sp, tonTai := nghiep_vu.LayChiTietSanPham(id)

	if !tonTai {
		c.JSON(http.StatusNotFound, gin.H{
			"trang_thai": "loi",
			"thong_bao":  "Không tìm thấy sản phẩm",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"du_lieu":    sp,
	})
}
