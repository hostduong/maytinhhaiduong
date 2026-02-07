package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func TrangDangKy(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "dang_ky", gin.H{})
}

func XuLyDangKy(c *gin.Context) {
	// Lấy dữ liệu và chuẩn hóa
	hoTen     := strings.TrimSpace(c.PostForm("ho_ten"))
	
	// User & Email phải ép về chữ thường
	user      := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email     := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	
	pass      := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin     := strings.TrimSpace(c.PostForm("ma_pin"))
	
	// Xử lý số điện thoại từ Intl-Input
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")) 
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }
	
	ngaySinh  := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinh  := strings.TrimSpace(c.PostForm("gioi_tinh"))

	// 1. VALIDATE SERVER-SIDE (Chốt chặn cuối cùng)
	// Các hàm KiemTra... trong bao_mat đã được cập nhật logic mới nhất
	if !bao_mat.KiemTraHoTen(hoTen) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Họ tên không hợp lệ!"})
		return
	}
	if !bao_mat.KiemTraTenDangNhap(user) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập không đúng quy tắc!"})
		return
	}
	if !bao_mat.KiemTraEmail(email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Email không hợp lệ!"})
		return
	}
	if !bao_mat.KiemTraMaPin(maPin) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Mã PIN phải đúng 8 số!"})
		return
	}
	if !bao_mat.KiemTraDinhDangMatKhau(pass) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Mật khẩu chứa ký tự không cho phép!"})
		return
	}

	// 2. Kiểm tra trùng lặp
	if nghiep_vu.KiemTraTonTaiUserEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	// 3. Logic tạo tài khoản (Admin đầu tiên hoặc Khách)
	var maKH, vaiTro, loaiKH string
	if nghiep_vu.DemSoLuongKhachHang() == 0 {
		maKH = "KH_0001"
		vaiTro = "admin"
		loaiKH = "quan_tri_vien"
	} else {
		maKH = nghiep_vu.TaoMaKhachHangMoi()
		vaiTro = "" 
		loaiKH = "khach_le"
	}

	// 4. Mã hóa và Lưu trữ
	passHash, _ := bao_mat.HashMatKhau(pass)
	
	// Lưu ý: Mã PIN cũng nên hash nếu muốn bảo mật cao, 
	// nhưng ở đây ta giữ nguyên logic truyền maPin (đã hash trong logic_khach_hang nếu có)
	// Dựa vào logic cũ của bạn: ThemKhachHangMoi sẽ tự hash PIN.
	
	cookie := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	newKH := &mo_hinh.KhachHang{
		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		DienThoai:      dienThoai,
		MatKhauHash:    passHash, // Hash pass
		MaPinHash:      maPin,    // Truyền PIN thô, hàm dưới sẽ hash
		TenKhachHang:   hoTen,
		NgaySinh:       ngaySinh,
		GioiTinh:       gioiTinh,
		LoaiKhachHang:  loaiKH,
		VaiTroQuyenHan: vaiTro,
		Cookie:         cookie,
		CookieExpired:  expiredTime,
		TrangThai:      1,
		NgayTao:        time.Now().Format("2006-01-02 15:04:05"),
	}

	// Hàm này sẽ Hash mã PIN và ghi vào Sheet
	nghiep_vu.ThemKhachHangMoi(newKH)
	
	// Set Cookie để auto login
	c.SetCookie("session_id", cookie, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	if vaiTro == "admin" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
