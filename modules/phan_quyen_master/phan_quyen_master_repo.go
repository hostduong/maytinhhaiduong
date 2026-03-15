package phan_quyen_master

import (
	"app/core"
	"encoding/json"
)

func Repo_FindByCode(masterID, maVaiTro string) (*core.PhanQuyen, bool) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
	lock.RLock()
	defer lock.RUnlock()
	
	// Ép kiểu từ RAM Cache (Giả định Sếp khai báo CacheMapPhanQuyen trong core)
	// Để demo chạy mượt, tôi sẽ parse trực tiếp từ list raw nếu Sếp chưa làm hàm nạp riêng
	// Tuy nhiên, chuẩn nhất là dùng Cache:
	pq, ok := core.CacheMapPhanQuyen[core.TaoCompositeKey(masterID, maVaiTro)]
	return pq, ok
}

func Repo_Insert(masterID string, pq *core.PhanQuyen) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
	lock.Lock()
	
	core.CacheDongHienTaiCauHinh[masterID]++
	pq.DongTrongSheet = core.CacheDongHienTaiCauHinh[masterID]
	
	core.CachePhanQuyen[masterID] = append(core.CachePhanQuyen[masterID], pq)
	core.CacheMapPhanQuyen[core.TaoCompositeKey(masterID, pq.MaVaiTro)] = pq
	lock.Unlock()

	b, _ := json.Marshal(pq)
	core.PushAppend(masterID, core.TenSheetCauHinhMaster, []interface{}{
		core.PrePhanQuyen + pq.MaVaiTro, 
		string(b),
	})
}

func Repo_Update(masterID string, pq *core.PhanQuyen) {
	lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
	lock.Lock()
	lock.Unlock()

	b, _ := json.Marshal(pq)
	core.PushUpdate(masterID, core.TenSheetCauHinhMaster, pq.DongTrongSheet, core.CotCH_DataJSON, string(b))
}
