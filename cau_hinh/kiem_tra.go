package cau_hinh

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

// --- PHẦN MÃ HÓA ---
const SECRET_KEY = "MayTinhShop_@2026_Secret_Key_!@#$$^KeepItSecret"

func HashMatKhau(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func KiemTraMatKhau(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TaoSessionIDAnToan() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil { return "" }
	return hex.EncodeToString(b)
}

func TaoChuKyBaoMat(sessionID string, userAgent string) string {
	data := sessionID + "___" + userAgent
	h := hmac.New(sha256.New, []byte(SECRET_KEY))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// --- PHẦN KIỂM TRA ĐẦU VÀO ---
func KiemTraTenDangNhap(user string) bool {
	if len(user) < 6 || len(user) > 30 { return false }
	
	// Chỉ cho phép chữ thường, số và dấu gạch ngang
	match, _ := regexp.MatchString(`^[a-z0-9\-]+$`, user)
	if !match { return false }
	
	// Không được bắt đầu, kết thúc bằng dấu gạch ngang, hoặc chứa 2 dấu gạch ngang liên tiếp
	if strings.HasPrefix(user, "-") || strings.HasSuffix(user, "-") || strings.Contains(user, "--") { 
		return false 
	}
	
	return true
}

func KiemTraEmail(email string) bool {
	if len(email) < 6 || len(email) > 100 { return false }
	match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@(?:[a-z0-9-]+\.)+[a-z]{2,}$`, email)
	return match
}

func KiemTraDinhDangMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 { return false }
	match, _ := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()\-+_.,?]+$`, pass)
	return match
}

func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

func KiemTraHoTen(name string) bool {
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 { return false }
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}
