package chuc_nang_master

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// =============================================================
// 1. TRANG CÀI ĐẶT CẤU HÌNH MASTER
// =============================================================
func TrangCaiDatCauHinhMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// Áp dụng chuẩn bảo vệ của bạn
	myLevel := core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.Redirect(http.StatusFound, "/")
		return
	}

	kh, _ := core.LayKhachHang(shopID, userID)

	c.HTML(http.StatusOK, "master_cai_dat_cau_hinh", gin.H{
		"TieuDe":       "Cấu Hình Hệ Thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,

		"ListDanhMuc":    core.LayDanhSachDanhMuc(shopID),
		"ListThuongHieu": core.LayDanhSachThuongHieu(shopID),
		"ListBLN":        core.LayDanhSachBienLoiNhuan(shopID),
		"ListNCC":        core.LayDanhSachNhaCungCap(shopID),
	})
}

// =============================================================
// 2. API DANH MỤC
// =============================================================
func API_LuuDanhMucMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	// Áp dụng chuẩn bảo vệ của bạn
	myLevel := core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maDM := strings.TrimSpace(c.PostForm("ma_danh_muc"))
	tenDM := strings.TrimSpace(c.PostForm("ten_danh_muc"))
	dmMe := strings.TrimSpace(c.PostForm("danh_muc_me"))
	thueVAT, _ := strconv.ParseFloat(c.PostForm("thue_vat"), 64)
	loiNhuan, _ := strconv.ParseFloat(c.PostForm("loi_nhuan"), 64)
	trangThai := 0; if c.PostForm("trang_thai") == "on" { trangThai = 1 }
	isNew := c.PostForm("is_new") == "true"

	if maDM == "" || tenDM == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên danh mục không được để trống!"})
		return
	}

	var targetRow int

	if isNew {
		if _, ok := core.LayChiTietDanhMuc(shopID, maDM); ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã danh mục này đã tồn tại!"})
			return
		}
		
		targetRow = core.DongBatDau_DanhMuc + len(core.LayDanhSachDanhMuc(shopID))
		
		newDM := &core.DanhMuc{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			MaDanhMuc:      strings.ToUpper(maDM),
			TenDanhMuc:     tenDM,
			DanhMucMe:      dmMe,
			ThueVAT:        thueVAT,
			LoiNhuan:       loiNhuan,
			Slot:           0,
			TrangThai:      trangThai,
		}
		core.ThemDanhMucVaoRam(newDM) 
	} else {
		found, ok := core.LayChiTietDanhMuc(shopID, maDM)
		if !ok {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy danh mục để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet
		
		func() {
			core.KhoaHeThong.Lock()
			defer core.KhoaHeThong.Unlock()
			found.TenDanhMuc = tenDM
			found.DanhMucMe = dmMe
			found.ThueVAT = thueVAT
			found.LoiNhuan = loiNhuan
			found.TrangThai = trangThai
		}()
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_MaDanhMuc, strings.ToUpper(maDM))
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_TenDanhMuc, tenDM)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_DanhMucMe, dmMe)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_ThueVAT, thueVAT)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_LoiNhuan, loiNhuan)
	ghi(shopID, "DANH_MUC", targetRow, core.CotDM_TrangThai, trangThai)
	
	if isNew { ghi(shopID, "DANH_MUC", targetRow, core.CotDM_Slot, 0) }

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Danh mục thành công!"})
}

func API_DongBoSlotDanhMucMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	// Áp dụng chuẩn bảo vệ của bạn
	myLevel := core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	listSP := core.LayDanhSachSanPhamMayTinh(shopID)
	mapMaxSlot := make(map[string]int) 
	
	for _, sp := range listSP {
		maSP := strings.TrimSpace(sp.MaSanPham)
		ketThucSo := len(maSP)
		batDauSo := ketThucSo
		
		for i := len(maSP) - 1; i >= 0; i-- {
			char := maSP[i]
			if char >= '0' && char <= '9' {
				batDauSo = i
			} else {
				break
			}
		}

		if batDauSo < ketThucSo {
			prefix := strings.ToUpper(maSP[0:batDauSo])
			soStr := maSP[batDauSo:ketThucSo]
			so, err := strconv.Atoi(soStr)
			if err == nil {
				if so > mapMaxSlot[prefix] {
					mapMaxSlot[prefix] = so
				}
			}
		}
	}

	listDM := core.LayDanhSachDanhMuc(shopID)

	core.KhoaHeThong.Lock()
	countUpdate := 0
	
	for _, dm := range listDM {
		maxThucTe, coDuLieu := mapMaxSlot[dm.MaDanhMuc]
		if coDuLieu && maxThucTe > dm.Slot {
			dm.Slot = maxThucTe
			core.ThemVaoHangCho(shopID, "DANH_MUC", dm.DongTrongSheet, core.CotDM_Slot, dm.Slot)
			countUpdate++
		}
	}
	core.KhoaHeThong.Unlock()

	msg := "Đã đồng bộ xong. Các bộ đếm đều chính xác."
	if countUpdate > 0 {
		msg = "Đã cập nhật lại Slot cho " + strconv.Itoa(countUpdate) + " danh mục."
	}

	c.JSON(200, gin.H{"status": "ok", "msg": msg})
}

// =============================================================
// 3. API THƯƠNG HIỆU
// =============================================================
func API_LuuThuongHieuMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	
	// Áp dụng chuẩn bảo vệ của bạn
	myLevel := core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	maTH := strings.TrimSpace(c.PostForm("ma_thuong_hieu"))
	tenTH := strings.TrimSpace(c.PostForm("ten_thuong_hieu"))
	logoUrl := strings.TrimSpace(c.PostForm("logo_url"))
	moTa := strings.TrimSpace(c.PostForm("mo_ta"))
	trangThai := 0; if c.PostForm("trang_thai") == "on" { trangThai = 1 }
	isNew := c.PostForm("is_new") == "true"

	if maTH == "" || tenTH == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã và Tên thương hiệu không được trống!"})
		return
	}

	var targetRow int

	if isNew {
		targetRow = core.DongBatDau_ThuongHieu + len(core.LayDanhSachThuongHieu(shopID))
		newTH := &core.ThuongHieu{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			MaThuongHieu:   strings.ToUpper(maTH),
			TenThuongHieu:  tenTH,
			LogoUrl:        logoUrl,
			MoTa:           moTa,
			TrangThai:      trangThai,
		}
		core.ThemThuongHieuVaoRam(newTH) 
	} else {
		var found *core.ThuongHieu
		for _, item := range core.LayDanhSachThuongHieu(shopID) {
			if item.MaThuongHieu == maTH { found = item; break }
		}
		if found == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy thương hiệu để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet
		func() {
			core.KhoaHeThong.Lock()
			defer core.KhoaHeThong.Unlock()
			found.TenThuongHieu = tenTH
			found.LogoUrl = logoUrl
			found.MoTa = moTa
			found.TrangThai = trangThai
		}()
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_MaThuongHieu, strings.ToUpper(maTH))
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_TenThuongHieu, tenTH)
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_LogoUrl, logoUrl)
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_MoTa, moTa)
	ghi(shopID, "THUONG_HIEU", targetRow, core.CotTH_TrangThai, trangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Thương hiệu thành công!"})
}

// ==============================================================================
// PHẦN 4: NHÀ CUNG CẤP
// ==============================================================================
const (
	TenSheetNhaCungCap    = "NHA_CUNG_CAP"
	DongBatDau_NhaCungCap = 2

	CotNCC_MaNhaCungCap       = 0  // A
	CotNCC_TenNhaCungCap      = 1  // B
	CotNCC_MaSoThue           = 2  // C
	CotNCC_DienThoai          = 3  // D
	CotNCC_Email              = 4  // E
	CotNCC_KhuVuc             = 5  // F
	CotNCC_DiaChi             = 6  // G
	CotNCC_NguoiLienHe        = 7  // H
	CotNCC_NganHang           = 8  // I
	CotNCC_NhomNhaCungCap     = 9  // J
	CotNCC_LoaiNhaCungCap     = 10 // K
	CotNCC_DieuKhoanThanhToan = 11 // L
	CotNCC_ChietKhauMacDinh   = 12 // M
	CotNCC_HanMucCongNo       = 13 // N
	CotNCC_CongNoDauKy        = 14 // O
	CotNCC_TongMua            = 15 // P
	CotNCC_NoCanTra           = 16 // Q
	CotNCC_ThongTinThemJson   = 17 // R
	CotNCC_TrangThai          = 18 // S
	CotNCC_GhiChu             = 19 // T
	CotNCC_NguoiTao           = 20 // U
	CotNCC_NgayTao            = 21 // V
	CotNCC_NgayCapNhat        = 22 // W
)

type NhaCungCap struct {
	SpreadsheetID      string  `json:"-"`
	DongTrongSheet     int     `json:"-"`
	MaNhaCungCap       string  `json:"ma_nha_cung_cap"`
	TenNhaCungCap      string  `json:"ten_nha_cung_cap"`
	MaSoThue           string  `json:"ma_so_thue"`
	DienThoai          string  `json:"dien_thoai"`
	Email              string  `json:"email"`
	KhuVuc             string  `json:"khu_vuc"`
	DiaChi             string  `json:"dia_chi"`
	NguoiLienHe        string  `json:"nguoi_lien_he"`
	NganHang           string  `json:"ngan_hang"`
	NhomNhaCungCap     string  `json:"nhom_nha_cung_cap"`
	LoaiNhaCungCap     string  `json:"loai_nha_cung_cap"`
	DieuKhoanThanhToan string  `json:"dieu_khoan_thanh_toan"`
	ChietKhauMacDinh   float64 `json:"chiet_khau_mac_dinh"`
	HanMucCongNo       float64 `json:"han_muc_cong_no"`
	CongNoDauKy        float64 `json:"cong_no_dau_ky"`
	TongMua            float64 `json:"tong_mua"`
	NoCanTra           float64 `json:"no_can_tra"`
	ThongTinThemJson   string  `json:"thong_tin_them_json"`
	TrangThai          int     `json:"trang_thai"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`
}

var (
	CacheNhaCungCap    = make(map[string][]*NhaCungCap)
	CacheMapNhaCungCap = make(map[string]*NhaCungCap)
)

func NapNhaCungCap(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, TenSheetNhaCungCap)
	if err != nil { return }
	
	list := []*NhaCungCap{}
	for i, r := range raw {
		if i < DongBatDau_NhaCungCap-1 { continue }
		maNCC := LayString(r, CotNCC_MaNhaCungCap)
		if maNCC == "" { continue }
		
		ncc := &NhaCungCap{
			SpreadsheetID:      shopID,
			DongTrongSheet:     i + 1,
			MaNhaCungCap:       maNCC,
			TenNhaCungCap:      LayString(r, CotNCC_TenNhaCungCap),
			MaSoThue:           LayString(r, CotNCC_MaSoThue),
			DienThoai:          LayString(r, CotNCC_DienThoai),
			Email:              LayString(r, CotNCC_Email),
			KhuVuc:             LayString(r, CotNCC_KhuVuc),
			DiaChi:             LayString(r, CotNCC_DiaChi),
			NguoiLienHe:        LayString(r, CotNCC_NguoiLienHe),
			NganHang:           LayString(r, CotNCC_NganHang),
			NhomNhaCungCap:     LayString(r, CotNCC_NhomNhaCungCap),
			LoaiNhaCungCap:     LayString(r, CotNCC_LoaiNhaCungCap),
			DieuKhoanThanhToan: LayString(r, CotNCC_DieuKhoanThanhToan),
			ChietKhauMacDinh:   LayFloat(r, CotNCC_ChietKhauMacDinh),
			HanMucCongNo:       LayFloat(r, CotNCC_HanMucCongNo),
			CongNoDauKy:        LayFloat(r, CotNCC_CongNoDauKy),
			TongMua:            LayFloat(r, CotNCC_TongMua),
			NoCanTra:           LayFloat(r, CotNCC_NoCanTra),
			ThongTinThemJson:   LayString(r, CotNCC_ThongTinThemJson),
			TrangThai:          LayInt(r, CotNCC_TrangThai),
			GhiChu:             LayString(r, CotNCC_GhiChu),
			NguoiTao:           LayString(r, CotNCC_NguoiTao),
			NgayTao:            LayString(r, CotNCC_NgayTao),
			NgayCapNhat:        LayString(r, CotNCC_NgayCapNhat),
		}
		list = append(list, ncc)
		key := TaoCompositeKey(shopID, maNCC)
		CacheMapNhaCungCap[key] = ncc
	}
	KhoaHeThong.Lock()
	CacheNhaCungCap[shopID] = list
	KhoaHeThong.Unlock()
}

func LayDanhSachNhaCungCap(shopID string) []*NhaCungCap {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	if list, ok := CacheNhaCungCap[shopID]; ok {
		return list
	}
	return []*NhaCungCap{}
}

func TaoMaNhaCungCapMoi(shopID string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	prefix := "NCC"
	maxNum := 0
	for _, ncc := range CacheNhaCungCap[shopID] {
		if strings.HasPrefix(ncc.MaNhaCungCap, prefix) {
			numStr := strings.TrimPrefix(ncc.MaNhaCungCap, prefix)
			if num, err := strconv.Atoi(numStr); err == nil {
				if num > maxNum { maxNum = num }
			}
		}
	}
	return fmt.Sprintf("%s%03d", prefix, maxNum+1)
}

// =============================================================
// 5. API NHÀ CUNG CẤP
// =============================================================
func API_LuuNhaCungCapMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	myLevel := core.LayCapBacVaiTro(shopID, userID, c.GetString("USER_ROLE"))
	if myLevel > 2 {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	// Lấy toàn bộ 23 trường
	maNCC := strings.TrimSpace(c.PostForm("ma_nha_cung_cap"))
	tenNCC := strings.TrimSpace(c.PostForm("ten_nha_cung_cap"))
	mst := strings.TrimSpace(c.PostForm("ma_so_thue"))
	sdt := strings.TrimSpace(c.PostForm("dien_thoai"))
	email := strings.TrimSpace(c.PostForm("email"))
	khuVuc := strings.TrimSpace(c.PostForm("khu_vuc"))
	diaChi := strings.TrimSpace(c.PostForm("dia_chi"))
	nguoiLh := strings.TrimSpace(c.PostForm("nguoi_lien_he"))
	nganHang := strings.TrimSpace(c.PostForm("ngan_hang"))
	nhomNCC := strings.TrimSpace(c.PostForm("nhom_nha_cung_cap"))
	loaiNCC := strings.TrimSpace(c.PostForm("loai_nha_cung_cap"))
	dkThanhToan := strings.TrimSpace(c.PostForm("dieu_khoan_thanh_toan"))

	ckMacDinh, _ := strconv.ParseFloat(c.PostForm("chiet_khau_mac_dinh"), 64)
	hanMuc, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("han_muc_cong_no"), ".", ""), 64)
	noDauKy, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("cong_no_dau_ky"), ".", ""), 64)

	thongTinJson := strings.TrimSpace(c.PostForm("thong_tin_them_json"))
	if thongTinJson == "" { thongTinJson = "{}" }
	
	ghiChu := strings.TrimSpace(c.PostForm("ghi_chu"))
	trangThai := 0
	if c.PostForm("trang_thai") == "on" || c.PostForm("trang_thai") == "1" {
		trangThai = 1
	}
	isNew := c.PostForm("is_new") == "true"

	if tenNCC == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên nhà cung cấp không được để trống!"})
		return
	}

	var targetRow int
	nowStr := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")

	if isNew {
		maMoi := maNCC
		if maMoi == "" { maMoi = core.TaoMaNhaCungCapMoi(shopID) }
		targetRow = core.DongBatDau_NhaCungCap + len(core.LayDanhSachNhaCungCap(shopID))

		newNCC := &core.NhaCungCap{
			SpreadsheetID:      shopID,
			DongTrongSheet:     targetRow,
			MaNhaCungCap:       maMoi,
			TenNhaCungCap:      tenNCC,
			MaSoThue:           mst,
			DienThoai:          sdt,
			Email:              email,
			KhuVuc:             khuVuc,
			DiaChi:             diaChi,
			NguoiLienHe:        nguoiLh,
			NganHang:           nganHang,
			NhomNhaCungCap:     nhomNCC,
			LoaiNhaCungCap:     loaiNCC,
			DieuKhoanThanhToan: dkThanhToan,
			ChietKhauMacDinh:   ckMacDinh,
			HanMucCongNo:       hanMuc,
			CongNoDauKy:        noDauKy,
			TongMua:            0,       // Mới tạo thì tổng mua bằng 0
			NoCanTra:           noDauKy, // Nợ cần trả ban đầu = nợ đầu kỳ
			ThongTinThemJson:   thongTinJson,
			TrangThai:          trangThai,
			GhiChu:             ghiChu,
			NguoiTao:           userID,
			NgayTao:            nowStr,
			NgayCapNhat:        nowStr,
		}
		
		core.KhoaHeThong.Lock()
		core.CacheNhaCungCap[shopID] = append(core.CacheNhaCungCap[shopID], newNCC)
		core.KhoaHeThong.Unlock()
		
		maNCC = maMoi

		// Ghi các cột đặc thù lúc tạo mới
		core.ThemVaoHangCho(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TongMua, 0)
		core.ThemVaoHangCho(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NoCanTra, noDauKy)
		core.ThemVaoHangCho(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NguoiTao, userID)
		core.ThemVaoHangCho(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NgayTao, nowStr)

	} else {
		list := core.LayDanhSachNhaCungCap(shopID)
		var found *core.NhaCungCap
		for _, item := range list {
			if item.MaNhaCungCap == maNCC { found = item; break }
		}

		if found == nil {
			c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy Nhà cung cấp để sửa!"})
			return
		}
		targetRow = found.DongTrongSheet

		// Logic Kế toán: Nếu sửa Nợ đầu kỳ, phải bù trừ vào Nợ hiện tại
		chenhLechNoDauKy := noDauKy - found.CongNoDauKy
		noCanTraMoi := found.NoCanTra + chenhLechNoDauKy

		core.KhoaHeThong.Lock()
		found.TenNhaCungCap = tenNCC
		found.MaSoThue = mst
		found.DienThoai = sdt
		found.Email = email
		found.KhuVuc = khuVuc
		found.DiaChi = diaChi
		found.NguoiLienHe = nguoiLh
		found.NganHang = nganHang
		found.NhomNhaCungCap = nhomNCC
		found.LoaiNhaCungCap = loaiNCC
		found.DieuKhoanThanhToan = dkThanhToan
		found.ChietKhauMacDinh = ckMacDinh
		found.HanMucCongNo = hanMuc
		found.CongNoDauKy = noDauKy
		found.NoCanTra = noCanTraMoi
		found.ThongTinThemJson = thongTinJson
		found.TrangThai = trangThai
		found.GhiChu = ghiChu
		found.NgayCapNhat = nowStr
		core.KhoaHeThong.Unlock()

		// Cập nhật lại cột Nợ cần trả do thay đổi nợ đầu kỳ
		core.ThemVaoHangCho(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NoCanTra, noCanTraMoi)
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_MaNhaCungCap, maNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TenNhaCungCap, tenNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_MaSoThue, mst)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DienThoai, sdt)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_Email, email)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_KhuVuc, khuVuc)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DiaChi, diaChi)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NguoiLienHe, nguoiLh)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NganHang, nganHang)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NhomNhaCungCap, nhomNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_LoaiNhaCungCap, loaiNCC)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_DieuKhoanThanhToan, dkThanhToan)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_ChietKhauMacDinh, ckMacDinh)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_HanMucCongNo, hanMuc)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_CongNoDauKy, noDauKy)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_ThongTinThemJson, thongTinJson)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_TrangThai, trangThai)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_GhiChu, ghiChu)
	ghi(shopID, core.TenSheetNhaCungCap, targetRow, core.CotNCC_NgayCapNhat, nowStr)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Nhà cung cấp thành công!"})
}
