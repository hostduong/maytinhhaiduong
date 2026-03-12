package ho_so

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"app/config"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangHoSoMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.HTML(http.StatusOK, "master_ho_so", gin.H{
		"TieuDe":   "Hồ sơ",
		"NhanVien": kh,
		"QuyenHan": vaiTro,
	})
}

func API_LuuHoSoMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy tài khoản!"})
		return
	}

	core.KhoaHeThong.Lock()
	kh.ThongTin.TenKhachHang = strings.TrimSpace(c.PostForm("ho_ten"))
	kh.ThongTin.DienThoai = strings.TrimSpace(c.PostForm("dien_thoai"))
	kh.ThongTin.NgaySinh = strings.TrimSpace(c.PostForm("ngay_sinh"))
	kh.ThongTin.DiaChi = strings.TrimSpace(c.PostForm("dia_chi"))
	kh.ThongTin.MaSoThue = strings.TrimSpace(c.PostForm("ma_so_thue"))
	kh.ThongTin.AnhDaiDien = strings.TrimSpace(c.PostForm("anh_dai_dien"))
	
	if kh.MangXaHoi == nil { kh.MangXaHoi = make(map[string]string) }
	kh.MangXaHoi["zalo"] = strings.TrimSpace(c.PostForm("zalo"))
	kh.MangXaHoi["facebook"] = strings.TrimSpace(c.PostForm("url_fb"))
	kh.MangXaHoi["tiktok"] = strings.TrimSpace(c.PostForm("url_tiktok"))
	
	gioiTinh := c.PostForm("gioi_tinh")
	if gioiTinh == "Nam" { kh.ThongTin.GioiTinh = 1 } else if gioiTinh == "Nữ" { kh.ThongTin.GioiTinh = 0 } else { kh.ThongTin.GioiTinh = -1 }
	
	kh.NgayCapNhat = time.Now().Unix()

	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	r := kh.DongTrongSheet
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(shopID, "KHACH_HANG", r, core.CotKH_DataJSON, jsonStr)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
}

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

	if !config.KiemTraMatKhau(passCu, kh.BaoMat.MatKhauHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu hiện tại không đúng!"})
		return
	}
	if !config.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	hash, _ := config.HashMatKhau(passMoi)
	core.KhoaHeThong.Lock()
	kh.BaoMat.MatKhauHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_DataJSON, jsonStr)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

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

	if !config.KiemTraMatKhau(pinCu, kh.BaoMat.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN hiện tại không đúng!"})
		return
	}
	if !config.KiemTraMaPin(pinMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN mới không hợp lệ!"})
		return
	}

	hash, _ := config.HashMatKhau(pinMoi)
	core.KhoaHeThong.Lock()
	kh.BaoMat.MaPinHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_DataJSON, jsonStr)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
}
