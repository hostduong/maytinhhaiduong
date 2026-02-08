package chuc_nang

import (
	"net/http"
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core" // [QUAN TRỌNG] Sử dụng Core

	"github.com/gin-gonic/gin"
)

// API_Admin_SuaThanhVien : Dành cho Admin/Sale sửa thông tin người khác
// Method: POST /admin/api/member/update
func API_Admin_SuaThanhVien(c *gin.Context) {
	// 1. KIỂM TRA QUYỀN HẠN (Logic thay thế nghiep_vu.KiemTraQuyen)
	vaiTro := c.GetString("USER_ROLE")
	choPhep := false
	if vaiTro == "admin_root" || vaiTro == "admin" || vaiTro == "quan_ly" {
		choPhep = true
	}

	if !choPhep {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Bạn không có quyền sửa thành viên!"})
		return
	}

	// 2. LẤY DỮ LIỆU
	maKhachHangCanSua := c.PostForm("ma_khach_hang")
	
	// Tìm khách hàng trong RAM Core
	khachHang, tonTai := core.LayKhachHang(maKhachHangCanSua)
	if !tonTai {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "msg": "Không tìm thấy khách hàng này!"})
		return
	}

	// 3. CẬP NHẬT THÔNG TIN
	hoTenMoi := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi   := strings.TrimSpace(c.PostForm("dien_thoai"))
	
	// Lấy ID Sheet chuẩn (Hỗ trợ đa Shop)
	idSheet := khachHang.SpreadsheetID
	if idSheet == "" { idSheet = cau_hinh.BienCauHinh.IdFileSheet }
	row := khachHang.DongTrongSheet

	// Cập nhật RAM Core & Đẩy vào Hàng Chờ Ghi
	if hoTenMoi != "" {
		// Lưu ý: Cần Lock nếu muốn thread-safe tuyệt đối, hoặc gán trực tiếp
		khachHang.TenKhachHang = hoTenMoi
		core.ThemVaoHangCho(idSheet, "KHACH_HANG", row, core.CotKH_TenKhachHang, hoTenMoi)
	}
	if sdtMoi != "" {
		khachHang.DienThoai = sdtMoi
		core.ThemVaoHangCho(idSheet, "KHACH_HANG", row, core.CotKH_DienThoai, sdtMoi)
	}

	// Reset Mật khẩu (Nếu có)
	passMoi := c.PostForm("new_password")
	if passMoi != "" {
		// Chỉ Admin Root hoặc Admin mới được reset pass
		if vaiTro == "admin_root" || vaiTro == "admin" {
			hash, _ := bao_mat.HashMatKhau(passMoi)
			khachHang.MatKhauHash = hash
			core.ThemVaoHangCho(idSheet, "KHACH_HANG", row, core.CotKH_MatKhauHash, hash)
		} else {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Bạn không có quyền reset mật khẩu!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Cập nhật thành công!"})
}
