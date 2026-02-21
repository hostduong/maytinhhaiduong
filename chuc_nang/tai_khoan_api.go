package chuc_nang

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// Helper tạo mã PIN 8 số
func taoMaPIN8So() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

// API_DoiThongTin : Cập nhật Full thông tin cá nhân
func API_DoiThongTin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]

	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhStr := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	diaChiMoi   := strings.TrimSpace(c.PostForm("dia_chi"))
	maSoThueMoi := strings.TrimSpace(c.PostForm("ma_so_thue"))
	
	// Mạng xã hội
	zaloMoi     := strings.TrimSpace(c.PostForm("zalo"))
	fbMoi       := strings.TrimSpace(c.PostForm("url_fb"))
	tiktokMoi   := strings.TrimSpace(c.PostForm("url_tiktok"))

	// Validate
	if !cau_hinh.KiemTraHoTen(hoTenMoi) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ!"})
		return 
	}
	
	cookie, _ := c.Cookie("session_id")
	
	// [SAAS] Tìm khách hàng trong Shop
	if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
		
		core.KhoaHeThong.Lock()
		kh.TenKhachHang = hoTenMoi
		kh.DienThoai = sdtMoi
		kh.NgaySinh = ngaySinhMoi
		
		// Convert giới tính
		if gioiTinhStr == "Nam" { 
			kh.GioiTinh = 1 
		} else if gioiTinhStr == "Nữ" { 
			kh.GioiTinh = 0 
		} else { 
			kh.GioiTinh = -1 
		}
		
		kh.DiaChi = diaChiMoi
		kh.MaSoThue = maSoThueMoi
		
		// Cập nhật Struct con (Mạng xã hội)
		kh.MangXaHoi.Zalo = zaloMoi
		kh.MangXaHoi.Facebook = fbMoi
		kh.MangXaHoi.Tiktok = tiktokMoi
		
		kh.NguoiCapNhat = kh.TenDangNhap // Lưu vết chính user tự sửa
		kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
		core.KhoaHeThong.Unlock()

		// Ghi xuống Sheet
		row := kh.DongTrongSheet
		sheet := "KHACH_HANG"
		ghi := core.ThemVaoHangCho

		// Ghi cột thường
		ghi(shopID, sheet, row, core.CotKH_TenKhachHang, kh.TenKhachHang)
		ghi(shopID, sheet, row, core.CotKH_DienThoai, kh.DienThoai)
		ghi(shopID, sheet, row, core.CotKH_NgaySinh, kh.NgaySinh)
		ghi(shopID, sheet, row, core.CotKH_GioiTinh, kh.GioiTinh)
		ghi(shopID, sheet, row, core.CotKH_DiaChi, kh.DiaChi)
		ghi(shopID, sheet, row, core.CotKH_MaSoThue, kh.MaSoThue)
		
		// Ghi JSON
		jsonMXH := core.ToJSON(kh.MangXaHoi)
		ghi(shopID, sheet, row, core.CotKH_MangXaHoiJson, jsonMXH)
		
		// Lưu vết
		ghi(shopID, sheet, row, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
		ghi(shopID, sheet, row, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"}) 
	}
}

// API_DoiMatKhau
func API_DoiMatKhau(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
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
		kh.NguoiCapNhat = kh.TenDangNhap
		kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
		core.KhoaHeThong.Unlock()
		
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
		
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) 
	}
}

// API_DoiMaPin
func API_DoiMaPin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
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
		kh.NguoiCapNhat = kh.TenDangNhap
		kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
		core.KhoaHeThong.Unlock()

		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashMoi)
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
		
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else { 
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}) 
	}
}

// API_GuiOTPPin : Gửi mã PIN mới (Giả lập)
func API_GuiOTPPin(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	cookie, _ := c.Cookie("session_id")
	
	kh, ok := core.TimKhachHangTheoCookie(shopID, cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	// Tạo PIN mới
	newPinRaw := taoMaPIN8So()
	
	// Gửi mail (Giả lập log console)
	log.Printf("📧 [MAIL MOCK] Shop [%s] - Gửi PIN mới '%s' đến %s", shopID, newPinRaw, kh.Email)

	// Lưu PIN mới
	hashNewPin, _ := cau_hinh.HashMatKhau(newPinRaw)
	
	core.KhoaHeThong.Lock()
	kh.MaPinHash = hashNewPin
	kh.NguoiCapNhat = "Hệ thống" // Reset tự động thì ghi là Hệ thống
	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()
	
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashNewPin)
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã PIN mới vào Email (Kiểm tra Log)!"})
}
