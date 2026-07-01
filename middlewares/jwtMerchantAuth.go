package middlewares

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtMerchantAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		common.Unauthorized(c, "未提供认证令牌")
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		common.Unauthorized(c, "认证令牌格式错误")
		c.Abort()
		return
	}

	tokenString := parts[1]
	claims, err := common.ParseMerchantToken(tokenString)
	if err != nil {
		common.Unauthorized(c, "认证令牌无效或已过期")
		c.Abort()
		return
	}

	// 验证商户状态是否正常
	merchant := models.Merchant{Id: claims.UserId}
	models.DB.Find(&merchant)
	if merchant.Status != 1 {
		common.Forbidden(c, "商户已被禁用")
		c.Abort()
		return
	}

	// 将商户信息存入 context
	c.Set("merchantId", claims.UserId)
	c.Set("username", claims.Username)
	c.Set("shopName", claims.ShopName)

	c.Next()
}
