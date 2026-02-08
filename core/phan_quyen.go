package core

import (
	"strings"
	"app/cau_hinh"
)

// =============================================================
// 1. KHO LƯU TRỮ QUYỀN (Ma trận 2 chiều)
// =============================================================
// Cấu trúc: Map[VaiTro][MaChucNang] -> true/false
// Ví dụ: _Map_Quyen["sale"]["product.view"] = true
var _Map_Quyen map[string]map[string]bool

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

	// Khởi tạo map mới
	tempMap := make(map[string]map[string]bool)

	// A. XỬ LÝ TIÊU ĐỀ (Dòng 1) ĐỂ TÌM CỘT VAI TRÒ
	// Cột 0: ma_chuc_nang, Cột 1: nhom, Cột 2: mo_ta
	// Từ Cột 3 trở đi là các Vai Trò (admin, sale, kho...)
	var danhSachVaiTro []string
	if len(raw) < 1 { return }
	
	header := raw[0]
	cotBatDauVaiTro := 3 // Cột D

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

		// Duyệt qua từng cột vai trò để lấy giá trị 1/0
		for j, roleName := range danhSachVaiTro {
			colIndex := cotBatDauVaiTro + j
			
			// Lấy giá trị (1 hoặc true là OK)
			val := layString(row, colIndex)
			isAllow := (val == "1" || strings.ToLower(val) == "true")

			// Gán vào Map
			if isAllow {
				tempMap[roleName][maChucNang] = true
			}
		}
	}

	// Cập nhật vào biến toàn cục (Thread-safe xử lý khi gọi)
	KhoaHeThong.Lock()
	_Map_Quyen = tempMap
	KhoaHeThong.Unlock()
}

// =============================================================
// 3. HÀM KIỂM TRA QUYỀN (Core Logic)
// =============================================================
func KiemTraQuyen(vaiTro string, maChucNang string) bool {
	// 1. Admin Root luôn luôn đúng (Quyền lực tuyệt đối)
	if vaiTro == "admin_root" {
		return true
	}

	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()

	// 2. Kiểm tra trong Map
	vaiTro = strings.ToLower(vaiTro)
	
	if listQuyen, ok := _Map_Quyen[vaiTro]; ok {
		if allowed, exist := listQuyen[maChucNang]; exist {
			return allowed
		}
	}

	// Mặc định là chặn (Deny All)
	return false
}
