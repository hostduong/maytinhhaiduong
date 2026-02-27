package chuc_nang_master

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// =========================================================
// 1. GIAO DIỆN TRANG CHAT MASTER
// =========================================================
func TrangTinNhanMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID") 
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	me, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	listAll := core.LayDanhSachKhachHang(masterShopID)
	
	// Lấy Danh sách Vai trò để map StyleLevel và StyleTheme
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[masterShopID]
	core.KhoaHeThong.RUnlock()

	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTro { mapStyle[v.MaVaiTro] = v }

	var listChat []*core.KhachHang
	for _, kh := range listAll {
		khCopy := *kh 
		khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
		
		// Bơm Style VIP vào object
		if khCopy.MaKhachHang == "0000000000000000000" {
			khCopy.StyleLevel, khCopy.StyleTheme = 0, 9 
		} else {
			if vInfo, ok := mapStyle[khCopy.VaiTroQuyenHan]; ok {
				khCopy.StyleLevel, khCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme
			} else {
				khCopy.StyleLevel, khCopy.StyleTheme = 9, 0 
			}
		}
		if khCopy.MaKhachHang == "0000000000000000001" { khCopy.StyleLevel = 0 }

		listChat = append(listChat, &khCopy)
	}

	c.HTML(http.StatusOK, "master_tin_nhan", gin.H{
		"TieuDe":   "Tin nhắn",
		"NhanVien": me,
		"ListChat": listChat, 
	})
}

func API_DanhDauDaDocMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	msgID := strings.TrimSpace(c.PostForm("msg_id"))

	core.DanhDauDocTinNhan(masterShopID, userID, msgID)
	c.JSON(200, gin.H{"status": "ok"})
}

// =========================================================
// 3. API GỬI TIN NHẮN TRỰC TIẾP (CHAT)
// =========================================================
func API_GuiTinNhanChat(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	nguoiNhanID := strings.TrimSpace(c.PostForm("nguoi_nhan_id"))
	noiDung := strings.TrimSpace(c.PostForm("noi_dung"))
	
	if nguoiNhanID == "" || noiDung == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Thiếu thông tin người nhận hoặc nội dung!"})
		return
	}

	loc := time.FixedZone("ICT", 7*3600)
	nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")
	msgID := fmt.Sprintf("MSG_%d_%s", time.Now().UnixNano(), nguoiNhanID) 

	newMsg := &core.TinNhan{
		MaTinNhan:      msgID,
		LoaiTinNhan:    "CHAT",
		NguoiGuiID:     userID,         
		NguoiNhanID:    nguoiNhanID,           
		TieuDe:         "Tin nhắn hệ thống",
		NoiDung:        noiDung,
		NgayTao:        nowStr,
	}
	core.ThemMoiTinNhan(shopID, newMsg)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi", "data": newMsg})
}
