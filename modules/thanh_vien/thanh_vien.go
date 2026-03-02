package thanh_vien

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"app/config"
	"app/core"
	"github.com/gin-gonic/gin"
)

func TrangQuanLyThanhVienMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	myLevel := core.LayCapBacVaiTro(masterShopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.Redirect(http.StatusFound, "/")
		return
	}
	
	me, _ := core.LayKhachHang(masterShopID, userID)
	listAll := core.LayDanhSachKhachHang(masterShopID)
	
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[masterShopID]
	core.KhoaHeThong.RUnlock()

	mapStyle := make(map[string]core.VaiTroInfo)
	for _, v := range listVaiTro { mapStyle[v.MaVaiTro] = v }

	var listView []*core.KhachHang
	for _, kh := range listAll {
		khCopy := *kh 
		khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
		if khCopy.MaKhachHang == "0000000000000000000" {
			khCopy.StyleLevel, khCopy.StyleTheme = 0, 9 
		} else {
			if vInfo, ok := mapStyle[khCopy.VaiTroQuyenHan]; ok {
				khCopy.StyleLevel, khCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme
			} else {
				khCopy.StyleLevel, khCopy.StyleTheme = 9, 0 
			}
		}
		listView = append(listView, &khCopy)
	}

	meCopy := *me
	if vInfo, ok := mapStyle[meCopy.VaiTroQuyenHan]; ok {
		meCopy.StyleLevel, meCopy.StyleTheme = vInfo.StyleLevel, vInfo.StyleTheme
	} else { meCopy.StyleLevel = 9 }
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" { meCopy.StyleLevel = 0 }

	if len(listVaiTro) == 0 {
		listVaiTro = []core.VaiTroInfo{
			{MaVaiTro: "quan_tri_he_thong", TenVaiTro: "Quản trị hệ thống", StyleLevel: 0, StyleTheme: 9},
			{MaVaiTro: "quan_tri_vien_he_thong", TenVaiTro: "Quản trị viên hệ thống", StyleLevel: 1, StyleTheme: 4},
			{MaVaiTro: "quan_tri_it_he_thong", TenVaiTro: "Quản trị IT hệ thống", StyleLevel: 2, StyleTheme: 7},
			{MaVaiTro: "quan_tri_cua_hang", TenVaiTro: "Quản trị cửa hàng", StyleLevel: 3, StyleTheme: 5},
		}
	}

	c.HTML(http.StatusOK, "master_thanh_vien", gin.H{
		"TieuDe": "Thành Viên",
		"NhanVien": &meCopy, 
		"DanhSach": listView, 
		"DanhSachVaiTro": listVaiTro, 
	})
}

func API_LuuThanhVienMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID") 
	
	// Ở đây đã chặn cứng: Chỉ có Level 0, 1, 2 mới lọt được qua khe cửa này
	myLevel := core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị Lõi 99K mới được sửa hồ sơ tại đây!"})
		return
	}

	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	admin, okAdmin := core.LayKhachHang(shopID, userID)
	if !okAdmin || admin.MaPinHash == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng thiết lập Mã PIN trước."})
		return
	}
	if !config.KiemTraMatKhau(pinXacNhan, admin.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN xác nhận không chính xác!"})
		return
	}

	maKH := c.PostForm("ma_khach_hang")
	kh, ok := core.LayKhachHang(shopID, maKH)
	if !ok { return }

	targetLevel := core.LayCapBacVaiTro(shopID, kh.MaKhachHang, kh.VaiTroQuyenHan)
	newRole := c.PostForm("vai_tro")

	// Bảo vệ ID 001 tuyệt đối
	if maKH == "0000000000000000001" && userID != "0000000000000000001" {
		c.JSON(200, gin.H{"status": "error", "msg": "BẢO MẬT: Không ai được chạm vào hồ sơ Người Sáng Lập!"})
		return
	}
	if newRole == "quan_tri_he_thong" && maKH != "0000000000000000001" {
		c.JSON(200, gin.H{"status": "error", "msg": "BẢO MẬT: Chỉ có duy nhất 1 vị trí Quản trị hệ thống (ID 001)!"})
		return
	}

	// [ĐÃ SỬA]: MIỄN TRỪ KIỂM TRA CẤP BẬC NẾU TÀI KHOẢN BỊ SỬA LÀ BOT (ID 000)
	if maKH != "0000000000000000000" {
		// Dành cho người dùng thật: Luật "Không được sửa người ngang hàng hoặc cấp cao hơn"
		if userID != maKH && myLevel >= targetLevel {
			c.JSON(200, gin.H{"status": "error", "msg": "Lỗi: Bạn không có quyền chỉnh sửa cấp trên hoặc người ngang hàng!"})
			return
		}
	}

	// Chặn việc tự phong tước hoặc nâng quyền người khác cao hơn mình
	if newRole != "" && newRole != kh.VaiTroQuyenHan {
		newLevel := core.LayCapBacVaiTro(shopID, "", newRole)
		if newLevel <= myLevel && myLevel != 0 {
			c.JSON(200, gin.H{"status": "error", "msg": "Lỗi: Không thể bổ nhiệm chức vụ ngang bằng hoặc cao hơn quyền hạn của bạn!"})
			return
		}
	}

	if maKH == userID && c.PostForm("trang_thai") == "0" {
		c.JSON(200, gin.H{"status": "error", "msg": "Hệ thống bảo vệ: Không thể tự khóa tài khoản chính mình!"})
		return
	}

	core.KhoaHeThong.Lock()
	if newRole != "" {
		kh.VaiTroQuyenHan = newRole
		if cv := strings.TrimSpace(c.PostForm("chuc_vu")); cv != "" { kh.ChucVu = cv 
		} else {
			kh.ChucVu = newRole 
			for _, v := range core.CacheDanhSachVaiTro[shopID] {
				if v.MaVaiTro == newRole { kh.ChucVu = v.TenVaiTro; break }
			}
		}
	}
	kh.TenKhachHang = strings.TrimSpace(c.PostForm("ten_khach_hang"))
	kh.DienThoai = strings.TrimSpace(c.PostForm("dien_thoai"))
	kh.NgaySinh = strings.TrimSpace(c.PostForm("ngay_sinh"))
	kh.DiaChi = strings.TrimSpace(c.PostForm("dia_chi"))
	kh.MaSoThue = strings.TrimSpace(c.PostForm("ma_so_thue"))
	kh.GhiChu = strings.TrimSpace(c.PostForm("ghi_chu"))
	kh.AnhDaiDien = strings.TrimSpace(c.PostForm("anh_dai_dien"))
	kh.NguonKhachHang = strings.TrimSpace(c.PostForm("nguon_khach_hang"))
	
	if gt := c.PostForm("gioi_tinh"); gt == "1" { kh.GioiTinh = 1 } else if gt == "0" { kh.GioiTinh = 0 } else { kh.GioiTinh = -1 }
	tt := c.PostForm("trang_thai")
	if tt == "1" { kh.TrangThai = 1 } else if tt == "-1" { kh.TrangThai = -1 } else { kh.TrangThai = 0 }

	kh.MangXaHoi.Zalo = strings.TrimSpace(c.PostForm("zalo"))
	kh.MangXaHoi.Facebook = strings.TrimSpace(c.PostForm("facebook"))
	kh.MangXaHoi.Tiktok = strings.TrimSpace(c.PostForm("tiktok"))

	if maKH != "0000000000000000000" {
		if passMoi := strings.TrimSpace(c.PostForm("mat_khau_moi")); passMoi != "" {
			hash, _ := cau_hinh.HashMatKhau(passMoi); kh.MatKhauHash = hash
			core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
		}
		if pinMoi := strings.TrimSpace(c.PostForm("pin_moi")); pinMoi != "" {
			hashPin, _ := cau_hinh.HashMatKhau(pinMoi); kh.MaPinHash = hashPin
			core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashPin)
		}
	}
	kh.NgayCapNhat = time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")
	kh.NguoiCapNhat = admin.TenDangNhap 
	core.KhoaHeThong.Unlock()

	ghi := core.ThemVaoHangCho; r := kh.DongTrongSheet; sh := "KHACH_HANG"
	ghi(shopID, sh, r, core.CotKH_TenKhachHang, kh.TenKhachHang); ghi(shopID, sh, r, core.CotKH_DienThoai, kh.DienThoai)
	ghi(shopID, sh, r, core.CotKH_NgaySinh, kh.NgaySinh); ghi(shopID, sh, r, core.CotKH_GioiTinh, kh.GioiTinh)
	ghi(shopID, sh, r, core.CotKH_DiaChi, kh.DiaChi); ghi(shopID, sh, r, core.CotKH_MaSoThue, kh.MaSoThue)
	ghi(shopID, sh, r, core.CotKH_GhiChu, kh.GhiChu); ghi(shopID, sh, r, core.CotKH_TrangThai, kh.TrangThai)
	ghi(shopID, sh, r, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan); ghi(shopID, sh, r, core.CotKH_ChucVu, kh.ChucVu)
	ghi(shopID, sh, r, core.CotKH_AnhDaiDien, kh.AnhDaiDien); ghi(shopID, sh, r, core.CotKH_NguonKhachHang, kh.NguonKhachHang) 
	ghi(shopID, sh, r, core.CotKH_NgayCapNhat, kh.NgayCapNhat); ghi(shopID, sh, r, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
	ghi(shopID, sh, r, core.CotKH_MangXaHoiJson, core.ToJSON(kh.MangXaHoi))

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thông tin thành công!"})
}

func API_GuiTinNhanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	if core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE")) > 2 {
		c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị Lõi 99K mới được phát sóng thông báo!"})
		return
	}

	tieuDe, noiDung := strings.TrimSpace(c.PostForm("tieu_de")), strings.TrimSpace(c.PostForm("noi_dung"))
	jsonIDs := c.PostForm("danh_sach_id") 
	if tieuDe == "" || noiDung == "" { c.JSON(200, gin.H{"status": "error", "msg": "Tiêu đề và Nội dung không được để trống!"}); return }

	var listMaKH []string
	if err := json.Unmarshal([]byte(jsonIDs), &listMaKH); err != nil || len(listMaKH) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Chưa chọn người nhận hợp lệ!"}); return
	}

	senderID := userID
	if c.PostForm("send_as_bot") == "1" {
		if bot, ok := core.LayKhachHang(shopID, "0000000000000000000"); ok {
			senderID = bot.MaKhachHang
		} else { 
            senderID = "0000000000000000000" 
        }
	}

	now := time.Now(); nowStr := now.In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")
	msgID := fmt.Sprintf("ALL_%d_%s", now.UnixNano(), senderID) 

	core.ThemMoiTinNhan(shopID, &core.TinNhan{
		MaTinNhan: msgID, LoaiTinNhan: "ALL", NguoiGuiID: senderID, NguoiNhanID: jsonIDs,       
		TieuDe: tieuDe, NoiDung: noiDung, NgayTao: nowStr,
	})

	c.JSON(200, gin.H{"status": "ok", "msg": fmt.Sprintf("Đã gửi thông báo thành công cho %d người!", len(listMaKH))})
}
