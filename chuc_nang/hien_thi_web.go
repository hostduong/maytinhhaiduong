package chuc_nang

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"app/core"
	data_pc "app/core/may_tinh"
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
	theme := c.GetString("THEME") // [SAAS] Lấy theme động
	
	// Tạm thời vẫn lấy danh sách chung, Chặng 4 ta sẽ đổi hàm này thành LayDanhSachSanPham_PC
	danhSachSP := data_pc.LayDanhSachSanPham(shopID) 
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	c.HTML(http.StatusOK, "/khung_giao_dien", gin.H{
		"TieuDe": "Trang Chủ", "DanhSachSanPham": danhSachSP,
		"DaDangNhap": daLogin, "TenNguoiDung": tenUser, "QuyenHan": quyen,
	})
}

func ChiTietSanPham(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	theme := c.GetString("THEME") // [SAAS] Lấy theme động
	id := c.Param("id")
	
	sp, tonTai := data_pc.LayChiTietSKU(shopID, id)
	if !tonTai { c.String(http.StatusNotFound, "Không tìm thấy!"); return }
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	
	c.HTML(http.StatusOK, "/chi_tiet_san_pham", gin.H{
		"TieuDe": sp.TenSanPham, "SanPham": sp,
		"DaDangNhap": daLogin, "TenNguoiDung": tenUser, "QuyenHan": quyen,
	})
}

func TrangHoSo(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	if !daLogin { c.Redirect(http.StatusFound, "/login"); return }
	
	cookie, _ := c.Cookie("session_id")
	kh, _ := core.TimKhachHangTheoCookie(shopID, cookie)

	templateName := theme + "/ho_so" // Giao diện khách
	if quyen != "customer" {
		templateName = "ho_so_admin" // Giao diện Admin (Dùng chung)
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"TieuDe": "Hồ sơ cá nhân", "DaDangNhap": daLogin,
		"TenNguoiDung": tenUser, "QuyenHan": quyen, "NhanVien": kh,
	})
}
