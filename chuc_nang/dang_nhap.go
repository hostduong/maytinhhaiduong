package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/cau_hinh" // [SỬA] Đã gộp bao_mat vào cau_hinh
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangDangNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS] Lấy ShopID
	theme := c.GetString("THEME") // [SAAS] Lấy theme động
	cookie, _ := c.Cookie("session_id")
	
	if cookie != "" {
		// [SAAS] Tìm theo ShopID
		if _, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, theme+"dang_nhap", gin.H{"TieuDe": "Đăng Nhập"})
}

func XuLyDangNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID") // [SAAS]
	theme := c.GetString("THEME") // [SAAS] Lấy theme động

	inputDinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))
	ghiNho        := c.PostForm("ghi_nho")

	// [SAAS] Tìm trong Shop hiện tại
	kh, ok := core.TimKhachHangTheoUserOrEmail(shopID, inputDinhDanh)
	if !ok {
		c.HTML(http.StatusOK, theme+"dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	// [SỬA] Dùng hàm từ cau_hinh
	if !cau_hinh.KiemTraMatKhau(pass, kh.MatKhauHash) {
		c.HTML(http.StatusOK, theme+"dang_nhap", gin.H{"Loi": "Mật khẩu không đúng!"})
		return
	}

	if kh.TrangThai == 0 {
		c.HTML(http.StatusOK, theme+"dang_nhap", gin.H{"Loi": "Tài khoản bị khóa!"})
		return
	}

	// --- LOGIC TẠO SESSION ---
	var thoiGianSong time.Duration
	if ghiNho == "on" {
		thoiGianSong = 30 * 24 * time.Hour
	} else {
		thoiGianSong = cau_hinh.ThoiGianHetHanCookie
	}

	sessionID := cau_hinh.TaoSessionIDAnToan()
	userAgent := c.Request.UserAgent()
	signature := cau_hinh.TaoChuKyBaoMat(sessionID, userAgent)
	
	expTime := time.Now().Add(thoiGianSong).Unix()
	maxAge  := int(thoiGianSong.Seconds())

	// --- [QUAN TRỌNG] CẬP NHẬT MAP TOKEN (JSON) ---
	core.KhoaHeThong.Lock()
	if kh.RefreshTokens == nil {
		kh.RefreshTokens = make(map[string]core.TokenInfo)
	}
	
	// Thêm thiết bị mới vào danh sách
	kh.RefreshTokens[sessionID] = core.TokenInfo{
		DeviceName: userAgent,
		ExpiresAt:  expTime,
	}
	core.KhoaHeThong.Unlock()
	
	// --- GHI XUỐNG SHEET (CỘT JSON F) ---
	go func() {
		sID := kh.SpreadsheetID
		if sID == "" { sID = shopID } // Fallback
		
		// Serialize Map thành JSON String
		jsonStr := core.ToJSON(kh.RefreshTokens)
		
		// Ghi vào cột CotKH_RefreshTokenJson (Cột F)
		core.ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_RefreshTokenJson, jsonStr)
	}()

	// Set Cookie trình duyệt
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)
	c.SetCookie("session_sign", signature, maxAge, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func DangXuat(c *gin.Context) {
	// Xóa cookie trình duyệt
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.SetCookie("session_sign", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
