package goi_dich_vu_master

import "app/core"

// ====================================================================
// REPO: CHUYÊN GIAO TIẾP VỚI RAM VÀ ĐẨY HÀNG CHỜ (QUEUE)
// ====================================================================

func Repo_FindByCode(shopID, maGoi string) (*core.GoiDichVu, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetGoiDichVuMaster)
	lock.RLock()
	defer lock.RUnlock()
	g, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(shopID, maGoi)]
	return g, ok
}

func Repo_Insert(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetGoiDichVuMaster)
	lock.Lock()
	gdv.DongTrongSheet = len(core.CacheGoiDichVu[masterID]) + core.DongBatDau_GoiDichVu
	core.CacheGoiDichVu[masterID] = append(core.CacheGoiDichVu[masterID], gdv)
	core.CacheMapGoiDichVu[core.TaoCompositeKey(masterID, gdv.MaGoi)] = gdv
	lock.Unlock()

	core.PushAppend(masterID, core.TenSheetGoiDichVuMaster, []interface{}{
		gdv.MaGoi, gdv.TenGoi, gdv.LoaiGoi, gdv.ThoiHanNgay, gdv.ThoiHanHienThi, gdv.NhanHienThi, 
		gdv.GiaNiemYet, gdv.GiaBan, gdv.MaCodeKichHoatJson, gdv.GioiHanJson, 
		gdv.MoTa, gdv.NgayBatDau, gdv.NgayKetThuc, gdv.SoLuongConLai, gdv.TrangThai,
	})
}

func Repo_Update(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetGoiDichVuMaster)
	lock.Lock()
	lock.Unlock()

	r := gdv.DongTrongSheet
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
	core.PushUpdate(masterID, core.TenSheetGoiDichVuMaster, r, core.CotGDV_SoLuongConLai, gdv.SoLuongConLai)
}
