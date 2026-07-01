package main

import (
	"ginshop58/models"
	"ginshop58/routers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 数据库迁移 + 播种内置角色和权限数据
	models.RunMigrationAndSeed()

	r := gin.Default()

	// 配置 gin 允许跨域请求
	r.Use(cors.Default())

	// 配置静态文件服务 - 仅保留 upload 目录（本地图片上传）
	r.Static("/static/upload", "./static/upload")

	// 注册 API 路由
	routers.ApiRoutersInit(r)

	// 生产环境：服务 Vue SPA 构建产物
	distPath := "./admin_frontend/dist"
	if _, err := os.Stat(distPath); err == nil {
		// 服务 assets 目录（JS/CSS）
		r.Static("/assets", filepath.Join(distPath, "assets"))
		// 服务 favicon
		r.StaticFile("/favicon.ico", filepath.Join(distPath, "favicon.ico"))
		// SPA 回退：所有非 API 的 GET 请求返回 index.html
		r.NoRoute(func(c *gin.Context) {
			if c.Request.Method == "GET" {
				c.File(filepath.Join(distPath, "index.html"))
			} else {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
			}
		})
	}

	r.Run()
}
