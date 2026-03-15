package phan_quyen_master

import (
	"app/core"
	"errors"
	"time"
)

func Service_XuLyLuu(masterID string, isNew bool, input *core.PhanQuyen) error {
	if input.MaVaiTro == "" || input.TenVaiTro == "" { return errors.New("Mã và Tên vai trò là bắt buộc") }

	// [LUẬT THÉP 2]: Đóng đinh các chỉ số cho QUAN_TRI_HE_THONG
	if input.MaVaiTro == "QUAN_TRI_HE_THONG" {
		input.Level = 0
		input.TrangThai = 1
	} else {
		if input.Level < 1 || input.Level > 9 { return errors.New("Cấp bậc (Level) phải từ 1 đến 9") }
	}

	if isNew {
		if _, exist := Repo_FindByCode(masterID, input.MaVaiTro); exist { 
			return errors.New("Mã vai trò đã tồn tại trên hệ thống") 
		}
		
		input.Version = 1
		input.CreatedAt = time.Now().Unix()
		input.UpdatedAt = input.CreatedAt
		
		Repo_Insert(masterID, input)
	} else {
		old, ok := Repo_FindByCode(masterID, input.MaVaiTro)
		if !ok { return errors.New("Không tìm thấy vai trò để cập nhật") }

		lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
		lock.Lock()
		
		input.SpreadsheetID = old.SpreadsheetID
		input.DongTrongSheet = old.DongTrongSheet
		input.CreatedAt = old.CreatedAt
		input.Version = old.Version + 1
		input.UpdatedAt = time.Now().Unix()
		
		*old = *input
		lock.Unlock()

		Repo_Update(masterID, old)
	}
	return nil
}
