package core

import (
	"strings"
	"sync"
	"app/cau_hinh"
)

const (
	DongBatDau_PhanQuyen = 11
	
	CotPQ_MaChucNang = 0
	CotPQ_Nhom       = 1
	CotPQ_MoTa       = 2
	CotPQ_StartRole  = 3
)

type VaiTroInfo struct {
	MaVaiTro   string
	TenVaiTro  string
	// [MỚI] THUẬT TOÁN TÁCH SỐ UI
	StyleLevel int // Level phân quyền (Hàng chục)
	StyleTheme int // Mã màu (Hàng đơn vị)
}

var (
	CachePhanQuyen = make(map[string]map[string]map[string]bool)
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)
	mtxQuyen sync.RWMutex
)

func NapPhanQuyen(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }

	raw, err := LoadSheetData(shopID, "PHAN_QUYEN")
	if err != nil { return }

	headerIndex := -1
	styleIndex := -1

	for i, row := range raw {
		if len(row) > 0 {
			firstCell := strings.TrimSpace(strings.ToLower(LayString(row, 0)))
			if firstCell == "ma_chuc_nang" {
				headerIndex = i
			} else if firstCell == "style" { // TÌM ĐÚNG DÒNG CHỨA STYLE
				styleIndex = i
			}
		}
	}

	if headerIndex == -1 { return } 

	tempMap := make(map[string]map[string]bool)
	var danhSachVaiTroCuaShop []VaiTroInfo 

	header := raw[headerIndex]
	var listMaVaiTro []string 

	for i := CotPQ_StartRole; i < len(header); i++ {
		headerText := strings.TrimSpace(LayString(header, i))
		if headerText == "" { continue }

		parts := strings.Split(headerText, "|")
		roleID := strings.ToLower(strings.TrimSpace(parts[0]))
		roleID = strings.ReplaceAll(roleID, " ", "_") 

		roleName := roleID 
		if len(parts) > 1 {
			roleName = strings.TrimSpace(parts[1]) 
		}

		if roleID != "" {
			listMaVaiTro = append(listMaVaiTro, roleID)
			tempMap[roleID] = make(map[string]bool)
			
			// [MỚI] THUẬT TOÁN ĐỌC STYLE MA TRẬN TỪ SHEETS
			styleCode := 5 // Mặc định Level 5 nếu ô rỗng
			if styleIndex != -1 {
				val := LayInt(raw[styleIndex], i)
				if val > 0 { styleCode = val }
			}

			lvl := styleCode
			thm := 0
			// Nếu style là số có 2 chữ số (VD: 42, 15)
			if styleCode >= 10 {
				lvl = styleCode / 10 // Lấy số chục làm Level
				thm = styleCode % 10 // Lấy số lẻ làm Theme màu
			}

			danhSachVaiTroCuaShop = append(danhSachVaiTroCuaShop, VaiTroInfo{
				MaVaiTro:   roleID,
				TenVaiTro:  roleName,
				StyleLevel: lvl,
				StyleTheme: thm,
			})
		}
	}

	for i, row := range raw {
		if i <= headerIndex { continue }

		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" || maChucNang == "style" { continue } // Bỏ qua dòng style khi duyệt data

		for j, roleID := range listMaVaiTro {
			colIndex := CotPQ_StartRole + j
			val := LayString(row, colIndex)
			isAllow := (val == "1" || strings.ToLower(val) == "true")

			if isAllow {
				tempMap[roleID][maChucNang] = true
			}
		}
	}

	mtxQuyen.Lock()
	CachePhanQuyen[shopID] = tempMap
	CacheDanhSachVaiTro[shopID] = danhSachVaiTroCuaShop
	mtxQuyen.Unlock()
}

func KiemTraQuyen(shopID string, vaiTro string, maChucNang string) bool {
	if vaiTro == "quan_tri_he_thong" { return true } 

	mtxQuyen.RLock()
	defer mtxQuyen.RUnlock()

	vaiTro = strings.ToLower(strings.TrimSpace(vaiTro))
	vaiTro = strings.ReplaceAll(vaiTro, " ", "_") 

	if shopMap, ok := CachePhanQuyen[shopID]; ok {
		if listQuyen, exists := shopMap[vaiTro]; exists {
			if allowed, has := listQuyen[maChucNang]; has {
				return allowed
			}
		}
	}
	return false
}

// =========================================================================
// HÀM TIỆN ÍCH: LẤY CẤP BẬC (LEVEL) CỦA MỘT USER ĐỂ KIỂM TRA BẢO MẬT BACKEND
// =========================================================================
func LayCapBacVaiTro(shopID string, maKH string, vaiTro string) int {
	// Ưu tiên cao nhất (GOD): Bot Hệ Thống và Người Sáng Lập mặc định là Tầng 1
	if maKH == "0000000000000000000" || vaiTro == "quan_tri_he_thong" {
		return 1
	}

	mtxQuyen.RLock()
	defer mtxQuyen.RUnlock()
	
	// Quét trong bảng Phân quyền để lấy Level
	for _, v := range CacheDanhSachVaiTro[shopID] {
		if v.MaVaiTro == vaiTro {
			if v.StyleLevel > 0 {
				return v.StyleLevel // Trả về số 1, 2, 3, 4, 5
			}
			return 5 // Nếu Admin quên điền ô style, mặc định ném xuống Tầng 5 thấp nhất
		}
	}
	return 5
}
