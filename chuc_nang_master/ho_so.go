package chuc_nang_master

import (
	"net/http"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// ==========================================================
// 1. TRANG HỒ SƠ MASTER
// ==========================================================
func TrangHoSoMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	// Lớp khiên bảo vệ chéo
	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Đẩy ra View của Master
	c.HTML(http.StatusOK, "master_ho_so", gin.H{
		"TieuDe":   "Hồ sơ cá nhân",
		"NhanVien": kh,
		"QuyenHan": vaiTro,
	})
}

// ==========================================================
// 2. API CẬP NHẬT THÔNG TIN
// ==========================================================
func API_LuuHoSoMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy tài khoản!"})
		return
	}

	core.KhoaHeThong.Lock()
	kh.TenKhachHang = strings.TrimSpace(c.PostForm("ho_ten"))
	kh.DienThoai = strings.TrimSpace(c.PostForm("dien_thoai"))
	kh.NgaySinh = strings.TrimSpace(c.PostForm("ngay_sinh"))
	kh.DiaChi = strings.TrimSpace(c.PostForm("dia_chi"))
	kh.MaSoThue = strings.TrimSpace(c.PostForm("ma_so_thue"))
	
	// [MỚI] Bắt dữ liệu Ảnh Đại Diện từ Form
	kh.AnhDaiDien = strings.TrimSpace(c.PostForm("anh_dai_dien"))
	
	kh.MangXaHoi.Zalo = strings.TrimSpace(c.PostForm("zalo"))
	kh.MangXaHoi.Facebook = strings.TrimSpace(c.PostForm("url_fb"))
	kh.MangXaHoi.Tiktok = strings.TrimSpace(c.PostForm("url_tiktok"))
	
	gioiTinh := c.PostForm("gioi_tinh")
	if gioiTinh == "Nam" { kh.GioiTinh = 1 } else if gioiTinh == "Nữ" { kh.GioiTinh = 0 } else { kh.GioiTinh = -1 }
	
	loc := time.FixedZone("ICT", 7*3600)
	kh.NgayCapNhat = time.Now().In(loc).Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

	ghi := core.ThemVaoHangCho
	r := kh.DongTrongSheet
	sh := "KHACH_HANG"

	ghi(shopID, sh, r, core.CotKH_TenKhachHang, kh.TenKhachHang)
	ghi(shopID, sh, r, core.CotKH_DienThoai, kh.DienThoai)
	ghi(shopID, sh, r, core.CotKH_NgaySinh, kh.NgaySinh)
	ghi(shopID, sh, r, core.CotKH_GioiTinh, kh.GioiTinh)
	ghi(shopID, sh, r, core.CotKH_DiaChi, kh.DiaChi)
	ghi(shopID, sh, r, core.CotKH_MaSoThue, kh.MaSoThue)
	ghi(shopID, sh, r, core.CotKH_MangXaHoiJson, core.ToJSON(kh.MangXaHoi))
	
	// [MỚI] Ghi dữ liệu Ảnh Đại Diện xuống Sheet
	ghi(shopID, sh, r, core.CotKH_AnhDaiDien, kh.AnhDaiDien)
	
	ghi(shopID, sh, r, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
}

// ==========================================================
// 3. API ĐỔI MẬT KHẨU
// ==========================================================
func API_DoiMatKhauMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	passCu := c.PostForm("pass_cu")
	passMoi := c.PostForm("pass_moi")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản không tồn tại!"})
		return
	}

	if !cau_hinh.KiemTraMatKhau(passCu, kh.MatKhauHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu hiện tại không đúng!"})
		return
	}
	if !cau_hinh.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	hash, _ := cau_hinh.HashMatKhau(passMoi)
	core.KhoaHeThong.Lock()
	kh.MatKhauHash = hash
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

// ==========================================================
// 4. API ĐỔI MÃ PIN
// ==========================================================
func API_DoiMaPinMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	pinCu := c.PostForm("pin_cu")
	pinMoi := c.PostForm("pin_moi")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản không tồn tại!"})
		return
	}

	if !cau_hinh.KiemTraMatKhau(pinCu, kh.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN hiện tại không đúng!"})
		return
	}
	if !cau_hinh.KiemTraMaPin(pinMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN mới không hợp lệ!"})
		return
	}

	hash, _ := cau_hinh.HashMatKhau(pinMoi)
	core.KhoaHeThong.Lock()
	kh.MaPinHash = hash
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
}
