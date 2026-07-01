package middlewares

import (
	"encoding/json"
	"ginshop58/models"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitMerchantAuthMiddleware(c *gin.Context) {

	pathname := strings.Split(c.Request.URL.String(), "?")[0]
	session := sessions.Default(c)
	merchantInfo := session.Get("merchantInfo")
	merchantInfoStr, ok := merchantInfo.(string)

	if ok {
		var merchantInfoStruct []models.Merchant
		err := json.Unmarshal([]byte(merchantInfoStr), &merchantInfoStruct)
		if err != nil || !(len(merchantInfoStruct) > 0 && merchantInfoStruct[0].Username != "") {
			if pathname != "/merchant/login" && pathname != "/merchant/doLogin" && pathname != "/merchant/captcha" {
				c.Redirect(302, "/merchant/login")
			}
		} else {
			// 商户已登录，继续执行
			c.Next()
		}
	} else {
		if pathname != "/merchant/login" && pathname != "/merchant/doLogin" && pathname != "/merchant/captcha" {
			c.Redirect(302, "/merchant/login")
		}
	}
}
