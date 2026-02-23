package chuc_nang_admin

import (
	"net/http"
	"time"
	"app/core"
	"github.com/gin-gonic/gin"
	data_pc "app/core/may_tinh"
)

// Struct hiển thị
type DuLieuDashboard struct {
	TongDoanhThu    float64
	DonHangHomNay   int
	TongSanPham     int
	TongKhachHang   int
	ChartNhan       []string
	ChartDoanhThu   []float64
}

func TrangTongQuan(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	// Lấy khách hàng theo shop
	kh, _ := core.LayKhachHang(shopID, userID)
	
	// Tính thống kê theo shop
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

func API_NapLaiDuLieu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	vaiTro := c.GetString("USER_ROLE")
	
	// Đã cập nhật đúng tên Role của hệ thống mới
	if vaiTro != "quan_tri_vien_he_thong" && vaiTro != "quan_tri_vien" {
		c.JSON(http.StatusOK, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	go func() {
		core.HeThongDangBan = true
		
		// [SAAS] Nạp lại dữ liệu của chính Shop này
		core.NapPhanQuyen(shopID) 
		core.NapDanhMuc(shopID)   
		core.NapThuongHieu(shopID)
		core.NapBienLoiNhuan(shopID)
		data_pc.NapDuLieu(shopID) // Nạp dữ liệu ngành máy tính
		core.NapKhachHang(shopID)
		
		core.HeThongDangBan = false
	}()

	// Đã cập nhật chuẩn Form JSON status/msg
	c.JSON(http.StatusOK, gin.H{
		"status": "ok", 
		"msg": "Đang tiến hành đồng bộ toàn bộ dữ liệu ngầm...",
	})
}

// Hàm thống kê (Cần ShopID)
func tinhToanThongKe(shopID string) DuLieuDashboard {
	var kq DuLieuDashboard

	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()

	// Đếm số lượng trong Shop
	kq.TongSanPham = len(data_pc.LayDanhSachSanPham(shopID))
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
