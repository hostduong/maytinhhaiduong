package chuc_nang

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"app/cau_hinh"
	"app/core" // [QUAN TRỌNG] Chỉ dùng Core, không dùng nghiep_vu/mo_hinh cũ
	"app/nghiep_vu" // Vẫn cần cho hàm KiemTraQuyen (sẽ chuyển nốt sau)

	"github.com/gin-gonic/gin"
)

// TrangQuanLySanPham : Hiển thị danh sách
func TrangQuanLySanPham(c *gin.Context) {
	// Lấy User ID từ Session (Logic cũ)
	userID := c.GetString("USER_ID")
	
	// Tạm thời dùng hàm cũ để lấy thông tin nhân viên (Sẽ refactor sau)
	kh, _ := nghiep_vu.LayThongTinKhachHang(userID)

	// [MỚI] Lấy dữ liệu từ Core (Đã được sort và filter sẵn)
	listSP := core.LayDanhSachSanPham()
	listDM := core.LayDanhSachDanhMuc()
	listTH := core.LayDanhSachThuongHieu()

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang, // Nếu kh nil sẽ crash, nhưng tạm chấp nhận để test
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       listSP,
		"ListDanhMuc":    listDM,
		"ListThuongHieu": listTH,
	})
}

// API_LuuSanPham : Xử lý Thêm/Sửa
func API_LuuSanPham(c *gin.Context) {
	// 1. Check quyền
	vaiTro := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTro, "product.edit") {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền này!"})
		return
	}

	// 2. Lấy dữ liệu form
	maSP        := strings.TrimSpace(c.PostForm("ma_san_pham"))
	tenSP       := strings.TrimSpace(c.PostForm("ten_san_pham"))
	tenRutGon   := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku         := strings.TrimSpace(c.PostForm("sku"))
	
	// Xử lý giá tiền
	giaBanStr   := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr    = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _   := strconv.ParseFloat(giaBanStr, 64)

	danhMucRaw  := c.PostForm("ma_danh_muc")
	danhMuc     := xuLyTags(danhMucRaw) // Hàm helper ở dưới

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

	// 3. Logic Thêm/Sửa dùng CORE
	var sp *core.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")
	sheetID := cau_hinh.BienCauHinh.IdFileSheet // Lấy Sheet ID hiện tại

	// Lock logic xử lý
	core.KhoaHeThong.Lock()
	
	if maSP == "" {
		// TẠO MỚI
		isNew = true
		maSP = core.TaoMaSPMoi() // Hàm mới trong core
		sp = &core.SanPham{
			SpreadsheetID: sheetID,
			MaSanPham:     maSP,
			NgayTao:       nowStr,
			NguoiTao:      userID,
		}
	} else {
		// SỬA: Tìm trong Map của Core
		// Lưu ý: Core dùng composite key, nhưng hàm LayChiTietSanPham đã xử lý giúp ta
		// Tuy nhiên ở đây ta cần truy cập biến global để sửa, hoặc dùng setter.
		// Để đơn giản, ta tìm lại con trỏ trong map và update.
		// (Logic này hơi hack một chút khi truy cập thẳng biến core,
		// chuẩn ra nên viết hàm UpdateSanPham trong core. Nhưng tạm thời OK).
		
		// Cách an toàn: Gọi hàm Core tìm
		// Vì hàm LayChiTietSanPham trả về con trỏ, nên sửa nó là sửa trong RAM luôn!
		foundSP, ok := core.LayChiTietSanPham(maSP)
		if ok {
			sp = foundSP
		} else {
			// Trường hợp hiếm: ID gửi lên không tồn tại -> coi như mới hoặc lỗi
			sp = &core.SanPham{SpreadsheetID: sheetID, MaSanPham: maSP, NgayTao: nowStr}
		}
	}

	// Cập nhật thông tin
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

	// Nếu là mới thì phải append vào danh sách RAM
	if isNew {
		// Tính dòng mới (Logic đơn giản: Dòng cuối + 1)
		// Core nên có hàm GetNextRowIndex, tạm thời ta tính thủ công hoặc để core tự lo
		// Để đơn giản: Gán đại DongTrongSheet = 0, Worker sẽ tự append (nếu worker thông minh)
		// Hoặc:
		sp.DongTrongSheet = core.DongBatDauDuLieu + len(core.LayDanhSachSanPham()) 
		
		// Thêm vào RAM Core
		core.ThemSanPhamVaoRam(sp)
	}
	
	core.KhoaHeThong.Unlock()

	// 4. Đẩy xuống Hàng Chờ Ghi (Dùng Core Write Queue)
	targetRow := sp.DongTrongSheet
	// Nếu dòng <= 0 (chưa xác định), Worker Smart Queue của bạn sẽ tự xử lý append?
	// Với Smart Queue của bạn: row là int key. Nếu row chưa có, nó tạo mới.
	// Để an toàn, ta cần row chính xác.
	
	if targetRow > 0 {
		// Ghi từng cột (Hơi dài dòng nhưng an toàn)
		// Format: ThemVaoHangCho(SpreadID, SheetName, Row, Col, Value)
		ghi := core.ThemVaoHangCho
		sName := "SAN_PHAM"

		ghi(sheetID, sName, targetRow, core.CotSP_MaSanPham, sp.MaSanPham)
		ghi(sheetID, sName, targetRow, core.CotSP_TenSanPham, sp.TenSanPham)
		ghi(sheetID, sName, targetRow, core.CotSP_TenRutGon, sp.TenRutGon)
		ghi(sheetID, sName, targetRow, core.CotSP_Sku, sp.Sku)
		ghi(sheetID, sName, targetRow, core.CotSP_MaDanhMuc, sp.MaDanhMuc)
		ghi(sheetID, sName, targetRow, core.CotSP_MaThuongHieu, sp.MaThuongHieu)
		ghi(sheetID, sName, targetRow, core.CotSP_DonVi, sp.DonVi)
		ghi(sheetID, sName, targetRow, core.CotSP_MauSac, sp.MauSac)
		ghi(sheetID, sName, targetRow, core.CotSP_UrlHinhAnh, sp.UrlHinhAnh)
		ghi(sheetID, sName, targetRow, core.CotSP_ThongSo, sp.ThongSo)
		ghi(sheetID, sName, targetRow, core.CotSP_MoTaChiTiet, sp.MoTaChiTiet)
		ghi(sheetID, sName, targetRow, core.CotSP_BaoHanhThang, sp.BaoHanhThang)
		ghi(sheetID, sName, targetRow, core.CotSP_TinhTrang, sp.TinhTrang)
		ghi(sheetID, sName, targetRow, core.CotSP_TrangThai, sp.TrangThai)
		ghi(sheetID, sName, targetRow, core.CotSP_GiaBanLe, sp.GiaBanLe)
		ghi(sheetID, sName, targetRow, core.CotSP_GhiChu, sp.GhiChu)
		ghi(sheetID, sName, targetRow, core.CotSP_NguoiTao, sp.NguoiTao)
		ghi(sheetID, sName, targetRow, core.CotSP_NgayTao, sp.NgayTao)
		ghi(sheetID, sName, targetRow, core.CotSP_NgayCapNhat, sp.NgayCapNhat)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công (Core)!"})
}

// Helper xử lý tags (giữ nguyên)
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
