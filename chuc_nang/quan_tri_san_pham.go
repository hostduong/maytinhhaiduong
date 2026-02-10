package chuc_nang

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// TrangQuanLySanPham : Hiển thị
func TrangQuanLySanPham(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.String(500, "LỖI HỆ THỐNG: %v", err)
		}
	}()

	userID := c.GetString("USER_ID")
	kh, found := core.LayKhachHang(userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	if !core.KiemTraQuyen(kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối</h3><a href="/">Về trang chủ</a>`))
		return
	}

	rawList := core.LayDanhSachSanPham()
	var cleanList []*core.SanPham
	
	uniqueDM := make(map[string]bool)
	uniqueTH := make(map[string]bool)

	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			cleanList = append(cleanList, sp)
			if sp.DanhMuc != "" {
				parts := strings.Split(sp.DanhMuc, "|")
				for _, p := range parts {
					if p = strings.TrimSpace(p); p != "" { uniqueDM[p] = true }
				}
			}
			if sp.ThuongHieu != "" { uniqueTH[strings.TrimSpace(sp.ThuongHieu)] = true }
		}
	}

	var listDM, listTH []string
	for k := range uniqueDM { listDM = append(listDM, k) }
	for k := range uniqueTH { listTH = append(listTH, k) }

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       cleanList,
		"ListDanhMuc":    listDM,
		"ListThuongHieu": listTH,
	})
}

// API_LuuSanPham : Xử lý Lưu
func API_LuuSanPham(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	
	if maSP == "" {
		if !core.KiemTraQuyen(vaiTro, "product.create") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thêm!"})
			return
		}
	} else {
		if !core.KiemTraQuyen(vaiTro, "product.edit") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa!"})
			return
		}
	}

	tenSP := strings.TrimSpace(c.PostForm("ten_san_pham"))
	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được trống!"})
		return
	}

	giaBanStr := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr  = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _ := strconv.ParseFloat(giaBanStr, 64)

	thuongHieu := strings.TrimSpace(c.PostForm("ma_thuong_hieu"))
	tenRutGon  := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku        := strings.TrimSpace(c.PostForm("sku"))
	danhMuc    := xuLyTags(c.PostForm("ma_danh_muc")) 
	donVi      := c.PostForm("don_vi")
	mauSac     := c.PostForm("mau_sac")
	tinhTrang  := c.PostForm("tinh_trang")
	moTa       := c.PostForm("mo_ta_chi_tiet")
	hinhAnh    := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	thongSo    := c.PostForm("thong_so")
	baoHanh, _ := strconv.Atoi(c.PostForm("bao_hanh_thang"))
	ghiChu     := c.PostForm("ghi_chu")
	slug       := taoSlug(tenSP)

	trangThai := 0
	if c.PostForm("trang_thai") == "on" { trangThai = 1 }

	var sp *core.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")
	sheetID := cau_hinh.BienCauHinh.IdFileSheet

	core.KhoaHeThong.Lock()
	
	if maSP == "" {
		isNew = true
		maSP = core.TaoMaSPMoi(thuongHieu) 
		sp = &core.SanPham{
			SpreadsheetID: sheetID,
			MaSanPham:     maSP,
			NgayTao:       nowStr,
			NguoiTao:      userID,
		}
	} else {
		foundSP, ok := core.LayChiTietSanPham(maSP)
		if ok {
			sp = foundSP
		} else {
			sp = &core.SanPham{SpreadsheetID: sheetID, MaSanPham: maSP, NgayTao: nowStr}
		}
	}

	// [QUAN TRỌNG] Cập nhật RAM (Vì sp là pointer, gán vào đây là RAM đổi luôn)
	sp.TenSanPham = tenSP
	sp.TenRutGon  = tenRutGon
	sp.Slug       = slug
	sp.Sku        = sku
	sp.DanhMuc    = danhMuc
	sp.ThuongHieu = thuongHieu
	sp.DonVi      = donVi
	sp.MauSac     = mauSac
	sp.TinhTrang  = tinhTrang
	sp.MoTaChiTiet= moTa
	sp.UrlHinhAnh = hinhAnh
	sp.ThongSo    = thongSo
	sp.BaoHanhThang = baoHanh
	sp.TrangThai  = trangThai
	sp.GiaBanLe   = giaBan
	sp.GhiChu     = ghiChu
	sp.NgayCapNhat= nowStr

	if isNew {
		sp.DongTrongSheet = core.DongBatDau_SanPham + len(core.LayDanhSachSanPham()) 
		core.ThemSanPhamVaoRam(sp)
	}
	core.KhoaHeThong.Unlock()

	targetRow := sp.DongTrongSheet
	if targetRow > 0 {
		ghi := core.ThemVaoHangCho
		sheet := "SAN_PHAM"

		ghi(sheetID, sheet, targetRow, core.CotSP_MaSanPham, sp.MaSanPham)
		ghi(sheetID, sheet, targetRow, core.CotSP_TenSanPham, sp.TenSanPham)
		ghi(sheetID, sheet, targetRow, core.CotSP_TenRutGon, sp.TenRutGon)
		ghi(sheetID, sheet, targetRow, core.CotSP_Slug, sp.Slug)
		ghi(sheetID, sheet, targetRow, core.CotSP_Sku, sp.Sku)
		ghi(sheetID, sheet, targetRow, core.CotSP_DanhMuc, sp.DanhMuc)
		ghi(sheetID, sheet, targetRow, core.CotSP_ThuongHieu, sp.ThuongHieu)
		ghi(sheetID, sheet, targetRow, core.CotSP_DonVi, sp.DonVi)
		ghi(sheetID, sheet, targetRow, core.CotSP_MauSac, sp.MauSac)
		ghi(sheetID, sheet, targetRow, core.CotSP_UrlHinhAnh, sp.UrlHinhAnh)
		ghi(sheetID, sheet, targetRow, core.CotSP_ThongSo, sp.ThongSo)
		ghi(sheetID, sheet, targetRow, core.CotSP_MoTaChiTiet, sp.MoTaChiTiet)
		ghi(sheetID, sheet, targetRow, core.CotSP_BaoHanhThang, sp.BaoHanhThang)
		ghi(sheetID, sheet, targetRow, core.CotSP_TinhTrang, sp.TinhTrang)
		ghi(sheetID, sheet, targetRow, core.CotSP_TrangThai, sp.TrangThai)
		ghi(sheetID, sheet, targetRow, core.CotSP_GiaBanLe, sp.GiaBanLe)
		ghi(sheetID, sheet, targetRow, core.CotSP_GhiChu, sp.GhiChu)
		ghi(sheetID, sheet, targetRow, core.CotSP_NguoiTao, sp.NguoiTao)
		ghi(sheetID, sheet, targetRow, core.CotSP_NgayTao, sp.NgayTao)
		ghi(sheetID, sheet, targetRow, core.CotSP_NgayCapNhat, sp.NgayCapNhat)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu thành công!"})
}

// Helpers
type TagifyItem struct { Value string `json:"value"` }

func xuLyTags(raw string) string {
	if raw == "" { return "" }
	if !strings.Contains(raw, "[") { return raw }
	var items []TagifyItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil { return raw }
	var values []string
	for _, item := range items {
		if v := strings.TrimSpace(item.Value); v != "" { values = append(values, v) }
	}
	return strings.Join(values, "|")
}

func taoSlug(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "/", "-")
	return text
}
