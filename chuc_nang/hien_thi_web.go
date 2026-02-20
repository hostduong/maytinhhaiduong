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

// =============================================================
// [MỚI] BỘ HÀM TIỆN ÍCH CHO HTML (VIEW HELPER)
// =============================================================
func LayBoHamHTML() template.FuncMap {
	return template.FuncMap{
		// 1. Hàm cắt chuỗi ảnh (Giải quyết vấn đề của bạn)
		"firstImg": func(s string) string {
			if s == "" { return "" }
			// Cắt lấy ảnh đầu tiên trước dấu |
			parts := strings.Split(s, "|")
			url := strings.TrimSpace(parts[0])
			return url
		},

		// 2. Định dạng tiền tệ (Mang từ main qua)
		"format_money": func(n float64) string {
			p := message.NewPrinter(language.Vietnamese)
			return p.Sprintf("%.0f", n)
		},

		// 3. Xuất JSON cho JS dùng (Mang từ main qua)
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
	danhSachSP := core.LayDanhSachSanPham(shopID)
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe": "Trang Chủ", "DanhSachSanPham": danhSachSP,
		"DaDangNhap": daLogin, "TenNguoiDung": tenUser, "QuyenHan": quyen,
	})
}

func ChiTietSanPham(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	id := c.Param("id")
	sp, tonTai := core.LayChiTietSKU(shopID, id)
	if !tonTai { c.String(http.StatusNotFound, "Không tìm thấy!"); return }
	daLogin, tenUser, quyen := layThongTinNguoiDung(c)
	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
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

	// --- [MỚI] RẼ NHÁNH TEMPLATE THEO QUYỀN ---
	templateName := "ho_so"         // Mặc định là cho Khách hàng
	if quyen != "customer" {
		templateName = "ho_so_admin" // Nếu là nhân sự -> Nạp file của Admin
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"TieuDe": "Hồ sơ cá nhân", "DaDangNhap": daLogin,
		"TenNguoiDung": tenUser, "QuyenHan": quyen, "NhanVien": kh,
	})
}
