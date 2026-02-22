package chuc_nang

import (
	"fmt"
	"strings"

	"app/cau_hinh"
	"app/core"
	"github.com/gin-gonic/gin"
)

// API_DoiThongTin : Cập nhật Full thông tin cá nhân
func API_DoiThongTin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // Lấy ID của cửa hàng hiện tại
	
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	diaChiMoi   := strings.TrimSpace(c.PostForm("dia_chi"))
	maSoThueMoi := strings.TrimSpace(c.PostForm("ma_so_thue"))
	zaloMoi     := strings.TrimSpace(c.PostForm("zalo"))
	fbMoi       := strings.TrimSpace(c.PostForm("url_fb"))
	tiktokMoi   := strings.TrimSpace(c.PostForm("url_tiktok"))

	if !cau_hinh.KiemTraHoTen(hoTenMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ!"})
		return 
	}
	
	cookie, _ := c.Cookie("session_id")
	if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
		core.KhoaHeThong.Lock()
		kh.TenKhachHang = hoTenMoi
		kh.DienThoai = sdtMoi
		kh.NgaySinh = ngaySinhMoi
		gioiTinh := -1
		if gioiTinhMoi == "Nam" { gioiTinh = 1 } else if gioiTinhMoi == "Nữ" { gioiTinh = 0 }
		kh.GioiTinh = gioiTinh
		kh.DiaChi = diaChiMoi
		kh.MaSoThue = maSoThueMoi
		kh.MangXaHoi.Zalo = zaloMoi
		kh.MangXaHoi.Facebook = fbMoi
		kh.MangXaHoi.Tiktok = tiktokMoi
		core.KhoaHeThong.Unlock()

		row := kh.DongTrongSheet
		sheet := "KHACH_HANG"

		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_TenKhachHang, hoTenMoi)
		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_DienThoai, sdtMoi)
		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_NgaySinh, ngaySinhMoi)
		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_GioiTinh, gioiTinh)
		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_DiaChi, diaChiMoi)
		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_MaSoThue, maSoThueMoi)
		core.ThemVaoHangCho(shopID, sheet, row, core.CotKH_MangXaHoiJson, core.ToJSON(kh.MangXaHoi))

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"}) 
	}
}

// API_DoiMatKhau
func API_DoiMatKhau(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	passCu := strings.TrimSpace(c.PostForm("pass_cu"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !cau_hinh.KiemTraDinhDangMatKhau(passMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không đúng quy tắc!"})
		return 
	}
	
	if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
		if !cau_hinh.KiemTraMatKhau(passCu, kh.MatKhauHash) { 
			c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu cũ không đúng!"})
			return 
		}
		hash, _ := cau_hinh.HashMatKhau(passMoi)
		core.KhoaHeThong.Lock()
		kh.MatKhauHash = hash
		core.KhoaHeThong.Unlock()
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) 
	}
}

// API_DoiMaPin
func API_DoiMaPin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !cau_hinh.KiemTraMaPin(pinMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "PIN phải đủ 8 số!"})
		return 
	}
	
	if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
		if !cau_hinh.KiemTraMatKhau(pinCu, kh.MaPinHash) {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN hiện tại không đúng!"})
			return
		}
		hashMoi, _ := cau_hinh.HashMatKhau(pinMoi)
		core.KhoaHeThong.Lock()
		kh.MaPinHash = hashMoi
		core.KhoaHeThong.Unlock()
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}) 
	}
}

// API_GuiOTPPin : Gửi mã PIN mới vào Email
func API_GuiOTPPin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	cookie, _ := c.Cookie("session_id")
	kh, ok := core.TimKhachHangTheoCookie(shopID, cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	theGui, msg := core.KiemTraRateLimit(kh.Email)
	if !theGui { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	newPinRaw := core.TaoMaOTP()
	
	body := fmt.Sprintf(`Xin chào,

Chúng tôi đã tạo mã PIN mới cho tài khoản %s theo yêu cầu của bạn trên hệ thống.

Mã PIN mới của bạn là: %s

Vì lý do bảo mật, vui lòng đổi mã PIN này ngay sau khi đăng nhập.

Nếu bạn không yêu cầu thay đổi mã PIN, bạn hãy thay đổi thông tin ngay lập tức.

Trân trọng,
Đội ngũ hỗ trợ`, kh.Email, newPinRaw)

	err := core.GuiMailThongBaoAPI(kh.Email, "Thông báo thay đổi mã PIN", "Hỗ trợ tài khoản", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	hashNewPin, _ := cau_hinh.HashMatKhau(newPinRaw)
	core.KhoaHeThong.Lock()
	kh.MaPinHash = hashNewPin
	core.KhoaHeThong.Unlock()
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashNewPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã PIN mới vào Email!"})
}
