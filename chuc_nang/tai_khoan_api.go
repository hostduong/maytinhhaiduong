package chuc_nang

import (
	"fmt"
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API_DoiThongTin : Cập nhật Full thông tin cá nhân
func API_DoiThongTin(c *gin.Context) {
	// 1. Lấy dữ liệu từ Form
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	// Các trường mới
	diaChiMoi   := strings.TrimSpace(c.PostForm("dia_chi"))
	maSoThueMoi := strings.TrimSpace(c.PostForm("ma_so_thue"))
	zaloMoi     := strings.TrimSpace(c.PostForm("zalo"))
	fbMoi       := strings.TrimSpace(c.PostForm("url_fb"))
	tiktokMoi   := strings.TrimSpace(c.PostForm("url_tiktok"))

	// 2. Validate cơ bản
	if !bao_mat.KiemTraHoTen(hoTenMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ!"})
		return 
	}
	
	// 3. Tìm user trong session
	cookie, _ := c.Cookie("session_id")
	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		// Cập nhật RAM
		kh.TenKhachHang = hoTenMoi
		kh.DienThoai = sdtMoi
		kh.NgaySinh = ngaySinhMoi
		kh.GioiTinh = gioiTinhMoi
		kh.DiaChi = diaChiMoi
		kh.MaSoThue = maSoThueMoi
		kh.Zalo = zaloMoi
		kh.UrlFb = fbMoi
		kh.UrlTiktok = tiktokMoi

		// Đẩy vào hàng chờ ghi xuống Sheet
		sID := cau_hinh.BienCauHinh.IdFileSheet
		row := kh.DongTrongSheet
		sheet := "KHACH_HANG"

		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_TenKhachHang, hoTenMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_DienThoai, sdtMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_NgaySinh, ngaySinhMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_GioiTinh, gioiTinhMoi)
		
		// Ghi các cột mới (Đảm bảo file chi_muc.go đã có const cho các cột này)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_DiaChi, diaChiMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_MaSoThue, maSoThueMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_Zalo, zaloMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_UrlFb, fbMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_UrlTiktok, tiktokMoi)

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"}) 
	}
}

// API_DoiMatKhau : (Giữ nguyên logic cũ)
func API_DoiMatKhau(c *gin.Context) {
	passCu := strings.TrimSpace(c.PostForm("pass_cu"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không đúng quy tắc!"})
		return 
	}
	
	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		if !bao_mat.KiemTraMatKhau(passCu, kh.MatKhauHash) { 
			c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu cũ không đúng!"})
			return 
		}
		hash, _ := bao_mat.HashMatKhau(passMoi)
		kh.MatKhauHash = hash
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hash)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) 
	}
}

// API_DoiMaPin : (Đã hash)
func API_DoiMaPin(c *gin.Context) {
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !bao_mat.KiemTraMaPin(pinMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "PIN phải đủ 8 số!"})
		return 
	}
	
	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		if !bao_mat.KiemTraMatKhau(pinCu, kh.MaPinHash) {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN hiện tại không đúng!"})
			return
		}
		hashMoi, _ := bao_mat.HashMatKhau(pinMoi)
		kh.MaPinHash = hashMoi
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, hashMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}) 
	}
}


// API_GuiOTPPin : Gửi mã PIN mới vào Email
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	// Check Rate Limit (Logic mới 1p/lần)
	theGui, msg := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	// Tạo PIN mới
	newPinRaw := nghiep_vu.TaoMaOTP()
	
	// [CẬP NHẬT BODY EMAIL THEO YÊU CẦU]
	body := fmt.Sprintf(`Xin chào,

Chúng tôi đã tạo mã PIN mới cho tài khoản %s theo yêu cầu của bạn trên hệ thống.

Mã PIN mới của bạn là: %s

Vì lý do bảo mật, vui lòng đổi mã PIN này ngay sau khi đăng nhập.

Nếu bạn không yêu cầu thay đổi mã PIN, bạn hãy thay đổi thông tin ngay lập tức.

Trân trọng,
Đội ngũ hỗ trợ`, kh.Email, newPinRaw)

	// Gửi mail
	err := nghiep_vu.GuiMailThongBaoAPI(kh.Email, "Thông báo thay đổi mã PIN", "Hỗ trợ tài khoản", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	// Lưu PIN mới (đã hash)
	hashNewPin, _ := bao_mat.HashMatKhau(newPinRaw)
	kh.MaPinHash = hashNewPin
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, hashNewPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã PIN mới vào Email!"})
}
