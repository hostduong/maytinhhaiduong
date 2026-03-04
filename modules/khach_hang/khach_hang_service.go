package khach_hang

import (
	"encoding/json"
	"fmt"
	"time"

	"app/core"
	"app/modules/thanh_toan" // <-- Gọi module thanh toán tập trung
)

type CustomerService struct {
	paySvc thanh_toan.PaymentService // Dùng chung chuyên gia tính tiền
}

func (s *CustomerService) BuyStarterPackage(masterShopID, userID, maGoi, maCode string) (string, error) {
	// HỎI GIÁ TỪ MODULE THANH TOÁN
	finalPrice, _, goi, err := s.paySvc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil { return "", err }

	// ĐÚNG NHƯ SẾP CHỐT: Nếu > 0đ thì chặn lại báo chờ tích hợp cổng thanh toán
	if finalPrice > 0 {
		return "", fmt.Errorf("Gói này cần thanh toán %v₫. Chức năng thanh toán tự động đang được bảo trì.", finalPrice)
	}

	// Xử lý nạp gói vào RAM KHACH_HANG như cũ...
    // (Bóc tách gioi_han_json (Cột H) nhét vào PlanInfo...) 
    // ...
    
	return fmt.Sprintf("https://%s.99k.vn/admin/database", "tendangnhap"), nil
}
