package core

import (
	"sync"
)

// ==============================================================================
// 1. QUẢN TRỊ KHÓA RAM (GRANULAR MUTEX LOCKING)
// Tuyệt đối không dùng 1 khóa Global cho toàn hệ thống.
// Mỗi Sheet của mỗi Shop sẽ có một ổ khóa riêng biệt.
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
// 2. BỘ NHỚ DATA ĐA NGƯỜI THUÊ (MULTI-TENANT CACHE)
// Cấu trúc chuẩn: Map[ShopID] -> Dữ liệu
// ==============================================================================

var (
	// --- HỆ THỐNG & PHÂN QUYỀN ---
	CachePhanQuyen      = make(map[string]map[string]map[string]bool) // [ShopID][Role][ma_chuc_nang] = true
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)               // [ShopID] -> List vai trò

	// --- KHÁCH HÀNG & NHÂN SỰ ---
	CacheKhachHang    = make(map[string][]*KhachHang)    // Danh sách theo Shop
	CacheMapKhachHang = make(map[string]*KhachHang)      // Tra cứu nhanh: Key = ShopID__MaKH

	// --- CẤU HÌNH KINH DOANH (MASTER DATA) ---
	CacheDanhMuc      = make(map[string][]*DanhMuc)      // Danh sách danh mục
	CacheMapDanhMuc   = make(map[string]*DanhMuc)        // Tra cứu: Key = ShopID__MaDM
	
	CacheThuongHieu   = make(map[string][]*ThuongHieu)   // Danh sách thương hiệu
	CacheMapThuongHieu = make(map[string]*ThuongHieu)    // Tra cứu: Key = ShopID__MaTH

	CacheBienLoiNhuan = make(map[string][]*BienLoiNhuan) // Khung lợi nhuận theo Shop

	// --- ĐỐI TÁC ---
	CacheNhaCungCap    = make(map[string][]*NhaCungCap)  // Danh sách NCC
	CacheMapNhaCungCap = make(map[string]*NhaCungCap)    // [FIX LỖI BUILD]: Tra cứu NCC nhanh

	// --- SẢN PHẨM (NGÀNH MÁY TÍNH) ---
	CacheSanPhamMayTinh      = make(map[string][]*SanPhamMayTinh)   // Toàn bộ SKU phẳng
	CacheMapSKUMayTinh       = make(map[string]*SanPhamMayTinh)     // Tra cứu SKU: Key = ShopID__MaSKU
	CacheGroupSanPhamMayTinh = make(map[string][]*SanPhamMayTinh)   // Nhóm theo Model: Key = ShopID__MaSP

	// --- KHO & PHIẾU NHẬP ---
	CachePhieuNhap    = make(map[string][]*PhieuNhap)    // Danh sách phiếu nhập
	CacheMapPhieuNhap = make(map[string]*PhieuNhap)      // Tra cứu: Key = ShopID__MaPN

	// --- GIAO TIẾP ---
	CacheTinNhan      = make(map[string][]*TinNhan)      // Danh sách tin nhắn/thông báo theo Shop
)

// ==============================================================================
// 3. TIỆN ÍCH DỮ LIỆU
// ==============================================================================

// TaoCompositeKey: Tạo khóa gộp để lưu trữ trong các Map phẳng (Flat Map)
// Đảm bảo dữ liệu các Shop không bao giờ bị đè lên nhau.
// Lưu ý: Nếu hàm này đã khai báo ở common.go thì bạn có thể xóa đoạn này.
func TaoCompositeKey(shopID, entityID string) string {
	return shopID + "__" + entityID
}
