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
	NapCauHinhThuocTinh, // [MỚI] Nạp ma trận Meta-data cấu hình EAV
}

// [TẦNG 2]: Tổng kho Admin
var CacHamNapAdmin = []func(string){
	NapPhanQuyenAdmin,
	napKhachHangAdminNoErr,
	NapDanhMuc,
	NapThuongHieu,
	NapBienLoiNhuan,
	NapSanPhamGeneric, // [MỚI] Quái vật nạp NoSQL Đa ngành
}

// [TẦNG 3]: Cửa hàng bán lẻ (Lazy Load)
var CacHamNapCuaHang = []func(string){
	// NapKhachHangCuaShop, // Sếp mở comment khi code xong
	// NapPhanQuyenCuaShop, // Sếp mở comment khi code xong
	NapSanPhamGeneric, // [MỚI] Quái vật nạp NoSQL Đa ngành
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
// [MỚI] NẠP CẤU HÌNH THUỘC TÍNH (EAV META-DATA) - CHẠY TỪ MASTER
// ==============================================================================
func NapCauHinhThuocTinh(masterID string) {
	raw := napDataGeneric(masterID, TenSheetCauHinhThuocTinh, nil)
	if raw == nil || len(raw) == 0 { return }

	var listNganh []ConfigNganhHang
	mapNganh := make(map[string]ConfigNganhHang)
	mapThuocTinh := make(map[string][]ThuocTinhNganh)
	colToNganh := make(map[int]string) // Ánh xạ tọa độ cột -> ma_nganh

	// Bước 1: Quét Dòng 1 (Header chứa JSON cấu hình ngành)
	headerRow := raw[DongBatDau_CauHinhThuocTinh-1]
	for i := CotCHTT_StartNganh; i < len(headerRow); i++ {
		jsonStr := LayString(headerRow, i)
		if jsonStr == "" { continue }
		
		var cfg ConfigNganhHang
		if err := json.Unmarshal([]byte(jsonStr), &cfg); err == nil && cfg.MaNganh != "" {
			listNganh = append(listNganh, cfg)
			mapNganh[cfg.MaNganh] = cfg
			colToNganh[i] = cfg.MaNganh
			mapThuocTinh[cfg.MaNganh] = []ThuocTinhNganh{} // Khởi tạo mảng rỗng
		}
	}

	// Bước 2: Quét các dòng dưới (Bắt đầu từ sau dòng Header)
	for i, row := range raw {
		if i < DongBatDau_CauHinhThuocTinh { continue } 
		
		maTT := LayString(row, CotCHTT_MaThuocTinh)
		if maTT == "" { continue }
		
		tt := ThuocTinhNganh{
			MaThuocTinh:  maTT,
			TenThuocTinh: LayString(row, CotCHTT_TenThuocTinh),
			KieuNhap:     LayString(row, CotCHTT_KieuNhap),
			DonVi:        LayString(row, CotCHTT_DonVi),
		}

		// Ngành nào đánh số "1" thì gắn thuộc tính này vào mảng của ngành đó
		for colIdx, maNganh := range colToNganh {
			val := LayString(row, colIdx)
			if val == "1" {
				mapThuocTinh[maNganh] = append(mapThuocTinh[maNganh], tt)
			}
		}
	}

	KhoaHeThong.Lock()
	defer KhoaHeThong.Unlock()
	CacheDanhSachNganh = listNganh
	CacheMapNganh = mapNganh
	CacheThuocTinh = mapThuocTinh
	log.Println("✅ [LOADER] Đã nạp thành công Ma trận Cấu hình Thuộc tính (EAV).")
}

// ==============================================================================
// [MỚI] NẠP TỔNG KHO SẢN PHẨM JSON (NOSQL - CHẠY CHUNG MỌI NGÀNH)
// ==============================================================================
func NapSanPhamGeneric(shopID string) {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }

	// 1. Thu thập danh sách các Sheet vật lý đang lưu trữ sản phẩm (Từ cấu hình)
	KhoaHeThong.RLock()
	uniqueSheets := make(map[string]bool)
	for _, nganh := range CacheDanhSachNganh {
		if nganh.TenSheet != "" {
			uniqueSheets[nganh.TenSheet] = true
		}
	}
	KhoaHeThong.RUnlock()

	// Khởi tạo các Map tạm để nhào nặn O(1)
	tempCacheSP := make(map[string][]*ProductJSON)
	tempMapSP := make(map[string]*ProductJSON)
	tempMapSKU := make(map[string]*ProductSKU)

	// 2. Chọc vào từng Sheet Vật lý tải JSON lên
	for sheetName := range uniqueSheets {
		raw := napDataGeneric(shopID, sheetName, nil)
		if raw == nil { continue }

		for i, r := range raw {
			if i < DongBatDau_Product-1 { continue } // Bỏ qua Header
			maSP := LayString(r, CotProd_MaSanPham)
			dataJSON := LayString(r, CotProd_DataJSON)
			
			if maSP == "" || dataJSON == "" { continue }

			var sp ProductJSON
			if err := json.Unmarshal([]byte(dataJSON), &sp); err == nil {
				// [BỔ SUNG] Gắn số dòng thực tế để Queue biết đường Update
				sp.SpreadsheetID = shopID
				sp.DongTrongSheet = i + 1
				
				spPtr := &sp
				// Phân lô dữ liệu theo Mã Ngành
				tempCacheSP[sp.MaNganh] = append(tempCacheSP[sp.MaNganh], spPtr)
				
				// Nạp vào Map O(1) siêu tốc
				tempMapSP[TaoCompositeKey(shopID, maSP)] = spPtr
				for sIdx := range sp.SKU {
					skuPtr := &sp.SKU[sIdx]
					tempMapSKU[TaoCompositeKey(shopID, skuPtr.MaSKU)] = skuPtr
				}
			} else {
				log.Printf("⚠️ [LOADER] Lỗi Unmarshal JSON Sản phẩm %s (Shop: %s): %v", maSP, shopID, err)
			}
		}
	}

	// 3. Đổ mẻ trộn vào Bồn chứa RAM chính thức
	lockSP := GetSheetLock(shopID, "PRODUCTS_CACHE")
	lockSP.Lock()
	CacheSanPham[shopID] = tempCacheSP
	lockSP.Unlock()

	// Ghi đè vào Map O(1) Global (Dọn rác cũ trước khi ghi để tránh dính bóng ma)
	KhoaHeThong.Lock()
	for k := range CacheMapSanPham {
		if strings.HasPrefix(k, shopID+"__") { delete(CacheMapSanPham, k) }
	}
	for k := range CacheMapSKU {
		if strings.HasPrefix(k, shopID+"__") { delete(CacheMapSKU, k) }
	}
	for k, v := range tempMapSP { CacheMapSanPham[k] = v }
	for k, v := range tempMapSKU { CacheMapSKU[k] = v }
	KhoaHeThong.Unlock()
}


// ==============================================================================
// 4. TẦNG MASTER: NẠP KÉT SẮT & QUẢN TRỊ LÕI
// ==============================================================================

func NapPhanQuyenMaster(masterID string) {
	raw, err := LoadSheetData(masterID, TenSheetCauHinhMaster)
	if err != nil || len(raw) == 0 { return }
	
	lock := GetSheetLock(masterID, TenSheetCauHinhMaster)
	lock.Lock()
	defer lock.Unlock()

	list := []*PhanQuyen{}
	for i, r := range raw {
		if i < DongBatDau_CauHinh-1 { continue }
		
		maCH := LayString(r, CotCH_MaCauHinh)
		dataJSON := LayString(r, CotCH_DataJSON)
		
		if !strings.HasPrefix(maCH, PrePhanQuyen) { continue }
		if dataJSON == "" { continue }
		
		var pq PhanQuyen
		if err := json.Unmarshal([]byte(dataJSON), &pq); err == nil {
			pq.SpreadsheetID = masterID
			pq.DongTrongSheet = i + 1
			list = append(list, &pq)
		} else {
			log.Printf("⚠️ [LOADER] Lỗi Unmarshal JSON Phân Quyền %s: %v", maCH, err)
		}
	}

	CachePhanQuyen[masterID] = list
	for _, p := range list { 
		CacheMapPhanQuyen[TaoCompositeKey(masterID, p.MaVaiTro)] = p 
	}
	log.Println("✅ [LOADER] Đã nạp Phân Quyền từ Két Sắt NoSQL.")
}
func NapKhachHangMaster(masterID string) error {
	raw, err := LoadSheetData(masterID, TenSheetKhachHangMaster)
	if err != nil || len(raw) == 0 { return err }
	return xulyNhanDuLieuKhachHang(masterID, TenSheetKhachHangMaster, raw)
}


func NapGoiDichVuMaster(masterID string) {
	raw, err := LoadSheetData(masterID, TenSheetCauHinhMaster)
	if err != nil || len(raw) == 0 { return }
	
	lock := GetSheetLock(masterID, TenSheetCauHinhMaster)
	lock.Lock()
	defer lock.Unlock()

	CacheDongHienTaiCauHinh[masterID] = len(raw) 
	
	list := []*GoiDichVu{}
	for i, r := range raw {
		if i < DongBatDau_CauHinh-1 { continue }
		
		maCH := LayString(r, CotCH_MaCauHinh)
		dataJSON := LayString(r, CotCH_DataJSON)
		
		if !strings.HasPrefix(maCH, PreGoiDichVu) { continue }
		if dataJSON == "" { continue }
		
		var gdv GoiDichVu
		if err := json.Unmarshal([]byte(dataJSON), &gdv); err == nil {
			gdv.SpreadsheetID = masterID
			gdv.DongTrongSheet = i + 1
			list = append(list, &gdv)
		} else {
			log.Printf("⚠️ [LOADER] Lỗi Unmarshal JSON Gói %s: %v", maCH, err)
		}
	}

	// [MỚI] Sắp xếp tự động theo trường XepHang
	sort.Slice(list, func(i, j int) bool {
		return list[i].XepHang < list[j].XepHang
	})

	CacheGoiDichVu[masterID] = list
	for _, g := range list { 
		CacheMapGoiDichVu[TaoCompositeKey(masterID, g.MaGoi)] = g 
	}
	log.Println("✅ [LOADER] Đã nạp & sắp xếp Gói dịch vụ từ Két Sắt NoSQL.")
}

// ==============================================================================
// 5. TẦNG ADMIN: NẠP TỔNG KHO CHỦ SHOP (CRM)
// ==============================================================================

func NapPhanQuyenAdmin(adminID string) {
	// Tạm thời để trống chờ nâng cấp Tầng Admin sang chuẩn JSON NoSQL giống Master.
	// (Đã xóa hàm xulyNhanDuLieuPhanQuyen cũ để chống lỗi Build)
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


// Nạp JSON Tin Nhắn NoSQL 2 Cột
func NapTinNhanMaster(masterID string) {
	raw, err := LoadSheetData(masterID, TenSheetTinNhanMaster)
	if err != nil || len(raw) == 0 { return }
	
	list := []*TinNhan{}
	for i, r := range raw {
		if i < DongBatDau_TinNhan-1 { continue }
		maTN := LayString(r, CotTN_MaTinNhan)
		dataJSON := LayString(r, CotTN_DataJSON)
		if maTN == "" { continue }
		
		var tn TinNhan
		if dataJSON != "" {
			if err := json.Unmarshal([]byte(dataJSON), &tn); err != nil {
				log.Printf("⚠️ Lỗi parse JSON Tin Nhắn %s: %v", maTN, err)
				continue
			}
		} else {
			tn.MaTinNhan = maTN 
		}
		
		tn.SpreadsheetID = masterID
		tn.DongTrongSheet = i + 1
		list = append(list, &tn)
	}
	lock := GetSheetLock(masterID, TenSheetTinNhanMaster)
	lock.Lock(); defer lock.Unlock()
	CacheTinNhan[masterID] = list
}

// Nạp Khách Hàng SaaS (NoSQL 2 Cột)
func xulyNhanDuLieuKhachHang(shopID string, sheetName string, raw [][]interface{}) error {
	list := []*KhachHang{}
	
	for i, r := range raw {
		if i < DongBatDau_KhachHang-1 { continue }
		maKH := LayString(r, CotKH_MaKhachHang)
		dataJSON := LayString(r, CotKH_DataJSON)
		
		if maKH == "" { continue }

		var kh KhachHang
		if dataJSON != "" {
			if err := json.Unmarshal([]byte(dataJSON), &kh); err != nil {
				log.Printf("⚠️ Lỗi parse JSON Khách hàng %s: %v", maKH, err)
				continue
			}
		} else {
			kh.MaKhachHang = maKH
		}

		kh.SpreadsheetID = shopID
		kh.DongTrongSheet = i + 1
		if kh.RefreshTokens == nil { kh.RefreshTokens = make(map[string]TenantDeviceToken) }
		if kh.NganHang.TenNganHang == "" { kh.NganHang = TenantNganHang{} } // Đảm bảo khởi tạo
		kh.Inbox = make([]*TinNhan, 0)
		
		if sheetName == TenSheetKhachHangAdmin && kh.System.GoogleAuthJson != "" && kh.System.SheetID != "" {
			KetNoiGoogleSheetRieng(kh.System.SheetID, kh.System.GoogleAuthJson)
		}

		list = append(list, &kh)
	}

	lock := GetSheetLock(shopID, sheetName)
	lock.Lock(); defer lock.Unlock()
	CacheKhachHang[shopID] = list
	
	KhoaHeThong.Lock()
	for _, kh := range list {
		CacheMapKhachHang[TaoCompositeKey(shopID, kh.MaKhachHang)] = kh
		
		if sheetName == TenSheetKhachHangAdmin && kh.System.SheetID != "" {
			if kh.TenDangNhap != "" {
				CacheDomainToSheetID[kh.TenDangNhap+".99k.vn"] = kh.System.SheetID
				CacheDomainToSheetID[kh.TenDangNhap] = kh.System.SheetID
			}
			if kh.Domain.Subdomain != "" { CacheDomainToSheetID[kh.Domain.Subdomain] = kh.System.SheetID }
			if kh.Domain.CustomDomain != "" { CacheDomainToSheetID[kh.Domain.CustomDomain] = kh.System.SheetID }
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
