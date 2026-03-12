package auth

import (
	"encoding/json"
	"app/core"
)

type Repo struct{}

func (r *Repo) FindByUserOrEmail(shopID, input string) (*core.KhachHang, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.RLock()
	defer lock.RUnlock()

	for _, kh := range core.CacheKhachHang[shopID] {
		if kh.TenDangNhap == input || kh.Email == input {
			return kh, true
		}
	}
	return nil, false
}

func (r *Repo) CountUsers(shopID string) int {
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.RLock()
	defer lock.RUnlock()
	return len(core.CacheKhachHang[shopID])
}

// Ghi chớp nhoáng RAM + Ghi Hàng Đợi Background (Chỉ 2 Cột)
func (r *Repo) InsertUser(shopID string, kh *core.KhachHang) {
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], kh)
	core.CacheMapKhachHang[core.TaoCompositeKey(shopID, kh.MaKhachHang)] = kh
	lock.Unlock()

	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	
	// Đóng gói JSON và nã 1 phát xuống cột A, B
	b, _ := json.Marshal(kh)
	core.PushAppend(sID, core.TenSheetKhachHang, []interface{}{kh.MaKhachHang, string(b)})
}

// 1 Hàm duy nhất thay thế cho mọi hàm Update riêng lẻ
func (r *Repo) UpdateUserJSON(shopID string, dong int, jsonStr string) {
	core.PushUpdate(shopID, core.TenSheetKhachHang, dong, core.CotKH_DataJSON, jsonStr)
}

func (r *Repo) SendWelcomeMessage(shopID string, msg *core.TinNhan) {
	core.ThemMoiTinNhan(shopID, msg)
}
