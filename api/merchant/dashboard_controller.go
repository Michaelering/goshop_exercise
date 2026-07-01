package merchant

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func (con DashboardController) Index(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	shopName, _ := c.Get("shopName")
	username, _ := c.Get("username")

	var goodsCount int64
	models.DB.Model(&models.Goods{}).Where("merchant_id=? AND is_delete=0", merchantId).Count(&goodsCount)

	var orderCount int64
	models.DB.Raw("SELECT COUNT(DISTINCT o.id) FROM `order` o "+
		"INNER JOIN order_item oi ON oi.order_id = o.id "+
		"INNER JOIN goods g ON g.id = oi.product_id "+
		"WHERE g.merchant_id=?", merchantId).Count(&orderCount)

	common.Success(c, gin.H{
		"username":    username,
		"shopName":    shopName,
		"goodsCount":  goodsCount,
		"orderCount":  orderCount,
	})
}
