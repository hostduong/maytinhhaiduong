package goi_dich_vu_master

import (
	"app/core"
	"encoding/json"
)

func Repo_FindByCode(shopID, maGoi string) (*core.GoiDichVu, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetCauHinh)
	lock.RLock()
	defer lock.RUnlock()
	g, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(shopID, maGoi)]
	return g, ok
}

func Repo_Insert(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinh)
	lock.Lock()
	
	core.CacheDongHienTaiCauHinh[masterID]++
	gdv.DongTrongSheet = core.CacheDongHienTaiCauHinh[masterID]
	
	core.CacheGoiDichVu[masterID] = append(core.CacheGoiDichVu[masterID], gdv)
	core.CacheMapGoiDichVu[core.TaoCompositeKey(masterID, gdv.MaGoi)] = gdv
	lock.Unlock()

	b, _ := json.Marshal(gdv)
	core.PushAppend(masterID, core.TenSheetCauHinh, []interface{}{
		core.PreGoiDichVu + gdv.MaGoi, 
		string(b),
	})
}

func Repo_Update(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinh)
	lock.Lock()
	lock.Unlock()

	b, _ := json.Marshal(gdv)
	core.PushUpdate(masterID, core.TenSheetCauHinh, gdv.DongTrongSheet, core.CotCH_DataJSON, string(b))
}
