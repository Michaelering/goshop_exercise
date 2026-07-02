package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type OrderController struct{}

// Index 列出所有订单（管理员视角，不限制商户）
func (con OrderController) Index(c *gin.Context) {
	page, _ := models.Int(c.DefaultQuery("page", "1"))
	pageSize, _ := models.Int(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	where := "1=1"
	var args []interface{}
	if keyword != "" {
		where += " AND (order_id LIKE ? OR name LIKE ? OR phone LIKE ?)"
		kw := "%" + keyword + "%"
		args = append(args, kw, kw, kw)
	}

	var totalCount int64
	models.DB.Table("`order`").Where(where, args...).Count(&totalCount)

	orderList := []models.Order{}
	models.DB.Where(where, args...).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").
		Find(&orderList)

	// 预加载订单明细
	for i := 0; i < len(orderList); i++ {
		models.DB.Where("order_id=?", orderList[i].Id).Find(&orderList[i].OrderItem)
	}

	common.List(c, orderList, totalCount, page, pageSize)
}

// Get 获取单个订单详情
func (con OrderController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	order := models.Order{Id: id}
	models.DB.Find(&order)
	if order.Id == 0 {
		common.BadRequest(c, "订单不存在")
		return
	}

	models.DB.Where("order_id=?", order.Id).Find(&order.OrderItem)

	common.Success(c, order)
}

// Update 更新订单状态、物流信息等
func (con OrderController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	order := models.Order{Id: id}
	models.DB.Find(&order)
	if order.Id == 0 {
		common.BadRequest(c, "订单不存在")
		return
	}

	// 更新订单状态
	orderStatus := c.PostForm("order_status")
	if orderStatus != "" {
		status, _ := models.Int(orderStatus)
		order.OrderStatus = status
	}
	payStatus := c.PostForm("pay_status")
	if payStatus != "" {
		status, _ := models.Int(payStatus)
		order.PayStatus = status
	}
	payType := c.PostForm("pay_type")
	if payType != "" {
		pt, _ := models.Int(payType)
		order.PayType = pt
	}

	// 物流信息
	logisticsCompany := c.PostForm("logistics_company")
	if logisticsCompany != "" {
		lc, _ := models.Int(logisticsCompany)
		order.LogisticsCompany = lc
	}
	waybillNo := c.PostForm("waybill_no")
	if waybillNo != "" {
		wn, _ := models.Int(waybillNo)
		order.WaybillNo = wn
	}

	// 地址信息
	name := c.PostForm("name")
	if name != "" {
		order.Name = name
	}
	phone := c.PostForm("phone")
	if phone != "" {
		order.Phone = phone
	}
	address := c.PostForm("address")
	if address != "" {
		order.Address = address
	}

	if err := models.DB.Save(&order).Error; err != nil {
		common.Error(c, 500, "更新订单失败")
		return
	}

	common.Success(c, order)
}
