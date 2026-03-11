package product_master

import (
	"encoding/json"
	"net/http"
	"strings"

	"app/config"
	"app/core"
	"github.com/gin-gonic/gin"
)

// DTO giả lập cấu trúc cũ để HTML không bị sụp đổ
type ViewProductDTO struct {
	DataJSON     string
	MaSanPham    string
	TenSanPham   string
	TenRutGon    string
	Slug         string
	MaSKU        string
	TenSKU       string
	SKUChinh     int
	TrangThai    int
	MaDanhMuc    string
	MaThuongHieu string
	DonVi        string
	UrlHinhAnh   string
	GiaBan       float64
	GiaNiemYet   float64
	NgayCapNhat  string
}

func TrangQuanLySanPhamMaster(c *gin.Context) {
	defer func() { if err := recover(); err != nil { c.String(500, "LỖI HỆ THỐNG: %v", err) } }()

	masterShopID := c.GetString("SHOP_ID") 
	adminShopID := config.BienCauHinh.IdFileSheetAdmin 
	
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	kh, found := core.LayKhachHang(masterShopID, userID)
	if !found || kh == nil { c.Redirect(http.StatusFound, "/login"); return }

	if vaiTro != "quan_tri_he_thong" && !core.KiemTraQuyen(masterShopID, vaiTro, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`<h3>⛔ Truy cập bị từ chối</h3><a href="/">Về trang chủ</a>`))
		return
	}

	// Tạm thời fix cứng ngành Điện tử, sau này lấy từ URL (c.Query("nganh"))
	maNganh := "dien_tu" 

	rawList := core.LayDanhSachSanPham(adminShopID, maNganh)
	
	var cleanList []ViewProductDTO 
	var fullList []ViewProductDTO  

	for _, sp := range rawList {
		if sp != nil && sp.MaSanPham != "" {
			b, _ := json.Marshal(sp)
			jsonStr := string(b)
			
			var mainSKU *core.ProductSKU
			for i := range sp.SKU {
				if sp.SKU[i].MaSKU == sp.SKUChinh { mainSKU = &sp.SKU[i]; break }
			}
			if mainSKU == nil && len(sp.SKU) > 0 { mainSKU = &sp.SKU[0] }

			dto := ViewProductDTO{
				DataJSON: jsonStr, // Bí mật ném JSON cho HTML tự bung
				MaSanPham: sp.MaSanPham,
				TenSanPham: sp.TenSanPham,
				TenRutGon: sp.TenRutGon,
				Slug: sp.Slug,
				TrangThai: sp.TrangThai,
				MaDanhMuc: strings.Join(sp.MaDanhMuc, "|"),
				MaThuongHieu: sp.MaThuongHieu,
				NgayCapNhat: sp.QuanLy.NgayCapNhat,
			}
			
			if mainSKU != nil {
				dto.MaSKU = mainSKU.MaSKU
				dto.TenSKU = mainSKU.TenSKU
				dto.SKUChinh = 1
				if mainSKU.TrangThai == -1 { dto.TrangThai = -1 }
				dto.DonVi = mainSKU.DonVi
				if len(mainSKU.HinhAnh) > 0 { dto.UrlHinhAnh = strings.Join(mainSKU.HinhAnh, "|") }
				dto.GiaBan = mainSKU.Gia.GiaBan
				dto.GiaNiemYet = mainSKU.Gia.GiaNiemYet
			}

			cleanList = append(cleanList, dto)
			fullList = append(fullList, dto)
		}
	}

	// Lấy Cấu hình Thuộc tính từ RAM (O(1))
	core.KhoaHeThong.RLock()
	thuocTinhData := core.CacheThuocTinh
	core.KhoaHeThong.RUnlock()

	c.HTML(http.StatusOK, "product_master", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       cleanList, 
		"DanhSachFull":   fullList,  
		"ListDanhMuc":    core.LayDanhSachDanhMuc(adminShopID),    
		"ListThuongHieu": core.LayDanhSachThuongHieu(adminShopID), 
		"ListBLN":        core.LayDanhSachBienLoiNhuan(adminShopID), 
		
		// [BỔ SUNG 2 DÒNG NÀY LÀ XONG]
		"CauHinhThuocTinh": thuocTinhData, 
		"MaNganh":          maNganh,
	})
}
