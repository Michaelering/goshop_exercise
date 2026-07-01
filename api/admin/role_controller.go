package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

func (con RoleController) Index(c *gin.Context) {
	roleList := []models.Role{}
	models.DB.Find(&roleList)
	common.Success(c, roleList)
}

func (con RoleController) Create(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")

	if title == "" {
		common.BadRequest(c, "角色的标题不能为空")
		return
	}

	role := models.Role{
		Title:       title,
		Description: description,
		Status:      1,
		IsBuiltin:   0, // 用户创建的角色不是内置角色
		AddTime:     int(models.GetUnix()),
	}
	err := models.DB.Create(&role).Error
	if err != nil {
		common.Error(c, 500, "增加角色失败")
		return
	}
	common.Success(c, role)
}

func (con RoleController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	title := strings.Trim(c.PostForm("title"), " ")
	description := strings.Trim(c.PostForm("description"), " ")

	if title == "" {
		common.BadRequest(c, "角色的标题不能为空")
		return
	}

	role := models.Role{Id: id}
	models.DB.Find(&role)
	if role.Id == 0 {
		common.BadRequest(c, "角色不存在")
		return
	}

	// 内置角色不允许修改标题（仅允许修改描述）
	if role.IsBuiltin == 1 {
		role.Description = description
	} else {
		role.Title = title
		role.Description = description
	}

	err = models.DB.Save(&role).Error
	if err != nil {
		common.Error(c, 500, "修改角色失败")
		return
	}
	common.Success(c, role)
}

func (con RoleController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	// 检查是否为内置角色，内置角色不可删除
	role := models.Role{Id: id}
	models.DB.Find(&role)
	if role.Id == 0 {
		common.BadRequest(c, "角色不存在")
		return
	}
	if role.IsBuiltin == 1 {
		common.BadRequest(c, "内置角色不可删除")
		return
	}

	// 检查是否有管理员正在使用此角色
	var managerCount int64
	models.DB.Model(&models.Manager{}).Where("role_id = ?", id).Count(&managerCount)
	if managerCount > 0 {
		common.BadRequest(c, "该角色下还有管理员，请先移除或转移管理员")
		return
	}

	models.DB.Delete(&role)
	common.Success(c, nil)
}

func (con RoleController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	role := models.Role{Id: id}
	models.DB.Find(&role)
	if role.Id == 0 {
		common.BadRequest(c, "角色不存在")
		return
	}
	common.Success(c, role)
}

func (con RoleController) Auth(c *gin.Context) {
	roleId, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	// 获取所有权限（改用 parent_id 和 children）
	accessList := []models.Access{}
	models.DB.Where("parent_id=?", 0).Preload("Children").Find(&accessList)

	// 获取角色已有的权限
	roleAccess := []models.RoleAccess{}
	models.DB.Where("role_id=?", roleId).Find(&roleAccess)
	roleAccessMap := make(map[int]int)
	for _, v := range roleAccess {
		roleAccessMap[v.AccessId] = v.AccessId
	}

	// 标记已有的权限
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

	// 获取角色名
	role := models.Role{Id: roleId}
	models.DB.Find(&role)

	common.Success(c, gin.H{
		"roleId":     roleId,
		"roleTitle":  role.Title,
		"accessList": accessList,
	})
}

func (con RoleController) DoAuth(c *gin.Context) {
	roleId, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	// 检查角色是否存在
	role := models.Role{Id: roleId}
	models.DB.Find(&role)
	if role.Id == 0 {
		common.BadRequest(c, "角色不存在")
		return
	}

	accessIds := c.PostFormArray("access_node[]")

	// 删除旧权限
	models.DB.Where("role_id=?", roleId).Delete(&models.RoleAccess{})

	// 添加新权限
	for _, v := range accessIds {
		accessId, _ := models.Int(v)
		roleAccess := models.RoleAccess{
			RoleId:   roleId,
			AccessId: accessId,
		}
		models.DB.Create(&roleAccess)
	}

	common.Success(c, nil)
}
