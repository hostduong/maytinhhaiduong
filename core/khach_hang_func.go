package core

import (
	"strings"
)

func TimKhachHangTheoUserOrEmail(shopID, input string) (*KhachHang, bool) {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	input = strings.ToLower(strings.TrimSpace(input))
	list := CacheKhachHang[shopID]
	for _, kh := range list {
		if strings.ToLower(kh.TenDangNhap) == input || (kh.Email != "" && strings.ToLower(kh.Email) == input) {
			if kh.MaKhachHang == "0000000000000000000" { return nil, false } // Chặn Bot
			return kh, true
		}
	}
	return nil, false
}

func TaoMaKhachHangMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	for {
		id := LayChuoiSoNgauNhien(19)
		key := TaoCompositeKey(shopID, id)
		if _, exist := CacheMapKhachHang[key]; !exist { return id }
	}
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	sID := kh.SpreadsheetID
	if sID == "" { sID = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8" }
	CacheKhachHang[sID] = append(CacheKhachHang[sID], kh)
	key := TaoCompositeKey(sID, kh.MaKhachHang)
	CacheMapKhachHang[key] = kh
}
