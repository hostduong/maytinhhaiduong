package cau_hinh

import (
	"errors"
	"time"
	"app/core"
)

type Service struct {
	repo Repo
}

type DTO_LuuNhaCungCap struct {
	IsNew              bool
	MaNhaCungCap       string
	TenNhaCungCap      string
	MaSoThue           string
	DienThoai          string
	Email              string
	KhuVuc             string
	DiaChi             string
	NguoiLienHe        string
	NganHang           string
	NhomNhaCungCap     string
	LoaiNhaCungCap     string
	DieuKhoanThanhToan string
	ChietKhauMacDinh   float64
	HanMucCongNo       float64
	CongNoDauKy        float64
	ThongTinThemJson   string
	TrangThai          int
	GhiChu             string
	NguoiThaoTac       string
}

func (s *Service) XuLyLuuNhaCungCap(shopID string, input DTO_LuuNhaCungCap) error {
	if input.TenNhaCungCap == "" { return errors.New("Tên nhà cung cấp không được để trống") }
	nowStr := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")

	if input.IsNew {
		maNCC := input.MaNhaCungCap
		if maNCC == "" { maNCC = s.repo.TaoMaNCCMoi(shopID) } else {
			if _, exist := s.repo.FindNCCByCode(shopID, maNCC); exist { return errors.New("Mã đã tồn tại") }
		}

		newNCC := &core.NhaCungCap{
			SpreadsheetID: shopID, MaNhaCungCap: maNCC, TenNhaCungCap: input.TenNhaCungCap,
			MaSoThue: input.MaSoThue, DienThoai: input.DienThoai, Email: input.Email,
			KhuVuc: input.KhuVuc, DiaChi: input.DiaChi, NguoiLienHe: input.NguoiLienHe,
			NganHang: input.NganHang, NhomNhaCungCap: input.NhomNhaCungCap, LoaiNhaCungCap: input.LoaiNhaCungCap,
			DieuKhoanThanhToan: input.DieuKhoanThanhToan, ChietKhauMacDinh: input.ChietKhauMacDinh,
			HanMucCongNo: input.HanMucCongNo, CongNoDauKy: input.CongNoDauKy,
			TongMua: 0, NoCanTra: input.CongNoDauKy, ThongTinThemJson: input.ThongTinThemJson,
			TrangThai: input.TrangThai, GhiChu: input.GhiChu, NguoiTao: input.NguoiThaoTac,
			NgayTao: nowStr, NgayCapNhat: nowStr,
		}
		s.repo.InsertNCC(shopID, newNCC)
	} else {
		found, exist := s.repo.FindNCCByCode(shopID, input.MaNhaCungCap)
		if !exist { return errors.New("Không tìm thấy dữ liệu") }

		chenhLech := input.CongNoDauKy - found.CongNoDauKy
		noCanTraMoi := found.NoCanTra + chenhLech

		lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
		lock.Lock()
		found.TenNhaCungCap = input.TenNhaCungCap; found.DienThoai = input.DienThoai
		found.CongNoDauKy = input.CongNoDauKy; found.NoCanTra = noCanTraMoi
		found.ThongTinThemJson = input.ThongTinThemJson; found.TrangThai = input.TrangThai
		found.NgayCapNhat = nowStr
		lock.Unlock()

		s.repo.UpdateNCC(shopID, found)
	}
	return nil
}
