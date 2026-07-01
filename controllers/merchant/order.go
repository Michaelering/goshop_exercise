package merchant

import (
	"ginshop58/models"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	BaseController
}

func (con OrderController) Index(c *gin.Context) {
	merchantId := getMerchantId(c)
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	pageSize := 8

	// 查询包含该商户商品的订单（去重）
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

	// 加载每个订单的订单项
	for i := 0; i < len(orderList); i++ {
		models.DB.Where("order_id=?", orderList[i].Id).Find(&orderList[i].OrderItem)
	}

	if len(orderList) > 0 {
		c.HTML(http.StatusOK, "merchant/order/index.html", gin.H{
			"orderList":  orderList,
			"totalPages": math.Ceil(float64(totalCount) / float64(pageSize)),
			"page":       page,
		})
	} else {
		if page != 1 {
			c.Redirect(302, "/merchant/order")
		} else {
			c.HTML(http.StatusOK, "merchant/order/index.html", gin.H{
				"orderList":  orderList,
				"totalPages": math.Ceil(float64(totalCount) / float64(pageSize)),
				"page":       page,
			})
		}
	}
}
