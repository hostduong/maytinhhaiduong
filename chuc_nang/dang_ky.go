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

func TrangDangKy(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	cookie, _ := c.Cookie("session_id")
	
	if cookie != "" {
		if _, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "dang_ky", gin.H{"TieuDe": "Đăng Ký Tài Khoản"})
}

func XuLyDangKy(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") 

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

	if dienThoai == "" || !cau_hinh.KiemTraHoTen(hoTen) || !cau_hinh.KiemTraTenDangNhap(user) || !cau_hinh.KiemTraEmail(email) || !cau_hinh.KiemTraMaPin(maPin) || !cau_hinh.KiemTraDinhDangMatKhau(pass) {
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

	listHienTai := core.LayDanhSachKhachHang(shopID)
	soLuong := len(listHienTai)
	var maKH, vaiTro, chucVu string

	loc := time.FixedZone("ICT", 7*3600) 
	nowVN := time.Now().In(loc)
	nowStr := nowVN.Format("2006-01-02 15:04:05")
	ghi := core.ThemVaoHangCho

	if theme == "theme_master" {
		if soLuong == 0 {
			// [MỚI] 1. TẠO TÀI KHOẢN HỆ THỐNG (BOT) TRƯỚC VÀO DÒNG 11
			bot := &core.KhachHang{
				SpreadsheetID:  shopID,
				DongTrongSheet: core.DongBatDau_KhachHang,
				MaKhachHang:    "0000000000000000000",
				TenDangNhap:    "admin",
				VaiTroQuyenHan: "admin",
				ChucVu:         "Hệ thống",
				TenKhachHang:   "Trợ lý ảo 99K",
				TrangThai:      1,
				NgayTao:        nowStr,
				NgayCapNhat:    nowStr,
			}
			core.ThemKhachHangVaoRam(bot)
			
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_MaKhachHang, bot.MaKhachHang)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_TenDangNhap, bot.TenDangNhap)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_VaiTroQuyenHan, bot.VaiTroQuyenHan)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_ChucVu, bot.ChucVu)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_TenKhachHang, bot.TenKhachHang)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_TrangThai, bot.TrangThai)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_NgayTao, bot.NgayTao)
			ghi(shopID, "KHACH_HANG", bot.DongTrongSheet, core.CotKH_NgayCapNhat, bot.NgayCapNhat)

			// 2. THIẾT LẬP BIẾN CHO TÀI KHOẢN ROOT (SẼ ĐƯỢC LƯU Ở DÒNG 12 DO soLuong ĐÃ TĂNG)
			maKH = "0000000000000000001"
			vaiTro = "quan_tri_he_thong" 
			chucVu = "Quản trị hệ thống"
			soLuong = 1 // Cập nhật lại số lượng để Root chiếm dòng tiếp theo
		} else {
			maKH = core.TaoMaKhachHangMoi(shopID)
			vaiTro = "khach_hang" 
			chucVu = "Khách hàng" 
		}
	} else {
		if soLuong == 0 {
			maKH = "0000000000000000001"
			vaiTro = "quan_tri_vien" 
			chucVu = "Quản trị viên"
		} else {
			maKH = core.TaoMaKhachHangMoi(shopID)
			vaiTro = "khach_hang"
			chucVu = "Khách hàng"
		}
	}

	passHash, _ := cau_hinh.HashMatKhau(pass)
	pinHash, _ := cau_hinh.HashMatKhau(maPin)
	
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
		TrangThai:      1, 
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
		Inbox:          make([]*core.TinNhan, 0), 
		NgayTao:        nowStr, 
		NguoiCapNhat:   user,
		NgayCapNhat:    nowStr, 
	}

	newKH.DongTrongSheet = core.DongBatDau_KhachHang + soLuong
	core.ThemKhachHangVaoRam(newKH)
	
	r := newKH.DongTrongSheet
	sh := "KHACH_HANG"
	
	ghi(shopID, sh, r, core.CotKH_MaKhachHang, newKH.MaKhachHang)
	ghi(shopID, sh, r, core.CotKH_TenDangNhap, newKH.TenDangNhap)
	ghi(shopID, sh, r, core.CotKH_Email, newKH.Email)
	ghi(shopID, sh, r, core.CotKH_MatKhauHash, newKH.MatKhauHash)
	ghi(shopID, sh, r, core.CotKH_MaPinHash, newKH.MaPinHash)
	ghi(shopID, sh, r, core.CotKH_VaiTroQuyenHan, newKH.VaiTroQuyenHan)
	ghi(shopID, sh, r, core.CotKH_ChucVu, newKH.ChucVu)
	ghi(shopID, sh, r, core.CotKH_TrangThai, newKH.TrangThai)
	ghi(shopID, sh, r, core.CotKH_TenKhachHang, newKH.TenKhachHang)
	ghi(shopID, sh, r, core.CotKH_DienThoai, newKH.DienThoai)
	ghi(shopID, sh, r, core.CotKH_NgaySinh, newKH.NgaySinh)
	ghi(shopID, sh, r, core.CotKH_GioiTinh, newKH.GioiTinh)
	ghi(shopID, sh, r, core.CotKH_NguonKhachHang, newKH.NguonKhachHang)
	ghi(shopID, sh, r, core.CotKH_NgayTao, newKH.NgayTao)
	ghi(shopID, sh, r, core.CotKH_NguoiCapNhat, newKH.NguoiCapNhat)
	ghi(shopID, sh, r, core.CotKH_NgayCapNhat, newKH.NgayCapNhat)
	ghi(shopID, sh, r, core.CotKH_NguoiCapNhat, newKH.NguoiCapNhat)
	ghi(shopID, sh, r, core.CotKH_NgayCapNhat, newKH.NgayCapNhat)

	// ==============================================================
	// TỰ ĐỘNG GỬI TIN NHẮN CHÀO MỪNG TỪ TRỢ LÝ ẢO 99K
	// ==============================================================
	msgID := fmt.Sprintf("AUTO_%d_0000000000000000000", time.Now().UnixNano())
	welcomeMsg := &core.TinNhan{
		MaTinNhan:   msgID,
		LoaiTinNhan: "AUTO",
		NguoiGuiID:  "0000000000000000000", // Gửi từ Bot Hệ Thống
		NguoiNhanID: maKH,                  // Gửi cho người vừa đăng ký
		TieuDe:      "Chào mừng gia nhập Nền tảng 99K",
		NoiDung:     "Chào mừng " + hoTen + " đến với hệ thống 99K.vn! Cửa hàng của bạn đã được khởi tạo thành công. Nếu cần hỗ trợ trong quá trình sử dụng, bạn có thể phản hồi trực tiếp tại cuộc trò chuyện này.",
		NgayTao:     nowStr,
	}
	core.ThemMoiTinNhan(shopID, welcomeMsg)

	if theme == "theme_master" && vaiTro != "quan_tri_he_thong" {
		code := core.TaoMaOTP6So() 
		core.LuuOTP(shopID + "_" + user, code) 
	}

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

	if vaiTro == "quan_tri_he_thong" || vaiTro == "quan_tri_vien_he_thong" || vaiTro == "quan_tri_vien" {
		c.Redirect(http.StatusFound, "/master/tong-quan")
	} else if theme == "theme_master" {
		c.Redirect(http.StatusFound, "/admin/tong-quan") 
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

func API_XacThucKichHoat(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID") 
	otp := strings.TrimSpace(c.PostForm("otp"))

	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok || !core.KiemTraOTP(masterShopID + "_" + kh.TenDangNhap, otp) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP không đúng hoặc đã hết hạn!"})
		return
	}

	loc := time.FixedZone("ICT", 7*3600)
	nowVN := time.Now().In(loc)

	core.KhoaHeThong.Lock()
	kh.GoiDichVu = append(kh.GoiDichVu, core.PlanInfo{
		MaGoi:      "TRIAL_3DAYS",
		TenGoi:     "Dùng thử 3 ngày",
		NgayHetHan: nowVN.AddDate(0, 0, 3).Format("2006-01-02 15:04:05"), 
		TrangThai:  "active",
	})
	core.KhoaHeThong.Unlock()

	core.ThemVaoHangCho(masterShopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_GoiDichVuJson, core.ToJSON(kh.GoiDichVu))

	go func(sub string) {
		TuDongThemSubdomain(sub)
	}(kh.TenDangNhap)

	c.JSON(200, gin.H{"status": "ok", "msg": "Xác thực thành công! Hệ thống đang khởi tạo Tên miền."})
}

func TuDongThemSubdomain(subdomain string) error {
	ctx := context.Background()
	jsonKey := cau_hinh.BienCauHinh.GoogleAuthJson 
	if jsonKey == "" { return nil }
	
	srv, err := run.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil { return err }
	
	fullDomain := subdomain + ".99k.vn"
	parent := "projects/project-47337221-fda1-48c7-b2f/locations/asia-southeast1" 

	req := &run.DomainMapping{
		Metadata: &run.ObjectMeta{ Name: fullDomain },
		Spec: &run.DomainMappingSpec{ RouteName: "maytinhhaiduong", CertificateMode: "AUTOMATIC" },
	}

	_, err = srv.Namespaces.Domainmappings.Create(parent, req).Do()
	if err != nil {
		log.Printf("❌ Lỗi tạo subdomain %s: %v", fullDomain, err)
		return err
	}
	
	log.Printf("✅ Đã lệnh cho Google tạo subdomain: %s", fullDomain)
	return nil
}
