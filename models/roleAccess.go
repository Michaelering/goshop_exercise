package models

type RoleAccess struct {
	AccessId int `json:"access_id"`
	RoleId   int `json:"role_id"`
}

func (RoleAccess) TableName() string {
	return "role_access"
}
