package auth_master

import (
	"strings"

	"app/config"
	"app/core"
)

type Repo struct{}

// Ép cứng ID File Master, không lấy từ tham số động để chống leo thang đặc quyền
func getMasterID() string { return config.BienCauHinh.IdFileSheetMaster }

func (r *Repo) FindByUserOrEmail(input string) (*core.KhachHang, bool) {
	masterID := getMasterID()
	lock := core.GetSheetLock(masterID, core.TenSheetKhachHangMaster)
	lock.RLock()
	defer lock.RUnlock()

	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range core.CacheKhachHang[masterID] {
		// Bỏ qua tài khoản BOT (000)
		if kh.MaKhachHang == "0000000000000000000" { continue }
		
		if strings.ToLower(kh.TenDangNhap) == input || (kh.Email != "" && strings.ToLower(kh.Email) == input) {
			return kh, true
		}
	}
	return nil, false
}

// Hàm đẩy JSON cập nhật thẳng vào File Master
func (r *Repo) UpdateUserJSON(dong int, jsonStr string) {
	core.PushUpdate(getMasterID(), core.TenSheetKhachHangMaster, dong, core.CotKH_DataJSON, jsonStr)
}
