package chuc_nang

import (
	"net/http"
	"sort"
	"time"

	"app/bo_nho_dem" // [MỚI] Import gói chứa dữ liệu gốc
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

type DuLieuDashboard struct {
	TongDoanhThu    float64
	DonHangHomNay   int
	TongSanPham     int
	TongKhachHang   int
	DonHangMoiNhat  []mo_hinh.PhieuXuat
	ChartNhan       []string
	ChartDoanhThu   []float64
}

func TrangTongQuan(c *gin.Context) {
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	kh, _ := nghiep_vu.LayThongTinKhachHang(userID)

	stats := tinhToanThongKe()

	c.HTML(http.StatusOK, "quan_tri", gin.H{
		"TieuDe":       "Tổng quan hệ thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     vaiTro,
		"ThongKe":      stats,
	})
}

func tinhToanThongKe() DuLieuDashboard {
	var kq DuLieuDashboard

	// [SỬA] Dùng Lock từ bo_nho_dem
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()

	// [SỬA] Truy cập trực tiếp vào bo_nho_dem để đếm
	// 1. Đếm sản phẩm (Dựa trên danh sách)
	kq.TongSanPham = len(bo_nho_dem.CacheSanPham.DanhSach)
	
	// 2. Đếm khách hàng (Dựa trên danh sách mới tạo -> ĐẾM ĐÚNG)
	kq.TongKhachHang = len(bo_nho_dem.CacheKhachHang.DanhSach)

	now := time.Now().Format("2006-01-02")
	mapDoanhThuNgay := make(map[string]float64)
	var listPX []mo_hinh.PhieuXuat

	// [SỬA] Duyệt CachePhieuXuat từ bo_nho_dem
	for _, px := range bo_nho_dem.CachePhieuXuat.DuLieu {
		listPX = append(listPX, px)
		
		if px.TrangThai != "Đã hủy" {
			kq.TongDoanhThu += px.TongTienPhieu
			if len(px.NgayTao) >= 10 {
				ngay := px.NgayTao[:10]
				mapDoanhThuNgay[ngay] += px.TongTienPhieu
			}
		}

		if len(px.NgayTao) >= 10 && px.NgayTao[:10] == now {
			kq.DonHangHomNay++
		}
	}

	sort.Slice(listPX, func(i, j int) bool {
		return listPX[i].NgayTao > listPX[j].NgayTao
	})

	limit := 5
	if len(listPX) < 5 { limit = len(listPX) }
	kq.DonHangMoiNhat = listPX[:limit]

	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		key := t.Format("2006-01-02")
		label := t.Format("02/01")
		
		kq.ChartNhan = append(kq.ChartNhan, label)
		kq.ChartDoanhThu = append(kq.ChartDoanhThu, mapDoanhThuNgay[key])
	}

	return kq
}

func API_NapLaiDuLieu(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTro, "system.reload") {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}

	// [SỬA] Gọi hàm reload từ bo_nho_dem
	bo_nho_dem.LamMoiHeThong()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong", 
		"thong_diep": "Đã đồng bộ dữ liệu mới nhất!",
	})
}
