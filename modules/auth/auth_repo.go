package auth

import "app/core"

type Repo struct{}

// FindByUserOrEmail: RLock (Chỉ Đọc) - Tốc độ ánh sáng
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

// CountUsers: RLock (Chỉ Đọc)
func (r *Repo) CountUsers(shopID string) int {
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.RLock()
	defer lock.RUnlock()
	return len(core.CacheKhachHang[shopID])
}

// InsertUser: Lock (Ghi chớp nhoáng RAM) + Ghi Hàng Đợi Background
func (r *Repo) InsertUser(shopID string, kh *core.KhachHang) {
	// 1. Khóa Độc Quyền để Ghi vào RAM
	lock := core.GetSheetLock(shopID, core.TenSheetKhachHang)
	lock.Lock()
	core.CacheKhachHang[shopID] = append(core.CacheKhachHang[shopID], kh)
	core.CacheMapKhachHang[core.TaoCompositeKey(shopID, kh.MaKhachHang)] = kh
	lock.Unlock()

	// 2. Đóng gói đẩy xuống Queue chạy ngầm (Không bắt người dùng đợi)
	sID := kh.SpreadsheetID; if sID == "" { sID = shopID }
	
	rowData := make([]interface{}, 26)
	rowData[core.CotKH_MaKhachHang] = kh.MaKhachHang; rowData[core.CotKH_TenDangNhap] = kh.TenDangNhap
	rowData[core.CotKH_Email] = kh.Email; rowData[core.CotKH_MatKhauHash] = kh.MatKhauHash
	rowData[core.CotKH_MaPinHash] = kh.MaPinHash; rowData[core.CotKH_VaiTroQuyenHan] = kh.VaiTroQuyenHan
	rowData[core.CotKH_ChucVu] = kh.ChucVu; rowData[core.CotKH_TrangThai] = kh.TrangThai
	rowData[core.CotKH_TenKhachHang] = kh.TenKhachHang; rowData[core.CotKH_DienThoai] = kh.DienThoai
	rowData[core.CotKH_NgaySinh] = kh.NgaySinh; rowData[core.CotKH_GioiTinh] = kh.GioiTinh
	rowData[core.CotKH_NguonKhachHang] = kh.NguonKhachHang; rowData[core.CotKH_NgayTao] = kh.NgayTao
	rowData[core.CotKH_NguoiCapNhat] = kh.NguoiCapNhat; rowData[core.CotKH_NgayCapNhat] = kh.NgayCapNhat
	rowData[core.CotKH_RefreshTokenJson] = core.ToJSON(kh.RefreshTokens)
	
	core.PushAppend(sID, core.TenSheetKhachHang, rowData)
}

func (r *Repo) UpdateTokens(shopID string, dong int, tokens map[string]core.TokenInfo) {
	core.PushUpdate(shopID, core.TenSheetKhachHang, dong, core.CotKH_RefreshTokenJson, core.ToJSON(tokens))
}

func (r *Repo) UpdatePassword(shopID string, dong int, hash string) {
	core.PushUpdate(shopID, core.TenSheetKhachHang, dong, core.CotKH_MatKhauHash, hash)
}

func (r *Repo) UpdateGoiDichVu(shopID string, dong int, goi []core.PlanInfo) {
	core.PushUpdate(shopID, core.TenSheetKhachHang, dong, core.CotKH_GoiDichVuJson, core.ToJSON(goi))
}

func (r *Repo) SendWelcomeMessage(shopID string, msg *core.TinNhan) {
	// (Tin nhắn sẽ được xử lý Lock riêng tại file core.ThemMoiTinNhan)
	core.ThemMoiTinNhan(shopID, msg)
}
