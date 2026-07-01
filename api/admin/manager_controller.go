package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type ManagerController struct{}

func (con ManagerController) Index(c *gin.Context) {
	page, _ := models.Int(c.DefaultQuery("page", "1"))
	pageSize, _ := models.Int(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	where := "1=1"
	var args []interface{}
	if keyword != "" {
		where += " AND username LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	managerList := []models.Manager{}
	var count int64
	models.DB.Where(where, args...).Preload("Role").Offset((page - 1) * pageSize).Limit(pageSize).Find(&managerList)
	models.DB.Where(where, args...).Table("manager").Count(&count)

	// 隐藏密码
	for i := range managerList {
		managerList[i].Password = ""
	}

	common.List(c, managerList, count, page, pageSize)
}

func (con ManagerController) Create(c *gin.Context) {
	roleId, _ := models.Int(c.PostForm("role_id"))
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	if len(username) < 2 || len(password) < 6 {
		common.BadRequest(c, "用户名或者密码的长度不合法")
		return
	}

	var count int64
	models.DB.Where("username=?", username).Table("manager").Count(&count)
	if count > 0 {
		common.BadRequest(c, "此管理员已存在")
		return
	}

	manager := models.Manager{
		Username: username,
		Password: models.Md5(password),
		Email:    email,
		Mobile:   mobile,
		RoleId:   roleId,
		Status:   1,
		AddTime:  int(models.GetUnix()),
	}
	err := models.DB.Create(&manager).Error
	if err != nil {
		common.Error(c, 500, "增加管理员失败")
		return
	}
	manager.Password = ""
	common.Success(c, manager)
}

func (con ManagerController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	roleId, _ := models.Int(c.PostForm("role_id"))
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	manager := models.Manager{Id: id}
	models.DB.Find(&manager)
	if manager.Id == 0 {
		common.BadRequest(c, "管理员不存在")
		return
	}

	manager.Username = username
	manager.Email = email
	manager.Mobile = mobile
	manager.RoleId = roleId

	if password != "" {
		if len(password) < 6 {
			common.BadRequest(c, "密码长度不能小于6位")
			return
		}
		manager.Password = models.Md5(password)
	}

	err = models.DB.Save(&manager).Error
	if err != nil {
		common.Error(c, 500, "修改失败")
		return
	}
	manager.Password = ""
	common.Success(c, manager)
}

func (con ManagerController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.Manager{Id: id})
	common.Success(c, nil)
}

func (con ManagerController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	manager := models.Manager{Id: id}
	models.DB.Find(&manager)
	if manager.Id == 0 {
		common.BadRequest(c, "管理员不存在")
		return
	}
	// 获取所有角色供下拉选择
	roleList := []models.Role{}
	models.DB.Find(&roleList)

	manager.Password = ""
	common.Success(c, gin.H{
		"manager":  manager,
		"roleList": roleList,
	})
}
