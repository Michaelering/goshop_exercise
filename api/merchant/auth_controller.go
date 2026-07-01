package merchant

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (con AuthController) Login(c *gin.Context) {
	captchaId := c.PostForm("captchaId")
	username := c.PostForm("username")
	password := c.PostForm("password")
	verifyValue := c.PostForm("verifyValue")

	if !models.VerifyCaptcha(captchaId, verifyValue) {
		common.BadRequest(c, "验证码验证失败")
		return
	}

	merchantList := []models.Merchant{}
	password = models.Md5(password)
	models.DB.Where("username=? AND password=? AND status=1", username, password).Find(&merchantList)

	if len(merchantList) == 0 {
		common.Error(c, 400, "用户名或者密码错误，或账号已被禁用")
		return
	}

	merchant := merchantList[0]

	token, err := common.GenerateMerchantToken(merchant.Id, merchant.Username, merchant.ShopName)
	if err != nil {
		common.Error(c, 500, "生成令牌失败")
		return
	}

	// 同步写入 session 兼容旧商户页面
	session := sessions.Default(c)
	sessionData, _ := models.JsonEncode(merchantList)
	if sessionData != "" {
		session.Set("merchantInfo", sessionData)
		session.Save()
	}

	common.Success(c, gin.H{
		"token":    token,
		"userId":   merchant.Id,
		"username": merchant.Username,
		"shopName": merchant.ShopName,
	})
}

func (con AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("merchantInfo")
	session.Save()
	common.Success(c, nil)
}

func (con AuthController) CurrentUser(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	username, _ := c.Get("username")
	shopName, _ := c.Get("shopName")

	common.Success(c, gin.H{
		"userId":   merchantId,
		"username": username,
		"shopName": shopName,
	})
}
