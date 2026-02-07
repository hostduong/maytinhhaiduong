package nghiep_vu

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const (
	URL_API_MAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_MAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
)

type ThongTinOTP struct { MaCode string; HetHanLuc int64 }
type BoDemRate struct { LanGuiCuoi int64 } // Rút gọn struct

var CacheOTP = make(map[string]ThongTinOTP)
var CacheRate = make(map[string]*BoDemRate)
var mtxOTP sync.Mutex

func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

func LuuOTP(userKey string, code string) {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	CacheOTP[userKey] = ThongTinOTP{MaCode: code, HetHanLuc: time.Now().Add(10 * time.Minute).Unix()}
}

func KiemTraOTP(userKey string, inputCode string) bool {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	otp, ok := CacheOTP[userKey]
	if !ok || time.Now().Unix() > otp.HetHanLuc { return false }
	if otp.MaCode == inputCode { delete(CacheOTP, userKey); return true }
	return false
}

func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

// [CẬP NHẬT]: Chỉ giới hạn 1 phút 1 lần, bỏ giới hạn 6h
func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	now := time.Now().Unix()
	
	rd, ok := CacheRate[email]
	if !ok {
		// Chưa gửi lần nào -> Tạo mới
		CacheRate[email] = &BoDemRate{LanGuiCuoi: now}
		return true, ""
	}

	// Kiểm tra thời gian chờ (60s)
	if now - rd.LanGuiCuoi < 60 {
		return false, fmt.Sprintf("Vui lòng đợi %d giây nữa.", 60-(now-rd.LanGuiCuoi))
	}

	// Cập nhật thời gian gửi cuối
	rd.LanGuiCuoi = now
	return true, ""
}

func GuiMailXacMinhAPI(email, code string) error {
	return callApi(map[string]string{"type": "sender_mail", "api_key": KEY_API_MAIL, "email": email, "code": code})
}

func GuiMailThongBaoAPI(email, subject, name, body string) error {
	return callApi(map[string]string{"type": "sender", "api_key": KEY_API_MAIL, "email": email, "subject": subject, "name": name, "body": body})
}

func callApi(payload interface{}) error {
	b, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(b))
	if err != nil { return fmt.Errorf("Lỗi kết nối đến Google: %v", err) }
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	
	var r struct { Status string `json:"status"`; Messenger string `json:"messenger"` }
	if err := json.Unmarshal(bodyBytes, &r); err != nil { return fmt.Errorf("Lỗi định dạng phản hồi") }

	if r.Status == "true" { return nil }
	if r.Messenger != "" { return fmt.Errorf("%s", r.Messenger) }
	return fmt.Errorf("Gửi mail thất bại")
}
