

// BuyStarterPackage: Thực thi nghiệp vụ mua gói và bẻ lái về Master
func (s *BangGiaService) BuyStarterPackage(masterShopID, userID, maGoi, maCode string) (string, error) {
	// Sử dụng dấu "_" cho biến codeHopLe để tránh lỗi "declared and not used"
	finalPrice, _, goi, err := s.paySvc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil { return "", err }

	// Chỉ cho phép gói 0đ (Dùng thử hoặc giảm giá 100%)
	if finalPrice > 0 {
		return "", fmt.Errorf("Gói này yêu cầu thanh toán %v VNĐ. Cổng thanh toán đang bảo trì.", finalPrice)
	}

	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok { return "", errors.New("Không tìm thấy thông tin tài khoản") }

	// Bóc tách giới hạn tài nguyên từ JSON
	var limits map[string]interface{}
	_ = json.Unmarshal([]byte(goi.GioiHanJson), &limits)
	maxSP, _ := limits["max_san_pham"].(float64)
	maxNV, _ := limits["max_nhan_vien"].(float64)

	newPlan := core.PlanInfo{
		MaGoi:       goi.MaGoi,
		TenGoi:      goi.TenGoi,
		LoaiGoi:     goi.LoaiGoi, // "STARTER"
		NgayHetHan:  time.Now().AddDate(0, 0, goi.ThoiHanNgay).Format("2006-01-02 15:04:05"),
		TrangThai:   "active",
		MaxSanPham:  int(maxSP),
		MaxNhanVien: int(maxNV),
	}

	// Cập nhật RAM và đẩy Queue
	lock := core.GetSheetLock(masterShopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.GoiDichVu = []core.PlanInfo{newPlan}
	jsonBytes, _ := json.Marshal(kh.GoiDichVu)
	currentRow, tenDangNhap := kh.DongTrongSheet, kh.TenDangNhap
	lock.Unlock()

	core.PushUpdate(masterShopID, core.TenSheetKhachHang, currentRow, core.CotKH_GoiDichVuJson, string(jsonBytes))

	// Kích hoạt hạ tầng chạy ngầm
	go s.KhoiTaoHaTangSubdomain(tenDangNhap)

	return "https://www.99k.vn/admin/database", nil
}
