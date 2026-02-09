package chuc_nang

import (
	"encoding/json"
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

	// [FIX LỖI TRẮNG TRANG]
	// Kiểm tra quyền xem (product.view)
	// Nếu không có quyền -> Trả về trang HTML báo lỗi trực tiếp (Không dùng template Dashboard để tránh Panic thiếu dữ liệu)
	if !core.KiemTraQuyen(kh.VaiTroQuyenHan, "product.view") {
		c.Writer.WriteHeader(http.StatusForbidden)
		c.Writer.Write([]byte(`
			<!DOCTYPE html>
			<html lang="vi">
			<head>
				<meta charset="UTF-8">
				<title>Truy cập bị từ chối</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<style>
					body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; background-color: #f3f4f6; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
					.card { background: white; padding: 2rem; border-radius: 1rem; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); text-align: center; max-width: 400px; width: 90%; }
					.icon { font-size: 3rem; margin-bottom: 1rem; }
					h2 { color: #dc2626; margin: 0 0 0.5rem 0; font-size: 1.5rem; }
					p { color: #4b5563; margin-bottom: 1.5rem; }
					a { display: inline-block; background-color: #2563eb; color: white; padding: 0.5rem 1rem; border-radius: 0.5rem; text-decoration: none; font-weight: 500; transition: background 0.2s; }
					a:hover { background-color: #1d4ed8; }
				</style>
			</head>
			<body>
				<div class="card">
					<div class="icon">⛔</div>
					<h2>Truy cập bị từ chối</h2>
					<p>Bạn không có quyền xem danh sách sản phẩm.<br><small>(Mã quyền: product.view)</small></p>
					<a href="/admin/tong-quan">Quay lại Dashboard</a>
				</div>
			</body>
			</html>
		`))
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
	vaiTro := c.GetString("USER_ROLE")

	maSP       := strings.TrimSpace(c.PostForm("ma_san_pham"))
	thuongHieu := c.PostForm("ma_thuong_hieu")

	giaBanStr := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr  = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _ := strconv.ParseFloat(giaBanStr, 64)

	// --- CHECK QUYỀN ---
	if maSP == "" {
		if !core.KiemTraQuyen(vaiTro, "product.create") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thêm sản phẩm mới!"})
			return
		}
	} else {
		if !core.KiemTraQuyen(vaiTro, "product.edit") {
			c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền sửa sản phẩm!"})
			return
		}
		
		spCu, ok := core.LayChiTietSanPham(maSP)
		if ok && spCu.GiaBanLe != giaBan {
			if !core.KiemTraQuyen(vaiTro, "product.edit_price") {
				c.JSON(200, gin.H{"status": "error", "msg": "Chỉ Quản trị viên được sửa giá bán!"})
				return
			}
		}
	}

	// --- LẤY DỮ LIỆU ---
	tenSP       := strings.TrimSpace(c.PostForm("ten_san_pham"))
	tenRutGon   := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku         := strings.TrimSpace(c.PostForm("sku"))
	
	// Xử lý Tags JSON từ Tagify -> Chuỗi "Tag1|Tag2"
	danhMucRaw  := c.PostForm("ma_danh_muc")
	danhMuc     := xuLyTags(danhMucRaw) 

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

	// --- LOGIC CORE ---
	var sp *core.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")
	sheetID := cau_hinh.BienCauHinh.IdFileSheet

	core.KhoaHeThong.Lock()
	
	if maSP == "" {
		isNew = true
		// Sinh mã mới (HD + Brand + YYMM...)
		maSP = core.TaoMaSPMoi(thuongHieu) 
		
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

	if isNew {
		sp.DongTrongSheet = core.DongBatDau_SanPham + len(core.LayDanhSachSanPham()) 
		core.ThemSanPhamVaoRam(sp)
	}
	
	core.KhoaHeThong.Unlock()

	// --- GHI XUỐNG SHEET ---
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

// Struct JSON Tagify
type TagifyItem struct {
	Value string `json:"value"`
}

// Hàm xử lý Tags: Input `[{"value":"A"}, {"value":"B"}]` -> Output `A|B`
func xuLyTags(raw string) string {
	if raw == "" { return "" }
	if !strings.Contains(raw, "[") { return raw }

	var items []TagifyItem
	err := json.Unmarshal([]byte(raw), &items)
	if err != nil { return raw }

	var values []string
	for _, item := range items {
		val := strings.TrimSpace(item.Value)
		if val != "" {
			values = append(values, val)
		}
	}
	return strings.Join(values, "|")
}

func taoMaSPMoi() string {
	return core.TaoMaSPMoi("") 
}
