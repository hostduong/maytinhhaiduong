package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/core"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

type Service struct { repo Repo }

func (s *Service) Login(shopID, dinhDanh, pass, userAgent string, ghiNho bool) (string, string, error) {
	if dinhDanh == "admin" { return "", "", errors.New("Tài khoản không tồn tại!") }
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return "", "", errors.New("Tài khoản không tồn tại!") }
	if !config.KiemTraMatKhau(pass, kh.MatKhauHash) { return "", "", errors.New("Mật khẩu không đúng!") }
	if kh.TrangThai == 0 { return "", "", errors.New("Tài khoản bị khóa!") }

	thoiGianSong := config.ThoiGianHetHanCookie
	if ghiNho { thoiGianSong = 30 * 24 * time.Hour }

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	expTime := time.Now().Add(thoiGianSong).Unix()

	core.KhoaHeThong.Lock()
	if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]core.TokenInfo) }
	
	// Dọn dẹp & Giới hạn 5 Thiết bị
	nowUnix := time.Now().Unix()
	for key, info := range kh.RefreshTokens { if info.ExpiresAt < nowUnix { delete(kh.RefreshTokens, key) } }
	if len(kh.RefreshTokens) >= 5 {
		var oldestKey string; var oldestTime int64 = 1<<63 - 1
		for key, info := range kh.RefreshTokens { if info.ExpiresAt < oldestTime { oldestTime = info.ExpiresAt; oldestKey = key } }
		if oldestKey != "" { delete(kh.RefreshTokens, oldestKey) }
	}
	kh.RefreshTokens[sessionID] = core.TokenInfo{DeviceName: userAgent, ExpiresAt: expTime}
	core.KhoaHeThong.Unlock()

	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	go s.repo.UpdateTokens(sID, kh.DongTrongSheet, kh.RefreshTokens)
	return sessionID, signature, nil
}

func (s *Service) Register(shopID, theme, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, string, error) {
	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	if _, ok := s.repo.FindByUserOrEmail(shopID, user); ok { return "", "", "", errors.New("Tên đăng nhập đã tồn tại!") }
	if _, ok := s.repo.FindByUserOrEmail(shopID, email); ok { return "", "", "", errors.New("Email đã được sử dụng!") }

	soLuong := s.repo.CountUsers(shopID)
	var maKH, vaiTro, chucVu string
	loc := time.FixedZone("ICT", 7*3600); nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")

	if theme == "theme_master" {
		if soLuong == 0 {
			bot := &core.KhachHang{
				SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang, MaKhachHang: "0000000000000000000",
				TenDangNhap: "admin", VaiTroQuyenHan: "admin", ChucVu: "Hệ thống", TenKhachHang: "Trợ lý ảo 99K",
				TrangThai: 1, NgayTao: nowStr, NgayCapNhat: nowStr,
			}
			s.repo.InsertUser(bot)
			maKH = "0000000000000000001"; vaiTro = "quan_tri_he_thong"; chucVu = "Quản trị hệ thống"; soLuong = 1
		} else { maKH = core.TaoMaKhachHangMoi(shopID); vaiTro = "khach_hang"; chucVu = "Khách hàng" }
	} else {
		if soLuong == 0 { maKH = "0000000000000000001"; vaiTro = "quan_tri_vien"; chucVu = "Quản trị viên" } else { maKH = core.TaoMaKhachHangMoi(shopID); vaiTro = "khach_hang"; chucVu = "Khách hàng" }
	}

	passHash, _ := config.HashMatKhau(pass); pinHash, _ := config.HashMatKhau(maPin)
	newKH := &core.KhachHang{
		SpreadsheetID: shopID, MaKhachHang: maKH, TenDangNhap: user, Email: email, MatKhauHash: passHash, MaPinHash: pinHash,
		RefreshTokens: make(map[string]core.TokenInfo), VaiTroQuyenHan: vaiTro, ChucVu: chucVu, TrangThai: 1,
		GoiDichVu: make([]core.PlanInfo, 0), CauHinh: core.UserConfig{Theme: "light", Language: "vi"},
		NguonKhachHang: "web_register", TenKhachHang: hoTen, DienThoai: dienThoai, NgaySinh: ngaySinh, GioiTinh: gioiTinh,
		NgayTao: nowStr, NguoiCapNhat: user, NgayCapNhat: nowStr, DongTrongSheet: core.DongBatDau_KhachHang + soLuong,
	}

	sessionID := config.TaoSessionIDAnToan(); signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	newKH.RefreshTokens[sessionID] = core.TokenInfo{DeviceName: userAgent, ExpiresAt: time.Now().Add(config.ThoiGianHetHanCookie).Unix()}
	s.repo.InsertUser(newKH)

	s.repo.SendWelcomeMessage(shopID, &core.TinNhan{
		MaTinNhan: fmt.Sprintf("AUTO_%d_000", time.Now().UnixNano()), LoaiTinNhan: "AUTO",
		NguoiGuiID: "0000000000000000000", NguoiNhanID: maKH, TieuDe: "Chào mừng gia nhập Nền tảng 99K",
		NoiDung: "Chào mừng " + hoTen + " đến với hệ thống 99K.vn! Nếu cần hỗ trợ, bạn có thể phản hồi trực tiếp tại đây.", NgayTao: nowStr,
	})

	if theme == "theme_master" && vaiTro != "quan_tri_he_thong" { core.LuuOTP(shopID+"_"+user, core.TaoMaOTP6So()) }
	return sessionID, signature, vaiTro, nil
}

func (s *Service) VerifyOTPAndActivate(shopID, userID, otp string) error {
	kh, ok := s.repo.FindByUserOrEmail(shopID, userID)
	if !ok || !core.KiemTraOTP(shopID+"_"+kh.TenDangNhap, otp) { return errors.New("Mã OTP không đúng hoặc đã hết hạn!") }

	core.KhoaHeThong.Lock()
	kh.GoiDichVu = append(kh.GoiDichVu, core.PlanInfo{MaGoi: "TRIAL_3DAYS", TenGoi: "Dùng thử 3 ngày", NgayHetHan: time.Now().In(time.FixedZone("ICT", 7*3600)).AddDate(0, 0, 3).Format("2006-01-02 15:04:05"), TrangThai: "active"})
	core.KhoaHeThong.Unlock()
	s.repo.UpdateGoiDichVu(shopID, kh.DongTrongSheet, kh.GoiDichVu)
	go s.CreateSubdomain(kh.TenDangNhap)
	return nil
}

func (s *Service) SendOtp(shopID, dinhDanh string) error {
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return nil } // Giả vờ thành công để chống dò tài khoản
	if kh.Email == "" || !strings.Contains(kh.Email, "@") { return errors.New("Tài khoản này chưa cập nhật Email, vui lòng dùng PIN.") }

	okLimit, msg := core.KiemTraRateLimit(kh.Email)
	if !okLimit { return errors.New(msg) }

	code := core.TaoMaOTP6So()
	if err := core.GuiMailXacMinhAPI(kh.Email, code); err != nil { return errors.New("Lỗi hệ thống gửi mail: " + err.Error()) }
	core.LuuOTP(shopID+"_"+kh.TenDangNhap, code)
	return nil
}

func (s *Service) ResetByOtp(shopID, dinhDanh, otp, passMoi string) error {
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok || !core.KiemTraOTP(shopID+"_"+kh.TenDangNhap, otp) { return errors.New("Mã OTP không đúng hoặc đã hết hạn!") }
	hash, _ := config.HashMatKhau(passMoi); core.KhoaHeThong.Lock(); kh.MatKhauHash = hash; core.KhoaHeThong.Unlock()
	s.repo.UpdatePassword(shopID, kh.DongTrongSheet, hash); return nil
}

func (s *Service) ResetByPin(shopID, dinhDanh, pinInput, passMoi string) error {
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok || !config.KiemTraMatKhau(pinInput, kh.MaPinHash) { return errors.New("Tài khoản hoặc mã PIN không chính xác!") }
	hash, _ := config.HashMatKhau(passMoi); core.KhoaHeThong.Lock(); kh.MatKhauHash = hash; core.KhoaHeThong.Unlock()
	s.repo.UpdatePassword(shopID, kh.DongTrongSheet, hash); return nil
}

func (s *Service) CreateSubdomain(subdomain string) error {
	ctx := context.Background(); jsonKey := config.BienCauHinh.GoogleAuthJson
	if jsonKey == "" { return nil }
	srv, err := run.NewService(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil { return err }
	fullDomain := subdomain + ".99k.vn"
	parent := "projects/project-47337221-fda1-48c7-b2f/locations/asia-southeast1"
	_, _ = srv.Namespaces.Domainmappings.Create(parent, &run.DomainMapping{ Metadata: &run.ObjectMeta{Name: fullDomain}, Spec: &run.DomainMappingSpec{RouteName: "maytinhhaiduong", CertificateMode: "AUTOMATIC"} }).Do()
	return nil
}
