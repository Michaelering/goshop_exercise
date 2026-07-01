package models

type Manager struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	RoleId   int    `json:"role_id"`
	AddTime  int    `json:"add_time"`
	Role     Role   `gorm:"foreignKey:RoleId;references:Id" json:"role"`
}

func (Manager) TableName() string {
	return "manager"
}
