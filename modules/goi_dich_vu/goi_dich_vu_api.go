package goi_dich_vu

import (
	"app/core"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
)

var (
	repo    = Repo{}
	service = Service{repo: repo}
)

func API_LuuGoiDichVu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	
	gn, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("gia_niem_yet"), ".", ""), 64)
	gb, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("gia_ban"), ".", ""), 64)
	th, _ := strconv.Atoi(c.PostForm("thoi_han_ngay"))
	sl, _ := strconv.Atoi(c.PostForm("so_luong_con_lai"))
	tt := 0; if c.PostForm("trang_thai") == "on" || c.PostForm("trang_thai") == "1" { tt = 1 }

	dto := DTO_LuuGoiDichVu{
		IsNew: c.PostForm("is_new") == "true",
		MaGoi: strings.ToUpper(strings.TrimSpace(c.PostForm("ma_goi"))),
		TenGoi: c.PostForm("ten_goi"), LoaiGoi: c.PostForm("loai_goi"),
		ThoiHanNgay: th, GiaNiemYet: gn, GiaBan: gb,
		CodesJson: c.PostForm("codes_json"), GioiHanJson: c.PostForm("gioi_han_json"),
		MoTa: c.PostForm("mo_ta"), NhanHienThi: c.PostForm("nhan_hien_thi"),
		NgayBatDau: c.PostForm("ngay_bat_dau"), NgayKetThuc: c.PostForm("ngay_ket_thuc"),
		SoLuongConLai: sl, TrangThai: tt,
	}

	if err := service.XuLyLuu(shopID, dto); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "ok", "msg": "Thành công!"})
}
