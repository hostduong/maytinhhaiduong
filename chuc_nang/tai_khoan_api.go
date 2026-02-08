package chuc_nang

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core" // [MỚI]

	"github.com/gin-gonic/gin"
)

// Helper tạo mã PIN 8 số
func taoMaPIN8So() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

// API_DoiThongTin : Cập nhật Full thông tin cá nhân
func API_DoiThongTin(c *gin.Context) {
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	diaChiMoi   := strings.TrimSpace(c.PostForm("dia_chi"))
	maSoThueMoi := strings.TrimSpace(c.PostForm("ma_so_thue"))
	zaloMoi     := strings.TrimSpace(c.PostForm("zalo"))
	fbMoi       := strings.TrimSpace(c.PostForm("url_fb"))
	tiktokMoi   := strings.TrimSpace(c.PostForm("url_tiktok"))

	if !bao_mat.KiemTraHoTen(hoTenMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ!"})
		return 
	}
	
	cookie, _ := c.Cookie("session_id")
	if kh, ok := core.TimKhachHangTheoCookie(cookie); ok {
		// Cập nhật RAM Core
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
		sID := kh.SpreadsheetID
		if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
		row := kh.DongTrongSheet
		sheet := "KHACH_HANG"

		// Wrapper ngắn gọn
		ghi := core.ThemVaoHangCho
		ghi(sID, sheet, row, core.CotKH_TenKhachHang, hoTenMoi)
		ghi(sID, sheet, row, core.CotKH_DienThoai, sdtMoi)
		ghi(sID, sheet, row, core.CotKH_NgaySinh, ngaySinhMoi)
		ghi(sID, sheet, row, core.CotKH_GioiTinh, gioiTinhMoi)
		ghi(sID, sheet, row, core.CotKH_DiaChi, diaChiMoi)
		ghi(sID, sheet, row, core.CotKH_MaSoThue, maSoThueMoi)
		ghi(sID, sheet, row, core.CotKH_Zalo, zaloMoi)
		ghi(sID, sheet, row, core.CotKH_UrlFb, fbMoi)
		ghi(sID, sheet, row, core.CotKH_UrlTiktok, tiktokMoi)

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"}) 
	}
}

// API_DoiMatKhau
func API_DoiMatKhau(c *gin.Context) {
	passCu := strings.TrimSpace(c.PostForm("pass_cu"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không đúng quy tắc!"})
		return 
	}
	
	if kh, ok := core.TimKhachHangTheoCookie(cookie); ok {
		if !bao_mat.KiemTraMatKhau(passCu, kh.MatKhauHash) { 
			c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu cũ không đúng!"})
			return 
		}
		hash, _ := bao_mat.HashMatKhau(passMoi)
		kh.MatKhauHash = hash
		
		sID := kh.SpreadsheetID
		if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
		
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) 
	}
}

// API_DoiMaPin
func API_DoiMaPin(c *gin.Context) {
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !bao_mat.KiemTraMaPin(pinMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "PIN phải đủ 8 số!"})
		return 
	}
	
	if kh, ok := core.TimKhachHangTheoCookie(cookie); ok {
		if !bao_mat.KiemTraMatKhau(pinCu, kh.MaPinHash) {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN hiện tại không đúng!"})
			return
		}
		hashMoi, _ := bao_mat.HashMatKhau(pinMoi)
		kh.MaPinHash = hashMoi
		
		sID := kh.SpreadsheetID
		if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }

		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}) 
	}
}


// API_GuiOTPPin : Gửi mã PIN mới (Giả lập)
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := core.TimKhachHangTheoCookie(cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	// Tạo PIN mới
	newPinRaw := taoMaPIN8So()
	
	// Gửi mail (Giả lập log console)
	log.Printf("📧 [MAIL MOCK] Gửi PIN mới '%s' đến %s", newPinRaw, kh.Email)

	// Lưu PIN mới (đã hash)
	hashNewPin, _ := bao_mat.HashMatKhau(newPinRaw)
	kh.MaPinHash = hashNewPin
	
	sID := kh.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
	
	core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashNewPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã PIN mới vào Email (Kiểm tra Log)!"})
}
