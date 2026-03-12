package auth_verify

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"app/config"
	"app/core"
)

const (
	URL_API_MAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_MAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
)

type Service struct { repo *Repo }

func NewService(r *Repo) *Service { return &Service{repo: r} }

// Sinh mã OTP 6 số bảo mật
func taoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

// Bắn API sang Google Apps Script
func guiMailXacMinhAPI(email, code string) error {
	payload := map[string]string{
		"type":    "sender_mail",
		"api_key": KEY_API_MAIL,
		"email":   email,
		"code":    code,
	}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(b))
	if err != nil { return errors.New("Lỗi kết nối máy chủ Mail") }
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var r struct { Status string `json:"status"`; Messenger string `json:"messenger"` }
	json.Unmarshal(bodyBytes, &r)

	if r.Status == "true" { return nil }
	if r.Messenger != "" { return errors.New(r.Messenger) }
	return errors.New("Gửi mail thất bại")
}

// --- CÁC HÀM NGHIỆP VỤ ---

func (s *Service) SendOtp(appMode, shopID, identifier string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }

	kh, ok := s.repo.FindUser(appMode, shopID, identifier)
	if !ok { return errors.New("Tài khoản không tồn tại trên hệ thống!") }
	if kh.Email == "" || !config.KiemTraEmail(kh.Email) { 
		return errors.New("Tài khoản chưa có Email hợp lệ, vui lòng dùng Mã PIN.") 
	}

	// Xử lý Rate Limit
	okLimit, waitTime := s.repo.CheckRateLimit(kh.Email)
	if !okLimit { return fmt.Errorf("Vui lòng đợi %d giây nữa trước khi gửi lại.", waitTime) }

	// Xử lý Gửi và Lưu OTP
	code := taoMaOTP6So()
	if err := guiMailXacMinhAPI(kh.Email, code); err != nil { return err }
	
	s.repo.LuuOTP(shopID+"_"+kh.TenDangNhap, code)
	return nil
}

func (s *Service) VerifyOtp(appMode, shopID, identifier, otp string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }
	
	kh, ok := s.repo.FindUser(appMode, shopID, identifier)
	if !ok { return errors.New("Tài khoản không tồn tại!") }

	if !s.repo.KiemTraVaXoaOTP(shopID+"_"+kh.TenDangNhap, otp) { 
		return errors.New("Mã OTP không đúng hoặc đã hết hạn!") 
	}
	return nil
}

func (s *Service) VerifyPin(appMode, shopID, identifier, pinInput string) error {
	if err := core.EnsureKhachHangLoaded(shopID); err != nil { return err }

	kh, ok := s.repo.FindUser(appMode, shopID, identifier)
	if !ok { return errors.New("Tài khoản không tồn tại!") }

	if !config.KiemTraMatKhau(pinInput, kh.BaoMat.MaPinHash) { 
		return errors.New("Mã PIN bảo mật không chính xác!") 
	}
	return nil
}
