package models

type Access struct {
	Id          int      `json:"id"`
	ModuleName  string   `json:"module_name"`  // 模块名称
	ActionName  string   `json:"action_name"`  // 操作名称
	Type        int      `json:"type"`         // 节点类型: 1=模块组  2=菜单项
	UrlPrefix   string   `json:"url_prefix"`   // 路由前缀，如 "goods"、"manager"
	HttpMethods string   `json:"http_methods"` // 允许的 HTTP 方法，逗号分隔，如 "GET,POST"
	ParentId    int      `json:"parent_id"`    // 父级ID，0=顶级模块
	Sort        int      `json:"sort"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	AddTime     int      `json:"add_time"`
	Children    []Access `gorm:"foreignKey:ParentId;references:Id" json:"children"`
	Checked     bool     `gorm:"-" json:"checked"` // 虚拟字段，不存数据库
}

func (Access) TableName() string {
	return "access"
}
