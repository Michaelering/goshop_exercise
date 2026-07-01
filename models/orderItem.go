package models

type OrderItem struct {
	Id           int     `json:"id"`
	OrderId      int     `json:"order_id"`
	Uid          int     `json:"uid"`
	ProductTitle string  `json:"product_title"`
	ProductId    int     `json:"product_id"`
	ProductImg   string  `json:"product_img"`
	ProductPrice float64 `json:"product_price"`
	ProductNum   int     `json:"product_num"`
	GoodsVersion string  `json:"goods_version"`
	GoodsColor   string  `json:"goods_color"`
	AddTime      int     `json:"add_time"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
