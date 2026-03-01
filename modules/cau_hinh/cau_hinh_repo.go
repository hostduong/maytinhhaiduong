package cau_hinh

import (
	"fmt"
	"strconv"
	"strings"
	"app/core"
)

type Repo struct{}

func (r *Repo) TaoMaNCCMoi(shopID string) string {
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

func (r *Repo) FindNCCByCode(shopID, maNCC string) (*core.NhaCungCap, bool) {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.RLock()
	defer lock.RUnlock()
	ncc, ok := core.CacheMapNhaCungCap[core.TaoCompositeKey(shopID, maNCC)]
	return ncc, ok
}

func (r *Repo) InsertNCC(shopID string, ncc *core.NhaCungCap) {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.Lock()
	core.CacheNhaCungCap[shopID] = append(core.CacheNhaCungCap[shopID], ncc)
	core.CacheMapNhaCungCap[core.TaoCompositeKey(shopID, ncc.MaNhaCungCap)] = ncc
	ncc.DongTrongSheet = core.DongBatDau_NhaCungCap + len(core.CacheNhaCungCap[shopID]) - 1
	lock.Unlock()

	rowData := []interface{}{
		ncc.MaNhaCungCap, ncc.TenNhaCungCap, ncc.MaSoThue, ncc.DienThoai, ncc.Email,
		ncc.KhuVuc, ncc.DiaChi, ncc.NguoiLienHe, ncc.NganHang, ncc.NhomNhaCungCap,
		ncc.LoaiNhaCungCap, ncc.DieuKhoanThanhToan, ncc.ChietKhauMacDinh, ncc.HanMucCongNo,
		ncc.CongNoDauKy, ncc.TongMua, ncc.NoCanTra, ncc.ThongTinThemJson, ncc.TrangThai,
		ncc.GhiChu, ncc.NguoiTao, ncc.NgayTao, ncc.NgayCapNhat,
	}
	core.PushAppend(shopID, core.TenSheetNhaCungCap, rowData)
}

func (r *Repo) UpdateNCC(shopID string, ncc *core.NhaCungCap) {
	lock := core.GetSheetLock(shopID, core.TenSheetNhaCungCap)
	lock.RLock()
	row := ncc.DongTrongSheet
	lock.RUnlock()

	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_TenNhaCungCap, ncc.TenNhaCungCap)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_DienThoai, ncc.DienThoai)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_CongNoDauKy, ncc.CongNoDauKy)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_NoCanTra, ncc.NoCanTra)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_ThongTinThemJson, ncc.ThongTinThemJson)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_TrangThai, ncc.TrangThai)
	core.PushUpdate(shopID, core.TenSheetNhaCungCap, row, core.CotNCC_NgayCapNhat, ncc.NgayCapNhat)
}
