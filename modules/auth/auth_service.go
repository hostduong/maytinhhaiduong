package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/core"
)

type Service struct { repo Repo }

func (s *Service) Login(shopID, dinhDanh, pass, userAgent string, ghiNho bool) (string, string, error) {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", err }

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

	// Khóa độc quyền Sheet Khách Hàng để cập nhật Token
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]core.TokenInfo) }
	
	nowUnix := time.Now().Unix()
	for key, info := range kh.RefreshTokens { if info.ExpiresAt < nowUnix { delete(kh.RefreshTokens, key) } }
	
	// Giới hạn 5 thiết bị
	if len(kh.RefreshTokens) >= 5 {
		var oldestKey string; var oldestTime int64 = 1<<63 - 1
		for key, info := range kh.RefreshTokens { if info.ExpiresAt < oldestTime { oldestTime = info.ExpiresAt; oldestKey = key } }
		if oldestKey != "" { delete(kh.RefreshTokens, oldestKey) }
	}
	kh.RefreshTokens[sessionID] = core.TokenInfo{DeviceName: userAgent, ExpiresAt: expTime}
	lock.Unlock() // Mở khóa siêu tốc

	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	go s.repo.UpdateTokens(sID, kh.DongTrongSheet, kh.RefreshTokens) // Đẩy Google Sheet chạy ngầm
	
	return sessionID, signature, nil
}

func (s *Service) Register(appMode, shopID, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, string, error) {
	// 1. Chặn đăng ký ở Vùng Tuyệt Mật
	if appMode == "MASTER_CORE" {
		return "", "", "", errors.New("Truy cập bị từ chối! Vùng này không cho phép đăng ký ngoại tuyến.")
	}

	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", "", err }

	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	if _, ok := s.repo.FindByUserOrEmail(shopID, user); ok { return "", "", "", errors.New("Tên đăng nhập đã tồn tại!") }
	if _, ok := s.repo.FindByUserOrEmail(shopID, email); ok { return "", "", "", errors.New("Email đã được sử dụng!") }

	// 2. Logic tạo Khách hàng (Dân thường)
	maKH := core.TaoMaKhachHangMoi(shopID)
	var vaiTro, chucVu, nguon string
	
	if appMode == "TENANT_ADMIN" {
		vaiTro = "khach_hang" // Khách của hệ thống
		chucVu = "Khách Hàng"
		nguon = "web_saas_register"
	} else if appMode == "STOREFRONT" {
		vaiTro = "khach_le" // Khách của cửa hàng
		chucVu = "Khách Lẻ"
		nguon = "web_store_register"
	}

	loc := time.FixedZone("ICT", 7*3600); nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")
	passHash, _ := config.HashMatKhau(pass); pinHash, _ := config.HashMatKhau(maPin)
	soLuong := s.repo.CountUsers(shopID)

	newKH := &core.KhachHang{
		SpreadsheetID: shopID, MaKhachHang: maKH, TenDangNhap: user, Email: email, MatKhauHash: passHash, MaPinHash: pinHash,
		RefreshTokens: make(map[string]core.TokenInfo), VaiTroQuyenHan: vaiTro, ChucVu: chucVu, TrangThai: 1,
		GoiDichVu: make([]core.PlanInfo, 0), CauHinh: core.UserConfig{Theme: "light", Language: "vi"},
		NguonKhachHang: nguon, TenKhachHang: hoTen, DienThoai: dienThoai, NgaySinh: ngaySinh, GioiTinh: gioiTinh,
		NgayTao: nowStr, NguoiCapNhat: user, NgayCapNhat: nowStr, DongTrongSheet: core.DongBatDau_KhachHang + soLuong,
	}

	sessionID := config.TaoSessionIDAnToan(); signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	newKH.RefreshTokens[sessionID] = core.TokenInfo{DeviceName: userAgent, ExpiresAt: time.Now().Add(config.ThoiGianHetHanCookie).Unix()}
	
	// Repo sẽ tự khóa RAM để Insert
	s.repo.InsertUser(shopID, newKH)

	// 3. Tin nhắn chào mừng chỉ dành cho Khách đăng ký SaaS
	if appMode == "TENANT_ADMIN" {
		s.repo.SendWelcomeMessage(shopID, &core.TinNhan{
			MaTinNhan: fmt.Sprintf("AUTO_%d_000", time.Now().UnixNano()), LoaiTinNhan: "AUTO",
			NguoiGuiID: "0000000000000000000", NguoiNhanID: maKH, TieuDe: "Chào mừng bạn",
			NoiDung: "Chào mừng " + hoTen + " đến với hệ thống 99K.VN! Nếu cần hỗ trợ, bạn có thể gửi tin nhắn tại đây.", NgayTao: nowStr,
		})
	}

	return sessionID, signature, vaiTro, nil
}

// ================= CÁC HÀM XỬ LÝ QUÊN MẬT KHẨU =================

func (s *Service) SendOtp(shopID, dinhDanh string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }

	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return nil } // Giấu lỗi để chống scan tài khoản
	if kh.Email == "" || !strings.Contains(kh.Email, "@") { return errors.New("Tài khoản này chưa cập nhật Email, vui lòng dùng PIN.") }

	okLimit, msg := core.KiemTraRateLimit(kh.Email)
	if !okLimit { return errors.New(msg) }

	code := core.TaoMaOTP6So()
	if err := core.GuiMailXacMinhAPI(kh.Email, code); err != nil { return errors.New("Lỗi hệ thống gửi mail: " + err.Error()) }
	core.LuuOTP(shopID+"_"+kh.TenDangNhap, code)
	return nil
}

func (s *Service) ResetByOtp(shopID, dinhDanh, otp, passMoi string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }

	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok || !core.KiemTraOTP(shopID+"_"+kh.TenDangNhap, otp) { return errors.New("Mã OTP không đúng hoặc đã hết hạn!") }
	hash, _ := config.HashMatKhau(passMoi)
	
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock(); kh.MatKhauHash = hash; lock.Unlock()
	
	s.repo.UpdatePassword(shopID, kh.DongTrongSheet, hash); return nil
}

func (s *Service) ResetByPin(shopID, dinhDanh, pinInput, passMoi string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }

	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok || !config.KiemTraMatKhau(pinInput, kh.MaPinHash) { return errors.New("Tài khoản hoặc mã PIN không chính xác!") }
	hash, _ := config.HashMatKhau(passMoi)
	
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock(); kh.MatKhauHash = hash; lock.Unlock()
	
	s.repo.UpdatePassword(shopID, kh.DongTrongSheet, hash); return nil
}
