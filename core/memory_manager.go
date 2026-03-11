package core

import (
	"log"
	"runtime"
	"runtime/debug"
	"sort"
	"sync"
	"time"

	"app/config"
)

var (
	memoryMutex       sync.Mutex
	ThoiGianTruyCap   = make(map[string]int64) // Lưu thời gian truy cập cuối (Unix nano) của từng Shop
	SoHanhDongDangCho = make(map[string]int)   // Đếm số lệnh đang nằm trong Queue chưa ghi xong
)

// 1. Lễ tân chấm công: Gọi hàm này mỗi khi có Request đập vào Shop
func DanhDauTruyCapShop(shopID string) {
	if shopID == config.BienCauHinh.IdFileSheetMaster || shopID == config.BienCauHinh.IdFileSheetAdmin {
		return // Sếp và Admin là Bất tử, không bị chấm công để xóa
	}
	memoryMutex.Lock()
	ThoiGianTruyCap[shopID] = time.Now().UnixNano()
	memoryMutex.Unlock()
}

// 2. Chốt chặn Queue: Gọi khi bắt đầu tống data vào Queue và khi Queue ghi xong
func TangTaskQueue(shopID string) {
	memoryMutex.Lock()
	SoHanhDongDangCho[shopID]++
	memoryMutex.Unlock()
}
func GiamTaskQueue(shopID string) {
	memoryMutex.Lock()
	SoHanhDongDangCho[shopID]--
	if SoHanhDongDangCho[shopID] < 0 { SoHanhDongDangCho[shopID] = 0 }
	memoryMutex.Unlock()
}

// 3. ĐỌC ĐỒNG HỒ NƯỚC (Kiểm tra xem RAM đang chiếm bao nhiêu)
func kiemTraMucRAM() (uint64, uint64, uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m) // Đọc thông số RAM từ lõi Golang
	
	maxBytes := uint64(config.BienCauHinh.MaxRamMB) * 1024 * 1024
	highBytes := maxBytes * uint64(config.BienCauHinh.HighWatermarkPct) / 100
	lowBytes := maxBytes * uint64(config.BienCauHinh.LowWatermarkPct) / 100
	
	return m.Alloc, highBytes, lowBytes
}

// 4. BÁC LAO CÔNG: Kiểm tra và dọn rác nếu vượt 75%
func KiemTraVaXoaRAMKhiDay() {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()

	currentAlloc, highWatermark, lowWatermark := kiemTraMucRAM()
	
	// Nếu RAM chưa chạm 75%, đi ngủ tiếp, không cần dọn
	if currentAlloc < highWatermark {
		return 
	}

	log.Printf("⚠️ [CẢNH BÁO BỘ NHỚ] RAM chạm mức %d MB. Bắt đầu xả lũ...", currentAlloc/1024/1024)

	// Lập danh sách các Shop đang có trên RAM (Trừ Master và Admin)
	type ShopTruyCap struct {
		ShopID string
		LastAccess int64
	}
	var danhSachLRU []ShopTruyCap
	for id, lastTime := range ThoiGianTruyCap {
		danhSachLRU = append(danhSachLRU, ShopTruyCap{ShopID: id, LastAccess: lastTime})
	}

	// Sắp xếp Shop từ CŨ NHẤT (ít truy cập nhất) lên đầu
	sort.Slice(danhSachLRU, func(i, j int) bool {
		return danhSachLRU[i].LastAccess < danhSachLRU[j].LastAccess
	})

	shopsXoa := 0
	for _, st := range danhSachLRU {
		// Dừng dọn nếu RAM đã tụt xuống dưới mức an toàn 60%
		cur, _, _ := kiemTraMucRAM()
		if cur <= lowWatermark { break }

		// [QUAN TRỌNG] Kiểm tra xem Shop này có đang kẹt lệnh Ghi nào không?
		if SoHanhDongDangCho[st.ShopID] > 0 {
			log.Printf("⏳ Bỏ qua Shop %s vì đang có %d lệnh chờ ghi.", st.ShopID, SoHanhDongDangCho[st.ShopID])
			continue
		}

		// TIẾN HÀNH "TRẢM" KHỎI RAM
		xoaShopKhoiRAM(st.ShopID)
		delete(ThoiGianTruyCap, st.ShopID)
		delete(SoHanhDongDangCho, st.ShopID)
		shopsXoa++
		
		// Ép bộ gom rác của Go chạy ngay lập tức để nhả RAM cho OS (Cloud Run)
		runtime.GC()
		debug.FreeOSMemory() 
	}

	log.Printf("🧹 [DỌN RÁC HOÀN TẤT] Đã xóa %d Shop. RAM hiện tại: %d MB", shopsXoa, kiemTraMucRAMThuCap()/1024/1024)
}

func kiemTraMucRAMThuCap() uint64 { var m runtime.MemStats; runtime.ReadMemStats(&m); return m.Alloc }

// HÀM SÁT THỦ: Hủy diệt toàn bộ dữ liệu của 1 Shop khỏi các biến Map
func xoaShopKhoiRAM(shopID string) {
	StatusMutex.Lock()
	delete(CacheStatusKhachHang, shopID)
	StatusMutex.Unlock()

	delete(CacheKhachHang, shopID)
	
	// XÓA BỘ NHỚ SẢN PHẨM MỚI
	KhoaHeThong.Lock()
	delete(CacheSanPham, shopID)
	KhoaHeThong.Unlock()

	delete(CachePhieuNhap, shopID)
	delete(CachePhieuXuat, shopID)
	delete(CacheSerialSanPham, shopID)
	delete(CacheNhaCungCap, shopID)
	delete(CacheDanhMuc, shopID)
	delete(CacheThuongHieu, shopID)
	
	// Quét và xóa các Composite Key (shopID__xyz)
	for k := range CacheMapKhachHang { if k[:len(shopID)] == shopID { delete(CacheMapKhachHang, k) } }
	
	// Quét dọn bộ nhớ O(1) mới của Sản Phẩm
	KhoaHeThong.Lock()
	for k := range CacheMapSanPham { if k[:len(shopID)] == shopID { delete(CacheMapSanPham, k) } }
	for k := range CacheMapSKU { if k[:len(shopID)] == shopID { delete(CacheMapSKU, k) } }
	KhoaHeThong.Unlock()
}
