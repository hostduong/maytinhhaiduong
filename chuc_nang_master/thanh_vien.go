package chuc_nang_master

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

func TrangQuanLyThanhVienMaster(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if vaiTro != "quan_tri_he_thong" && vaiTro != "quan_tri_vien_he_thong" {
		c.Redirect(http.StatusFound, "/")
		return
	}
	
	me, _ := core.LayKhachHang(masterShopID, userID)
	listAll := core.LayDanhSachKhachHang(masterShopID)
	
	var listView []*core.KhachHang
	for _, kh := range listAll {
		khCopy := *kh 
		khCopy.Inbox = core.LayHopThuNguoiDung(masterShopID, khCopy.MaKhachHang, khCopy.VaiTroQuyenHan)
		listView = append(listView, &khCopy)
	}
	
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[masterShopID]
	core.KhoaHeThong.RUnlock()

	if len(listVaiTro) == 0 {
		listVaiTro = []core.VaiTroInfo{
			{MaVaiTro: "quan_tri_he_thong", TenVaiTro: "Quản trị hệ thống"},
			{MaVaiTro: "quan_tri_vien_he_thong", TenVaiTro: "Quản trị viên hệ thống"},
			{MaVaiTro: "quan_tri_vien", TenVaiTro: "Quản trị viên"},
			{MaVaiTro: "khach_hang", TenVaiTro: "Khách hàng"},
		}
	}

	c.HTML(http.StatusOK, "master_thanh_vien", gin.H{
		"TieuDe":         "Core Team",
		"NhanVien":       me,
		"DanhSach":       listView, 
		"DanhSachVaiTro": listVaiTro, 
	})
}

func API_LuuThanhVienMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID") 
	myRole := c.GetString("USER_ROLE")
	
	if myRole != "quan_tri_he_thong" && myRole != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền quản trị nhân sự!"})
		return
	}

	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	if pinXacNhan == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập mã PIN xác nhận!"})
		return
	}

	admin, okAdmin := core.LayKhachHang(shopID, userID)
	if !okAdmin || admin.MaPinHash == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn chưa thiết lập Mã PIN! Vui lòng cài đặt trước."})
		return
	}

	if !cau_hinh.KiemTraMatKhau(pinXacNhan, admin.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN xác nhận không chính xác!"})
		return
	}

	maKH := c.PostForm("ma_khach_hang")
	kh, ok := core.LayKhachHang(shopID, maKH)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thành viên này!"})
		return
	}

	isTargetRootLevel1 := (maKH == "0000000000000000001" || kh.VaiTroQuyenHan == "quan_tri_he_thong")
	isMeRootLevel1 := (userID == "0000000000000000001" || myRole == "quan_tri_he_thong")
	isTargetRootLevel2 := (kh.VaiTroQuyenHan == "quan_tri_vien_he_thong")
	newRole := c.PostForm("vai_tro")

	if isTargetRootLevel1 && !isMeRootLevel1 {
		c.JSON(200, gin.H{"status": "error", "msg": "BẢO MẬT TỐI CAO: Không ai có thể chỉnh sửa thông tin của Chủ tịch!"})
		return
	}
	if isTargetRootLevel1 && isMeRootLevel1 && newRole != "" && newRole != "quan_tri_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn là Người sáng lập, không thể tự giáng chức chính mình!"})
		return
	}
	if isTargetRootLevel2 && !isMeRootLevel1 && userID != maKH {
		c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị hệ thống (Cấp 1) mới sửa được Quản trị viên hệ thống (Cấp 2) khác!"})
		return
	}
	if (newRole == "quan_tri_he_thong" || newRole == "quan_tri_vien_he_thong") && !isMeRootLevel1 {
		if userID != maKH || (userID == maKH && kh.VaiTroQuyenHan != newRole) {
			c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị hệ thống mới có quyền bổ nhiệm chức vụ này!"})
			return
		}
	}

	trangThaiMoi := c.PostForm("trang_thai")
	if maKH == userID && trangThaiMoi == "0" {
		c.JSON(200, gin.H{"status": "error", "msg": "Hệ thống bảo vệ: Bạn không thể tự khóa tài khoản của chính mình!"})
		return
	}

	core.KhoaHeThong.Lock()
	if newRole != "" {
		kh.VaiTroQuyenHan = newRole
		chucVuTuY := strings.TrimSpace(c.PostForm("chuc_vu"))
		if chucVuTuY != "" {
			kh.ChucVu = chucVuTuY 
		} else {
			kh.ChucVu = newRole 
			for _, v := range core.CacheDanhSachVaiTro[shopID] {
				if v.MaVaiTro == newRole {
					kh.ChucVu = v.TenVaiTro
					break
				}
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
	
	gioiTinh := c.PostForm("gioi_tinh")
	if gioiTinh == "1" { kh.GioiTinh = 1 } else if gioiTinh == "0" { kh.GioiTinh = 0 } else { kh.GioiTinh = -1 }
	if trangThaiMoi == "1" { kh.TrangThai = 1 } else { kh.TrangThai = 0 }

	kh.MangXaHoi.Zalo = strings.TrimSpace(c.PostForm("zalo"))
	kh.MangXaHoi.Facebook = strings.TrimSpace(c.PostForm("facebook"))
	kh.MangXaHoi.Tiktok = strings.TrimSpace(c.PostForm("tiktok"))

	// [MỚI BỔ SUNG BẢO MẬT]: Chỉ cho phép đổi Pass/PIN nếu KHÔNG PHẢI là Bot Hệ Thống
	if maKH != "0000000000000000000" {
		passMoi := strings.TrimSpace(c.PostForm("mat_khau_moi"))
		if passMoi != "" {
			hash, _ := cau_hinh.HashMatKhau(passMoi)
			kh.MatKhauHash = hash
			core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MatKhauHash, hash)
		}

		pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
		if pinMoi != "" {
			hashPin, _ := cau_hinh.HashMatKhau(pinMoi)
			kh.MaPinHash = hashPin
			core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_MaPinHash, hashPin)
		}
	}

	loc := time.FixedZone("ICT", 7*3600)
	kh.NgayCapNhat = time.Now().In(loc).Format("2006-01-02 15:04:05")
	
	// [ĐÃ SỬA]: Lưu thông tin User ID (Tên đăng nhập) của người thực hiện hành động này
	kh.NguoiCapNhat = admin.TenDangNhap 
	core.KhoaHeThong.Unlock()

	ghi := core.ThemVaoHangCho
	r := kh.DongTrongSheet
	sh := "KHACH_HANG"
	
	ghi(shopID, sh, r, core.CotKH_TenKhachHang, kh.TenKhachHang)
	ghi(shopID, sh, r, core.CotKH_DienThoai, kh.DienThoai)
	ghi(shopID, sh, r, core.CotKH_NgaySinh, kh.NgaySinh)
	ghi(shopID, sh, r, core.CotKH_GioiTinh, kh.GioiTinh)
	ghi(shopID, sh, r, core.CotKH_DiaChi, kh.DiaChi)
	ghi(shopID, sh, r, core.CotKH_MaSoThue, kh.MaSoThue)
	ghi(shopID, sh, r, core.CotKH_GhiChu, kh.GhiChu)
	ghi(shopID, sh, r, core.CotKH_TrangThai, kh.TrangThai)
	ghi(shopID, sh, r, core.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	ghi(shopID, sh, r, core.CotKH_ChucVu, kh.ChucVu)
	ghi(shopID, sh, r, core.CotKH_AnhDaiDien, kh.AnhDaiDien)         
	ghi(shopID, sh, r, core.CotKH_NguonKhachHang, kh.NguonKhachHang) 
	ghi(shopID, sh, r, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
	
	// [ĐÃ SỬA]: Ghi Người Cập Nhật xuống Sheet
	ghi(shopID, sh, r, core.CotKH_NguoiCapNhat, kh.NguoiCapNhat)
	
	jsonMXH := core.ToJSON(kh.MangXaHoi)
	ghi(shopID, sh, r, core.CotKH_MangXaHoiJson, jsonMXH)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thông tin thành công!"})
}

func API_GuiTinNhanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	myRole := c.GetString("USER_ROLE")
	
	if myRole != "quan_tri_he_thong" && myRole != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền gửi thông báo!"})
		return
	}

	tieuDe := strings.TrimSpace(c.PostForm("tieu_de"))
	noiDung := strings.TrimSpace(c.PostForm("noi_dung"))
	jsonIDs := c.PostForm("danh_sach_id")

	if tieuDe == "" || noiDung == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tiêu đề và Nội dung không được để trống!"})
		return
	}

	var listMaKH []string
	if err := json.Unmarshal([]byte(jsonIDs), &listMaKH); err != nil || len(listMaKH) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Chưa chọn người nhận hợp lệ!"})
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

	count := 0
	for _, maKH := range listMaKH {
		if _, ok := core.LayKhachHang(shopID, maKH); ok {
			msgID := fmt.Sprintf("MSG_%d_%s", time.Now().UnixNano(), maKH) 

			newMsg := &core.TinNhan{
				MaTinNhan:      msgID,
				LoaiTinNhan:    "SYSTEM",
				NguoiGuiID:     userID,         
				NguoiNhanID:    maKH,           
				TieuDe:         tieuDe,
				NoiDung:        noiDung,
				NgayTao:        nowStr,
				TenNguoiGui:    tenNguoiGui,
				ChucVuNguoiGui: chucVuNguoiGui,
			}
			core.ThemMoiTinNhan(shopID, newMsg)
			count++
		}
	}

	c.JSON(200, gin.H{"status": "ok", "msg": fmt.Sprintf("Đã gửi thông báo thành công cho %d người!", count)})
}
