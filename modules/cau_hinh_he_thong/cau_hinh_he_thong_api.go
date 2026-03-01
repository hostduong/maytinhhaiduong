package cau_hinh_he_thong

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Khởi tạo Dependency Injection (Tiêm phụ thuộc)
// Đảm bảo API gọi Service, Service gọi Repo
var (
	repo    = CauHinhRepo{}
	service = CauHinhService{repo: repo}
)

// ==============================================================================
// API: LƯU NHÀ CUNG CẤP
// ==============================================================================
func API_LuuNhaCungCap(c *gin.Context) {
	// 1. Lấy thân phận người dùng từ Đường ống Middleware
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// 2. Xử lý ép kiểu dữ liệu từ Text sang Số (Loại bỏ dấu chấm ngàn)
	hanMuc, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("han_muc_cong_no"), ".", ""), 64)
	noDauKy, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("cong_no_dau_ky"), ".", ""), 64)
	ckMacDinh, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("chiet_khau_mac_dinh"), ".", ""), 64)

	// Xử lý Checkbox Trạng thái
	trangThai := 0
	if c.PostForm("trang_thai") == "on" || c.PostForm("trang_thai") == "1" {
		trangThai = 1
	}

	// Xử lý Json mở rộng
	jsonInfo := strings.TrimSpace(c.PostForm("thong_tin_them_json"))
	if jsonInfo == "" {
		jsonInfo = "{}"
	}

	// 3. Đóng gói vào DTO (Data Transfer Object) để chuyển xuống Service
	dto := DTO_LuuNhaCungCap{
		IsNew:              c.PostForm("is_new") == "true",
		MaNhaCungCap:       strings.TrimSpace(c.PostForm("ma_nha_cung_cap")),
		TenNhaCungCap:      strings.TrimSpace(c.PostForm("ten_nha_cung_cap")),
		MaSoThue:           strings.TrimSpace(c.PostForm("ma_so_thue")),
		DienThoai:          strings.TrimSpace(c.PostForm("dien_thoai")),
		Email:              strings.TrimSpace(c.PostForm("email")),
		KhuVuc:             strings.TrimSpace(c.PostForm("khu_vuc")),
		DiaChi:             strings.TrimSpace(c.PostForm("dia_chi")),
		NguoiLienHe:        strings.TrimSpace(c.PostForm("nguoi_lien_he")),
		NganHang:           strings.TrimSpace(c.PostForm("ngan_hang")),
		NhomNhaCungCap:     strings.TrimSpace(c.PostForm("nhom_nha_cung_cap")),
		LoaiNhaCungCap:     strings.TrimSpace(c.PostForm("loai_nha_cung_cap")),
		DieuKhoanThanhToan: strings.TrimSpace(c.PostForm("dieu_khoan_thanh_toan")),
		ChietKhauMacDinh:   ckMacDinh,
		HanMucCongNo:       hanMuc,
		CongNoDauKy:        noDauKy,
		ThongTinThemJson:   jsonInfo,
		TrangThai:          trangThai,
		GhiChu:             strings.TrimSpace(c.PostForm("ghi_chu")),
		NguoiThaoTac:       userID,
	}

	// 4. Bàn giao cho Service tính toán và lưu DB
	err := service.XuLyLuuNhaCungCap(shopID, dto)

	// 5. Trả kết quả JSON về cho Frontend
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Nhà Cung Cấp thành công!"})
}

// Bổ sung các hàm API khác ở dưới này sau (VD: API_LuuDanhMuc, API_LuuThuongHieu...)
