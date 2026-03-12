package auth_admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/core"
	"app/modules/auth_verify" // Tích hợp module xác thực dùng chung
)

type Service struct { repo Repo }

func (s *Service) Login(dinhDanh, pass, userAgent string, ghiNho bool) (string, string, error) {
	adminID := config.BienCauHinh.IdFileSheetAdmin
	if err := core.EnsureKhachHangLoaded(adminID); err != nil { return "", "", err }

	kh, ok := s.repo.FindByUserOrEmail(dinhDanh)
	if !ok { return "", "", errors.New("Tài khoản không tồn tại trên hệ thống!") }
	if !config.KiemTraMatKhau(pass, kh.BaoMat.MatKhauHash) { return "", "", errors.New("Mật khẩu không đúng!") }
	if kh.TrangThai == 0 { return "", "", errors.New("Tài khoản bị khóa!") }

	thoiGianSong := config.ThoiGianHetHanCookie
	if ghiNho { thoiGianSong = 30 * 24 * time.Hour }

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	expTime := time.Now().Add(thoiGianSong).Unix()

	lock := core.GetSheetLock(adminID, core.TenSheetKhachHangAdmin)
	lock.Lock()
	if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]core.TenantDeviceToken) }
	
	// Dọn rác Session cũ
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

	go s.repo.UpdateUserJSON(kh.DongTrongSheet, jsonStr) 
	return sessionID, signature, nil
}

func (s *Service) Register(hoTen, user, email, pass, maPin, dienThoai, userAgent string) (string, string, error) {
	adminID := config.BienCauHinh.IdFileSheetAdmin
	if err := core.EnsureKhachHangLoaded(adminID); err != nil { return "", "", err }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	if _, ok := s.repo.FindByUserOrEmail(user); ok { return "", "", errors.New("Tên đăng nhập đã tồn tại!") }
	if _, ok := s.repo.FindByUserOrEmail(email); ok { return "", "", errors.New("Email đã được sử dụng!") }

	maKH := core.TaoMaKhachHangMoi(adminID)
	passHash, _ := config.HashMatKhau(pass)
	pinHash, _ := config.HashMatKhau(maPin)
	nowUnix := time.Now().Unix()
	soLuongUser := s.repo.CountUsers()
	
	newKH := &core.KhachHang{
		SpreadsheetID: adminID, DongTrongSheet: core.DongBatDau_KhachHang + soLuongUser,
		Version: 1, MaKhachHang: maKH, TenDangNhap: user, Email: email,
		BaoMat: core.TenantBaoMat{MatKhauHash: passHash, MaPinHash: pinHash},
		RefreshTokens: make(map[string]core.TenantDeviceToken),
		VaiTroQuyenHan: "khach_hang", ChucVu: "Chủ Cửa Hàng", TrangThai: 1,
		GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
		CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
		ThongTin: core.TenantThongTin{NguonKhachHang: "web_saas_register", TenKhachHang: hoTen, DienThoai: dienThoai},
		NgayTao: nowUnix, NguoiCapNhat: user, NgayCapNhat: nowUnix, 
	}

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	newKH.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: time.Now().Add(config.ThoiGianHetHanCookie).Unix(), Created: nowUnix}
	
	s.repo.InsertUser(newKH)

	s.repo.SendWelcomeMessage(&core.TinNhan{
		MaTinNhan: fmt.Sprintf("MSG_%d_000", time.Now().UnixNano()), LoaiTinNhan: "AUTO",
		NguoiGuiID: "0000000000000000000", NguoiNhanID: []string{maKH}, TieuDe: "Chào mừng Chủ Shop",
		NoiDung: "Chào mừng " + hoTen + " đến với hệ thống 99K.VN! Hãy bắt đầu thiết lập Cửa hàng của bạn ngay.", 
		NgayTao: nowUnix, NguoiDoc: []string{}, TrangThaiXoa: []string{},
	})

	return sessionID, signature, nil
}

// Hàm đổi mật khẩu mới (Sử dụng sau khi đã Verify PIN/OTP thành công)
func (s *Service) ResetPassword(dinhDanh, passMoi string) error {
	adminID := config.BienCauHinh.IdFileSheetAdmin
	kh, ok := s.repo.FindByUserOrEmail(dinhDanh)
	if !ok { return errors.New("Không tìm thấy tài khoản") }

	hash, _ := config.HashMatKhau(passMoi)
	lock := core.GetSheetLock(adminID, core.TenSheetKhachHangAdmin)
	lock.Lock()
	kh.BaoMat.MatKhauHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()
	
	s.repo.UpdateUserJSON(kh.DongTrongSheet, jsonStr)
	return nil
}
