package goi_dich_vu

import "app/core"

// ====================================================================
// REPO: CHUYÊN GIAO TIẾP VỚI RAM VÀ ĐẨY HÀNG CHỜ (QUEUE)
// ====================================================================

func Repo_LayDanhSachGoiDichVu(masterID string) []*core.GoiDichVu {
	lock := core.GetSheetLock(masterID, core.TenSheetGoiDichVuMaster)
	lock.RLock()
	defer lock.RUnlock()
	return core.CacheGoiDichVu[masterID]
}

func Repo_ThemGoiDichVu(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetGoiDichVuMaster)
	lock.Lock()
	
	// Tính dòng tiếp theo (Dòng hiện tại + số dòng header)
	gdv.DongTrongSheet = len(core.CacheGoiDichVu[masterID]) + core.DongBatDau_GoiDichVu
	core.CacheGoiDichVu[masterID] = append(core.CacheGoiDichVu[masterID], gdv)
	core.CacheMapGoiDichVu[core.TaoCompositeKey(masterID, gdv.MaGoi)] = gdv
	
	lock.Unlock()

	// Ghi xuống Queue để chạy ngầm (Sử dụng đúng biến Master)
	core.PushAppend(masterID, core.TenSheetGoiDichVuMaster, []interface{}{
		gdv.MaGoi, gdv.TenGoi, gdv.LoaiGoi, gdv.ThoiHanNgay, gdv.ThoiHanHienThi, gdv.NhanHienThi, 
		gdv.GiaNiemYet, gdv.GiaBan, gdv.MaCodeKichHoatJson, gdv.GioiHanJson, 
		gdv.MoTa, gdv.NgayBatDau, gdv.NgayKetThuc, gdv.SoLuongConLai, gdv.TrangThai,
	})
}

func Repo_CapNhatGoiDichVu(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetGoiDichVuMaster)
	lock.Lock()
	// RAM đã được update từ Service trước đó
	lock.Unlock()

	r := gdv.DongTrongSheet
	// Ghi từng ô xuống Queue ngầm (Sử dụng đúng biến Master)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_TenGoi, gdv.TenGoi)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_LoaiGoi, gdv.LoaiGoi)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_ThoiHanNgay, gdv.ThoiHanNgay)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_ThoiHanHienThi, gdv.ThoiHanHienThi)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_NhanHienThi, gdv.NhanHienThi)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_GiaNiemYet, gdv.GiaNiemYet)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_GiaBan, gdv.GiaBan)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_MaCodeKichHoatJson, gdv.MaCodeKichHoatJson)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_GioiHanJson, gdv.GioiHanJson)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_MoTa, gdv.MoTa)
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_TrangThai, gdv.TrangThai)
}
