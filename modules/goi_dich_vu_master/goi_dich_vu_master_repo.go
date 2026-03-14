package goi_dich_vu_master

import (
	"app/core"
	"encoding/json"
)

func Repo_FindByCode(masterID, maGoi string) (*core.GoiDichVu, bool) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
	lock.RLock()
	defer lock.RUnlock()
	g, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(masterID, maGoi)]
	return g, ok
}

func Repo_Insert(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
	lock.Lock()
	
	// Tự động đẩy dòng xuống cuối
	core.CacheDongHienTaiCauHinh[masterID]++
	gdv.DongTrongSheet = core.CacheDongHienTaiCauHinh[masterID]
	
	core.CacheGoiDichVu[masterID] = append(core.CacheGoiDichVu[masterID], gdv)
	core.CacheMapGoiDichVu[core.TaoCompositeKey(masterID, gdv.MaGoi)] = gdv
	lock.Unlock()

	// Đóng gói JSON và nã súng 2 cột
	b, _ := json.Marshal(gdv)
	core.PushAppend(masterID, core.TenSheetCauHinhMaster, []interface{}{
		core.PreGoiDichVu + gdv.MaGoi, 
		string(b),
	})
}

func Repo_Update(masterID string, gdv *core.GoiDichVu) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
	lock.Lock()
	lock.Unlock()

	// Ghi đè cục JSON mới vào Cột B
	b, _ := json.Marshal(gdv)
	core.PushUpdate(masterID, core.TenSheetCauHinhMaster, gdv.DongTrongSheet, core.CotCH_DataJSON, string(b))
}
