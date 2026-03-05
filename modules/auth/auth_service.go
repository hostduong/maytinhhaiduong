package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/config"
	"app/core"

	"golang.org/x/oauth2/google"
)

type Service struct { repo Repo }

func (s *Service) Login(shopID, dinhDanh, pass, userAgent string, ghiNho bool) (string, string, error) {
	// Check tường lửa: Đảm bảo RAM đã được nạp
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", err }

	if dinhDanh == "admin" { return "", "", errors.New("Tài khoản không tồn tại!") }
	
	// Repo đã tự dùng RLock
	kh, ok := s.repo.FindByUserOrEmail(shopID, dinhDanh)
	if !ok { return "", "", errors.New("Tài khoản không tồn tại!") }
	if !config.KiemTraMatKhau(pass, kh.MatKhauHash) { return "", "", errors.New("Mật khẩu không đúng!") }
	if kh.TrangThai == 0 { return "", "", errors.New("Tài khoản bị khóa!") }

	thoiGianSong := config.ThoiGianHetHanCookie
	if ghiNho { thoiGianSong = 30 * 24 * time.Hour }

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	expTime := time.Now().Add(thoiGianSong).Unix()

	// [LOCK CHUẨN MỰC]: Khóa độc quyền Sheet Khách Hàng của Shop để ghi Token
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
	lock.Unlock() // Mở khóa ngay lập tức!

	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	go s.repo.UpdateTokens(sID, kh.DongTrongSheet, kh.RefreshTokens) // Đẩy Google Sheet chạy ngầm
	
	return sessionID, signature, nil
}

func (s *Service) Register(shopID, theme, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, string, error) {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return "", "", "", err }

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

	if theme == "theme_master" || theme == "template_admin" {
		if soLuong == 0 {
			bot := &core.KhachHang{
				SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang, MaKhachHang: "0000000000000000000",
				TenDangNhap: "admin", VaiTroQuyenHan: "admin", ChucVu: "Hệ thống", TenKhachHang: "Trợ lý ảo 99K",
				TrangThai: 1, NgayTao: nowStr, NgayCapNhat: nowStr,
			}
			s.repo.InsertUser(shopID, bot)
			maKH = "0000000000000000001"; vaiTro = "quan_tri_he_thong"; chucVu = "Quản trị hệ thống"; soLuong = 1
		} else { maKH = core.TaoMaKhachHangMoi(shopID); vaiTro = "khach_hang"; chucVu = "Khách hàng" }
	} else {
		// Dành cho Web Cửa hàng (cuahang.99k.vn)
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
	
	// Repo sẽ tự khóa RAM để Insert
	s.repo.InsertUser(shopID, newKH)

	s.repo.SendWelcomeMessage(shopID, &core.TinNhan{
		MaTinNhan: fmt.Sprintf("AUTO_%d_000", time.Now().UnixNano()), LoaiTinNhan: "AUTO",
		NguoiGuiID: "0000000000000000000", NguoiNhanID: maKH, TieuDe: "Chào mừng gia nhập",
		NoiDung: "Chào mừng " + hoTen + " đến với hệ thống! Nếu cần hỗ trợ, bạn có thể phản hồi trực tiếp tại đây.", NgayTao: nowStr,
	})

	// Sinh mã OTP nếu là khách đăng ký gói
	if (theme == "theme_master" || theme == "template_admin") && vaiTro != "quan_tri_he_thong" { 
		core.LuuOTP(shopID+"_"+user, core.TaoMaOTP6So()) 
	}
	return sessionID, signature, vaiTro, nil
}

func (s *Service) VerifyOTPAndActivate(shopID, userID, otp string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }

	kh, ok := s.repo.FindByUserOrEmail(shopID, userID)
	if !ok || !core.KiemTraOTP(shopID+"_"+kh.TenDangNhap, otp) { return errors.New("Mã OTP không đúng hoặc đã hết hạn!") }

	// [LOCK CHUẨN MỰC] Ghi gói Trial
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.GoiDichVu = append(kh.GoiDichVu, core.PlanInfo{MaGoi: "TRIAL_3DAYS", TenGoi: "Dùng thử 3 ngày", NgayHetHan: time.Now().In(time.FixedZone("ICT", 7*3600)).AddDate(0, 0, 3).Format("2006-01-02 15:04:05"), TrangThai: "active"})
	lock.Unlock()
	
	s.repo.UpdateGoiDichVu(shopID, kh.DongTrongSheet, kh.GoiDichVu)
	
	// Kích hoạt Subdomain ngầm
	go s.CreateSubdomainADC(kh.TenDangNhap)
	return nil
}

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

// =========================================================================
// CreateSubdomainADC: Dùng Cơ chế ADC tự xin quyền nội bộ Cloud Run
// (Xóa bỏ thư viện "run" cũ yêu cầu file JSON)
// =========================================================================
func (s *Service) CreateSubdomainADC(tenDangNhap string) {
	subdomain := fmt.Sprintf("%s.99k.vn", tenDangNhap)
	
	payload := map[string]interface{}{
		"apiVersion": "domains.cloudrun.com/v1",
		"kind":       "DomainMapping",
		"metadata": map[string]string{"name": subdomain},
		"spec":     map[string]interface{}{"routeName": "maytinhhaiduong", "certificateMode": "AUTOMATIC"},
	}
	
	body, _ := json.Marshal(payload)
	apiURL := "https://asia-southeast1-run.googleapis.com/apis/domains.cloudrun.com/v1/namespaces/project-47337221-fda1-48c7-b2f/domainmappings"
	
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Dùng ADC lấy Token ẩn
	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err == nil {
		if token, err := creds.TokenSource.Token(); err == nil {
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		}
	} else {
		fmt.Println("[AUTH] Cảnh báo: Chạy môi trường Local, sẽ bỏ qua cấp Domain.")
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil { return }
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Printf("✅ [AUTH] Kích hoạt thành công Subdomain Trial: %s\n", subdomain)
	} else if resp.StatusCode == 409 {
		fmt.Printf("⚡ [AUTH] Subdomain Trial %s đã tồn tại, bỏ qua.\n", subdomain)
	}
}
