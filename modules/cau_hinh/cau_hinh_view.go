package cau_hinh

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func TrangCauHinhView(c *gin.Context) {
	// Tạm thời trả về trang rỗng để đảm bảo hệ thống Build Pass
	// Giao diện sẽ được nạp sau khi hoàn tất quy hoạch thư mục HTML
	c.String(http.StatusOK, "Trang cấu hình hệ thống đang được nâng cấp...")
}
