package goi_dich_vu

import (
	"app/core"
)

type Repo struct{}

func (r *Repo) FindByCode(shopID, maGoi string) (*core.GoiDichVu, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetGoiDichVu)
	lock.RLock(); defer lock.RUnlock()
	g, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(shopID, maGoi)]
	return g, ok
}

func (r *Repo) Insert(shopID string, g *core.GoiDichVu) {
	lock := core.GetSheetLock(shopID, core.TenSheetGoiDichVu)
	lock.Lock()
	core.CacheGoiDichVu[shopID] = append(core.CacheGoiDichVu[shopID], g)
	core.CacheMapGoiDichVu[core.TaoCompositeKey(shopID, g.MaGoi)] = g
	g.DongTrongSheet = core.DongBatDau_GoiDichVu + len(core.CacheGoiDichVu[shopID]) - 1
	lock.Unlock()

	rowData := []interface{}{
		g.MaGoi, g.TenGoi, g.LoaiGoi, g.ThoiHanNgay, g.GiaNiemYet, g.GiaBan,
		g.MaCodeKichHoatJson, g.GioiHanJson, g.MoTa, g.NhanHienThi,
		g.NgayBatDau, g.NgayKetThuc, g.SoLuongConLai, g.TrangThai,
	}
	core.PushAppend(shopID, core.TenSheetGoiDichVu, rowData)
}

func (r *Repo) Update(shopID string, g *core.GoiDichVu) {
	lock := core.GetSheetLock(shopID, core.TenSheetGoiDichVu)
	lock.RLock()
	row := g.DongTrongSheet
	lock.RUnlock()

	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_TenGoi, g.TenGoi)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_LoaiGoi, g.LoaiGoi)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_ThoiHanNgay, g.ThoiHanNgay)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_GiaNiemYet, g.GiaNiemYet)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_GiaBan, g.GiaBan)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_MaCodeKichHoatJson, g.MaCodeKichHoatJson)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_GioiHanJson, g.GioiHanJson)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_MoTa, g.MoTa)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_NhanHienThi, g.NhanHienThi)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_SoLuongConLai, g.SoLuongConLai)
	core.PushUpdate(shopID, core.TenSheetGoiDichVu, row, core.CotGDV_TrangThai, g.TrangThai)
}
