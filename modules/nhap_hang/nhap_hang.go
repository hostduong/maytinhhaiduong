package nhap_hang

import (
	"net/http"

	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangNhapHangMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// 1. Lấy thông tin người dùng hiện tại (Để hiển thị Avatar/Tên trên Header & Sidebar)
	me, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	meCopy := *me

	// Bổ sung Level & Theme cho user để render icon trên giao diện (giống bên thành viên)
	meCopy.StyleLevel = core.LayCapBacVaiTro(shopID, userID, meCopy.VaiTroQuyenHan)
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" {
		meCopy.StyleLevel, meCopy.StyleTheme = 0, 9
	}

	// 2. Lấy dữ liệu Master Data ném ra form Nhập hàng (Lấy từ RAM siêu tốc)
	danhSachNCC := core.LayDanhSachNhaCungCap(shopID)
	danhSachSP := core.LayDanhSachSanPhamMayTinh(shopID)

	// 3. Render giao diện
	c.HTML(http.StatusOK, "master_nhap_hang", gin.H{
		"TieuDe":      "Nhập Hàng", // Tên này phải khớp với menu bên Sidebar để nó sáng màu Tím lên
		"NhanVien":    &meCopy,
		"DanhSachNCC": danhSachNCC,
		"DanhSachSP":  danhSachSP,
	})
}
