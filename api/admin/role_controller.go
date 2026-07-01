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

	role.Title = title
	role.Description = description

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
	models.DB.Delete(&models.Role{Id: id})
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

	// 获取所有权限
	accessList := []models.Access{}
	models.DB.Where("module_id=?", 0).Preload("AccessItem").Find(&accessList)

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
		for j := 0; j < len(accessList[i].AccessItem); j++ {
			if _, ok := roleAccessMap[accessList[i].AccessItem[j].Id]; ok {
				accessList[i].AccessItem[j].Checked = true
			}
		}
	}

	common.Success(c, gin.H{
		"roleId":     roleId,
		"accessList": accessList,
	})
}

func (con RoleController) DoAuth(c *gin.Context) {
	roleId, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
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
