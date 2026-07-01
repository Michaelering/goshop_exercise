package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 标准 JSON 响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}

// 分页列表响应
func List(c *gin.Context, data interface{}, total int64, page int, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"message":  "success",
		"data":     data,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未登录或登录已过期"
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    401,
		"message": message,
		"data":    nil,
	})
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "没有权限"
	}
	c.JSON(http.StatusForbidden, gin.H{
		"code":    403,
		"message": message,
		"data":    nil,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":    400,
		"message": message,
		"data":    nil,
	})
}
