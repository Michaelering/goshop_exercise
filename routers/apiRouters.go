package routers

import (
	"ginshop58/api/admin"
	"ginshop58/api/merchant"
	"ginshop58/middlewares"

	"github.com/gin-gonic/gin"
)

func ApiRoutersInit(r *gin.Engine) {
	apiGroup := r.Group("/api/v1")
	{
		// ====== 公共接口（无需认证）======
		apiGroup.GET("/public/captcha", admin.AuthController{}.Captcha)

		// ====== 管理端接口 ======
		adminGroup := apiGroup.Group("/admin")
		{
			// 登录（无需 JWT）
			adminGroup.POST("/login", admin.AuthController{}.Login)

			// 以下需要 JWT 认证
			adminAuth := adminGroup.Group("", middlewares.JwtAdminAuthMiddleware)
			{
				adminAuth.GET("/logout", admin.AuthController{}.Logout)
				adminAuth.GET("/currentUser", admin.AuthController{}.CurrentUser)

				// 仪表盘
				adminAuth.GET("/dashboard", admin.DashboardController{}.Index)
				adminAuth.GET("/menu", admin.DashboardController{}.Menu)

				// 管理员管理
				adminAuth.GET("/manager", admin.ManagerController{}.Index)
				adminAuth.POST("/manager", admin.ManagerController{}.Create)
				adminAuth.PUT("/manager/:id", admin.ManagerController{}.Update)
				adminAuth.DELETE("/manager/:id", admin.ManagerController{}.Delete)
				adminAuth.GET("/manager/:id", admin.ManagerController{}.Get)

				// 角色管理
				adminAuth.GET("/role", admin.RoleController{}.Index)
				adminAuth.POST("/role", admin.RoleController{}.Create)
				adminAuth.PUT("/role/:id", admin.RoleController{}.Update)
				adminAuth.DELETE("/role/:id", admin.RoleController{}.Delete)
				adminAuth.GET("/role/:id", admin.RoleController{}.Get)
				adminAuth.GET("/role/:id/auth", admin.RoleController{}.Auth)
				adminAuth.POST("/role/:id/auth", admin.RoleController{}.DoAuth)

				// 权限管理
				adminAuth.GET("/access", admin.AccessController{}.Index)
				adminAuth.POST("/access", admin.AccessController{}.Create)
				adminAuth.PUT("/access/:id", admin.AccessController{}.Update)
				adminAuth.DELETE("/access/:id", admin.AccessController{}.Delete)
				adminAuth.GET("/access/:id", admin.AccessController{}.Get)
				adminAuth.GET("/access/topModules", admin.AccessController{}.TopModules)

				// 商品分类
				adminAuth.GET("/goodsCate", admin.GoodsCateController{}.Index)
				adminAuth.POST("/goodsCate", admin.GoodsCateController{}.Create)
				adminAuth.PUT("/goodsCate/:id", admin.GoodsCateController{}.Update)
				adminAuth.DELETE("/goodsCate/:id", admin.GoodsCateController{}.Delete)
				adminAuth.GET("/goodsCate/:id", admin.GoodsCateController{}.Get)
				adminAuth.GET("/goodsCate/tree", admin.GoodsCateController{}.Tree)

				// 商品类型
				adminAuth.GET("/goodsType", admin.GoodsTypeController{}.Index)
				adminAuth.POST("/goodsType", admin.GoodsTypeController{}.Create)
				adminAuth.PUT("/goodsType/:id", admin.GoodsTypeController{}.Update)
				adminAuth.DELETE("/goodsType/:id", admin.GoodsTypeController{}.Delete)
				adminAuth.GET("/goodsType/:id", admin.GoodsTypeController{}.Get)

				// 商品类型属性
				adminAuth.GET("/goodsTypeAttr", admin.GoodsTypeAttrController{}.Index)
				adminAuth.POST("/goodsTypeAttr", admin.GoodsTypeAttrController{}.Create)
				adminAuth.PUT("/goodsTypeAttr/:id", admin.GoodsTypeAttrController{}.Update)
				adminAuth.DELETE("/goodsTypeAttr/:id", admin.GoodsTypeAttrController{}.Delete)
				adminAuth.GET("/goodsTypeAttr/:id", admin.GoodsTypeAttrController{}.Get)
				adminAuth.GET("/goodsTypeAttr/byType/:cateId", admin.GoodsTypeAttrController{}.ByType)

				// 商品管理
				adminAuth.GET("/goods", admin.GoodsController{}.Index)
				adminAuth.POST("/goods", admin.GoodsController{}.Create)
				adminAuth.PUT("/goods/:id", admin.GoodsController{}.Update)
				adminAuth.DELETE("/goods/:id", admin.GoodsController{}.Delete)
				adminAuth.GET("/goods/:id", admin.GoodsController{}.Get)
				adminAuth.POST("/goods/:id/status", admin.GoodsController{}.ToggleStatus)
				adminAuth.POST("/goods/uploadImage", admin.GoodsController{}.UploadImage)
				adminAuth.POST("/goods/uploadEditorImage", admin.GoodsController{}.UploadEditorImage)
				adminAuth.POST("/goodsImage/color", admin.GoodsController{}.ChangeImageColor)
				adminAuth.DELETE("/goodsImage/:id", admin.GoodsController{}.RemoveImage)

				// 导航管理
				adminAuth.GET("/nav", admin.NavController{}.Index)
				adminAuth.POST("/nav", admin.NavController{}.Create)
				adminAuth.PUT("/nav/:id", admin.NavController{}.Update)
				adminAuth.DELETE("/nav/:id", admin.NavController{}.Delete)
				adminAuth.GET("/nav/:id", admin.NavController{}.Get)

				// 轮播图管理
				adminAuth.GET("/focus", admin.FocusController{}.Index)
				adminAuth.POST("/focus", admin.FocusController{}.Create)
				adminAuth.PUT("/focus/:id", admin.FocusController{}.Update)
				adminAuth.DELETE("/focus/:id", admin.FocusController{}.Delete)
				adminAuth.GET("/focus/:id", admin.FocusController{}.Get)

				// 系统设置
				adminAuth.GET("/setting", admin.SettingController{}.Index)
				adminAuth.PUT("/setting", admin.SettingController{}.Update)

				// 商户管理（后台管理员管理商户）
				adminAuth.GET("/merchant", admin.MerchantMgmtController{}.Index)
				adminAuth.POST("/merchant", admin.MerchantMgmtController{}.Create)
				adminAuth.PUT("/merchant/:id", admin.MerchantMgmtController{}.Update)
				adminAuth.DELETE("/merchant/:id", admin.MerchantMgmtController{}.Delete)
				adminAuth.GET("/merchant/:id", admin.MerchantMgmtController{}.Get)

				// 订单管理（管理员查看/管理全部订单）
				adminAuth.GET("/order", admin.OrderController{}.Index)
				adminAuth.GET("/order/:id", admin.OrderController{}.Get)
				adminAuth.PUT("/order/:id", admin.OrderController{}.Update)

				// 公共修改状态和排序
				adminAuth.PUT("/changeStatus", admin.DashboardController{}.ChangeStatus)
				adminAuth.PUT("/changeNum", admin.DashboardController{}.ChangeNum)
				adminAuth.POST("/flushAll", admin.DashboardController{}.FlushAll)
			}
		}

		// ====== 商户端接口 ======
		merchantGroup := apiGroup.Group("/merchant")
		{
			// 登录
			merchantGroup.POST("/login", merchant.AuthController{}.Login)

			// 需要 JWT 认证
			merchantAuth := merchantGroup.Group("", middlewares.JwtMerchantAuthMiddleware)
			{
				merchantAuth.GET("/logout", merchant.AuthController{}.Logout)
				merchantAuth.GET("/currentUser", merchant.AuthController{}.CurrentUser)

				// 仪表盘
				merchantAuth.GET("/dashboard", merchant.DashboardController{}.Index)

				// 商品管理
				merchantAuth.GET("/goods", merchant.GoodsController{}.Index)
				merchantAuth.POST("/goods", merchant.GoodsController{}.Create)
				merchantAuth.PUT("/goods/:id", merchant.GoodsController{}.Update)
				merchantAuth.DELETE("/goods/:id", merchant.GoodsController{}.Delete)
				merchantAuth.GET("/goods/:id", merchant.GoodsController{}.Get)
				merchantAuth.POST("/goods/uploadImage", merchant.GoodsController{}.UploadImage)
				merchantAuth.POST("/goods/uploadEditorImage", merchant.GoodsController{}.UploadEditorImage)
				merchantAuth.POST("/goodsImage/color", merchant.GoodsController{}.ChangeImageColor)
				merchantAuth.DELETE("/goodsImage/:id", merchant.GoodsController{}.RemoveImage)

				// 订单管理
				merchantAuth.GET("/order", merchant.OrderController{}.Index)
			}
		}
	}
}
