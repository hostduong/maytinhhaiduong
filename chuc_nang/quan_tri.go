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
	// [SAAS] Lấy Context
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	// Lấy thông tin user (truyền shopID)
	kh, _ := core.LayKhachHang(shopID, userID)
	
	// Tính thống kê cho Shop này
	stats := tinhToanThongKe(shopID)

	c.HTML(http.StatusOK, "quan_tri", gin.H{
		"TieuDe":       "Tổng quan hệ thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": c.GetString("USER_NAME"),
		"QuyenHan":     vaiTro,
		"ThongKe":      stats,
	})
}

// API Reload dữ liệu (Đồng bộ cho Shop hiện tại)
func API_NapLaiDuLieu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}

	go func() {
		core.HeThongDangBan = true
		
		// Nạp lại toàn bộ dữ liệu của Shop này
		core.NapPhanQuyen(shopID) 
		core.NapDanhMuc(shopID)   
		core.NapThuongHieu(shopID)
		core.NapBienLoiNhuan(shopID)
		core.NapSanPham(shopID)
		core.NapKhachHang(shopID)
		
		core.HeThongDangBan = false
	}()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong", 
		"thong_diep": "Đang tiến hành đồng bộ dữ liệu từ Sheet...",
	})
}

// Hàm tính toán thống kê (Theo Shop)
func tinhToanThongKe(shopID string) DuLieuDashboard {
	var kq DuLieuDashboard

	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()

	// Lấy danh sách của Shop ID
	kq.TongSanPham = len(core.LayDanhSachSanPham(shopID))
	kq.TongKhachHang = len(core.LayDanhSachKhachHang(shopID))
	
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
