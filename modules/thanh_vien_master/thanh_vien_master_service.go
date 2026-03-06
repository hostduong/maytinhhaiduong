package thanh_vien_master

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"app/config"
	"app/core"
)

// DTO (Data Transfer Object) để nhận dữ liệu từ API truyền sang
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
	if myLevel > 2 {
		return errors.New("Chỉ Quản trị Lõi 99K mới được sửa hồ sơ tại đây!")
	}

	admin, okAdmin := Repo_LayKhachHang(dto.ShopID, dto.AdminID)
	if !okAdmin || admin.MaPinHash == "" { return errors.New("Vui lòng thiết lập Mã PIN trước.") }
	if !config.KiemTraMatKhau(dto.PinXacNhan, admin.MaPinHash) { return errors.New("Mã PIN xác nhận không chính xác!") }

	kh, ok := Repo_LayKhachHang(dto.ShopID, dto.MaKH)
	if !ok { return errors.New("Tài khoản không tồn tại!") }

	targetLevel := Repo_LayCapBac(dto.ShopID, kh.MaKhachHang, kh.VaiTroQuyenHan)

	// BẢO VỆ LÕI (SECURITY CHECKS)
	if dto.MaKH == "0000000000000000001" && dto.AdminID != "0000000000000000001" {
		return errors.New("BẢO MẬT: Không ai được chạm vào hồ sơ Người Sáng Lập!")
	}
	if dto.VaiTro == "quan_tri_he_thong" && dto.MaKH != "0000000000000000001" {
		return errors.New("BẢO MẬT: Chỉ có duy nhất 1 vị trí Quản trị hệ thống (ID 001)!")
	}
	if dto.MaKH != "0000000000000000000" && dto.AdminID != dto.MaKH && myLevel >= targetLevel {
		return errors.New("Lỗi: Bạn không có quyền chỉnh sửa cấp trên hoặc người ngang hàng!")
	}
	if dto.VaiTro != "" && dto.VaiTro != kh.VaiTroQuyenHan {
		newLevel := Repo_LayCapBac(dto.ShopID, "", dto.VaiTro)
		if newLevel <= myLevel && myLevel != 0 {
			return errors.New("Lỗi: Không thể bổ nhiệm chức vụ ngang bằng hoặc cao hơn quyền hạn của bạn!")
		}
	}
	if dto.MaKH == dto.AdminID && dto.TrangThai == "0" {
		return errors.New("Hệ thống bảo vệ: Không thể tự khóa tài khoản chính mình!")
	}

	// XỬ LÝ LƯU RAM
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
	
	kh.TenKhachHang = dto.TenKhachHang
	kh.DienThoai = dto.DienThoai
	kh.NgaySinh = dto.NgaySinh
	kh.DiaChi = dto.DiaChi
	kh.MaSoThue = dto.MaSoThue
	kh.GhiChu = dto.GhiChu
	kh.AnhDaiDien = dto.AnhDaiDien
	kh.NguonKhachHang = dto.NguonKhachHang
	kh.GioiTinh = dto.GioiTinh

	if dto.TrangThai == "1" { kh.TrangThai = 1 } else if dto.TrangThai == "-1" { kh.TrangThai = -1 } else { kh.TrangThai = 0 }

	kh.MangXaHoi.Zalo = dto.Zalo
	kh.MangXaHoi.Facebook = dto.Facebook
	kh.MangXaHoi.Tiktok = dto.Tiktok

	if dto.MaKH != "0000000000000000000" {
		if dto.MatKhauMoi != "" {
			hash, _ := config.HashMatKhau(dto.MatKhauMoi); kh.MatKhauHash = hash
			Repo_GhiCapNhatXuongQueue(dto.ShopID, kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
		}
		if dto.PinMoi != "" {
			hashPin, _ := config.HashMatKhau(dto.PinMoi); kh.MaPinHash = hashPin
			Repo_GhiCapNhatXuongQueue(dto.ShopID, kh.DongTrongSheet, core.CotKH_MaPinHash, hashPin)
		}
	}
	kh.NgayCapNhat = time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")
	kh.NguoiCapNhat = admin.TenDangNhap
	core.KhoaHeThong.Unlock()

	// ĐẨY RA HÀNG ĐỢI BACKGROUND ĐỂ GHI SHEET NGẦM
	r := kh.DongTrongSheet
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_TenKhachHang, kh.TenKhachHang)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_DienThoai, kh.DienThoai)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_NgaySinh, kh.NgaySinh)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_GioiTinh, kh.GioiTinh)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_DiaChi, kh.DiaChi)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_MaSoThue, kh.MaSoThue)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_GhiChu, kh.GhiChu)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_TrangThai, kh.TrangThai)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_ChucVu, kh.ChucVu)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_AnhDaiDien, kh.AnhDaiDien)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_NguonKhachHang, kh.NguonKhachHang)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
	Repo_GhiCapNhatXuongQueue(dto.ShopID, r, core.CotKH_MangXaHoiJson, core.ToJSON(kh.MangXaHoi))

	return nil
}

func Service_GuiTinNhan(shopID, adminID, adminRole, tieuDe, noiDung, jsonIDs, sendAsBot string) (int, error) {
	if Repo_LayCapBac(shopID, adminID, adminRole) > 2 {
		return 0, errors.New("Chỉ Quản trị Lõi 99K mới được phát sóng thông báo!")
	}

	var listMaKH []string
	if err := json.Unmarshal([]byte(jsonIDs), &listMaKH); err != nil || len(listMaKH) == 0 {
		return 0, errors.New("Chưa chọn người nhận hợp lệ!")
	}

	senderID := adminID
	if sendAsBot == "1" {
		if bot, ok := Repo_LayKhachHang(shopID, "0000000000000000000"); ok {
			senderID = bot.MaKhachHang
		} else { senderID = "0000000000000000000" }
	}

	now := time.Now()
	msgID := fmt.Sprintf("ALL_%d_%s", now.UnixNano(), senderID)
	nowStr := now.In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")

	Repo_ThemTinNhanMoi(shopID, &core.TinNhan{
		MaTinNhan: msgID, LoaiTinNhan: "ALL", NguoiGuiID: senderID, NguoiNhanID: jsonIDs,
		TieuDe: tieuDe, NoiDung: noiDung, NgayTao: nowStr,
	})

	return len(listMaKH), nil
}
