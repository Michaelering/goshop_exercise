package merchant

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

func (con OrderController) Index(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	page, _ := models.Int(c.DefaultQuery("page", "1"))
	pageSize, _ := models.Int(c.DefaultQuery("pageSize", "10"))

	var totalCount int64
	models.DB.Raw("SELECT COUNT(DISTINCT o.id) FROM `order` o "+
		"INNER JOIN order_item oi ON oi.order_id = o.id "+
		"INNER JOIN goods g ON g.id = oi.product_id "+
		"WHERE g.merchant_id=?", merchantId).Count(&totalCount)

	orderList := []models.Order{}
	models.DB.Raw("SELECT DISTINCT o.* FROM `order` o "+
		"INNER JOIN order_item oi ON oi.order_id = o.id "+
		"INNER JOIN goods g ON g.id = oi.product_id "+
		"WHERE g.merchant_id=? ORDER BY o.id DESC LIMIT ? OFFSET ?",
		merchantId, pageSize, (page-1)*pageSize).Scan(&orderList)

	for i := 0; i < len(orderList); i++ {
		models.DB.Where("order_id=?", orderList[i].Id).Find(&orderList[i].OrderItem)
	}

	common.List(c, orderList, totalCount, page, pageSize)
}
