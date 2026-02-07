package bao_mat

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// [CẤU HÌNH] Mã bí mật chỉ Server biết (Tuyệt đối không lộ ra ngoài)
// Dùng để tạo chữ ký cho Cookie, chống làm giả hoặc copy cookie sang máy khác
const SECRET_KEY = "MayTinhShop_@2026_Secret_Key_!@#$$^KeepItSecret"

// HashMatKhau : Băm mật khẩu bằng Bcrypt
func HashMatKhau(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// KiemTraMatKhau : So sánh mật khẩu nhập vào và Hash
func KiemTraMatKhau(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TaoSessionIDAnToan : Tạo chuỗi ngẫu nhiên dài 128 ký tự (64 byte hex)
func TaoSessionIDAnToan() string {
	b := make([]byte, 64) // 64 byte
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// --- [MỚI] Hàm Tạo Chữ Ký Bảo Mật (HMAC-SHA256) ---
// Input: SessionID + UserAgent (Trình duyệt)
// Output: Chuỗi ký tự (Signature)
func TaoChuKyBaoMat(sessionID string, userAgent string) string {
	// Dữ liệu cần ký = SessionID + UserAgent
	// Kẻ tấn công có sessionID nhưng khác UserAgent sẽ tạo ra chữ ký sai
	data := sessionID + "___" + userAgent
	
	// Tạo HMAC dùng SHA256 và Secret Key
	h := hmac.New(sha256.New, []byte(SECRET_KEY))
	h.Write([]byte(data))
	
	// Trả về dạng Hex string
	return hex.EncodeToString(h.Sum(nil))
}
