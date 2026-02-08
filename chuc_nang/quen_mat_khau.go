package chuc_nang

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core" // [MỚI]

	"github.com/gin-gonic/gin"
)

// --- BỘ NHỚ OTP TẠM THỜI (Thay thế nghiep_vu) ---
var (
	cacheOTPMemory = make(map[string]string) // Map[User]OTP
	mtxOTP         sync.Mutex
)

// Helper: Lưu OTP
func luuOTPCucBo(user, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	cacheOTPMemory[user] = code
	// Tự xóa sau 5 phút
	go func(u string) {
		time.Sleep(5 * time.Minute)
		mtxOTP.Lock()
		delete(cacheOTPMemory, u)
		mtxOTP.Unlock()
	}(user)
}

// Helper: Kiểm tra OTP
func kiemTraOTPCucBo(user, code string) bool {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	if val, ok := cacheOTPMemory[user]; ok && val == code {
		delete(cacheOTPMemory, user) // Xóa sau khi dùng
		return true
	}
	return false
}

// Helper: Tạo mã 6 số
func taoMaOTP6So() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// --- LOGIC CHÍNH ---

func TrangQuenMatKhau(c *gin.Context) { c.HTML(http.StatusOK, "quen_mat_khau", gin.H{}) }

// [CÁCH 1]: Đổi mật khẩu bằng Mã PIN
func XuLyQuenPassBangPIN(c *gin.Context) {
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	pinInput := strings.TrimSpace(c.PostForm("pin"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))
	
	if !bao_mat.KiemTraMaPin(pinInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN phải đúng 8 chữ số!"})
		return
	}
	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	kh, ok := core.TimKhachHangTheoUserOrEmail(dinhDanh)
	
	if !ok || !bao_mat.KiemTraMatKhau(pinInput, kh.MaPinHash) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản hoặc mã PIN không chính xác!"})
		return 
	}
	
	// Hash mật khẩu mới và lưu vào Core
	hash, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hash
	
	sID := kh.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
	
	core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

// [CÁCH 2 - BƯỚC 1]: Gửi OTP
func XuLyGuiOTPEmail(c *gin.Context) {
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	
	kh, ok := core.TimKhachHangTheoUserOrEmail(dinhDanh)
	if !ok { 
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu tài khoản tồn tại, mã OTP sẽ được gửi đến Email."})
		return 
	}

	if kh.Email == "" || !strings.Contains(kh.Email, "@") {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản chưa có Email, vui lòng dùng PIN."})
		return
	}

	code := taoMaOTP6So()
	
	// Gửi mail (GIẢ LẬP - In ra console để test)
	// TODO: Tích hợp thư viện mail thật ở đây
	log.Printf("📧 [MAIL MOCK] Gửi OTP '%s' đến %s", code, kh.Email)
	
	// Lưu OTP vào bộ nhớ cục bộ
	luuOTPCucBo(kh.TenDangNhap, code)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã OTP (Kiểm tra Console Log nếu đang test)!"})
}

// [CÁCH 2 - BƯỚC 2]: Xác nhận OTP và Đổi Pass
func XuLyQuenPassBangOTP(c *gin.Context) {
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	otp      := strings.TrimSpace(c.PostForm("otp"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))

	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	kh, ok := core.TimKhachHangTheoUserOrEmail(dinhDanh)
	
	if !ok || !kiemTraOTPCucBo(kh.TenDangNhap, otp) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP không đúng hoặc đã hết hạn!"})
		return 
	}

	hash, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hash
	
	sID := kh.SpreadsheetID
	if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }

	core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
