package dang_ky_master

import (
	"encoding/json"
	"errors"
	"time"

	"app/config"
	"app/core"
)

func Service_DangKyMaster(shopID, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, error) {
	// [TƯỜNG LỬA CHỐNG LỖI MẠNG]: Nếu lỗi, chặn đứng không cho tạo ID 001
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { 
		return "", "", errors.New("Hệ thống gián đoạn kết nối. Không thể khởi tạo ngay lúc này!") 
	}

	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	core.KhoaHeThong.RLock()
	var exists bool
	for _, kh := range core.CacheKhachHang[shopID] {
		if kh.TenDangNhap == user || kh.Email == email { exists = true; break }
	}
	core.KhoaHeThong.RUnlock()
	if exists { return "", "", errors.New("Tên đăng nhập hoặc Email đã tồn tại!") }

	var maKH, vaiTro, chucVu string
	isFirstMaster := false

	// KIỂM TRA SỰ TỒN TẠI CỦA SÁNG LẬP VIÊN (ID: 001)
	_, hasGod := core.LayKhachHang(shopID, "0000000000000000001")
	
	if !hasGod {
		maKH = "0000000000000000001"
		vaiTro = "quan_tri_he_thong"
		chucVu = "Quản trị hệ thống"
		isFirstMaster = true 
	} else {
		maKH = core.TaoMaKhachHangMoi(shopID)
		vaiTro = "quan_tri_vien_he_thong"
		chucVu = "Quản trị viên hệ thống"
	}

	passHash, _ := config.HashMatKhau(pass)
	pinHash, _ := config.HashMatKhau(maPin)
	nowUnix := time.Now().Unix()

	lock := core.GetSheetLock(shopID, core.TenSheetKhachHangMaster)
	lock.Lock()
	defer lock.Unlock()

	// 1. SINH BOT NẾU LÀ NGƯỜI ĐẦU TIÊN
	if isFirstMaster {
		soLuongBot := len(core.CacheKhachHang[shopID])
		botKH := &core.KhachHang{
			SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + soLuongBot,
			Version: 1, MaKhachHang: "0000000000000000000", TenDangNhap: "admin", Email: "bot@99k.vn",
			BaoMat: core.TenantBaoMat{MatKhauHash: "", MaPinHash: ""}, 
			RefreshTokens: make(map[string]core.TenantDeviceToken),
			VaiTroQuyenHan: "quan_tri_vien_he_thong", ChucVu: "Hệ thống", TrangThai: 1,
			GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
			CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
			ThongTin: core.TenantThongTin{NguonKhachHang: "system_bot", TenKhachHang: "Hệ thống", GioiTinh: 1},
			NganHang: core.TenantNganHang{}, MangXaHoi: make(map[string]string),
			NgayTao: nowUnix, NguoiCapNhat: "system", NgayCapNhat: nowUnix, 
		}
		core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], botKH)
		core.CacheMapKhachHang[core.TaoCompositeKey(shopID, botKH.MaKhachHang)] = botKH
		bBot, _ := json.Marshal(botKH)
		core.PushAppend(shopID, core.TenSheetKhachHangMaster, []interface{}{botKH.MaKhachHang, string(bBot)})
	}

	// 2. SINH TÀI KHOẢN MASTER (HOẶC ADMIN)
	soLuongUser := len(core.CacheKhachHang[shopID])
	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	
	newKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + soLuongUser,
		Version: 1, MaKhachHang: maKH, TenDangNhap: user, Email: email,
		BaoMat: core.TenantBaoMat{MatKhauHash: passHash, MaPinHash: pinHash},
		RefreshTokens: map[string]core.TenantDeviceToken{
			sessionID: {DeviceID: sessionID, Dev: userAgent, Exp: time.Now().Add(config.ThoiGianHetHanCookie).Unix(), Created: nowUnix},
		},
		VaiTroQuyenHan: vaiTro, ChucVu: chucVu, TrangThai: 1,
		GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
		CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
		ThongTin: core.TenantThongTin{NguonKhachHang: "master_core_register", TenKhachHang: hoTen, DienThoai: dienThoai, NgaySinh: ngaySinh, GioiTinh: gioiTinh},
		NganHang: core.TenantNganHang{}, MangXaHoi: make(map[string]string),
		NgayTao: nowUnix, NguoiCapNhat: user, NgayCapNhat: nowUnix, 
	}

	core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], newKH)
	core.CacheMapKhachHang[core.TaoCompositeKey(shopID, newKH.MaKhachHang)] = newKH
	bUser, _ := json.Marshal(newKH)
	core.PushAppend(shopID, core.TenSheetKhachHangMaster, []interface{}{newKH.MaKhachHang, string(bUser)})

	return sessionID, signature, nil
}
