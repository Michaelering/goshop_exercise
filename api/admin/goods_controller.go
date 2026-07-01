package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"math"

	"github.com/gin-gonic/gin"
)

type GoodsController struct{}

func (con GoodsController) Index(c *gin.Context) {
	page, _ := models.Int(c.DefaultQuery("page", "1"))
	pageSize, _ := models.Int(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	where := "is_delete=0"
	var args []interface{}
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
	params := ParseGoodsParams(c, 0)

	goods, err := CreateGoods(params)
	if err != nil {
		common.Error(c, 500, "增加商品失败")
		return
	}

	// 处理图片上传
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
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	params := ParseGoodsParams(c, 0)
	params.GoodsImg, _ = models.UploadImg(c, "goods_img")
	if params.GoodsImg != "" && models.GetOssStatus() != 1 {
		go func() {
			models.ResizeGoodsImage(params.GoodsImg)
		}()
	}

	err = UpdateGoods(id, params)
	if err != nil {
		common.Error(c, 500, "修改商品失败")
		return
	}
	common.Success(c, nil)
}

func (con GoodsController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "商品不存在")
		return
	}
	goods.IsDelete = 1
	goods.Status = 0
	models.DB.Save(&goods)
	common.Success(c, nil)
}

func (con GoodsController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	if goods.Id == 0 {
		common.BadRequest(c, "商品不存在")
		return
	}

	// 获取关联数据
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)
	goodsColorSlice := splitString(goods.GoodsColor, ",")
	for i := 0; i < len(goodsColorList); i++ {
		for _, cid := range goodsColorSlice {
			if models.String(goodsColorList[i].Id) == cid {
				goodsColorList[i].Checked = true
				break
			}
		}
	}

	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsImageList)

	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	goodsAttrStr := BuildGoodsAttrHtml(goods.Id)

	common.Success(c, gin.H{
		"goods":          goods,
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
		"goodsAttrStr":   goodsAttrStr,
		"goodsImageList": goodsImageList,
	})
}

func (con GoodsController) ToggleStatus(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	field := c.PostForm("field")
	if field == "" {
		field = "status"
	}

	err = models.DB.Exec("UPDATE goods SET "+field+"=ABS("+field+"-1) WHERE id=?", id).Error
	if err != nil {
		common.Error(c, 500, "修改失败")
		return
	}
	common.Success(c, nil)
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
	goodsImageId, _ := models.Int(c.PostForm("goods_image_id"))
	colorId, _ := models.Int(c.PostForm("color_id"))
	if goodsImageId == 0 || colorId == 0 {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsImage := models.GoodsImage{Id: goodsImageId}
	models.DB.Find(&goodsImage)
	if goodsImage.Id == 0 {
		common.BadRequest(c, "图片不存在")
		return
	}
	goodsImage.ColorId = colorId
	models.DB.Save(&goodsImage)
	common.Success(c, nil)
}

func (con GoodsController) RemoveImage(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.GoodsImage{Id: id})
	common.Success(c, nil)
}

// 解析商品表单参数（通用，admin 和 merchant 共用）
func ParseGoodsParams(c *gin.Context, merchantId int) *GoodsCreateParams {
	return &GoodsCreateParams{
		MerchantId:    merchantId,
		Title:         c.PostForm("title"),
		SubTitle:      c.PostForm("sub_title"),
		GoodsSn:       c.PostForm("goods_sn"),
		CateId:        getFormInt(c, "cate_id"),
		GoodsNumber:   getFormInt(c, "goods_number"),
		MarketPrice:   getFormFloat(c, "market_price"),
		Price:         getFormFloat(c, "price"),
		RelationGoods: c.PostForm("relation_goods"),
		GoodsAttr:     c.PostForm("goods_attr"),
		GoodsVersion:  c.PostForm("goods_version"),
		GoodsGift:     c.PostForm("goods_gift"),
		GoodsFitting:  c.PostForm("goods_fitting"),
		GoodsColor:    c.PostFormArray("goods_color"),
		GoodsKeywords: c.PostForm("goods_keywords"),
		GoodsDesc:     c.PostForm("goods_desc"),
		GoodsContent:  c.PostForm("goods_content"),
		IsHot:         getFormInt(c, "is_hot"),
		IsBest:        getFormInt(c, "is_best"),
		IsNew:         getFormInt(c, "is_new"),
		GoodsTypeId:   getFormInt(c, "goods_type_id"),
		Sort:          getFormInt(c, "sort"),
		Status:        getFormInt(c, "status"),
		GoodsImageList: c.PostFormArray("goods_image_list"),
		AttrIdList:     c.PostFormArray("attr_id_list"),
		AttrValueList:  c.PostFormArray("attr_value_list"),
	}
}

func getFormInt(c *gin.Context, key string) int {
	val, _ := models.Int(c.PostForm(key))
	return val
}

func getFormFloat(c *gin.Context, key string) float64 {
	val, _ := models.Float(c.PostForm(key))
	return val
}

func splitString(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	for _, v := range split(s, sep) {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func split(s string, sep string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if string(c) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	result = append(result, current)
	return result
}
