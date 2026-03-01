package cau_hinh_he_thong

import (
	"fmt"
	"strconv"
	"strings"
	"app/core"
)

type CauHinhRepo struct{}

// 1. Lấy danh sách NCC an toàn
func (r *CauHinhRepo) LayDanhSachNCC(shopID string) []*core.NhaCungCap {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.RLock()
	defer lock.RUnlock()
	return core.CacheNhaCungCap[shopID]
}

// 2. Tìm chi tiết 1 NCC
func (r *CauHinhRepo) FindNCCByCode(shopID, maNCC string) (*core.NhaCungCap, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.RLock()
	defer lock.RUnlock()
	// Giả định bạn đã khai báo CacheMapNhaCungCap trong core/ram_cache.go
	ncc, ok := core.CacheMapNhaCungCap[core.TaoCompositeKey(shopID, maNCC)]
	return ncc, ok
}

// 3. Tự động sinh mã NCC (Tìm số lớn nhất)
func (r *CauHinhRepo) TaoMaNCCMoi(shopID string) string {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.RLock()
	defer lock.RUnlock()
	
	prefix := "NCC"
	maxNum := 0
	for _, ncc := range core.CacheNhaCungCap[shopID] {
		if strings.HasPrefix(ncc.MaNhaCungCap, prefix) {
			numStr := strings.TrimPrefix(ncc.MaNhaCungCap, prefix)
			if num, err := strconv.Atoi(numStr); err == nil && num > maxNum {
				maxNum = num
			}
		}
	}
	return fmt.Sprintf("%s%03d", prefix, maxNum+1)
}

// 4. Lưu NCC Mới (Cập nhật RAM + Đẩy lệnh APPEND)
func (r *CauHinhRepo) InsertNCC(shopID string, ncc *core.NhaCungCap) {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.Lock()
	
	// Cập nhật RAM
	core.CacheNhaCungCap[shopID] = append(core.CacheNhaCungCap[shopID], ncc)
	core.CacheMapNhaCungCap[core.TaoCompositeKey(shopID, ncc.MaNhaCungCap)] = ncc
	
	// Xác định dòng để báo cho RAM (phục vụ update sau này)
	ncc.DongTrongSheet = core.DongBatDau_NhaCungCap + len(core.CacheNhaCungCap[shopID]) - 1
	lock.Unlock()

	// Đẩy xuống Queue APPEND 23 cột
	rowData := []interface{}{
		ncc.MaNhaCungCap, ncc.TenNhaCungCap, ncc.MaSoThue, ncc.DienThoai, ncc.Email,
		ncc.KhuVuc, ncc.DiaChi, ncc.NguoiLienHe, ncc.NganHang, ncc.NhomNhaCungCap,
		ncc.LoaiNhaCungCap, ncc.DieuKhoanThanhToan, ncc.ChietKhauMacDinh, ncc.HanMucCongNo,
		ncc.CongNoDauKy, ncc.TongMua, ncc.NoCanTra, ncc.ThongTinThemJson, ncc.TrangThai,
		ncc.GhiChu, ncc.NguoiTao, ncc.NgayTao, ncc.NgayCapNhat,
	}
	core.PushAppend(shopID, core.TenSheetNhaCungCap, rowData)
}

// 5. Cập nhật NCC cũ (Cập nhật RAM + Đẩy lệnh UPDATE)
func (r *CauHinhRepo) UpdateNCC(shopID string, ncc *core.NhaCungCap) {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.Lock()
	// Dữ liệu trong ncc (con trỏ) đã được Service tính toán và gán, ta chỉ cần gọi Update Queue
	row := ncc.DongTrongSheet
	lock.Unlock()

	// Ghi đè từng ô cần thiết (Batch Update)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_TenNhaCungCap, ncc.TenNhaCungCap)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_DienThoai, ncc.DienThoai)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_CongNoDauKy, ncc.CongNoDauKy)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_NoCanTra, ncc.NoCanTra)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_ThongTinThemJson, ncc.ThongTinThemJson)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_NgayCapNhat, ncc.NgayCapNhat)
	// (Bạn bổ sung các lệnh PushUpdate cho các cột khác tương tự)
}
