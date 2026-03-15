package core

import (
	"sync"
	"app/config"
)

// ==============================================================================
// 1. QUẢN TRỊ KHÓA RAM (GRANULAR MUTEX LOCKING)
// Mô hình: [SpreadsheetID] -> [SheetName] -> RWMutex
// ==============================================================================
var (
	mutexRegistryLock sync.Mutex
	sheetLocks        = make(map[string]map[string]*sync.RWMutex)
	KhoaHeThong       sync.RWMutex // Khóa toàn cục cho các biến config cực kỳ hiếm khi đổi
)

// GetSheetLock: Thuật toán cấp phát ổ khóa độc lập O(1)
func GetSheetLock(shopID, sheetName string) *sync.RWMutex {
	mutexRegistryLock.Lock()
	defer mutexRegistryLock.Unlock()

	if sheetLocks[shopID] == nil {
		sheetLocks[shopID] = make(map[string]*sync.RWMutex)
	}
	if sheetLocks[shopID][sheetName] == nil {
		sheetLocks[shopID][sheetName] = &sync.RWMutex{}
	}
	return sheetLocks[shopID][sheetName]
}

// TaoCompositeKey: Hàm tiện ích sinh ID chuẩn cho các Map dữ liệu
func TaoCompositeKey(shopID, objectID string) string {
	return shopID + "__" + objectID
}

// ==============================================================================
// 2. BỘ NHỚ ĐỊNH TUYẾN TÊN MIỀN & SAAS (SỬ DỤNG CHO MIDDLEWARE TẦNG 3)
// ==============================================================================

// CacheDomainToSheetID: Map siêu tốc O(1) để tra cứu Tên Miền -> Spreadsheet ID
var CacheDomainToSheetID = make(map[string]string)

// ==============================================================================
// 3. BỘ NHỚ DATA ĐA NGƯỜI THUÊ (MULTI-TENANT CACHE)
// ==============================================================================

const (
	FlagEmpty   = 0 // Chưa nạp
	FlagLoading = 1 // Đang nạp
	FlagOK      = 2 // Đã nạp thành công
	FlagError   = 3 // Lỗi nạp
)

var (
	// --- CỜ TRẠNG THÁI (BỨC TƯỜNG LỬA BẢO VỆ RAM) ---
	CacheStatusKhachHang = make(map[string]int)
	StatusMutex          sync.RWMutex

	// --- HỆ THỐNG & PHÂN QUYỀN ---
	CachePhanQuyen    = make(map[string][]*PhanQuyen)
	CacheMapPhanQuyen = make(map[string]*PhanQuyen)          

	// --- KHÁCH HÀNG & NHÂN SỰ ---
	CacheKhachHang    = make(map[string][]*KhachHang)   
	CacheMapKhachHang = make(map[string]*KhachHang)     

	// --- CẤU HÌNH KINH DOANH (MASTER DATA) ---
	CacheDanhMuc      = make(map[string][]*DanhMuc)     
	CacheMapDanhMuc   = make(map[string]*DanhMuc)       
	CacheThuongHieu   = make(map[string][]*ThuongHieu)  
	CacheMapThuongHieu = make(map[string]*ThuongHieu)   
	CacheBienLoiNhuan = make(map[string][]*BienLoiNhuan)

	// --- GÓI DỊCH VỤ SAAS ---
	CacheGoiDichVu    = make(map[string][]*GoiDichVu)
	CacheMapGoiDichVu = make(map[string]*GoiDichVu)
	CacheDongHienTaiCauHinh = make(map[string]int)

	// --- ĐỐI TÁC ---
	CacheNhaCungCap    = make(map[string][]*NhaCungCap) 
	CacheMapNhaCungCap = make(map[string]*NhaCungCap)   

	// ==============================================================
	// [MỚI] BỘ NHỚ SẢN PHẨM NOSQL & EAV META-DATA (DYNAMIC UI)
	// ==============================================================
	
	// 1. Cấu hình Ngành hàng & Thuộc tính (Nạp từ Master)
	CacheDanhSachNganh = make([]ConfigNganhHang, 0)
	CacheMapNganh      = make(map[string]ConfigNganhHang)    // Key: ma_nganh
	CacheThuocTinh     = make(map[string][]ThuocTinhNganh) // Key: ma_nganh -> ds thuộc tính

	// 2. Tổng kho Sản Phẩm (Dùng chung vĩnh viễn)
	// Mô hình: ShopID -> MaNganh -> Danh sách sản phẩm
	CacheSanPham       = make(map[string]map[string][]*ProductJSON) 
	
	// Lookup O(1) đa ngành (Tìm nhanh không cần biết thuộc ngành nào)
	CacheMapSanPham    = make(map[string]*ProductJSON) // Key: ShopID__MaSanPham
	CacheMapSKU        = make(map[string]*ProductSKU)  // Key: ShopID__MaSKU
	// ==============================================================

	// --- KHO & PHIẾU NHẬP ---
	CachePhieuNhap    = make(map[string][]*PhieuNhap)   
	CacheMapPhieuNhap = make(map[string]*PhieuNhap)     

	// --- BÁN HÀNG & PHIẾU XUẤT ---
	CachePhieuXuat    = make(map[string][]*PhieuXuat)
	CacheMapPhieuXuat = make(map[string]*PhieuXuat)

	// --- SERIAL & BẢO HÀNH ---
	CacheSerialSanPham = make(map[string][]*SerialSanPham)
	CacheMapSerial     = make(map[string]*SerialSanPham) 

	// --- GIAO TIẾP ---
	CacheTinNhan      = make(map[string][]*TinNhan)     
)


// ==============================================================================
// THUẬT TOÁN TÍNH TOÁN SAAS (Chỉ đọc từ RAM Admin, cực kỳ tối ưu)
// ==============================================================================

// LayGioiHanSanPhamCuaShop: Quét danh sách Khách Hàng trong Tầng Admin
func LayGioiHanSanPhamCuaShop(shopID string) int {
	adminID := config.BienCauHinh.IdFileSheetAdmin
	
	// Master và Admin miễn nhiễm với giới hạn
	if shopID == adminID || shopID == config.BienCauHinh.IdFileSheetMaster {
		return -1 
	}

	// [LOCK CHUẨN]: Chỉ khóa Sheet Khách Hàng của File Admin để đọc
	lock := GetSheetLock(adminID, TenSheetKhachHang)
	lock.RLock()
	defer lock.RUnlock()

	for _, kh := range CacheKhachHang[adminID] {
		// [ĐÃ FIX]: Sử dụng cấu trúc JSON mới (kh.System.SheetID)
		if kh.System.SheetID == shopID {
			maxSP := 0
			hasActivePlan := false
			
			// Tính tổng giới hạn của tất cả các gói đang active
			for _, goi := range kh.GoiDichVu {
				if goi.TrangThai == "active" {
					hasActivePlan = true
					if goi.MaxSanPham == -1 { return -1 } // Một gói Vô hạn -> Vô hạn
					maxSP += goi.MaxSanPham
				}
			}
			
			if !hasActivePlan { return 0 } // Không có gói cước -> 0 sản phẩm
			return maxSP
		}
	}
	
	return 0 
}
