package chuc_nang_admin

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

// ==============================================================
// 1. TRANG HIỂN THỊ DANH SÁCH
// ==============================================================
func TrangQuanLyThanhVien(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	me, _ := core.LayKhachHang(shopID, userID)
	listAll := core.LayDanhSachKhachHang(shopID)
	
	core.KhoaHeThong.RLock()
	listVaiTro := core.CacheDanhSachVaiTro[shopID]
	core.KhoaHeThong.RUnlock()

	if len(listVaiTro) == 0 {
		listVaiTro = []core.VaiTroInfo{
			{MaVaiTro: "quan_tri_vien_he_thong", TenVaiTro: "Quản trị hệ thống (Chưa đồng bộ)"},
			{MaVaiTro: "quan_tri_vien", TenVaiTro: "Quản trị viên (Chưa đồng bộ)"},
			{MaVaiTro: "khach_hang", TenVaiTro: "Khách hàng (Chưa đồng bộ)"},
		}
	}

	c.HTML(http.StatusOK, "quan_tri_thanh_vien", gin.H{
		"TieuDe":         "Quản lý thành viên",
		"NhanVien":       me,
		"DanhSach":       listAll,
		"DanhSachVaiTro": listVaiTro, 
	})
}

// ==============================================================
// 2. API: CẬP NHẬT FULL THÔNG TIN THÀNH VIÊN
// ==============================================================
func API_Admin_LuuThanhVien(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID") 
	myRole := c.GetString("USER_ROLE")
	
	if myRole != "quan_tri_vien_he_thong" && myRole != "quan_tri_vien" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền quản trị nhân sự!"})
		return
	}

	// =====================================================
	// [BẢO MẬT KÉP]: KIỂM TRA MÃ PIN CỦA ADMIN ĐANG THAO TÁC
	// =====================================================
	pinXacNhan := strings.TrimSpace(c.PostForm("pin_xac_nhan"))
	if pinXacNhan == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng nhập mã PIN xác nhận!"})
		return
	}
	if !cau_hinh.KiemTraMaPin(pinXacNhan) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN xác nhận không hợp lệ (Phải đúng 8 số)!"})
		return
	}

	admin, okAdmin := core.LayKhachHang(shopID, userID)
	if !okAdmin {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi phiên đăng nhập!"})
		return
	}

	if admin.MaPinHash == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn chưa thiết lập Mã PIN! Vui lòng vào Hồ sơ cá nhân để cài đặt Mã PIN trước."})
		return
	}

	// Gọi hàm KiemTraMatKhau giống hệt bên quen_mat_khau.go
	if !cau_hinh.KiemTraMatKhau(pinXacNhan, admin.MaPinHash) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN xác nhận không chính xác!"})
		return
	}
	// =====================================================

	maKH := c.PostForm("ma_khach_hang")
	kh, ok := core.LayKhachHang(shopID, maKH)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thành viên này!"})
		return
	}

	if kh.VaiTroQuyenHan == "quan_tri_vien_he_thong" && myRole != "quan_tri_vien_he_thong" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không thể chỉnh sửa thông tin của Quản trị viên hệ thống!"})
		return
	}

	trangThaiMoi := c.PostForm("trang_thai")
	if maKH == userID && trangThaiMoi == "0" {
		c.JSON(200, gin.H{"status": "error", "msg": "Hệ thống bảo vệ: Bạn không thể tự khóa tài khoản của chính mình!"})
		return
	}

	core.KhoaHeThong.Lock()
	
	newRole := c.PostForm("vai_tro")
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
	
	gioiTinh := c.PostForm("gioi_tinh")
	if gioiTinh == "1" { kh.GioiTinh = 1 } else if gioiTinh == "0" { kh.GioiTinh = 0 } else { kh.GioiTinh = -1 }

	if trangThaiMoi == "1" { kh.TrangThai = 1 } else { kh.TrangThai = 0 }

	kh.MangXaHoi.Zalo = strings.TrimSpace(c.PostForm("zalo"))
	kh.MangXaHoi.Facebook = strings.TrimSpace(c.PostForm("facebook"))
	kh.MangXaHoi.Tiktok = strings.TrimSpace(c.PostForm("tiktok"))

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

	kh.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")
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
	ghi(shopID, sh, r, core.CotKH_NgayCapNhat, kh.NgayCapNhat)
	
	jsonMXH := core.ToJSON(kh.MangXaHoi)
	ghi(shopID, sh, r, core.CotKH_MangXaHoiJson, jsonMXH)

	c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thông tin thành công!"})
}

// ==============================================================
// 3. API: GỬI TIN NHẮN (BULK & SINGLE)
// ==============================================================
func API_Admin_GuiTinNhan(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	myRole := c.GetString("USER_ROLE")
	
	if myRole != "quan_tri_vien_he_thong" && myRole != "quan_tri_vien" {
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

	nowStr := time.Now().Format("2006-01-02 15:04:05")
	msgID := fmt.Sprintf("MSG_%d", time.Now().Unix()) 

	count := 0
	core.KhoaHeThong.Lock()
	for _, maKH := range listMaKH {
		if kh, ok := core.CacheMapKhachHang[core.TaoCompositeKey(shopID, maKH)]; ok {
			newMsg := core.MessageInfo{
				ID:      msgID,
				TieuDe:  tieuDe,
				NoiDung: noiDung,
				DaDoc:   false,
				NgayTao: nowStr,
			}
			kh.Inbox = append(kh.Inbox, newMsg)
			
			jsonInbox := core.ToJSON(kh.Inbox)
			core.ThemVaoHangCho(shopID, "KHACH_HANG", kh.DongTrongSheet, core.CotKH_InboxJson, jsonInbox)
			count++
		}
	}
	core.KhoaHeThong.Unlock()

	c.JSON(200, gin.H{"status": "ok", "msg": fmt.Sprintf("Đã gửi thông báo thành công cho %d người!", count)})
}
