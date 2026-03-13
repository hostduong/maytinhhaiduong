package thanh_vien_master

import "app/core"

// [FIX] Tự viết hàm Lấy Khách Hàng chuyên biệt cho Master để Lock đúng Sheet
func Repo_LayKhachHangMaster(masterID, userID string) (*core.KhachHang, bool) {
	lock := core.GetSheetLock(masterID, core.TenSheetKhachHangMaster)
	lock.RLock()
	defer lock.RUnlock()
	kh, ok := core.CacheMapKhachHang[core.TaoCompositeKey(masterID, userID)]
	return kh, ok
}

func Repo_LayDanhSachMaster(masterID string) []*core.KhachHang {
	lock := core.GetSheetLock(masterID, core.TenSheetKhachHangMaster)
	lock.RLock()
	defer lock.RUnlock()
	return core.CacheKhachHang[masterID]
}

func Repo_LayCapBac(shopID, userID, role string) int {
	return core.LayCapBacVaiTro(shopID, userID, role)
}

func Repo_GhiCapNhatJSONXuongQueue(shopID string, dong int, jsonStr string) {
	// [FIX LỚN] Bắn đúng vào KHACH_HANG_MASTER thay vì KHACH_HANG
	core.ThemVaoHangCho(shopID, core.TenSheetKhachHangMaster, dong, core.CotKH_DataJSON, jsonStr)
}

func Repo_ThemTinNhanMoi(shopID string, tn *core.TinNhan) {
	core.ThemMoiTinNhan(shopID, tn)
}
