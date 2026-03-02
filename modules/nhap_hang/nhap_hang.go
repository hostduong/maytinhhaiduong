package nhap_hang

import (
	"fmt"
	"net/http"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// ============================================================================
// 1. RENDER GIAO DIỆN HTML
// ============================================================================
func TrangNhapHangMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	// Lấy thông tin người dùng hiện tại (Để hiển thị Avatar/Tên trên Header & Sidebar)
	me, ok := core.LayKhachHang(shopID, userID)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	meCopy := *me

	meCopy.StyleLevel = core.LayCapBacVaiTro(shopID, userID, meCopy.VaiTroQuyenHan)
	if meCopy.MaKhachHang == "0000000000000000000" || meCopy.MaKhachHang == "0000000000000000001" || meCopy.VaiTroQuyenHan == "quan_tri_he_thong" {
		meCopy.StyleLevel, meCopy.StyleTheme = 0, 9
	}

	// Lấy dữ liệu Master Data ném ra form Nhập hàng
	danhSachNCC := core.LayDanhSachNhaCungCap(shopID)
	danhSachSP := core.LayDanhSachSanPhamMayTinh(shopID)

	c.HTML(http.StatusOK, "master_nhap_hang", gin.H{
		"TieuDe":      "Nhập Hàng",
		"NhanVien":    &meCopy,
		"DanhSachNCC": danhSachNCC,
		"DanhSachSP":  danhSachSP,
	})
}

// ============================================================================
// 2. CẤU TRÚC NHẬN JSON TỪ GIAO DIỆN
// ============================================================================
type ChiTietInput struct {
	MaSKU      string  `json:"ma_sku"`
	SoLuong    int     `json:"so_luong"`
	DonGiaNhap float64 `json:"don_gia_nhap"`
}

type PhieuNhapInput struct {
	MaNhaCungCap        string         `json:"ma_nha_cung_cap"`
	MaKho               string         `json:"ma_kho"`
	NgayNhap            string         `json:"ngay_nhap"`
	SoHoaDon            string         `json:"so_hoa_don"`
	GhiChuPhieu         string         `json:"ghi_chu_phieu"`
	GiamGiaPhieu        float64        `json:"giam_gia_phieu"`
	PhuongThucThanhToan string         `json:"phuong_thuc_thanh_toan"`
	DaTra               float64        `json:"da_tra"`
	ChiTiet             []ChiTietInput `json:"chi_tiet"`
}

// ============================================================================
// 3. API XỬ LÝ LƯU PHIẾU NHẬP
// ============================================================================
func API_LuuPhieuNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

	var input PhieuNhapInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"status": "error", "msg": "Dữ liệu đầu vào không hợp lệ!"})
		return
	}

	if len(input.ChiTiet) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiếu nhập chưa có sản phẩm nào!"})
		return
	}

	loc := time.FixedZone("ICT", 7*3600)
	now := time.Now().In(loc)
	nowStr := now.Format("2006-01-02 15:04:05")

	// 1. Tự động sinh Mã Phiếu Nhập
	maPN := fmt.Sprintf("PN%s%s", now.Format("060102"), core.LayChuoiSoNgauNhien(4))

	// 2. Tính toán tổng tiền
	var tongTienHang float64 = 0
	for _, ct := range input.ChiTiet {
		tongTienHang += ct.DonGiaNhap * float64(ct.SoLuong)
	}

	thanhToan := tongTienHang - input.GiamGiaPhieu
	if thanhToan < 0 { thanhToan = 0 }
	conNo := thanhToan - input.DaTra

	trangThaiThanhToan := "CHUA_THANH_TOAN"
	if input.DaTra >= thanhToan && thanhToan > 0 {
		trangThaiThanhToan = "DA_THANH_TOAN"
	} else if input.DaTra > 0 {
		trangThaiThanhToan = "THANH_TOAN_MOT_PHAN"
	}

	nguoiTao, _ := core.LayKhachHang(shopID, userID)
	tenNguoiTao := "Hệ thống"
	if nguoiTao != nil { tenNguoiTao = nguoiTao.TenDangNhap }

	// 3. Khởi tạo Object Phiếu Nhập (Header)
	pn := &core.PhieuNhap{
		SpreadsheetID: shopID, MaPhieuNhap: maPN, MaNhaCungCap: input.MaNhaCungCap, MaKho: input.MaKho,
		NgayNhap: input.NgayNhap, TrangThai: 1, SoHoaDon: input.SoHoaDon, NgayHoaDon: "", UrlChungTu: "",
		TongTienPhieu: tongTienHang, GiamGiaPhieu: input.GiamGiaPhieu, DaThanhToan: input.DaTra, ConNo: conNo,
		PhuongThucThanhToan: input.PhuongThucThanhToan, TrangThaiThanhToan: trangThaiThanhToan,
		GhiChu: input.GhiChuPhieu, NguoiTao: tenNguoiTao, NgayTao: nowStr, NgayCapNhat: nowStr,
		ChiTiet: make([]*core.ChiTietPhieuNhap, 0),
	}

	// Đẩy Header vào Queue
	rowPN := make([]interface{}, 22)
	rowPN[core.CotPN_MaPhieuNhap] = pn.MaPhieuNhap; rowPN[core.CotPN_MaNhaCungCap] = pn.MaNhaCungCap; rowPN[core.CotPN_MaKho] = pn.MaKho
	rowPN[core.CotPN_NgayNhap] = pn.NgayNhap; rowPN[core.CotPN_TrangThai] = pn.TrangThai; rowPN[core.CotPN_SoHoaDon] = pn.SoHoaDon
	rowPN[core.CotPN_TongTienPhieu] = pn.TongTienPhieu; rowPN[core.CotPN_GiamGiaPhieu] = pn.GiamGiaPhieu; rowPN[core.CotPN_DaThanhToan] = pn.DaThanhToan
	rowPN[core.CotPN_ConNo] = pn.ConNo; rowPN[core.CotPN_PhuongThucThanhToan] = pn.PhuongThucThanhToan; rowPN[core.CotPN_TrangThaiThanhToan] = pn.TrangThaiThanhToan
	rowPN[core.CotPN_GhiChu] = pn.GhiChu; rowPN[core.CotPN_NguoiTao] = pn.NguoiTao; rowPN[core.CotPN_NgayTao] = pn.NgayTao; rowPN[core.CotPN_NgayCapNhat] = pn.NgayCapNhat
	core.PushAppend(shopID, core.TenSheetPhieuNhap, rowPN)

	// 4. Xử lý Chi tiết & Sinh Serial
	core.KhoaHeThong.Lock() // Khóa RAM để thao tác an toàn
	for _, item := range input.ChiTiet {
		spCache, ok := core.LayChiTietSKUMayTinh(shopID, item.MaSKU)
		tenSP, donVi, maSP := "Sản phẩm không xác định", "Cái", ""
		if ok && spCache != nil {
			tenSP = spCache.TenSanPham
			donVi = spCache.DonVi
			maSP = spCache.MaSanPham
		}

		thanhTienDong := item.DonGiaNhap * float64(item.SoLuong)
		giaVonThucTe := item.DonGiaNhap // Tương lai có thể chia tỉ lệ giảm giá phiếu vào đây

		ct := &core.ChiTietPhieuNhap{
			SpreadsheetID: shopID, MaPhieuNhap: maPN, MaSanPham: maSP, MaSKU: item.MaSKU,
			TenSanPham: tenSP, DonVi: donVi, SoLuong: item.SoLuong, DonGiaNhap: item.DonGiaNhap,
			ThanhTienDong: thanhTienDong, GiaVonThucTe: giaVonThucTe,
		}
		pn.ChiTiet = append(pn.ChiTiet, ct)

		// Đẩy dòng Chi tiết vào Queue
		rowCT := make([]interface{}, 15)
		rowCT[core.CotCTPN_MaPhieuNhap] = ct.MaPhieuNhap; rowCT[core.CotCTPN_MaSanPham] = ct.MaSanPham; rowCT[core.CotCTPN_MaSKU] = ct.MaSKU
		rowCT[core.CotCTPN_TenSanPham] = ct.TenSanPham; rowCT[core.CotCTPN_DonVi] = ct.DonVi; rowCT[core.CotCTPN_SoLuong] = ct.SoLuong
		rowCT[core.CotCTPN_DonGiaNhap] = ct.DonGiaNhap; rowCT[core.CotCTPN_ThanhTienDong] = ct.ThanhTienDong; rowCT[core.CotCTPN_GiaVonThucTe] = ct.GiaVonThucTe
		core.PushAppend(shopID, core.TenSheetChiTietPhieuNhap, rowCT)

		// TỰ ĐỘNG SINH SERIAL CHO TỪNG SẢN PHẨM NHẬP (Ví dụ: Nhập 5 => Sinh 5 Serial)
		for i := 0; i < item.SoLuong; i++ {
			imei := fmt.Sprintf("SN%s%s", now.Format("060102"), core.LayChuoiSoNgauNhien(6))
			sr := &core.SerialSanPham{
				SpreadsheetID: shopID, SerialIMEI: imei, MaSanPham: maSP, MaSKU: item.MaSKU,
				MaNhaCungCap: input.MaNhaCungCap, MaPhieuNhap: maPN, TrangThai: 1, // 1 = Đang trong kho
				NgayNhapKho: input.NgayNhap, GiaVonNhap: giaVonThucTe, MaKho: input.MaKho, NgayCapNhat: nowStr,
			}
			
			// Lưu vào RAM
			core.CacheSerialSanPham[shopID] = append(core.CacheSerialSanPham[shopID], sr)
			core.CacheMapSerial[core.TaoCompositeKey(shopID, imei)] = sr

			// Đẩy Serial vào Queue
			rowSR := make([]interface{}, 19)
			rowSR[core.CotSR_SerialIMEI] = sr.SerialIMEI; rowSR[core.CotSR_MaSanPham] = sr.MaSanPham; rowSR[core.CotSR_MaSKU] = sr.MaSKU
			rowSR[core.CotSR_MaNhaCungCap] = sr.MaNhaCungCap; rowSR[core.CotSR_MaPhieuNhap] = sr.MaPhieuNhap; rowSR[core.CotSR_TrangThai] = sr.TrangThai
			rowSR[core.CotSR_NgayNhapKho] = sr.NgayNhapKho; rowSR[core.CotSR_GiaVonNhap] = sr.GiaVonNhap; rowSR[core.CotSR_MaKho] = sr.MaKho; rowSR[core.CotSR_NgayCapNhat] = sr.NgayCapNhat
			core.PushAppend(shopID, core.TenSheetSerial, rowSR)
		}
	}

	// Cập nhật Cache Phiếu Nhập
	core.CachePhieuNhap[shopID] = append(core.CachePhieuNhap[shopID], pn)
	core.CacheMapPhieuNhap[core.TaoCompositeKey(shopID, maPN)] = pn
	core.KhoaHeThong.Unlock()

	c.JSON(200, gin.H{"status": "ok", "msg": "Tạo Phiếu Nhập và Khởi tạo Serial thành công!", "ma_phieu": maPN})
}
