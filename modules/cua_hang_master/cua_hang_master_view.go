package cua_hang_master

import (
	"net/http"
	"app/config"
	"app/core"
	"github.com/gin-gonic/gin"
)

// Bổ sung Struct để đổ dữ liệu ra View dễ dàng hơn
type CuaHangView struct {
	*core.KhachHang
	GoiDichVuHienTai string
	GoiDichVuBadge   string
	HanSuDung        string
	DomainHienThi    string
}

func TrangQuanLyCuaHangMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") // Đây là Master
	adminShopID := config.BienCauHinh.IdFileSheetAdmin // Đây là nơi chứa danh sách cửa hàng
	userID := c.GetString("USER_ID")
	
	// Lấy thông tin Sếp
	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok || core.LayCapBacVaiTro(masterShopID, userID, c.GetString("USER_ROLE")) > 2 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	
	// Lấy danh sách toàn bộ Cửa Hàng (Từ Tầng Admin)
	listAll := core.LayDanhSachKhachHang(adminShopID)

	var listView []CuaHangView
	for _, kh := range listAll {

		cv := CuaHangView{
			KhachHang: kh,
			GoiDichVuHienTai: "Chưa Đăng Ký",
			GoiDichVuBadge: "bg-slate-100 text-slate-500",
			HanSuDung: "---",
			DomainHienThi: kh.TenDangNhap + ".99k.vn",
		}

		// Xử lý Domain hiển thị
		if kh.CauHinh.CustomDomain != "" {
			cv.DomainHienThi = kh.CauHinh.CustomDomain
		}

		// Xử lý Gói Dịch Vụ (Tìm gói Starter đang Active)
		for _, p := range kh.GoiDichVu {
			if p.LoaiGoi == "STARTER" && p.TrangThai == "active" {
				cv.GoiDichVuHienTai = p.TenGoi
				cv.HanSuDung = p.NgayHetHan
				cv.GoiDichVuBadge = "bg-purple-100 text-purple-700"
				break
			}
		}

		listView = append(listView, cv)
	}

	c.HTML(http.StatusOK, "cua_hang_master", gin.H{
		"TieuDe": "Quản Lý Cửa Hàng",
		"NhanVien": me, 
		"DanhSach": listView, 
	})
}
