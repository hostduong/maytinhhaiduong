package chuc_nang

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core"

	"github.com/gin-gonic/gin"
)

// =============================================================
// 1. TRANG QUẢN LÝ (HIỂN THỊ)
// =============================================================
func TrangQuanLySanPham(c *gin.Context) {
	userID := c.GetString("USER_ID")
	
	// Lấy thông tin user & Chốt chặn lỗi
	kh, found := core.LayKhachHang(userID)
	if !found || kh == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// [MỚI] CHECK QUYỀN XEM (product.view)
	// Theo bảng phân quyền: Khách hàng = 0, Admin/Sale/Kho = 1
	if !core.KiemTraQuyen(kh.VaiTroQuyenHan, "product.view") {
		c.HTML(http.StatusForbidden, "quan_tri", gin.H{ // Hoặc dùng template báo lỗi riêng
			"TieuDe":   "Từ chối truy cập",
			"Error":    "Bạn không có quyền xem danh sách sản phẩm (product.view)!",
			"NhanVien": kh,
		})
		return
	}

	// Lấy dữ liệu từ Core
	listSP := core.LayDanhSachSanPham()
	listDM := core.LayDanhSachDanhMuc()
	listTH := core.LayDanhSachThuongHieu()

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
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

// =============================================================
// 2. API XỬ LÝ LƯU (THÊM / SỬA)
// =============================================================
func API_LuuSanPham(c *gin.Context) {
	// Lấy vai trò hiện tại
	vaiTro := c.GetString("USER_ROLE")

	// Lấy dữ liệu form cơ bản để check logic
	maSP      := strings.TrimSpace(c.PostForm("ma_san_pham"))
	giaBanStr := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr  = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _ := strconv.ParseFloat(giaBanStr, 64)

	// [MỚI] LOGIC PHÂN QUYỀN CHI TIẾT (MATRIX RBAC)
	if maSP == "" {
		// --- TRƯỜNG HỢP: TẠO MỚI ---
		if !core.KiemTraQuyen(vaiTro, "product.create") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thêm sản phẩm mới (product.create)!"})
			return
		}
	} else {
		// --- TRƯỜNG HỢP: CẬP NHẬT ---
		if !core.KiemTraQuyen(vaiTro, "product.edit") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa thông tin sản phẩm (product.edit)!"})
			return
		}

		// [ĐẶC BIỆT] CHECK QUYỀN SỬA GIÁ (product.edit_price)
		// Logic: Load sản phẩm cũ lên so sánh giá
		spCu, ok := core.LayChiTietSanPham(maSP)
		if ok {
			if spCu.GiaBanLe != giaBan {
				// Nếu giá thay đổi -> Phải có quyền edit_price
				if !core.KiemTraQuyen(vaiTro, "product.edit_price") {
					c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị viên mới được quyền sửa giá bán (product.edit_price)!"})
					return
				}
			}
		}
	}

	// ---------------------------------------------------------
	// NẾU QUA ĐƯỢC CÁC CHECK TRÊN THÌ MỚI XỬ LÝ DỮ LIỆU
	// ---------------------------------------------------------
	tenSP       := strings.TrimSpace(c.PostForm("ten_san_pham"))
	tenRutGon   := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku         := strings.TrimSpace(c.PostForm("sku"))
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

	// Tương tác Core (Giữ nguyên logic cũ của bạn)
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

	// Gán dữ liệu
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
	// Cập nhật người sửa cuối (nếu cần)
	// sp.NguoiCapNhat = userID 

	if isNew {
		sp.DongTrongSheet = core.DongBatDau_SanPham + len(core.LayDanhSachSanPham()) 
		core.ThemSanPhamVaoRam(sp)
	}
	
	core.KhoaHeThong.Unlock()

	// Đẩy xuống hàng chờ ghi Sheet
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

// Helper xử lý tags
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
