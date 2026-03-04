package khach_hang

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/core"
)

type Service struct{}

// TinhGiaCuoiCung: Hàm lõi để tính toán số tiền khách phải trả
func (s *Service) TinhGiaCuoiCung(masterShopID, maGoi, maCode string) (float64, string, *core.GoiDichVu, error) {
	core.KhoaHeThong.RLock()
	defer core.KhoaHeThong.RUnlock()

	// 1. Tìm gói trong RAM Master
	goi, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(masterShopID, maGoi)]
	if !ok || goi.TrangThai != 1 {
		return 0, "", nil, errors.New("Gói dịch vụ không tồn tại hoặc đã ngừng kinh doanh")
	}

	giaCuoi := goi.GiaBan
	codeApDung := ""

	// 2. Nếu khách có nhập mã, kiểm tra trong mảng DanhSachCode đã nạp RAM
	if maCode != "" {
		for _, c := range goi.DanhSachCode {
			if c.Code == maCode {
				// Kiểm tra lượt dùng (nếu có giới hạn)
				if c.SoLuong != -1 && c.SoLuong <= 0 {
					return giaCuoi, "", goi, errors.New("Mã giảm giá này đã hết lượt sử dụng")
				}
				giaCuoi -= c.GiamTien
				codeApDung = c.Code
				break
			}
		}
	}

	if giaCuoi < 0 { giaCuoi = 0 }
	return giaCuoi, codeApDung, goi, nil
}

func (s *Service) ThucThiKichHoatGoi(masterShopID, userID, maGoi, maCode string) (string, error) {
	// GỌI LẠI HÀM TÍNH GIÁ (Bảo mật 2 lớp, không tin dữ liệu từ Client gửi lên)
	giaCuoi, codeHopLe, goi, err := s.TinhGiaCuoiCung(masterShopID, maGoi, maCode)
	if err != nil { return "", err }

	if giaCuoi > 0 {
		return "", fmt.Errorf("Gói này yêu cầu thanh toán %v VNĐ. Vui lòng liên hệ Admin.", giaCuoi)
	}

	// Lấy hồ sơ khách hàng
	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok { return "", errors.New("Không tìm thấy thông tin tài khoản") }

	// Ép phẳng dữ liệu giới hạn từ GioiHanJson
	var limits map[string]interface{}
	_ = json.Unmarshal([]byte(goi.GioiHanJson), &limits)
	
	maxSP, _ := limits["max_san_pham"].(float64)
	maxNV, _ := limits["max_nhan_vien"].(float64)

	newPlan := core.PlanInfo{
		MaGoi:       goi.MaGoi,
		TenGoi:      goi.TenGoi,
		LoaiGoi:     goi.LoaiGoi, // STARTER
		NgayHetHan:  time.Now().AddDate(0, 0, goi.ThoiHanNgay).Format("2006-01-02 15:04:05"),
		TrangThai:   "active",
		MaxSanPham:  int(maxSP),
		MaxNhanVien: int(maxNV),
	}

	// Ghi vào RAM và đẩy Queue
	lock := core.GetSheetLock(masterShopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.GoiDichVu = []core.PlanInfo{newPlan} // STARTER luôn là gói nền duy nhất
	jsonStr, _ := json.Marshal(kh.GoiDichVu)
	row, tenDangNhap := kh.DongTrongSheet, kh.TenDangNhap
	lock.Unlock()

	core.PushUpdate(masterShopID, core.TenSheetKhachHang, row, core.CotKH_GoiDichVuJson, string(jsonStr))

	// Trả về URL bẻ lái tuyệt đối
	return fmt.Sprintf("https://%s.99k.vn/admin/database", tenDangNhap), nil
}
