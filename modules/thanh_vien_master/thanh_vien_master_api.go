package thanh_vien_master

import (
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
)

func API_LuuThanhVienMaster(c *gin.Context) {
	gt := -1
	if c.PostForm("gioi_tinh") == "1" { gt = 1 } else if c.PostForm("gioi_tinh") == "0" { gt = 0 }

	dto := DTO_UpdateThanhVien{
		ShopID:         c.GetString("SHOP_ID"),
		AdminID:        c.GetString("USER_ID"),
		AdminRole:      c.GetString("USER_ROLE"),
		PinXacNhan:     strings.TrimSpace(c.PostForm("pin_xac_nhan")),
		MaKH:           c.PostForm("ma_khach_hang"),
		VaiTro:         c.PostForm("vai_tro"),
		ChucVu:         strings.TrimSpace(c.PostForm("chuc_vu")),
		TrangThai:      c.PostForm("trang_thai"),
		TenCuaHang:     strings.TrimSpace(c.PostForm("ten_cua_hang")),
		TenKhachHang:   strings.TrimSpace(c.PostForm("ten_khach_hang")),
		DienThoai:      strings.TrimSpace(c.PostForm("dien_thoai")),
		NgaySinh:       strings.TrimSpace(c.PostForm("ngay_sinh")),
		DiaChi:         strings.TrimSpace(c.PostForm("dia_chi")),
		MaSoThue:       strings.TrimSpace(c.PostForm("ma_so_thue")),
		GhiChu:         strings.TrimSpace(c.PostForm("ghi_chu")),
		AnhDaiDien:     strings.TrimSpace(c.PostForm("anh_dai_dien")),
		NguonKhachHang: strings.TrimSpace(c.PostForm("nguon_khach_hang")),
		Zalo:           strings.TrimSpace(c.PostForm("zalo")),
		Facebook:       strings.TrimSpace(c.PostForm("facebook")),
		Tiktok:         strings.TrimSpace(c.PostForm("tiktok")),
		MatKhauMoi:     strings.TrimSpace(c.PostForm("mat_khau_moi")),
		PinMoi:         strings.TrimSpace(c.PostForm("pin_moi")),
		GioiTinh:       gt,
	}

	err := Service_LuuThanhVien(dto)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thông tin thành công!"})
}

func API_GuiTinNhanMaster(c *gin.Context) {
	tieuDe := strings.TrimSpace(c.PostForm("tieu_de"))
	noiDung := strings.TrimSpace(c.PostForm("noi_dung"))
	jsonIDs := c.PostForm("danh_sach_id")
	
	if tieuDe == "" || noiDung == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tiêu đề và Nội dung không được để trống!"})
		return
	}

	soLuong, err := Service_GuiTinNhan(
		c.GetString("SHOP_ID"), 
		c.GetString("USER_ID"), 
		c.GetString("USER_ROLE"), 
		tieuDe, noiDung, jsonIDs, c.PostForm("send_as_bot"),
	)

	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": fmt.Sprintf("Đã gửi thông báo thành công cho %d người!", soLuong)})
}
