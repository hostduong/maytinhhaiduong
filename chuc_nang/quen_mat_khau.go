package chuc_nang

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"app/cau_hinh" // Chứa hàm Hash/Check và Validation
	"app/core"

	"github.com/gin-gonic/gin"
)

// --- BỘ NHỚ OTP TẠM THỜI (RAM) ---
// Map[Key]OTP. Key sẽ là "ShopID_Username" để tránh trùng giữa các shop
var (
	cacheOTPMemory = make(map[string]string) 
	mtxOTP         sync.Mutex
)

// Helper: Lưu OTP (Key = ShopID + User)
func luuOTPCucBo(shopID, user, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	
	key := shopID + "_" + user
	cacheOTPMemory[key] = code
	
	// Tự xóa sau 5 phút
	go func(k string) {
		time.Sleep(5 * time.Minute)
		mtxOTP.Lock()
		delete(cacheOTPMemory, k)
		mtxOTP.Unlock()
	}(key)
}

// Helper: Kiểm tra OTP
func kiemTraOTPCucBo(shopID, user, code string) bool {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	
	key := shopID + "_" + user
	if val, ok := cacheOTPMemory[key]; ok && val == code {
		delete(cacheOTPMemory, key) // Xóa sau khi dùng xong
		return true
	}
	return false
}

// Helper: Tạo mã 6 số
func taoMaOTP6So() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// =============================================================
// LOGIC CHÍNH
// =============================================================

func TrangQuenMatKhau(c *gin.Context) { 
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{}) 
}

// [CÁCH 1]: Đổi mật khẩu bằng Mã PIN
func XuLyQuenPassBangPIN(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	theme := c.GetString("THEME") // [SAAS] Lấy theme động

	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	pinInput := strings.TrimSpace(c.PostForm("pin"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))
	
	// Validate
	if !cau_hinh.KiemTraMaPin(pinInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN phải đúng 8 chữ số!"})
		return
	}
	if !cau_hinh.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	// [SAAS] Tìm user trong Shop
	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, dinhDanh)
	
	if !ok || !cau_hinh.KiemTraMatKhau(pinInput, kh.MaPinHash) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản hoặc mã PIN không chính xác!"})
		return 
	}
	
	// Hash mật khẩu mới và lưu vào Core RAM
	hash, _ := cau_hinh.HashMatKhau(passMoi)
	
	core.KhoaHeThong.Lock()
	kh.MatKhauHash = hash
	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()
	
	// Ghi Sheet (Truyền shopID)
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

// [CÁCH 2 - BƯỚC 1]: Gửi OTP
func XuLyGuiOTPEmail(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	
	// Tìm user
	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, dinhDanh)
	if !ok { 
		// Fake success để tránh dò user
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu tài khoản tồn tại, mã OTP sẽ được gửi đến Email."})
		return 
	}

	if kh.Email == "" || !strings.Contains(kh.Email, "@") {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản chưa có Email, vui lòng dùng PIN."})
		return
	}

	code := taoMaOTP6So()
	
	// Gửi mail (GIẢ LẬP - In ra console để test)
	log.Printf("📧 [MAIL MOCK] Shop [%s] - Gửi OTP '%s' đến %s", shopID, code, kh.Email)
	
	// Lưu OTP vào bộ nhớ cục bộ (Kèm ShopID)
	luuOTPCucBo(shopID, kh.TenDangNhap, code)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã OTP (Kiểm tra Console Log nếu đang test)!"})
}

// [CÁCH 2 - BƯỚC 2]: Xác nhận OTP và Đổi Pass
func XuLyQuenPassBangOTP(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]

	dinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("dinh_danh")))
	otp      := strings.TrimSpace(c.PostForm("otp"))
	passMoi  := strings.TrimSpace(c.PostForm("pass_moi"))

	if !cau_hinh.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	// Tìm user
	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, dinhDanh)
	
	// Check OTP (Kèm ShopID)
	if !ok || !kiemTraOTPCucBo(shopID, kh.TenDangNhap, otp) { 
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP không đúng hoặc đã hết hạn!"})
		return 
	}

	// Đổi pass
	hash, _ := cau_hinh.HashMatKhau(passMoi)
	
	core.KhoaHeThong.Lock()
	kh.MatKhauHash = hash
	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
	core.KhoaHeThong.Unlock()
	
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
	core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_NgayCapNhat, kh.NgayCapNhat)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
