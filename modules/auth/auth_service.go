package auth

import (
	"encoding/json"
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
	if !config.KiemTraMatKhau(pass, kh.BaoMat.MatKhauHash) { return "", "", errors.New("Mật khẩu không đúng!") }
	if kh.TrangThai == 0 { return "", "", errors.New("Tài khoản bị khóa!") }

	thoiGianSong := config.ThoiGianHetHanCookie
	if ghiNho { thoiGianSong = 30 * 24 * time.Hour }

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	expTime := time.Now().Add(thoiGianSong).Unix()

	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]core.TenantDeviceToken) }
	
	nowUnix := time.Now().Unix()
	for key, info := range kh.RefreshTokens { if info.Exp < nowUnix { delete(kh.RefreshTokens, key) } }
	
	if len(kh.RefreshTokens) >= 5 {
		var oldestKey string; var oldestTime int64 = 1<<63 - 1
		for key, info := range kh.RefreshTokens { if info.Exp < oldestTime { oldestTime = info.Exp; oldestKey = key } }
		if oldestKey != "" { delete(kh.RefreshTokens, oldestKey) }
	}
	
	// Lưu Token theo cấu trúc mới
	kh.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: expTime, Created: nowUnix}
	
	// Cập nhật JSON
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()

	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	go s.repo.UpdateUserJSON(sID, kh.DongTrongSheet, jsonStr) 
	
	return sessionID, signature, nil
}

func (s *Service) Register(appMode, shopID, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, string, error) {
	if appMode == "MASTER_CORE" { return "", "", "", errors.New("Truy cập bị từ chối! Vùng này không cho phép đăng ký ngoại tuyến.") }
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", "", err }

	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	if _, ok := s.repo.FindByUserOrEmail(shopID, user); ok { return "", "", "", errors.New("Tên đăng nhập đã tồn tại!") }
	if _, ok := s.repo.FindByUserOrEmail(shopID, email); ok { return "", "", "", errors.New("Email đã được sử dụng!") }

	maKH := core.TaoMaKhachHangMoi(shopID)
	var vaiTro, chucVu, nguon string
	if appMode == "TENANT_ADMIN" {
		vaiTro = "khach_hang"; chucVu = "Khách Hàng"; nguon = "web_saas_register"
	} else if appMode == "STOREFRONT" {
		vaiTro = "khach_le"; chucVu = "Khách Lẻ"; nguon = "web_store_register"
	}

	passHash, _ := config.HashMatKhau(pass); pinHash, _ := config.HashMatKhau(maPin)
	soLuong := s.repo.CountUsers(shopID)
	nowUnix := time.Now().Unix()

	// Khởi tạo Struct Đỉnh cao
	newKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + soLuong,
		Version: 1, MaKhachHang: maKH, TenDangNhap: user, Email: email,
		BaoMat: core.TenantBaoMat{MatKhauHash: passHash, MaPinHash: pinHash},
		RefreshTokens: make(map[string]core.TenantDeviceToken),
		VaiTroQuyenHan: vaiTro, ChucVu: chucVu, TrangThai: 1,
		GoiDichVu: make([]core.TenantGoiDichVu, 0),
		Modules: make(map[string]bool),
		CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
		ThongTin: core.TenantThongTin{NguonKhachHang: nguon, TenKhachHang: hoTen, DienThoai: dienThoai, NgaySinh: ngaySinh, GioiTinh: gioiTinh},
		NganHang: core.TenantNganHang{},
		MangXaHoi: make(map[string]string),
		NgayTao: nowUnix, NguoiCapNhat: user, NgayCapNhat: nowUnix, 
	}

	sessionID := config.TaoSessionIDAnToan(); signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	newKH.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: time.Now().Add(config.ThoiGianHetHanCookie).Unix(), Created: nowUnix}
	
	s.repo.InsertUser(shopID, newKH)

	if appMode == "TENANT_ADMIN" {
		s.repo.SendWelcomeMessage(shopID, &core.TinNhan{
			MaTinNhan: fmt.Sprintf("MSG_%d_000", time.Now().UnixNano()), LoaiTinNhan: "AUTO",
			NguoiGuiID: "0000000000000000000", NguoiNhanID: []string{maKH}, TieuDe: "Chào mừng bạn",
			NoiDung: "Chào mừng " + hoTen + " đến với hệ thống 99K.VN! Nếu cần hỗ trợ, bạn có thể gửi tin nhắn tại đây.", 
			NgayTao: nowUnix, NguoiDoc: []string{}, TrangThaiXoa: []string{}, DinhKem: []core.FileDinhKem{}, ThamChieuID: []string{},
		})
	}

	return sessionID, signature, vaiTro, nil
}

// CÁC HÀM RESET
func (s *Service) SendOtp(shopID, dinhDanh string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return nil } 
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
	lock.Lock()
	kh.BaoMat.MatKhauHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()
	
	s.repo.UpdateUserJSON(shopID, kh.DongTrongSheet, jsonStr); return nil
}

func (s *Service) ResetByPin(shopID, dinhDanh, pinInput, passMoi string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok || !config.KiemTraMatKhau(pinInput, kh.BaoMat.MaPinHash) { return errors.New("Tài khoản hoặc mã PIN không chính xác!") }
	hash, _ := config.HashMatKhau(passMoi)
	
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.BaoMat.MatKhauHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()
	
	s.repo.UpdateUserJSON(shopID, kh.DongTrongSheet, jsonStr); return nil
}
