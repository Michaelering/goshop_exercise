package middlewares

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// 超级管理员角色ID
const SuperAdminRoleId = 1

// 无需权限校验的白名单路径（不含 /api/v1/admin/ 前缀）
var whitelistPaths = map[string]bool{
	"logout":               true,
	"currentUser":          true,
	"dashboard":            true,
	"menu":                 true,
	"changeStatus":         true,
	"changeNum":            true,
	"flushAll":             true,
	"goods/uploadImage":    true,
	"goods/uploadEditorImage": true,
	"goodsImage/color":     true,
}

// JwtAdminAuthMiddleware JWT 认证 + RBAC 权限校验中间件
func JwtAdminAuthMiddleware(c *gin.Context) {
	// 1. 从 Authorization header 获取 token
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

	// 2. 解析 JWT
	tokenString := parts[1]
	claims, err := common.ParseAdminToken(tokenString)
	if err != nil {
		common.Unauthorized(c, "认证令牌无效或已过期")
		c.Abort()
		return
	}

	// 3. 将用户信息注入 context
	c.Set("adminId", claims.UserId)
	c.Set("username", claims.Username)
	c.Set("roleId", claims.RoleId)
	c.Set("roleTitle", claims.RoleTitle)
	// 4. 获取有效的 roleId（兼容旧 token 中 roleId=0 的情况）
	effectiveRoleId := claims.RoleId
	if effectiveRoleId == 0 && claims.UserId > 0 {
		// 旧 token（is_super 用户升级前），从 DB 获取真实 roleId
		var mgr models.Manager
		models.DB.Where("id = ?", claims.UserId).Find(&mgr)
		if mgr.RoleId > 0 {
			effectiveRoleId = mgr.RoleId
			c.Set("roleId", mgr.RoleId)
		}
	}

	// 5. 超级管理员 roleId=1 直接放行（不论 isBuiltin，因为旧 token 可能缺失此字段）
	if effectiveRoleId == SuperAdminRoleId {
		c.Next()
		return
	}

	// 5. 提取 URL 第一段路径作为模块前缀
	urlPath := strings.TrimPrefix(c.Request.URL.Path, "/api/v1/admin/")
	urlPath = strings.TrimPrefix(urlPath, "/")

	// 如果路径为空（根路径），放行
	if urlPath == "" {
		c.Next()
		return
	}

	// 6. 白名单路径直接放行
	if whitelistPaths[urlPath] {
		c.Next()
		return
	}

	// 7. 提取第一段路径作为模块前缀
	// 例如: "goods/123" → "goods", "goodsCate/tree" → "goodsCate", "goodsImage/5" → "goodsImage"
	modulePrefix := extractModulePrefix(urlPath)

	// 8. 查找该模块的 access 记录
	access := models.Access{}
	models.DB.Where("url_prefix = ? AND type = 1 AND status = 1", modulePrefix).Find(&access)
	if access.Id == 0 {
		// 模块未在 access 表中注册，尝试匹配多段前缀
		// 例如 "goods/uploadImage" → 在 whitelist 中已处理，此处为兜底
		common.Forbidden(c, "没有权限访问该资源")
		c.Abort()
		return
	}

	// 9. 检查角色是否拥有此模块权限
	roleAccess := models.RoleAccess{}
	models.DB.Where("role_id = ? AND access_id = ?", effectiveRoleId, access.Id).Find(&roleAccess)
	if roleAccess.AccessId == 0 {
		common.Forbidden(c, "没有权限")
		c.Abort()
		return
	}

	// 10. 检查 HTTP 方法是否在允许列表中
	if access.HttpMethods != "" {
		if !strings.Contains(access.HttpMethods, c.Request.Method) {
			common.Forbidden(c, "没有该操作的权限")
			c.Abort()
			return
		}
	}

	c.Next()
}

// extractModulePrefix 提取 URL 的第一段路径作为模块前缀
// "goods" → "goods"
// "goods/123" → "goods"
// "goodsCate/tree" → "goodsCate"
// "goodsImage/5" → "goodsImage"
// "role/1/auth" → "role"
func extractModulePrefix(urlPath string) string {
	idx := strings.Index(urlPath, "/")
	if idx == -1 {
		return urlPath
	}
	return urlPath[:idx]
}
