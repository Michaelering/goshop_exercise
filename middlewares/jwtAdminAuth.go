package middlewares

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAdminAuthMiddleware(c *gin.Context) {
	// 从 Authorization header 获取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		common.Unauthorized(c, "未提供认证令牌")
		c.Abort()
		return
	}

	// 格式: "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		common.Unauthorized(c, "认证令牌格式错误")
		c.Abort()
		return
	}

	tokenString := parts[1]
	claims, err := common.ParseAdminToken(tokenString)
	if err != nil {
		common.Unauthorized(c, "认证令牌无效或已过期")
		c.Abort()
		return
	}

	// 将用户信息存入 context
	c.Set("adminId", claims.UserId)
	c.Set("username", claims.Username)
	c.Set("roleId", claims.RoleId)
	c.Set("isSuper", claims.IsSuper)

	// RBAC 权限校验（非超级管理员时）
	if claims.IsSuper == 0 {
		// 获取当前请求路径，去掉 /api/v1/admin/ 前缀
		urlPath := strings.Replace(c.Request.URL.Path, "/api/v1/admin/", "", 1)
		// 排除不需要权限的路径
		if !excludeApiAuthPath(urlPath) {
			// 获取角色的权限
			roleAccess := []models.RoleAccess{}
			models.DB.Where("role_id=?", claims.RoleId).Find(&roleAccess)
			roleAccessMap := make(map[int]int)
			for _, v := range roleAccess {
				roleAccessMap[v.AccessId] = v.AccessId
			}

			// 查找当前路径是否有对应的权限节点
			access := models.Access{}
			models.DB.Where("url = ?", urlPath).Find(&access)
			if _, ok := roleAccessMap[access.Id]; !ok {
				common.Forbidden(c, "没有权限")
				c.Abort()
				return
			}
		}
	}

	c.Next()
}

// API 中不需要权限校验的路径
func excludeApiAuthPath(urlPath string) bool {
	excludePaths := []string{"", "login", "captcha", "loginOut", "welcome"}
	for _, v := range excludePaths {
		if v == urlPath {
			return true
		}
	}
	return false
}
