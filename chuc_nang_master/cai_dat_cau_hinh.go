package chuc_nang_master

import (
	"net/http"
	"strconv"
	"strings"

	"app/core"
	"github.com/gin-gonic/gin"
)

// =============================================================
// 1. TRANG CÀI ĐẶT CẤU HÌNH MASTER
// =============================================================
func TrangCaiDatCauHinhMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	kh, _ := core.LayKhachHang(shopID, userID)

	c.HTML(http.StatusOK, "master_cai_dat_cau_hinh", gin.H{
		"TieuDe":       "Cấu Hình Hệ Thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,

		// Nạp đủ 4 mảng Master Data
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
	vaiTro := c.GetString("USER_ROLE")
	
	if vaiTro != "admin_root" && vaiTro != "admin" {
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
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
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
			if maSP[i] >= '0' && maSP[i] <= '9' {
				batDauSo = i
			} else { break }
		}
		if batDauSo < ketThucSo {
			prefix := strings.ToUpper(maSP[0:batDauSo])
			soStr := maSP[batDauSo:ketThucSo]
			if so, err := strconv.Atoi(soStr); err == nil {
				if so > mapMaxSlot[prefix] { mapMaxSlot[prefix] = so }
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

	msg := "Đã đồng bộ xong."
	if countUpdate > 0 { msg = "Đã cập nhật lại Slot cho " + strconv.Itoa(countUpdate) + " danh mục." }
	c.JSON(200, gin.H{"status": "ok", "msg": msg})
}

// =============================================================
// 3. API THƯƠNG HIỆU
// =============================================================
func API_LuuThuongHieuMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
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

// =============================================================
// 4. API BIÊN LỢI NHUẬN
// =============================================================
func API_LuuBienLoiNhuanMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thao tác!"})
		return
	}

	khungGia, _ := strconv.ParseFloat(strings.ReplaceAll(c.PostForm("khung_gia_nhap"), ".", ""), 64)
	loiNhuan, _ := strconv.ParseFloat(c.PostForm("bien_loi_nhuan"), 64)
	trangThai := 0; if c.PostForm("trang_thai") == "on" { trangThai = 1 }
	isNew := c.PostForm("is_new") == "true"
	dongCu, _ := strconv.Atoi(c.PostForm("dong_cu")) 

	if khungGia <= 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Khung giá phải lớn hơn 0!"})
		return
	}

	var targetRow int
	if isNew {
		targetRow = core.DongBatDau_BienLoiNhuan + len(core.LayDanhSachBienLoiNhuan(shopID))
		newBLN := &core.BienLoiNhuan{
			SpreadsheetID:  shopID,
			DongTrongSheet: targetRow,
			KhungGiaNhap:   khungGia,
			BienLoiNhuan:   loiNhuan,
			TrangThai:      trangThai,
		}
		core.ThemBienLoiNhuanVaoRam(newBLN) 
	} else {
		targetRow = dongCu
		core.SuaBienLoiNhuanTrongRam(shopID, dongCu, khungGia, loiNhuan, trangThai) 
	}

	ghi := core.ThemVaoHangCho
	ghi(shopID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_KhungGiaNhap, khungGia)
	ghi(shopID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_BienLoiNhuan, loiNhuan)
	ghi(shopID, "BIEN_LOI_NHUAN", targetRow, core.CotBLN_TrangThai, trangThai)

	c.JSON(200, gin.H{"status": "ok", "msg": "Lưu Khung giá thành công!"})
}

/// ==============================================================================
// PHẦN 5: PHÂN QUYỀN VÀ BẢO MẬT
// ==============================================================================
const (
	TenSheetPhanQuyen = "PHAN_QUYEN"
	DongBatDau_PhanQuyen = 2

	CotPQ_Id           = 0
	CotPQ_TenVaiTro    = 1
	CotPQ_CapBac       = 2
	CotPQ_QuyenQuanTri = 3
	CotPQ_QuyenHangHoa = 4
	CotPQ_QuyenDonHang = 5
	CotPQ_TrangThai    = 6
)

type PhanQuyen struct {
	Id           string `json:"id"`
	TenVaiTro    string `json:"ten_vai_tro"`
	CapBac       int    `json:"cap_bac"`
	QuyenQuanTri string `json:"quyen_quan_tri"`
	QuyenHangHoa string `json:"quyen_hang_hoa"`
	QuyenDonHang string `json:"quyen_don_hang"`
	TrangThai    int    `json:"trang_thai"`
}

type VaiTroInfo struct {
	MaVaiTro   string
	TenVaiTro  string
	CapBac     int
	TrangThai  int
	StyleLevel string
	StyleTheme string
}

var (
	CachePhanQuyen      = make(map[string]map[string]*PhanQuyen)
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)
	lockPhanQuyen       sync.RWMutex
)

func NapPhanQuyen(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, TenSheetPhanQuyen)
	if err != nil { return }

	mapPQ := make(map[string]*PhanQuyen)
	listVT := make([]VaiTroInfo, 0)

	for i, r := range raw {
		if i < DongBatDau_PhanQuyen-1 { continue }
		id := LayString(r, CotPQ_Id)
		if id == "" { continue }

		pq := &PhanQuyen{
			Id:           id,
			TenVaiTro:    LayString(r, CotPQ_TenVaiTro),
			CapBac:       LayInt(r, CotPQ_CapBac),
			QuyenQuanTri: LayString(r, CotPQ_QuyenQuanTri),
			QuyenHangHoa: LayString(r, CotPQ_QuyenHangHoa),
			QuyenDonHang: LayString(r, CotPQ_QuyenDonHang),
			TrangThai:    LayInt(r, CotPQ_TrangThai),
		}
		mapPQ[id] = pq
		
		styleLvl := "bg-slate-100 text-slate-600"
		styleTheme := "border-slate-200"
		if pq.CapBac >= 90 {
			styleLvl = "bg-red-100 text-red-700"
			styleTheme = "border-red-200"
		} else if pq.CapBac >= 50 {
			styleLvl = "bg-blue-100 text-blue-700"
			styleTheme = "border-blue-200"
		} else if pq.CapBac >= 10 {
			styleLvl = "bg-green-100 text-green-700"
			styleTheme = "border-green-200"
		}

		listVT = append(listVT, VaiTroInfo{
			MaVaiTro:   id,
			TenVaiTro:  pq.TenVaiTro,
			CapBac:     pq.CapBac,
			TrangThai:  pq.TrangThai,
			StyleLevel: styleLvl,
			StyleTheme: styleTheme,
		})
	}

	sort.Slice(listVT, func(i, j int) bool {
		return listVT[i].CapBac > listVT[j].CapBac
	})

	lockPhanQuyen.Lock()
	CachePhanQuyen[shopID] = mapPQ
	CacheDanhSachVaiTro[shopID] = listVT
	lockPhanQuyen.Unlock()
}

func KiemTraQuyen(shopID, idVaiTro, keyQuyen string) bool {
	if idVaiTro == "admin_root" { return true }
	
	lockPhanQuyen.RLock()
	defer lockPhanQuyen.RUnlock()

	mapPQ, okShop := CachePhanQuyen[shopID]
	if !okShop { return false }

	pq, okPQ := mapPQ[idVaiTro]
	if !okPQ || pq.TrangThai != 1 { return false }

	switch keyQuyen {
	case "quan_tri": return pq.QuyenQuanTri == "xem_sua"
	case "hang_hoa": return pq.QuyenHangHoa == "xem_sua" || pq.QuyenHangHoa == "xem"
	case "sua_hang_hoa": return pq.QuyenHangHoa == "xem_sua"
	case "don_hang": return pq.QuyenDonHang == "xem_sua" || pq.QuyenDonHang == "xem"
	case "sua_don_hang": return pq.QuyenDonHang == "xem_sua"
	}
	return false
}

// Khôi phục lại hàm LayCapBacVaiTro với 3 tham số để tránh lỗi ở file thanh_vien.go
func LayCapBacVaiTro(shopID, userID, idVaiTro string) int {
	if idVaiTro == "admin_root" { return 9999 }
	
	lockPhanQuyen.RLock()
	defer lockPhanQuyen.RUnlock()

	mapPQ, okShop := CachePhanQuyen[shopID]
	if !okShop { return 0 }

	pq, okPQ := mapPQ[idVaiTro]
	if !okPQ { return 0 }
	return pq.CapBac
}
