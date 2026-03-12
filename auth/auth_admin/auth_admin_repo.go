package auth_admin

import (
	"encoding/json"
	"app/config"
	"app/core"
)

type Repo struct{}

func getAdminID() string { return config.BienCauHinh.IdFileSheetAdmin }

func (r *Repo) FindByUserOrEmail(input string) (*core.KhachHang, bool) {
	adminID := getAdminID()
	lock := core.GetSheetLock(adminID, core.TenSheetKhachHangAdmin)
	lock.RLock()
	defer lock.RUnlock()

	for _, kh := range core.CacheKhachHang[adminID] {
		if kh.TenDangNhap == input || kh.Email == input {
			return kh, true
		}
	}
	return nil, false
}

func (r *Repo) CountUsers() int {
	adminID := getAdminID()
	lock := core.GetSheetLock(adminID, core.TenSheetKhachHangAdmin)
	lock.RLock()
	defer lock.RUnlock()
	return len(core.CacheKhachHang[adminID])
}

// Bắn thẳng xuống Sheet KHACH_HANG_ADMIN
func (r *Repo) InsertUser(kh *core.KhachHang) {
	adminID := getAdminID()
	lock := core.GetSheetLock(adminID, core.TenSheetKhachHangAdmin)
	lock.Lock()
	core.CacheKhachHang[adminID] = append(core.CacheKhachHang[adminID], kh)
	core.CacheMapKhachHang[core.TaoCompositeKey(adminID, kh.MaKhachHang)] = kh
	lock.Unlock()

	b, _ := json.Marshal(kh)
	core.PushAppend(adminID, core.TenSheetKhachHangAdmin, []interface{}{kh.MaKhachHang, string(b)})
}

func (r *Repo) UpdateUserJSON(dong int, jsonStr string) {
	core.PushUpdate(getAdminID(), core.TenSheetKhachHangAdmin, dong, core.CotKH_DataJSON, jsonStr)
}

func (r *Repo) SendWelcomeMessage(msg *core.TinNhan) {
	core.ThemMoiTinNhan(getAdminID(), msg)
}
