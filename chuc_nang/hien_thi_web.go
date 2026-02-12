package chuc_nang

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

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
	sp, tonTai := core.LayChiTietSanPham(shopID, id)
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
	c.HTML(http.StatusOK, "ho_so", gin.H{
		"TieuDe": "Hồ sơ cá nhân", "DaDangNhap": daLogin,
		"TenNguoiDung": tenUser, "QuyenHan": quyen, "NhanVien": kh,
	})
}
