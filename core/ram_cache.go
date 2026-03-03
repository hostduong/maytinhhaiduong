package core

import (
	"sync"
)

// ==============================================================================
// 1. QUẢN TRỊ KHÓA RAM (GRANULAR MUTEX LOCKING)
// ==============================================================================
var (
	mutexRegistryLock sync.Mutex
	sheetLocks = make(map[string]map[string]*sync.RWMutex)
)

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
	CachePhanQuyen      = make(map[string]map[string]map[string]bool)
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)               

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

	// --- ĐỐI TÁC ---
	CacheNhaCungCap    = make(map[string][]*NhaCungCap)  
	CacheMapNhaCungCap = make(map[string]*NhaCungCap)    

	// --- SẢN PHẨM (NGÀNH MÁY TÍNH) ---
	CacheSanPhamMayTinh      = make(map[string][]*SanPhamMayTinh)   
	CacheMapSKUMayTinh       = make(map[string]*SanPhamMayTinh)     
	CacheGroupSanPhamMayTinh = make(map[string][]*SanPhamMayTinh)   

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
