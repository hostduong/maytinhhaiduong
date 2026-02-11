package chuc_nang

import (
	"net/http"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// Hiển thị trang
func TrangQuanLyDanhMuc(c *gin.Context) {
	userID := c.GetString("USER_ID")
	kh, _ := core.LayKhachHang(userID)

	c.HTML(http.StatusOK, "quan_tri_danh_muc", gin.H{
		"TieuDe":         "Cấu hình Danh mục & Thương hiệu",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"ListDanhMuc":    core.LayDanhSachDanhMuc(),
		"ListThuongHieu": core.LayDanhSachThuongHieu(),
	})
}

// XÓA DANH MỤC BẰNG PIN
func API_XoaDanhMuc(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	maDM := c.PostForm("ma_danh_muc")
	maPin := c.PostForm("ma_pin")

	kh, _ := core.LayKhachHang(c.GetString("USER_ID"))
	if !bao_mat.KiemTraMatKhau(maPin, kh.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}

	dm, ok := core.LayChiTietDanhMuc(maDM)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy!"})
		return
	}

	// Đổi Mã thành rỗng để ẩn khỏi hệ thống (Soft delete cho cấu hình)
	core.KhoaHeThong.Lock()
	dm.MaDanhMuc = ""
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "DANH_MUC", dm.DongTrongSheet, core.CotDM_MaDanhMuc, "")

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã xóa danh mục!"})
}

// (Tương tự, bạn có thể tự viết thêm API_XoaThuongHieu theo form trên)
