package chuc_nang

import (
	"net/http"
	"app/core" // [MỚI] Dùng Core

	"github.com/gin-gonic/gin"
)

// API_LayDanhSachSanPham : Trả về JSON list sản phẩm
// Method: GET /api/san-pham
func API_LayDanhSachSanPham(c *gin.Context) {
	// Lấy từ Core (Đã lọc theo Shop hiện tại)
	danhSach := core.LayDanhSachSanPham()

	// Trả về JSON thành công
	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"so_luong":   len(danhSach),
		"du_lieu":    danhSach,
	})
}

// API_LayMenu : Trả về danh mục
// Method: GET /api/cau-hinh
func API_LayMenu(c *gin.Context) {
	menu := core.LayDanhSachDanhMuc()
	
	// Banner/Cấu hình web chưa chuyển sang Core nên tạm trả về rỗng
	// để không phụ thuộc vào code cũ.
	banner := map[string]interface{}{} 

	c.JSON(http.StatusOK, gin.H{
		"danh_muc": menu,
		"cau_hinh": banner,
	})
}

// API_ChiTietSanPham : Lấy 1 SP cụ thể
// Method: GET /api/san-pham/:id
func API_ChiTietSanPham(c *gin.Context) {
	id := c.Param("id")
	
	// Lấy từ Core
	sp, tonTai := core.LayChiTietSanPham(id)

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
