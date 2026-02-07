package cau_hinh

import "time"

// CẤU HÌNH THỜI GIAN & BẢO MẬT
const (
	// Thời gian sống của phiên làm việc (30 phút)
	ThoiGianHetHanCookie = 30 * time.Minute

	// Thời gian "ân hạn" (Grace Period)
	ThoiGianAnHan = 5 * time.Minute 

	// RATE LIMIT (GIỚI HẠN TỐC ĐỘ)
	GioiHanNguoiDung   = 10   // request / giây
)

// Mapping Cột Trong Sheet NHAN_VIEN
// (Để code Middleware gọi được các biến CotNV_...)
const (
	CotNV_Cookie          = 7 // Cột H
	CotNV_CookieExpired   = 8 // Cột I
)
