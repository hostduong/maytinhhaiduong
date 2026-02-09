package chuc_nang

import (
	"encoding/json"
	"fmt" // [MỚI] Thêm fmt để in lỗi
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
	// [QUAN TRỌNG] BẮT LỖI PANIC ĐỂ KHÔNG BỊ TRẮNG TRANG
	defer func() {
		if err := recover(); err != nil {
			logMsg := fmt.Sprintf("❌ LỖI HỆ THỐNG (PANIC): %v", err)
			fmt.Println(logMsg) // In ra terminal
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(fmt.Sprintf(`
				<h3>Hệ thống gặp lỗi khi hiển thị dữ liệu!</h3>
				<p>Chi tiết: %v</p>
				<p>Vui lòng chụp ảnh này gửi kỹ thuật viên.</p>
				<a href="/admin/tong-quan">Quay lại Dashboard</a>
			`, err)))
		}
	}()

	userID := c.GetString("USER_ID")
	
	kh, found := core.LayKhachHang(userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 1. Check quyền
	if !core.KiemTraQuyen(kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối: Bạn không có quyền xem sản phẩm.</h3><a href="/">Về trang chủ</a>`))
		return
	}

	// 2. LẤY DỮ LIỆU & LỌC SẠCH (TRÁNH NIL POINTER)
	
	// A. Lọc Sản Phẩm
	rawList := core.LayDanhSachSanPham()
	var cleanList []*core.SanPham
	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			cleanList = append(cleanList, sp)
		}
	}

	// B. Lọc Danh Mục (Nguyên nhân tiềm ẩn gây trắng trang)
	rawDM := core.LayDanhSachDanhMuc()
	var cleanDM []*core.DanhMuc
	for _, dm := range rawDM {
		if dm != nil && dm.MaDanhMuc != "" {
			cleanDM = append(cleanDM, dm)
		}
	}

	// C. Lọc Thương Hiệu (Nguyên nhân tiềm ẩn gây trắng trang)
	rawTH := core.LayDanhSachThuongHieu()
	var cleanTH []*core.ThuongHieu
	for _, th := range rawTH {
		if th != nil && th.MaThuongHieu != "" {
			cleanTH = append(cleanTH, th)
		}
	}

	// 3. Render
	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh, 
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       cleanList, // Đã sạch
		"ListDanhMuc":    cleanDM,   // Đã sạch
		"ListThuongHieu": cleanTH,   // Đã sạch
	})
}

// =============================================================
// 2. API XỬ LÝ LƯU (THÊM / SỬA)
// =============================================================
func API_LuuSanPham(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")

	maSP       := strings.TrimSpace(c.PostForm("ma_san_pham"))
	thuongHieu := c.PostForm("ma_thuong_hieu")

	giaBanStr := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr  = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _ := strconv.ParseFloat(giaBanStr, 64)

	// --- CHECK QUYỀN ---
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
		
		spCu, ok := core.LayChiTietSanPham(maSP)
		if ok && spCu.GiaBanLe != giaBan {
			if !core.KiemTraQuyen(vaiTro, "product.edit_price") {
				c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị viên được sửa giá bán!"})
				return
			}
		}
	}

	// --- LẤY DỮ LIỆU ---
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

	// --- LOGIC CORE ---
	var sp *core.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")
	sheetID := cau_hinh.BienCauHinh.IdFileSheet

	core.KhoaHeThong.Lock()
	
	if maSP == "" {
		isNew = true
		// Sinh mã mới (HD + Brand + YYMM...)
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

	// Gán dữ liệu
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

	// --- GHI XUỐNG SHEET ---
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

// Struct JSON Tagify
type TagifyItem struct {
	Value string `json:"value"`
}

// Hàm xử lý Tags: Input `[{"value":"A"}, {"value":"B"}]` -> Output `A|B`
func xuLyTags(raw string) string {
	if raw == "" { return "" }
	if !strings.Contains(raw, "[") { return raw }

	var items []TagifyItem
	err := json.Unmarshal([]byte(raw), &items)
	if err != nil { return raw }

	var values []string
	for _, item := range items {
		val := strings.TrimSpace(item.Value)
		if val != "" {
			values = append(values, val)
		}
	}
	return strings.Join(values, "|")
}

func taoMaSPMoi() string {
	return core.TaoMaSPMoi("") 
}
