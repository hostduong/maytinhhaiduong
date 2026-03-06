package cua_hang_master

import (
	"errors"
	"time"

	"app/config"
	"app/core"
)

type DTO_UpdateCuaHang struct {
	AdminID, AdminRole, PinXacNhan         string
	MaKH, TrangThai                        string
	TenKhachHang, DienThoai, NgaySinh      string
	DiaChi, MaSoThue, GhiChu, AnhDaiDien   string
	Zalo, Facebook, Tiktok                 string
	GioiTinh                               int
	MatKhauMoi                             string
	
	// [ĐÃ VÁ LỖI]: Bổ sung 4 trường vào DTO
	VaiTro, ChucVu, NguonKhachHang, PinMoi string
	
	SpreadsheetID                          string
	CustomDomain                           string
}

func Service_LuuCuaHang(dto DTO_UpdateCuaHang) error {
	masterID := config.BienCauHinh.IdFileSheetMaster
	adminID := config.BienCauHinh.IdFileSheetAdmin

	// Kiểm tra quyền của người sửa (Chỉ Sếp mới được sửa)
	myLevel := core.LayCapBacVaiTro(masterID, dto.AdminID, dto.AdminRole)
	if myLevel > 2 {
		return errors.New("Chỉ Quản trị Lõi 99K mới được sửa hồ sơ tại đây!")
	}

	// Xác thực PIN
	admin, okAdmin := core.LayKhachHang(masterID, dto.AdminID)
	if !okAdmin || admin.MaPinHash == "" { return errors.New("Vui lòng thiết lập Mã PIN trước.") }
	if !config.KiemTraMatKhau(dto.PinXacNhan, admin.MaPinHash) { return errors.New("Mã PIN xác nhận không chính xác!") }

	// Lấy thông tin Cửa hàng (Từ Tầng Admin)
	kh, ok := core.LayKhachHang(adminID, dto.MaKH)
	if !ok { return errors.New("Cửa hàng không tồn tại!") }

	// Bắt đầu nhào nặn RAM
	core.KhoaHeThong.Lock()
	
	// [ĐÃ VÁ LỖI]: Xử lý cập nhật Vai trò và Chức vụ tùy chỉnh
	if dto.VaiTro != "" {
		kh.VaiTroQuyenHan = dto.VaiTro
		if dto.ChucVu != "" { 
			kh.ChucVu = dto.ChucVu 
		} else {
			kh.ChucVu = dto.VaiTro
			for _, v := range core.CacheDanhSachVaiTro[adminID] {
				if v.MaVaiTro == dto.VaiTro { kh.ChucVu = v.TenVaiTro; break }
			}
		}
	}

	kh.TenKhachHang = dto.TenKhachHang
	kh.DienThoai = dto.DienThoai
	kh.NgaySinh = dto.NgaySinh
	kh.DiaChi = dto.DiaChi
	kh.MaSoThue = dto.MaSoThue
	kh.GhiChu = dto.GhiChu
	kh.AnhDaiDien = dto.AnhDaiDien
	kh.GioiTinh = dto.GioiTinh
	kh.NguonKhachHang = dto.NguonKhachHang // [ĐÃ VÁ LỖI]

	if dto.TrangThai == "1" { kh.TrangThai = 1 } else if dto.TrangThai == "-1" { kh.TrangThai = -1 } else { kh.TrangThai = 0 }

	kh.MangXaHoi.Zalo = dto.Zalo
	kh.MangXaHoi.Facebook = dto.Facebook
	kh.MangXaHoi.Tiktok = dto.Tiktok
	
	kh.DataSheets.SpreadsheetID = dto.SpreadsheetID
	kh.CauHinh.CustomDomain = dto.CustomDomain

	if dto.MatKhauMoi != "" {
		hash, _ := config.HashMatKhau(dto.MatKhauMoi); kh.MatKhauHash = hash
		core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	}

	// [ĐÃ VÁ LỖI]: Xử lý cập nhật mã PIN mới
	if dto.PinMoi != "" {
		hashPin, _ := config.HashMatKhau(dto.PinMoi); kh.MaPinHash = hashPin
		core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, kh.DongTrongSheet, core.CotKH_MaPinHash, hashPin)
	}

	kh.NgayCapNhat = time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")
	kh.NguoiCapNhat = admin.TenDangNhap
	core.KhoaHeThong.Unlock()

	// Ghi ngầm
	r := kh.DongTrongSheet
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_TenKhachHang, kh.TenKhachHang)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_DienThoai, kh.DienThoai)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_NgaySinh, kh.NgaySinh)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_GioiTinh, kh.GioiTinh)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_DiaChi, kh.DiaChi)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_MaSoThue, kh.MaSoThue)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_GhiChu, kh.GhiChu)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_TrangThai, kh.TrangThai)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_AnhDaiDien, kh.AnhDaiDien)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
	
	// [ĐÃ VÁ LỖI]: Bắn lệnh ghi 3 cột mới xuống Google Sheets
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_ChucVu, kh.ChucVu)
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_NguonKhachHang, kh.NguonKhachHang)

	// JSON Fields
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_MangXaHoiJson, core.ToJSON(kh.MangXaHoi))
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_DataSheetsJson, core.ToJSON(kh.DataSheets))
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_CauHinhJson, core.ToJSON(kh.CauHinh))

	return nil
}
