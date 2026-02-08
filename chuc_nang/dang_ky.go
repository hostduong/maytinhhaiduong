package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core" // [QUAN TRỌNG] Sử dụng Core

	"github.com/gin-gonic/gin"
)

func TrangDangKy(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		// Kiểm tra bằng Core
		if _, ok := core.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "dang_ky", gin.H{"TieuDe": "Đăng Ký Tài Khoản"})
}

func XuLyDangKy(c *gin.Context) {
	// Lấy dữ liệu và chuẩn hóa
	hoTen     := strings.TrimSpace(c.PostForm("ho_ten"))
	
	// User & Email phải ép về chữ thường
	user      := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email     := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	
	pass      := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin     := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")) 
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }
	
	ngaySinh  := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinh  := strings.TrimSpace(c.PostForm("gioi_tinh"))

	// 1. VALIDATE SERVER-SIDE (Giữ nguyên logic bảo mật)
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

	// 2. Kiểm tra trùng lặp (Dùng Core)
	if core.KiemTraTonTaiUserEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	// 3. Logic tạo tài khoản (Dùng Core)
	var maKH, vaiTro, loaiKH string
	
	// Đếm số lượng từ Core
	soLuongUser := len(core.LayDanhSachKhachHang())

	if soLuongUser == 0 {
		maKH = "KH_0001"
		vaiTro = "admin_root" // Admin đầu tiên
		loaiKH = "Quản trị viên"
	} else {
		maKH = core.TaoMaKhachHangMoi()
		vaiTro = "customer" 
		loaiKH = "Khách lẻ"
	}

	// 4. Mã hóa
	passHash, _ := bao_mat.HashMatKhau(pass)
	// Hash luôn mã PIN để lưu an toàn
	pinHash, _ := bao_mat.HashMatKhau(maPin)
	
	cookie := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// 5. Tạo Struct Core
	sID := cau_hinh.BienCauHinh.IdFileSheet
	
	// Tính dòng mới = Dòng bắt đầu + Số lượng hiện có
	newRow := core.DongBatDauDuLieu + soLuongUser

	newKH := &core.KhachHang{
		SpreadsheetID:  sID,
		DongTrongSheet: newRow,

		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		DienThoai:      dienThoai,
		MatKhauHash:    passHash,
		MaPinHash:      pinHash, // Lưu đã hash
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

	// 6. Lưu vào RAM Core
	core.ThemKhachHangVaoRam(newKH)
	
	// 7. Đẩy xuống Hàng Chờ Ghi (Dùng Core Queue)
	ghi := core.ThemVaoHangCho
	sheet := "KHACH_HANG"
	
	ghi(sID, sheet, newRow, core.CotKH_MaKhachHang, newKH.MaKhachHang)
	ghi(sID, sheet, newRow, core.CotKH_TenDangNhap, newKH.TenDangNhap)
	ghi(sID, sheet, newRow, core.CotKH_MatKhauHash, newKH.MatKhauHash)
	ghi(sID, sheet, newRow, core.CotKH_MaPinHash, newKH.MaPinHash)
	ghi(sID, sheet, newRow, core.CotKH_Email, newKH.Email)
	ghi(sID, sheet, newRow, core.CotKH_DienThoai, newKH.DienThoai)
	ghi(sID, sheet, newRow, core.CotKH_TenKhachHang, newKH.TenKhachHang)
	ghi(sID, sheet, newRow, core.CotKH_NgaySinh, newKH.NgaySinh)
	ghi(sID, sheet, newRow, core.CotKH_GioiTinh, newKH.GioiTinh)
	ghi(sID, sheet, newRow, core.CotKH_LoaiKhachHang, newKH.LoaiKhachHang)
	ghi(sID, sheet, newRow, core.CotKH_VaiTroQuyenHan, newKH.VaiTroQuyenHan)
	ghi(sID, sheet, newRow, core.CotKH_TrangThai, newKH.TrangThai)
	ghi(sID, sheet, newRow, core.CotKH_Cookie, newKH.Cookie)
	ghi(sID, sheet, newRow, core.CotKH_CookieExpired, newKH.CookieExpired)
	ghi(sID, sheet, newRow, core.CotKH_NgayTao, newKH.NgayTao)
	
	// 8. Set Cookie và Redirect
	c.SetCookie("session_id", cookie, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	if vaiTro == "admin_root" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
