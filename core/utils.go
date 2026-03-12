package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"app/config"
)

var ThemVaoHangCho = PushUpdate

// =======================================================
// KHÁCH HÀNG & NHÂN SỰ (Dùng khóa TenSheetKhachHang)
// =======================================================
func LayKhachHang(shopID, userID string) (*KhachHang, bool) {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	kh, ok := CacheMapKhachHang[TaoCompositeKey(shopID, userID)]
	return kh, ok
}

func LayDanhSachKhachHang(shopID string) []*KhachHang {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	return CacheKhachHang[shopID]
}

func TimKhachHangTheoCookie(shopID, cookie string) (*KhachHang, bool) {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	for _, kh := range CacheKhachHang[shopID] {
		if info, ok := kh.RefreshTokens[cookie]; ok {
			// [ĐÃ FIX]: Đổi .ExpiresAt thành .Exp theo đúng chuẩn JSON NoSQL
			if time.Now().Unix() <= info.Exp { return kh, true }
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(shopID, input string) (*KhachHang, bool) {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	input = strings.ToLower(strings.TrimSpace(input))
	for _, kh := range CacheKhachHang[shopID] {
		if strings.ToLower(kh.TenDangNhap) == input || (kh.Email != "" && strings.ToLower(kh.Email) == input) {
			if kh.MaKhachHang == "0000000000000000000" { return nil, false } 
			return kh, true
		}
	}
	return nil, false
}

func TaoMaKhachHangMoi(shopID string) string {
	lock := GetSheetLock(shopID, TenSheetKhachHang)
	lock.RLock(); defer lock.RUnlock()
	for {
		id := LayChuoiSoNgauNhien(19)
		if _, exist := CacheMapKhachHang[TaoCompositeKey(shopID, id)]; !exist { return id }
	}
}

func ThemKhachHangVaoRam(kh *KhachHang) {
	sID := kh.SpreadsheetID
	if sID == "" { sID = config.BienCauHinh.IdFileSheetAdmin } // Fallback chuẩn xác
	lock := GetSheetLock(sID, TenSheetKhachHang)
	lock.Lock(); defer lock.Unlock()
	
	CacheKhachHang[sID] = append(CacheKhachHang[sID], kh)
	CacheMapKhachHang[TaoCompositeKey(sID, kh.MaKhachHang)] = kh
}

// =======================================================
// PHÂN QUYỀN & HỆ THỐNG (Dùng khóa TenSheetPhanQuyen)
// =======================================================
func LayCapBacVaiTro(shopID, userID, role string) int {
	if userID == "0000000000000000001" || role == "quan_tri_he_thong" { return 0 }
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.RLock(); defer lock.RUnlock()
	for _, v := range CacheDanhSachVaiTro[shopID] {
		if v.MaVaiTro == role { return v.StyleLevel }
	}
	return 9
}

func KiemTraQuyen(shopID, role, maChucNang string) bool {
	if role == "quan_tri_he_thong" { return true }
	lock := GetSheetLock(shopID, TenSheetPhanQuyen)
	lock.RLock(); defer lock.RUnlock()
	if shopMap, ok := CachePhanQuyen[shopID]; ok {
		if listQuyen, exists := shopMap[role]; exists {
			if allowed, has := listQuyen[maChucNang]; has && allowed { return true }
		}
	}
	return false
}

// =======================================================
// SẢN PHẨM & MASTER DATA (Khóa phân mảnh theo bảng)
// =======================================================

// Hàm sinh ID động không cần phân biệt ngành
func TaoMaSanPhamMoi(shopID, prefix string) string {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	if prefix == "" { prefix = "SP" }
	max := 0
	for k := range CacheMapSanPham {
		if strings.HasPrefix(k, shopID+"__"+prefix) {
			idOnly := strings.TrimPrefix(k, shopID+"__")
			var num int
			fmt.Sscanf(strings.TrimPrefix(idOnly, prefix), "%d", &num)
			if num > max { max = num }
		}
	}
	return fmt.Sprintf("%s%04d", prefix, max+1)
}

func CapNhatSlotThuCong(shopID, dmMa string, slotMoi int) {
	lock := GetSheetLock(shopID, TenSheetDanhMuc)
	lock.Lock(); defer lock.Unlock()
	for _, dm := range CacheDanhMuc[shopID] {
		if dm.MaDanhMuc == dmMa {
			if slotMoi > dm.Slot { dm.Slot = slotMoi; PushUpdate(shopID, TenSheetDanhMuc, dm.DongTrongSheet, CotDM_Slot, slotMoi) }
			break
		}
	}
}

// Lấy danh sách sản phẩm Generic theo từng ngành
func LayDanhSachSanPham(shopID string, maNganh string) []*ProductJSON {
	KhoaHeThong.RLock()
	defer KhoaHeThong.RUnlock()
	if list, ok := CacheSanPham[shopID][maNganh]; ok {
		return list
	}
	return []*ProductJSON{}
}

func LayDanhSachDanhMuc(shopID string) []*DanhMuc { lock := GetSheetLock(shopID, TenSheetDanhMuc); lock.RLock(); defer lock.RUnlock(); return CacheDanhMuc[shopID] }
func LayDanhSachThuongHieu(shopID string) []*ThuongHieu { lock := GetSheetLock(shopID, TenSheetThuongHieu); lock.RLock(); defer lock.RUnlock(); return CacheThuongHieu[shopID] }
func LayDanhSachBienLoiNhuan(shopID string) []*BienLoiNhuan { lock := GetSheetLock(shopID, TenSheetBienLoiNhuan); lock.RLock(); defer lock.RUnlock(); return CacheBienLoiNhuan[shopID] }
func LayDanhSachNhaCungCap(shopID string) []*NhaCungCap { lock := GetSheetLock(shopID, TenSheetNhaCungCap); lock.RLock(); defer lock.RUnlock(); return CacheNhaCungCap[shopID] }

// TÌM VÀ THAY THẾ KHỐI TIN NHẮN TRONG FILE UTILS.GO

// =======================================================
// TIN NHẮN (NoSQL 2 Cột)
// =======================================================
func LayHopThuNguoiDung(shopID, userID, role string) []*TinNhan {
	lock := GetSheetLock(shopID, TenSheetTinNhanMaster)
	lock.RLock(); defer lock.RUnlock()
	var rs []*TinNhan
	
	for _, tn := range CacheTinNhan[shopID] {
		// 1. Kiểm tra xem người này có lỡ ấn "Xóa tàng hình" tin nhắn này chưa?
		isDeleted := false
		for _, id := range tn.TrangThaiXoa {
			if id == userID { isDeleted = true; break }
		}
		if isDeleted { continue }

		// 2. Xác định xem mình có phải Người Nhận / Người Gửi không
		isRecipient := false
		if tn.NguoiGuiID == userID {
			isRecipient = true
		} else {
			for _, nhanID := range tn.NguoiNhanID {
				if nhanID == "ALL" || nhanID == userID {
					isRecipient = true; break
				}
				if tn.LoaiTinNhan == "ROLE" && strings.Contains(nhanID, role) {
					isRecipient = true; break
				}
			}
		}
		
		// 3. Đóng gói cho lên UI
		if isRecipient {
			// Check xem ID đã chui vào rổ "Đã Seen" chưa
			tn.DaDoc = false
			for _, docID := range tn.NguoiDoc {
				if docID == userID { tn.DaDoc = true; break }
			}
			rs = append(rs, tn)
		}
	}
	return rs
}

func ThemMoiTinNhan(shopID string, tn *TinNhan) {
	lock := GetSheetLock(shopID, TenSheetTinNhanMaster)
	lock.Lock()
	tn.DongTrongSheet = DongBatDau_TinNhan + len(CacheTinNhan[shopID])
	CacheTinNhan[shopID] = append(CacheTinNhan[shopID], tn)
	lock.Unlock()
	
	// Nã 1 phát duy nhất 2 cột xuống Google Sheet
	b, _ := json.Marshal(tn)
	PushAppend(shopID, TenSheetTinNhanMaster, []interface{}{ tn.MaTinNhan, string(b) })
}

func DanhDauDocTinNhan(shopID, userID, msgID string) {
	lock := GetSheetLock(shopID, TenSheetTinNhanMaster)
	lock.Lock(); defer lock.Unlock()
	
	for _, tn := range CacheTinNhan[shopID] {
		if tn.MaTinNhan == msgID {
			daDoc := false
			for _, u := range tn.NguoiDoc { 
				if u == userID { daDoc = true; break } 
			}
			
			// Nếu chưa nằm trong rổ Seen thì nhét vào và update Sheet
			if !daDoc { 
				tn.NguoiDoc = append(tn.NguoiDoc, userID)
				b, _ := json.Marshal(tn)
				PushUpdate(shopID, TenSheetTinNhanMaster, tn.DongTrongSheet, CotTN_DataJSON, string(b)) 
			}
			break
		}
	}
}
func ToJSON(v interface{}) string { b, _ := json.Marshal(v); return string(b) }

func LayIntStr(s string) int {
	if s == "" { return 0 }
	s = strings.TrimSpace(s)
	val, err := strconv.Atoi(s)
	if err != nil { return 0 }
	return val
}

// =======================================================
// TRẠM KIỂM SOÁT BẢO VỆ DỮ LIỆU (GATEKEEPER)
// =======================================================
func EnsureKhachHangLoaded(shopID string) error {
	if shopID == "" { shopID = config.BienCauHinh.IdFileSheetAdmin }
	
	StatusMutex.RLock()
	status := CacheStatusKhachHang[shopID]
	StatusMutex.RUnlock()

	if status == FlagOK {
		return nil
	}
	
	if status == FlagLoading {
		for i := 0; i < 6; i++ {
			time.Sleep(500 * time.Millisecond)
			StatusMutex.RLock()
			s := CacheStatusKhachHang[shopID]
			StatusMutex.RUnlock()
			if s == FlagOK { return nil }
		}
		return fmt.Errorf("Hệ thống đang đồng bộ dữ liệu. Vui lòng thử lại sau giây lát!")
	}

	// [ĐÃ FIX]: Sử dụng Động cơ Định tuyến chuẩn thay vì gọi hàm cũ đã xóa
	if shopID == config.BienCauHinh.IdFileSheetMaster {
		return NapKhachHangMaster(shopID)
	} else if shopID == config.BienCauHinh.IdFileSheetAdmin {
		return NapKhachHangAdmin(shopID)
	} else {
		NapDuLieuCuaMotShop(shopID)
		return nil
	}
}
