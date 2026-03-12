package auth_verify

import (
	"strings"
	"sync"
	"time"

	"app/core"
)

type ThongTinOTP struct {
	MaCode    string
	HetHanLuc int64
}

type BoDemRate struct {
	LanGuiCuoi int64
}

type Repo struct {
	muOTP     sync.Mutex
	cacheOTP  map[string]ThongTinOTP
	cacheRate map[string]*BoDemRate
}

// Khởi tạo Repo với bản đồ bộ nhớ riêng cho OTP
func NewRepo() *Repo {
	return &Repo{
		cacheOTP:  make(map[string]ThongTinOTP),
		cacheRate: make(map[string]*BoDemRate),
	}
}

// --- TƯƠNG TÁC GOOGLE SHEETS (TÌM USER) ---
func (r *Repo) FindUser(appMode, shopID, input string) (*core.KhachHang, bool) {
	sheetName := core.TenSheetKhachHang
	targetID := shopID

	if appMode == "MASTER_CORE" {
		sheetName = core.TenSheetKhachHangMaster
	} else if appMode == "TENANT_ADMIN" {
		sheetName = core.TenSheetKhachHangAdmin
	}

	lock := core.GetSheetLock(targetID, sheetName)
	lock.RLock()
	defer lock.RUnlock()

	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range core.CacheKhachHang[targetID] {
		if strings.ToLower(kh.TenDangNhap) == input || (kh.Email != "" && strings.ToLower(kh.Email) == input) {
			if kh.MaKhachHang != "0000000000000000000" { return kh, true }
		}
	}
	return nil, false
}

// --- QUẢN TRỊ BỘ NHỚ OTP TẠI CHỖ ---
func (r *Repo) LuuOTP(userKey string, code string) {
	r.muOTP.Lock()
	defer r.muOTP.Unlock()
	r.cacheOTP[userKey] = ThongTinOTP{
		MaCode:    code,
		HetHanLuc: time.Now().Add(10 * time.Minute).Unix(),
	}
}

func (r *Repo) KiemTraVaXoaOTP(userKey string, inputCode string) bool {
	r.muOTP.Lock()
	defer r.muOTP.Unlock()
	
	otp, ok := r.cacheOTP[userKey]
	if !ok || time.Now().Unix() > otp.HetHanLuc { return false }
	
	if otp.MaCode == inputCode {
		delete(r.cacheOTP, userKey) // Dùng xong là xóa ngay (Tiêu hủy)
		return true
	}
	return false
}

func (r *Repo) CheckRateLimit(email string) (bool, int64) {
	r.muOTP.Lock()
	defer r.muOTP.Unlock()
	
	now := time.Now().Unix()
	rd, ok := r.cacheRate[email]
	
	if !ok {
		r.cacheRate[email] = &BoDemRate{LanGuiCuoi: now}
		return true, 0
	}

	if now-rd.LanGuiCuoi < 60 {
		return false, 60 - (now - rd.LanGuiCuoi) // Trả về số giây phải đợi
	}

	rd.LanGuiCuoi = now
	return true, 0
}
