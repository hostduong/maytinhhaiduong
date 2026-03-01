package chuc_nang

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

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
	
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[masterShopID]
	core.KhoaHeThong.RUnlock()

	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTro { mapStyle[v.MaVaiTro] = v }

	var listChat []*core.KhachHang
	for _, kh := range listAll {
		khCopy := *kh 
		
		if khCopy.MaKhachHang == "0000000000000000000" {
			// [FIX]: Chặn Bot hiển thị toàn bộ tin nhắn của hệ thống. Chỉ hiện tin của Mình và Bot.
			khCopy.StyleLevel, khCopy.StyleTheme = 0, 9 
			var myBotInbox []*core.TinNhan
			allMyMsgs := core.LayHopThuNguoiDung(masterShopID, userID, vaiTro)
			for _, m := range allMyMsgs {
				// Lọc riêng tin nhắn hệ thống hoặc chat trực tiếp với Bot
				if m.NguoiGuiID == "0000000000000000000" || m.NguoiNhanID == "0000000000000000000" || m.LoaiTinNhan == "SYSTEM" || m.LoaiTinNhan == "ALL" {
					myBotInbox = append(myBotInbox, m)
				}
			}
			khCopy.Inbox = myBotInbox
		} else {
			khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
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

func API_GuiTinNhanChat(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	nguoiNhanID := strings.TrimSpace(c.PostForm("nguoi_nhan_id"))
	noiDung := strings.TrimSpace(c.PostForm("noi_dung"))
	
	if nguoiNhanID == "" || noiDung == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Thiếu thông tin người nhận hoặc nội dung!"})
		return
	}

	// [MỚI]: Kích hoạt chức năng Trả lời dưới danh nghĩa Hệ Thống (Bot)
	sendAsBot := c.PostForm("send_as_bot")
	tieuDe := strings.TrimSpace(c.PostForm("tieu_de"))
	
	senderID := userID
	msgType := "CHAT"
	
	if sendAsBot == "1" {
		// Chỉ Level 0, 1, 2 mới được dùng quyền này
		if core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE")) <= 2 {
			senderID = "0000000000000000000" // Ép ID Bot
			msgType = "AUTO" // Đổi sang AUTO để Front-end vẽ khung Card ở giữa
			if tieuDe == "" { tieuDe = "Phản hồi từ Hệ thống" }
		}
	} else {
		tieuDe = "" // Chat 1-1 thường thì không cần tiêu đề
	}

	loc := time.FixedZone("ICT", 7*3600)
	nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")
	msgID := fmt.Sprintf("MSG_%d_%s", time.Now().UnixNano(), nguoiNhanID) 

	newMsg := &core.TinNhan{
		MaTinNhan:      msgID,
		LoaiTinNhan:    msgType,
		NguoiGuiID:     senderID,         
		NguoiNhanID:    nguoiNhanID,           
		TieuDe:         tieuDe,
		NoiDung:        noiDung,
		NgayTao:        nowStr,
	}
	core.ThemMoiTinNhan(shopID, newMsg)
	
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi", "data": newMsg})
}
