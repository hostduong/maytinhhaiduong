package bang_gia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"app/core"
	"app/modules/thanh_toan"

	"golang.org/x/oauth2/google"
)

type BangGiaService struct {
	paySvc thanh_toan.PaymentService
}

// KhoiTaoHaTangSubdomain: Chạy ngầm gọi API Google Cloud Run bằng ADC
func (s *BangGiaService) KhoiTaoHaTangSubdomain(tenDangNhap string) {
	subdomain := fmt.Sprintf("%s.99k.vn", tenDangNhap)
	
	payload := map[string]interface{}{
		"apiVersion": "domains.cloudrun.com/v1",
		"kind":       "DomainMapping",
		"metadata": map[string]string{"name": subdomain},
		"spec":     map[string]interface{}{"routeName": "maytinhhaiduong"},
	}
	
	body, _ := json.Marshal(payload)
	apiURL := "https://asia-southeast1-run.googleapis.com/apis/domains.cloudrun.com/v1/namespaces/project-47337221-fda1-48c7-b2f/domainmappings"
	
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// =========================================================================
	// [THUẬT TOÁN ĐÁM MÂY]: DÙNG ADC XIN QUYỀN TRỰC TIẾP TỪ GOOGLE METADATA SERVER
	// =========================================================================
	ctx := context.Background()
	// Hàm này sẽ tự động biết nó đang ở Cloud Run và lấy Token cực mượt
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err == nil {
		token, err := creds.TokenSource.Token()
		if err == nil {
			req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		} else {
			fmt.Println("❌ [HẠ TẦNG] ADC tìm thấy quyền nhưng không thể sinh Access Token:", err)
		}
	} else {
		fmt.Println("❌ [HẠ TẦNG] Lỗi ADC! Máy chủ không tìm thấy Default Credentials:", err)
	}
	// =========================================================================

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ [HẠ TẦNG] Lỗi kết nối đến Google Cloud API: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// =========================================================================
	// XỬ LÝ LƯU LƯỢNG THÔNG MINH (KỂ CẢ KHI KHÁCH GIA HẠN/MUA LẠI)
	// =========================================================================
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Printf("✅ [HẠ TẦNG] Kích hoạt thành công Subdomain: %s\n", subdomain)
	} else if resp.StatusCode == 409 {
		// BẮT MÃ 409 CONFLICT: DOMAIN ĐÃ TỒN TẠI, BỎ QUA KHÔNG BÁO LỖI
		fmt.Printf("⚡ [HẠ TẦNG] Subdomain %s đã tồn tại trên hệ thống. Bỏ qua bước tạo mới.\n", subdomain)
	} else {
		// Các lỗi khác (thiếu quyền IAM, sai URL...)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		fmt.Printf("❌ [HẠ TẦNG] Từ chối kích hoạt %s (Mã lỗi: %d). Chi tiết: %s\n", subdomain, resp.StatusCode, buf.String())
	}
}

// BuyStarterPackage: Thực thi nghiệp vụ mua gói và bẻ lái về Master
func (s *BangGiaService) BuyStarterPackage(masterShopID, userID, maGoi, maCode string) (string, error) {
	finalPrice, _, goi, err := s.paySvc.GetFinalPrice(masterShopID, maGoi, maCode)
	if err != nil { return "", err }

	if finalPrice > 0 {
		return "", fmt.Errorf("Gói này yêu cầu thanh toán %v VNĐ. Cổng thanh toán đang bảo trì.", finalPrice)
	}

	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok { return "", errors.New("Không tìm thấy thông tin tài khoản") }

	var limits map[string]interface{}
	_ = json.Unmarshal([]byte(goi.GioiHanJson), &limits)
	maxSP, _ := limits["max_san_pham"].(float64)
	maxNV, _ := limits["max_nhan_vien"].(float64)

	newPlan := core.PlanInfo{
		MaGoi:       goi.MaGoi,
		TenGoi:      goi.TenGoi,
		LoaiGoi:     goi.LoaiGoi, 
		NgayHetHan:  time.Now().AddDate(0, 0, goi.ThoiHanNgay).Format("2006-01-02 15:04:05"),
		TrangThai:   "active",
		MaxSanPham:  int(maxSP),
		MaxNhanVien: int(maxNV),
	}

	lock := core.GetSheetLock(masterShopID, core.TenSheetKhachHang)
	lock.Lock()
	kh.GoiDichVu = []core.PlanInfo{newPlan}
	jsonBytes, _ := json.Marshal(kh.GoiDichVu)
	currentRow, tenDangNhap := kh.DongTrongSheet, kh.TenDangNhap
	lock.Unlock()

	core.PushUpdate(masterShopID, core.TenSheetKhachHang, currentRow, core.CotKH_GoiDichVuJson, string(jsonBytes))

	// Chạy Goroutine tạo Subdomain ngầm để giao diện không bị đơ
	go s.KhoiTaoHaTangSubdomain(tenDangNhap)

	return "https://www.99k.vn/admin/database", nil
}
