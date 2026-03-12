package thanh_vien_master

import "app/core"

func Repo_LayKhachHang(shopID, userID string) (*core.KhachHang, bool) {
	return core.LayKhachHang(shopID, userID)
}

func Repo_LayDanhSach(shopID string) []*core.KhachHang {
	return core.LayDanhSachKhachHang(shopID)
}

func Repo_LayCapBac(shopID, userID, role string) int {
	return core.LayCapBacVaiTro(shopID, userID, role)
}

// Bắn thẳng 1 lệnh JSON xuống Queue (Tạm biệt 26 lệnh update cũ)
func Repo_GhiCapNhatJSONXuongQueue(shopID string, dong int, jsonStr string) {
	core.ThemVaoHangCho(shopID, core.TenSheetKhachHang, dong, core.CotKH_DataJSON, jsonStr)
}

func Repo_ThemTinNhanMoi(shopID string, tn *core.TinNhan) {
	core.ThemMoiTinNhan(shopID, tn)
}
