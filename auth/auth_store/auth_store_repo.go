package auth_store

import (
	"encoding/json"
	"strings"

	"app/core"
)

type Repo struct{}

// Hàm tìm user linh hoạt theo ShopID của từng cửa hàng
func (r *Repo) FindByUserOrEmail(shopID, input string) (*core.KhachHang, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.RLock()
	defer lock.RUnlock()

	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range core.CacheKhachHang[shopID] {
		if strings.ToLower(kh.TenDangNhap) == input || (kh.Email != "" && strings.ToLower(kh.Email) == input) {
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

// Bơm dữ liệu JSON thẳng xuống RAM và File của đúng cửa hàng đó
func (r *Repo) InsertUser(shopID string, kh *core.KhachHang) {
	// 1. Khóa Phân Mảnh
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], kh)
	lock.Unlock()

	// 2. [BẢN VÁ] Khóa Toàn Cục
	core.KhoaHeThong.Lock()
	core.CacheMapKhachHang[core.TaoCompositeKey(shopID, kh.MaKhachHang)] = kh
	core.KhoaHeThong.Unlock()

	b, _ := json.Marshal(kh)
	core.PushAppend(shopID, core.TenSheetKhachHang, []interface{}{kh.MaKhachHang, string(b)})
}
func (r *Repo) UpdateUserJSON(shopID string, dong int, jsonStr string) {
	core.PushUpdate(shopID, core.TenSheetKhachHang, dong, core.CotKH_DataJSON, jsonStr)
}
