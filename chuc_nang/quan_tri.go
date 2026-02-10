package chuc_nang

import (
	"net/http"
	"time"

	"app/core"

	"github.com/gin-gonic/gin"
)

// Struct dữ liệu hiển thị
type DuLieuDashboard struct {
	TongDoanhThu    float64
	DonHangHomNay   int
	TongSanPham     int
	TongKhachHang   int
	ChartNhan       []string
	ChartDoanhThu   []float64
}

func TrangTongQuan(c *gin.Context) {
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	kh, _ := core.LayKhachHang(userID)
	stats := tinhToanThongKe()

	c.HTML(http.StatusOK, "quan_tri", gin.H{
		"TieuDe":       "Tổng quan hệ thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": c.GetString("USER_NAME"),
		"QuyenHan":     vaiTro,
		"ThongKe":      stats,
	})
}

package chuc_nang

import (
	"net/http"
	"time"

	"app/core"

	"github.com/gin-gonic/gin"
)

// API Reload dữ liệu
func API_NapLaiDuLieu(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}

	// [CẬP NHẬT] Nạp thêm Danh Mục và Thương Hiệu
	go func() {
		core.HeThongDangBan = true
		
		core.NapPhanQuyen("") // Nạp lại quyền (nếu cần)
		core.NapDanhMuc("")   // [MỚI]
		core.NapThuongHieu("")// [MỚI]
		core.NapSanPham("")
		core.NapKhachHang("")
		
		core.HeThongDangBan = false
	}()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong", 
		"thong_diep": "Đang tiến hành đồng bộ toàn bộ dữ liệu (Danh mục, Thương hiệu, Sản phẩm, Khách hàng)...",
	})
}

func tinhToanThongKe() DuLieuDashboard {
	var kq DuLieuDashboard

	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()

	kq.TongSanPham = len(core.LayDanhSachSanPham())
	kq.TongKhachHang = len(core.LayDanhSachKhachHang())
	kq.TongDoanhThu = 0
	kq.DonHangHomNay = 0
	
	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		label := t.Format("02/01")
		kq.ChartNhan = append(kq.ChartNhan, label)
		kq.ChartDoanhThu = append(kq.ChartDoanhThu, 0)
	}

	return kq
}

// API Reload dữ liệu
func API_NapLaiDuLieu(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}

	// [ĐÃ SỬA] Chỉ nạp lại Sản phẩm và Khách hàng (Bỏ Danh mục/Thương hiệu cũ)
	go func() {
		core.HeThongDangBan = true
		core.NapSanPham("")
		core.NapKhachHang("")
		core.HeThongDangBan = false
	}()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong", 
		"thong_diep": "Đang tiến hành đồng bộ dữ liệu...",
	})
}
