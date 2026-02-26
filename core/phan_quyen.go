package core

import (
	"strings"
	"sync"
	"app/cau_hinh"
)

const (
	DongBatDau_PhanQuyen = 11
	CotPQ_MaChucNang    = 0
	CotPQ_Nhom          = 1
	CotPQ_MoTa          = 2
	CotPQ_StartRole     = 3
)

type VaiTroInfo struct {
	MaVaiTro   string
	TenVaiTro  string
	StyleLevel int // Tầng quyền lực (0 đến 9)
	StyleTheme int // Mã màu sắc (0 đến 9)
}

var (
	CachePhanQuyen      = make(map[string]map[string]map[string]bool)
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)
	mtxQuyen            sync.RWMutex
)

func NapPhanQuyen(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }
	raw, err := LoadSheetData(shopID, "PHAN_QUYEN")
	if err != nil { return }

	headerIndex, styleIndex := -1, -1
	for i, row := range raw {
		if len(row) > 0 {
			firstCell := strings.TrimSpace(strings.ToLower(LayString(row, 0)))
			if firstCell == "ma_chuc_nang" { headerIndex = i } else if firstCell == "style" { styleIndex = i }
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
		roleID := strings.ReplaceAll(strings.ToLower(strings.TrimSpace(parts[0])), " ", "_") 
		roleName := roleID 
		if len(parts) > 1 { roleName = strings.TrimSpace(parts[1]) }

		if roleID != "" {
			listMaVaiTro = append(listMaVaiTro, roleID)
			tempMap[roleID] = make(map[string]bool)
			
			// MA TRẬN 10x10 TỰ ĐỘNG
			styleCode := 90 // Mặc định Level 9, Màu 0 (Khách hàng)
			if styleIndex != -1 {
				val := LayInt(raw[styleIndex], i)
				if val >= 0 { styleCode = val } // Hỗ trợ cả số 0 (Level 0)
			}
			
			lvl, thm := styleCode, 0
			if styleCode >= 10 { lvl = styleCode / 10; thm = styleCode % 10 }

			danhSachVaiTroCuaShop = append(danhSachVaiTroCuaShop, VaiTroInfo{
				MaVaiTro: roleID, TenVaiTro: roleName, StyleLevel: lvl, StyleTheme: thm,
			})
		}
	}

	for i, row := range raw {
		if i <= headerIndex { continue }
		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" || maChucNang == "style" { continue } 

		for j, roleID := range listMaVaiTro {
			val := LayString(row, CotPQ_StartRole+j)
			if val == "1" || strings.ToLower(val) == "true" { tempMap[roleID][maChucNang] = true }
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
	vaiTro = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(vaiTro)), " ", "_") 
	if shopMap, ok := CachePhanQuyen[shopID]; ok {
		if listQuyen, exists := shopMap[vaiTro]; exists {
			if allowed, has := listQuyen[maChucNang]; has { return allowed }
		}
	}
	return false
}

// HÀM LẤY LEVEL QUYỀN LỰC ĐỂ CHỐT CHẶN BẢO MẬT API
func LayCapBacVaiTro(shopID string, maKH string, vaiTro string) int {
	if maKH == "0000000000000000000" || vaiTro == "quan_tri_he_thong" { return 0 } // Đỉnh chóp
	mtxQuyen.RLock()
	defer mtxQuyen.RUnlock()
	for _, v := range CacheDanhSachVaiTro[shopID] {
		if v.MaVaiTro == vaiTro { return v.StyleLevel }
	}
	return 9 // Mặc định tầng chót
}
