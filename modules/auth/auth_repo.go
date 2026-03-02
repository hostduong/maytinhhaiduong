package auth

import "app/core"

type Repo struct{}

func (r *Repo) FindByUserOrEmail(shopID, input string) (*core.KhachHang, bool) {
	return core.TimKhachHangTheoUserOrEmail(shopID, input)
}

func (r *Repo) CountUsers(shopID string) int {
	return len(core.LayDanhSachKhachHang(shopID))
}

func (r *Repo) InsertUser(kh *core.KhachHang) {
	core.ThemKhachHangVaoRam(kh)
	sID := kh.SpreadsheetID; if sID == "" { sID = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8" }
	
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
	core.ThemMoiTinNhan(shopID, msg)
}
