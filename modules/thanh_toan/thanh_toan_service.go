package thanh_toan

import (
	"errors"
	"app/core"
)

type PaymentService struct{}

// GetFinalPrice: Tính toán số tiền thực tế dựa trên Gói và Mã giảm giá
func (s *PaymentService) GetFinalPrice(masterShopID, maGoi, maCode string) (float64, string, *core.GoiDichVu, error) {
	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()

	// 1. Tìm gói dựa trên cột ma_goi (Cột A) 
	goi, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(masterShopID, maGoi)]
	if !ok || goi.TrangThai != 1 {
		return 0, "", nil, errors.New("Gói dịch vụ không khả dụng")
	}

	// 2. Lấy giá bán (Cột F) 
	giaCuoi := goi.GiaBan
	appliedCode := ""

	// 3. Xử lý mảng mã khuyến mãi (Cột G) 
	if maCode != "" {
		for _, c := range goi.DanhSachCode {
			if c.Code == maCode {
				if c.SoLuong != -1 && c.SoLuong <= 0 {
					return giaCuoi, "", goi, errors.New("Mã giảm giá đã hết lượt dùng")
				}
				giaCuoi -= c.GiamTien // Trừ tiền từ RAM Master
				appliedCode = c.Code
				break
			}
		}
	}

	if giaCuoi < 0 { giaCuoi = 0 }
	return giaCuoi, appliedCode, goi, nil
}
