package models

type Merchant struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	ShopName string `json:"shop_name"` //店铺名称
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Status   int    `json:"status"` // 1=启用 0=禁用
	AddTime  int    `json:"add_time"`
}

func (Merchant) TableName() string {
	return "merchant"
}
