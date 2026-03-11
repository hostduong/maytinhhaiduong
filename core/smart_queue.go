package core

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"
)

type SyncNote struct {
	ShopID    string
	TableName string
	Action    string // "APPEND" hoặc "UPDATE"
	ObjectID  string // Mã SP
}

var (
	queueMutex sync.Mutex
	NoteQueue  = make(map[string]*SyncNote) // Map để tự động gộp các lệnh trùng nhau (Debounce)
)

// GhiChuDongBo: Phát giấy Note báo mộng cho Worker
func GhiChuDongBo(shopID, tableName, action, objectID string) {
	key := shopID + "_" + tableName + "_" + objectID
	
	queueMutex.Lock()
	if _, exists := NoteQueue[key]; !exists {
		TangTaskQueue(shopID) // Khóa Shop lại, cấm Lễ tân xóa khỏi RAM
	}
	NoteQueue[key] = &SyncNote{
		ShopID:    shopID,
		TableName: tableName,
		Action:    action,
		ObjectID:  objectID,
	}
	queueMutex.Unlock()
}

// Chạy vòng lặp 5 giây (Đã được gọi ngầm ở main.go)
func KhoiTaoWorkerGhiSheet() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			ProcessQueue()
		}
	}()
}

// Xử lý tờ Note
func ProcessQueue() {
	queueMutex.Lock()
	currentNotes := make(map[string]*SyncNote)
	for k, v := range NoteQueue { currentNotes[k] = v }
	queueMutex.Unlock()

	// 1. [MỚI] Thu thập danh sách Tên Sheet Sản Phẩm hợp lệ từ Cấu Hình RAM
	KhoaHeThong.RLock()
	mapSheetSanPham := make(map[string]bool)
	for _, nganh := range CacheDanhSachNganh {
		if nganh.TenSheet != "" {
			mapSheetSanPham[nganh.TenSheet] = true
		}
	}
	KhoaHeThong.RUnlock()

	// 2. Chạy vòng lặp xử lý
	for key, note := range currentNotes {
		success := false
		
		// [ĐÃ SỬA]: So khớp động 100% với tên Sheet trong file cấu hình, bất chấp tiền tố là gì
		if mapSheetSanPham[note.TableName] {
			success = ghiSanPhamNoSQL(note)
		} else {
			success = true // Bỏ qua các bảng cũ chưa áp dụng Smart Queue
		}

		if success {
			queueMutex.Lock()
			delete(NoteQueue, key) // Xé bỏ Note
			GiamTaskQueue(note.ShopID) // Giải phóng RAM cho phép xóa
			queueMutex.Unlock()
		}
	}
}
func ghiSanPhamNoSQL(note *SyncNote) bool {
	KhoaHeThong.RLock()
	sp, ok := CacheMapSanPham[TaoCompositeKey(note.ShopID, note.ObjectID)]
	KhoaHeThong.RUnlock()

	if !ok || sp == nil { return true } // Bị xóa khỏi RAM trước khi kịp ghi

	b, _ := json.Marshal(sp)
	jsonStr := string(b)

	if note.Action == "APPEND" {
		err := PushAppend(note.ShopID, note.TableName, []interface{}{sp.MaSanPham, jsonStr})
		if err != nil {
			log.Printf("❌ [QUEUE] Lỗi APPEND %s: %v", sp.MaSanPham, err)
			return false // Giữ lại Note để Retry
		}
	} else {
		// UPDATE đúng dòng
		err1 := PushUpdate(note.ShopID, note.TableName, sp.DongTrongSheet, CotProd_MaSanPham, sp.MaSanPham)
		err2 := PushUpdate(note.ShopID, note.TableName, sp.DongTrongSheet, CotProd_DataJSON, jsonStr)
		if err1 != nil || err2 != nil {
			log.Printf("❌ [QUEUE] Lỗi UPDATE %s", sp.MaSanPham)
			return false // Giữ lại Note để Retry
		}
	}
	return true // Xong!
}
