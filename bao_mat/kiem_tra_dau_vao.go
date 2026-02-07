package bao_mat

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 1. Tên Đăng Nhập
// Quy tắc: 6-30 ký tự, a-z, 0-9, . _
// Không bắt đầu/kết thúc bằng . _
// Không chứa .. hoặc __ hoặc ._ hoặc _.
func KiemTraTenDangNhap(user string) bool {
	// 1. Check độ dài (6-30)
	if len(user) < 6 || len(user) > 30 {
		return false
	}

	// 2. Check ký tự hợp lệ (Chỉ a-z, 0-9, ., _)
	// (Regex này an toàn với RE2 của Go)
	match, _ := regexp.MatchString(`^[a-z0-9._]+$`, user)
	if !match {
		return false
	}

	// 3. Check ký tự đầu và cuối (Phải là chữ hoặc số, không được là . hoặc _)
	firstChar := user[0]
	lastChar := user[len(user)-1]
	if !isAlphaNumeric(firstChar) || !isAlphaNumeric(lastChar) {
		return false
	}

	// 4. Check các cặp ký tự cấm liên tiếp (Thay thế cho Lookahead)
	// Chặn .. và __ và ._ và _.
	forbiddenPairs := []string{"..", "__", "._", "_."}
	for _, pair := range forbiddenPairs {
		if strings.Contains(user, pair) {
			return false
		}
	}

	return true
}

// 2. Email
// Quy tắc: 6-100 ký tự, định dạng chuẩn, chặn @mail..com
func KiemTraEmail(email string) bool {
	// 1. Check độ dài (6-100)
	if len(email) < 6 || len(email) > 100 {
		return false
	}

	// 2. Regex định dạng chuẩn
	// Phần domain: (?:[a-z0-9-]+\.)+ nghĩa là (Cụm-từ + Chấm) lặp lại ít nhất 1 lần
	// Điều này giúp chặn đứng trường hợp "mail..com" hoặc "mail" (không có chấm)
	match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@(?:[a-z0-9-]+\.)+[a-z]{2,}$`, email)
	return match
}

// 3. Mật khẩu
// Quy tắc: 8-30 ký tự, Whitelist ký tự cho phép
func KiemTraDinhDangMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 {
		return false
	}
	// Whitelist: a-z, A-Z, 0-9 và các ký tự đặc biệt cụ thể
	match, _ := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()\-+_.,?]+$`, pass)
	return match
}

// 4. Mã PIN
// Quy tắc: Đúng 8 số
func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

// 5. Họ Tên
// Quy tắc: 6-50 ký tự, Unicode (Tiếng Việt), khoảng trắng.
func KiemTraHoTen(name string) bool {
	name = strings.TrimSpace(name)
	// Đếm số ký tự thực (Rune) thay vì byte để hỗ trợ tiếng Việt chính xác
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 {
		return false
	}
	
	// \p{L} là chữ cái Unicode, \s là khoảng trắng
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}

// Helper kiểm tra chữ/số (byte - dùng cho ASCII a-z 0-9)
func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
