package chuc_nang_master

import (
	"net/http"
	"strconv"
	"strings"

	"app/core"
	"github.com/gin-gonic/gin"
)

// =============================================================
// 1. TRANG CÀI ĐẶT CẤU HÌNH MASTER
// =============================================================
func TrangCaiDatCauHinhMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, _ := core.LayKhachHang(shopID, userID)

	// Đổi tên template mapping cho đồng bộ
	c.HTML(http.StatusOK, "master_cai_dat_cau_hinh", gin.H{
		"TieuDe":       "Cấu Hình Hệ Thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,

		// Nạp dữ liệu
		"ListDanhMuc":    core.LayDanhSachDanhMuc(shopID),
		"ListThuongHieu": core.LayDanhSachThuongHieu(shopID),
		"ListBLN":        core.LayDanhSachBienLoiNhuan(shopID),
		"ListNCC":        core.LayDanhSachNhaCungCap(shopID),
	})
}

// =============================================================
// 2. API LƯU NHÀ CUNG CẤP
// =============================================================
func API_LuuNhaCungCapMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maNCC := strings.TrimSpace(c.PostForm("ma_nha_cung_cap"))
	tenNCC := strings.TrimSpace(c.PostForm("ten_nha_cung_cap"))
	sdt := strings.TrimSpace(c.PostForm("dien_thoai"))
	email := strings.TrimSpace(c.PostForm("email"))
	diaChi := strings.TrimSpace(c.PostForm("dia_chi"))
	mst := strings.TrimSpace(c.PostForm("ma_so_thue"))
	nguoiLh := strings.TrimSpace(c.PostForm("nguoi_lien_he"))
	nganHang := strings.TrimSpace(c.PostForm("ngan_hang"))
	ghiChu := strings.TrimSpace(c.PostForm("ghi_chu"))
	trangThai := 0
	if c.PostForm("trang_thai") == "on" {
		trangThai = 1
	}
	isNew := c.PostForm("is_new") == "true"

	if tenNCC == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên nhà cung cấp không được để trống!"})
		return
	}

	var targetRow int

	if isNew {
		maMoi := maNCC
		if maMoi == "" { maMoi = core.TaoMaNhaCungCapMoi(shopID) }
		targetRow = core.DongBatDau_NhaCungCap + len(core.LayDanhSachNhaCungCap(shopID))

		newNCC := &core.NhaCungCap{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			MaNhaCungCap:   maMoi,
			TenNhaCungCap:  tenNCC,
			DienThoai:      sdt,
			Email:          email,
			DiaChi:         diaChi,
			MaSoThue:       mst,
			NguoiLienHe:    nguoiLh,
			NganHang:       nganHang,
			TrangThai:      trangThai,
			GhiChu:         ghiChu,
		}
		
		core.KhoaHeThong.Lock()
		core.CacheNhaCungCap[shopID] = append(core.CacheNhaCungCap[shopID], newNCC)
		core.KhoaHeThong.Unlock()
		
		maNCC = maMoi
	} else {
		list := core.LayDanhSachNhaCungCap(shopID)
		var found *core.NhaCungCap
		for _, item := range list {
			if item.MaNhaCungCap == maNCC { found = item; break }
		}

		if found == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy Nhà cung cấp để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet

		core.KhoaHeThong.Lock()
		found.TenNhaCungCap = tenNCC
		found.DienThoai = sdt
		found.Email = email
		found.DiaChi = diaChi
		found.MaSoThue = mst
		found.NguoiLienHe = nguoiLh
		found.NganHang = nganHang
		found.TrangThai = trangThai
		found.GhiChu = ghiChu
		core.KhoaHeThong.Unlock()
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_MaNhaCungCap, maNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TenNhaCungCap, tenNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DienThoai, sdt)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_Email, email)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DiaChi, diaChi)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_MaSoThue, mst)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NguoiLienHe, nguoiLh)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NganHang, nganHang)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TrangThai, trangThai)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_GhiChu, ghiChu)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Nhà cung cấp thành công!"})
}
