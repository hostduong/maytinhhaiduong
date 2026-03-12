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
	cookie, _ := c.Cookie("session_token") 
	if cookie != "" {
		if kh, ok := core.TimKhachHangTheoCookie(shopID, cookie); ok {
			return true, kh.ThongTin.TenKhachHang, kh.VaiTroQuyenHan
		}
	}
	return false, "", ""
}

func TrangChu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") 
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	if theme == "theme_master" || theme == "" {
		c.HTML(http.StatusOK, "trang_chu", gin.H{
			"TieuDe": "Nền tảng tạo Website & POS chỉ với 99K",
			"DaDangNhap": daLogin, 
			"TenNguoiDung": tenUser, 
			"QuyenHan": quyen,
		})
		return
	}

	var danhSachSP []*core.ProductJSON
	core.KhoaHeThong.RLock()
	if nganhMap, exists := core.CacheSanPham[shopID]; exists {
		for _, ds := range nganhMap {
			danhSachSP = append(danhSachSP, ds...)
		}
	}
	core.KhoaHeThong.RUnlock()
	
	tenantVal, exists := c.Get("TENANT_INFO")
	var cauHinh core.TenantCauHinh
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
	theme := c.GetString("THEME") 
	id := c.Param("id")
	
	var sp *core.ProductJSON
	tonTai := false

	core.KhoaHeThong.RLock()
	if foundSp, ok := core.CacheMapSanPham[core.TaoCompositeKey(shopID, id)]; ok {
		sp = foundSp
		tonTai = true
	} else {
		for _, nganhMap := range core.CacheSanPham[shopID] {
			for _, item := range nganhMap {
				if item.Slug == id && item.TrangThai == 1 {
					sp = item
					tonTai = true
					break
				}
			}
			if tonTai { break }
		}
	}
	core.KhoaHeThong.RUnlock()

	if !tonTai || sp == nil { c.String(http.StatusNotFound, "Không tìm thấy!"); return }
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	tenantVal, exists := c.Get("TENANT_INFO")
	var cauHinh core.TenantCauHinh
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
