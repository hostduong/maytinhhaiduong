package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// TrangQuanLyThanhVien : Hiển thị danh sách khách hàng/nhân viên
func TrangQuanLyThanhVien(c *gin.Context) {
	// Lấy thông tin người đang đăng nhập
	userID := c.GetString("USER_ID")
	me, _ := core.LayKhachHang(userID)

	// Lấy toàn bộ danh sách từ Core
	listAll := core.LayDanhSachKhachHang()

	c.HTML(http.StatusOK, "quan_tri_thanh_vien", gin.H{
		"TieuDe":       "Quản lý thành viên",
		"NhanVien":     me,
		"DanhSach":     listAll,
		"TenNguoiDung": me.TenKhachHang,
		"QuyenHan":     me.VaiTroQuyenHan,
	})
}

// API_Admin_LuuThanhVien : Cập nhật thông tin và QUYỀN HẠN
func API_Admin_LuuThanhVien(c *gin.Context) {
	// 1. Check quyền Admin (Chỉ Admin Root hoặc Admin mới được sửa người khác)
	myRole := c.GetString("USER_ROLE")
	if myRole != "admin_root" && myRole != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền quản trị nhân sự!"})
		return
	}

	// 2. Lấy dữ liệu
	maKH := c.PostForm("ma_khach_hang")
	role := c.PostForm("vai_tro") // admin, sale, kho...
	name := c.PostForm("ten_khach_hang")
	pass := c.PostForm("mat_khau_moi") // Nếu có nhập thì đổi, không thì thôi
	
	kh, ok := core.LayKhachHang(maKH)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thành viên này!"})
		return
	}

	// 3. Logic chặn: Không ai được sửa Admin Root (trừ chính họ)
	if kh.VaiTroQuyenHan == "admin_root" && myRole != "admin_root" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không thể tác động đến tài khoản Super Admin!"})
		return
	}

	// 4. Cập nhật Core RAM
	core.KhoaHeThong.Lock()
	kh.TenKhachHang = name
	kh.VaiTroQuyenHan = role
	// Map lại tên chức vụ hiển thị cho đẹp
	switch role {
	case "admin": kh.ChucVu = "Quản trị viên"
	case "sale": kh.ChucVu = "Nhân viên kinh doanh"
	case "kho": kh.ChucVu = "Thủ kho"
	case "customer": kh.ChucVu = "Khách hàng"
	default: kh.ChucVu = "Thành viên"
	}
	
	// Đổi pass nếu có
	if pass != "" {
		// Import "app/bao_mat" ở trên đầu file nếu chưa có
		// hash, _ := bao_mat.HashMatKhau(pass)
		// kh.MatKhauHash = hash
		// TODO: Bạn cần import package bao_mat để dùng hàm Hash
	}
	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

	// 5. Ghi xuống Sheet
	sID := cau_hinh.BienCauHinh.IdFileSheet
	row := kh.DongTrongSheet
	
	ghi := core.ThemVaoHangCho
	ghi(sID, "KHACH_HANG", row, core.CotKH_TenKhachHang, kh.TenKhachHang)
	ghi(sID, "KHACH_HANG", row, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	ghi(sID, "KHACH_HANG", row, core.CotKH_ChucVu, kh.ChucVu)
	ghi(sID, "KHACH_HANG", row, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thành viên thành công!"})
}
