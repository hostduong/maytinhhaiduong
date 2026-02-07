package chuc_nang

import (
	"net/http"
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API_Admin_SuaThanhVien : Dành cho Admin/Sale sửa thông tin người khác
// Method: POST /admin/api/member/update
func API_Admin_SuaThanhVien(c *gin.Context) {
	// 1. KIỂM TRA QUYỀN HẠN
	vaiTroNguoiSua := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTroNguoiSua, "member.edit") {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Bạn không có quyền sửa thành viên!"})
		return
	}

	// 2. LẤY DỮ LIỆU
	maKhachHangCanSua := c.PostForm("ma_khach_hang") // ID của người bị sửa
	
	// Tìm khách hàng trong RAM
	khachHang, tonTai := nghiep_vu.LayThongTinKhachHang(maKhachHangCanSua)
	if !tonTai {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "msg": "Không tìm thấy khách hàng này!"})
		return
	}

	// 3. CẬP NHẬT THÔNG TIN (Chỉ cập nhật những gì gửi lên)
	hoTenMoi := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi   := strings.TrimSpace(c.PostForm("dien_thoai"))
	
	// [ĐÃ XÓA BIẾN EMAIL THỪA TẠI ĐÂY ĐỂ FIX LỖI BUILD]
	// Logic sửa Email phức tạp hơn vì dính đến Key Map, tạm thời chưa cho sửa Email ở đây
	
	idSheet := cau_hinh.BienCauHinh.IdFileSheet
	row := khachHang.DongTrongSheet

	// Cập nhật RAM & Đẩy vào Hàng Chờ Ghi
	if hoTenMoi != "" {
		khachHang.TenKhachHang = hoTenMoi
		nghiep_vu.ThemVaoHangCho(idSheet, "KHACH_HANG", row, mo_hinh.CotKH_TenKhachHang, hoTenMoi)
	}
	if sdtMoi != "" {
		khachHang.DienThoai = sdtMoi
		nghiep_vu.ThemVaoHangCho(idSheet, "KHACH_HANG", row, mo_hinh.CotKH_DienThoai, sdtMoi)
	}

	// Reset Mật khẩu (Nếu có)
	passMoi := c.PostForm("new_password")
	if passMoi != "" {
		// Chỉ Admin mới được reset pass (Ví dụ thêm logic check sâu hơn)
		// Hoặc dùng mã quyền riêng: member.reset_pass
		if nghiep_vu.KiemTraQuyen(vaiTroNguoiSua, "member.reset_pass") {
			hash, _ := bao_mat.HashMatKhau(passMoi)
			khachHang.MatKhauHash = hash
			nghiep_vu.ThemVaoHangCho(idSheet, "KHACH_HANG", row, mo_hinh.CotKH_MatKhauHash, hash)
		} else {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "Bạn không có quyền reset mật khẩu!"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "msg": "Cập nhật thành công!"})
}
