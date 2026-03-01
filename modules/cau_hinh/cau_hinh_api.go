package cau_hinh

import (
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

var (
	repo    = Repo{}
	service = Service{repo: repo}
)

func API_LuuNhaCungCap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	hanMuc, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("han_muc_cong_no"), ".", ""), 64)
	noDauKy, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("cong_no_dau_ky"), ".", ""), 64)
	ck, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("chiet_khau_mac_dinh"), ".", ""), 64)
	tt := 0; if c.PostForm("trang_thai") == "on" || c.PostForm("trang_thai") == "1" { tt = 1 }
	js := strings.TrimSpace(c.PostForm("thong_tin_them_json")); if js == "" { js = "{}" }

	dto := DTO_LuuNhaCungCap{
		IsNew: c.PostForm("is_new") == "true", MaNhaCungCap: strings.TrimSpace(c.PostForm("ma_nha_cung_cap")),
		TenNhaCungCap: strings.TrimSpace(c.PostForm("ten_nha_cung_cap")), MaSoThue: strings.TrimSpace(c.PostForm("ma_so_thue")),
		DienThoai: strings.TrimSpace(c.PostForm("dien_thoai")), Email: strings.TrimSpace(c.PostForm("email")),
		KhuVuc: strings.TrimSpace(c.PostForm("khu_vuc")), DiaChi: strings.TrimSpace(c.PostForm("dia_chi")),
		NguoiLienHe: strings.TrimSpace(c.PostForm("nguoi_lien_he")), NganHang: strings.TrimSpace(c.PostForm("ngan_hang")),
		NhomNhaCungCap: strings.TrimSpace(c.PostForm("nhom_nha_cung_cap")), LoaiNhaCungCap: strings.TrimSpace(c.PostForm("loai_nha_cung_cap")),
		DieuKhoanThanhToan: strings.TrimSpace(c.PostForm("dieu_khoan_thanh_toan")), ChietKhauMacDinh: ck,
		HanMucCongNo: hanMuc, CongNoDauKy: noDauKy, ThongTinThemJson: js, TrangThai: tt,
		GhiChu: strings.TrimSpace(c.PostForm("ghi_chu")), NguoiThaoTac: userID,
	}

	if err := service.XuLyLuuNhaCungCap(shopID, dto); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Thành công!"})
}
