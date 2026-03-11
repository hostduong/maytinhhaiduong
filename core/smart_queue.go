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

	for key, note := range currentNotes {
		success := false
		
		if strings.HasPrefix(note.TableName, "SP_") {
			success = ghiSanPhamNoSQL(note)
		} else {
			success = true // Bỏ qua các bảng cũ chưa tích hợp
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
