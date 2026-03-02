package nhap_hang

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"app/core"
	"github.com/gin-gonic/gin"
)

// ============================================================================
// 1. RENDER GIAO DIỆN HTML VÀ NẠP PHIẾU NHÁP
// ============================================================================
func TrangNhapHangMaster(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")

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

	danhSachNCC := core.LayDanhSachNhaCungCap(shopID)
	danhSachSP := core.LayDanhSachSanPhamMayTinh(shopID)

	core.GetSheetLock(shopID, core.TenSheetPhieuNhap).RLock()
	allPhieu := core.CachePhieuNhap[shopID]
	var danhSachNhaps []*core.PhieuNhap
	for _, p := range allPhieu {
		if p.TrangThai == 0 || p.TrangThai == 2 || p.TrangThai == -1 {
			danhSachNhaps = append(danhSachNhaps, p)
		}
	}
	core.GetSheetLock(shopID, core.TenSheetPhieuNhap).RUnlock()

	c.HTML(http.StatusOK, "master_nhap_hang", gin.H{
		"TieuDe":        "Nhập Hàng",
		"NhanVien":      &meCopy,
		"DanhSachNCC":   danhSachNCC,
		"DanhSachSP":    danhSachSP,
		"DanhSachNhaps": danhSachNhaps,
	})
}

// ============================================================================
// 2. CẤU TRÚC NHẬN JSON TỪ GIAO DIỆN
// ============================================================================
type ChiTietInput struct {
	MaSKU      string   `json:"ma_sku"`
	SoLuong    int      `json:"so_luong"`
	DonGiaNhap float64  `json:"don_gia_nhap"`
	Serials    []string `json:"serials"`
}

type PhieuNhapInput struct {
	MaPhieuNhap         string         `json:"ma_phieu_nhap"`
	MaNhaCungCap        string         `json:"ma_nha_cung_cap"`
	MaKho               string         `json:"ma_kho"`
	NgayNhap            string         `json:"ngay_nhap"`
	SoHoaDon            string         `json:"so_hoa_don"`
	GhiChuPhieu         string         `json:"ghi_chu_phieu"`
	GiamGiaPhieu        float64        `json:"giam_gia_phieu"`
	ChiPhiNhap          float64        `json:"chi_phi_nhap"`
	DaTra               float64        `json:"da_tra"`
	PhuongThucThanhToan string         `json:"phuong_thuc_thanh_toan"`
	TrangThai           int            `json:"trang_thai"`
	ChiTiet             []ChiTietInput `json:"chi_tiet"`
}

// ============================================================================
// HÀM KIỂM TRA QUYỀN (HỖ TRỢ BYPASS CHO ADMIN)
// ============================================================================
func checkQuyenNhapHang(vaiTro string, userID string) bool {
	// 1. Bypass cho các tài khoản Root/Admin
	if vaiTro == "quan_tri_he_thong" || vaiTro == "quan_tri_vien_he_thong" || 
	   vaiTro == "quan_tri_cua_hang" || vaiTro == "quan_tri_vien_cua_hang" || 
	   userID == "0000000000000000001" {
		return true
	}
	
	// 2. Tương lai: Kiểm tra logic quyền mềm (VD: "stock.import") ở đây
	// Hiện tại mặc định trả về true nếu đã lọt qua middleware
	return true 
}

// ============================================================================
// 3. API XỬ LÝ LƯU PHIẾU NHẬP
// ============================================================================
func API_LuuPhieuNhap(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if !checkQuyenNhapHang(vaiTro, userID) {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thực hiện thao tác nhập hàng!"})
		return
	}

	var input PhieuNhapInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"status": "error", "msg": "Dữ liệu không hợp lệ!"})
		return
	}

	if len(input.ChiTiet) == 0 {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiếu nhập chưa có sản phẩm!"})
		return
	}

	loc := time.FixedZone("ICT", 7*3600)
	nowStr := time.Now().In(loc).Format("2006-01-02 15:04:05")

	nguoiThaoTac, _ := core.LayKhachHang(shopID, userID)
	tenNguoiThaoTac := "Hệ thống"
	if nguoiThaoTac != nil {
		tenNguoiThaoTac = nguoiThaoTac.TenDangNhap
	}

	var tongTienHang float64 = 0
	for _, ct := range input.ChiTiet {
		tongTienHang += ct.DonGiaNhap * float64(ct.SoLuong)
	}

	thanhToan := tongTienHang - input.GiamGiaPhieu + input.ChiPhiNhap
	if thanhToan < 0 { thanhToan = 0 }
	conNo := thanhToan - input.DaTra

	trangThaiThanhToan := "CHUA_THANH_TOAN"
	if input.DaTra >= thanhToan && thanhToan > 0 {
		trangThaiThanhToan = "DA_THANH_TOAN"
	} else if input.DaTra > 0 {
		trangThaiThanhToan = "THANH_TOAN_MOT_PHAN"
	}

	chiTietBytes, _ := json.Marshal(input.ChiTiet)
	chiTietJsonStr := string(chiTietBytes)

	core.GetSheetLock(shopID, core.TenSheetPhieuNhap).Lock()
	defer core.GetSheetLock(shopID, core.TenSheetPhieuNhap).Unlock()

	var pn *core.PhieuNhap
	isUpdate := false

	if input.MaPhieuNhap != "" && !strings.HasPrefix(input.MaPhieuNhap, "TEMP_") {
		if oldPn, ok := core.CacheMapPhieuNhap[core.TaoCompositeKey(shopID, input.MaPhieuNhap)]; ok {
			if oldPn.TrangThai == 1 {
				c.JSON(200, gin.H{"status": "error", "msg": "Phiếu này đã hoàn thành, không thể sửa!"})
				return
			}
			pn = oldPn
			isUpdate = true
		}
	}

	if !isUpdate {
		maPN := fmt.Sprintf("PN%s%s", time.Now().In(loc).Format("060102"), core.LayChuoiSoNgauNhien(4))
		
		pn = &core.PhieuNhap{
			SpreadsheetID: shopID, MaPhieuNhap: maPN, MaNhaCungCap: input.MaNhaCungCap, MaKho: input.MaKho, 
			NgayNhap: input.NgayNhap, ChiTietJson: chiTietJsonStr, TrangThai: input.TrangThai,
			SoHoaDon: input.SoHoaDon, TongTienPhieu: tongTienHang, GiamGiaPhieu: input.GiamGiaPhieu, 
			ChiPhiNhap: input.ChiPhiNhap, DaThanhToan: input.DaTra, ConNo: conNo,
			PhuongThucThanhToan: input.PhuongThucThanhToan, TrangThaiThanhToan: trangThaiThanhToan, 
			GhiChu: input.GhiChuPhieu, NguoiTao: tenNguoiThaoTac, NgayTao: nowStr, 
			NguoiCapNhat: tenNguoiThaoTac, NgayCapNhat: nowStr,
			ChiTiet: make([]*core.ChiTietPhieuNhap, 0),
		}

		if input.TrangThai == 1 {
			pn.NguoiDuyet = tenNguoiThaoTac
			pn.NgayDuyet = nowStr
		}

		rowPN := make([]interface{}, 23)
		rowPN[core.CotPN_MaPhieuNhap] = pn.MaPhieuNhap; rowPN[core.CotPN_MaNhaCungCap] = pn.MaNhaCungCap
		rowPN[core.CotPN_MaKho] = pn.MaKho; rowPN[core.CotPN_NgayNhap] = pn.NgayNhap
		rowPN[core.CotPN_ChiTietJson] = pn.ChiTietJson; rowPN[core.CotPN_TrangThai] = pn.TrangThai
		rowPN[core.CotPN_SoHoaDon] = pn.SoHoaDon; rowPN[core.CotPN_NgayHoaDon] = pn.NgayHoaDon; rowPN[core.CotPN_UrlChungTu] = pn.UrlChungTu
		rowPN[core.CotPN_TongTienPhieu] = pn.TongTienPhieu; rowPN[core.CotPN_GiamGiaPhieu] = pn.GiamGiaPhieu
		rowPN[core.CotPN_ChiPhiNhap] = pn.ChiPhiNhap; rowPN[core.CotPN_DaThanhToan] = pn.DaThanhToan
		rowPN[core.CotPN_ConNo] = pn.ConNo; rowPN[core.CotPN_PhuongThucThanhToan] = pn.PhuongThucThanhToan
		rowPN[core.CotPN_TrangThaiThanhToan] = pn.TrangThaiThanhToan; rowPN[core.CotPN_GhiChu] = pn.GhiChu
		rowPN[core.CotPN_NguoiTao] = pn.NguoiTao; rowPN[core.CotPN_NgayTao] = pn.NgayTao
		rowPN[core.CotPN_NguoiDuyet] = pn.NguoiDuyet; rowPN[core.CotPN_NgayDuyet] = pn.NgayDuyet
		rowPN[core.CotPN_NguoiCapNhat] = pn.NguoiCapNhat; rowPN[core.CotPN_NgayCapNhat] = pn.NgayCapNhat
		
		core.PushAppend(shopID, core.TenSheetPhieuNhap, rowPN)
		
		core.KhoaHeThong.Lock()
		pn.DongTrongSheet = core.DongBatDau_PhieuNhap + len(core.CachePhieuNhap[shopID])
		core.CachePhieuNhap[shopID] = append(core.CachePhieuNhap[shopID], pn)
		core.CacheMapPhieuNhap[core.TaoCompositeKey(shopID, pn.MaPhieuNhap)] = pn
		core.KhoaHeThong.Unlock()

	} else {
		core.KhoaHeThong.Lock()
		pn.MaNhaCungCap = input.MaNhaCungCap; pn.MaKho = input.MaKho; pn.NgayNhap = input.NgayNhap
		pn.ChiTietJson = chiTietJsonStr; pn.TrangThai = input.TrangThai; pn.SoHoaDon = input.SoHoaDon
		pn.TongTienPhieu = tongTienHang; pn.GiamGiaPhieu = input.GiamGiaPhieu; pn.ChiPhiNhap = input.ChiPhiNhap
		pn.DaThanhToan = input.DaTra; pn.ConNo = conNo; pn.PhuongThucThanhToan = input.PhuongThucThanhToan
		pn.TrangThaiThanhToan = trangThaiThanhToan; pn.GhiChu = input.GhiChuPhieu
		pn.NguoiCapNhat = tenNguoiThaoTac; pn.NgayCapNhat = nowStr

		if input.TrangThai == 1 {
			pn.NguoiDuyet = tenNguoiThaoTac; pn.NgayDuyet = nowStr
		}
		core.KhoaHeThong.Unlock()

		r := pn.DongTrongSheet
		sheet := core.TenSheetPhieuNhap
		ghi := core.ThemVaoHangCho

		ghi(shopID, sheet, r, core.CotPN_MaNhaCungCap, pn.MaNhaCungCap); ghi(shopID, sheet, r, core.CotPN_MaKho, pn.MaKho)
		ghi(shopID, sheet, r, core.CotPN_NgayNhap, pn.NgayNhap); ghi(shopID, sheet, r, core.CotPN_ChiTietJson, pn.ChiTietJson)
		ghi(shopID, sheet, r, core.CotPN_TrangThai, pn.TrangThai); ghi(shopID, sheet, r, core.CotPN_SoHoaDon, pn.SoHoaDon)
		ghi(shopID, sheet, r, core.CotPN_TongTienPhieu, pn.TongTienPhieu); ghi(shopID, sheet, r, core.CotPN_GiamGiaPhieu, pn.GiamGiaPhieu)
		ghi(shopID, sheet, r, core.CotPN_ChiPhiNhap, pn.ChiPhiNhap); ghi(shopID, sheet, r, core.CotPN_DaThanhToan, pn.DaThanhToan)
		ghi(shopID, sheet, r, core.CotPN_ConNo, pn.ConNo); ghi(shopID, sheet, r, core.CotPN_PhuongThucThanhToan, pn.PhuongThucThanhToan)
		ghi(shopID, sheet, r, core.CotPN_TrangThaiThanhToan, pn.TrangThaiThanhToan); ghi(shopID, sheet, r, core.CotPN_GhiChu, pn.GhiChu)
		ghi(shopID, sheet, r, core.CotPN_NguoiCapNhat, pn.NguoiCapNhat); ghi(shopID, sheet, r, core.CotPN_NgayCapNhat, pn.NgayCapNhat)
		
		if input.TrangThai == 1 {
			ghi(shopID, sheet, r, core.CotPN_NguoiDuyet, pn.NguoiDuyet); ghi(shopID, sheet, r, core.CotPN_NgayDuyet, pn.NgayDuyet)
		}
	}

	if input.TrangThai == 1 {
		core.KhoaHeThong.Lock()
		pn.ChiTiet = make([]*core.ChiTietPhieuNhap, 0) 
		
		for _, item := range input.ChiTiet {
			spCache, ok := core.LayChiTietSKUMayTinh(shopID, item.MaSKU)
			tenSP, donVi, maSP := "Sản phẩm không xác định", "Cái", ""
			if ok && spCache != nil {
				tenSP = spCache.TenSanPham; donVi = spCache.DonVi; maSP = spCache.MaSanPham
			}

			ct := &core.ChiTietPhieuNhap{
				SpreadsheetID: shopID, MaPhieuNhap: pn.MaPhieuNhap, MaSanPham: maSP, MaSKU: item.MaSKU,
				TenSanPham: tenSP, DonVi: donVi, SoLuong: item.SoLuong, DonGiaNhap: item.DonGiaNhap,
				ThanhTienDong: item.DonGiaNhap * float64(item.SoLuong), GiaVonThucTe: item.DonGiaNhap,
			}
			pn.ChiTiet = append(pn.ChiTiet, ct)

			rowCT := make([]interface{}, 15)
			rowCT[core.CotCTPN_MaPhieuNhap] = ct.MaPhieuNhap; rowCT[core.CotCTPN_MaSanPham] = ct.MaSanPham
			rowCT[core.CotCTPN_MaSKU] = ct.MaSKU; rowCT[core.CotCTPN_TenSanPham] = ct.TenSanPham
			rowCT[core.CotCTPN_DonVi] = ct.DonVi; rowCT[core.CotCTPN_SoLuong] = ct.SoLuong
			rowCT[core.CotCTPN_DonGiaNhap] = ct.DonGiaNhap; rowCT[core.CotCTPN_ThanhTienDong] = ct.ThanhTienDong
			rowCT[core.CotCTPN_GiaVonThucTe] = ct.GiaVonThucTe
			core.PushAppend(shopID, core.TenSheetChiTietPhieuNhap, rowCT)

			for i := 0; i < item.SoLuong; i++ {
				imei := ""
				if i < len(item.Serials) && strings.TrimSpace(item.Serials[i]) != "" {
					imei = strings.TrimSpace(item.Serials[i])
				} else {
					imei = fmt.Sprintf("SN%s%s", time.Now().Format("060102"), core.LayChuoiSoNgauNhien(6))
				}

				sr := &core.SerialSanPham{
					SpreadsheetID: shopID, SerialIMEI: imei, MaSanPham: maSP, MaSKU: item.MaSKU,
					MaNhaCungCap: input.MaNhaCungCap, MaPhieuNhap: pn.MaPhieuNhap, TrangThai: 1, 
					NgayNhapKho: input.NgayNhap, GiaVonNhap: item.DonGiaNhap, MaKho: input.MaKho, NgayCapNhat: nowStr,
				}

				core.CacheSerialSanPham[shopID] = append(core.CacheSerialSanPham[shopID], sr)
				core.CacheMapSerial[core.TaoCompositeKey(shopID, imei)] = sr

				rowSR := make([]interface{}, 19)
				rowSR[core.CotSR_SerialIMEI] = sr.SerialIMEI; rowSR[core.CotSR_MaSanPham] = sr.MaSanPham
				rowSR[core.CotSR_MaSKU] = sr.MaSKU; rowSR[core.CotSR_MaNhaCungCap] = sr.MaNhaCungCap
				rowSR[core.CotSR_MaPhieuNhap] = sr.MaPhieuNhap; rowSR[core.CotSR_TrangThai] = sr.TrangThai
				rowSR[core.CotSR_NgayNhapKho] = sr.NgayNhapKho; rowSR[core.CotSR_GiaVonNhap] = sr.GiaVonNhap
				rowSR[core.CotSR_MaKho] = sr.MaKho; rowSR[core.CotSR_NgayCapNhat] = sr.NgayCapNhat
				core.PushAppend(shopID, core.TenSheetSerial, rowSR)
			}
		}
		core.KhoaHeThong.Unlock()
	}

	loai := "Ghi sổ & Hoàn thành Nhập kho"
	if input.TrangThai == 0 {
		loai = "Đã lưu nháp Phiếu Nhập"
	} else if input.TrangThai == 2 {
		loai = "Đã gửi Yêu cầu Duyệt"
	}
	
	c.JSON(200, gin.H{"status": "ok", "msg": loai + " thành công!", "ma_phieu": pn.MaPhieuNhap})
}

// ============================================================================
// 4. API ĐỔI TRẠNG THÁI PHIẾU (XÓA / KHÔI PHỤC)
// ============================================================================
func API_DoiTrangThaiPhieu(c *gin.Context) {
	shopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")

	if !checkQuyenNhapHang(vaiTro, userID) {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền thực hiện thao tác này!"})
		return
	}

	maPhieu := c.PostForm("ma_phieu_nhap")
	trangThaiStr := c.PostForm("trang_thai")
	
	trangThaiMoi, err := strconv.Atoi(trangThaiStr)
	if err != nil {
		c.JSON(400, gin.H{"status": "error", "msg": "Trạng thái không hợp lệ!"})
		return
	}

	core.GetSheetLock(shopID, core.TenSheetPhieuNhap).Lock()
	defer core.GetSheetLock(shopID, core.TenSheetPhieuNhap).Unlock()

	pn, ok := core.CacheMapPhieuNhap[core.TaoCompositeKey(shopID, maPhieu)]
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Không tìm thấy phiếu này trên hệ thống!"})
		return
	}
	
	if pn.TrangThai == 1 {
		c.JSON(200, gin.H{"status": "error", "msg": "Phiếu đã chốt sổ, không thể xóa hoặc thay đổi!"})
		return
	}

	nguoiThaoTac, _ := core.LayKhachHang(shopID, userID)
	tenNguoiThaoTac := "Hệ thống"
	if nguoiThaoTac != nil { tenNguoiThaoTac = nguoiThaoTac.TenDangNhap }

	nowStr := time.Now().In(time.FixedZone("ICT", 7*3600)).Format("2006-01-02 15:04:05")

	core.KhoaHeThong.Lock()
	pn.TrangThai = trangThaiMoi
	pn.NguoiCapNhat = tenNguoiThaoTac
	pn.NgayCapNhat = nowStr
	core.KhoaHeThong.Unlock()

	core.PushUpdate(shopID, core.TenSheetPhieuNhap, pn.DongTrongSheet, core.CotPN_TrangThai, trangThaiMoi)
	core.PushUpdate(shopID, core.TenSheetPhieuNhap, pn.DongTrongSheet, core.CotPN_NguoiCapNhat, tenNguoiThaoTac)
	core.PushUpdate(shopID, core.TenSheetPhieuNhap, pn.DongTrongSheet, core.CotPN_NgayCapNhat, nowStr)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã cập nhật trạng thái phiếu thành công!"})
}
