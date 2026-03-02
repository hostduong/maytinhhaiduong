package hien_thi_web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"app/core"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Giữ nguyên bộ hàm Format HTML chuẩn của bạn
func LayBoHamHTML() template.FuncMap {
	return template.FuncMap{
		"firstImg": func(s string) string {
			if s == "" { return "" }
			parts := strings.Split(s, "|")
			return strings.TrimSpace(parts[0])
		},
		"format_money": func(n float64) string {
			p := message.NewPrinter(language.Vietnamese)
			return p.Sprintf("%.0f", n)
		},
		"json": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
		"split": strings.Split,
	}
}

// Phục hồi nguyên vẹn hàm đọc Cookie và lấy thông tin User
func layThongTinNguoiDung(c *gin.Context) (bool, string, string) {
	shopID := c.GetString("SHOP_ID")
	cookie, _ := c.Cookie("session_token") // Đã đồng bộ key cookie với m_auth
	if cookie != "" {
		if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			return true, kh.TenKhachHang, kh.VaiTroQuyenHan
		}
	}
	return false, "", ""
}

// Phục hồi nguyên vẹn logic điều hướng Trang chủ đa luồng (B2B & B2C)
func TrangChu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") 
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	// 1. NGÃ RẼ 1: DÀNH CHO TRANG CHỦ NỀN TẢNG (99k.vn)
	if theme == "theme_master" || theme == "" {
		c.HTML(http.StatusOK, "trang_chu", gin.H{
			"TieuDe": "Nền tảng tạo Website & POS chỉ với 99K",
			"DaDangNhap": daLogin, 
			"TenNguoiDung": tenUser, 
			"QuyenHan": quyen,
		})
		return
	}

	// 2. NGÃ RẼ 2: DÀNH CHO CỬA HÀNG B2C (VD: cuahang.99k.vn)
	danhSachSP := core.LayDanhSachSanPhamMayTinh(shopID) 
	
	tenantVal, exists := c.Get("TENANT_INFO")
	var cauHinh core.UserConfig
	if exists {
		chuShop := tenantVal.(*core.KhachHang)
		cauHinh = chuShop.CauHinh
	}

	c.HTML(http.StatusOK, theme+"/trang_chu", gin.H{
		"TieuDe": "Trang Chủ", 
		"DanhSachSanPham": danhSachSP,
		"DaDangNhap": daLogin, 
		"TenNguoiDung": tenUser, 
		"QuyenHan": quyen,
		"CauHinhShop": cauHinh,
	})
}

// Phục hồi hàm chi tiết sản phẩm 
func ChiTietSanPham(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") 
	id := c.Param("id")
	
	sp, tonTai := core.LayChiTietSKUMayTinh(shopID, id)
	if !tonTai { c.String(http.StatusNotFound, "Không tìm thấy!"); return }
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	tenantVal, exists := c.Get("TENANT_INFO")
	var cauHinh core.UserConfig
	if exists {
		chuShop := tenantVal.(*core.KhachHang)
		cauHinh = chuShop.CauHinh
	}

	c.HTML(http.StatusOK, theme+"/chi_tiet_san_pham", gin.H{
		"TieuDe": sp.TenSanPham, "SanPham": sp,
		"DaDangNhap": daLogin, "TenNguoiDung": tenUser, "QuyenHan": quyen,
		"CauHinhShop": cauHinh,
	})
}

// ========================================================
// CÁC HÀM RENDER GIAO DIỆN BẢO MẬT (MẶT TIỀN PUBLIC)
// ========================================================

func TrangDangNhap(c *gin.Context) {
	c.HTML(http.StatusOK, "dang_nhap", gin.H{"TieuDe": "Đăng Nhập Hệ Thống"})
}

func TrangDangKy(c *gin.Context) {
	c.HTML(http.StatusOK, "dang_ky", gin.H{"TieuDe": "Đăng Ký Tài Khoản Mới"})
}

func TrangQuenMatKhau(c *gin.Context) {
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{"TieuDe": "Khôi phục Mật Khẩu"})
}

func TrangXacThucOTP(c *gin.Context) {
	c.HTML(http.StatusOK, "xac_thuc_otp", gin.H{"TieuDe": "Xác thực OTP"})
}
