package routers

import (
	"ginshop58/controllers/merchant"
	"ginshop58/middlewares"

	"github.com/gin-gonic/gin"
)

func MerchantRoutersInit(r *gin.Engine) {
	merchantRouters := r.Group("/merchant", middlewares.InitMerchantAuthMiddleware)
	{
		merchantRouters.GET("/", merchant.MainController{}.Index)
		merchantRouters.GET("/welcome", merchant.MainController{}.Welcome)

		merchantRouters.GET("/login", merchant.LoginController{}.Index)
		merchantRouters.GET("/captcha", merchant.LoginController{}.Captcha)
		merchantRouters.POST("/doLogin", merchant.LoginController{}.DoLogin)
		merchantRouters.GET("/loginOut", merchant.LoginController{}.LoginOut)

		merchantRouters.GET("/goods", merchant.GoodsController{}.Index)
		merchantRouters.GET("/goods/add", merchant.GoodsController{}.Add)
		merchantRouters.POST("/goods/doAdd", merchant.GoodsController{}.DoAdd)
		merchantRouters.GET("/goods/edit", merchant.GoodsController{}.Edit)
		merchantRouters.POST("/goods/doEdit", merchant.GoodsController{}.DoEdit)
		merchantRouters.GET("/goods/delete", merchant.GoodsController{}.Delete)
		merchantRouters.GET("/goods/goodsTypeAttribute", merchant.GoodsController{}.GoodsTypeAttribute)
		merchantRouters.POST("/goods/editorImageUpload", merchant.GoodsController{}.EditorImageUpload)
		merchantRouters.POST("/goods/goodsImageUpload", merchant.GoodsController{}.GoodsImageUpload)
		merchantRouters.GET("/goods/changeGoodsImageColor", merchant.GoodsController{}.ChangeGoodsImageColor)
		merchantRouters.GET("/goods/removeGoodsImage", merchant.GoodsController{}.RemoveGoodsImage)

		merchantRouters.GET("/order", merchant.OrderController{}.Index)
	}
}
