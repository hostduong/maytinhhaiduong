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
	
	kh.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: expTime, Created: nowUnix}
	
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()

	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	go s.repo.UpdateUserJSON(sID, kh.DongTrongSheet, jsonStr) 
	
	return sessionID, signature, nil
}

func (s *Service) Register(appMode, shopID, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, string, error) {
	// 1. TƯỜNG LỬA CHỐNG LỖI MẠNG (Đảm bảo RAM đồng bộ 100% với Sheet)
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { 
		// Tải không thành công do lỗi mạng, lỗi server... => Khóa chặn không cho đăng ký để bảo vệ ID 001
		return "", "", "", errors.New("Hệ thống đang bận hoặc gián đoạn kết nối dữ liệu gốc. Vui lòng thử lại sau!") 
	}

	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	if _, ok := s.repo.FindByUserOrEmail(shopID, user); ok { return "", "", "", errors.New("Tên đăng nhập đã tồn tại!") }
	if _, ok := s.repo.FindByUserOrEmail(shopID, email); ok { return "", "", "", errors.New("Email đã được sử dụng!") }

	var maKH, vaiTro, chucVu, nguon string
	isFirstMaster := false

	// 2. ĐỊNH HƯỚNG QUYỀN LỰC CHUẨN MỰC
	// Kiểm tra xem hệ thống đã có Chúa tể (001) chưa
	_, hasGod := core.LayKhachHang(shopID, "0000000000000000001")
	
	if !hasGod {
		// TRƯỜNG HỢP A: NGƯỜI ĐẦU TIÊN (Chỉ kích hoạt 1 lần duy nhất trong đời hệ thống)
		maKH = "0000000000000000001"
		vaiTro = "quan_tri_he_thong"
		chucVu = "Quản trị hệ thống"
		isFirstMaster = true 
		if appMode == "MASTER_CORE" { nguon = "master_core_register" } else { nguon = "web_saas_register" }
	} else {
		// TRƯỜNG HỢP B: TỪ NGƯỜI THỨ 2 TRỞ ĐI CHẮC CHẮN LÀ KHÁCH HÀNG (An toàn tuyệt đối)
		maKH = core.TaoMaKhachHangMoi(shopID)
		vaiTro = "khach_hang"
		chucVu = "Khách hàng"
		
		if appMode == "MASTER_CORE" { 
			nguon = "master_core_register" 
		} else if appMode == "TENANT_ADMIN" { 
			nguon = "web_saas_register" 
		} else { 
			nguon = "web_store_register" 
		}
	}

	passHash, _ := config.HashMatKhau(pass)
	pinHash, _ := config.HashMatKhau(maPin)
	nowUnix := time.Now().Unix()

	// 3. ĐỒNG THỜI SINH BOT NẾU LÀ NGƯỜI ĐẦU TIÊN
	if isFirstMaster {
		soLuongBot := s.repo.CountUsers(shopID)
		botKH := &core.KhachHang{
			SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + soLuongBot,
			Version: 1, MaKhachHang: "0000000000000000000", TenDangNhap: "admin", Email: "bot@99k.vn",
			BaoMat: core.TenantBaoMat{MatKhauHash: "", MaPinHash: ""}, // Không có pass, pin => Cấm đăng nhập
			RefreshTokens: make(map[string]core.TenantDeviceToken),
			VaiTroQuyenHan: "quan_tri_vien_he_thong", ChucVu: "Hệ thống", TrangThai: 1,
			GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
			CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
			// Set giới tính là Nam (1) theo yêu cầu
			ThongTin: core.TenantThongTin{NguonKhachHang: "system_bot", TenKhachHang: "Hệ thống", GioiTinh: 1},
			NganHang: core.TenantNganHang{}, MangXaHoi: make(map[string]string),
			NgayTao: nowUnix, NguoiCapNhat: "system", NgayCapNhat: nowUnix, 
		}
		s.repo.InsertUser(shopID, botKH)
	}

	// 4. LƯU TÀI KHOẢN NGƯỜI DÙNG THỰC TẾ
	soLuongUser := s.repo.CountUsers(shopID)
	
	newKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + soLuongUser,
		Version: 1, MaKhachHang: maKH, TenDangNhap: user, Email: email,
		BaoMat: core.TenantBaoMat{MatKhauHash: passHash, MaPinHash: pinHash},
		RefreshTokens: make(map[string]core.TenantDeviceToken),
		VaiTroQuyenHan: vaiTro, ChucVu: chucVu, TrangThai: 1,
		GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
		CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
		ThongTin: core.TenantThongTin{NguonKhachHang: nguon, TenKhachHang: hoTen, DienThoai: dienThoai, NgaySinh: ngaySinh, GioiTinh: gioiTinh},
		NganHang: core.TenantNganHang{}, MangXaHoi: make(map[string]string),
		NgayTao: nowUnix, NguoiCapNhat: user, NgayCapNhat: nowUnix, 
	}

	sessionID := config.TaoSessionIDAnToan(); signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	newKH.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: time.Now().Add(config.ThoiGianHetHanCookie).Unix(), Created: nowUnix}
	
	s.repo.InsertUser(shopID, newKH)

	// Nếu là Cửa hàng hoặc Master -> Gắn tin nhắn chào mừng
	if appMode == "TENANT_ADMIN" || appMode == "MASTER_CORE" {
		s.repo.SendWelcomeMessage(shopID, &core.TinNhan{
			MaTinNhan: fmt.Sprintf("MSG_%d_000", time.Now().UnixNano()), LoaiTinNhan: "AUTO",
			NguoiGuiID: "0000000000000000000", NguoiNhanID: []string{maKH}, TieuDe: "Chào mừng bạn",
			NoiDung: "Chào mừng " + hoTen + " đến với hệ thống 99K.VN! Nếu cần hỗ trợ, bạn có thể gửi tin nhắn tại đây.", 
			NgayTao: nowUnix, NguoiDoc: []string{}, TrangThaiXoa: []string{}, DinhKem: []core.FileDinhKem{}, ThamChieuID: []string{},
		})
	}

	return sessionID, signature, vaiTro, nil
}

// ================= CÁC HÀM XỬ LÝ QUÊN MẬT KHẨU =================

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
