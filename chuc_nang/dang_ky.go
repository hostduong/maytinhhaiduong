package chuc_nang

import (
	"context" // [MỚI THÊM]
	"log"     // [MỚI THÊM]
	"net/http"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option" // [MỚI THÊM]
	"google.golang.org/api/run/v1"
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

	// 2. VALIDATE DỮ LIỆU
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

	// 3. KIỂM TRA TRÙNG LẶP
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
		maKH = "0000000000000000001"
		vaiTro = "admin_root"
		chucVu = "Quản trị cấp cao"
	} else {
		maKH = core.TaoMaKhachHangMoi(shopID)
		vaiTro = "customer"
		chucVu = "Khách hàng"
	}

	// 5. MÃ HÓA MẬT KHẨU
	passHash, _ := cau_hinh.HashMatKhau(pass)
	pinHash, _ := cau_hinh.HashMatKhau(maPin)
	
	sessionID := cau_hinh.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	ttl := cau_hinh.ThoiGianHetHanCookie
	expTime := time.Now().Add(ttl).Unix()
	
	nowStr := time.Now().Format("2006-01-02 15:04:05")

	// 6. KHỞI TẠO CÁC STRUCT CON (JSON DATA AN TOÀN)
	tokens := make(map[string]core.TokenInfo)
	tokens[sessionID] = core.TokenInfo{ DeviceName: userAgent, ExpiresAt: expTime }

	dsInfo := core.DataSheetInfo{}
	plans  := make([]core.PlanInfo, 0)
	conf   := core.UserConfig{ Theme: "light", Language: "vi" }
	mxh    := core.SocialInfo{} 
	vi     := core.WalletInfo{ SoDuHienTai: 0 }
	inbox  := make([]core.MessageInfo, 0)

	// 7. TẠO STRUCT KHACH HANG HOÀN CHỈNH
	newKH := &core.KhachHang{
		SpreadsheetID:  shopID,
		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		MatKhauHash:    passHash,
		MaPinHash:      pinHash,
		RefreshTokens:  tokens, 
		
		VaiTroQuyenHan: vaiTro,
		ChucVu:         chucVu,
		TrangThai:      1,
		
		DataSheets:     dsInfo,
		GoiDichVu:      plans,
		CauHinh:        conf,

		NguonKhachHang: "web_register",
		TenKhachHang:   hoTen,
		DienThoai:      dienThoai,
		AnhDaiDien:     "",
		MangXaHoi:      mxh,
		DiaChi:         "",
		NgaySinh:       ngaySinh,
		GioiTinh:       gioiTinh,
		MaSoThue:       "",
		ViTien:         vi,
		Inbox:          inbox,
		
		GhiChu:         "",
		NgayTao:        nowStr,
		NguoiCapNhat:   user, // Chính người này tạo
		NgayCapNhat:    nowStr,
	}

	// 8. LƯU VÀO RAM
	newKH.DongTrongSheet = core.DongBatDau_KhachHang + soLuong
	core.ThemKhachHangVaoRam(newKH)

	// 9. GHI XUỐNG SHEET (ĐÚNG THỨ TỰ 27 CỘT MỚI)
	ghi := core.ThemVaoHangCho
	row := newKH.DongTrongSheet
	sheet := "KHACH_HANG"

	// Cột A -> I
	ghi(shopID, sheet, row, core.CotKH_MaKhachHang, newKH.MaKhachHang)
	ghi(shopID, sheet, row, core.CotKH_TenDangNhap, newKH.TenDangNhap)
	ghi(shopID, sheet, row, core.CotKH_Email, newKH.Email)
	ghi(shopID, sheet, row, core.CotKH_MatKhauHash, newKH.MatKhauHash)
	ghi(shopID, sheet, row, core.CotKH_MaPinHash, newKH.MaPinHash)
	ghi(shopID, sheet, row, core.CotKH_RefreshTokenJson, core.ToJSON(newKH.RefreshTokens))
	ghi(shopID, sheet, row, core.CotKH_VaiTroQuyenHan, newKH.VaiTroQuyenHan)
	ghi(shopID, sheet, row, core.CotKH_ChucVu, newKH.ChucVu)
	ghi(shopID, sheet, row, core.CotKH_TrangThai, newKH.TrangThai)
	
	// Cột J, K, L (Core SaaS JSON)
	ghi(shopID, sheet, row, core.CotKH_DataSheetsJson, core.ToJSON(newKH.DataSheets))
	ghi(shopID, sheet, row, core.CotKH_GoiDichVuJson, core.ToJSON(newKH.GoiDichVu))
	ghi(shopID, sheet, row, core.CotKH_CauHinhJson, core.ToJSON(newKH.CauHinh))
	
	// Cột M -> U
	ghi(shopID, sheet, row, core.CotKH_NguonKhachHang, newKH.NguonKhachHang)
	ghi(shopID, sheet, row, core.CotKH_TenKhachHang, newKH.TenKhachHang)
	ghi(shopID, sheet, row, core.CotKH_DienThoai, newKH.DienThoai)
	ghi(shopID, sheet, row, core.CotKH_AnhDaiDien, newKH.AnhDaiDien)
	ghi(shopID, sheet, row, core.CotKH_MangXaHoiJson, core.ToJSON(newKH.MangXaHoi))
	ghi(shopID, sheet, row, core.CotKH_DiaChi, newKH.DiaChi)
	ghi(shopID, sheet, row, core.CotKH_NgaySinh, newKH.NgaySinh)
	ghi(shopID, sheet, row, core.CotKH_GioiTinh, newKH.GioiTinh)
	ghi(shopID, sheet, row, core.CotKH_MaSoThue, newKH.MaSoThue)
	
	// Cột V, W, X, Y, Z, AA
	ghi(shopID, sheet, row, core.CotKH_ViTienJson, core.ToJSON(newKH.ViTien))
	ghi(shopID, sheet, row, core.CotKH_InboxJson, core.ToJSON(newKH.Inbox))
	ghi(shopID, sheet, row, core.CotKH_GhiChu, newKH.GhiChu)
	ghi(shopID, sheet, row, core.CotKH_NgayTao, newKH.NgayTao)
	ghi(shopID, sheet, row, core.CotKH_NguoiCapNhat, newKH.NguoiCapNhat)
	ghi(shopID, sheet, row, core.CotKH_NgayCapNhat, newKH.NgayCapNhat)
	
	// 10. SET COOKIE VÀ CHUYỂN HƯỚNG
	signature := cau_hinh.TaoChuKyBaoMat(sessionID, userAgent)
	maxAge := int(ttl.Seconds())

	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	if vaiTro == "admin_root" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

// Code minh họa logic cốt lõi
func TuDongThemSubdomain(subdomain string) error {
	ctx := context.Background()
	
    // 1. Dùng file JSON đang có để xác thực với Google
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson 
	srv, err := run.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	
	// 2. Tên miền đầy đủ
	fullDomain := subdomain + ".99k.vn"
    // ID của project và khu vực (vd: projects/my-project/locations/asia-southeast1)
	parent := "projects/project-47337221-fda1-48c7-b2f/locations/asia-southeast1" 

	// 3. Cấu hình yêu cầu Add Mapping
	req := &run.DomainMapping{
		Metadata: &run.ObjectMeta{
			Name: fullDomain,
		},
		Spec: &run.DomainMappingSpec{
			RouteName:       "maytinhhaiduong", // Tên service Cloud Run của bạn
			CertificateMode: "AUTOMATIC",     // Google tự lo SSL
		},
	}

	// 4. Gửi lệnh lên Google Cloud
	_, err = srv.Namespaces.Domainmappings.Create(parent, req).Do()
	if err != nil {
		log.Printf("❌ Lỗi tạo subdomain %s: %v", fullDomain, err)
		return err
	}
	
	log.Printf("✅ Đã lệnh cho Google tạo subdomain: %s", fullDomain)
	return nil
}
