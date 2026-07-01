package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessController struct{}

func (con AccessController) Index(c *gin.Context) {
	accessList := []models.Access{}
	models.DB.Where("parent_id=?", 0).Preload("Children").Find(&accessList)
	common.Success(c, accessList)
}

func (con AccessController) TopModules(c *gin.Context) {
	accessList := []models.Access{}
	models.DB.Where("parent_id=?", 0).Find(&accessList)
	common.Success(c, accessList)
}

func (con AccessController) Create(c *gin.Context) {
	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := c.PostForm("action_name")
	accessType, _ := models.Int(c.PostForm("type"))
	urlPrefix := c.PostForm("url_prefix")
	httpMethods := c.PostForm("http_methods")
	parentId, _ := models.Int(c.PostForm("parent_id"))
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))
	description := c.PostForm("description")

	if moduleName == "" {
		common.BadRequest(c, "模块名称不能为空")
		return
	}

	access := models.Access{
		ModuleName:  moduleName,
		Type:        accessType,
		ActionName:  actionName,
		UrlPrefix:   urlPrefix,
		HttpMethods: httpMethods,
		ParentId:    parentId,
		Sort:        sort,
		Description: description,
		Status:      status,
		AddTime:     int(models.GetUnix()),
	}
	err := models.DB.Create(&access).Error
	if err != nil {
		common.Error(c, 500, "增加数据失败")
		return
	}
	common.Success(c, access)
}

func (con AccessController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	moduleName := strings.Trim(c.PostForm("module_name"), " ")
	actionName := c.PostForm("action_name")
	accessType, _ := models.Int(c.PostForm("type"))
	urlPrefix := c.PostForm("url_prefix")
	httpMethods := c.PostForm("http_methods")
	parentId, _ := models.Int(c.PostForm("parent_id"))
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))
	description := c.PostForm("description")

	if moduleName == "" {
		common.BadRequest(c, "模块名称不能为空")
		return
	}

	access := models.Access{Id: id}
	models.DB.Find(&access)
	if access.Id == 0 {
		common.BadRequest(c, "数据不存在")
		return
	}

	access.ModuleName = moduleName
	access.Type = accessType
	access.ActionName = actionName
	access.UrlPrefix = urlPrefix
	access.HttpMethods = httpMethods
	access.ParentId = parentId
	access.Sort = sort
	access.Description = description
	access.Status = status

	err = models.DB.Save(&access).Error
	if err != nil {
		common.Error(c, 500, "修改数据失败")
		return
	}
	common.Success(c, access)
}

func (con AccessController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	access := models.Access{Id: id}
	models.DB.Find(&access)
	if access.ParentId == 0 {
		// 顶级模块，检查是否有子节点
		var childCount int64
		models.DB.Where("parent_id=?", id).Table("access").Count(&childCount)
		if childCount > 0 {
			common.BadRequest(c, "当前模块下面有菜单或者操作，请先删除子节点")
			return
		}
	}
	models.DB.Delete(&access)
	common.Success(c, nil)
}

func (con AccessController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	access := models.Access{Id: id}
	models.DB.Find(&access)
	if access.Id == 0 {
		common.BadRequest(c, "数据不存在")
		return
	}

	// 获取顶级模块供下拉选择
	accessList := []models.Access{}
	models.DB.Where("parent_id=?", 0).Find(&accessList)

	common.Success(c, gin.H{
		"access":     access,
		"accessList": accessList,
	})
}
