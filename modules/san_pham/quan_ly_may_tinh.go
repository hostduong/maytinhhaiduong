package san_pham

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

type InputSKUMayTinh struct {
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
	PhanTramLai  float64 `json:"phan_tram_lai"`
	GiaNiemYet   float64 `json:"gia_niem_yet"`
	PhanTramGiam float64 `json:"phan_tram_giam"`
	SoTienGiam   float64 `json:"so_tien_giam"`
	GiaBan       float64 `json:"gia_ban"`
	GhiChu       string  `json:"ghi_chu"`
}

func TrangQuanLyMayTinhMaster(c *gin.Context) {
	defer func() { if err := recover(); err != nil { c.String(500, "LỖI HỆ THỐNG: %v", err) } }()

	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, found := core.LayKhachHang(shopID, userID)
	if !found || kh == nil { c.Redirect(http.StatusFound, "/login"); return }

	if !core.KiemTraQuyen(shopID, kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối</h3><a href="/">Về trang chủ</a>`))
		return
	}

	rawList := core.LayDanhSachSanPhamMayTinh(shopID)
	
	var cleanList []*core.SanPhamMayTinh 
	var fullList []*core.SanPhamMayTinh  
	groupSP := make(map[string][]*core.SanPhamMayTinh)

	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			fullList = append(fullList, sp)
			groupSP[sp.MaSanPham] = append(groupSP[sp.MaSanPham], sp)
		}
	}

	for _, dsSKU := range groupSP {
		var spChinh *core.SanPhamMayTinh
		for _, sp := range dsSKU { if sp.SKUChinh == 1 { spChinh = sp; break } }
		if spChinh == nil && len(dsSKU) > 0 { spChinh = dsSKU[0] }
		if spChinh != nil { cleanList = append(cleanList, spChinh) }
	}

	// TRỎ TÊN ĐÚNG VÀO TEMPLATE MỚI
	c.HTML(http.StatusOK, "master_quan_ly_may_tinh", gin.H{
		"TieuDe":         "Quản lý sản phẩm (Máy Tính)",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       cleanList, 
		"DanhSachFull":   fullList,  
		"ListDanhMuc":    core.LayDanhSachDanhMuc(shopID),    
		"ListThuongHieu": core.LayDanhSachThuongHieu(shopID), 
		"ListBLN":        core.LayDanhSachBienLoiNhuan(shopID), 
	})
}

func API_LayChiTietMayTinhMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	maSP := c.Param("ma_sp")
	if maSP == "" { c.JSON(200, gin.H{"status": "error", "msg": "Thiếu mã sản phẩm!"}); return }

	core.KhoaHeThong.RLock()
	listSKU := core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(masterShopID, maSP)]
	core.KhoaHeThong.RUnlock()

	if len(listSKU) == 0 { c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy sản phẩm!"}); return }
	c.JSON(200, gin.H{"status": "ok", "data": listSKU})
}

func API_LuuMayTinhMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	userID := c.GetString("USER_ID")
	
	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	if maSP == "" {
		if !core.KiemTraQuyen(shopID, vaiTro, "product.create") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thêm sản phẩm!"}); return
		}
	} else {
		if !core.KiemTraQuyen(shopID, vaiTro, "product.edit") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa sản phẩm!"}); return
		}
	}

	dataJSON := c.PostForm("data_skus")
	var inputSKUs []InputSKUMayTinh
	if err := json.Unmarshal([]byte(dataJSON), &inputSKUs); err != nil || len(inputSKUs) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Dữ liệu không hợp lệ!"}); return
	}

	hasMain := false
	for _, sku := range inputSKUs { if sku.SKUChinh == 1 { hasMain = true; break } }
	if !hasMain { inputSKUs[0].SKUChinh = 1 }

	loc := time.FixedZone("ICT", 7*3600)
	nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")

	core.KhoaHeThong.RLock()
	existingSKUs := core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(shopID, maSP)]
	core.KhoaHeThong.RUnlock()

	if maSP == "" {
		firstCodeDM := ""
		if inputSKUs[0].MaDanhMuc != "" { 
			parsedDM := xuLyTags(inputSKUs[0].MaDanhMuc)
			if parsedDM != "" { firstCodeDM = strings.Split(parsedDM, "|")[0] }
		}
		maSP = core.TaoMaSPMayTinhMoi(shopID, firstCodeDM)
	} else {
		if len(existingSKUs) == 0 {
			firstCodeDM := ""
			if inputSKUs[0].MaDanhMuc != "" { 
				parsedDM := xuLyTags(inputSKUs[0].MaDanhMuc)
				if parsedDM != "" { firstCodeDM = strings.Split(parsedDM, "|")[0] }
			}
			re := regexp.MustCompile(`[0-9]+`)
			nums := re.FindAllString(maSP, -1)
			if len(nums) > 0 {
				lastNumStr := nums[len(nums)-1] 
				if slotMoi, err := strconv.Atoi(lastNumStr); err == nil {
					if firstCodeDM != "" { core.CapNhatSlotThuCong(shopID, firstCodeDM, slotMoi) }
				}
			}
		}
	}

	existingMap := make(map[string]*core.SanPhamMayTinh)
	for _, sp := range existingSKUs { existingMap[sp.LayIDDuyNhat()] = sp }
	processedSKUs := make(map[string]bool) 

	core.KhoaHeThong.Lock()
	defer core.KhoaHeThong.Unlock()

	for i, in := range inputSKUs {
		skuID := in.MaSKU
		if skuID == "" { skuID = fmt.Sprintf("%s-%02d", maSP, i+1) }
		
		var sp *core.SanPhamMayTinh
		isNewSKU := false
		
		if exist, ok := existingMap[skuID]; ok {
			sp = exist; processedSKUs[skuID] = true
		} else {
			isNewSKU = true
			currentList := core.CacheSanPhamMayTinh[shopID]
			sp = &core.SanPhamMayTinh{
				SpreadsheetID:  shopID,
				DongTrongSheet: core.DongBatDau_SanPhamMayTinh + len(currentList),
				MaSanPham:      maSP,
				MaSKU:          skuID,
			}
		}

		newTenSanPham   := strings.TrimSpace(in.TenSanPham)
		newTenRutGon    := strings.TrimSpace(in.TenRutGon)
		newSlug         := taoSlugChuan(newTenSanPham)
		newTenSKU       := strings.TrimSpace(in.TenSKU)
		newMaDanhMuc    := xuLyTags(in.MaDanhMuc)
		newMaThuongHieu := xuLyTags(in.MaThuongHieu)
		newDonVi        := xuLyTags(in.DonVi)
		newMauSac       := xuLyTags(in.MauSac)
		newUrlHinhAnh   := strings.TrimSpace(in.UrlHinhAnh)
		newTinhTrang    := xuLyTags(in.TinhTrang)

		isChanged := false

		if isNewSKU {
			isChanged = true
			sp.NgayTao = nowStr; sp.NguoiTao = userID
			sp.NgayCapNhat = nowStr; sp.NguoiCapNhat = userID
		} else {
			if sp.TenSanPham != newTenSanPham || sp.TenRutGon != newTenRutGon || sp.Slug != newSlug || sp.TenSKU != newTenSKU ||
				sp.SKUChinh != in.SKUChinh || sp.TrangThai != in.TrangThai || sp.MaDanhMuc != newMaDanhMuc || sp.MaThuongHieu != newMaThuongHieu ||
				sp.DonVi != newDonVi || sp.MauSac != newMauSac || sp.KhoiLuong != in.KhoiLuong || sp.KichThuoc != in.KichThuoc ||
				sp.UrlHinhAnh != newUrlHinhAnh || sp.ThongSoHTML != in.ThongSoHTML || sp.MoTaHTML != in.MoTaHTML || sp.BaoHanh != in.BaoHanh ||
				sp.TinhTrang != newTinhTrang || sp.GiaNhap != in.GiaNhap || sp.PhanTramLai != in.PhanTramLai || sp.GiaNiemYet != in.GiaNiemYet ||
				sp.PhanTramGiam != in.PhanTramGiam || sp.SoTienGiam != in.SoTienGiam || sp.GiaBan != in.GiaBan || sp.GhiChu != in.GhiChu {
				
				isChanged = true
				sp.NgayCapNhat = nowStr; sp.NguoiCapNhat = userID
			}
		}

		if isChanged {
			sp.TenSanPham = newTenSanPham; sp.TenRutGon = newTenRutGon; sp.Slug = newSlug; sp.TenSKU = newTenSKU
			sp.SKUChinh = in.SKUChinh; sp.TrangThai = in.TrangThai; sp.MaDanhMuc = newMaDanhMuc; sp.MaThuongHieu = newMaThuongHieu
			sp.DonVi = newDonVi; sp.MauSac = newMauSac; sp.KhoiLuong = in.KhoiLuong; sp.KichThuoc = in.KichThuoc
			sp.UrlHinhAnh = newUrlHinhAnh; sp.ThongSoHTML = in.ThongSoHTML; sp.MoTaHTML = in.MoTaHTML; sp.BaoHanh = in.BaoHanh
			sp.TinhTrang = newTinhTrang; sp.GiaNhap = in.GiaNhap; sp.PhanTramLai = in.PhanTramLai; sp.GiaNiemYet = in.GiaNiemYet
			sp.PhanTramGiam = in.PhanTramGiam; sp.SoTienGiam = in.SoTienGiam; sp.GiaBan = in.GiaBan; sp.GhiChu = in.GhiChu

			if isNewSKU {
				core.CacheSanPhamMayTinh[shopID] = append(core.CacheSanPhamMayTinh[shopID], sp)
				core.CacheMapSKUMayTinh[core.TaoCompositeKey(shopID, sp.LayIDDuyNhat())] = sp
				core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(shopID, sp.MaSanPham)] = append(core.CacheGroupSanPhamMayTinh[core.TaoCompositeKey(shopID, sp.MaSanPham)], sp)
			}

			ghi := core.ThemVaoHangCho
			sheet := core.TenSheetMayTinh
			r := sp.DongTrongSheet
			
			ghi(shopID, sheet, r, core.CotPC_MaSanPham, sp.MaSanPham); ghi(shopID, sheet, r, core.CotPC_TenSanPham, sp.TenSanPham)
			ghi(shopID, sheet, r, core.CotPC_TenRutGon, sp.TenRutGon); ghi(shopID, sheet, r, core.CotPC_Slug, sp.Slug)
			ghi(shopID, sheet, r, core.CotPC_MaSKU, sp.MaSKU); ghi(shopID, sheet, r, core.CotPC_TenSKU, sp.TenSKU)
			ghi(shopID, sheet, r, core.CotPC_SKUChinh, sp.SKUChinh); ghi(shopID, sheet, r, core.CotPC_TrangThai, sp.TrangThai)
			ghi(shopID, sheet, r, core.CotPC_MaDanhMuc, sp.MaDanhMuc); ghi(shopID, sheet, r, core.CotPC_MaThuongHieu, sp.MaThuongHieu)
			ghi(shopID, sheet, r, core.CotPC_DonVi, sp.DonVi); ghi(shopID, sheet, r, core.CotPC_MauSac, sp.MauSac)
			ghi(shopID, sheet, r, core.CotPC_KhoiLuong, sp.KhoiLuong); ghi(shopID, sheet, r, core.CotPC_KichThuoc, sp.KichThuoc)
			ghi(shopID, sheet, r, core.CotPC_UrlHinhAnh, sp.UrlHinhAnh); ghi(shopID, sheet, r, core.CotPC_ThongSoHTML, sp.ThongSoHTML)
			ghi(shopID, sheet, r, core.CotPC_MoTaHTML, sp.MoTaHTML); ghi(shopID, sheet, r, core.CotPC_BaoHanh, sp.BaoHanh)
			ghi(shopID, sheet, r, core.CotPC_TinhTrang, sp.TinhTrang); ghi(shopID, sheet, r, core.CotPC_GiaNhap, sp.GiaNhap)
			ghi(shopID, sheet, r, core.CotPC_PhanTramLai, sp.PhanTramLai); ghi(shopID, sheet, r, core.CotPC_GiaNiemYet, sp.GiaNiemYet)
			ghi(shopID, sheet, r, core.CotPC_PhanTramGiam, sp.PhanTramGiam); ghi(shopID, sheet, r, core.CotPC_SoTienGiam, sp.SoTienGiam)
			ghi(shopID, sheet, r, core.CotPC_GiaBan, sp.GiaBan); ghi(shopID, sheet, r, core.CotPC_GhiChu, sp.GhiChu)
			
			if isNewSKU {
				ghi(shopID, sheet, r, core.CotPC_NguoiTao, sp.NguoiTao); ghi(shopID, sheet, r, core.CotPC_NgayTao, sp.NgayTao)
			}
			ghi(shopID, sheet, r, core.CotPC_NguoiCapNhat, sp.NguoiCapNhat); ghi(shopID, sheet, r, core.CotPC_NgayCapNhat, sp.NgayCapNhat)
		}
	}

	for skuID, sp := range existingMap {
		if !processedSKUs[skuID] {
			if sp.TrangThai != -1 {
				sp.TrangThai = -1; sp.SKUChinh = 0; sp.NgayCapNhat = nowStr; sp.NguoiCapNhat = userID
				core.ThemVaoHangCho(shopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_TrangThai, -1)
				core.ThemVaoHangCho(shopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_SKUChinh, 0)
				core.ThemVaoHangCho(shopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_NgayCapNhat, nowStr)
				core.ThemVaoHangCho(shopID, core.TenSheetMayTinh, sp.DongTrongSheet, core.CotPC_NguoiCapNhat, userID)
			}
		}
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu thành công!"})
}

// Helper Functions
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

func taoSlugChuan(s string) string {
	s = strings.ToLower(s); s = strings.ReplaceAll(s, "đ", "d")
	patterns := map[string]string{ "[áàảãạăắằẳẵặâấầẩẫậ]": "a", "[éèẻẽẹêếềểễệ]": "e", "[iíìỉĩị]": "i", "[óòỏõọôốồổỗộơớờởỡợ]": "o", "[úùủũụưứừửữự]": "u", "[ýỳỷỹỵ]": "y" }
	for p, r := range patterns { re := regexp.MustCompile(p); s = re.ReplaceAllString(s, r) }
	reInvalid := regexp.MustCompile(`[^a-z0-9]+`); s = reInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
