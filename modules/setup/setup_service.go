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
	// 1. Phải tải thành công Sheet mới cho đăng ký
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { 
		return "", "", errors.New("Không thể kết nối đến Database. Vui lòng thử lại sau!") 
	}

	// [CHỐT CHẶN CHỐNG SPAM] Khóa RAM để đảm bảo chỉ 1 người lọt vào
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHangMaster)
	lock.Lock()
	defer lock.Unlock()

	_, hasGod := core.LayKhachHang(shopID, "0000000000000000001")
	if hasGod {
		return "", "", errors.New("Hệ thống đã được khởi tạo! Không thể đăng ký thêm tài khoản Sáng lập.")
	}

	// ==============================================================
	// [ZERO-TRUST VALIDATION]: KIỂM TRA ĐẦU VÀO NGHIÊM NGẶT TẠI BACKEND
	// ==============================================================
	if dienThoai == "" { return "", "", errors.New("Số điện thoại không được để trống!") }
	if !config.KiemTraHoTen(hoTen) { return "", "", errors.New("Họ tên chứa ký tự không hợp lệ hoặc quá ngắn!") }
	if !config.KiemTraTenDangNhap(user) { return "", "", errors.New("Tên đăng nhập sai định dạng (Chỉ a-z, 0-9, dấu gạch ngang)!") }
	if !config.KiemTraEmail(email) { return "", "", errors.New("Định dạng Email không hợp lệ!") }
	if !config.KiemTraDinhDangMatKhau(pass) { return "", "", errors.New("Mật khẩu không đạt chuẩn an toàn (Cần chữ hoa, số, ký tự đặc biệt)!") }
	if !config.KiemTraMaPin(maPin) { return "", "", errors.New("Mã PIN bắt buộc phải là 8 chữ số!") }

	gioiTinh := -1
	if gioiTinhStr == "Nam" { gioiTinh = 1 } else if gioiTinhStr == "Nữ" { gioiTinh = 0 }

	passHash, _ := config.HashMatKhau(pass)
	pinHash, _ := config.HashMatKhau(maPin)
	nowUnix := time.Now().Unix()

	// 2. KHAI SINH BOT HỆ THỐNG (ID: 000)
	botKH := &core.KhachHang{
		SpreadsheetID: shopID, DongTrongSheet: core.DongBatDau_KhachHang,
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
		VaiTroQuyenHan: "quan_tri_he_thong", ChucVu: "Quản trị hệ thống", TrangThai: 1, 
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

	// 4. Gửi tin nhắn chào mừng
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
