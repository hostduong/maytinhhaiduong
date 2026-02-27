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

	// Lấy danh sách toàn bộ nhân sự để làm Danh bạ Chat
	listAll := core.LayDanhSachKhachHang(masterShopID)
	var listChat []*core.KhachHang

	for _, kh := range listAll {
		khCopy := *kh 
		// Lấy lịch sử hộp thư liên quan đến nhân sự này
		khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
		listChat = append(listChat, &khCopy)
	}

	c.HTML(http.StatusOK, "master_tin_nhan", gin.H{
		"TieuDe":   "Tin nhắn",
		"NhanVien": me,
		"ListChat": listChat, 
	})
}

// =========================================================
// 2. API ĐÁNH DẤU ĐÃ ĐỌC
// =========================================================
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
	
	sender, _ := core.LayKhachHang(shopID, userID)
	chucVuNguoiGui := "Hệ Thống Master"
	tenNguoiGui := "Nền tảng 99k.vn"
	if sender != nil {
		if sender.ChucVu != "" { chucVuNguoiGui = sender.ChucVu }
		if sender.TenKhachHang != "" { tenNguoiGui = sender.TenKhachHang }
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
		TenNguoiGui:    tenNguoiGui,
		ChucVuNguoiGui: chucVuNguoiGui,
	}
	core.ThemMoiTinNhan(shopID, newMsg)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi", "data": newMsg})
}
