package thanh_vien_master

import (
	"net/http"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangQuanLyThanhVienMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	myRole := c.GetString("USER_ROLE")
	
	// Lấy cấp bậc (Level) của người đang đăng nhập
	myLevel := Repo_LayCapBac(masterShopID, userID, myRole)

	if myLevel > 2 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	
	me, ok := Repo_LayKhachHangMaster(masterShopID, userID)
	if !ok || me == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	listAll := Repo_LayDanhSachMaster(masterShopID)
	
	core.KhoaHeThong.RLock()
	listVaiTroRaw := core.CacheDanhSachVaiTro[masterShopID]
	core.KhoaHeThong.RUnlock()

	if len(listVaiTroRaw) == 0 {
		listVaiTroRaw = []core.VaiTroInfo{
			{MaVaiTro: "quan_tri_he_thong", TenVaiTro: "Sáng Lập Viên", StyleLevel: 0, StyleTheme: 9},
			{MaVaiTro: "quan_tri_vien_he_thong", TenVaiTro: "Quản trị viên hệ thống", StyleLevel: 1, StyleTheme: 4},
			{MaVaiTro: "quan_tri_it_he_thong", TenVaiTro: "Quản trị IT hệ thống", StyleLevel: 2, StyleTheme: 7},
		}
	}

	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTroRaw { mapStyle[v.MaVaiTro] = v }

	// ========================================================
	// ĐỘNG CƠ LỌC QUYỀN TỪ SERVER (BẢO MẬT BACKEND)
	// ========================================================
	var listVaiTroFiltered []core.VaiTroInfo
	for _, v := range listVaiTroRaw {
		if myLevel == 0 {
			// Sáng Lập Viên (Level 0) thấy được mọi quyền
			listVaiTroFiltered = append(listVaiTroFiltered, v)
		} else if v.StyleLevel > myLevel {
			// Các quản lý khác chỉ thấy được quyền THẤP HƠN mình (Level số to hơn)
			listVaiTroFiltered = append(listVaiTroFiltered, v)
		}
	}

	var listView []*core.KhachHang
	for _, kh := range listAll {
		if kh == nil { continue }
		khCopy := *kh 
		khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
		
		if khCopy.MangXaHoi == nil { khCopy.MangXaHoi = make(map[string]string) }

		if khCopy.MaKhachHang == "0000000000000000000" || khCopy.MaKhachHang == "0000000000000000001" || khCopy.VaiTroQuyenHan == "quan_tri_he_thong" {
			khCopy.StyleLevel, khCopy.StyleTheme = 0, 9 
		} else {
			if vInfo, ok := mapStyle[khCopy.VaiTroQuyenHan]; ok { 
				khCopy.StyleLevel, khCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme 
			} else { 
				khCopy.StyleLevel, khCopy.StyleTheme = 9, 0 
			}
		}
		listView = append(listView, &khCopy)
	}

	meCopy := *me
	if meCopy.MangXaHoi == nil { meCopy.MangXaHoi = make(map[string]string) }
	
	if vInfo, ok := mapStyle[meCopy.VaiTroQuyenHan]; ok { 
		meCopy.StyleLevel, meCopy.StyleTheme = vInfo.StyleLevel, meCopy.StyleTheme 
	} else { 
		meCopy.StyleLevel, meCopy.StyleTheme = 9, 0 
	}
	
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" { 
		meCopy.StyleLevel, meCopy.StyleTheme = 0, 9 
	}

	c.HTML(http.StatusOK, "thanh_vien_master", gin.H{
		"TieuDe": "Thành Viên",
		"NhanVien": &meCopy, 
		"DanhSach": listView, 
		"DanhSachVaiTro": listVaiTroFiltered, // CHỈ BƠM DANH SÁCH ĐÃ LỌC
	})
}
