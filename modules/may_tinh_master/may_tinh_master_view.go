package may_tinh_master

import (
	"net/http"

	"app/config"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangQuanLyMayTinhMaster(c *gin.Context) {
	defer func() { if err := recover(); err != nil { c.String(500, "LỖI HỆ THỐNG: %v", err) } }()

	masterShopID := c.GetString("SHOP_ID") 
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	kh, found := core.LayKhachHang(masterShopID, userID)
	if !found || kh == nil { c.Redirect(http.StatusFound, "/login"); return }

	if vaiTro != "quan_tri_he_thong" && !core.KiemTraQuyen(masterShopID, vaiTro, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối</h3><a href="/">Về trang chủ</a>`))
		return
	}

	rawList := core.LayDanhSachSanPhamMayTinh(adminShopID)
	
	var cleanList []*core.SanPhamMayTinh 
	var fullList []*core.SanPhamMayTinh  
	groupSP := make(map[string][]*core.SanPhamMayTinh)

	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			fullList = append(fullList, sp)
			groupSP[sp.MaSanPham] = append(groupSP[sp.MaSanPham], sp)
		}
	}

	for _, dsSKU := range groupSP {
		var spChinh *core.SanPhamMayTinh
		for _, sp := range dsSKU { if sp.SKUChinh == 1 { spChinh = sp; break } }
		if spChinh == nil && len(dsSKU) > 0 { spChinh = dsSKU[0] }
		if spChinh != nil { cleanList = append(cleanList, spChinh) }
	}

	c.HTML(http.StatusOK, "may_tinh_master", gin.H{
		"TieuDe":         "Quản lý sản phẩm (Máy Tính)",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"DanhSach":       cleanList, 
		"DanhSachFull":   fullList,  
		"ListDanhMuc":    core.LayDanhSachDanhMuc(adminShopID),    
		"ListThuongHieu": core.LayDanhSachThuongHieu(adminShopID), 
		"ListBLN":        core.LayDanhSachBienLoiNhuan(adminShopID), 
	})
}
