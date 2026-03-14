package goi_dich_vu_master

import (
	"app/core"
	"errors"
	"time"
)

func Service_XuLyLuu(masterID string, isNew bool, input *core.GoiDichVu) error {
	if input.MaGoi == "" || input.TenGoi == "" { return errors.New("Mã và Tên gói là bắt buộc") }

	if isNew {
		if _, exist := Repo_FindByCode(masterID, input.MaGoi); exist { 
			return errors.New("Mã gói đã tồn tại trên hệ thống") 
		}
		
		input.Version = 1
		input.CreatedAt = time.Now().Unix()
		input.UpdatedAt = input.CreatedAt
		if input.NgayBatDau == 0 { input.NgayBatDau = input.CreatedAt }
		
		Repo_Insert(masterID, input)
	} else {
		old, ok := Repo_FindByCode(masterID, input.MaGoi)
		if !ok { return errors.New("Không tìm thấy gói cước để cập nhật") }

		lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
		lock.Lock()
		
		// Kế thừa các chỉ số hệ thống từ bản ghi cũ
		input.SpreadsheetID = old.SpreadsheetID
		input.DongTrongSheet = old.DongTrongSheet
		input.CreatedAt = old.CreatedAt
		input.Version = old.Version + 1
		input.UpdatedAt = time.Now().Unix()
		
		// Ghi đè con trỏ trong RAM Cache
		*old = *input
		lock.Unlock()

		Repo_Update(masterID, old)
	}
	return nil
}
