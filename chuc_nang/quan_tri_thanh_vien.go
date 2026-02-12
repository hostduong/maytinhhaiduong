package chuc_nang

import (
	"net/http"
	"time"

	"app/cau_hinh" // [ĐÃ SỬA] Dùng cấu hình thay vì bao_mat
	"app/core"

	"github.com/gin-gonic/gin"
)

// TrangQuanLyThanhVien : Hiển thị danh sách khách hàng/nhân viên
func TrangQuanLyThanhVien(c *gin.Context) {
	// [SAAS] Lấy Context
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	// Lấy thông tin admin đang đăng nhập (Truyền shopID)
	// Lưu ý: Hàm LayKhachHang phải được định nghĩa trong core/khach_hang.go
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

// API_Admin_LuuThanhVien : Cập nhật thông tin và QUYỀN HẠN
func API_Admin_LuuThanhVien(c *gin.Context) {
	// [SAAS] Lấy Context
	shopID := c.GetString("SHOP_ID")
	myRole := c.GetString("USER_ROLE")
	
	// 1. Check quyền Admin
	if myRole != "admin_root" && myRole != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền quản trị nhân sự!"})
		return
	}

	// 2. Lấy dữ liệu
	maKH := c.PostForm("ma_khach_hang")
	role := c.PostForm("vai_tro") // admin, sale, kho...
	name := c.PostForm("ten_khach_hang")
	pass := c.PostForm("mat_khau_moi") // Nếu có nhập thì đổi, không thì thôi
	
	// Tìm khách hàng trong Shop
	kh, ok := core.LayKhachHang(shopID, maKH)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thành viên này!"})
		return
	}

	// 3. Logic chặn: Không ai được sửa Admin Root (trừ chính họ)
	if kh.VaiTroQuyenHan == "admin_root" && myRole != "admin_root" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không thể tác động đến tài khoản Super Admin!"})
		return
	}

	// 4. Cập nhật Core RAM
	core.KhoaHeThong.Lock()
	kh.TenKhachHang = name
	kh.VaiTroQuyenHan = role
	
	// Map lại tên chức vụ hiển thị cho đẹp
	switch role {
	case "admin": kh.ChucVu = "Quản trị viên"
	case "sale": kh.ChucVu = "Nhân viên kinh doanh"
	case "kho": kh.ChucVu = "Thủ kho"
	case "customer": kh.ChucVu = "Khách hàng"
	default: kh.ChucVu = "Thành viên"
	}
	
	// Đổi pass nếu có
	if pass != "" {
		// [ĐÃ SỬA] Dùng hàm từ cau_hinh
		hash, _ := cau_hinh.HashMatKhau(pass)
		kh.MatKhauHash = hash
		
		// Ghi đè pass xuống Sheet (Cột MatKhauHash)
		// Lưu ý: shopID chính là SpreadsheetID trong mô hình SaaS
		core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	}
	
	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()

	// 5. Ghi xuống Sheet (Các thông tin khác)
	ghi := core.ThemVaoHangCho
	row := kh.DongTrongSheet
	sheet := "KHACH_HANG"
	
	ghi(shopID, sheet, row, core.CotKH_TenKhachHang, kh.TenKhachHang)
	ghi(shopID, sheet, row, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	ghi(shopID, sheet, row, core.CotKH_ChucVu, kh.ChucVu)
	ghi(shopID, sheet, row, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thành viên thành công!"})
}
