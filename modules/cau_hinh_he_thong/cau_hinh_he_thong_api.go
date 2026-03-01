package cau_hinh_he_thong

import (
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

// Khởi tạo Dependency Injection (Tiêm phụ thuộc)
var (
	repo    = CauHinhRepo{}
	service = CauHinhService{repo: repo}
)

func API_LuuNhaCungCap(c *gin.Context) {
	// 1. Lấy thông tin từ Middleware 5 lớp
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// 2. Thu gom dữ liệu rác (Raw Data) từ Form chuyển thành DTO Sạch
	hanMuc, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("han_muc_cong_no"), ".", ""), 64)
	noDauKy, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("cong_no_dau_ky"), ".", ""), 64)
	tt := 0
	if c.PostForm("trang_thai") == "on" || c.PostForm("trang_thai") == "1" { tt = 1 }

	jsonInfo := strings.TrimSpace(c.PostForm("thong_tin_them_json"))
	if jsonInfo == "" { jsonInfo = "{}" }

	dto := DTO_LuuNhaCungCap{
		IsNew:              c.PostForm("is_new") == "true",
		MaNhaCungCap:       strings.TrimSpace(c.PostForm("ma_nha_cung_cap")),
		TenNhaCungCap:      strings.TrimSpace(c.PostForm("ten_nha_cung_cap")),
		MaSoThue:           strings.TrimSpace(c.PostForm("ma_so_thue")),
		DienThoai:          strings.TrimSpace(c.PostForm("dien_thoai")),
		KhuVuc:             strings.TrimSpace(c.PostForm("khu_vuc")),
		HanMucCongNo:       hanMuc,
		CongNoDauKy:        noDauKy,
		ThongTinThemJson:   jsonInfo,
		TrangThai:          tt,
		NguoiThaoTac:       userID,
		// ... map nốt các trường còn lại
	}

	// 3. Ném cho Service xử lý (Không dùng lệnh IF check quyền gì ở đây cả)
	err := service.XuLyLuuNhaCungCap(shopID, dto)

	// 4. Trả kết quả
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Nhà Cung Cấp thành công!"})
}
