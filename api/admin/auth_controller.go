package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (con AuthController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha(34, 100, 2)
	if err != nil {
		common.Error(c, 500, "验证码生成失败")
		return
	}
	c.JSON(200, gin.H{
		"code":         200,
		"message":      "success",
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con AuthController) Login(c *gin.Context) {
	captchaId := c.PostForm("captchaId")
	username := c.PostForm("username")
	password := c.PostForm("password")
	verifyValue := c.PostForm("verifyValue")

	if !models.VerifyCaptcha(captchaId, verifyValue) {
		common.BadRequest(c, "验证码验证失败")
		return
	}

	managerList := []models.Manager{}
	password = models.Md5(password)
	models.DB.Where("username=? AND password=?", username, password).Preload("Role").Find(&managerList)

	if len(managerList) == 0 {
		common.Error(c, 400, "用户名或者密码错误")
		return
	}

	manager := managerList[0]
	if manager.Status == 0 {
		common.Error(c, 400, "该管理员已被禁用")
		return
	}

	// 获取角色信息
	roleTitle := ""
	if manager.Role.Id > 0 {
		roleTitle = manager.Role.Title
	}

	// 生成 JWT
	token, err := common.GenerateAdminToken(manager.Id, manager.Username, manager.RoleId, roleTitle)
	if err != nil {
		common.Error(c, 500, "生成令牌失败")
		return
	}

	common.Success(c, gin.H{
		"token":     token,
		"userId":    manager.Id,
		"username":  manager.Username,
		"roleId":    manager.RoleId,
		"roleTitle": roleTitle,
	})
}

func (con AuthController) Logout(c *gin.Context) {
	common.Success(c, nil)
}

func (con AuthController) CurrentUser(c *gin.Context) {
	adminId, _ := c.Get("adminId")
	username, _ := c.Get("username")
	roleId, _ := c.Get("roleId")
	roleTitle, _ := c.Get("roleTitle")

	common.Success(c, gin.H{
		"userId":    adminId,
		"username":  username,
		"roleId":    roleId,
		"roleTitle": roleTitle,
	})
}
