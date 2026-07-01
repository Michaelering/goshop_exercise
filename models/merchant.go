package models

type Merchant struct {
	Id       int
	Username string
	Password string
	ShopName string //店铺名称
	Mobile   string
	Email    string
	Status   int // 1=启用 0=禁用
	AddTime  int
}

func (Merchant) TableName() string {
	return "merchant"
}
