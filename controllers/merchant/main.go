package merchant

import (
	"encoding/json"
	"ginshop58/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MainController struct {
	BaseController
}

func (con MainController) Index(c *gin.Context) {
	session := sessions.Default(c)
	merchantInfo := session.Get("merchantInfo")
	merchantInfoStr, ok := merchantInfo.(string)

	if ok {
		var merchantInfoStruct []models.Merchant
		json.Unmarshal([]byte(merchantInfoStr), &merchantInfoStruct)
		merchant := merchantInfoStruct[0]

		// 统计该商户的商品数量
		var goodsCount int64
		models.DB.Model(&models.Goods{}).Where("merchant_id=? AND is_delete=0", merchant.Id).Count(&goodsCount)

		// 统计该商户的订单数量（通过商品关联）
		var orderCount int64
		models.DB.Raw("SELECT COUNT(DISTINCT o.id) FROM `order` o "+
			"INNER JOIN order_item oi ON oi.order_id = o.id "+
			"INNER JOIN goods g ON g.id = oi.product_id "+
			"WHERE g.merchant_id=?", merchant.Id).Count(&orderCount)

		c.HTML(http.StatusOK, "merchant/main/index.html", gin.H{
			"username":    merchant.Username,
			"shopName":    merchant.ShopName,
			"goodsCount":  goodsCount,
			"orderCount":  orderCount,
		})
	} else {
		c.Redirect(302, "/merchant/login")
	}
}

func (con MainController) Welcome(c *gin.Context) {
	session := sessions.Default(c)
	merchantInfo := session.Get("merchantInfo")
	merchantInfoStr, ok := merchantInfo.(string)

	if ok {
		var merchantInfoStruct []models.Merchant
		json.Unmarshal([]byte(merchantInfoStr), &merchantInfoStruct)
		merchant := merchantInfoStruct[0]

		c.HTML(http.StatusOK, "merchant/main/welcome.html", gin.H{
			"username": merchant.Username,
			"shopName": merchant.ShopName,
		})
	} else {
		c.Redirect(302, "/merchant/login")
	}
}
