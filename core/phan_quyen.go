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

// BỘ NHỚ ĐA SHOP
// Map[ShopID] -> Map[TenVaiTro] -> Map[MaChucNang] -> true/false
var (
	CachePhanQuyen = make(map[string]map[string]map[string]bool)
	mtxQuyen       sync.RWMutex
)

func NapPhanQuyen(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }

	raw, err := LoadSheetData(shopID, "PHAN_QUYEN")
	if err != nil { return }
	if len(raw) < DongBatDau_PhanQuyen { return }

	tempMap := make(map[string]map[string]bool)

	// A. QUÉT HEADER (Dòng 0)
	header := raw[0]
	var danhSachVaiTro []string 

	for i := CotPQ_StartRole; i < len(header); i++ {
		roleName := strings.ToLower(strings.TrimSpace(LayString(header, i)))
		roleName = strings.ReplaceAll(roleName, " ", "_") 
		if roleName != "" {
			danhSachVaiTro = append(danhSachVaiTro, roleName)
			tempMap[roleName] = make(map[string]bool)
		}
	}

	// B. DUYỆT DỮ LIỆU
	for i, row := range raw {
		if i < DongBatDau_PhanQuyen-1 { continue }

		maChucNang := strings.TrimSpace(LayString(row, CotPQ_MaChucNang))
		if maChucNang == "" { continue }

		for j, roleName := range danhSachVaiTro {
			colIndex := CotPQ_StartRole + j
			val := LayString(row, colIndex)
			isAllow := (val == "1" || strings.ToLower(val) == "true")

			if isAllow {
				tempMap[roleName][maChucNang] = true
			}
		}
	}

	mtxQuyen.Lock()
	CachePhanQuyen[shopID] = tempMap
	mtxQuyen.Unlock()
}

func KiemTraQuyen(shopID string, vaiTro string, maChucNang string) bool {
	// Super Admin luôn đúng
	if vaiTro == "admin_root" { return true }

	mtxQuyen.RLock()
	defer mtxQuyen.RUnlock()

	vaiTro = strings.ToLower(strings.TrimSpace(vaiTro))
	vaiTro = strings.ReplaceAll(vaiTro, " ", "_") 

	// Lấy bảng quyền của Shop đó
	if shopMap, ok := CachePhanQuyen[shopID]; ok {
		if listQuyen, exists := shopMap[vaiTro]; exists {
			if allowed, has := listQuyen[maChucNang]; has {
				return allowed
			}
		}
	}

	return false
}
