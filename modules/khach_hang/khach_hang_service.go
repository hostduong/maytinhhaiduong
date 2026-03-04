package khach_hang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"app/core"
)

// KhoiTaoHaTangSubdomain: Chạy ngầm để map domain vào Cloud Run
func (s *CustomerService) KhoiTaoHaTangSubdomain(tenDangNhap string) {
	subdomain := fmt.Sprintf("%s.99k.vn", tenDangNhap)
	
	// 1. Kiểm tra xem Subdomain này đã tồn tại trong hệ thống chưa (Tránh tạo trùng)
	// (Logic này có thể kiểm tra trong RAM Master hoặc gọi GET sang Google API)

	// 2. Cấu trúc Body JSON gửi sang Google Cloud Run Admin API
	// Tài liệu: https://cloud.google.com/run/docs/reference/rest/v1/namespaces.domainmappings/create
	payload := map[string]interface{}{
		"apiVersion": "domains.cloudrun.com/v1",
		"kind":       "DomainMapping",
		"metadata": map[string]string{
			"name": subdomain, // Ví dụ: duongpc.99k.vn
		},
		"spec": map[string]interface{}{
			"routeName": "service-99k-core", // Tên Service Cloud Run hiện tại của sếp
		},
	}
	
	body, _ := json.Marshal(payload)

	// 3. Gọi API sang Google (Cần có Service Account Token)
	// [LƯU Ý]: Đây là đoạn sếp cần cấu hình PROJECT_ID và REGION của sếp
	apiURL := "https://asia-southeast1-run.googleapis.com/apis/domains.cloudrun.com/v1/namespaces/99k-project/domainmappings"
	
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	
	// Gắn Token xác thực (Sếp sẽ lấy token này từ file JSON Service Account)
	// req.Header.Set("Authorization", "Bearer " + getGoogleAccessToken())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	
	if err != nil || resp.StatusCode >= 400 {
		// Ghi log lỗi vào RAM Master để Admin theo dõi nếu tạo thất bại
		fmt.Printf("Lỗi tạo Subdomain %s: %v\n", subdomain, err)
		return
	}
	
	fmt.Printf("Kích hoạt tạo Subdomain thành công: %s\n", subdomain)
}

// Cập nhật lại hàm BuyStarterPackage để kích hoạt chạy ngầm
func (s *CustomerService) BuyStarterPackage(masterShopID, userID, maGoi, maCode string) (string, error) {
	// ... (Phần logic tính giá và kiểm tra 0đ của sếp) ...

	// Sau khi mọi thứ OK, ghi vào RAM và Queue xong 
	// [BƯỚC CHỐT]: Kích hoạt tạo hạ tầng chạy ngầm
	go s.KhoiTaoHaTangSubdomain(kh.TenDangNhap) 

	// Bẻ lái về trang cài đặt Database trên tên miền Master
	return "https://www.99k.vn/admin/database", nil
}
