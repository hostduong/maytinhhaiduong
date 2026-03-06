package thanh_vien_master

import "app/core"

// Giao tiếp với RAM
func Repo_LayKhachHang(shopID, userID string) (*core.KhachHang, bool) {
	return core.LayKhachHang(shopID, userID)
}

func Repo_LayDanhSach(shopID string) []*core.KhachHang {
	return core.LayDanhSachKhachHang(shopID)
}

func Repo_LayCapBac(shopID, userID, role string) int {
	return core.LayCapBacVaiTro(shopID, userID, role)
}

// Giao tiếp với Background Queue (Đẩy lệnh lưu sheet ngầm)
func Repo_GhiCapNhatXuongQueue(shopID string, dong int, cot int, giatri interface{}) {
	core.ThemVaoHangCho(shopID, core.TenSheetKhachHang, dong, cot, giatri)
}

// Giao tiếp tạo Tin nhắn
func Repo_ThemTinNhanMoi(shopID string, tn *core.TinNhan) {
	core.ThemMoiTinNhan(shopID, tn)
}
