package core

import (
	"strings"
	"sync"
	"app/cau_hinh"
)

// =============================================================
// 1. CẤU HÌNH DÒNG & CỘT
// =============================================================
const (
	DongBatDau_PhanQuyen = 11
	
	CotPQ_MaChucNang = 0
	CotPQ_Nhom       = 1
	CotPQ_MoTa       = 2
	CotPQ_StartRole  = 3 // Cột D (Index 3) bắt đầu là các Vai trò
)

// =============================================================
// 2. KHO LƯU TRỮ QUYỀN (Ma trận 2 chiều)
// =============================================================
// Map[TenVaiTro][MaChucNang] -> true/false
var (
	_Map_Quyen map[string]map[string]bool
	mtxQuyen   sync.RWMutex
)

// =============================================================
// 3. LOGIC NẠP DỮ LIỆU "ĐỘNG" TỪ SHEET
// =============================================================
func NapPhanQuyen(targetSpreadsheetID string) {
	if targetSpreadsheetID == "" {
		targetSpreadsheetID = cau_hinh.BienCauHinh.IdFileSheet
	}

	raw, err := loadSheetData(targetSpreadsheetID, "PHAN_QUYEN")
	if err != nil { return }

	// Phải có ít nhất Header (Dòng 1) và Dữ liệu
	if len(raw) < DongBatDau_PhanQuyen { return }

	tempMap := make(map[string]map[string]bool)

	// A. QUÉT HEADER (Dòng 0) ĐỂ TÌM CÁC VAI TRÒ
	header := raw[0]
	var danhSachVaiTro []string 

	// Duyệt từ Cột D trở đi
	for i := CotPQ_StartRole; i < len(header); i++ {
		// Chuẩn hóa tên vai trò: "Thu Kho " -> "thu_kho"
		roleName := strings.ToLower(strings.TrimSpace(layString(header, i)))
		roleName = strings.ReplaceAll(roleName, " ", "_") 
		
		if roleName != "" {
			danhSachVaiTro = append(danhSachVaiTro, roleName)
			tempMap[roleName] = make(map[string]bool)
		}
	}

	// B. DUYỆT CÁC DÒNG CHỨC NĂNG (Từ dòng bắt đầu)
	for i, row := range raw {
		if i < DongBatDau_PhanQuyen-1 { continue }

		maChucNang := strings.TrimSpace(layString(row, CotPQ_MaChucNang))
		if maChucNang == "" { continue }

		// Duyệt qua từng cột vai trò tương ứng với Header đã tìm được
		for j, roleName := range danhSachVaiTro {
			colIndex := CotPQ_StartRole + j
			
			// Lấy giá trị (1 hoặc TRUE là được phép)
			val := layString(row, colIndex)
			isAllow := (val == "1" || strings.ToLower(val) == "true")

			if isAllow {
				tempMap[roleName][maChucNang] = true
			}
		}
	}

	mtxQuyen.Lock()
	_Map_Quyen = tempMap
	mtxQuyen.Unlock()
}

// =============================================================
// 4. HÀM KIỂM TRA (CHECK)
// =============================================================
func KiemTraQuyen(vaiTro string, maChucNang string) bool {
	// Super Admin luôn đúng (Hardcode để tránh bị khóa nếu config sai)
	if vaiTro == "admin_root" { return true }

	mtxQuyen.RLock()
	defer mtxQuyen.RUnlock()

	// Chuẩn hóa input
	vaiTro = strings.ToLower(strings.TrimSpace(vaiTro))
	vaiTro = strings.ReplaceAll(vaiTro, " ", "_") 

	if listQuyen, ok := _Map_Quyen[vaiTro]; ok {
		// Kiểm tra quyền cụ thể
		if allowed, exist := listQuyen[maChucNang]; exist {
			return allowed
		}
	}

	// Mặc định CHẶN nếu không tìm thấy
	return false
}
