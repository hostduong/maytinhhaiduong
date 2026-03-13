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

	admin, okAdmin := Repo_LayKhachHangMaster(dto.ShopID, dto.AdminID)
	if !okAdmin || admin.BaoMat.MaPinHash == "" { return errors.New("Vui lòng thiết lập Mã PIN trước.") }
	if !config.KiemTraMatKhau(dto.PinXacNhan, admin.BaoMat.MaPinHash) { return errors.New("Mã PIN xác nhận không chính xác!") }

	kh, ok := Repo_LayKhachHangMaster(dto.ShopID, dto.MaKH)
	if !ok { return errors.New("Tài khoản không tồn tại!") }

	targetLevel := Repo_LayCapBac(dto.ShopID, kh.MaKhachHang, kh.VaiTroQuyenHan)

	// [LUẬT THÉP 1]: BẢO VỆ TUYỆT ĐỐI SÁNG LẬP VIÊN
	if kh.MaKhachHang == "0000000000000000001" || kh.VaiTroQuyenHan == "quan_tri_he_thong" {
		if dto.AdminID != "0000000000000000001" {
			return errors.New("Bất khả xâm phạm: Không ai được phép chạm vào thông tin của Sáng Lập Viên!")
		}
	}

	// [LUẬT THÉP 2]: ĐỘC QUYỀN VƯƠNG MIỆN
	if dto.VaiTro == "quan_tri_he_thong" && dto.MaKH != "0000000000000000001" { 
		return errors.New("Quyền Sáng Lập Viên là duy nhất, không thể cấp cho người khác!") 
	}

	// [LUẬT THÉP 3]: MASTER SỬA CẢ THIÊN HẠ, KẺ KHÁC PHẢI TUÂN THEO CẤP BẬC
	// (Nếu không phải là 001 đang thao tác, thì tiến hành chặn theo Level)
	if dto.AdminID != "0000000000000000001" {
		if dto.MaKH == "0000000000000000000" { 
			return errors.New("Chỉ Sáng Lập Viên mới được phép cấu hình BOT!") 
		}
		if dto.AdminID != dto.MaKH && myLevel >= targetLevel { 
			return errors.New("Thao tác bị từ chối: Bạn không có quyền chỉnh sửa tài khoản cấp trên hoặc ngang quyền!") 
		}
	}

	if dto.MaKH == dto.AdminID && dto.TrangThai == "0" { return errors.New("Không thể tự khóa tài khoản chính mình!") }
	
	lockMaster := core.GetSheetLock(dto.ShopID, core.TenSheetKhachHangMaster)
	lockMaster.Lock()
	
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
		if dto.MatKhauMoi != "" { hash, _ := config.HashMatKhau(dto.MatKhauMoi); kh.BaoMat.MatKhauHash = hash }
		if dto.PinMoi != "" { hashPin, _ := config.HashMatKhau(dto.PinMoi); kh.BaoMat.MaPinHash = hashPin }
	}
	
	kh.NgayCapNhat = time.Now().Unix() 
	kh.NguoiCapNhat = admin.TenDangNhap

	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lockMaster.Unlock() 

	Repo_GhiCapNhatJSONXuongQueue(dto.ShopID, kh.DongTrongSheet, jsonStr)

	return nil
}

func Service_GuiTinNhan(shopID, adminID, adminRole, tieuDe, noiDung, jsonIDs, sendAsBot string) (int, error) {
	if Repo_LayCapBac(shopID, adminID, adminRole) > 2 { return 0, errors.New("Chỉ Quản trị Lõi mới được phát sóng!") }
	var listMaKH []string
	if err := json.Unmarshal([]byte(jsonIDs), &listMaKH); err != nil || len(listMaKH) == 0 { return 0, errors.New("Chưa chọn người nhận hợp lệ!") }

	senderID := adminID
	if sendAsBot == "1" {
		if bot, ok := Repo_LayKhachHangMaster(shopID, "0000000000000000000"); ok { senderID = bot.MaKhachHang } else { senderID = "0000000000000000000" }
	}

	now := time.Now()
	msgID := fmt.Sprintf("MSG_%d_%s", now.UnixNano(), senderID)

	tn := &core.TinNhan{ 
		MaTinNhan: msgID, LoaiTinNhan: "GROUP", NguoiGuiID: senderID, 
		NguoiNhanID: listMaKH, TieuDe: tieuDe, NoiDung: noiDung, 
		NgayTao: now.Unix(), NguoiDoc: []string{}, TrangThaiXoa: []string{},
		DinhKem: []core.FileDinhKem{}, ThamChieuID: []string{},
	}
	
	if len(listMaKH) == 1 && listMaKH[0] == "ALL" { tn.LoaiTinNhan = "ALL" }

	Repo_ThemTinNhanMoi(shopID, tn)
	return len(listMaKH), nil
}
