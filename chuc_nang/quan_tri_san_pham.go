package chuc_nang

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"app/core"
	"github.com/gin-gonic/gin"
)

// =============================================================
// 1. TRANG QUẢN LÝ (LIST)
// =============================================================
func TrangQuanLySanPham(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.String(500, "LỖI HỆ THỐNG: %v", err)
		}
	}()

	// [SAAS] Lấy Context
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// Lấy thông tin người dùng (Cần ShopID để biết user của shop nào)
	kh, found := core.LayKhachHang(shopID, userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// [SAAS] Check quyền theo Shop
	if !core.KiemTraQuyen(shopID, kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối</h3><a href="/">Về trang chủ</a>`))
		return
	}

	// [SAAS] Lấy danh sách của riêng Shop này
	rawList := core.LayDanhSachSanPham(shopID)
	
	// Lọc và làm sạch danh sách hiển thị
	var cleanList []*core.SanPham
	uniqueDM := make(map[string]bool)
	uniqueTH := make(map[string]bool)

	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			cleanList = append(cleanList, sp)
			
			// Thu thập danh mục để filter nhanh (nếu cần)
			if sp.DanhMuc != "" {
				parts := strings.Split(sp.DanhMuc, "|")
				for _, p := range parts {
					if p = strings.TrimSpace(p); p != "" { uniqueDM[p] = true }
				}
			}
			// Thu thập thương hiệu
			th := strings.TrimSpace(sp.ThuongHieu)
			if th != "" { uniqueTH[th] = true }
		}
	}

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       cleanList,
		
		// [SAAS] Load các danh sách bổ trợ theo ShopID
		"ListDanhMuc":    core.LayDanhSachDanhMuc(shopID),    
		"ListThuongHieu": core.LayDanhSachThuongHieu(shopID), 
		"ListBLN":        core.LayDanhSachBienLoiNhuan(shopID), 
	})
}

// =============================================================
// 2. API LƯU SẢN PHẨM (THÊM / SỬA)
// =============================================================
func API_LuuSanPham(c *gin.Context) {
	// [SAAS] Lấy Context
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	
	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	
	// Check Quyền (Truyền ShopID)
	if maSP == "" {
		if !core.KiemTraQuyen(shopID, vaiTro, "product.create") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thêm sản phẩm!"})
			return
		}
	} else {
		if !core.KiemTraQuyen(shopID, vaiTro, "product.edit") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa sản phẩm!"})
			return
		}
	}

	tenSP := strings.TrimSpace(c.PostForm("ten_san_pham"))
	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được trống!"})
		return
	}

	// --- 1. XỬ LÝ DỮ LIỆU ĐẦU VÀO ---
	
	// Xử lý tiền tệ (Xóa dấu chấm, phẩy)
	giaNhap, _  := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(c.PostForm("gia_nhap"), ".", ""), ",", ""), 64)
	giaBanLe, _ := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", ""), ",", ""), 64)
	giamGia, _  := strconv.ParseFloat(c.PostForm("giam_gia"), 64)
	
	// Tính giá bán thực
	giaBanThuc := giaBanLe
	if giamGia > 0 {
		giaBanThuc = giaBanLe * (1 - giamGia/100)
	}

	// Xử lý Tags (JSON Tagify)
	thuongHieu := xuLyTags(c.PostForm("ma_thuong_hieu"))
	danhMuc    := xuLyTags(c.PostForm("ma_danh_muc")) 
	donVi      := xuLyTags(c.PostForm("don_vi"))
	mauSac     := xuLyTags(c.PostForm("mau_sac"))
	tinhTrang  := xuLyTags(c.PostForm("tinh_trang"))
	
	tenRutGon  := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku        := strings.TrimSpace(c.PostForm("sku"))
	moTa       := c.PostForm("mo_ta_chi_tiet")
	hinhAnh    := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	thongSo    := c.PostForm("thong_so")
	ghiChu     := c.PostForm("ghi_chu")
	
	// Xử lý bảo hành
	bhNum := c.PostForm("bao_hanh_num")
	bhUnit := c.PostForm("bao_hanh_unit")
	baoHanh := ""
	if bhNum != "" {
		baoHanh = bhNum + " " + bhUnit
	}

	// Tạo Slug và Trạng thái
	slug := taoSlugChuan(tenSP)
	trangThai := 0
	if c.PostForm("trang_thai") == "on" { trangThai = 1 }

	// --- 2. KHỞI TẠO HOẶC TÌM SẢN PHẨM ---
	var sp *core.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")

	if maSP == "" {
		// --- THÊM MỚI ---
		isNew = true
		
		// Tìm mã danh mục đầu tiên để sinh mã SP (VD: MAIN -> MAIN0001)
		firstDM := ""
		if danhMuc != "" { firstDM = strings.Split(danhMuc, "|")[0] }
		
		// [SAAS] Tìm mã danh mục trong Shop hiện tại
		maCodeDM := core.TimMaDanhMucTheoTen(shopID, firstDM) 
		
		// [SAAS] Sinh mã mới dựa trên ShopID
		maSP = core.TaoMaSPMoi(shopID, maCodeDM) 
		
		sp = &core.SanPham{
			SpreadsheetID: shopID, // Gán vào Shop hiện tại
			MaSanPham:     maSP,
			NgayTao:       nowStr,
			NguoiTao:      userID,
		}
	} else {
		// --- CẬP NHẬT ---
		// [SAAS] Tìm chi tiết trong Shop hiện tại
		foundSP, ok := core.LayChiTietSanPham(shopID, maSP)
		if ok { 
			sp = foundSP 
		} else { 
			// Fallback nếu không tìm thấy trong RAM (hiếm)
			sp = &core.SanPham{SpreadsheetID: shopID, MaSanPham: maSP} 
		}
	}

	// Khóa RAM để cập nhật an toàn
	if !isNew { core.KhoaHeThong.Lock() }

	// --- 3. GÁN DỮ LIỆU VÀO STRUCT ---
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
	sp.BaoHanh    = baoHanh
	sp.TrangThai  = trangThai
	
	sp.GiaNhap    = giaNhap
	sp.GiaBanLe   = giaBanLe
	sp.GiamGia    = giamGia
	sp.GiaBanThuc = giaBanThuc
	
	sp.GhiChu     = ghiChu
	sp.NgayCapNhat= nowStr

	if !isNew {
		core.KhoaHeThong.Unlock()
	} else {
		// Nếu là mới, tính dòng tiếp theo và thêm vào RAM
		currentList := core.LayDanhSachSanPham(shopID)
		sp.DongTrongSheet = core.DongBatDau_SanPham + len(currentList)
		
		// [SAAS] Thêm vào RAM của Shop
		core.ThemSanPhamVaoRam(sp)
	}

	// --- 4. ĐẨY VÀO HÀNG CHỜ GHI (QUEUE) ---
	targetRow := sp.DongTrongSheet
	if targetRow > 0 {
		ghi := core.ThemVaoHangCho
		sheet := "SAN_PHAM"

		// Sử dụng shopID làm ID file sheet
		ghi(shopID, sheet, targetRow, core.CotSP_MaSanPham, sp.MaSanPham)
		ghi(shopID, sheet, targetRow, core.CotSP_TenSanPham, sp.TenSanPham)
		ghi(shopID, sheet, targetRow, core.CotSP_TenRutGon, sp.TenRutGon)
		ghi(shopID, sheet, targetRow, core.CotSP_Slug, sp.Slug)
		ghi(shopID, sheet, targetRow, core.CotSP_Sku, sp.Sku)
		ghi(shopID, sheet, targetRow, core.CotSP_DanhMuc, sp.DanhMuc)
		ghi(shopID, sheet, targetRow, core.CotSP_ThuongHieu, sp.ThuongHieu)
		ghi(shopID, sheet, targetRow, core.CotSP_DonVi, sp.DonVi)
		ghi(shopID, sheet, targetRow, core.CotSP_MauSac, sp.MauSac)
		ghi(shopID, sheet, targetRow, core.CotSP_UrlHinhAnh, sp.UrlHinhAnh)
		ghi(shopID, sheet, targetRow, core.CotSP_ThongSo, sp.ThongSo)
		ghi(shopID, sheet, targetRow, core.CotSP_MoTaChiTiet, sp.MoTaChiTiet)
		ghi(shopID, sheet, targetRow, core.CotSP_BaoHanh, sp.BaoHanh)
		ghi(shopID, sheet, targetRow, core.CotSP_TinhTrang, sp.TinhTrang)
		ghi(shopID, sheet, targetRow, core.CotSP_TrangThai, sp.TrangThai)
		
		ghi(shopID, sheet, targetRow, core.CotSP_GiaNhap, sp.GiaNhap)
		ghi(shopID, sheet, targetRow, core.CotSP_GiaBanLe, sp.GiaBanLe)
		ghi(shopID, sheet, targetRow, core.CotSP_GiamGia, sp.GiamGia)
		ghi(shopID, sheet, targetRow, core.CotSP_GiaBanThuc, sp.GiaBanThuc)
		
		ghi(shopID, sheet, targetRow, core.CotSP_GhiChu, sp.GhiChu)
		ghi(shopID, sheet, targetRow, core.CotSP_NguoiTao, sp.NguoiTao)
		ghi(shopID, sheet, targetRow, core.CotSP_NgayTao, sp.NgayTao)
		ghi(shopID, sheet, targetRow, core.CotSP_NgayCapNhat, sp.NgayCapNhat)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}

// =============================================================
// 3. CÁC HÀM HELPER (XỬ LÝ CHUỖI/TAGS)
// =============================================================

type TagifyItem struct { Value string `json:"value"` }

func xuLyTags(raw string) string {
	if raw == "" { return "" }
	// Nếu không phải JSON (người dùng nhập tay không dùng Tagify)
	if !strings.Contains(raw, "[") { return raw }
	
	var items []TagifyItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil { return raw }
	
	var values []string
	for _, item := range items {
		if v := strings.TrimSpace(item.Value); v != "" {
			values = append(values, v)
		}
	}
	return strings.Join(values, "|")
}

func taoSlugChuan(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "đ", "d")
	
	// Map ký tự có dấu sang không dấu
	patterns := map[string]string{
		"[áàảãạăắằẳẵặâấầẩẫậ]": "a",
		"[éèẻẽẹêếềểễệ]":       "e",
		"[iíìỉĩị]":            "i",
		"[óòỏõọôốồổỗộơớờởỡợ]": "o",
		"[úùủũụưứừửữự]":       "u",
		"[ýỳỷỹỵ]":             "y",
	}
	for p, r := range patterns {
		re := regexp.MustCompile(p)
		s = re.ReplaceAllString(s, r)
	}
	
	// Xóa ký tự đặc biệt
	reInvalid := regexp.MustCompile(`[^a-z0-9]+`)
	s = reInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
