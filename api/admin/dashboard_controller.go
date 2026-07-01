package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct{}

// 超级管理员角色ID
const superAdminRoleId = 1

// SQL 注入防护：允许的表名白名单
var allowedTables = map[string]bool{
	"manager": true, "role": true, "access": true,
	"goods": true, "goods_cate": true, "goods_type": true,
	"goods_type_attribute": true, "goods_image": true,
	"nav": true, "focus": true, "merchant": true, "setting": true,
}

// SQL 注入防护：允许的字段名白名单
var allowedFields = map[string]bool{
	"status": true, "sort": true, "is_hot": true, "is_best": true, "is_new": true,
}

func (con DashboardController) Index(c *gin.Context) {
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
	roleId, _ := c.Get("roleId")
	isBuiltin, _ := c.Get("isBuiltin")

	accessList := []models.Access{}
	models.DB.Where("parent_id=?", 0).Preload("Children", func(db *gorm.DB) *gorm.DB {
		return db.Order("access.sort DESC")
	}).Order("sort DESC").Find(&accessList)

	// 超级管理员判定：
	// 1. roleId == 1（超管角色ID）+ is_builtin == 1 → 直接返回全部菜单
	// 2. roleId == 1 但 is_builtin == 0 → 兼容旧 token，查 DB 确认角色权限
	// 3. roleId == 0 → 旧 token 可能无 role_id，查 DB 获取真实 roleId
	effectiveRoleId := 0
	if rid, ok := roleId.(int); ok && rid > 0 {
		effectiveRoleId = rid
	}

	// 兼容：token 中 roleId 可能为 0（旧 is_super 用户），从 DB 重新查询
	if effectiveRoleId == 0 {
		adminId, _ := c.Get("adminId")
		if aid, ok := adminId.(int); ok && aid > 0 {
			var mgr models.Manager
			models.DB.Where("id = ?", aid).Find(&mgr)
			if mgr.RoleId > 0 {
				effectiveRoleId = mgr.RoleId
			}
		}
	}

	// 查询角色是否为超管
	isSuperAdmin := false
	if effectiveRoleId == superAdminRoleId {
		// roleId=1 就是超管角色，不论 isBuiltin
		isSuperAdmin = true
	} else if builtinVal, ok := isBuiltin.(int); ok && builtinVal == 1 && effectiveRoleId == superAdminRoleId {
		isSuperAdmin = true
	}

	if isSuperAdmin {
		common.Success(c, accessList)
		return
	}

	// 获取角色权限（使用修正后的 effectiveRoleId）
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id=?", effectiveRoleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = v.AccessId
	}

	// 标记有权限的节点
	for i := 0; i < len(accessList); i++ {
		if _, ok := roleAccessMap[accessList[i].Id]; ok {
			accessList[i].Checked = true
		}
		for j := 0; j < len(accessList[i].Children); j++ {
			if _, ok := roleAccessMap[accessList[i].Children[j].Id]; ok {
				accessList[i].Children[j].Checked = true
			}
		}
	}

	// 只返回有权限的模块
	checkedAccessList := []models.Access{}
	for _, v := range accessList {
		if v.Checked {
			// 过滤子菜单：只保留 Type=2 且有权限的
			checkedItems := []models.Access{}
			for _, item := range v.Children {
				if item.Checked && item.Type == 2 {
					checkedItems = append(checkedItems, item)
				}
			}
			v.Children = checkedItems
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

	// SQL 注入防护：白名单校验
	if !allowedTables[table] || !allowedFields[field] {
		common.BadRequest(c, "参数不合法")
		return
	}

	// 使用白名单校验后的参数（安全）
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

	// SQL 注入防护：白名单校验
	if !allowedTables[table] || !allowedFields[field] {
		common.BadRequest(c, "参数不合法")
		return
	}
	// num 已经通过 PostForm 获取且经过了白名单校验的 table/field
	// num 也需要做基本校验（确保是数字）
	numInt, err := models.Int(num)
	if err != nil {
		common.BadRequest(c, "排序值必须是数字")
		return
	}

	err = models.DB.Exec("UPDATE "+table+" SET "+field+"=? WHERE id=?", numInt, id).Error
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
