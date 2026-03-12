package setup

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/config"
	"app/core"
)

func Service_ThucThiSetup(shopID, hoTen, user, email, pass, maPin, dienThoai, ngaySinh, gioiTinhStr, userAgent string) (string, string, error) {
	// 1. Phải tải thành công Sheet mới cho đăng ký (Chống mất mạng sinh lỗi)
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { 
		return "", "", errors.New("Không thể kết nối đến Database. Vui lòng thử lại sau!") 
	}

	// [CHỐT CHẶN 2 - CHỐNG SPAM TUYỆT ĐỐI]
	// Khóa RAM lại để đảm bảo không có 2 luồng Spam lọt vào cùng lúc
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHangMaster)
	lock.Lock()
	defer lock.Unlock()

	_, hasGod := core.LayKhachHang(shopID, "0000000000000000001")
	if hasGod {
		return "", "", errors.New("Hệ thống đã được khởi tạo! Không thể đăng ký thêm tài khoản Sáng lập.")
	}

	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	if dienThoai == "" || !config.KiemTraHoTen(hoTen) || !config.KiemTraTenDangNhap(user) || !config.KiemTraEmail(email) || !config.KiemTraMaPin(maPin) || !config.KiemTraDinhDangMatKhau(pass) {
		return "", "", errors.New("Dữ liệu nhập vào không hợp lệ!")
	}
	
	passHash, _ := config.HashMatKhau(pass)
	pinHash, _ := config.HashMatKhau(maPin)
	nowUnix := time.Now().Unix()

	// 2. KHAI SINH BOT HỆ THỐNG (ID: 000)
	botKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang,
		Version: 1, MaKhachHang: "0000000000000000000", TenDangNhap: "admin", Email: "bot@99k.vn",
		BaoMat: core.TenantBaoMat{MatKhauHash: "", MaPinHash: ""}, // Không có pass/pin -> Khóa login
		RefreshTokens: make(map[string]core.TenantDeviceToken),
		VaiTroQuyenHan: "quan_tri_vien_he_thong", ChucVu: "Hệ thống", TrangThai: 1,
		GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
		CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
		ThongTin: core.TenantThongTin{NguonKhachHang: "system_bot", TenKhachHang: "Hệ thống", GioiTinh: 1}, // Nam
		NganHang: core.TenantNganHang{}, MangXaHoi: make(map[string]string),
		NgayTao: nowUnix, NguoiCapNhat: "system", NgayCapNhat: nowUnix, 
	}
	core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], botKH)
	core.CacheMapKhachHang[core.TaoCompositeKey(shopID, botKH.MaKhachHang)] = botKH
	bBot, _ := json.Marshal(botKH)
	core.PushAppend(shopID, core.TenSheetKhachHangMaster, []interface{}{botKH.MaKhachHang, string(bBot)})

	// 3. KHAI SINH SÁNG LẬP VIÊN (ID: 001)
	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	
	newKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang + 1,
		Version: 1, MaKhachHang: "0000000000000000001", TenDangNhap: user, Email: email,
		BaoMat: core.TenantBaoMat{MatKhauHash: passHash, MaPinHash: pinHash},
		RefreshTokens: map[string]core.TenantDeviceToken{
			sessionID: {DeviceID: sessionID, Dev: userAgent, Exp: time.Now().Add(config.ThoiGianHetHanCookie).Unix(), Created: nowUnix},
		},
		VaiTroQuyenHan: "quan_tri_he_thong", ChucVu: "Quản trị hệ thống", TrangThai: 1, // Quyền lực tuyệt đối
		GoiDichVu: make([]core.TenantGoiDichVu, 0), Modules: make(map[string]bool),
		CauHinh: core.TenantCauHinh{Theme: "light", Lang: "vi"},
		ThongTin: core.TenantThongTin{NguonKhachHang: "master_core_setup", TenKhachHang: hoTen, DienThoai: dienThoai, NgaySinh: ngaySinh, GioiTinh: gioiTinh},
		NganHang: core.TenantNganHang{}, MangXaHoi: make(map[string]string),
		NgayTao: nowUnix, NguoiCapNhat: user, NgayCapNhat: nowUnix, 
	}

	core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], newKH)
	core.CacheMapKhachHang[core.TaoCompositeKey(shopID, newKH.MaKhachHang)] = newKH
	bUser, _ := json.Marshal(newKH)
	core.PushAppend(shopID, core.TenSheetKhachHangMaster, []interface{}{newKH.MaKhachHang, string(bUser)})

	// 4. Gửi tin nhắn chào mừng (Chạy ngầm không cần khóa RAM)
	go func() {
		msgID := fmt.Sprintf("MSG_%d_000", time.Now().UnixNano())
		core.ThemMoiTinNhan(shopID, &core.TinNhan{
			MaTinNhan: msgID, LoaiTinNhan: "AUTO",
			NguoiGuiID: "0000000000000000000", NguoiNhanID: []string{"0000000000000000001"}, TieuDe: "Khởi tạo thành công",
			NoiDung: "Hệ thống đã được thiết lập thành công! Chào mừng Quản trị viên tối cao.", 
			NgayTao: nowUnix, NguoiDoc: []string{}, TrangThaiXoa: []string{}, DinhKem: []core.FileDinhKem{}, ThamChieuID: []string{},
		})
	}()

	return sessionID, signature, nil
}
