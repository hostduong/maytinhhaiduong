package nghiep_vu

import (
	"strings"
	"sync"
	"app/bo_nho_dem"
)

// [SỬA] Dùng biến của bo_nho_dem, không cần khai báo lại biến CachePhanQuyen ở đây nữa
// Hàm NapDuLieuPhanQuyen cũng đã được chuyển sang bo_nho_dem/cau_hinh.go
// Ở đây chỉ giữ lại hàm KiemTraQuyen

var mtxPhanQuyen sync.RWMutex 
// Lưu ý: CachePhanQuyen giờ nằm ở bo_nho_dem.CachePhanQuyen

func KiemTraQuyen(vaiTroHienTai string, maChucNang string) bool {
	if vaiTroHienTai == "admin_root" || vaiTroHienTai == "quan_tri_vien_cap_cao" {
		return true
	}

	// [CẦN LOCK CỦA STORE MỚI]
	// Vì CachePhanQuyen giờ là biến toàn cục của bo_nho_dem, 
	// nhưng nó được nạp 1 lần lúc boot, ít thay đổi.
	// Tuy nhiên để an toàn, vẫn nên RLock nếu nó có thể bị reload.
	// Hàm reload LamMoiHeThong có gọi NapDuLieuPhanQuyen, nên phải RLock từ KhoaHeThong.
	
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()

	roleKey := strings.ToUpper(vaiTroHienTai)

	if listQuyen, ok := bo_nho_dem.CachePhanQuyen[roleKey]; ok {
		if allowed, exist := listQuyen[maChucNang]; exist {
			return allowed
		}
	}

	return false
}
