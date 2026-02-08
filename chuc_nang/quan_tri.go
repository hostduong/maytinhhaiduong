package chuc_nang

import (
	"net/http"
	"time"

	"app/core" // [MỚI] Chỉ dùng Core

	"github.com/gin-gonic/gin"
)

// Struct dữ liệu hiển thị (Giữ nguyên cấu trúc để View không lỗi)
type DuLieuDashboard struct {
	TongDoanhThu    float64
	DonHangHomNay   int
	TongSanPham     int
	TongKhachHang   int
	// DonHangMoiNhat []mo_hinh.PhieuXuat // Tạm đóng vì chưa có Struct PhieuXuat trong Core
	ChartNhan       []string
	ChartDoanhThu   []float64
}

func TrangTongQuan(c *gin.Context) {
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	// 1. Lấy thông tin Admin từ Core
	kh, _ := core.LayKhachHang(userID)

	// 2. Tính toán thống kê
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

func tinhToanThongKe() DuLieuDashboard {
	var kq DuLieuDashboard

	// Lock để đọc an toàn
	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()

	// 1. Đếm sản phẩm từ Core
	kq.TongSanPham = len(core.LayDanhSachSanPham())
	
	// 2. Đếm khách hàng từ Core
	kq.TongKhachHang = len(core.LayDanhSachKhachHang())

	// 3. Phần Thống Kê Doanh Thu & Đơn Hàng
	// TODO: Mở lại phần này sau khi chuyển đổi module DonHang sang Core
	// Hiện tại để giá trị 0 để hệ thống chạy được mà không cần file cũ.
	kq.TongDoanhThu = 0
	kq.DonHangHomNay = 0
	
	// Giả lập biểu đồ rỗng để giao diện không bị méo
	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		label := t.Format("02/01")
		kq.ChartNhan = append(kq.ChartNhan, label)
		kq.ChartDoanhThu = append(kq.ChartDoanhThu, 0)
	}

	return kq
}

// API Reload dữ liệu (Nút đồng bộ trên Menu)
func API_NapLaiDuLieu(c *gin.Context) {
	// Logic check quyền đơn giản
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}

	// Gọi Core nạp lại dữ liệu (Mặc định ID Config)
	go func() {
		core.HeThongDangBan = true
		core.NapDanhMuc("")
		core.NapThuongHieu("")
		core.NapSanPham("")
		core.NapKhachHang("")
		core.HeThongDangBan = false
	}()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong", 
		"thong_diep": "Đang tiến hành đồng bộ dữ liệu...",
	})
}
