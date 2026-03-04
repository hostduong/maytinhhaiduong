package khach_hang

import (
	"github.com/gin-gonic/gin"
)

var service = Service{}

func API_MuaGoiKhachHang(c *gin.Context) {
	masterShopID := c.GetString("SHOP_ID")
	userID := c.GetString("USER_ID") // Trạm CheckAuth đảm bảo dòng này luôn có dữ liệu

	maGoi := c.PostForm("ma_goi")
	maCode := c.PostForm("ma_code")

	// Đẩy vào đường ống Service xử lý
	redirectURL, err := service.XuLyMuaGoiStarter(masterShopID, userID, maGoi, maCode)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	// [NƠI GỌI API GOOGLE CLOUD RUN TẠO SUBDOMAIN SẼ NẰM Ở ĐÂY LÚC TRIỂN KHAI THỰC TẾ]
	// errCloud := CallCloudRunAPI(tenDangNhap)
    // ...

	c.JSON(200, gin.H{
		"status":       "ok",
		"msg":          "Khởi tạo gói cước thành công!",
		"redirect_url": redirectURL,
	})
}
