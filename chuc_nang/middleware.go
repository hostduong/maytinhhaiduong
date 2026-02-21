package chuc_nang

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// =============================================================
// PHẦN 0: GATEWAY SAAS & TRẠM KIỂM SOÁT DỊCH VỤ (CHẶNG 2)
// =============================================================

// 1. GatewaySaaS: Điều phối Subdomain và Khởi tạo Shop
func GatewaySaaS(c *gin.Context) {
	host := c.Request.Host
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0] // Bỏ port nếu chạy local (vd: localhost:8080)
	}

	masterShopID := cau_hinh.BienCauHinh.IdFileSheet // ID của Nền tảng www.99k.vn

	// --- TRƯỜNG HỢP 1: TẦNG 0 (TRANG CHỦ NỀN TẢNG) ---
	if host == "www.99k.vn" || host == "99k.vn" || host == "localhost" {
		c.Set("SHOP_ID", masterShopID)
		c.Set("THEME", "theme_master")
		c.Next()
		return
	}

	// --- TRƯỜNG HỢP 2: TẦNG 1 & TẦNG 3 (CỬA HÀNG) ---
	subdomain := strings.Split(host, ".")[0] // Lấy "cuahang1" từ "cuahang1.99k.vn"
	
	// Quét RAM của Master để tìm Chủ Shop
	danhSachChung := core.LayDanhSachKhachHang(masterShopID)
	var tenant *core.KhachHang

	for _, kh := range danhSachChung {
		// Tìm theo Subdomain (Tên đăng nhập) hoặc Domain riêng
		if strings.ToLower(kh.TenDangNhap) == subdomain || kh.CauHinh.CustomDomain == host {
			tenant = kh
			break
		}
	}

	// Nếu không tìm thấy Chủ Shop nào khớp
	if tenant == nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte(`
			<div style="text-align:center; padding: 50px; font-family: sans-serif;">
				<h1 style="color:#ef4444;">Cửa hàng không tồn tại</h1>
				<p>Địa chỉ trang web này không thuộc hệ thống hoặc đã bị xóa.</p>
				<p>Truy cập <a href="https://www.99k.vn" style="color:#3b82f6;">99K.vn</a> để tạo cửa hàng mới.</p>
			</div>
		`))
		c.Abort()
		return
	}

	// Lấy ID Sheet của cửa hàng đó
	shopID := tenant.DataSheets.SpreadsheetID
	if shopID == "" {
		c.Data(http.StatusServiceUnavailable, "text/html; charset=utf-8", []byte(`
			<div style="text-align:center; padding: 50px; font-family: sans-serif;">
				<h1 style="color:#f59e0b;">Đang khởi tạo Dữ liệu</h1>
				<p>Hệ thống đang chuẩn bị Database cho cửa hàng này. Vui lòng thử lại sau vài phút.</p>
			</div>
		`))
		c.Abort()
		return
	}

	// Lấy Theme
	theme := tenant.CauHinh.Theme
	if theme == "" { theme = "may_tinh" } // Default Theme

	// Đẩy thông tin vào luồng để các Middleware và Controller sau xài
	c.Set("SHOP_ID", shopID)
	c.Set("THEME", theme)
	c.Set("TENANT_INFO", tenant)

	c.Next()
}

// 2. KiemTraGoiDichVu: Trạm thu phí tự động
func KiemTraGoiDichVu(c *gin.Context) {
	tenantVal, exists := c.Get("TENANT_INFO")
	if !exists {
		// Không có tenant info -> Đang ở Nền tảng mẹ (Tầng 0) -> Miễn phí qua trạm
		c.Next()
		return
	}

	tenant := tenantVal.(*core.KhachHang)
	hasActivePlan := false
	now := time.Now()

	// Quét mảng Gói dịch vụ của Chủ shop
	for _, plan := range tenant.GoiDichVu {
		if plan.TrangThai == "active" || plan.TrangThai == "trial" {
			// Nếu không có ngày hết hạn (Gói vĩnh viễn)
			if plan.NgayHetHan == "" {
				hasActivePlan = true
				break
			}
			// Parse ngày và so sánh
			expDate, err := time.Parse("2006-01-02", plan.NgayHetHan)
			if err == nil && (expDate.After(now) || expDate.Equal(now)) {
				hasActivePlan = true
				break
			}
		}
	}

	// Nếu hết hạn -> Khóa chặn cửa
	if !hasActivePlan {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{"status": "error", "msg": "Cửa hàng đã hết hạn dịch vụ. Vui lòng gia hạn."})
		} else {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
				<div style="text-align:center; padding: 50px; font-family: sans-serif; background: #f8fafc; height: 100vh;">
					<div style="max-width: 500px; margin: 0 auto; background: white; padding: 40px; border-radius: 16px; box-shadow: 0 10px 25px rgba(0,0,0,0.05);">
						<div style="font-size: 48px; margin-bottom: 20px;">🚧</div>
						<h1 style="color:#334155; margin-bottom: 10px;">Cửa Hàng Tạm Ngưng</h1>
						<p style="color:#64748b; line-height: 1.6;">Cửa hàng này đang tạm ngưng hoạt động do hết hạn gói dịch vụ.</p>
						<p style="color:#64748b; line-height: 1.6;">Nếu bạn là chủ cửa hàng, vui lòng đăng nhập vào hệ thống quản trị mẹ để gia hạn.</p>
						<a href="https://www.99k.vn/login" style="display:inline-block; margin-top:20px; padding: 12px 24px; background: #2563eb; color: white; text-decoration: none; font-weight: bold; border-radius: 8px;">Quản lý thanh toán</a>
					</div>
				</div>
			`))
		}
		c.Abort()
		return
	}

	c.Next()
}


// =============================================================
// PHẦN 1: RATE LIMIT (BẢO VỆ CHỐNG SPAM)
// =============================================================
var boDem = make(map[string]int)
var mtx sync.Mutex

func KhoiTaoBoDemRateLimit() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mtx.Lock()
			boDem = make(map[string]int) 
			mtx.Unlock()
		}
	}()
}

func xoaCookie(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
}

// =============================================================
// PHẦN 2: MIDDLEWARE XÁC THỰC (AUTH)
// =============================================================
func KiemTraDangNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")

	if core.HeThongDangBan && c.Request.Method != "GET" {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"status": "error", "msg": "Hệ thống đang đồng bộ, vui lòng thử lại sau 5 giây."})
		return
	}

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

	if err1 != nil || cookieID == "" {
		if c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" || c.Request.Method == "POST" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Vui lòng đăng nhập!"})
		} else {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		}
		return
	}

	userAgent := c.Request.UserAgent()
	signatureServer := cau_hinh.TaoChuKyBaoMat(cookieID, userAgent) 

	if err2 != nil || cookieSign != signatureServer {
		xoaCookie(c)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Phát hiện bất thường (Cookie Mismatch)!"})
		return
	}

	khachHang, timThay := core.TimKhachHangTheoCookie(shopID, cookieID)
	if !timThay {
		xoaCookie(c)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	tokenInfo, ok := khachHang.RefreshTokens[cookieID]
	if !ok {
		xoaCookie(c)
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	now := time.Now().Unix()
	if now > tokenInfo.ExpiresAt {
		xoaCookie(c)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
		return
	}

	// Gắn thông tin User vào Context
	c.Set("USER_ID", khachHang.MaKhachHang)
	c.Set("USER_ROLE", khachHang.VaiTroQuyenHan)
	c.Set("USER_NAME", khachHang.TenKhachHang)
	
	c.Next()
}

// =============================================================
// PHẦN 3: PHÂN QUYỀN (ADMIN GATEKEEPER)
// =============================================================
func KiemTraQuyenHan(c *gin.Context) {
	role := c.GetString("USER_ROLE")

	if role == "" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	if role == "khach_hang" || role == "customer" {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}

	c.Next()
}
