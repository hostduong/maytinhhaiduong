package chuc_nang

import (
	"encoding/json" // [THÊM] Để xử lý Tags JSON
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

	// [SỬA LẠI ĐOẠN NÀY ĐỂ TRÁNH LỖI PANIC]
	// Nếu không có quyền, hiển thị thông báo đơn giản thay vì load trang Dashboard
	if !core.KiemTraQuyen(kh.VaiTroQuyenHan, "product.view") {
		c.HTML(http.StatusForbidden, "khung_giao_dien", gin.H{
			"TieuDe":       "Từ chối truy cập",
			"DaDangNhap":   true,
			"TenNguoiDung": kh.TenKhachHang,
			"QuyenHan":     kh.VaiTroQuyenHan,
			"NoiDung":      "<h3>⛔ Bạn không có quyền truy cập trang này (product.view).</h3><p>Vui lòng liên hệ quản trị viên.</p>",
		})
		return
	}

	listSP := core.LayDanhSachSanPham()
	listDM := core.LayDanhSachDanhMuc()
	listTH := core.LayDanhSachThuongHieu()

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh, 
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       listSP,
		"ListDanhMuc":    listDM, 
		"ListThuongHieu": listTH,
	})
}

// =============================================================
// 2. API XỬ LÝ LƯU (THÊM / SỬA)
// =============================================================
func API_LuuSanPham(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")

	maSP      := strings.TrimSpace(c.PostForm("ma_san_pham"))
	thuongHieu := c.PostForm("ma_thuong_hieu") // [LẤY SỚM] Để dùng cho sinh mã

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
	
	// [SỬA LẠI] Xử lý danh mục ngăn cách bởi dấu |
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
		
		// [QUAN TRỌNG] Truyền thương hiệu vào để sinh mã (HD + Brand + YYMM...)
		maSP = core.TaoMaSPMoi(thuongHieu) 
		
		sp = &core.SanPham{
			SpreadsheetID: sheetID,
			MaSanPham:     maSP,
			NgayTao:       nowStr,
			NguoiTao:      userID, // Người tạo ban đầu
		}
	} else {
		foundSP, ok := core.LayChiTietSanPham(maSP)
		if ok {
			sp = foundSP
		} else {
			sp = &core.SanPham{SpreadsheetID: sheetID, MaSanPham: maSP, NgayTao: nowStr}
		}
	}

	// Gán dữ liệu mới
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
	
	// Luôn cập nhật ngày sửa (Người tạo không đổi)
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

// [MỚI] Struct để hứng JSON từ Tagify
type TagifyItem struct {
	Value string `json:"value"`
}

// [MỚI] Hàm xử lý Tags chuẩn JSON -> Pipe (|)
// Input:  [{"value":"Laptop"}, {"value":"Dell"}]
// Output: Laptop|Dell (Giữ nguyên thứ tự ưu tiên)
func xuLyTags(raw string) string {
	if raw == "" { return "" }
	
	// 1. Nếu không phải JSON (trường hợp sửa thủ công hoặc lỗi), trả về nguyên gốc
	if !strings.Contains(raw, "[") { return raw }

	// 2. Parse JSON
	var items []TagifyItem
	err := json.Unmarshal([]byte(raw), &items)
	if err != nil {
		return raw // Fallback
	}

	// 3. Gom lại thành mảng string
	var values []string
	for _, item := range items {
		val := strings.TrimSpace(item.Value)
		if val != "" {
			values = append(values, val)
		}
	}

	// 4. Nối bằng dấu gạch đứng
	return strings.Join(values, "|")
}

func taoMaSPMoi() string {
	// Wrapper cũ, không dùng nữa nhưng để đó tránh lỗi compile nếu có chỗ nào lỡ gọi
	return core.TaoMaSPMoi("") 
}
