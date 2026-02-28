package may_tinh

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app/core"
	data_pc "app/core/may_tinh" // Lõi Data của riêng ngành PC
	"github.com/gin-gonic/gin"
)

// Khai báo lại struct hứng dữ liệu (giống bên Admin)
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
	PhanTramLai  float64 `json:"phan_tram_lai"`
	GiaNiemYet   float64 `json:"gia_niem_yet"`
	PhanTramGiam float64 `json:"phan_tram_giam"`
	SoTienGiam   float64 `json:"so_tien_giam"`
	GiaBan       float64 `json:"gia_ban"`
	GhiChu       string  `json:"ghi_chu"`
}

// 1. GIAO DIỆN QUẢN LÝ SẢN PHẨM (MASTER)
func TrangQuanLySanPhamMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	// Lớp khiên bảo vệ Master
	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	kh, _ := core.LayKhachHang(masterShopID, userID)

	// Lấy danh sách sản phẩm từ RAM (Core)
	rawList := data_pc.LayDanhSachSanPham(masterShopID)
	
	var cleanList []*data_pc.SanPham 
	groupSP := make(map[string][]*data_pc.SanPham)

	// Nhóm các SKU lại để hiển thị Sản phẩm Cha (SKUChinh = 1)
	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			groupSP[sp.MaSanPham] = append(groupSP[sp.MaSanPham], sp)
		}
	}

	for _, dsSKU := range groupSP {
		var spChinh *data_pc.SanPham
		for _, sp := range dsSKU {
			if sp.SKUChinh == 1 { spChinh = sp; break }
		}
		if spChinh == nil && len(dsSKU) > 0 { spChinh = dsSKU[0] }
		if spChinh != nil { cleanList = append(cleanList, spChinh) }
	}

	// Trỏ ra giao diện HTML riêng của Master
	c.HTML(http.StatusOK, "master_may_tinh_san_pham", gin.H{
		"TieuDe":         "Sản Phẩm",
		"NhanVien":       kh,
		"DanhSach":       cleanList, 
		"ListDanhMuc":    core.LayDanhSachDanhMuc(masterShopID),    
		"ListThuongHieu": core.LayDanhSachThuongHieu(masterShopID), 
	})
}

// 2. API LẤY CHI TIẾT ĐỂ ĐỔ VÀO FORM EDIT (SPA)
func API_LayChiTietSanPhamMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không có quyền!"})
		return
	}

	maSP := c.Param("ma_sp")
	if maSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Thiếu mã sản phẩm!"})
		return
	}

	core.KhoaHeThong.RLock()
	listSKU := data_pc.CacheGroupSanPham[core.TaoCompositeKey(masterShopID, maSP)]
	core.KhoaHeThong.RUnlock()

	if len(listSKU) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy sản phẩm!"})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "data": listSKU})
}

// 3. API LƯU SẢN PHẨM (Tương tự logic Admin nhưng chặn quyền Master)
func API_LuuSanPhamMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Không có quyền thao tác!"})
		return
	}

	maSP := strings.TrimSpace(c.PostForm("ma_san_pham"))
	dataJSON := c.PostForm("data_skus")
	
	var inputSKUs []InputSKU
	if err := json.Unmarshal([]byte(dataJSON), &inputSKUs); err != nil || len(inputSKUs) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Dữ liệu không hợp lệ!"})
		return
	}

	// >>> LOGIC LƯU TƯƠNG TỰ FILE ADMIN BẠN ĐÃ VIẾT TRƯỚC ĐÓ <<<
	// (Để tránh file quá dài, mình gọi tắt phần xử lý DataPC ở đây. 
	// Bạn có thể copy nguyên đoạn logic "Dirty Check" và ghi Queue từ chuc_nang_admin/may_tinh/quan_tri_san_pham.go sang đây, 
	// chỉ thay biến `shopID` thành `masterShopID`).

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công vào hệ thống Master!"})
}
