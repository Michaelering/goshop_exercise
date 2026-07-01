package merchant

import (
	ginshop58_admin "ginshop58/api/admin"
	"ginshop58/api/common"
	"ginshop58/models"
	"math"

	"github.com/gin-gonic/gin"
)

type GoodsController struct{}

func (con GoodsController) Index(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	page, _ := models.Int(c.DefaultQuery("page", "1"))
	pageSize, _ := models.Int(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	where := "is_delete=0 AND merchant_id=?"
	var args []interface{}
	args = append(args, merchantId)
	if keyword != "" {
		where += " AND title LIKE ?"
		args = append(args, "%"+keyword+"%")
	}

	goodsList := []models.Goods{}
	var count int64
	models.DB.Where(where, args...).Offset((page-1)*pageSize).Limit(pageSize).Order("id desc").Find(&goodsList)
	models.DB.Where(where, args...).Table("goods").Count(&count)

	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	common.List(c, gin.H{
		"list":       goodsList,
		"totalPages": totalPages,
	}, count, page, pageSize)
}

func (con GoodsController) Create(c *gin.Context) {
	merchantId := getMerchantId(c)
	params := ginshop58_admin.ParseGoodsParams(c, merchantId)

	goods, err := ginshop58_admin.CreateGoods(params)
	if err != nil {
		common.Error(c, 500, "增加商品失败")
		return
	}

	goodsImg, _ := models.UploadImg(c, "goods_img")
	if len(goodsImg) > 0 {
		goods.GoodsImg = goodsImg
		if models.GetOssStatus() != 1 {
			go func() {
				models.ResizeGoodsImage(goodsImg)
			}()
		}
		models.DB.Save(&goods)
	}

	common.Success(c, goods)
}

func (con GoodsController) Update(c *gin.Context) {
	merchantId := getMerchantId(c)
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	// 验证归属
	goods := models.Goods{Id: id}
	models.DB.Where("merchant_id=?", merchantId).Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "商品不存在或不属于当前商户")
		return
	}

	params := ginshop58_admin.ParseGoodsParams(c, merchantId)
	params.GoodsImg, _ = models.UploadImg(c, "goods_img")
	if params.GoodsImg != "" && models.GetOssStatus() != 1 {
		go func() {
			models.ResizeGoodsImage(params.GoodsImg)
		}()
	}

	err = ginshop58_admin.UpdateGoods(id, params)
	if err != nil {
		common.Error(c, 500, "修改商品失败")
		return
	}
	common.Success(c, nil)
}

func (con GoodsController) Delete(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goods := models.Goods{Id: id}
	models.DB.Where("merchant_id=?", merchantId).Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "商品不存在或不属于当前商户")
		return
	}
	goods.IsDelete = 1
	goods.Status = 0
	models.DB.Save(&goods)
	common.Success(c, nil)
}

func (con GoodsController) Get(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goods := models.Goods{Id: id}
	models.DB.Where("merchant_id=?", merchantId).Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "商品不存在")
		return
	}

	// 获取关联数据（简化版，仅返回商品基本信息和图片）
	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsImageList)

	common.Success(c, gin.H{
		"goods":          goods,
		"goodsImageList": goodsImageList,
	})
}

func (con GoodsController) UploadImage(c *gin.Context) {
	imgDir, err := models.UploadImg(c, "file")
	if err != nil {
		common.Error(c, 500, "上传失败")
		return
	}

	if models.GetOssStatus() != 1 {
		go func() {
			models.ResizeGoodsImage(imgDir)
		}()
		common.Success(c, gin.H{"link": "/" + imgDir})
	} else {
		common.Success(c, gin.H{"link": models.GetSettingFromColumn("OssDomain") + imgDir})
	}
}

func (con GoodsController) UploadEditorImage(c *gin.Context) {
	imgDir, err := models.UploadImg(c, "file")
	if err != nil {
		common.Error(c, 500, "上传失败")
		return
	}

	if models.GetOssStatus() != 1 {
		go func() {
			models.ResizeGoodsImage(imgDir)
		}()
		common.Success(c, gin.H{"link": "/" + imgDir})
	} else {
		common.Success(c, gin.H{"link": models.GetSettingFromColumn("OssDomain") + imgDir})
	}
}

func (con GoodsController) ChangeImageColor(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	goodsImageId, _ := models.Int(c.PostForm("goods_image_id"))
	colorId, _ := models.Int(c.PostForm("color_id"))
	if goodsImageId == 0 || colorId == 0 {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsImage := models.GoodsImage{Id: goodsImageId}
	models.DB.Find(&goodsImage)

	// 验证图片归属
	var goods models.Goods
	models.DB.Where("id=? AND merchant_id=?", goodsImage.GoodsId, merchantId).Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "无权操作此图片")
		return
	}

	goodsImage.ColorId = colorId
	models.DB.Save(&goodsImage)
	common.Success(c, nil)
}

func (con GoodsController) RemoveImage(c *gin.Context) {
	merchantId, _ := c.Get("merchantId")
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsImage := models.GoodsImage{Id: id}
	models.DB.Find(&goodsImage)

	var goods models.Goods
	models.DB.Where("id=? AND merchant_id=?", goodsImage.GoodsId, merchantId).Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "无权操作此图片")
		return
	}

	models.DB.Delete(&goodsImage)
	common.Success(c, nil)
}

// 安全从 context 获取 merchantId
func getMerchantId(c *gin.Context) int {
	if v, ok := c.Get("merchantId"); ok {
		if id, ok := v.(int); ok {
			return id
		}
	}
	return 0
}
