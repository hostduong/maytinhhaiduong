package core

import (
	"strings"
	"sync"
	"app/cau_hinh"
)

// =============================================================
// 1. KHO LƯU TRỮ QUYỀN (Ma trận 2 chiều)
// =============================================================
// Cấu trúc: Map[VaiTro][MaChucNang] -> true/false
// Ví dụ: _Map_Quyen["sale"]["product.view"] = true
var (
	_Map_Quyen map[string]map[string]bool
	mtxQuyen   sync.RWMutex
)

// =============================================================
// 2. LOGIC NẠP DỮ LIỆU TỪ SHEET
// =============================================================
func NapPhanQuyen(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" {
		targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	// Đọc sheet PHAN_QUYEN
	raw, err := loadSheetData(targetSpreadsheetID, "PHAN_QUYEN")
	if err != nil { return }

	if len(raw) < 2 { return } // Phải có ít nhất Header và 1 dòng dữ liệu

	// Khởi tạo map mới
	tempMap := make(map[string]map[string]bool)

	// A. XỬ LÝ TIÊU ĐỀ (Dòng 1) ĐỂ TÌM CỘT VAI TRÒ
	// Cấu trúc: [0]ma_chuc_nang, [1]nhom, [2]mo_ta, [3...]CÁC VAI TRÒ
	header := raw[0]
	cotBatDauVaiTro := 3 // Cột D (Index 3)

	// Danh sách các vai trò tìm được từ Header
	var danhSachVaiTro []string 

	for i := cotBatDauVaiTro; i < len(header); i++ {
		roleName := strings.ToLower(layString(header, i))
		if roleName != "" {
			danhSachVaiTro = append(danhSachVaiTro, roleName)
			// Khởi tạo map cho vai trò này
			tempMap[roleName] = make(map[string]bool)
		}
	}

	// B. DUYỆT CÁC DÒNG DỮ LIỆU (Từ dòng 2)
	for i := 1; i < len(raw); i++ {
		row := raw[i]
		maChucNang := layString(row, 0) // Cột A
		
		if maChucNang == "" { continue }

		// Duyệt qua từng cột vai trò tương ứng
		for j, roleName := range danhSachVaiTro {
			colIndex := cotBatDauVaiTro + j
			
			// Lấy giá trị (1 hoặc true là OK)
			val := layString(row, colIndex)
			isAllow := (val == "1" || strings.ToLower(val) == "true")

			// Gán vào Map nếu được phép
			if isAllow {
				tempMap[roleName][maChucNang] = true
			}
		}
	}

	// Cập nhật vào biến toàn cục an toàn
	mtxQuyen.Lock()
	_Map_Quyen = tempMap
	mtxQuyen.Unlock()
}

// =============================================================
// 3. HÀM KIỂM TRA QUYỀN (Core Logic)
// =============================================================
func KiemTraQuyen(vaiTro string, maChucNang string) bool {
	// 1. Admin Root (Code cứng) luôn đúng
	if vaiTro == "admin_root" {
		return true
	}

	mtxQuyen.RLock()
	defer mtxQuyen.RUnlock()

	// 2. Chuẩn hóa vai trò về chữ thường để so sánh
	vaiTro = strings.ToLower(vaiTro)
	
	// 3. Tra cứu trong Map
	if listQuyen, ok := _Map_Quyen[vaiTro]; ok {
		if allowed, exist := listQuyen[maChucNang]; exist {
			return allowed
		}
	}

	// Mặc định là CHẶN (Deny All)
	return false
}
