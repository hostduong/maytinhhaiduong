package cau_hinh_he_thong

import (
	"errors"
	"time"
	"app/core"
)

type CauHinhService struct {
	repo CauHinhRepo
}

// Input DTO (Data Transfer Object) để Service hứng dữ liệu "Sạch" từ API
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

// Logic thực thi
func (s *CauHinhService) XuLyLuuNhaCungCap(shopID string, input DTO_LuuNhaCungCap) error {
	if input.TenNhaCungCap == "" {
		return errors.New("Tên nhà cung cấp không được để trống")
	}

	nowStr := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")

	if input.IsNew {
		maNCC := input.MaNhaCungCap
		if maNCC == "" {
			maNCC = s.repo.TaoMaNCCMoi(shopID)
		} else {
			if _, exist := s.repo.FindNCCByCode(shopID, maNCC); exist {
				return errors.New("Mã nhà cung cấp đã tồn tại")
			}
		}

		newNCC := &core.NhaCungCap{
			SpreadsheetID:      shopID,
			MaNhaCungCap:       maNCC,
			TenNhaCungCap:      input.TenNhaCungCap,
			MaSoThue:           input.MaSoThue,
			DienThoai:          input.DienThoai,
			Email:              input.Email,
			KhuVuc:             input.KhuVuc,
			DiaChi:             input.DiaChi,
			NguoiLienHe:        input.NguoiLienHe,
			NganHang:           input.NganHang,
			NhomNhaCungCap:     input.NhomNhaCungCap,
			LoaiNhaCungCap:     input.LoaiNhaCungCap,
			DieuKhoanThanhToan: input.DieuKhoanThanhToan,
			ChietKhauMacDinh:   input.ChietKhauMacDinh,
			HanMucCongNo:       input.HanMucCongNo,
			CongNoDauKy:        input.CongNoDauKy,
			TongMua:            0,                 // Mặc định tạo mới
			NoCanTra:           input.CongNoDauKy, // Nợ cần trả = Nợ đầu kỳ
			ThongTinThemJson:   input.ThongTinThemJson,
			TrangThai:          input.TrangThai,
			GhiChu:             input.GhiChu,
			NguoiTao:           input.NguoiThaoTac,
			NgayTao:            nowStr,
			NgayCapNhat:        nowStr,
		}
		s.repo.InsertNCC(shopID, newNCC)

	} else {
		found, exist := s.repo.FindNCCByCode(shopID, input.MaNhaCungCap)
		if !exist {
			return errors.New("Không tìm thấy Nhà cung cấp để sửa")
		}

		// LOGIC KẾ TOÁN: Bù trừ nợ khi sửa số dư đầu kỳ
		chenhLechNoDauKy := input.CongNoDauKy - found.CongNoDauKy
		noCanTraMoi := found.NoCanTra + chenhLechNoDauKy

		// Khóa RAM để thay đổi an toàn
		lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
		lock.Lock()
		found.TenNhaCungCap = input.TenNhaCungCap
		found.DienThoai = input.DienThoai
		// ... (Gán các trường còn lại)
		found.CongNoDauKy = input.CongNoDauKy
		found.NoCanTra = noCanTraMoi
		found.ThongTinThemJson = input.ThongTinThemJson
		found.NgayCapNhat = nowStr
		lock.Unlock()

		s.repo.UpdateNCC(shopID, found)
	}

	return nil
}
