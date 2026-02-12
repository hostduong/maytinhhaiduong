package chuc_nang

import (
	"net/http"
	"time"

	"app/bao_mat" // Dùng hàm Hash
	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangQuanLyThanhVien(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	// Lấy thông tin admin
	me, _ := core.LayKhachHang(shopID, userID)

	// Lấy danh sách thành viên của Shop
	listAll := core.LayDanhSachKhachHang(shopID)

	c.HTML(http.StatusOK, "quan_tri_thanh_vien", gin.H{
		"TieuDe":       "Quản lý thành viên",
		"NhanVien":     me,
		"DanhSach":     listAll,
		"TenNguoiDung": me.TenKhachHang,
		"QuyenHan":     me.VaiTroQuyenHan,
	})
}

func API_Admin_LuuThanhVien(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	myRole := c.GetString("USER_ROLE")
	
	if myRole != "admin_root" && myRole != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền quản trị nhân sự!"})
		return
	}

	maKH := c.PostForm("ma_khach_hang")
	role := c.PostForm("vai_tro")
	name := c.PostForm("ten_khach_hang")
	pass := c.PostForm("mat_khau_moi")
	
	// Tìm khách hàng trong Shop
	kh, ok := core.LayKhachHang(shopID, maKH)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thành viên này!"})
		return
	}

	if kh.VaiTroQuyenHan == "admin_root" && myRole != "admin_root" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không thể tác động đến tài khoản Super Admin!"})
		return
	}

	core.KhoaHeThong.Lock()
	kh.TenKhachHang = name
	kh.VaiTroQuyenHan = role
	
	switch role {
	case "admin": kh.ChucVu = "Quản trị viên"
	case "sale": kh.ChucVu = "Nhân viên kinh doanh"
	case "kho": kh.ChucVu = "Thủ kho"
	case "customer": kh.ChucVu = "Khách hàng"
	default: kh.ChucVu = "Thành viên"
	}
	
	if pass != "" {
		hash, _ := cau_hinh.HashMatKhau(pass) // Sửa lại dùng cau_hinh nếu đã gộp
		kh.MatKhauHash = hash
		
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	}
	
	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

	// Ghi Sheet (Truyền shopID)
	ghi := core.ThemVaoHangCho
	row := kh.DongTrongSheet
	
	ghi(shopID, "KHACH_HANG", row, core.CotKH_TenKhachHang, kh.TenKhachHang)
	ghi(shopID, "KHACH_HANG", row, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	ghi(shopID, "KHACH_HANG", row, core.CotKH_ChucVu, kh.ChucVu)
	ghi(shopID, "KHACH_HANG", row, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thành viên thành công!"})
}
