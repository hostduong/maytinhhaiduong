package cua_hang_master

import (
	"net/http"
	"time"

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
	
	core.KhoaHeThong.RLock()
	var listVaiTro []core.VaiTroInfo
	for _, pq := range core.CachePhanQuyen[adminShopID] {
		listVaiTro = append(listVaiTro, core.VaiTroInfo{
			MaVaiTro: pq.MaVaiTro, TenVaiTro: pq.TenVaiTro, StyleLevel: pq.Level, StyleTheme: 5,
		})
	}
	core.KhoaHeThong.RUnlock()

	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTro { mapStyle[v.MaVaiTro] = v }
	mapStyle["quan_tri_he_thong"] = core.VaiTroInfo{StyleLevel: 0, StyleTheme: 9}
	mapStyle["quan_tri_cua_hang"] = core.VaiTroInfo{StyleLevel: 3, StyleTheme: 5}

	var listView []CuaHangView
	for _, kh := range listAll {
		khCopy := *kh

		if khCopy.MaKhachHang == "0000000000000000000" || khCopy.MaKhachHang == "0000000000000000001" || khCopy.VaiTroQuyenHan == "quan_tri_he_thong" {
			khCopy.StyleLevel = 0
			khCopy.StyleTheme = 9
		} else if vInfo, ok := mapStyle[khCopy.VaiTroQuyenHan]; ok { 
			khCopy.StyleLevel = vInfo.StyleLevel
			khCopy.StyleTheme = vInfo.StyleTheme 
		} else { 
			khCopy.StyleLevel = 9
			khCopy.StyleTheme = 0 
		}

		cv := CuaHangView{
			KhachHang: &khCopy,
			GoiDichVuHienTai: "Chưa Đăng Ký",
			GoiDichVuBadge: "bg-slate-100 text-slate-500 border-slate-200",
			HanSuDung: "---",
			DomainHienThi: khCopy.TenDangNhap + ".99k.vn", 
			MaxSanPham: 0,
			MaxNhanVien: 0,
		}

		// [ĐÃ FIX]: Sử dụng khCopy.Domain thay cho khCopy.CauHinh
		if khCopy.Domain.CustomDomain != "" {
			cv.DomainHienThi = khCopy.Domain.CustomDomain
		} else if khCopy.Domain.Subdomain != "" {
			cv.DomainHienThi = khCopy.Domain.Subdomain
		}

		for _, p := range khCopy.GoiDichVu {
			if p.LoaiGoi == "STARTER" && p.TrangThai == "active" {
				cv.GoiDichVuHienTai = p.TenGoi
				// [ĐÃ FIX]: Convert Unix Timestamp sang String
				cv.HanSuDung = time.Unix(p.NgayHetHan, 0).Format("02/01/2006")
				cv.GoiDichVuBadge = "bg-purple-100 text-purple-700 border-purple-200"
				cv.MaxSanPham = p.MaxSanPham
				cv.MaxNhanVien = p.MaxNhanVien
				break
			}
		}

		listView = append(listView, cv)
	}

	meCopy := *me
	if vInfo, ok := mapStyle[meCopy.VaiTroQuyenHan]; ok { meCopy.StyleLevel, meCopy.StyleTheme = vInfo.StyleLevel, meCopy.StyleTheme } else { meCopy.StyleLevel, meCopy.StyleTheme = 9, 0 }
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" { meCopy.StyleLevel, meCopy.StyleTheme = 0, 9 }

	c.HTML(http.StatusOK, "cua_hang_master", gin.H{
		"TieuDe": "Quản Lý Cửa Hàng",
		"NhanVien": &meCopy, 
		"DanhSach": listView, 
		"DanhSachVaiTro": listVaiTro, 
	})
}
