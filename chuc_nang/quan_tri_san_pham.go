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

// =============================================================
// 1. TRANG QUẢN LÝ (HIỂN THỊ)
// =============================================================
func TrangQuanLySanPham(c *gin.Context) {
	userID := c.GetString("USER_ID")
	
	kh, found := core.LayKhachHang(userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 1. Check quyền (Dùng HTML đơn giản để tránh lỗi Template phức tạp)
	if !core.KiemTraQuyen(kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối: Bạn không có quyền xem sản phẩm.</h3><a href="/">Về trang chủ</a>`))
		return
	}

	// 2. LẤY VÀ LỌC SẠCH DANH SÁCH SẢN PHẨM (FIX TRẮNG TRANG)
	rawList := core.LayDanhSachSanPham()
	var cleanList []*core.SanPham

	for _, sp := range rawList {
		// Logic lọc an toàn:
		// 1. Phải khác nil
		// 2. Mã sản phẩm phải có giá trị (đã trim space từ Core)
		if sp != nil && sp.MaSanPham != "" {
			cleanList = append(cleanList, sp)
		}
	}

	listDM := core.LayDanhSachDanhMuc()
	listTH := core.LayDanhSachThuongHieu()

	// 3. Render
	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh, 
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       cleanList, // Truyền danh sách đã lọc sạch
		"ListDanhMuc":    listDM, 
		"ListThuongHieu": listTH,
	})
}

// ... (Phần API Lưu giữ nguyên như cũ) ...
// =============================================================
// 2. API XỬ LÝ LƯU (THÊM / SỬA)
// =============================================================
func API_LuuSanPham(c *gin.Context) {
	// ... (Code cũ của bạn giữ nguyên phần này) ...
    // Copy lại nội dung hàm API_LuuSanPham, xuLyTags, taoMaSPMoi từ file cũ vào đây
    // Hoặc giữ nguyên phần dưới nếu bạn chỉ copy đè hàm TrangQuanLySanPham
    
    vaiTro := c.GetString("USER_ROLE")
    // ...
    // (Để tiết kiệm không gian, tôi chỉ hiển thị phần thay đổi quan trọng ở trên)
    // Bạn hãy giữ nguyên logic lưu, chỉ thay thế hàm TrangQuanLySanPham là đủ.
    
    // --- CODE DƯỚI ĐÂY LÀ PHẦN LƯU (COPY LẠI NẾU CẦN) ---
	maSP       := strings.TrimSpace(c.PostForm("ma_san_pham"))
	thuongHieu := c.PostForm("ma_thuong_hieu")

	giaBanStr := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr  = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _ := strconv.ParseFloat(giaBanStr, 64)

	if maSP == "" {
		if !core.KiemTraQuyen(vaiTro, "product.create") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thêm sản phẩm mới!"})
			return
		}
	} else {
		if !core.KiemTraQuyen(vaiTro, "product.edit") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa sản phẩm!"})
			return
		}
	}

	tenSP       := strings.TrimSpace(c.PostForm("ten_san_pham"))
	tenRutGon   := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku         := strings.TrimSpace(c.PostForm("sku"))
	danhMucRaw  := c.PostForm("ma_danh_muc")
	danhMuc     := xuLyTags(danhMucRaw) 
	donVi       := c.PostForm("don_vi")
	mauSac      := c.PostForm("mau_sac")
	hinhAnh     := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	thongSo     := c.PostForm("thong_so")
	moTa        := c.PostForm("mo_ta_chi_tiet")
	baoHanh, _  := strconv.Atoi(c.PostForm("bao_hanh_thang"))
	tinhTrang   := c.PostForm("tinh_trang")
	ghiChu      := c.PostForm("ghi_chu")
	
	trangThai := 0
	if c.PostForm("trang_thai") == "on" { trangThai = 1 }

	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được để trống!"})
		return
	}

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

	sp.TenSanPham = tenSP
	sp.TenRutGon = tenRutGon
	sp.Sku = sku
	sp.GiaBanLe = giaBan
	sp.MaDanhMuc = danhMuc
	sp.MaThuongHieu = thuongHieu
	sp.DonVi = donVi
	sp.MauSac = mauSac
	sp.UrlHinhAnh = hinhAnh
	sp.ThongSo = thongSo
	sp.MoTaChiTiet = moTa
	sp.BaoHanhThang = baoHanh
	sp.TinhTrang = tinhTrang
	sp.TrangThai = trangThai
	sp.GhiChu = ghiChu
	sp.NgayCapNhat = nowStr

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
		ghi(sheetID, sheet, targetRow, core.CotSP_Sku, sp.Sku)
		ghi(sheetID, sheet, targetRow, core.CotSP_MaDanhMuc, sp.MaDanhMuc)
		ghi(sheetID, sheet, targetRow, core.CotSP_MaThuongHieu, sp.MaThuongHieu)
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

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}

type TagifyItem struct {
	Value string `json:"value"`
}

func xuLyTags(raw string) string {
	if raw == "" { return "" }
	if !strings.Contains(raw, "[") { return raw }
	var items []TagifyItem
	err := json.Unmarshal([]byte(raw), &items)
	if err != nil { return raw }
	var values []string
	for _, item := range items {
		val := strings.TrimSpace(item.Value)
		if val != "" { values = append(values, val) }
	}
	return strings.Join(values, "|")
}

func taoMaSPMoi() string {
	return core.TaoMaSPMoi("") 
}
