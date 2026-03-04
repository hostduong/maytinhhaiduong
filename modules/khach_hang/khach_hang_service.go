package khach_hang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"app/core"
	"app/modules/thanh_toan" // Gọi sang module thanh toán tập trung để tính giá
)

type CustomerService struct {
	paySvc thanh_toan.PaymentService
}

// ==============================================================================
// 1. KÍCH HOẠT HẠ TẦNG SUBDOMAIN (CHẠY NGẦM)
// ==============================================================================

// KhoiTaoHaTangSubdomain: Thực hiện gọi API sang Google Cloud Run để Map Domain
func (s *CustomerService) KhoiTaoHaTangSubdomain(tenDangNhap string) {
	subdomain := fmt.Sprintf("%s.99k.vn", tenDangNhap)
	
	// Cấu trúc Body JSON theo chuẩn Google Cloud Run Domain Mapping API
	payload := map[string]interface{}{
		"apiVersion": "domains.cloudrun.com/v1",
		"kind":       "DomainMapping",
		"metadata": map[string]string{
			"name": subdomain,
		},
		"spec": map[string]interface{}{
			"routeName": "maytinhhaiduong", // Tên Service Cloud Run của sếp
		},
	}
	
	body, _ := json.Marshal(payload)

	// URL API Google Cloud Run (Namespace là Project ID của sếp)
	apiURL := "https://asia-southeast1-run.googleapis.com/apis/domains.cloudrun.com/v1/namespaces/project-47337221-fda1-48c7-b2f/domainmappings"
	
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("[HẠ TẦNG] Lỗi tạo Request cho %s: %v\n", subdomain, err)
		return
	}
	
	// [QUAN TRỌNG]: Token này sếp lấy từ Service Account có quyền Cloud Run Admin
	// req.Header.Set("Authorization", "Bearer " + s.getGoogleAccessToken()) 
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	
	if err != nil {
		fmt.Printf("[HẠ TẦNG] Lỗi kết nối Google Cloud API cho %s: %v\n", subdomain, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("[HẠ TẦNG] Google từ chối tạo Domain %s. Mã lỗi: %d\n", subdomain, resp.StatusCode)
		return
	}
	
	fmt.Printf("[HẠ TẦNG] Kích hoạt tạo Subdomain thành công: %s\n", subdomain)
}

// ==============================================================================
// 2. NGHIỆP VỤ MUA GÓI STARTER
// ==============================================================================

// BuyStarterPackage: Xử lý mua gói, nạp RAM và kích hoạt hạ tầng
func (s *CustomerService) BuyStarterPackage(masterShopID, userID, maGoi, maCode string) (string, error) {
	// Bước 1: Gọi Module Thanh Toán để kiểm tra giá cuối cùng (Zero-Trust)
	finalPrice, codeHopLe, goi, err := s.paySvc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil {
		return "", err
	}

	// Bước 2: Chỉ cho phép đi tiếp nếu giá bằng 0đ (Dùng thử hoặc có mã giảm 100%)
	if finalPrice > 0 {
		return "", fmt.Errorf("Gói này yêu cầu thanh toán %v VNĐ. Cổng thanh toán đang bảo trì.", finalPrice)
	}

	// Bước 3: Tìm hồ sơ khách hàng trong RAM Master
	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		return "", errors.New("Không tìm thấy thông tin tài khoản trên hệ thống")
	}

	// Bước 4: Bóc tách giới hạn tài nguyên từ GioiHanJson (Cột H sheet GOI_DICH_VU)
	// Giả sử JSON: {"max_san_pham": 100, "max_nhan_vien": 5}
	var limits map[string]interface{}
	if err := json.Unmarshal([]byte(goi.GioiHanJson), &limits); err != nil {
		return "", errors.New("Dữ liệu giới hạn gói cước không hợp lệ")
	}
	
	maxSP, _ := limits["max_san_pham"].(float64)
	maxNV, _ := limits["max_nhan_vien"].(float64)

	// Bước 5: Khởi tạo cấu trúc PlanInfo mới (Ép phẳng dữ liệu vào RAM)
	ngayHetHan := time.Now().AddDate(0, 0, goi.ThoiHanNgay).Format("2006-01-02 15:04:05")
	newPlan := core.PlanInfo{
		MaGoi:       goi.MaGoi,
		TenGoi:      goi.TenGoi,
		LoaiGoi:     goi.LoaiGoi, // "STARTER" 
		NgayHetHan:  ngayHetHan,
		TrangThai:   "active",
		MaxSanPham:  int(maxSP),
		MaxNhanVien: int(maxNV),
	}

	// Bước 6: Cập nhật RAM Khách Hàng (Dùng Lock để đảm bảo an toàn dữ liệu)
	lock := core.GetSheetLock(masterShopID, core.TenSheetKhachHang)
	lock.Lock()
	
	// Thay thế hoặc thêm mới gói STARTER (Mỗi khách chỉ có 1 gói nền STARTER)
	updatedPlans := []core.PlanInfo{newPlan}
	kh.GoiDichVu = updatedPlans
	
	// Chuẩn bị dữ liệu JSON để ghi xuống Google Sheets
	jsonBytes, _ := json.Marshal(updatedPlans)
	goiDichVuJsonStr := string(jsonBytes)
	
	currentRow := kh.DongTrongSheet
	tenDangNhap := kh.TenDangNhap
	lock.Unlock()

	// Bước 7: Đẩy vào Hàng đợi (Queue) để Worker ghi xuống Sheet KHACH_HANG (Cột 11 - K) 
	core.PushUpdate(masterShopID, core.TenSheetKhachHang, currentRow, core.CotKH_GoiDichVuJson, goiDichVuJsonStr)

	// Bước 8: Kích hoạt tạo hạ tầng Subdomain chạy ngầm (Goroutine)
	go s.KhoiTaoHaTangSubdomain(tenDangNhap)

	// Bước 9: Trả về URL bẻ lái về Master để khách hàng cài đặt database
	return "https://www.99k.vn/admin/database", nil
}
