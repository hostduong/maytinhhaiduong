package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/cau_hinh" // Chứa hàm kiểm tra (KiemTraHoTen...)
	"app/core"     // Chứa Struct, Const, Hàm DB

	"github.com/gin-gonic/gin"
)

// Trang Đăng Ký (View)
func TrangDangKy(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	cookie, _ := c.Cookie("session_id")
	
	// Check nếu đã đăng nhập thì đá về trang chủ
	if cookie != "" {
		if _, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "dang_ky", gin.H{"TieuDe": "Đăng Ký Tài Khoản"})
}

// Xử Lý Đăng Ký (Logic)
func XuLyDangKy(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // Lấy ShopID từ Middleware

	// 1. LẤY DỮ LIỆU TỪ FORM
	hoTen     := strings.TrimSpace(c.PostForm("ho_ten"))
	user      := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email     := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass      := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin     := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")) 
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }
	
	ngaySinh  := strings.TrimSpace(c.PostForm("ngay_sinh"))
	
	// Convert giới tính
	gioiTinhStr := c.PostForm("gioi_tinh")
	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	// 2. VALIDATE DỮ LIỆU (Dùng hàm từ /cau_hinh/kiem_tra.go)
	if !cau_hinh.KiemTraHoTen(hoTen) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Họ tên không hợp lệ!"})
		return
	}
	if !cau_hinh.KiemTraTenDangNhap(user) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập không đúng quy tắc!"})
		return
	}
	if !cau_hinh.KiemTraEmail(email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Email không hợp lệ!"})
		return
	}
	if !cau_hinh.KiemTraMaPin(maPin) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Mã PIN phải đúng 8 số!"})
		return
	}
	if !cau_hinh.KiemTraDinhDangMatKhau(pass) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Mật khẩu chứa ký tự không cho phép!"})
		return
	}

	// 3. KIỂM TRA TRÙNG LẶP (Trong phạm vi Shop)
	if _, ok := core.TimKhachHangTheoUserOrEmail(shopID, user); ok {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập đã tồn tại!"})
		return
	}
	if _, ok := core.TimKhachHangTheoUserOrEmail(shopID, email); ok {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Email đã được sử dụng!"})
		return
	}

	// 4. LOGIC TẠO NGƯỜI DÙNG (ADMIN ĐẦU TIÊN)
	listHienTai := core.LayDanhSachKhachHang(shopID)
	soLuong := len(listHienTai)
	
	var maKH, vaiTro, chucVu string
	
	if soLuong == 0 {
		// --- NGƯỜI ĐẦU TIÊN CỦA SHOP -> LÀM CHỦ ---
		maKH = "0000000000000000001" // Giữ nguyên mã Admin huyền thoại
		vaiTro = "admin_root"
		chucVu = "Quản trị cấp cao"
	} else {
		// --- KHÁCH BÌNH THƯỜNG ---
		maKH = core.TaoMaKhachHangMoi(shopID)
		vaiTro = "customer"
		chucVu = "Khách hàng"
	}

	// 5. MÃ HÓA MẬT KHẨU
	passHash, _ := cau_hinh.HashMatKhau(pass)
	pinHash, _ := cau_hinh.HashMatKhau(maPin)
	
	// Tạo Session đầu tiên
	sessionID := cau_hinh.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	ttl := cau_hinh.ThoiGianHetHanCookie
	expTime := time.Now().Add(ttl).Unix()
	
	nowStr := time.Now().Format("2006-01-02 15:04:05")

	// 6. KHỞI TẠO CÁC STRUCT CON (JSON DATA)
	
	// Map Token (Cột F)
	tokens := make(map[string]core.TokenInfo)
	tokens[sessionID] = core.TokenInfo{
		DeviceName: userAgent,
		ExpiresAt:  expTime,
	}

	// Mạng xã hội (Cột P) - Mặc định rỗng
	mxh := core.SocialInfo{} 
	
	// Ví tiền (Cột U) - Mặc định 0đ
	vi := core.WalletInfo{ SoDuHienTai: 0 }
	
	// Cấu hình (Cột V)
	conf := core.UserConfig{ Theme: "light", Language: "vi" }

	// 7. TẠO STRUCT KHACH HANG HOÀN CHỈNH
	newKH := &core.KhachHang{
		SpreadsheetID:  shopID,
		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		MatKhauHash:    passHash,
		MaPinHash:      pinHash,
		RefreshTokens:  tokens, // Map token
		
		VaiTroQuyenHan: vaiTro,
		ChucVu:         chucVu,
		TrangThai:      1,
		
		NguonKhachHang: "web_register",
		TenKhachHang:   hoTen,
		DienThoai:      dienThoai,
		NgaySinh:       ngaySinh,
		GioiTinh:       gioiTinh,
		
		MangXaHoi:      mxh,
		ViTien:         vi,
		CauHinh:        conf,
		
		NgayTao:        nowStr,
		NgayCapNhat:    nowStr,
	}

	// 8. LƯU VÀO RAM (Cache)
	// Tính dòng tiếp theo (Header + Số lượng hiện tại)
	newKH.DongTrongSheet = core.DongBatDau_KhachHang + soLuong
	core.ThemKhachHangVaoRam(newKH) // Hàm này trong core/khach_hang.go

	// 9. GHI XUỐNG SHEET (QUEUE)
	// Dùng hàm Helper core.ThemVaoHangCho
	ghi := core.ThemVaoHangCho
	row := newKH.DongTrongSheet
	sheet := "KHACH_HANG"

	// --- NHÓM CỘT THƯỜNG ---
	ghi(shopID, sheet, row, core.CotKH_MaKhachHang, newKH.MaKhachHang)
	ghi(shopID, sheet, row, core.CotKH_TenDangNhap, newKH.TenDangNhap)
	ghi(shopID, sheet, row, core.CotKH_Email, newKH.Email)
	ghi(shopID, sheet, row, core.CotKH_MatKhauHash, newKH.MatKhauHash)
	ghi(shopID, sheet, row, core.CotKH_MaPinHash, newKH.MaPinHash)
	
	ghi(shopID, sheet, row, core.CotKH_VaiTroQuyenHan, newKH.VaiTroQuyenHan)
	ghi(shopID, sheet, row, core.CotKH_ChucVu, newKH.ChucVu)
	ghi(shopID, sheet, row, core.CotKH_TrangThai, newKH.TrangThai)
	
	ghi(shopID, sheet, row, core.CotKH_NguonKhachHang, newKH.NguonKhachHang)
	ghi(shopID, sheet, row, core.CotKH_TenKhachHang, newKH.TenKhachHang)
	ghi(shopID, sheet, row, core.CotKH_DienThoai, newKH.DienThoai)
	ghi(shopID, sheet, row, core.CotKH_NgaySinh, newKH.NgaySinh)
	ghi(shopID, sheet, row, core.CotKH_GioiTinh, newKH.GioiTinh)
	ghi(shopID, sheet, row, core.CotKH_NgayTao, newKH.NgayTao)

	// --- NHÓM CỘT JSON (QUAN TRỌNG) ---
	ghi(shopID, sheet, row, core.CotKH_RefreshTokenJson, core.ToJSON(newKH.RefreshTokens))
	ghi(shopID, sheet, row, core.CotKH_MangXaHoiJson, core.ToJSON(newKH.MangXaHoi))
	ghi(shopID, sheet, row, core.CotKH_ViTienJson, core.ToJSON(newKH.ViTien))
	ghi(shopID, sheet, row, core.CotKH_CauHinhJson, core.ToJSON(newKH.CauHinh))
	
	// 10. SET COOKIE VÀ CHUYỂN HƯỚNG
	// Tạo chữ ký bảo mật (Signature)
	signature := cau_hinh.TaoChuKyBaoMat(sessionID, userAgent)
	maxAge := int(ttl.Seconds())

	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	// Nếu là Admin -> Vào trang quản trị luôn
	if vaiTro == "admin_root" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
