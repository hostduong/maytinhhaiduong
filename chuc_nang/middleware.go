package chuc_nang

import (
	"net/http"
	"sync"
	"time"

	"app/bao_mat"
	"app/bo_nho_dem" // [MỚI] Import để check cờ HeThongDangBan
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// Bộ nhớ đếm Request cho Rate Limit
var boDem = make(map[string]int)
var mtx sync.Mutex

// Khởi chạy bộ đếm (Reset mỗi giây)
func KhoiTaoBoDemRateLimit() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mtx.Lock()
			boDem = make(map[string]int) // Xóa sạch bộ đếm cũ
			mtx.Unlock()
		}
	}()
}

// MIDDLEWARE CHÍNH
func KiemTraQuyenHan(c *gin.Context) {
	// [MỚI] 1. CHỐT CHẶN BẢO TRÌ (Sửa biến trỏ về bo_nho_dem)
	if bo_nho_dem.HeThongDangBan && c.Request.Method != "GET" {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"trang_thai": "he_thong_ban",
			"thong_diep": "Hệ thống đang đồng bộ dữ liệu, vui lòng thử lại sau 5 giây.",
		})
		return
	}

	// 2. KIỂM TRA RATE LIMIT (CHỐNG SPAM)
	cookieID, err1 := c.Cookie("session_id")
	cookieSign, err2 := c.Cookie("session_sign")
	
	keyLimit := ""
	if err1 != nil || cookieID == "" {
		keyLimit = "LIMIT__IP__" + c.ClientIP()
	} else {
		keyLimit = "LIMIT__COOKIE__" + cookieID
	}

	mtx.Lock()
	boDem[keyLimit]++
	soLanGoi := boDem[keyLimit]
	mtx.Unlock()

	if soLanGoi > cau_hinh.GioiHanNguoiDung {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"loi": "Thao tác quá nhanh! Vui lòng chậm lại."})
		return
	}

	// 3. KIỂM TRA ĐĂNG NHẬP & BẢO MẬT (AUTH)
	if cookieID == "" {
		c.Next()
		return
	}

	// KIỂM TRA TÍNH TOÀN VẸN (SECURITY CHECK)
	userAgent := c.Request.UserAgent()
	signatureServer := bao_mat.TaoChuKyBaoMat(cookieID, userAgent)

	if err2 != nil || cookieSign != signatureServer {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.SetCookie("session_sign", "", -1, "/", "", false, true)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"loi": "Phát hiện bất thường (Cookie Mismatch)! Vui lòng đăng nhập lại."})
		return
	}

	// 4. TÌM USER TRONG RAM (Gọi qua nghiep_vu - logic này đã được update bên trong nghiep_vu ở Bước 2)
	khachHang, timThay := nghiep_vu.TimKhachHangTheoCookie(cookieID)

	if !timThay {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.SetCookie("session_sign", "", -1, "/", "", false, true)
		c.Next()
		return
	}

	// 5. LOGIC GIA HẠN THÔNG MINH (Auto-Renew)
	thoiGianHetHan := khachHang.CookieExpired // Dạng int64
	now := time.Now().Unix()

	// Nếu đã hết hạn -> Đá ra
	if now > thoiGianHetHan {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.SetCookie("session_sign", "", -1, "/", "", false, true)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"loi": "Phiên đăng nhập hết hạn"})
		return
	}

	// Nếu còn hạn nhưng sắp hết (trong vùng ân hạn) -> GIA HẠN
	thoiGianConLai := time.Duration(thoiGianHetHan - now) * time.Second
	if thoiGianConLai < cau_hinh.ThoiGianAnHan {
		
		newExp := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()
		khachHang.CookieExpired = newExp

		// Gọi qua nghiep_vu để đẩy vào hàng chờ
		rowID := nghiep_vu.LayDongKhachHang(khachHang.MaKhachHang)
		if rowID > 0 {
			nghiep_vu.ThemVaoHangCho(
				cau_hinh.BienCauHinh.IdFileSheet,
				"KHACH_HANG",
				rowID,
				mo_hinh.CotKH_CookieExpired,
				newExp,
			)
		}

		maxAge := int(cau_hinh.ThoiGianHetHanCookie.Seconds())
		c.SetCookie("session_id", cookieID, maxAge, "/", "", false, true)
		c.SetCookie("session_sign", cookieSign, maxAge, "/", "", false, true)
	}

	c.Set("USER_ID", khachHang.MaKhachHang)
	c.Set("USER_ROLE", khachHang.VaiTroQuyenHan)
	
	c.Next()
}
