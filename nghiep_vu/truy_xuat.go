package nghiep_vu

import (
	"app/bo_nho_dem" // [MỚI]
	"app/mo_hinh"
)

// [THAY ĐỔI QUAN TRỌNG]
// Thay vì dùng Lock mịn (BoQuanLyKhoa), ta dùng bo_nho_dem.KhoaHeThong.RLock()
// Để đảm bảo khi Reload (Stop-the-world), tất cả các hàm đọc này đều phải dừng lại chờ.

// =================================================================================
// NHÓM 1: MASTER DATA
// =================================================================================

func LayDanhSachSanPham() []mo_hinh.SanPham {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()

	ketQua := make([]mo_hinh.SanPham, len(bo_nho_dem.CacheSanPham.DanhSach))
	copy(ketQua, bo_nho_dem.CacheSanPham.DanhSach)
	return ketQua
}

func LayChiTietSanPham(maSP string) (mo_hinh.SanPham, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	sp, tonTai := bo_nho_dem.CacheSanPham.DuLieu[maSP]
	return sp, tonTai
}

func LayDanhSachDanhMuc() map[string]mo_hinh.DanhMuc {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	
	kq := make(map[string]mo_hinh.DanhMuc)
	for k, v := range bo_nho_dem.CacheDanhMuc.DuLieu { kq[k] = v }
	return kq
}

func LayDanhSachThuongHieu() map[string]mo_hinh.ThuongHieu {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	kq := make(map[string]mo_hinh.ThuongHieu)
	for k, v := range bo_nho_dem.CacheThuongHieu.DuLieu { kq[k] = v }
	return kq
}

func LayDanhSachNhaCungCap() map[string]mo_hinh.NhaCungCap {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	kq := make(map[string]mo_hinh.NhaCungCap)
	for k, v := range bo_nho_dem.CacheNhaCungCap.DuLieu { kq[k] = v }
	return kq
}

func LayThongTinKhachHang(maKH string) (*mo_hinh.KhachHang, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	kh, tonTai := bo_nho_dem.CacheKhachHang.DuLieu[maKH]
	return kh, tonTai
}

func LayCauHinhWeb() map[string]mo_hinh.CauHinhWeb {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	kq := make(map[string]mo_hinh.CauHinhWeb)
	for k, v := range bo_nho_dem.CacheCauHinhWeb.DuLieu { kq[k] = v }
	return kq
}

// =================================================================================
// NHÓM 2: GIAO DỊCH
// =================================================================================

func LayThongTinDonHang(maPX string) (mo_hinh.PhieuXuat, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	px, tonTai := bo_nho_dem.CachePhieuXuat.DuLieu[maPX]
	return px, tonTai
}

func LayChiTietDonHang(maPX string) []mo_hinh.ChiTietPhieuXuat {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	var ketQua []mo_hinh.ChiTietPhieuXuat
	for _, ct := range bo_nho_dem.CacheChiTietXuat.DanhSach {
		if ct.MaPhieuXuat == maPX { ketQua = append(ketQua, ct) }
	}
	return ketQua
}

func LayThongTinVoucher(maVoucher string) (mo_hinh.KhuyenMai, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	km, tonTai := bo_nho_dem.CacheKhuyenMai.DuLieu[maVoucher]
	return km, tonTai
}

func LayThongTinPhieuNhap(maPN string) (mo_hinh.PhieuNhap, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	pn, tonTai := bo_nho_dem.CachePhieuNhap.DuLieu[maPN]
	return pn, tonTai
}

func LayChiTietPhieuNhap(maPN string) []mo_hinh.ChiTietPhieuNhap {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	var ketQua []mo_hinh.ChiTietPhieuNhap
	for _, ct := range bo_nho_dem.CacheChiTietNhap.DanhSach {
		if ct.MaPhieuNhap == maPN { ketQua = append(ketQua, ct) }
	}
	return ketQua
}

func TraCuuSerial(imei string) (mo_hinh.SerialSanPham, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	serial, tonTai := bo_nho_dem.CacheSerial.DuLieu[imei]
	return serial, tonTai
}

func LayThongTinBaoHanh(maPBH string) (mo_hinh.PhieuBaoHanh, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	pbh, tonTai := bo_nho_dem.CachePhieuBaoHanh.DuLieu[maPBH]
	return pbh, tonTai
}

func LayThongTinHoaDon(maHD string) (mo_hinh.HoaDon, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	hd, tonTai := bo_nho_dem.CacheHoaDon.DuLieu[maHD]
	return hd, tonTai
}

func LayChiTietHoaDon(maHD string) []mo_hinh.HoaDonChiTiet {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	var ketQua []mo_hinh.HoaDonChiTiet
	for _, ct := range bo_nho_dem.CacheHoaDonChiTiet.DanhSach {
		if ct.MaHoaDon == maHD { ketQua = append(ketQua, ct) }
	}
	return ketQua
}

func LayPhieuThuChi(maPTC string) (mo_hinh.PhieuThuChi, bool) {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	ptc, tonTai := bo_nho_dem.CachePhieuThuChi.DuLieu[maPTC]
	return ptc, tonTai
}

func LayDanhSachThuChi() []mo_hinh.PhieuThuChi {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()
	ketQua := make([]mo_hinh.PhieuThuChi, len(bo_nho_dem.CachePhieuThuChi.DanhSach))
	copy(ketQua, bo_nho_dem.CachePhieuThuChi.DanhSach)
	return ketQua
}
