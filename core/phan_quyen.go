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

// Struct mới chứa trọn bộ thông tin Chức vụ
type VaiTroInfo struct {
	MaVaiTro  string
	TenVaiTro string
}

var (
	// CachePhanQuyen: Map[ShopID] -> Map[MaVaiTro] -> Map[MaChucNang] -> true/false
	CachePhanQuyen = make(map[string]map[string]map[string]bool)
	
	// [MỚI] Lưu danh sách chức vụ để đổ ra Giao diện Dropdown
	// Map[ShopID] -> []VaiTroInfo
	CacheDanhSachVaiTro = make(map[string][]VaiTroInfo)
	
	mtxQuyen sync.RWMutex
)

func NapPhanQuyen(shopID string) {
	if shopID == "" { shopID = cau_hinh.BienCauHinh.IdFileSheet }

	raw, err := LoadSheetData(shopID, "PHAN_QUYEN")
	if err != nil { return }
	if len(raw) < DongBatDau_PhanQuyen { return }

	tempMap := make(map[string]map[string]bool)
	var danhSachVaiTroCuaShop []VaiTroInfo // Mảng chứa list chức vụ

	// A. QUÉT HEADER (Dòng số 0 trong mảng raw)
	header := raw[0]
	var listMaVaiTro []string 

	for i := CotPQ_StartRole; i < len(header); i++ {
		headerText := strings.TrimSpace(LayString(header, i))
		if headerText == "" { continue }

		// Tách chuỗi theo dấu "|" (Ví dụ: "giam_doc | Giám đốc")
		parts := strings.Split(headerText, "|")
		
		roleID := strings.ToLower(strings.TrimSpace(parts[0]))
		roleID = strings.ReplaceAll(roleID, " ", "_") 

		roleName := roleID // Mặc định tên = mã
		if len(parts) > 1 {
			roleName = strings.TrimSpace(parts[1]) // Nếu có dấu | thì lấy vế sau làm tên
		}

		if roleID != "" {
			listMaVaiTro = append(listMaVaiTro, roleID)
			tempMap[roleID] = make(map[string]bool)
			
			// Thêm vào mảng giao diện
			danhSachVaiTroCuaShop = append(danhSachVaiTroCuaShop, VaiTroInfo{
				MaVaiTro:  roleID,
				TenVaiTro: roleName,
			})
		}
	}

	// B. DUYỆT DỮ LIỆU CÁC DÒNG QUYỀN (Từ dòng 11)
	for i, row := range raw {
		if i < DongBatDau_PhanQuyen-1 { continue }

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

	// LƯU VÀO RAM
	mtxQuyen.Lock()
	CachePhanQuyen[shopID] = tempMap
	CacheDanhSachVaiTro[shopID] = danhSachVaiTroCuaShop
	mtxQuyen.Unlock()
}

func KiemTraQuyen(shopID string, vaiTro string, maChucNang string) bool {
	if vaiTro == "quan_tri_vien_he_thong" { return true } // Quyền tối cao

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
