package goi_dich_vu_master

import (
	"app/core"
	"errors"
	"time"
)

func Service_XuLyLuu(shopID string, isNew bool, input *core.GoiDichVu) error {
	if input.MaGoi == "" || input.TenGoi == "" { return errors.New("Mã và Tên gói là bắt buộc") }

	if isNew {
		if _, exist := Repo_FindByCode(shopID, input.MaGoi); exist { return errors.New("Mã gói đã tồn tại trên hệ thống") }
		
		input.Version = 1
		input.CreatedAt = time.Now().Unix()
		input.UpdatedAt = input.CreatedAt
		if input.NgayBatDau == 0 { input.NgayBatDau = input.CreatedAt }
		
		Repo_Insert(shopID, input)
	} else {
		old, ok := Repo_FindByCode(shopID, input.MaGoi)
		if !ok { return errors.New("Không tìm thấy gói cước để cập nhật") }

		lock := core.GetSheetLock(shopID, core.TenSheetCauHinhMaster)
		lock.Lock()
		
		// Kế thừa dữ liệu hệ thống
		input.SpreadsheetID = old.SpreadsheetID
		input.DongTrongSheet = old.DongTrongSheet
		input.CreatedAt = old.CreatedAt
		input.Version = old.Version + 1
		input.UpdatedAt = time.Now().Unix()
		
		lock.Unlock()

		Repo_Update(shopID, input)
	}
	return nil
}
