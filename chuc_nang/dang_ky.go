package chuc_nang

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

// ==========================================================
// 1. TRANG ĐĂNG KÝ
// ==========================================================
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

// ==========================================================
// 2. XỬ LÝ ĐĂNG KÝ (PHÂN LUỒNG MASTER VÀ TENANT)
// ==========================================================
func XuLyDangKy(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") // Để biết đang ở 99k.vn hay shopA.99k.vn

	hoTen     := strings.TrimSpace(c.PostForm("ho_ten"))
	user      := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	email     := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	pass      := strings.TrimSpace(c.PostForm("mat_khau"))
	maPin     := strings.TrimSpace(c.PostForm("ma_pin"))
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")) 
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }
	ngaySinh  := strings.TrimSpace(c.PostForm("ngay_sinh"))
	
	gioiTinhStr := c.PostForm("gioi_tinh")
	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	// Validate Dữ liệu
	if !cau_hinh.KiemTraHoTen(hoTen) || !cau_hinh.KiemTraTenDangNhap(user) || !cau_hinh.KiemTraEmail(email) || !cau_hinh.KiemTraMaPin(maPin) || !cau_hinh.KiemTraDinhDangMatKhau(pass) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Dữ liệu nhập vào không hợp lệ!"})
		return
	}

	if _, ok := core.TimKhachHangTheoUserOrEmail(shopID, user); ok {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập đã tồn tại!"})
		return
	}
	if _, ok := core.TimKhachHangTheoUserOrEmail(shopID, email); ok {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Email đã được sử dụng!"})
		return
	}

	// -----------------------------------------------------
	// [MỚI] LOGIC PHÂN QUYỀN VÀ TRẠNG THÁI THEO HỆ THỐNG
	// -----------------------------------------------------
	listHienTai := core.LayDanhSachKhachHang(shopID)
	soLuong := len(listHienTai)
	var maKH, vaiTro, chucVu string
	var trangThai int

	if theme == "theme_master" {
		// Dành cho Nền tảng 99k.vn
		if soLuong == 0 {
			maKH = "0000000000000000001"
			vaiTro = "quan_tri_vien_he_thong" // Admin tối cao của nền tảng
			chucVu = "Quản trị hệ thống"
			trangThai = 1 // Không cần OTP
		} else {
			maKH = core.TaoMaKhachHangMoi(shopID)
			vaiTro = "khach_hang" // Đối với nền tảng, họ là khách mua phần mềm
			chucVu = "Chủ cửa hàng"
			trangThai = 0 // Bắt buộc chờ xác thực OTP
		}
	} else {
		// Dành cho Cửa hàng (B2C)
		if soLuong == 0 {
			maKH = "0000000000000000001"
			vaiTro = "quan_tri_vien" // Áp dụng theo đúng File Phân Quyền PDF
			chucVu = "Quản trị viên"
			trangThai = 1
		} else {
			maKH = core.TaoMaKhachHangMoi(shopID)
			vaiTro = "khach_hang"
			chucVu = "Khách hàng"
			trangThai = 1 // Khách mua lẻ không cần xác thực rườm rà
		}
	}

	passHash, _ := cau_hinh.HashMatKhau(pass)
	pinHash, _ := cau_hinh.HashMatKhau(maPin)
	
	nowStr := time.Now().Format("2006-01-02 15:04:05")

	newKH := &core.KhachHang{
		SpreadsheetID:  shopID,
		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		MatKhauHash:    passHash,
		MaPinHash:      pinHash,
		RefreshTokens:  make(map[string]core.TokenInfo), 
		VaiTroQuyenHan: vaiTro,
		ChucVu:         chucVu,
		TrangThai:      trangThai,
		DataSheets:     core.DataSheetInfo{},
		GoiDichVu:      make([]core.PlanInfo, 0),
		CauHinh:        core.UserConfig{ Theme: "light", Language: "vi" },
		NguonKhachHang: "web_register",
		TenKhachHang:   hoTen,
		DienThoai:      dienThoai,
		MangXaHoi:      core.SocialInfo{},
		NgaySinh:       ngaySinh,
		GioiTinh:       gioiTinh,
		ViTien:         core.WalletInfo{ SoDuHienTai: 0 },
		Inbox:          make([]core.MessageInfo, 0),
		NgayTao:        nowStr,
		NguoiCapNhat:   user,
		NgayCapNhat:    nowStr,
	}

	// LƯU VÀO RAM & GHI XUỐNG SHEET (Lược bớt code ghi sheet cho gọn, bạn giữ nguyên hàm ghi 27 cột của bạn nhé)
	newKH.DongTrongSheet = core.DongBatDau_KhachHang + soLuong
	core.ThemKhachHangVaoRam(newKH)
	
	ghi := core.ThemVaoHangCho
	sh := "KHACH_HANG"
	r := newKH.DongTrongSheet
	ghi(shopID, sh, r, core.CotKH_MaKhachHang, newKH.MaKhachHang)
	ghi(shopID, sh, r, core.CotKH_TenDangNhap, newKH.TenDangNhap)
	ghi(shopID, sh, r, core.CotKH_Email, newKH.Email)
	ghi(shopID, sh, r, core.CotKH_MatKhauHash, newKH.MatKhauHash)
	ghi(shopID, sh, r, core.CotKH_MaPinHash, newKH.MaPinHash)
	ghi(shopID, sh, r, core.CotKH_VaiTroQuyenHan, newKH.VaiTroQuyenHan)
	ghi(shopID, sh, r, core.CotKH_ChucVu, newKH.ChucVu)
	ghi(shopID, sh, r, core.CotKH_TrangThai, newKH.TrangThai)
	// (GHI CÁC CỘT CÒN LẠI VÀO ĐÂY THEO CODE CŨ...)

	// -----------------------------------------------------
	// [ĐIỀU HƯỚNG]: NẾU TRẠNG THÁI = 0 -> GỬI MAIL VÀ XÁC THỰC
	// -----------------------------------------------------
	if trangThai == 0 {
		code := taoMaOTP6So() // Gọi hàm từ quen_mat_khau.go
		luuOTPCucBo(shopID, user, code)
		
		log.Printf("📧 [MAIL MOCK] Gửi OTP KÍCH HOẠT '%s' đến %s", code, email)
		// Đá sang trang nhập OTP
		c.Redirect(http.StatusFound, "/xac-thuc?u=" + user)
		return
	}

	// NẾU TRẠNG THÁI = 1 -> ĐĂNG NHẬP LUÔN
	sessionID := cau_hinh.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	ttl := cau_hinh.ThoiGianHetHanCookie
	expTime := time.Now().Add(ttl).Unix()
	
	newKH.RefreshTokens[sessionID] = core.TokenInfo{ DeviceName: userAgent, ExpiresAt: expTime }
	core.ThemVaoHangCho(shopID, sh, r, core.CotKH_RefreshTokenJson, core.ToJSON(newKH.RefreshTokens))

	signature := cau_hinh.TaoChuKyBaoMat(sessionID, userAgent)
	maxAge := int(ttl.Seconds())
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	if vaiTro == "quan_tri_vien_he_thong" || vaiTro == "quan_tri_vien" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}


// ==========================================================
// 3. API XÁC THỰC OTP, BƠM GÓI TRIAL VÀ CẤP SUBDOMAIN
// ==========================================================
func TrangXacThuc(c *gin.Context) {
	c.HTML(http.StatusOK, "xac_thuc_otp", gin.H{"User": c.Query("u")})
}

func XuLyXacThucOTP(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	user := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	otp := strings.TrimSpace(c.PostForm("otp"))

	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, user)
	if !ok || !kiemTraOTPCucBo(shopID, user, otp) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP không đúng hoặc đã hết hạn!"})
		return
	}

	// 1. MỞ KHÓA TÀI KHOẢN VÀ BƠM GÓI TRIAL
	core.KhoaHeThong.Lock()
	kh.TrangThai = 1
	kh.GoiDichVu = append(kh.GoiDichVu, core.PlanInfo{
		MaGoi:          "TRIAL_3DAYS",
		TenGoiDichVu:   "Dùng thử 3 ngày", // <-- SỬA THÀNH TenGoiDichVu
		NgayHetHan:     time.Now().AddDate(0, 0, 3).Format("2006-01-02 15:04:05"),
		TrangThai:      "active",
	})
	
	// Tạo Session Đăng nhập
	sessionID := cau_hinh.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	ttl := cau_hinh.ThoiGianHetHanCookie
	expTime := time.Now().Add(ttl).Unix()
	kh.RefreshTokens[sessionID] = core.TokenInfo{ DeviceName: userAgent, ExpiresAt: expTime }
	core.KhoaHeThong.Unlock()

	// 2. GHI XUỐNG SHEET
	ghi := core.ThemVaoHangCho
	r := kh.DongTrongSheet
	sh := "KHACH_HANG"
	ghi(shopID, sh, r, core.CotKH_TrangThai, 1)
	ghi(shopID, sh, r, core.CotKH_GoiDichVuJson, core.ToJSON(kh.GoiDichVu))
	ghi(shopID, sh, r, core.CotKH_RefreshTokenJson, core.ToJSON(kh.RefreshTokens))

	// 3. CHẠY NGẦM TẠO SUBDOMAIN
	go func(sub string) {
		TuDongThemSubdomain(sub)
	}(kh.TenDangNhap)

	// 4. SET COOKIE
	signature := cau_hinh.TaoChuKyBaoMat(sessionID, userAgent)
	maxAge := int(ttl.Seconds())
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	c.JSON(200, gin.H{"status": "ok", "msg": "Xác thực thành công! Hệ thống đang khởi tạo..."})
}

// Code tự động kích hoạt Subdomain
func TuDongThemSubdomain(subdomain string) error {
	ctx := context.Background()
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson 
	if jsonKey == "" { return nil } // Bỏ qua nếu chưa config
	
	srv, err := run.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil { return err }
	
	fullDomain := subdomain + ".99k.vn"
	parent := "projects/project-47337221-fda1-48c7-b2f/locations/asia-southeast1" 

	req := &run.DomainMapping{
		Metadata: &run.ObjectMeta{ Name: fullDomain },
		Spec: &run.DomainMappingSpec{
			RouteName:       "maytinhhaiduong",
			CertificateMode: "AUTOMATIC",
		},
	}

	_, err = srv.Namespaces.Domainmappings.Create(parent, req).Do()
	if err != nil {
		log.Printf("❌ Lỗi tạo subdomain %s: %v", fullDomain, err)
		return err
	}
	
	log.Printf("✅ Đã lệnh cho Google tạo subdomain: %s", fullDomain)
	return nil
}
