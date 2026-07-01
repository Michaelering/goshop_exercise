package models

type Role struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	IsBuiltin   int    `json:"is_builtin"` // 1=内置角色，不可删除
	AddTime     int    `json:"add_time"`
}

func (Role) TableName() string {
	return "role"
}
