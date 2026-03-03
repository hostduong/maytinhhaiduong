package core

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// ID của file Sheet mẫu trên Drive của 99K
const TemplateSheetID = "1dBYqJ-O9pI_WD4iS4ymEtK6GoCbckT2sPQ6yDEkC60w" 

// HamCloneVaCapQuyenSheet nhân bản file mẫu và share quyền cho email khách
func HamCloneVaCapQuyenSheet(emailKhach string, tenShop string, authJson string) (string, error) {
	ctx := context.Background()
	
	// Dùng Service Account (authJson) để xác thực
	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(authJson)))
	if err != nil {
		return "", fmt.Errorf("lỗi xác thực Drive API: %v", err)
	}

	// 1. Nhân bản file
	tenFileMoi := fmt.Sprintf("[99K.VN] Database - %s", tenShop)
	fileCopy := &drive.File{
		Name: tenFileMoi,
	}
	
	newFile, err := driveService.Files.Copy(TemplateSheetID, fileCopy).Do()
	if err != nil {
		return "", fmt.Errorf("lỗi khi clone file mẫu: %v", err)
	}

	// 2. Cấp quyền Editor cho Email của khách
	perm := &drive.Permission{
		Type:         "user",
		Role:         "writer", // Quyền chỉnh sửa
		EmailAddress: emailKhach,
	}
	_, err = driveService.Permissions.Create(newFile.Id, perm).SendNotificationEmail(false).Do()
	if err != nil {
		return newFile.Id, fmt.Errorf("clone thành công nhưng lỗi cấp quyền cho %s: %v", emailKhach, err)
	}

	return newFile.Id, nil
}
