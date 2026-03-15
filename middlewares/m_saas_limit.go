package middlewares

import (
	"net/http"
	"app/core"
	"fmt"
	"net/http"
	"sync"
	"time"	

	"github.com/gin-gonic/gin"
)

// CheckSaaSLimit: Trạm kiểm soát tài nguyên gói cước
func CheckSaaSLimit(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID := c.GetString("SHOP_ID")
		
		if resourceType == "san_pham" {
			maxSP := core.LayGioiHanSanPhamCuaShop(shopID)
			if maxSP != -1 {
				// Khóa bộ nhớ siêu tốc O(1) mới để đếm tổng sản phẩm
				lock := core.GetSheetLock(shopID, "PRODUCTS_CACHE")
				lock.RLock()
				currentCount := 0
				if nganhMap, exists := core.CacheSanPham[shopID]; exists {
					for _, ds := range nganhMap {
						currentCount += len(ds)
					}
				}
				lock.RUnlock()

				if currentCount >= maxSP {
					TuChoiTruyCap(c, http.StatusPaymentRequired, "Gói dịch vụ của bạn đã đạt giới hạn tối đa. Vui lòng nâng cấp gói cước!")
					return
				}
			}
		}
		
		c.Next()
	}
}



type blockRecord struct {
	Attempts  int
	LockUntil int64
}

var (
	bruteForceMap = make(map[string]*blockRecord)
	bfMutex       sync.Mutex
)

// BruteForceDefender: Chặn đứng các đợt tấn công vét cạn (Brute-force) mã PIN/OTP
func BruteForceDefender(maxAttempts int, lockMinutes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		dinhDanh := c.PostForm("dinh_danh")
		key := ip + "_" + dinhDanh // Gắn chặt IP với ID đang cố tấn công

		bfMutex.Lock()
		record, exists := bruteForceMap[key]
		if !exists {
			record = &blockRecord{}
			bruteForceMap[key] = record
		}

		now := time.Now().Unix()

		// 1. Kiểm tra xem có đang trong thời gian bị phạt (Lockout) không
		if record.LockUntil > now {
			bfMutex.Unlock()
			remain := (record.LockUntil - now) / 60
			if remain == 0 { remain = 1 }
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": "error",
				"msg":    fmt.Sprintf("Phát hiện truy cập bất thường! Hệ thống đang khóa bảo vệ. Vui lòng thử lại sau %d phút.", remain),
			})
			c.Abort()
			return
		}

		// 2. Nếu đã hết hạn khóa -> Reset lại bộ đếm
		if record.LockUntil != 0 && record.LockUntil <= now {
			record.Attempts = 0
			record.LockUntil = 0
		}

		// 3. Tăng số lần thử. Nếu vượt rào -> Kích hoạt khóa
		record.Attempts++
		if record.Attempts > maxAttempts {
			record.LockUntil = now + (lockMinutes * 60)
			bfMutex.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": "error",
				"msg":    fmt.Sprintf("Bạn đã thao tác sai quá %d lần. Để bảo vệ dữ liệu, hệ thống tạm khóa chức năng này trong %d phút!", maxAttempts, lockMinutes),
			})
			c.Abort()
			return
		}
		bfMutex.Unlock()

		c.Next()
	}
}
