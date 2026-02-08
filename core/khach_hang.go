package core

import (
	"fmt"
	"strings"
	"time"

	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH CỘT
// =============================================================
const (
	// [CHUẨN HÓA] Tên biến theo format: DongBatDau_<TênSheet>
	DongBatDau_KhachHang = 11 

	CotKH_MaKhachHang      = 0
	CotKH_TenDangNhap      = 1
	CotKH_MatKhauHash      = 2
	CotKH_Cookie           = 3
	CotKH_CookieExpired    = 4
	CotKH_MaPinHash        = 5
	CotKH_LoaiKhachHang    = 6
	CotKH_TenKhachHang     = 7
	CotKH_DienThoai        = 8
	CotKH_Email            = 9
	CotKH_UrlFb            = 10
	CotKH_Zalo             = 11
	CotKH_UrlTele          = 12
	CotKH_UrlTiktok        = 13
	CotKH_DiaChi           = 14
	CotKH_NgaySinh         = 15
	CotKH_GioiTinh         = 16
	CotKH_MaSoThue         = 17
	CotKH_DangNo           = 18
	CotKH_TongMua          = 19
	CotKH_ChucVu           = 20
	CotKH_VaiTroQuyenHan   = 21
	CotKH_TrangThai        = 22
	CotKH_GhiChu           = 23
	CotKH_NguoiTao         = 24
	CotKH_NgayTao          = 25
	CotKH_NgayCapNhat      = 26
)

// ... (Phần Struct và Code bên dưới GIỮ NGUYÊN như cũ) ...
// ... (Lưu ý: Trong hàm NapKhachHang, nhớ sửa dòng kiểm tra i < DongBatDau... nhé)

// Tôi viết lại hàm NapKhachHang để bạn copy đè cho chắc:
func NapKhachHang(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" { targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := loadSheetData(targetSpreadsheetID, "KHACH_HANG")
	if err != nil { return }

	_Map_KhachHang = make(map[string]*KhachHang)
	_DS_KhachHang = []*KhachHang{}

	for i, r := range raw {
		// [SỬA] Dùng biến chuẩn DongBatDau_KhachHang
		if i < DongBatDau_KhachHang-1 { continue }
		
		maKH := layString(r, CotKH_MaKhachHang)
		if maKH == "" { continue } // Logic lọc rác nới lỏng (chỉ cần có mã)
		
		key := TaoCompositeKey(targetSpreadsheetID, maKH)
		if _, daTonTai := _Map_KhachHang[key]; daTonTai { continue }

		kh := &KhachHang{
			SpreadsheetID:  targetSpreadsheetID,
			DongTrongSheet: i + 1,
			MaKhachHang:    maKH,
			TenDangNhap:    layString(r, CotKH_TenDangNhap),
			MatKhauHash:    layString(r, CotKH_MatKhauHash),
			Cookie:         layString(r, CotKH_Cookie),
			CookieExpired:  int64(layFloat(r, CotKH_CookieExpired)),
			MaPinHash:      layString(r, CotKH_MaPinHash),
			LoaiKhachHang:  layString(r, CotKH_LoaiKhachHang),
			TenKhachHang:   layString(r, CotKH_TenKhachHang),
			DienThoai:      layString(r, CotKH_DienThoai),
			Email:          layString(r, CotKH_Email),
			UrlFb:          layString(r, CotKH_UrlFb),
			Zalo:           layString(r, CotKH_Zalo),
			UrlTele:        layString(r, CotKH_UrlTele),
			UrlTiktok:      layString(r, CotKH_UrlTiktok),
			DiaChi:         layString(r, CotKH_DiaChi),
			NgaySinh:       layString(r, CotKH_NgaySinh),
			GioiTinh:       layString(r, CotKH_GioiTinh),
			MaSoThue:       layString(r, CotKH_MaSoThue),
			DangNo:         layFloat(r, CotKH_DangNo),
			TongMua:        layFloat(r, CotKH_TongMua),
			ChucVu:         layString(r, CotKH_ChucVu),
			VaiTroQuyenHan: layString(r, CotKH_VaiTroQuyenHan),
			TrangThai:      layInt(r, CotKH_TrangThai),
			GhiChu:         layString(r, CotKH_GhiChu),
			NguoiTao:       layString(r, CotKH_NguoiTao),
			NgayTao:        layString(r, CotKH_NgayTao),
			NgayCapNhat:    layString(r, CotKH_NgayCapNhat),
		}
		_DS_KhachHang = append(_DS_KhachHang, kh)
		_Map_KhachHang[key] = kh
	}
}

func LayDanhSachKhachHang() []*KhachHang {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	kq := make([]*KhachHang, len(_DS_KhachHang))
	copy(kq, _DS_KhachHang)
	return kq
}

func LayKhachHang(maKH string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	sheetID := cau_hinh.BienCauHinh.IdFileSheet
	key := TaoCompositeKey(sheetID, maKH)
	kh, ok := _Map_KhachHang[key]
	return kh, ok
}

func TimKhachHangTheoCookie(cookie string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	for _, kh := range _DS_KhachHang {
		if kh.Cookie == cookie && cookie != "" {
			if time.Now().Unix() > kh.CookieExpired { return nil, false }
			return kh, true
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(input string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range _DS_KhachHang {
		if strings.ToLower(kh.TenDangNhap) == input { return kh, true }
		if kh.Email != "" && strings.ToLower(kh.Email) == input { return kh, true }
	}
	return nil, false
}

func KiemTraTonTaiUserEmail(user, email string) bool {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	user = strings.ToLower(strings.TrimSpace(user))
	email = strings.ToLower(strings.TrimSpace(email))
	for _, kh := range _DS_KhachHang {
		if strings.ToLower(kh.TenDangNhap) == user { return true }
		if email != "" && strings.ToLower(kh.Email) == email { return true }
	}
	return false
}

func TaoMaKhachHangMoi() string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	maxID := 0
	for _, kh := range _DS_KhachHang {
		if kh.SpreadsheetID != cau_hinh.BienCauHinh.IdFileSheet { continue }
		parts := strings.Split(kh.MaKhachHang, "_")
		if len(parts) == 2 {
			var id int
			fmt.Sscanf(parts[1], "%d", &id)
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	if kh.SpreadsheetID == "" { kh.SpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet }
	_DS_KhachHang = append(_DS_KhachHang, kh)
	key := TaoCompositeKey(kh.SpreadsheetID, kh.MaKhachHang)
	_Map_KhachHang[key] = kh
}
