package chuc_nang

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// =============================================================
// STRUCT HỨNG DỮ LIỆU TỪ GIAO DIỆN MỚI (TỪNG TAB SKU)
// =============================================================
type InputSKU struct {
	MaSKU        string  `json:"ma_sku"`
	TenSKU       string  `json:"ten_sku"`
	SKUChinh     int     `json:"sku_chinh"`
	TrangThai    int     `json:"trang_thai"`
	TenSanPham   string  `json:"ten_san_pham"`
	TenRutGon    string  `json:"ten_rut_gon"`
	MaDanhMuc    string  `json:"ma_danh_muc"`
	MaThuongHieu string  `json:"ma_thuong_hieu"`
	DonVi        string  `json:"don_vi"`
	MauSac       string  `json:"mau_sac"`
	KhoiLuong    float64 `json:"khoi_luong"`
	KichThuoc    string  `json:"kich_thuoc"`
	UrlHinhAnh   string  `json:"url_hinh_anh"`
	ThongSoHTML  string  `json:"thong_so_html"`
	MoTaHTML     string  `json:"mo_ta_html"`
	BaoHanh      string  `json:"bao_hanh"`
	TinhTrang    string  `json:"tinh_trang"`
	GiaNhap      float64 `json:"gia_nhap"`
	PhanTramLai  float64 `json:"phan_tram_lai"` // <-- BỔ SUNG CỘT LÃI (U)
	GiaNiemYet   float64 `json:"gia_niem_yet"`
	PhanTramGiam float64 `json:"phan_tram_giam"`
	SoTienGiam   float64 `json:"so_tien_giam"`
	GiaBan       float64 `json:"gia_ban"`
	GhiChu       string  `json:"ghi_chu"`
}

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

	// Lấy thông tin người dùng
	kh, found := core.LayKhachHang(shopID, userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// [SAAS] Check quyền
	if !core.KiemTraQuyen(shopID, kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối</h3><a href="/">Về trang chủ</a>`))
		return
	}

	// [SAAS] Lấy TOÀN BỘ danh sách dòng (SKU) của Shop này
	rawList := core.LayDanhSachSanPham(shopID)
	
	var cleanList []*core.SanPham // Danh sách hiển thị Bảng (Chỉ chứa các SKU Chính)
	var fullList []*core.SanPham  // Danh sách đầy đủ để Frontend dựng lại Modal

	// Gom nhóm theo MaSanPham
	groupSP := make(map[string][]*core.SanPham)

	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			fullList = append(fullList, sp)
				groupSP[sp.MaSanPham] = append(groupSP[sp.MaSanPham], sp)
		}
	}

	// Lọc ra SKU Chính để đại diện hiển thị ra Bảng HTML
	for _, dsSKU := range groupSP {
		var spChinh *core.SanPham
		for _, sp := range dsSKU {
			if sp.SKUChinh == 1 {
				spChinh = sp
				break
			}
		}
		// Fallback: Lỗi người dùng không set SKU Chính, ép lấy dòng đầu tiên làm đại diện
		if spChinh == nil && len(dsSKU) > 0 {
			spChinh = dsSKU[0]
		}
		if spChinh != nil {
			cleanList = append(cleanList, spChinh)
		}
	}

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		
		"DanhSach":       cleanList, // Truyền cho Table (1 dòng/1 SP)
		"DanhSachFull":   fullList,  // Truyền để dựng lại dữ liệu các Tab biến thể
		
		"ListDanhMuc":    core.LayDanhSachDanhMuc(shopID),    
		"ListThuongHieu": core.LayDanhSachThuongHieu(shopID), 
		"ListBLN":        core.LayDanhSachBienLoiNhuan(shopID), 
	})
}

// =============================================================
// 2. API LƯU SẢN PHẨM (XỬ LÝ DỮ LIỆU NHIỀU TAB)
// =============================================================
func API_LuuSanPham(c *gin.Context) {
	// [SAAS] Lấy Context
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	userID := c.GetString("USER_ID")
	
	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	
	// Check Quyền
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

	dataJSON := c.PostForm("data_skus")
	var inputSKUs []InputSKU
	if err := json.Unmarshal([]byte(dataJSON), &inputSKUs); err != nil || len(inputSKUs) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Dữ liệu không hợp lệ!"})
		return
	}

	// Đảm bảo SKU Chính
	hasMain := false
	for _, sku := range inputSKUs {
		if sku.SKUChinh == 1 { hasMain = true; break }
	}
	if !hasMain { inputSKUs[0].SKUChinh = 1 }

	nowStr := time.Now().Format("2006-01-02 15:04:05")

	// 1. NẾU LÀ TẠO MỚI (Frontend để trống mã) -> Tự sinh
	if maSP == "" {
		firstCodeDM := ""
		if inputSKUs[0].MaDanhMuc != "" { 
			parsedDM := xuLyTags(inputSKUs[0].MaDanhMuc)
			if parsedDM != "" { firstCodeDM = strings.Split(parsedDM, "|")[0] }
		}
		maSP = core.TaoMaSPMoi(shopID, firstCodeDM) 
	} else {
		// 2. LOGIC MỚI: NẾU CÓ MÃ RỒI (Frontend tự điền) -> KIỂM TRA XEM CÓ PHẢI MỚI KHÔNG?
		// Kiểm tra trong RAM xem mã này đã tồn tại chưa
		listCheck := core.LayNhomSanPham(shopID, maSP)
		
		if len(listCheck) == 0 {
			// AHA! Mã này chưa từng tồn tại -> Đây là TẠO MỚI nhưng Frontend tự điền mã.
			// Chúng ta cần "Cập nhật ngược" lại Slot cho Danh mục để lần sau nó không bị trùng.
			
			// Bước A: Lấy Mã Danh Mục từ dữ liệu gửi lên
			firstCodeDM := ""
			if inputSKUs[0].MaDanhMuc != "" { 
				parsedDM := xuLyTags(inputSKUs[0].MaDanhMuc)
				if parsedDM != "" { firstCodeDM = strings.Split(parsedDM, "|")[0] }
			}

			// Bước B: Tách số từ Mã SP (Ví dụ MAIN0005 -> lấy số 5)
			// Giả sử quy tắc là Ký tự + Số. Ta xoá hết chữ, chỉ lấy số.
			re := regexp.MustCompile(`[0-9]+`)
			nums := re.FindAllString(maSP, -1)
			if len(nums) > 0 {
				// Lấy cụm số cuối cùng (thường là slot)
				lastNumStr := nums[len(nums)-1] 
				if slotMoi, err := strconv.Atoi(lastNumStr); err == nil {
					// Bước C: Gọi hàm cập nhật Slot cưỡng bức
					if firstCodeDM != "" {
						core.CapNhatSlotThuCong(shopID, firstCodeDM, slotMoi)
					}
				}
			}
		}
	}

	// ... (Đoạn dưới giữ nguyên logic lưu bình thường) ...
	existingSKUs := core.LayNhomSanPham(shopID, maSP)
	existingMap := make(map[string]*core.SanPham)
	for _, sp := range existingSKUs {
		existingMap[sp.LấyIDDuyNhat()] = sp
	}
	processedSKUs := make(map[string]bool) 

	core.KhoaHeThong.Lock()
	defer core.KhoaHeThong.Unlock()

	for i, in := range inputSKUs {
		skuID := in.MaSKU
		if skuID == "" { skuID = fmt.Sprintf("%s-%02d", maSP, i+1) }
		
		var sp *core.SanPham
		isNewSKU := false
		
		if exist, ok := existingMap[skuID]; ok {
			sp = exist 
			processedSKUs[skuID] = true
		} else {
			isNewSKU = true
			currentList := core.CacheSanPham[shopID]
			sp = &core.SanPham{
				SpreadsheetID:  shopID,
				DongTrongSheet: core.DongBatDau_SanPham + len(currentList),
				MaSanPham:      maSP,
				MaSKU:          skuID,
				NgayTao:        nowStr,
				NguoiTao:       userID,
			}
		}

		sp.TenSanPham   = strings.TrimSpace(in.TenSanPham)
		sp.TenRutGon    = strings.TrimSpace(in.TenRutGon)
		sp.Slug         = taoSlugChuan(sp.TenSanPham)
		sp.TenSKU       = strings.TrimSpace(in.TenSKU)
		sp.SKUChinh     = in.SKUChinh
		sp.TrangThai    = in.TrangThai
		sp.MaDanhMuc    = xuLyTags(in.MaDanhMuc)
		sp.MaThuongHieu = xuLyTags(in.MaThuongHieu)
		sp.DonVi        = xuLyTags(in.DonVi)
		sp.MauSac       = xuLyTags(in.MauSac)
		sp.KhoiLuong    = in.KhoiLuong
		sp.KichThuoc    = in.KichThuoc
		sp.UrlHinhAnh   = strings.TrimSpace(in.UrlHinhAnh)
		sp.ThongSoHTML  = in.ThongSoHTML
		sp.MoTaHTML     = in.MoTaHTML
		sp.BaoHanh      = in.BaoHanh
		sp.TinhTrang    = xuLyTags(in.TinhTrang)
		sp.GiaNhap      = in.GiaNhap
		sp.PhanTramLai  = in.PhanTramLai 
		sp.GiaNiemYet   = in.GiaNiemYet
		sp.PhanTramGiam = in.PhanTramGiam
		sp.SoTienGiam   = in.SoTienGiam
		sp.GiaBan       = in.GiaBan
		sp.GhiChu       = in.GhiChu
		sp.NgayCapNhat  = nowStr
		sp.NguoiCapNhat = userID

		if isNewSKU {
			core.CacheSanPham[shopID] = append(core.CacheSanPham[shopID], sp)
			core.CacheMapSKU[core.TaoCompositeKey(shopID, sp.LấyIDDuyNhat())] = sp
			core.CacheGroupSanPham[core.TaoCompositeKey(shopID, sp.MaSanPham)] = append(core.CacheGroupSanPham[core.TaoCompositeKey(shopID, sp.MaSanPham)], sp)
		}

		ghi := core.ThemVaoHangCho
		sheet := "SAN_PHAM"
		r := sp.DongTrongSheet
		
		ghi(shopID, sheet, r, core.CotSP_MaSanPham, sp.MaSanPham)
		ghi(shopID, sheet, r, core.CotSP_TenSanPham, sp.TenSanPham)
		ghi(shopID, sheet, r, core.CotSP_TenRutGon, sp.TenRutGon)
		ghi(shopID, sheet, r, core.CotSP_Slug, sp.Slug)
		ghi(shopID, sheet, r, core.CotSP_MaSKU, sp.MaSKU)
		ghi(shopID, sheet, r, core.CotSP_TenSKU, sp.TenSKU)
		ghi(shopID, sheet, r, core.CotSP_SKUChinh, sp.SKUChinh)
		ghi(shopID, sheet, r, core.CotSP_TrangThai, sp.TrangThai)
		ghi(shopID, sheet, r, core.CotSP_MaDanhMuc, sp.MaDanhMuc)
		ghi(shopID, sheet, r, core.CotSP_MaThuongHieu, sp.MaThuongHieu)
		ghi(shopID, sheet, r, core.CotSP_DonVi, sp.DonVi)
		ghi(shopID, sheet, r, core.CotSP_MauSac, sp.MauSac)
		ghi(shopID, sheet, r, core.CotSP_KhoiLuong, sp.KhoiLuong)
		ghi(shopID, sheet, r, core.CotSP_KichThuoc, sp.KichThuoc)
		ghi(shopID, sheet, r, core.CotSP_UrlHinhAnh, sp.UrlHinhAnh)
		ghi(shopID, sheet, r, core.CotSP_ThongSoHTML, sp.ThongSoHTML)
		ghi(shopID, sheet, r, core.CotSP_MoTaHTML, sp.MoTaHTML)
		ghi(shopID, sheet, r, core.CotSP_BaoHanh, sp.BaoHanh)
		ghi(shopID, sheet, r, core.CotSP_TinhTrang, sp.TinhTrang)
		ghi(shopID, sheet, r, core.CotSP_GiaNhap, sp.GiaNhap)
		ghi(shopID, sheet, r, core.CotSP_PhanTramLai, sp.PhanTramLai) 
		ghi(shopID, sheet, r, core.CotSP_GiaNiemYet, sp.GiaNiemYet)
		ghi(shopID, sheet, r, core.CotSP_PhanTramGiam, sp.PhanTramGiam)
		ghi(shopID, sheet, r, core.CotSP_SoTienGiam, sp.SoTienGiam)
		ghi(shopID, sheet, r, core.CotSP_GiaBan, sp.GiaBan)
		ghi(shopID, sheet, r, core.CotSP_GhiChu, sp.GhiChu)
		ghi(shopID, sheet, r, core.CotSP_NguoiTao, sp.NguoiTao)
		ghi(shopID, sheet, r, core.CotSP_NgayTao, sp.NgayTao)
		ghi(shopID, sheet, r, core.CotSP_NguoiCapNhat, sp.NguoiCapNhat)
		ghi(shopID, sheet, r, core.CotSP_NgayCapNhat, sp.NgayCapNhat)
	}

	for skuID, sp := range existingMap {
		if !processedSKUs[skuID] {
			sp.TrangThai = -1 
			sp.SKUChinh = 0
			sp.NgayCapNhat = nowStr
			sp.NguoiCapNhat = userID
			
			core.ThemVaoHangCho(shopID, "SAN_PHAM", sp.DongTrongSheet, core.CotSP_TrangThai, -1)
			core.ThemVaoHangCho(shopID, "SAN_PHAM", sp.DongTrongSheet, core.CotSP_SKUChinh, 0)
		}
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu toàn bộ biến thể sản phẩm thành công!"})
}

// =============================================================
// 3. CÁC HÀM HELPER (XỬ LÝ CHUỖI/TAGS)
// =============================================================
type TagifyItem struct { Value string `json:"value"` }

// Hàm này nhặt giá trị "value" từ mảng JSON của Tagify
// Vì Frontend đã cấu hình truyền Mã vào Value, nên hàm này tự động trả về đúng Mã (VD: MAIN, GIGA)
func xuLyTags(raw string) string {
	if raw == "" { return "" }
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
	
	reInvalid := regexp.MustCompile(`[^a-z0-9]+`)
	s = reInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
