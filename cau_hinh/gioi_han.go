package cau_hinh

import "time"

// CẤU HÌNH THỜI GIAN & BẢO MẬT
const (
	// Thời gian sống của phiên làm việc (30 phút)
	ThoiGianHetHanCookie = 30 * time.Minute

	// Thời gian "ân hạn" (Grace Period)
	// Nếu user thao tác khi còn < 5 phút thì tự động gia hạn thêm
	ThoiGianAnHan = 5 * time.Minute 

	// RATE LIMIT (GIỚI HẠN TỐC ĐỘ)
	GioiHanNguoiDung   = 100   // request / giây
)

// [ĐÃ ĐỒNG BỘ CHÍNH XÁC]
// Mapping Cột Trong Sheet KHACH_HANG (Dùng cho cấu hình nhanh nếu cần)
// Lưu ý: Giá trị này phải khớp với app/core/khach_hang.go
const (
	CotKH_Cookie          = 3 // Cột D (Index 3)
	CotKH_CookieExpired   = 4 // Cột E (Index 4)
)
