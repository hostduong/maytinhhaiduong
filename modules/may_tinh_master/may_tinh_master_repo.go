package may_tinh_master

import (
	"encoding/json"
	"regexp"
	"strings"
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

type TagifyItem struct {
	Value string `json:"value"`
}

func Repo_XuLyTags(raw string) string {
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

func Repo_TaoSlugChuan(s string) string {
	s = strings.ToLower(s); s = strings.ReplaceAll(s, "đ", "d")
	patterns := map[string]string{ "[áàảãạăắằẳẵặâấầẩẫậ]": "a", "[éèẻẽẹêếềểễệ]": "e", "[iíìỉĩị]": "i", "[óòỏõọôốồổỗộơớờởỡợ]": "o", "[úùủũụưứừửữự]": "u", "[ýỳỷỹỵ]": "y" }
	for p, r := range patterns { re := regexp.MustCompile(p); s = re.ReplaceAllString(s, r) }
	reInvalid := regexp.MustCompile(`[^a-z0-9]+`); s = reInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
