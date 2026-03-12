package auth_store

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"app/config"
	"app/core"
)

type Service struct { repo Repo }

func (s *Service) Login(shopID, dinhDanh, pass, userAgent string, ghiNho bool) (string, string, error) {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", err }

	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return "", "", errors.New("Tài khoản không tồn tại trên cửa hàng này!") }
	if !config.KiemTraMatKhau(pass, kh.BaoMat.MatKhauHash) { return "", "", errors.New("Mật khẩu không đúng!") }
	if kh.TrangThai == 0 { return "", "", errors.New("Tài khoản của bạn đã bị chủ shop khóa!") }

	thoiGianSong := config.ThoiGianHetHanCookie
	if ghiNho { thoiGianSong = 30 * 24 * time.Hour }

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	expTime := time.Now().Add(thoiGianSong).Unix()

	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]core.TenantDeviceToken) }
	
	// Dọn rác Session
	nowUnix := time.Now().Unix()
	for key, info := range kh.RefreshTokens { if info.Exp < nowUnix { delete(kh.RefreshTokens, key) } }
	if len(kh.RefreshTokens) >= 5 {
		var oldestKey string; var oldestTime int64 = 1<<63 - 1
		for key, info := range kh.RefreshTokens { if info.Exp < oldestTime { oldestTime = info.Exp; oldestKey = key } }
		if oldestKey != "" { delete(kh.RefreshTokens, oldestKey) }
	}
	
	kh.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: expTime, Created: nowUnix}
	
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()

	go s.repo.UpdateUserJSON(shopID, kh.DongTrongSheet, jsonStr) 
	return sessionID, signature, nil
}

func (s *Service) Register(shopID, hoTen, user, email, pass, maPin, dienThoai, userAgent string) (string, string, error) {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", err }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	if _, ok := s.repo.FindByUserOrEmail(shopID, user); ok { return "", "", errors.New("Tên đăng nhập đã tồn tại trên hệ thống của Shop!") }
	if _, ok := s.repo.FindByUserOrEmail(shopID, email); ok { return "", "", errors.New("Email đã được sử dụng!") }

	maKH := core.TaoMaKhachHangMoi(shopID)
	passHash, _ := config.HashMatKhau(pass)
	pinHash, _ := config.HashMatKhau(maPin)
	nowUnix := time.Now().Unix()
	soLuongUser := s.repo.CountUsers(shopID)
	
	// Gắn Quyền "Khách Lẻ" cho User mới
	newKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + soLuongUser,
		Version: 1, MaKhachHang: maKH, TenDangNhap: user, Email: email,
		BaoMat: core.TenantBaoMat{MatKhauHash: passHash, MaPinHash: pinHash},
		RefreshTokens: make(map[string]core.TenantDeviceToken),
		VaiTroQuyenHan: "khach_le", ChucVu: "Khách Lẻ", TrangThai: 1,
		ThongTin: core.TenantThongTin{NguonKhachHang: "web_store_register", TenKhachHang: hoTen, DienThoai: dienThoai},
		NgayTao: nowUnix, NguoiCapNhat: user, NgayCapNhat: nowUnix, 
	}

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	newKH.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: time.Now().Add(config.ThoiGianHetHanCookie).Unix(), Created: nowUnix}
	
	s.repo.InsertUser(shopID, newKH)
	return sessionID, signature, nil
}

func (s *Service) ResetPassword(shopID, dinhDanh, passMoi string) error {
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return errors.New("Không tìm thấy tài khoản") }

	hash, _ := config.HashMatKhau(passMoi)
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.BaoMat.MatKhauHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()
	
	s.repo.UpdateUserJSON(shopID, kh.DongTrongSheet, jsonStr)
	return nil
}
