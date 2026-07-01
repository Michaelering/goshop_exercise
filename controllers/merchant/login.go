package merchant

import (
	"encoding/json"
	"fmt"
	"ginshop58/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
}

func (con LoginController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "merchant/login/login.html", gin.H{})
}

func (con LoginController) DoLogin(c *gin.Context) {
	captchaId := c.PostForm("captchaId")
	username := c.PostForm("username")
	password := c.PostForm("password")
	verifyValue := c.PostForm("verifyValue")

	if flag := models.VerifyCaptcha(captchaId, verifyValue); flag {
		merchantList := []models.Merchant{}
		password = models.Md5(password)

		models.DB.Where("username=? AND password=? AND status=1", username, password).Find(&merchantList)

		if len(merchantList) > 0 {
			session := sessions.Default(c)
			merchantSlice, _ := json.Marshal(merchantList)
			session.Set("merchantInfo", string(merchantSlice))
			session.Save()
			con.Success(c, "登录成功", "/merchant")
		} else {
			con.Error(c, "用户名或者密码错误，或账号已被禁用", "/merchant/login")
		}
	} else {
		con.Error(c, "验证码验证失败", "/merchant/login")
	}
}

func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha(34, 100, 2)

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con LoginController) LoginOut(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("merchantInfo")
	session.Save()
	con.Success(c, "退出登录成功", "/merchant/login")
}
