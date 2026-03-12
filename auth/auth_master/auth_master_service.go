package auth_master

import (
	"encoding/json"
	"errors"
	"time"

	"app/config"
	"app/core"
)

type Service struct { repo Repo }

func (s *Service) Login(dinhDanh, pass, userAgent string, ghiNho bool) (string, string, error) {
	masterID := config.BienCauHinh.IdFileSheetMaster
	if err := core.EnsureKhachHangLoaded(masterID); err != nil { return "", "", err }

	kh, ok := s.repo.FindByUserOrEmail(dinhDanh)
	if !ok { return "", "", errors.New("Tài khoản Master không tồn tại!") }
	if !config.KiemTraMatKhau(pass, kh.BaoMat.MatKhauHash) { return "", "", errors.New("Mật khẩu không đúng!") }
	if kh.TrangThai == 0 { return "", "", errors.New("Tài khoản Master đang bị phong tỏa!") }

	thoiGianSong := config.ThoiGianHetHanCookie
	if ghiNho { thoiGianSong = 30 * 24 * time.Hour }

	sessionID := config.TaoSessionIDAnToan()
	signature := config.TaoChuKyBaoMat(sessionID, userAgent)
	expTime := time.Now().Add(thoiGianSong).Unix()

	lock := core.GetSheetLock(masterID, core.TenSheetKhachHangMaster)
	lock.Lock()
	if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]core.TenantDeviceToken) }
	
	nowUnix := time.Now().Unix()
	for key, info := range kh.RefreshTokens { if info.Exp < nowUnix { delete(kh.RefreshTokens, key) } }
	if len(kh.RefreshTokens) >= 5 {
		var oldestKey string; var oldestTime int64 = 1<<63 - 1
		for key, info := range kh.RefreshTokens { if info.Exp < oldestTime { oldestTime = info.Exp; oldestKey = key } }
		if oldestKey != "" { delete(kh.RefreshTokens, oldestKey) }
	}
	
	kh.RefreshTokens[sessionID] = core.TenantDeviceToken{DeviceID: sessionID, Dev: userAgent, Exp: expTime, Created: nowUnix}
	
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()

	go s.repo.UpdateUserJSON(kh.DongTrongSheet, jsonStr) 
	return sessionID, signature, nil
}

func (s *Service) ResetPassword(dinhDanh, passMoi string) error {
	masterID := config.BienCauHinh.IdFileSheetMaster
	kh, ok := s.repo.FindByUserOrEmail(dinhDanh)
	if !ok { return errors.New("Không tìm thấy tài khoản Master!") }

	hash, _ := config.HashMatKhau(passMoi)
	lock := core.GetSheetLock(masterID, core.TenSheetKhachHangMaster)
	lock.Lock()
	kh.BaoMat.MatKhauHash = hash
	b, _ := json.Marshal(kh)
	jsonStr := string(b)
	lock.Unlock()
	
	s.repo.UpdateUserJSON(kh.DongTrongSheet, jsonStr)
	return nil
}
