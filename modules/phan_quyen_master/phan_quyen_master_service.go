package phan_quyen_master

import (
	"app/core"
	"errors"
	"html"
	"strings"
	"time"
)

// [ĐÃ CẬP NHẬT]: Thêm cờ isMasterUser để phân quyền người đang thực thi lệnh lưu
func Service_XuLyLuu(masterID string, isNew bool, input *core.PhanQuyen, isMasterUser bool) error {
	
	// BƯỚC 1: LỌC RÁC & CHỐNG INJECTION (XSS)
	input.MaVaiTro = strings.ToUpper(strings.TrimSpace(html.EscapeString(input.MaVaiTro)))
	input.TenVaiTro = strings.TrimSpace(html.EscapeString(input.TenVaiTro))
	input.MoTa = strings.TrimSpace(html.EscapeString(input.MoTa))

	if input.MaVaiTro == "" || input.TenVaiTro == "" { 
		return errors.New("Mã và Tên vai trò là bắt buộc") 
	}

	// BƯỚC 2: QUYỀN LỰC CỦA MASTER (ID 001 MỚI ĐƯỢC KHÓA)
	if !isMasterUser {
		input.IsLocked = false // Không phải Chủ tịch thì không có quyền khóa
	}

	// BƯỚC 3: MÀNG LỌC QUYỀN HẠN (LỌC SẠCH DATA FAKE TỪ TRÌNH DUYỆT)
	var quyenHanSanhSach []string
	for _, q := range input.QuyenHan {
		q = strings.TrimSpace(q)
		// Chỉ những quyền có trong Bảng DanhSachQuyenHanChuan ở models.go mới được phép giữ lại
		if core.DanhSachQuyenHanChuan[q] { 
			quyenHanSanhSach = append(quyenHanSanhSach, q)
		}
	}
	input.QuyenHan = quyenHanSanhSach

	// BƯỚC 4: ĐẶC QUYỀN CỦA QUẢN TRỊ HỆ THỐNG
	if input.MaVaiTro == "QUAN_TRI_HE_THONG" {
		if !isMasterUser { 
			return errors.New("Chỉ Sáng Lập Viên mới được phép thao tác quyền QUAN_TRI_HE_THONG") 
		}
		input.Level = 0
		input.TrangThai = 1
		
		// Tự động nhồi Full 100% quyền hạn chuẩn từ Master List
		input.QuyenHan = []string{}
		for q := range core.DanhSachQuyenHanChuan {
			input.QuyenHan = append(input.QuyenHan, q)
		}
	} else {
		if input.Level < 1 || input.Level > 9 { 
			return errors.New("Cấp bậc (Level) phải từ 1 đến 9") 
		}
	}

	// BƯỚC 5: XỬ LÝ LƯU HOẶC CẬP NHẬT
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

		// CHỐT CHẶN TỬ THẦN: Vai trò đã bị Master khóa thì cấm người khác sửa
		if old.IsLocked && !isMasterUser {
			return errors.New("Vai trò này đã bị Sáng Lập Viên khóa cứng. Bạn không có quyền chỉnh sửa!")
		}

		lock := core.GetSheetLock(masterID, core.TenSheetCauHinhMaster)
		lock.Lock()
		
		input.SpreadsheetID = old.SpreadsheetID
		input.DongTrongSheet = old.DongTrongSheet
		input.CreatedAt = old.CreatedAt
		input.Version = old.Version + 1
		input.UpdatedAt = time.Now().Unix()
		
		// Giữ nguyên cờ khóa nếu người đang sửa không phải là Master
		if !isMasterUser {
			input.IsLocked = old.IsLocked
		}
		
		*old = *input
		lock.Unlock()

		Repo_Update(masterID, old)
	}
	return nil
}
