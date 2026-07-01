package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct{}

func (con DashboardController) Index(c *gin.Context) {
	// 获取统计数据
	var managerCount int64
	models.DB.Model(&models.Manager{}).Count(&managerCount)

	var goodsCount int64
	models.DB.Model(&models.Goods{}).Where("is_delete=0").Count(&goodsCount)

	var orderCount int64
	models.DB.Model(&models.Order{}).Count(&orderCount)

	adminId, _ := c.Get("adminId")
	username, _ := c.Get("username")

	common.Success(c, gin.H{
		"adminId":      adminId,
		"username":     username,
		"managerCount": managerCount,
		"goodsCount":   goodsCount,
		"orderCount":   orderCount,
	})
}

func (con DashboardController) Menu(c *gin.Context) {
	isSuper, _ := c.Get("isSuper")
	roleId, _ := c.Get("roleId")

	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Preload("AccessItem", func(db *gorm.DB) *gorm.DB {
		return db.Order("access.sort DESC")
	}).Order("sort DESC").Find(&accessList)

	if isSuperVal, ok := isSuper.(int); ok && isSuperVal == 1 {
		common.Success(c, accessList)
		return
	}

	// 获取角色权限
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id=?", roleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = v.AccessId
	}

	// 标记有权限的节点
	for i := 0; i < len(accessList); i++ {
		if _, ok := roleAccessMap[accessList[i].Id]; ok {
			accessList[i].Checked = true
		}
		for j := 0; j < len(accessList[i].AccessItem); j++ {
			if _, ok := roleAccessMap[accessList[i].AccessItem[j].Id]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}

	// 只返回有权限的模块
	checkedAccessList := []models.Access{}
	for _, v := range accessList {
		if v.Checked {
			// 过滤子菜单
			checkedItems := []models.Access{}
			for _, item := range v.AccessItem {
				if item.Checked && item.Type == 2 {
					checkedItems = append(checkedItems, item)
				}
			}
			v.AccessItem = checkedItems
			checkedAccessList = append(checkedAccessList, v)
		}
	}

	common.Success(c, checkedAccessList)
}

func (con DashboardController) ChangeStatus(c *gin.Context) {
	id := c.PostForm("id")
	table := c.PostForm("table")
	field := c.PostForm("field")

	if id == "" || table == "" || field == "" {
		common.BadRequest(c, "参数错误")
		return
	}

	err := models.DB.Exec("UPDATE "+table+" SET "+field+"=ABS("+field+"-1) WHERE id=?", id).Error
	if err != nil {
		common.Error(c, 500, "修改失败")
		return
	}
	common.Success(c, nil)
}

func (con DashboardController) ChangeNum(c *gin.Context) {
	id := c.PostForm("id")
	table := c.PostForm("table")
	field := c.PostForm("field")
	num := c.PostForm("num")

	if id == "" || table == "" || field == "" {
		common.BadRequest(c, "参数错误")
		return
	}

	err := models.DB.Exec("UPDATE "+table+" SET "+field+"="+num+" WHERE id=?", id).Error
	if err != nil {
		common.Error(c, 500, "修改失败")
		return
	}
	common.Success(c, nil)
}

func (con DashboardController) FlushAll(c *gin.Context) {
	models.CacheDb.FlushAll()

	common.Success(c, nil)
}
