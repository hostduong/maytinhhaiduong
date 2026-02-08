package chuc_nang

import (
	"fmt" // [THÊM] Để in log debug
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// TrangQuanLySanPham : Hiển thị danh sách
func TrangQuanLySanPham(c *gin.Context) {
	// [DEBUG] In ra console để biết hàm này ĐÃ ĐƯỢC GỌI
	fmt.Println(">>> [DEBUG] Đang truy cập TrangQuanLySanPham...")

	userID := c.GetString("USER_ID")
	
	// Check user
	kh, found := core.LayKhachHang(userID)
	if !found || kh == nil {
		fmt.Println(">>> [DEBUG] User không hợp lệ -> Redirect Login")
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Lấy dữ liệu
	listSP := core.LayDanhSachSanPham()
	
	// [DEBUG] Kiểm tra số lượng sản phẩm load được
	fmt.Printf(">>> [DEBUG] Load được %d sản phẩm.\n", len(listSP))

	listDM := core.LayDanhSachDanhMuc()
	listTH := core.LayDanhSachThuongHieu()

	// [SỬA TÊN TEMPLATE] Gọi đúng "view_quan_tri_san_pham"
	c.HTML(http.StatusOK, "view_quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       listSP,
		"ListDanhMuc":    listDM, 
		"ListThuongHieu": listTH,
	})
}

// API_LuuSanPham : Xử lý Thêm/Sửa
func API_LuuSanPham(c *gin.Context) {
	// Check quyền chặt chẽ
	vaiTro := c.GetString("USER_ROLE")
	if vaiTro != "admin_root" && vaiTro != "admin" && vaiTro != "quan_ly" {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa sản phẩm!"})
		return
	}

	// Lấy dữ liệu form
	maSP        := strings.TrimSpace(c.PostForm("ma_san_pham"))
	tenSP       := strings.TrimSpace(c.PostForm("ten_san_pham"))
	tenRutGon   := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku         := strings.TrimSpace(c.PostForm("sku"))
	
	giaBanStr   := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr    = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _   := strconv.ParseFloat(giaBanStr, 64)

	danhMucRaw  := c.PostForm("ma_danh_muc")
	danhMuc     := xuLyTags(danhMucRaw)

	thuongHieu  := c.PostForm("ma_thuong_hieu")
	donVi       := c.PostForm("don_vi")
	mauSac      := c.PostForm("mau_sac")
	hinhAnh     := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	thongSo     := c.PostForm("thong_so")
	moTa        := c.PostForm("mo_ta_chi_tiet")
	baoHanh, _  := strconv.Atoi(c.PostForm("bao_hanh_thang"))
	tinhTrang   := c.PostForm("tinh_trang")
	ghiChu      := c.PostForm("ghi_chu")
	
	trangThai := 0
	if c.PostForm("trang_thai") == "on" { trangThai = 1 }

	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được để trống!"})
		return
	}

	// Logic Core
	var sp *core.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")
	sheetID := cau_hinh.BienCauHinh.IdFileSheet

	core.KhoaHeThong.Lock()
	
	if maSP == "" {
		isNew = true
		maSP = core.TaoMaSPMoi()
		sp = &core.SanPham{
			SpreadsheetID: sheetID,
			MaSanPham:     maSP,
			NgayTao:       nowStr,
			NguoiTao:      userID,
		}
	} else {
		foundSP, ok := core.LayChiTietSanPham(maSP)
		if ok {
			sp = foundSP
		} else {
			sp = &core.SanPham{SpreadsheetID: sheetID, MaSanPham: maSP, NgayTao: nowStr}
		}
	}

	sp.TenSanPham = tenSP
	sp.TenRutGon = tenRutGon
	sp.Sku = sku
	sp.GiaBanLe = giaBan
	sp.MaDanhMuc = danhMuc
	sp.MaThuongHieu = thuongHieu
	sp.DonVi = donVi
	sp.MauSac = mauSac
	sp.UrlHinhAnh = hinhAnh
	sp.ThongSo = thongSo
	sp.MoTaChiTiet = moTa
	sp.BaoHanhThang = baoHanh
	sp.TinhTrang = tinhTrang
	sp.TrangThai = trangThai
	sp.GhiChu = ghiChu
	sp.NgayCapNhat = nowStr

	if isNew {
		// [CHUẨN] Dùng biến DongBatDau_SanPham
		sp.DongTrongSheet = core.DongBatDau_SanPham + len(core.LayDanhSachSanPham()) 
		core.ThemSanPhamVaoRam(sp)
	}
	
	core.KhoaHeThong.Unlock()

	// Ghi xuống Sheet
	targetRow := sp.DongTrongSheet
	if targetRow > 0 {
		ghi := core.ThemVaoHangCho
		sheet := "SAN_PHAM"

		ghi(sheetID, sheet, targetRow, core.CotSP_MaSanPham, sp.MaSanPham)
		ghi(sheetID, sheet, targetRow, core.CotSP_TenSanPham, sp.TenSanPham)
		ghi(sheetID, sheet, targetRow, core.CotSP_TenRutGon, sp.TenRutGon)
		ghi(sheetID, sheet, targetRow, core.CotSP_Sku, sp.Sku)
		ghi(sheetID, sheet, targetRow, core.CotSP_MaDanhMuc, sp.MaDanhMuc)
		ghi(sheetID, sheet, targetRow, core.CotSP_MaThuongHieu, sp.MaThuongHieu)
		ghi(sheetID, sheet, targetRow, core.CotSP_DonVi, sp.DonVi)
		ghi(sheetID, sheet, targetRow, core.CotSP_MauSac, sp.MauSac)
		ghi(sheetID, sheet, targetRow, core.CotSP_UrlHinhAnh, sp.UrlHinhAnh)
		ghi(sheetID, sheet, targetRow, core.CotSP_ThongSo, sp.ThongSo)
		ghi(sheetID, sheet, targetRow, core.CotSP_MoTaChiTiet, sp.MoTaChiTiet)
		ghi(sheetID, sheet, targetRow, core.CotSP_BaoHanhThang, sp.BaoHanhThang)
		ghi(sheetID, sheet, targetRow, core.CotSP_TinhTrang, sp.TinhTrang)
		ghi(sheetID, sheet, targetRow, core.CotSP_TrangThai, sp.TrangThai)
		ghi(sheetID, sheet, targetRow, core.CotSP_GiaBanLe, sp.GiaBanLe)
		ghi(sheetID, sheet, targetRow, core.CotSP_GhiChu, sp.GhiChu)
		ghi(sheetID, sheet, targetRow, core.CotSP_NguoiTao, sp.NguoiTao)
		ghi(sheetID, sheet, targetRow, core.CotSP_NgayTao, sp.NgayTao)
		ghi(sheetID, sheet, targetRow, core.CotSP_NgayCapNhat, sp.NgayCapNhat)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}

func xuLyTags(raw string) string {
	if !strings.Contains(raw, "[") { return raw }
	res := strings.ReplaceAll(raw, "[", "")
	res = strings.ReplaceAll(res, "]", "")
	res = strings.ReplaceAll(res, "{", "")
	res = strings.ReplaceAll(res, "}", "")
	res = strings.ReplaceAll(res, "\"value\":", "")
	res = strings.ReplaceAll(res, "\"", "")
	return res
}

func taoMaSPMoi() string {
	return core.TaoMaSPMoi()
}
