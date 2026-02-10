package chuc_nang

import (
	"net/http"
	"strings"
	"app/core"

	"github.com/gin-gonic/gin"
)

// API_LayDanhSachSanPham
func API_LayDanhSachSanPham(c *gin.Context) {
	danhSach := core.LayDanhSachSanPham()
	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"so_luong":   len(danhSach),
		"du_lieu":    danhSach,
	})
}

// API_LayMenu (TỰ ĐỘNG TẠO TỪ SẢN PHẨM)
func API_LayMenu(c *gin.Context) {
	// Logic mới: Quét toàn bộ sản phẩm để lấy danh sách Category duy nhất
	dsSP := core.LayDanhSachSanPham()
	uniqueDM := make(map[string]bool)

	for _, sp := range dsSP {
		if sp != nil && sp.DanhMuc != "" {
			parts := strings.Split(sp.DanhMuc, "|")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" { uniqueDM[p] = true }
			}
		}
	}

	// Chuyển map thành list struct để Frontend dễ dùng (giả lập cấu trúc cũ)
	var menu []map[string]string
	for dm := range uniqueDM {
		menu = append(menu, map[string]string{
			"ten_danh_muc": dm,
			"slug":         strings.ReplaceAll(strings.ToLower(dm), " ", "-"),
		})
	}

	banner := map[string]interface{}{} 

	c.JSON(http.StatusOK, gin.H{
		"danh_muc": menu, // Trả về list danh mục tự động
		"cau_hinh": banner,
	})
}

// API_ChiTietSanPham
func API_ChiTietSanPham(c *gin.Context) {
	id := c.Param("id")
	sp, tonTai := core.LayChiTietSanPham(id)
	if !tonTai {
		c.JSON(http.StatusNotFound, gin.H{"trang_thai": "loi", "thong_bao": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trang_thai": "thanh_cong", "du_lieu": sp})
}
