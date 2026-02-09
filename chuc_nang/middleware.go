package chuc_nang

import (
	"net/http"
	"sync"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/core" // [MỚI] Sử dụng Core

	"github.com/gin-gonic/gin"
)

// =============================================================
// PHẦN 1: RATE LIMIT (GIỮ NGUYÊN LOGIC CŨ)
// =============================================================
var boDem = make(map[string]int)
var mtx sync.Mutex

func KhoiTaoBoDemRateLimit() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mtx.Lock()
			boDem = make(map[string]int) // Reset bộ đếm mỗi giây
			mtx.Unlock()
		}
	}()
}

// Helper xóa cookie
func xoaCookie(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
}

// =============================================================
// PHẦN 2: MIDDLEWARE XÁC THỰC (AUTH & RENEW)
// =============================================================

// KiemTraDangNhap: Dùng cho API User & các trang cần đăng nhập
func KiemTraDangNhap(c *gin.Context) {
	// 1. CHỐT CHẶN BẢO TRÌ (Dùng biến từ Core)
	if core.HeThongDangBan && c.Request.Method != "GET" {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
			"status": "error", "msg": "Hệ thống đang đồng bộ dữ liệu, vui lòng thử lại sau 5 giây.",
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
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"status": "error", "msg": "Thao tác quá nhanh! Vui lòng chậm lại."})
		return
	}

	// 3. KIỂM TRA COOKIE CƠ BẢN
	if err1 != nil || cookieID == "" {
		// Nếu là API (AJAX) -> Trả lỗi JSON
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" || c.Request.Method == "POST" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Vui lòng đăng nhập!"})
		} else {
			c.Redirect(http.StatusFound, "/login") // Nếu là trình duyệt -> Chuyển hướng
			c.Abort()
		}
		return
	}

	// 4. KIỂM TRA CHỮ KÝ BẢO MẬT (SECURITY CHECK)
	userAgent := c.Request.UserAgent()
	signatureServer := bao_mat.TaoChuKyBaoMat(cookieID, userAgent)

	if err2 != nil || cookieSign != signatureServer {
		xoaCookie(c)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Phát hiện bất thường (Cookie Mismatch)!"})
		return
	}

	// 5. TÌM USER TRONG CORE (Thay thế nghiep_vu cũ)
	khachHang, timThay := core.TimKhachHangTheoCookie(cookieID)

	if !timThay {
		xoaCookie(c)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// 6. LOGIC GIA HẠN THÔNG MINH (Auto-Renew)
	thoiGianHetHan := khachHang.CookieExpired
	now := time.Now().Unix()

	// Nếu đã hết hạn -> Đá ra
	if now > thoiGianHetHan {
		xoaCookie(c)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
		return
	}

	// Nếu còn hạn nhưng sắp hết (trong vùng ân hạn) -> GIA HẠN
	thoiGianConLai := time.Duration(thoiGianHetHan - now) * time.Second
	if thoiGianConLai < cau_hinh.ThoiGianAnHan {
		
		newExp := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()
		
		// Cập nhật RAM Core an toàn
		core.KhoaHeThong.Lock()
		khachHang.CookieExpired = newExp
		core.KhoaHeThong.Unlock()

		// Đẩy vào hàng chờ (Sử dụng hàm Core mới)
		sID := khachHang.SpreadsheetID
		if sID == "" { sID = cau_hinh.BienCauHinh.IdFileSheet }
		
		// Ghi đè cột Cookie Expired (Cột E / Index 4)
		// Lưu ý: Core dùng const CotKH_... nên rất chuẩn
		core.ThemVaoHangCho(
			sID,
			"KHACH_HANG",
			khachHang.DongTrongSheet,
			core.CotKH_CookieExpired, 
			newExp,
		)

		// Set lại Cookie mới cho trình duyệt
		maxAge := int(cau_hinh.ThoiGianHetHanCookie.Seconds())
		c.SetCookie("session_id", cookieID, maxAge, "/", "", false, true)
		c.SetCookie("session_sign", cookieSign, maxAge, "/", "", false, true)
	}

	// 7. LƯU THÔNG TIN VÀO CONTEXT (Để Controller dùng)
	c.Set("USER_ID", khachHang.MaKhachHang)
	c.Set("USER_ROLE", khachHang.VaiTroQuyenHan)
	c.Set("USER_NAME", khachHang.TenKhachHang)
	
	c.Next()
}


// =============================================================
// PHẦN 3: MIDDLEWARE PHÂN QUYỀN (ADMIN GATEKEEPER)
// =============================================================


func KiemTraQuyenHan(c *gin.Context) {
	role := c.GetString("USER_ROLE")

	if role == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	// [SỬA LẠI] Chặn cả "khach_hang" (theo PDF) và "customer" (theo code đăng ký)
	if role == "khach_hang" || role == "customer" {
		c.HTML(http.StatusForbidden, "khung_giao_dien", gin.H{ // Dùng khung_giao_dien cho nhẹ, tránh lỗi
			"TieuDe": "Không có quyền truy cập",
			"DaDangNhap": true,
			"NoiDung": "<h3>⛔ KHÔNG CÓ QUYỀN TRUY CẬP</h3><p>Tài khoản khách hàng không thể vào trang quản trị.</p>",
		})
		c.Abort()
		return
	}

	c.Next()
}
