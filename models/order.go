package models

type Order struct {
	Id               int         `json:"id"`
	OrderId          string      `json:"order_id"`
	Uid              int         `json:"uid"`
	AllPrice         float64     `json:"all_price"`
	Phone            string      `json:"phone"`
	Name             string      `json:"name"`
	Address          string      `json:"address"`
	PayStatus        int         `json:"pay_status"`   // 支付状态： 0 表示未支付     1 已经支付
	PayType          int         `json:"pay_type"`     // 支付类型： 0 alipay    1 wechat
	OrderStatus      int         `json:"order_status"` // 订单状态： 0 已下单  1 已付款  2 已配货  3、发货   4、交易成功   5、退货   6、取消
	AddTime          int         `json:"add_time"`
	PayTime          int         `json:"pay_time"`
	DistributionTime int         `json:"distribution_time"`
	ExwarehouseTime  int         `json:"exwarehouse_time"`
	SuccessfulTime   int         `json:"successful_time"`
	CancelTime       int         `json:"cancel_time"`
	ReturnTime       int         `json:"return_time"`
	LogisticsCompany int         `json:"logistics_company"`
	WaybillNo        int         `json:"waybill_no"`
	OrderItem        []OrderItem `gorm:"foreignKey:OrderId;references:Id" json:"order_item"`
}

func (Order) TableName() string {
	return "order"
}
