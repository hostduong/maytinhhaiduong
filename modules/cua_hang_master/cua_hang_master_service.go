package cua_hang_master

import (
	"encoding/json"
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
	VaiTro, ChucVu, NguonKhachHang, PinMoi string
	SpreadsheetID                          string
	CustomDomain                           string
}

func Service_LuuCuaHang(dto DTO_UpdateCuaHang) error {
	masterID := config.BienCauHinh.IdFileSheetMaster
	adminID := config.BienCauHinh.IdFileSheetAdmin

	myLevel := core.LayCapBacVaiTro(masterID, dto.AdminID, dto.AdminRole)
	if myLevel > 2 {
		return errors.New("Chỉ Quản trị Lõi 99K mới được sửa hồ sơ tại đây!")
	}

	admin, okAdmin := core.LayKhachHang(masterID, dto.AdminID)
	if !okAdmin || admin.BaoMat.MaPinHash == "" { return errors.New("Vui lòng thiết lập Mã PIN trước.") }
	if !config.KiemTraMatKhau(dto.PinXacNhan, admin.BaoMat.MaPinHash) { return errors.New("Mã PIN xác nhận không chính xác!") }

	kh, ok := core.LayKhachHang(adminID, dto.MaKH)
	if !ok { return errors.New("Cửa hàng không tồn tại!") }

	core.KhoaHeThong.Lock()
	
	if dto.VaiTro != "" {
		kh.VaiTroQuyenHan = dto.VaiTro
		if dto.ChucVu != "" { kh.ChucVu = dto.ChucVu 
		} else {
			kh.ChucVu = dto.VaiTro
			if pq, ok := core.CacheMapPhanQuyen[core.TaoCompositeKey(dto.ShopID, dto.VaiTro)]; ok {
				kh.ChucVu = pq.TenVaiTro
			}
		}
	}

	kh.ThongTin.TenKhachHang = dto.TenKhachHang
	kh.ThongTin.DienThoai = dto.DienThoai
	kh.ThongTin.NgaySinh = dto.NgaySinh
	kh.ThongTin.DiaChi = dto.DiaChi
	kh.ThongTin.MaSoThue = dto.MaSoThue
	kh.ThongTin.AnhDaiDien = dto.AnhDaiDien
	kh.ThongTin.GioiTinh = dto.GioiTinh
	kh.ThongTin.NguonKhachHang = dto.NguonKhachHang
	kh.GhiChu = dto.GhiChu

	if dto.TrangThai == "1" { kh.TrangThai = 1 } else if dto.TrangThai == "-1" { kh.TrangThai = -1 } else { kh.TrangThai = 0 }

	if kh.MangXaHoi == nil { kh.MangXaHoi = make(map[string]string) }
	kh.MangXaHoi["zalo"] = dto.Zalo
	kh.MangXaHoi["facebook"] = dto.Facebook
	kh.MangXaHoi["tiktok"] = dto.Tiktok
	
	kh.System.SheetID = dto.SpreadsheetID
	kh.Domain.CustomDomain = dto.CustomDomain

	if dto.MatKhauMoi != "" {
		hash, _ := config.HashMatKhau(dto.MatKhauMoi); kh.BaoMat.MatKhauHash = hash
	}
	if dto.PinMoi != "" {
		hashPin, _ := config.HashMatKhau(dto.PinMoi); kh.BaoMat.MaPinHash = hashPin
	}

	kh.NgayCapNhat = time.Now().Unix()
	kh.NguoiCapNhat = admin.TenDangNhap
	
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	r := kh.DongTrongSheet
	core.KhoaHeThong.Unlock()

	// 1 Lệnh thay cho 15 lệnh
	core.ThemVaoHangCho(adminID, core.TenSheetKhachHangAdmin, r, core.CotKH_DataJSON, jsonStr)

	return nil
}
