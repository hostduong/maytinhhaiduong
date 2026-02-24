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
	MaVaiTro  string
	TenVaiTro string
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

	// [ĐÃ FIX LỖI TYPE]: Sử dụng hàm LayString để ép kiểu an toàn từ interface{} sang string
	headerIndex := -1
	for i, row := range raw {
		if len(row) > 0 {
			firstCell := LayString(row, 0)
			if strings.TrimSpace(strings.ToLower(firstCell)) == "ma_chuc_nang" {
				headerIndex = i
				break
			}
		}
	}

	if headerIndex == -1 { return } 

	tempMap := make(map[string]map[string]bool)
	var danhSachVaiTroCuaShop []VaiTroInfo 

	// A. QUÉT HEADER
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
			
			danhSachVaiTroCuaShop = append(danhSachVaiTroCuaShop, VaiTroInfo{
				MaVaiTro:  roleID,
				TenVaiTro: roleName,
			})
		}
	}

	// B. DUYỆT DỮ LIỆU
	for i, row := range raw {
		if i <= headerIndex { continue }

		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" { continue }

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
