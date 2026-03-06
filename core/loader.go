package core

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"app/config"
)

// ==============================================================================
// 0. TRẠM ĐĂNG KÝ HÀM NẠP (REGISTRY PATTERN) - NƠI DUY NHẤT CẦN SỬA KHI THÊM SHEET
// ==============================================================================

// Bọc 2 hàm này lại vì bản gốc của Sếp đang trả về error, ta cần nó về dạng chuẩn func(string)
func napKhachHangMasterNoErr(id string) { _ = NapKhachHangMaster(id) }
func napKhachHangAdminNoErr(id string)  { _ = NapKhachHangAdmin(id) }

// [TẦNG 1]: Két sắt Master
var CacHamNapMaster = []func(string){
	NapPhanQuyenMaster,
	napKhachHangMasterNoErr,
	NapGoiDichVuMaster,
	NapTinNhanMaster,
}

// [TẦNG 2]: Tổng kho Admin
var CacHamNapAdmin = []func(string){
	NapPhanQuyenAdmin,
	napKhachHangAdminNoErr,
}

// [TẦNG 3]: Cửa hàng bán lẻ (Lazy Load)
var CacHamNapCuaHang = []func(string){
	// NapKhachHangCuaShop, // Sếp mở comment khi code xong
	// NapPhanQuyenCuaShop, // Sếp mở comment khi code xong
	NapMayTinh,
	NapDanhMuc,
	NapThuongHieu,
	NapBienLoiNhuan,
	NapNhaCungCap,
	NapPhieuNhap,
	NapSerial,
	// NapPhieuXuat,        // Mở comment khi cần
	// NapHoaDon,           // Mở comment khi cần
	// NapPhieuThuChi,      // Mở comment khi cần
	// NapPhieuBaoHanh,     // Mở comment khi cần
}

// ==============================================================================
// 1. ĐỘNG CƠ NẠP SONG SONG (THE ENGINE)
// ==============================================================================
func ChayDanhSachNapSongSong(shopID string, danhSachHam []func(string)) {
	var wg sync.WaitGroup
	wg.Add(len(danhSachHam))
	
	for _, hamNap := range danhSachHam {
		go func(f func(string)) {
			defer wg.Done()
			f(shopID) // Đẩy ID vào để hàm tự nạp
		}(hamNap)
	}
	
	wg.Wait()
}

// ==============================================================================
// 2. KHỞI ĐỘNG HỆ THỐNG (Gọi khi bật Server hoặc Đồng bộ Lõi)
// ==============================================================================
func KhoiDongHeThongNapDuLieu() {
	HeThongDangBan = true
	defer func() { HeThongDangBan = false }()
	log.Println("[LOADER] Bắt đầu khởi động Động cơ nạp Master & Admin...")

	masterID := config.BienCauHinh.IdFileSheetMaster
	adminID := config.BienCauHinh.IdFileSheetAdmin

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ChayDanhSachNapSongSong(masterID, CacHamNapMaster)
		log.Println("[LOADER] ✔ Hoàn tất nạp Két sắt MASTER.")
	}()

	go func() {
		defer wg.Done()
		ChayDanhSachNapSongSong(adminID, CacHamNapAdmin)
		log.Println("[LOADER] ✔ Hoàn tất nạp Tổng kho ADMIN.")
	}()

	wg.Wait()
	log.Println("[LOADER] Động cơ sẵn sàng! Hệ thống đã lên mây ☁️")
}

// ==============================================================================
// 3. CƠ CHẾ NẠP ĐỘNG TẦNG 3 (LAZY LOADING KẾT HỢP DỌN RAM)
// ==============================================================================
func NapDuLieuCuaMotShop(shopID string) {
	StatusMutex.RLock()
	status := CacheStatusKhachHang[shopID]
	StatusMutex.RUnlock()

	// Tránh đàn ngựa giẫm đạp: Nếu đã nạp hoặc đang nạp thì bỏ qua
	if status == FlagOK || status == FlagLoading {
		return
	}

	// Đóng cổng, treo biển "Đang tải"
	StatusMutex.Lock()
	CacheStatusKhachHang[shopID] = FlagLoading
	StatusMutex.Unlock()

	log.Printf("⏳ [LAZY LOAD] Đang kéo Cửa hàng (ID: %s) lên RAM...", shopID)

	// Gọi bác lao công kiểm tra RAM (Mức 75%), dọn bớt Shop cũ nếu cần (File: memory_manager.go)
	KiemTraVaXoaRAMKhiDay()

	// Khởi động Động cơ nạp Cửa hàng
	ChayDanhSachNapSongSong(shopID, CacHamNapCuaHang)

	// Mở cổng
	StatusMutex.Lock()
	CacheStatusKhachHang[shopID] = FlagOK
	StatusMutex.Unlock()
	
	log.Printf("✅ [LAZY LOAD] Hoàn tất nạp Cửa hàng (ID: %s).", shopID)
}

// ==============================================================================
// 4. TẦNG MASTER: NẠP KÉT SẮT & QUẢN TRỊ LÕI
// ==============================================================================

func NapPhanQuyenMaster(masterID string) {
	raw, err := LoadSheetData(masterID, TenSheetPhanQuyenMaster)
	if err != nil || len(raw) == 0 { return }
	xulyNhanDuLieuPhanQuyen(masterID, TenSheetPhanQuyenMaster, raw)
}

func NapKhachHangMaster(masterID string) error {
	raw, err := LoadSheetData(masterID, TenSheetKhachHangMaster)
	if err != nil || len(raw) == 0 { return err }
	return xulyNhanDuLieuKhachHang(masterID, TenSheetKhachHangMaster, raw)
}

func NapGoiDichVuMaster(masterID string) {
	raw, err := LoadSheetData(masterID, TenSheetGoiDichVuMaster)
	if err != nil || len(raw) == 0 { return }
	
	list := []*GoiDichVu{}
	for i, r := range raw {
		if i < DongBatDau_GoiDichVu-1 { continue }
		maGoi := LayString(r, CotGDV_MaGoi)
		if maGoi == "" { continue }
		
		gdv := &GoiDichVu{
			SpreadsheetID:      masterID, 
			DongTrongSheet:     i + 1, 
			MaGoi:              maGoi,
			TenGoi:             LayString(r, CotGDV_TenGoi),
			LoaiGoi:            LayString(r, CotGDV_LoaiGoi),
			ThoiHanNgay:        LayIntStr(LayString(r, CotGDV_ThoiHanNgay)),
			ThoiHanHienThi:     LayString(r, CotGDV_ThoiHanHienThi), 
			NhanHienThi:        LayString(r, CotGDV_NhanHienThi),     
			GiaNiemYet:         LayFloat(r, CotGDV_GiaNiemYet),       
			GiaBan:             LayFloat(r, CotGDV_GiaBan),           
			MaCodeKichHoatJson: LayString(r, CotGDV_MaCodeKichHoatJson), 
			GioiHanJson:        LayString(r, CotGDV_GioiHanJson),     
			MoTa:               LayString(r, CotGDV_MoTa),            
			NgayBatDau:         LayString(r, CotGDV_NgayBatDau),      
			NgayKetThuc:        LayString(r, CotGDV_NgayKetThuc),     
			SoLuongConLai:      -1, 
			TrangThai:          LayInt(r, CotGDV_TrangThai),          
			DanhSachCode:       make([]CodeKichHoat, 0),
		}
		slStr := LayString(r, CotGDV_SoLuongConLai)
		if slStr != "" { gdv.SoLuongConLai = LayIntStr(slStr) }
		if gdv.MaCodeKichHoatJson != "" { _ = json.Unmarshal([]byte(gdv.MaCodeKichHoatJson), &gdv.DanhSachCode) }
		list = append(list, gdv)
	}

	lock := GetSheetLock(masterID, TenSheetGoiDichVuMaster)
	lock.Lock(); defer lock.Unlock()
	CacheGoiDichVu[masterID] = list
	for _, g := range list { CacheMapGoiDichVu[TaoCompositeKey(masterID, g.MaGoi)] = g }
}

func NapTinNhanMaster(masterID string) {
	raw, err := LoadSheetData(masterID, TenSheetTinNhanMaster)
	if err != nil || len(raw) == 0 { return }
	
	list := []*TinNhan{}
	for i, r := range raw {
		if i < DongBatDau_TinNhan-1 { continue }
		maTN := LayString(r, CotTN_MaTinNhan)
		if maTN == "" { continue }
		tn := &TinNhan{
			SpreadsheetID: masterID, DongTrongSheet: i + 1, MaTinNhan: maTN,
			LoaiTinNhan: LayString(r, CotTN_LoaiTinNhan), NguoiGuiID: LayString(r, CotTN_NguoiGuiID),
			NguoiNhanID: LayString(r, CotTN_NguoiNhanID), TieuDe: LayString(r, CotTN_TieuDe),
			NoiDung: LayString(r, CotTN_NoiDung), ThamChieuID: LayString(r, CotTN_ThamChieuID),
			ReplyChoID: LayString(r, CotTN_ReplyChoID), NgayTao: LayString(r, CotTN_NgayTao),
		}
		json.Unmarshal([]byte(LayString(r, CotTN_DinhKemJson)), &tn.DinhKem)
		json.Unmarshal([]byte(LayString(r, CotTN_NguoiDocJson)), &tn.NguoiDoc)
		json.Unmarshal([]byte(LayString(r, CotTN_TrangThaiXoa)), &tn.TrangThaiXoa)
		list = append(list, tn)
	}
	lock := GetSheetLock(masterID, TenSheetTinNhanMaster)
	lock.Lock(); defer lock.Unlock()
	CacheTinNhan[masterID] = list
}

// ==============================================================================
// 5. TẦNG ADMIN: NẠP TỔNG KHO CHỦ SHOP (CRM)
// ==============================================================================

func NapPhanQuyenAdmin(adminID string) {
	raw, err := LoadSheetData(adminID, TenSheetPhanQuyenAdmin)
	if err != nil || len(raw) == 0 { return }
	xulyNhanDuLieuPhanQuyen(adminID, TenSheetPhanQuyenAdmin, raw)
}

func NapKhachHangAdmin(adminID string) error {
	StatusMutex.Lock()
	CacheStatusKhachHang[adminID] = FlagLoading
	StatusMutex.Unlock()

	raw, err := LoadSheetData(adminID, TenSheetKhachHangAdmin)
	if err != nil || len(raw) == 0 { 
		StatusMutex.Lock()
		CacheStatusKhachHang[adminID] = FlagError
		StatusMutex.Unlock()
		return fmt.Errorf("không thể đọc dữ liệu Tầng Admin từ Google Sheets")
	}
	return xulyNhanDuLieuKhachHang(adminID, TenSheetKhachHangAdmin, raw)
}

// ==============================================================================
// HÀM XỬ LÝ LÕI DÙNG CHUNG (Tái sử dụng cho cả Master và Admin)
// ==============================================================================

// Xử lý nạp Phân Quyền vào RAM
func xulyNhanDuLieuPhanQuyen(shopID string, sheetName string, raw [][]interface{}) {
	headerIndex, styleIndex := -1, -1
	for i, row := range raw {
		if len(row) > 0 {
			firstCell := strings.TrimSpace(strings.ToLower(LayString(row, 0)))
			if firstCell == "ma_chuc_nang" { headerIndex = i 
			} else if firstCell == "style" || firstCell == "level" { styleIndex = i }
		}
	}
	if headerIndex == -1 { return }

	tempMap := make(map[string]map[string]bool)
	var danhSachVaiTroCuaShop []VaiTroInfo
	header := raw[headerIndex]
	var listMaVaiTro []string

	for i := CotPQ_StartRole; i < len(header); i++ {
		headerText := strings.TrimSpace(LayString(header, i))
		if headerText == "" { continue }
		parts := strings.Split(headerText, "|")
		roleID := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(parts[0])), " ", "_")
		roleName := roleID
		if len(parts) > 1 { roleName = strings.TrimSpace(parts[1]) }

		if roleID != "" {
			listMaVaiTro = append(listMaVaiTro, roleID)
			tempMap[roleID] = make(map[string]bool)
			
			styleCode := 90 
			if styleIndex != -1 {
				if valStr := LayString(raw[styleIndex], i); valStr != "" {
					if parsedVal, err := strconv.Atoi(valStr); err == nil { styleCode = parsedVal }
				}
			}
			
			var lvl, thm int
			if styleCode >= 10 { lvl = styleCode / 10; thm = styleCode % 10 } else {
				lvl = styleCode
				switch lvl {
				case 0: thm = 9; case 1: thm = 4; case 2: thm = 7; case 3: thm = 5
				case 4: thm = 4; case 5: thm = 6; case 6: thm = 2; case 7: thm = 1
				default: thm = 0
				}
			}
			danhSachVaiTroCuaShop = append(danhSachVaiTroCuaShop, VaiTroInfo{ MaVaiTro: roleID, TenVaiTro: roleName, StyleLevel: lvl, StyleTheme: thm })
		}
	}

	for i, row := range raw {
		if i <= headerIndex { continue }
		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" || maChucNang == "style" || maChucNang == "level" { continue }

		for j, roleID := range listMaVaiTro {
			val := LayString(row, CotPQ_StartRole+j)
			if val == "1" || strings.ToLower(val) == "true" { tempMap[roleID][maChucNang] = true }
		}
	}

	lock := GetSheetLock(shopID, sheetName)
	lock.Lock(); defer lock.Unlock()
	CachePhanQuyen[shopID] = tempMap
	CacheDanhSachVaiTro[shopID] = danhSachVaiTroCuaShop
}

// Xử lý nạp Khách Hàng (Tạo Định tuyến Subdomain luôn)
func xulyNhanDuLieuKhachHang(shopID string, sheetName string, raw [][]interface{}) error {
	list := []*KhachHang{}
	
	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := LayString(r, CotKH_MaKhachHang)
		if maKH == "" { continue }

		kh := &KhachHang{
			SpreadsheetID:  shopID,
			DongTrongSheet: i + 1,
			MaKhachHang:    maKH,
			TenDangNhap:    LayString(r, CotKH_TenDangNhap),
			Email:          LayString(r, CotKH_Email),
			MatKhauHash:    LayString(r, CotKH_MatKhauHash),
			MaPinHash:      LayString(r, CotKH_MaPinHash),
			VaiTroQuyenHan: strings.TrimSpace(LayString(r, CotKH_VaiTroQuyenHan)),
			ChucVu:         strings.TrimSpace(LayString(r, CotKH_ChucVu)),
			TrangThai:      LayInt(r, CotKH_TrangThai),
			NguonKhachHang: LayString(r, CotKH_NguonKhachHang),
			TenKhachHang:   LayString(r, CotKH_TenKhachHang),
			DienThoai:      LayString(r, CotKH_DienThoai),
			AnhDaiDien:     LayString(r, CotKH_AnhDaiDien),
			DiaChi:         LayString(r, CotKH_DiaChi),
			NgaySinh:       LayString(r, CotKH_NgaySinh),
			GioiTinh:       LayInt(r, CotKH_GioiTinh),
			MaSoThue:       LayString(r, CotKH_MaSoThue),
			GhiChu:         LayString(r, CotKH_GhiChu),
			NgayTao:        LayString(r, CotKH_NgayTao),
			NguoiCapNhat:   LayString(r, CotKH_NguoiCapNhat),
			NgayCapNhat:    LayString(r, CotKH_NgayCapNhat),
			Inbox:          make([]*TinNhan, 0),
		}

		_ = json.Unmarshal([]byte(LayString(r, CotKH_RefreshTokenJson)), &kh.RefreshTokens)
		if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]TokenInfo) }
		
		_ = json.Unmarshal([]byte(LayString(r, CotKH_DataSheetsJson)), &kh.DataSheets)
		// Chỉ kết nối DB riêng nếu ở file Admin và Khách đã setup
		if sheetName == TenSheetKhachHangAdmin && kh.DataSheets.GoogleAuthJson != "" && kh.DataSheets.SpreadsheetID != "" {
			KetNoiGoogleSheetRieng(kh.DataSheets.SpreadsheetID, kh.DataSheets.GoogleAuthJson)
		}

		_ = json.Unmarshal([]byte(LayString(r, CotKH_GoiDichVuJson)), &kh.GoiDichVu)
		if kh.GoiDichVu == nil { kh.GoiDichVu = make([]PlanInfo, 0) } 
		
		_ = json.Unmarshal([]byte(LayString(r, CotKH_CauHinhJson)), &kh.CauHinh)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_MangXaHoiJson)), &kh.MangXaHoi)
		_ = json.Unmarshal([]byte(LayString(r, CotKH_ViTienJson)), &kh.ViTien)

		list = append(list, kh)
	}

	lock := GetSheetLock(shopID, sheetName)
	lock.Lock(); defer lock.Unlock()
	CacheKhachHang[shopID] = list
	
	// Khóa Global để cập nhật bản đồ Tên Miền an toàn
	KhoaHeThong.Lock()
	for _, kh := range list {
		CacheMapKhachHang[TaoCompositeKey(shopID, kh.MaKhachHang)] = kh
		
		// Map Subdomain (Chỉ áp dụng cho Chủ shop ở file Admin)
		if sheetName == TenSheetKhachHangAdmin && kh.DataSheets.SpreadsheetID != "" {
			
			// 1. Tên đăng nhập mặc định (Mức ưu tiên 3)
			if kh.TenDangNhap != "" {
				subDefault := kh.TenDangNhap + ".99k.vn"
				CacheDomainToSheetID[subDefault] = kh.DataSheets.SpreadsheetID
				CacheDomainToSheetID[kh.TenDangNhap] = kh.DataSheets.SpreadsheetID
			}

			// 2. Subdomain cấp riêng (Mức ưu tiên 2 - Ghi đè ưu tiên 3)
			if kh.CauHinh.Subdomain != "" {
				CacheDomainToSheetID[kh.CauHinh.Subdomain] = kh.DataSheets.SpreadsheetID
			}

			// 3. Custom Domain chính chủ (Mức ưu tiên 1 - Lệnh tối cao)
			if kh.CauHinh.CustomDomain != "" {
				CacheDomainToSheetID[kh.CauHinh.CustomDomain] = kh.DataSheets.SpreadsheetID
			}
		}
	}
	KhoaHeThong.Unlock()

	StatusMutex.Lock()
	CacheStatusKhachHang[shopID] = FlagOK
	StatusMutex.Unlock()

	return nil
}

// 3. NẠP NHÀ CUNG CẤP (FULL CỘT)
func NapNhaCungCap(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	raw := napDataGeneric(shopID, TenSheetNhaCungCap, nil)
	if raw == nil { return }
	list := []*NhaCungCap{}
	for i, r := range raw {
		if i < DongBatDau_NhaCungCap-1 { continue }
		maNCC := LayString(r, CotNCC_MaNhaCungCap)
		if maNCC == "" { continue }
		ncc := &NhaCungCap{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaNhaCungCap: maNCC,
			TenNhaCungCap: LayString(r, CotNCC_TenNhaCungCap), MaSoThue: LayString(r, CotNCC_MaSoThue),
			DienThoai: LayString(r, CotNCC_DienThoai), Email: LayString(r, CotNCC_Email),
			KhuVuc: LayString(r, CotNCC_KhuVuc), DiaChi: LayString(r, CotNCC_DiaChi),
			NguoiLienHe: LayString(r, CotNCC_NguoiLienHe), AnhDaiDien: LayString(r, CotNCC_AnhDaiDien),
			NganHang: LayString(r, CotNCC_NganHang), NhomNhaCungCap: LayString(r, CotNCC_NhomNhaCungCap), 
			LoaiNhaCungCap: LayString(r, CotNCC_LoaiNhaCungCap), DieuKhoanThanhToan: LayString(r, CotNCC_DieuKhoanThanhToan), 
			ChietKhauMacDinh: LayFloat(r, CotNCC_ChietKhauMacDinh), HanMucCongNo: LayFloat(r, CotNCC_HanMucCongNo), 
			CongNoDauKy: LayFloat(r, CotNCC_CongNoDauKy), TongMua: LayFloat(r, CotNCC_TongMua), 
			NoCanTra: LayFloat(r, CotNCC_NoCanTra), ThongTinThemJson: LayString(r, CotNCC_ThongTinThemJson), 
			TrangThai: LayInt(r, CotNCC_TrangThai), GhiChu: LayString(r, CotNCC_GhiChu), 
			NguoiTao: LayString(r, CotNCC_NguoiTao), NgayTao: LayString(r, CotNCC_NgayTao), 
			NgayCapNhat: LayString(r, CotNCC_NgayCapNhat),
		}
		list = append(list, ncc)
	}
	lock := GetSheetLock(shopID, TenSheetNhaCungCap)
	lock.Lock()
	defer lock.Unlock()
	CacheNhaCungCap[shopID] = list
	for _, ncc := range list {
		CacheMapNhaCungCap[TaoCompositeKey(shopID, ncc.MaNhaCungCap)] = ncc
	}
}

// 4. NẠP MÁY TÍNH (FULL TỪNG SKUS)
func NapMayTinh(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	raw := napDataGeneric(shopID, TenSheetMayTinh, nil)
	if raw == nil { return }
	list := []*SanPhamMayTinh{}
	for i, r := range raw {
		if i < DongBatDau_SanPhamMayTinh-1 { continue }
		maSP := LayString(r, CotPC_MaSanPham)
		if maSP == "" { continue }
		sp := &SanPhamMayTinh{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaSanPham: maSP,
			TenSanPham: LayString(r, CotPC_TenSanPham), TenRutGon: LayString(r, CotPC_TenRutGon),
			Slug: LayString(r, CotPC_Slug), MaSKU: LayString(r, CotPC_MaSKU),
			TenSKU: LayString(r, CotPC_TenSKU), SKUChinh: LayInt(r, CotPC_SKUChinh),
			TrangThai: LayInt(r, CotPC_TrangThai), MaDanhMuc: LayString(r, CotPC_MaDanhMuc),
			MaThuongHieu: LayString(r, CotPC_MaThuongHieu), DonVi: LayString(r, CotPC_DonVi),
			MauSac: LayString(r, CotPC_MauSac), KhoiLuong: LayFloat(r, CotPC_KhoiLuong),
			KichThuoc: LayString(r, CotPC_KichThuoc), UrlHinhAnh: LayString(r, CotPC_UrlHinhAnh),
			ThongSoHTML: LayString(r, CotPC_ThongSoHTML), MoTaHTML: LayString(r, CotPC_MoTaHTML),
			BaoHanh: LayString(r, CotPC_BaoHanh), TinhTrang: LayString(r, CotPC_TinhTrang),
			GiaNhap: LayFloat(r, CotPC_GiaNhap), PhanTramLai: LayFloat(r, CotPC_PhanTramLai),
			GiaNiemYet: LayFloat(r, CotPC_GiaNiemYet), PhanTramGiam: LayFloat(r, CotPC_PhanTramGiam),
			SoTienGiam: LayFloat(r, CotPC_SoTienGiam), GiaBan: LayFloat(r, CotPC_GiaBan),
			GhiChu: LayString(r, CotPC_GhiChu), NguoiTao: LayString(r, CotPC_NguoiTao),
			NgayTao: LayString(r, CotPC_NgayTao), NguoiCapNhat: LayString(r, CotPC_NguoiCapNhat),
			NgayCapNhat: LayString(r, CotPC_NgayCapNhat),
		}
		list = append(list, sp)
	}

	lock := GetSheetLock(shopID, TenSheetMayTinh)
	lock.Lock()
	defer lock.Unlock()
	CacheSanPhamMayTinh[shopID] = list
	
	for k := range CacheGroupSanPhamMayTinh {
		if strings.HasPrefix(k, shopID+"__") { delete(CacheGroupSanPhamMayTinh, k) }
	}
	for _, sp := range list {
		CacheMapSKUMayTinh[TaoCompositeKey(shopID, sp.LayIDDuyNhat())] = sp
		kGroup := TaoCompositeKey(shopID, sp.MaSanPham)
		if CacheGroupSanPhamMayTinh[kGroup] == nil {
			CacheGroupSanPhamMayTinh[kGroup] = []*SanPhamMayTinh{}
		}
		CacheGroupSanPhamMayTinh[kGroup] = append(CacheGroupSanPhamMayTinh[kGroup], sp)
	}
}

// 5. NẠP DANH MỤC
func NapDanhMuc(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	raw := napDataGeneric(shopID, TenSheetDanhMuc, nil)
	if raw == nil { return }
	list := []*DanhMuc{}
	for i, r := range raw {
		if i < DongBatDau_DanhMuc-1 { continue }
		maDM := LayString(r, CotDM_MaDanhMuc)
		if maDM == "" { continue }
		dm := &DanhMuc{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaDanhMuc: maDM,
			TenDanhMuc: LayString(r, CotDM_TenDanhMuc), DanhMucMe: LayString(r, CotDM_DanhMucMe),
			ThueVAT: LayFloat(r, CotDM_ThueVAT), LoiNhuan: LayFloat(r, CotDM_LoiNhuan),
			Slot: LayInt(r, CotDM_Slot), TrangThai: LayInt(r, CotDM_TrangThai),
		}
		list = append(list, dm)
	}
	lock := GetSheetLock(shopID, TenSheetDanhMuc)
	lock.Lock(); defer lock.Unlock()
	CacheDanhMuc[shopID] = list
	for _, dm := range list { CacheMapDanhMuc[TaoCompositeKey(shopID, dm.MaDanhMuc)] = dm }
}

// 6. NẠP THƯƠNG HIỆU
func NapThuongHieu(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	raw := napDataGeneric(shopID, TenSheetThuongHieu, nil)
	if raw == nil { return }
	list := []*ThuongHieu{}
	for i, r := range raw {
		if i < DongBatDau_ThuongHieu-1 { continue }
		maTH := LayString(r, CotTH_MaThuongHieu)
		if maTH == "" { continue }
		th := &ThuongHieu{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaThuongHieu: maTH,
			TenThuongHieu: LayString(r, CotTH_TenThuongHieu), LogoUrl: LayString(r, CotTH_LogoUrl),
			MoTa: LayString(r, CotTH_MoTa), TrangThai: LayInt(r, CotTH_TrangThai),
		}
		list = append(list, th)
	}
	lock := GetSheetLock(shopID, TenSheetThuongHieu)
	lock.Lock(); defer lock.Unlock()
	CacheThuongHieu[shopID] = list
	for _, th := range list { CacheMapThuongHieu[TaoCompositeKey(shopID, th.MaThuongHieu)] = th }
}

// 7. NẠP BIÊN LỢI NHUẬN
func NapBienLoiNhuan(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	raw := napDataGeneric(shopID, TenSheetBienLoiNhuan, nil)
	if raw == nil { return }
	list := []*BienLoiNhuan{}
	for i, r := range raw {
		if i < DongBatDau_BienLoiNhuan-1 { continue }
		khungGia := LayFloat(r, CotBLN_KhungGiaNhap)
		if khungGia <= 0 { continue }
		bln := &BienLoiNhuan{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, KhungGiaNhap: khungGia,
			BienLoiNhuan: LayFloat(r, CotBLN_BienLoiNhuan), TrangThai: LayInt(r, CotBLN_TrangThai),
		}
		list = append(list, bln)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].KhungGiaNhap < list[j].KhungGiaNhap })
	var prev float64 = 0
	for _, b := range list { b.GiaTu = prev; prev = b.KhungGiaNhap + 1 }
	
	lock := GetSheetLock(shopID, TenSheetBienLoiNhuan)
	lock.Lock(); defer lock.Unlock()
	CacheBienLoiNhuan[shopID] = list
}


// 9. NẠP KHO HÀNG (Phiếu Nhập + Chi Tiết Phiếu Nhập)
func NapPhieuNhap(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	
	// 1. Nạp Bảng Cha (Header)
	rawPN, errPN := LoadSheetData(shopID, TenSheetPhieuNhap)
	if errPN != nil { return }
	
	listPN := []*PhieuNhap{}
	mapPNLocal := make(map[string]*PhieuNhap) // Map cục bộ để ráp Chi tiết siêu tốc
	
	for i, r := range rawPN {
		if i < 10 { continue } // Bỏ qua 10 dòng đầu (Tiêu đề)
		maPN := LayString(r, CotPN_MaPhieuNhap)
		if maPN == "" { continue }
		
		pn := &PhieuNhap{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, MaPhieuNhap: maPN,
			MaNhaCungCap: LayString(r, CotPN_MaNhaCungCap), MaKho: LayString(r, CotPN_MaKho),
			NgayNhap: LayString(r, CotPN_NgayNhap), ChiTietJson: LayString(r, CotPN_ChiTietJson),
			TrangThai: LayInt(r, CotPN_TrangThai), SoHoaDon: LayString(r, CotPN_SoHoaDon), 
			NgayHoaDon: LayString(r, CotPN_NgayHoaDon), UrlChungTu: LayString(r, CotPN_UrlChungTu), 
			TongTienPhieu: LayFloat(r, CotPN_TongTienPhieu), GiamGiaPhieu: LayFloat(r, CotPN_GiamGiaPhieu), 
			ChiPhiNhap: LayFloat(r, CotPN_ChiPhiNhap), DaThanhToan: LayFloat(r, CotPN_DaThanhToan),
			ConNo: LayFloat(r, CotPN_ConNo), PhuongThucThanhToan: LayString(r, CotPN_PhuongThucThanhToan),
			TrangThaiThanhToan: LayString(r, CotPN_TrangThaiThanhToan), GhiChu: LayString(r, CotPN_GhiChu),
			NguoiTao: LayString(r, CotPN_NguoiTao), NgayTao: LayString(r, CotPN_NgayTao), 
			NguoiDuyet: LayString(r, CotPN_NguoiDuyet), NgayDuyet: LayString(r, CotPN_NgayDuyet),
			NguoiCapNhat: LayString(r, CotPN_NguoiCapNhat), NgayCapNhat: LayString(r, CotPN_NgayCapNhat),
			ChiTiet: make([]*ChiTietPhieuNhap, 0),
		}

		if pn.TrangThai <= 0 && pn.ChiTietJson != "" {
			_ = json.Unmarshal([]byte(pn.ChiTietJson), &pn.ChiTiet)
		}

		listPN = append(listPN, pn)
		mapPNLocal[maPN] = pn
	}

	// 2. Nạp Bảng Con (Chi tiết) và ráp vào Bảng Cha
	rawCTPN, errCTPN := LoadSheetData(shopID, TenSheetChiTietPhieuNhap)
	if errCTPN == nil {
		for i, r := range rawCTPN {
			if i < 10 { continue }
			maPN := LayString(r, CotCTPN_MaPhieuNhap)
			if maPN == "" { continue }
			
			if parent, ok := mapPNLocal[maPN]; ok && parent.TrangThai > 0 {
				ct := &ChiTietPhieuNhap{
					SpreadsheetID: shopID, DongTrongSheet: i + 1, MaPhieuNhap: maPN,
					MaSanPham: LayString(r, CotCTPN_MaSanPham), MaSKU: LayString(r, CotCTPN_MaSKU),
					MaNganhHang: LayString(r, CotCTPN_MaNganhHang), TenSanPham: LayString(r, CotCTPN_TenSanPham),
					DonVi: LayString(r, CotCTPN_DonVi), SoLuong: LayInt(r, CotCTPN_SoLuong),
					DonGiaNhap: LayFloat(r, CotCTPN_DonGiaNhap), VATPercent: LayFloat(r, CotCTPN_VATPercent),
					GiaSauVAT: LayFloat(r, CotCTPN_GiaSauVAT), ChietKhauDong: LayFloat(r, CotCTPN_ChietKhauDong),
					ThanhTienDong: LayFloat(r, CotCTPN_ThanhTienDong), GiaVonThucTe: LayFloat(r, CotCTPN_GiaVonThucTe),
					BaoHanhThang: LayInt(r, CotCTPN_BaoHanhThang), GhiChuDong: LayString(r, CotCTPN_GhiChuDong),
				}
				parent.ChiTiet = append(parent.ChiTiet, ct)
			}
		}
	}

	lock := GetSheetLock(shopID, TenSheetPhieuNhap)
	lock.Lock(); defer lock.Unlock()
	CachePhieuNhap[shopID] = listPN
	for _, pn := range listPN { CacheMapPhieuNhap[TaoCompositeKey(shopID, pn.MaPhieuNhap)] = pn }
}

// 10. NẠP SERIAL SẢN PHẨM
func NapSerial(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	raw, err := LoadSheetData(shopID, TenSheetSerial)
	if err != nil { return }
	
	list := []*SerialSanPham{}
	for i, r := range raw {
		if i < 10 { continue }
		imei := LayString(r, CotSR_SerialIMEI)
		if imei == "" { continue }
		
		sr := &SerialSanPham{
			SpreadsheetID: shopID, DongTrongSheet: i + 1, SerialIMEI: imei,
			MaSanPham: LayString(r, CotSR_MaSanPham), MaSKU: LayString(r, CotSR_MaSKU),
			MaNganhHang: LayString(r, CotSR_MaNganhHang), MaNhaCungCap: LayString(r, CotSR_MaNhaCungCap),
			MaPhieuNhap: LayString(r, CotSR_MaPhieuNhap), MaPhieuXuat: LayString(r, CotSR_MaPhieuXuat),
			TrangThai: LayInt(r, CotSR_TrangThai), BaoHanhNhaCungCap: LayInt(r, CotSR_BaoHanhNhaCungCap),
			HanBaoHanhNhaCungCap: LayString(r, CotSR_HanBaoHanhNhaCungCap), MaKhachHangHienTai: LayString(r, CotSR_MaKhachHangHienTai),
			NgayNhapKho: LayString(r, CotSR_NgayNhapKho), NgayXuatKho: LayString(r, CotSR_NgayXuatKho),
			GiaVonNhap: LayFloat(r, CotSR_GiaVonNhap), KichHoatBaoHanhKhach: LayString(r, CotSR_KichHoatBaoHanhKhach),
			HanBaoHanhKhach: LayString(r, CotSR_HanBaoHanhKhach), MaKho: LayString(r, CotSR_MaKho),
			GhiChu: LayString(r, CotSR_GhiChu), NgayCapNhat: LayString(r, CotSR_NgayCapNhat),
		}
		list = append(list, sr)
	}

	lock := GetSheetLock(shopID, TenSheetSerial)
	lock.Lock(); defer lock.Unlock()
	CacheSerialSanPham[shopID] = list
	for _, sr := range list { CacheMapSerial[TaoCompositeKey(shopID, sr.SerialIMEI)] = sr }
}


// ==============================================================================
// HÀM GỐC TƯƠNG TÁC GOOGLE API (KHÔNG ĐƯỢC XÓA)
// ==============================================================================
func napDataGeneric(shopID, sheetName string, target interface{}) [][]interface{} {
	// [ĐÃ SỬA]: Không gán bừa idAdmin nếu shopID rỗng nữa. Phải báo lỗi.
	if shopID == "" { 
        log.Printf("⚠️ Lỗi: Đang cố nạp Sheet %s nhưng không có ShopID!", sheetName)
        return nil 
    }
	raw, err := LoadSheetData(shopID, sheetName)
	if err != nil {
		log.Printf("❌ Lỗi LoadSheetData (%s): %v", sheetName, err)
		return nil
	}
	return raw
}
