package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-contrib/sessions"
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
	models.DB.Where("username=? AND password=?", username, password).Find(&managerList)

	if len(managerList) == 0 {
		common.Error(c, 400, "用户名或者密码错误")
		return
	}

	manager := managerList[0]
	if manager.Status == 0 {
		common.Error(c, 400, "该管理员已被禁用")
		return
	}

	// 生成 JWT
	token, err := common.GenerateAdminToken(manager.Id, manager.Username, manager.RoleId, manager.IsSuper)
	if err != nil {
		common.Error(c, 500, "生成令牌失败")
		return
	}

	// 同时也写入 session 以兼容旧的后台页面
	session := sessions.Default(c)
	sessionData, _ := models.JsonEncode(managerList)
	if sessionData != "" {
		session.Set("userinfo", sessionData)
		session.Save()
	}

	common.Success(c, gin.H{
		"token":      token,
		"userId":     manager.Id,
		"username":   manager.Username,
		"isSuper":    manager.IsSuper,
		"roleId":     manager.RoleId,
	})
}

func (con AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("userinfo")
	session.Save()
	common.Success(c, nil)
}

func (con AuthController) CurrentUser(c *gin.Context) {
	adminId, _ := c.Get("adminId")
	username, _ := c.Get("username")
	isSuper, _ := c.Get("isSuper")
	roleId, _ := c.Get("roleId")

	// 获取角色名称
	roleTitle := ""
	if roleIdVal, ok := roleId.(int); ok && roleIdVal != 0 {
		role := models.Role{Id: roleIdVal}
		models.DB.Find(&role)
		roleTitle = role.Title
	}

	common.Success(c, gin.H{
		"userId":    adminId,
		"username":  username,
		"isSuper":   isSuper,
		"roleId":    roleId,
		"roleTitle": roleTitle,
	})
}

// 用于 JSON 序列化辅助（给 models 包加一个 JsonEncode 方法）
func init() {
	// 注册 JsonEncode 到 models 包（在 models/tools.go 中已存在类似功能）
}

// 注意：这里用到的 models.JsonEncode 需要补充，如果 models 中没有的话
// 使用 encoding/json 替代
