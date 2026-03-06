package cua_hang_master

import (
	"net/http"
	"app/config"
	"app/core"
	"github.com/gin-gonic/gin"
)

type CuaHangView struct {
	*core.KhachHang
	GoiDichVuHienTai string
	GoiDichVuBadge   string
	HanSuDung        string
	DomainHienThi    string
	MaxSanPham       int
	MaxNhanVien      int
}

func TrangQuanLyCuaHangMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	userID := c.GetString("USER_ID")
	
	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok || core.LayCapBacVaiTro(masterShopID, userID, c.GetString("USER_ROLE")) > 2 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	
	listAll := core.LayDanhSachKhachHang(adminShopID)
	
	// Nạp bản đồ màu sắc chức vụ (Style Theme) y hệt trang Thành viên
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[adminShopID]
	core.KhoaHeThong.RUnlock()
	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTro { mapStyle[v.MaVaiTro] = v }

	var listView []CuaHangView
	for _, kh := range listAll {
		// [ĐÃ FIX LỖI BẢNG TRẮNG]: Xóa điều kiện lọc vai trò, nạp toàn bộ danh sách

		khCopy := *kh
		if vInfo, ok := mapStyle[khCopy.VaiTroQuyenHan]; ok { 
			khCopy.StyleLevel = vInfo.StyleLevel; khCopy.StyleTheme = vInfo.StyleTheme 
		} else { 
			khCopy.StyleLevel = 9; khCopy.StyleTheme = 0 
		}

		cv := CuaHangView{
			KhachHang: &khCopy,
			GoiDichVuHienTai: "Chưa Đăng Ký",
			GoiDichVuBadge: "bg-slate-100 text-slate-500 border-slate-200",
			HanSuDung: "---",
			DomainHienThi: khCopy.TenDangNhap + ".99k.vn", // Mặc định là Subdomain
			MaxSanPham: 0,
			MaxNhanVien: 0,
		}

		// 💡 ƯU TIÊN CUSTOM DOMAIN
		if khCopy.CauHinh.CustomDomain != "" {
			cv.DomainHienThi = khCopy.CauHinh.CustomDomain
		} else if khCopy.CauHinh.Subdomain != "" {
			cv.DomainHienThi = khCopy.CauHinh.Subdomain
		}

		// Quét gói cước Starter đang dùng
		for _, p := range khCopy.GoiDichVu {
			if p.LoaiGoi == "STARTER" && p.TrangThai == "active" {
				cv.GoiDichVuHienTai = p.TenGoi
				cv.HanSuDung = p.NgayHetHan
				cv.GoiDichVuBadge = "bg-purple-100 text-purple-700 border-purple-200"
				cv.MaxSanPham = p.MaxSanPham
				cv.MaxNhanVien = p.MaxNhanVien
				break
			}
		}

		listView = append(listView, cv)
	}

	meCopy := *me
	if vInfo, ok := mapStyle[meCopy.VaiTroQuyenHan]; ok { meCopy.StyleLevel, meCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme } else { meCopy.StyleLevel, meCopy.StyleTheme = 9, 0 }
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" { meCopy.StyleLevel, meCopy.StyleTheme = 0, 9 }

	c.HTML(http.StatusOK, "cua_hang_master", gin.H{
		"TieuDe": "Quản Lý Cửa Hàng",
		"NhanVien": &meCopy, 
		"DanhSach": listView, 
		"DanhSachVaiTro": listVaiTro, 
	})
}
