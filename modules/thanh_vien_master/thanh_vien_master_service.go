package thanh_vien_master

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/config"
	"app/core"
)

type DTO_UpdateThanhVien struct {
	ShopID, AdminID, AdminRole, PinXacNhan string
	MaKH, VaiTro, ChucVu, TrangThai        string
	TenKhachHang, DienThoai, NgaySinh      string
	DiaChi, MaSoThue, GhiChu, AnhDaiDien   string
	NguonKhachHang, Zalo, Facebook, Tiktok string
	GioiTinh                               int
	MatKhauMoi, PinMoi                     string
}

func Service_LuuThanhVien(dto DTO_UpdateThanhVien) error {
	myLevel := Repo_LayCapBac(dto.ShopID, dto.AdminID, dto.AdminRole)
	if myLevel > 2 { return errors.New("Chỉ Quản trị Lõi mới được sửa hồ sơ tại đây!") }

	admin, okAdmin := Repo_LayKhachHang(dto.ShopID, dto.AdminID)
	if !okAdmin || admin.BaoMat.MaPinHash == "" { return errors.New("Vui lòng thiết lập Mã PIN trước.") }
	if !config.KiemTraMatKhau(dto.PinXacNhan, admin.BaoMat.MaPinHash) { return errors.New("Mã PIN xác nhận không chính xác!") }

	kh, ok := Repo_LayKhachHang(dto.ShopID, dto.MaKH)
	if !ok { return errors.New("Tài khoản không tồn tại!") }

	targetLevel := Repo_LayCapBac(dto.ShopID, kh.MaKhachHang, kh.VaiTroQuyenHan)

	// BẢO VỆ LÕI
	if dto.MaKH == "0000000000000000001" && dto.AdminID != "0000000000000000001" { 
		return errors.New("Không ai được chạm vào hồ sơ Sáng Lập Viên!") 
	}
	// [BỔ SUNG CHỐT CHẶN BẢO VỆ BOT]
	if dto.MaKH == "0000000000000000000" && dto.AdminID != "0000000000000000001" { 
		return errors.New("Chỉ Sáng Lập Viên (ID: 001) mới được quyền cấu hình BOT Hệ thống!") 
	}
	if dto.VaiTro == "quan_tri_he_thong" && dto.MaKH != "0000000000000000001" { 
		return errors.New("Chỉ có duy nhất 1 vị trí Quản trị hệ thống!") 
	}
	if dto.MaKH != "0000000000000000000" && dto.AdminID != dto.MaKH && myLevel >= targetLevel { 
		return errors.New("Bạn không có quyền chỉnh sửa cấp trên!") 
	}
	if dto.MaKH == dto.AdminID && dto.TrangThai == "0" { 
		return errors.New("Không thể tự khóa tài khoản chính mình!") 
	}
	core.KhoaHeThong.Lock()
	if dto.VaiTro != "" {
		kh.VaiTroQuyenHan = dto.VaiTro
		if dto.ChucVu != "" { kh.ChucVu = dto.ChucVu 
		} else {
			kh.ChucVu = dto.VaiTro
			for _, v := range core.CacheDanhSachVaiTro[dto.ShopID] {
				if v.MaVaiTro == dto.VaiTro { kh.ChucVu = v.TenVaiTro; break }
			}
		}
	}
	
	// Đẩy vào các Struct con
	kh.ThongTin.TenKhachHang = dto.TenKhachHang
	kh.ThongTin.DienThoai = dto.DienThoai
	kh.ThongTin.NgaySinh = dto.NgaySinh
	kh.ThongTin.DiaChi = dto.DiaChi
	kh.ThongTin.MaSoThue = dto.MaSoThue
	kh.ThongTin.AnhDaiDien = dto.AnhDaiDien
	kh.ThongTin.NguonKhachHang = dto.NguonKhachHang
	kh.ThongTin.GioiTinh = dto.GioiTinh
	kh.GhiChu = dto.GhiChu

	if dto.TrangThai == "1" { kh.TrangThai = 1 } else if dto.TrangThai == "-1" { kh.TrangThai = -1 } else { kh.TrangThai = 0 }

	if kh.MangXaHoi == nil { kh.MangXaHoi = make(map[string]string) }
	kh.MangXaHoi["zalo"] = dto.Zalo
	kh.MangXaHoi["facebook"] = dto.Facebook
	kh.MangXaHoi["tiktok"] = dto.Tiktok

	if dto.MaKH != "0000000000000000000" {
		if dto.MatKhauMoi != "" {
			hash, _ := config.HashMatKhau(dto.MatKhauMoi); kh.BaoMat.MatKhauHash = hash
		}
		if dto.PinMoi != "" {
			hashPin, _ := config.HashMatKhau(dto.PinMoi); kh.BaoMat.MaPinHash = hashPin
		}
	}
	
	kh.NgayCapNhat = time.Now().Unix() 
	kh.NguoiCapNhat = admin.TenDangNhap

	// Đóng gói JSON và nã đạn
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	core.KhoaHeThong.Unlock()

	Repo_GhiCapNhatJSONXuongQueue(dto.ShopID, kh.DongTrongSheet, jsonStr)

	return nil
}

func Service_GuiTinNhan(shopID, adminID, adminRole, tieuDe, noiDung, jsonIDs, sendAsBot string) (int, error) {
	if Repo_LayCapBac(shopID, adminID, adminRole) > 2 { return 0, errors.New("Chỉ Quản trị Lõi mới được phát sóng!") }
	var listMaKH []string
	if err := json.Unmarshal([]byte(jsonIDs), &listMaKH); err != nil || len(listMaKH) == 0 { return 0, errors.New("Chưa chọn người nhận hợp lệ!") }

	senderID := adminID
	if sendAsBot == "1" {
		if bot, ok := Repo_LayKhachHang(shopID, "0000000000000000000"); ok { senderID = bot.MaKhachHang } else { senderID = "0000000000000000000" }
	}

	now := time.Now()
	msgID := fmt.Sprintf("MSG_%d_%s", now.UnixNano(), senderID)

	// Tạo struct Tin Nhắn chuẩn NoSQL
	tn := &core.TinNhan{ 
		MaTinNhan: msgID, 
		LoaiTinNhan: "GROUP", // Chuyển thành GROUP vì danh sách nhận là mảng
		NguoiGuiID: senderID, 
		NguoiNhanID: listMaKH, 
		TieuDe: tieuDe, 
		NoiDung: noiDung, 
		NgayTao: now.Unix(),
		NguoiDoc: []string{}, // Khởi tạo mảng rỗng để hứng ID người seen
		TrangThaiXoa: []string{},
		DinhKem: []core.FileDinhKem{},
		ThamChieuID: []string{},
	}
	
	if len(listMaKH) == 1 && listMaKH[0] == "ALL" {
		tn.LoaiTinNhan = "ALL"
	}

	Repo_ThemTinNhanMoi(shopID, tn)
	return len(listMaKH), nil
}
