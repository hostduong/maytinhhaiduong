package chuc_nang

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

func layThongTinNguoiDung(c *gin.Context) (bool, string, string) {
	shopID := c.GetString("SHOP_ID")
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			return true, kh.TenKhachHang, kh.VaiTroQuyenHan
		}
	}
	return false, "", ""
}

func TrangChu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") 
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	// 1. NGÃ RẼ 1: DÀNH CHO TRANG CHỦ NỀN TẢNG (99k.vn)
	if theme == "theme_master" {
		c.HTML(http.StatusOK, "theme_master/trang_chu", gin.H{
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

func ChiTietSanPham(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") // [SAAS] Lấy theme động
	id := c.Param("id")
	
	sp, tonTai := core.LayChiTietSKUMayTinh(shopID, id)
	if !tonTai { c.String(http.StatusNotFound, "Không tìm thấy!"); return }
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	tenantVal, _ := c.Get("TENANT_INFO")
	chuShop := tenantVal.(*core.KhachHang)

	c.HTML(http.StatusOK, theme+"/chi_tiet_san_pham", gin.H{
		"TieuDe": sp.TenSanPham, "SanPham": sp,
		"DaDangNhap": daLogin, "TenNguoiDung": tenUser, "QuyenHan": quyen,
		"CauHinhShop": chuShop.CauHinh,
	})
}

func TrangHoSo(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") // Khai báo lại để dùng
	
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	if !daLogin { c.Redirect(http.StatusFound, "/login"); return }
	
	cookie, _ := c.Cookie("session_id")
	kh, _ := core.TimKhachHangTheoCookie(shopID, cookie)

	tenantVal, _ := c.Get("TENANT_INFO")
	chuShop := tenantVal.(*core.KhachHang)

	templateName := "ho_so" // Gọi form chung
	if quyen != "customer" {
		templateName = "ho_so_admin" 
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"TieuDe": "Hồ sơ cá nhân", "DaDangNhap": daLogin,
		"TenNguoiDung": tenUser, "QuyenHan": quyen, "NhanVien": kh,
		"CauHinhShop": chuShop.CauHinh, "Theme": theme, // Bơm theme vào để HTML tự điều hướng (nếu cần)
	})
}
