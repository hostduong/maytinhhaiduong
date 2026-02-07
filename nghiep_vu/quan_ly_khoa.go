package nghiep_vu

import (
	"sync"
)

// ========================================================
// QUẢN LÝ KHÓA ĐỘNG (LOCK MANAGER)
// Giúp tạo ra hàng nghìn ổ khóa riêng biệt dựa trên String Key
// ========================================================

type QuanLyKhoa struct {
	mu       sync.Mutex              // Khóa bảo vệ cái Map chứa khóa (Meta-lock)
	khoaMaps map[string]*sync.RWMutex // Map chứa các ổ khóa: Key -> RWMutex
}

// Biến toàn cục để dùng chung
var BoQuanLyKhoa = &QuanLyKhoa{
	khoaMaps: make(map[string]*sync.RWMutex),
}

// LayKhoa : Hàm quan trọng nhất
// Input: "SheetID__SAN_PHAM"
// Output: Ổ khóa RWMutex dành riêng cho Key đó
func (ql *QuanLyKhoa) LayKhoa(tenKey string) *sync.RWMutex {
	// 1. Khóa bảo vệ map để tìm kiếm an toàn
	ql.mu.Lock()
	defer ql.mu.Unlock()

	// 2. Kiểm tra xem đã có khóa cho Key này chưa
	if khoa, tonTai := ql.khoaMaps[tenKey]; tonTai {
		return khoa // Có rồi thì trả về dùng luôn
	}

	// 3. Nếu chưa có, tạo ổ khóa mới
	khoaMoi := &sync.RWMutex{}
	ql.khoaMaps[tenKey] = khoaMoi
	
	return khoaMoi
}
