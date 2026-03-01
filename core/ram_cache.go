package core

import (
	"sync"
)

// ==============================================================================
// QUẢN TRỊ KHÓA RAM (GRANULAR MUTEX LOCKING)
// ==============================================================================
var (
	mutexRegistryLock sync.Mutex
	// Cấu trúc: [ShopID][SheetName]*sync.RWMutex
	sheetLocks = make(map[string]map[string]*sync.RWMutex)
)

// GetSheetLock: Lấy ra ổ khóa độc lập của 1 Bảng thuộc 1 Shop
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

// ==============================================================================
// BỘ NHỚ DATA ĐA NGƯỜI THUÊ (MULTI-TENANT CACHE)
// Cấu trúc chuẩn: Map[ShopID] -> Dữ liệu
// ==============================================================================

var (
	// Dữ liệu Hệ thống & Phân quyền
	CachePhanQuyen      = make(map[string]map[string]map[string]bool) // [ShopID][Role][ma_chuc_nang] = true
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)               // [ShopID] -> List

	// Dữ liệu Khách hàng & Nhân sự
	CacheKhachHang    = make(map[string][]*KhachHang)
	CacheMapKhachHang = make(map[string]*KhachHang) // Key = TaoCompositeKey(ShopID, MaKH)

	// Dữ liệu Master
	CacheNhaCungCap = make(map[string][]*NhaCungCap)
	
	// Khai báo tiếp Cache DanhMuc, SanPham, PhieuNhap tại đây...
)

// Hàm tạo khóa gộp (Chống trùng lặp trên Map phẳng)
func TaoCompositeKey(shopID, entityID string) string {
	return shopID + "__" + entityID
}
