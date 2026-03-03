package goi_dich_vu

import (
	"app/core"
	"encoding/json"
	"errors"
)

type Service struct { repo Repo }

type DTO_LuuGoiDichVu struct {
	IsNew         bool
	MaGoi         string
	TenGoi        string
	LoaiGoi       string
	ThoiHanNgay   int
	GiaNiemYet    float64
	GiaBan        float64
	CodesJson     string // Hứng mảng Code từ UI gửi lên
	GioiHanJson   string
	MoTa          string
	NhanHienThi   string
	NgayBatDau    string
	NgayKetThuc   string
	SoLuongConLai int
	TrangThai     int
}

func (s *Service) XuLyLuu(shopID string, input DTO_LuuGoiDichVu) error {
	if input.MaGoi == "" || input.TenGoi == "" { return errors.New("Mã và Tên gói là bắt buộc") }

	if input.IsNew {
		if _, exist := s.repo.FindByCode(shopID, input.MaGoi); exist { return errors.New("Mã gói đã tồn tại") }
		
		newG := &core.GoiDichVu{
			MaGoi: input.MaGoi, TenGoi: input.TenGoi, LoaiGoi: input.LoaiGoi,
			ThoiHanNgay: input.ThoiHanNgay, GiaNiemYet: input.GiaNiemYet, GiaBan: input.GiaBan,
			MaCodeKichHoatJson: input.CodesJson, GioiHanJson: input.GioiHanJson,
			MoTa: input.MoTa, NhanHienThi: input.NhanHienThi, NgayBatDau: input.NgayBatDau,
			NgayKetThuc: input.NgayKetThuc, SoLuongConLai: input.SoLuongConLai, TrangThai: input.TrangThai,
		}
		// Parse ngược lại vào RAM mảng Codes
		_ = json.Unmarshal([]byte(input.CodesJson), &newG.DanhSachCode)
		s.repo.Insert(shopID, newG)
	} else {
		g, ok := s.repo.FindByCode(shopID, input.MaGoi)
		if !ok { return errors.New("Không tìm thấy gói để cập nhật") }

		lock := core.GetSheetLock(shopID, core.TenSheetGoiDichVu)
		lock.Lock()
		g.TenGoi = input.TenGoi; g.LoaiGoi = input.LoaiGoi; g.ThoiHanNgay = input.ThoiHanNgay
		g.GiaNiemYet = input.GiaNiemYet; g.GiaBan = input.GiaBan; g.MaCodeKichHoatJson = input.CodesJson
		g.GioiHanJson = input.GioiHanJson; g.MoTa = input.MoTa; g.NhanHienThi = input.NhanHienThi
		g.SoLuongConLai = input.SoLuongConLai; g.TrangThai = input.TrangThai
		_ = json.Unmarshal([]byte(input.CodesJson), &g.DanhSachCode)
		lock.Unlock()

		s.repo.Update(shopID, g)
	}
	return nil
}
