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
// 0. BỘ CHỈ HUY TỔNG (Hàm kích hoạt khi khởi động Server hoặc bấm Đồng bộ)
// ==============================================================================
func KhoiDongHeThongNapDuLieu() {
	HeThongDangBan = true
	defer func() { HeThongDangBan = false }()
	log.Println("[LOADER] Bắt đầu nạp dữ liệu Tầng Master & Admin...")

	masterID := config.BienCauHinh.IdFileSheetMaster
	adminID := config.BienCauHinh.IdFileSheetAdmin

	// Dùng WaitGroup để nạp song song cho nhanh (Tăng tốc độ khởi động x3 lần)
	var wg sync.WaitGroup
	wg.Add(2)

	// --- LUỒNG 1: Nạp Két Sắt Master ---
	go func() {
		defer wg.Done()
		NapPhanQuyenMaster(masterID)
		NapKhachHangMaster(masterID)
		NapGoiDichVuMaster(masterID)
		NapTinNhanMaster(masterID)
		log.Println("[LOADER] ✔ Nạp xong Tầng MASTER.")
	}()

	// --- LUỒNG 2: Nạp Tổng kho Admin ---
	go func() {
		defer wg.Done()
		NapPhanQuyenAdmin(adminID)
		NapKhachHangAdmin(adminID)
		log.Println("[LOADER] ✔ Nạp xong Tầng ADMIN.")
	}()

	wg.Wait()
	log.Println("[LOADER] Hệ thống đã sẵn sàng!")
}

// ==============================================================================
// 1. TẦNG MASTER: NẠP KÉT SẮT & QUẢN TRỊ LÕI
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
// 2. TẦNG ADMIN: NẠP TỔNG KHO CHỦ SHOP (CRM)
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
		if sheetName == TenSheetKhachHangAdmin {
			if kh.TenDangNhap != "" && kh.DataSheets.SpreadsheetID != "" {
				subdomain := kh.TenDangNhap + ".99k.vn"
				CacheDomainToSheetID[subdomain] = kh.DataSheets.SpreadsheetID
				CacheDomainToSheetID[kh.TenDangNhap] = kh.DataSheets.SpreadsheetID
			}
			if kh.CauHinh.CustomDomain != "" && kh.DataSheets.SpreadsheetID != "" {
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
