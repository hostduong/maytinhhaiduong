package chuc_nang

import (
	"net/http"
	"strings"
	
	"app/core"
	"github.com/gin-gonic/gin"
)

func API_LayDanhSachSanPham(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 
	danhSach := core.LayDanhSachSanPhamMayTinh(shopID)
	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"so_luong":   len(danhSach),
		"du_lieu":    danhSach,
	})
}

func API_LayMenu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") 
	dsSP := core.LayDanhSachSanPhamMayTinh(shopID)
	uniqueDM := make(map[string]bool)

	for _, sp := range dsSP {
		if sp != nil && sp.MaDanhMuc != "" {
			parts := strings.Split(sp.MaDanhMuc, "|")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" { uniqueDM[p] = true }
			}
		}
	}

	var menu []map[string]string
	for dm := range uniqueDM {
		menu = append(menu, map[string]string{
			"ten_danh_muc": dm,
			"slug":         strings.ReplaceAll(strings.ToLower(dm), " ", "-"),
		})
	}
	c.JSON(http.StatusOK, gin.H{"danh_muc": menu, "cau_hinh": map[string]interface{}{}})
}

func API_ChiTietSanPham(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	id := c.Param("id")
	
	sp, tonTai := core.LayChiTietSKUMayTinh(shopID, id)
	if !tonTai {
		c.JSON(http.StatusNotFound, gin.H{"trang_thai": "loi", "thong_bao": "Không tìm thấy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trang_thai": "thanh_cong", "du_lieu": sp})
}
