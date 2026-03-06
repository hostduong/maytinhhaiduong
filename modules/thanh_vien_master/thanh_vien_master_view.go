package thanh_vien_master

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangQuanLyThanhVienMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	if Repo_LayCapBac(masterShopID, userID, c.GetString("USER_ROLE")) > 2 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	
	me, _ := Repo_LayKhachHang(masterShopID, userID)
	listAll := Repo_LayDanhSach(masterShopID)
	
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[masterShopID]
	core.KhoaHeThong.RUnlock()

	if len(listVaiTro) == 0 {
		listVaiTro = []core.VaiTroInfo{
			{MaVaiTro: "quan_tri_he_thong", TenVaiTro: "Quản trị hệ thống", StyleLevel: 0, StyleTheme: 9},
			{MaVaiTro: "quan_tri_vien_he_thong", TenVaiTro: "Quản trị viên hệ thống", StyleLevel: 1, StyleTheme: 4},
			{MaVaiTro: "quan_tri_it_he_thong", TenVaiTro: "Quản trị IT hệ thống", StyleLevel: 2, StyleTheme: 7},
			{MaVaiTro: "quan_tri_cua_hang", TenVaiTro: "Quản trị cửa hàng", StyleLevel: 3, StyleTheme: 5},
		}
	}

	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTro { mapStyle[v.MaVaiTro] = v }

	var listView []*core.KhachHang
	for _, kh := range listAll {
		khCopy := *kh 
		khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
		if khCopy.MaKhachHang == "0000000000000000000" || khCopy.MaKhachHang == "0000000000000000001" || khCopy.VaiTroQuyenHan == "quan_tri_he_thong" {
			khCopy.StyleLevel, khCopy.StyleTheme = 0, 9 
		} else {
			if vInfo, ok := mapStyle[khCopy.VaiTroQuyenHan]; ok { khCopy.StyleLevel, khCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme } else { khCopy.StyleLevel, khCopy.StyleTheme = 9, 0 }
		}
		listView = append(listView, &khCopy)
	}

	meCopy := *me
	if vInfo, ok := mapStyle[meCopy.VaiTroQuyenHan]; ok { meCopy.StyleLevel, meCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme } else { meCopy.StyleLevel, meCopy.StyleTheme = 9, 0 }
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" { meCopy.StyleLevel, meCopy.StyleTheme = 0, 9 }

	c.HTML(http.StatusOK, "thanh_vien_master", gin.H{
		"TieuDe": "Thành Viên",
		"NhanVien": &meCopy, 
		"DanhSach": listView, 
		"DanhSachVaiTro": listVaiTro, 
	})
}
