package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core" 

	"github.com/gin-gonic/gin"
)

func TrangDangKy(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
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
	user      := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email     := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass      := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin     := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")) 
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }
	
	ngaySinh  := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinh  := strings.TrimSpace(c.PostForm("gioi_tinh"))

	// 1. VALIDATE SERVER-SIDE
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
	if core.KiemTraTonTaiUserEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	// 3. [UPDATE] Logic tạo Người dùng mới (Chuẩn hóa ID, Chức vụ, Loại KH)
	var maKH, vaiTro, chucVu string
	
	// Loại khách hàng luôn là "web" (để biết nguồn khách đến từ website)
	loaiKH := "web" 
	
	soLuongUser := len(core.LayDanhSachKhachHang())

	if soLuongUser == 0 {
		// --- NGƯỜI ĐẦU TIÊN (SUPER ADMIN) ---
		maKH   = "0000000000000000001"
		vaiTro = "admin_root" 
		chucVu = "Quản trị viên"
	} else {
		// --- NGƯỜI THỨ 2 TRỞ ĐI (KHÁCH HÀNG) ---
		// Dùng hàm sinh mã 19 số ngẫu nhiên từ Core
		maKH   = core.TaoMaKhachHangMoi()
		vaiTro = "customer" 
		chucVu = "Khách hàng"
	}

	// 4. Mã hóa
	passHash, _ := bao_mat.HashMatKhau(pass)
	pinHash, _ := bao_mat.HashMatKhau(maPin)
	
	cookie := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// 5. Tạo Struct Core
	sID := cau_hinh.BienCauHinh.IdFileSheet
	newRow := core.DongBatDau_KhachHang + soLuongUser

	newKH := &core.KhachHang{
		SpreadsheetID:  sID,
		DongTrongSheet: newRow,
		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		DienThoai:      dienThoai,
		MatKhauHash:    passHash,
		MaPinHash:      pinHash,
		TenKhachHang:   hoTen,
		NgaySinh:       ngaySinh,
		GioiTinh:       gioiTinh,
		LoaiKhachHang:  loaiKH, // web
		ChucVu:         chucVu, // Quản trị viên hoặc Khách hàng
		VaiTroQuyenHan: vaiTro, // admin_root hoặc customer
		Cookie:         cookie,
		CookieExpired:  expiredTime,
		TrangThai:      1,
		NgayTao:        time.Now().Format("2006-01-02 15:04:05"),
	}

	// 6. Lưu vào RAM & Sheet
	core.ThemKhachHangVaoRam(newKH)
	
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
	
	// [QUAN TRỌNG] Ghi đúng cột Loại KH và Chức vụ
	ghi(sID, sheet, newRow, core.CotKH_LoaiKhachHang, newKH.LoaiKhachHang)
	ghi(sID, sheet, newRow, core.CotKH_ChucVu, newKH.ChucVu)
	ghi(sID, sheet, newRow, core.CotKH_VaiTroQuyenHan, newKH.VaiTroQuyenHan)
	
	ghi(sID, sheet, newRow, core.CotKH_TrangThai, newKH.TrangThai)
	ghi(sID, sheet, newRow, core.CotKH_Cookie, newKH.Cookie)
	ghi(sID, sheet, newRow, core.CotKH_CookieExpired, newKH.CookieExpired)
	ghi(sID, sheet, newRow, core.CotKH_NgayTao, newKH.NgayTao)
	
	// 7. [FIX LỖI COOKIE MISMATCH]
	// Phải tạo chữ ký và set cookie session_sign
	userAgent := c.Request.UserAgent()
	signature := bao_mat.TaoChuKyBaoMat(cookie, userAgent)
	
	maxAge := int(cau_hinh.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_id", cookie, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true) // <-- Dòng này sửa lỗi Mismatch

	// Điều hướng
	if vaiTro == "admin_root" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
