package product_master

import (
	"regexp"
	"strings"
	"app/core"
)

// Ép tên SP thành URL hợp lệ (VD: Laptop Dell -> laptop-dell)
func Repo_TaoSlugChuan(s string) string {
	s = strings.ToLower(s); s = strings.ReplaceAll(s, "đ", "d")
	patterns := map[string]string{ "[áàảãạăắằẳẵặâấầẩẫậ]": "a", "[éèẻẽẹêếềểễệ]": "e", "[iíìỉĩị]": "i", "[óòỏõọôốồổỗộơớờởỡợ]": "o", "[úùủũụưứừửữự]": "u", "[ýỳỷỹỵ]": "y" }
	for p, r := range patterns { re := regexp.MustCompile(p); s = re.ReplaceAllString(s, r) }
	reInvalid := regexp.MustCompile(`[^a-z0-9]+`); s = reInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// Băm toàn bộ data thành 1 chuỗi text không dấu để Frontend tìm kiếm cực nhanh
func Repo_BuildSearchText(sp *core.ProductJSON) string {
	var parts []string
	parts = append(parts, sp.MaSanPham, sp.TenSanPham, sp.TenRutGon)
	for _, tag := range sp.Tags { parts = append(parts, tag) }
	for _, sku := range sp.SKU {
		parts = append(parts, sku.MaSKU, sku.TenSKU, sku.Barcode)
	}
	
	raw := strings.Join(parts, " ")
	raw = strings.ToLower(raw)
	return raw
}
