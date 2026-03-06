package cua_hang_master

import (
	"strings"
	"github.com/gin-gonic/gin"
)

func API_LuuCuaHangMaster(c *gin.Context) {
	gt := -1
	if c.PostForm("gioi_tinh") == "1" { gt = 1 } else if c.PostForm("gioi_tinh") == "0" { gt = 0 }

	dto := DTO_UpdateCuaHang{
		AdminID:        c.GetString("USER_ID"),
		AdminRole:      c.GetString("USER_ROLE"),
		PinXacNhan:     strings.TrimSpace(c.PostForm("pin_xac_nhan")),
		MaKH:           c.PostForm("ma_khach_hang"),
		TrangThai:      c.PostForm("trang_thai"),
		TenKhachHang:   strings.TrimSpace(c.PostForm("ten_khach_hang")),
		DienThoai:      strings.TrimSpace(c.PostForm("dien_thoai")),
		NgaySinh:       strings.TrimSpace(c.PostForm("ngay_sinh")),
		DiaChi:         strings.TrimSpace(c.PostForm("dia_chi")),
		MaSoThue:       strings.TrimSpace(c.PostForm("ma_so_thue")),
		GhiChu:         strings.TrimSpace(c.PostForm("ghi_chu")),
		AnhDaiDien:     strings.TrimSpace(c.PostForm("anh_dai_dien")),
		Zalo:           strings.TrimSpace(c.PostForm("zalo")),
		Facebook:       strings.TrimSpace(c.PostForm("facebook")),
		Tiktok:         strings.TrimSpace(c.PostForm("tiktok")),
		MatKhauMoi:     strings.TrimSpace(c.PostForm("mat_khau_moi")),
		GioiTinh:       gt,
		
		// [ĐÃ VÁ LỖI]: Hứng 4 trường bị thiếu từ HTML
		VaiTro:         c.PostForm("vai_tro"),
		ChucVu:         strings.TrimSpace(c.PostForm("chuc_vu")),
		NguonKhachHang: strings.TrimSpace(c.PostForm("nguon_khach_hang")),
		PinMoi:         strings.TrimSpace(c.PostForm("pin_moi")),
		
		SpreadsheetID:  strings.TrimSpace(c.PostForm("spreadsheet_id")),
		CustomDomain:   strings.TrimSpace(c.PostForm("custom_domain")),
	}

	err := Service_LuuCuaHang(dto)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thông tin cửa hàng thành công!"})
}
