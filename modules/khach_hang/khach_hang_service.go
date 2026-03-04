package khach_hang

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/core"
)

type Service struct{}

// Hàm vệ tinh: Bóc tách con số từ chuỗi JSON (Giới hạn tài nguyên)
func getIntFromJson(jsonStr, key string) int {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
		if val, ok := data[key].(float64); ok {
			return int(val)
		}
	}
	return 0
}

func (s *Service) XuLyMuaGoiStarter(masterShopID, userID, maGoi, maCode string) (string, error) {
	// 1. Dò tìm Gói Dịch Vụ trong RAM Master
	core.KhoaHeThong.RLock()
	goiDichVu, ok := core.CacheMapGoiDichVu[core.TaoCompositeKey(masterShopID, maGoi)]
	core.KhoaHeThong.RUnlock()

	if !ok || goiDichVu.LoaiGoi != "STARTER" || goiDichVu.TrangThai != 1 {
		return "", errors.New("Gói cước không hợp lệ hoặc đã ngừng bán")
	}

	// TODO: Tương lai xử lý logic đối chiếu trừ lượt maCode ở đây

	// 2. Tóm cổ Khách Hàng đang yêu cầu mua gói
	kh, ok := core.LayKhachHang(masterShopID, userID)
	if !ok {
		return "", errors.New("Phiên đăng nhập không hợp lệ")
	}

	// 3. Bóc tách & Ép Phẳng JSON (Nước đi quyết định tốc độ O(1))
	maxSanPham := getIntFromJson(goiDichVu.GioiHanJson, "max_san_pham")
	maxNhanVien := getIntFromJson(goiDichVu.GioiHanJson, "max_nhan_vien")
	ngayHetHan := time.Now().AddDate(0, 0, goiDichVu.ThoiHanNgay).Format("2006-01-02 15:04:05")

	newPlan := core.PlanInfo{
		MaGoi:       goiDichVu.MaGoi,
		TenGoi:      goiDichVu.TenGoi,
		LoaiGoi:     "STARTER",
		NgayHetHan:  ngayHetHan,
		TrangThai:   "active",
		MaxSanPham:  maxSanPham,
		MaxNhanVien: maxNhanVien,
	}

	// 4. Nhồi vé vào RAM Khách Hàng (Mutex Lock an toàn tuyệt đối)
	lockKH := core.GetSheetLock(masterShopID, core.TenSheetKhachHang)
	lockKH.Lock()
	
	// Tìm xem khách đã có STARTER chưa để ghi đè (Reset), nếu chưa thì thêm mới
	hasStarter := false
	for i, p := range kh.GoiDichVu {
		if p.LoaiGoi == "STARTER" {
			kh.GoiDichVu[i] = newPlan
			hasStarter = true
			break
		}
	}
	if !hasStarter {
		kh.GoiDichVu = append(kh.GoiDichVu, newPlan)
	}
	
	// Đóng gói lại thành chuỗi JSON để Ghi Sheet
	jsonBytes, _ := json.Marshal(kh.GoiDichVu)
	goiDichVuJsonStr := string(jsonBytes)
	row := kh.DongTrongSheet
	tenDangNhap := kh.TenDangNhap 
	lockKH.Unlock()

	// 5. Ném xuống Hàng Đợi (Queue) để Worker lo việc ghi Google Sheets
	core.PushUpdate(masterShopID, core.TenSheetKhachHang, row, core.CotKH_GoiDichVuJson, goiDichVuJsonStr)

	// 6. Trả về Cú Bẻ Lái Tuyệt Đối (Domain Mapping)
	redirectURL := fmt.Sprintf("https://%s.99k.vn/admin/database", tenDangNhap)
	return redirectURL, nil
}
