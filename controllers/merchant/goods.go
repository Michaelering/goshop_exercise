package merchant

import (
	"encoding/json"
	"fmt"
	"ginshop58/models"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

type GoodsController struct {
	BaseController
}

// 从 session 获取当前商户 ID
func getMerchantId(c *gin.Context) int {
	session := sessions.Default(c)
	merchantInfo := session.Get("merchantInfo")
	merchantInfoStr, _ := merchantInfo.(string)
	var merchantInfoStruct []models.Merchant
	json.Unmarshal([]byte(merchantInfoStr), &merchantInfoStruct)
	if len(merchantInfoStruct) > 0 {
		return merchantInfoStruct[0].Id
	}
	return 0
}

func (con GoodsController) Index(c *gin.Context) {
	merchantId := getMerchantId(c)
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}

	where := "is_delete=0 AND merchant_id=?"
	keyword := c.Query("keyword")
	if len(keyword) > 0 {
		where += " AND title like \"%" + keyword + "%\""
	}

	pageSize := 8
	goodsList := []models.Goods{}
	models.DB.Where(where, merchantId).Offset((page-1)*pageSize).Limit(pageSize).Order("id desc").Find(&goodsList)

	var count int64
	models.DB.Where(where, merchantId).Table("goods").Count(&count)

	if len(goodsList) > 0 {
		c.HTML(http.StatusOK, "merchant/goods/index.html", gin.H{
			"goodsList":  goodsList,
			"totalPages": math.Ceil(float64(count) / float64(pageSize)),
			"page":       page,
			"keyword":    keyword,
		})
	} else {
		if page != 1 {
			c.Redirect(302, "/merchant/goods")
		} else {
			c.HTML(http.StatusOK, "merchant/goods/index.html", gin.H{
				"goodsList":  goodsList,
				"totalPages": math.Ceil(float64(count) / float64(pageSize)),
				"page":       page,
				"keyword":    keyword,
			})
		}
	}
}

func (con GoodsController) Add(c *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)

	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	c.HTML(http.StatusOK, "merchant/goods/add.html", gin.H{
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
	})
}

func (con GoodsController) GoodsTypeAttribute(c *gin.Context) {
	cateId, err1 := models.Int(c.Query("cateId"))
	goodsTypeAttributeList := []models.GoodsTypeAttribute{}
	err2 := models.DB.Where("cate_id = ?", cateId).Find(&goodsTypeAttributeList).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"result":  "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  goodsTypeAttributeList,
		})
	}
}

func (con GoodsController) DoAdd(c *gin.Context) {
	merchantId := getMerchantId(c)

	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	goodsSn := c.PostForm("goods_sn")
	cateId, _ := models.Int(c.PostForm("cate_id"))
	goodsNumber, _ := models.Int(c.PostForm("goods_number"))
	marketPrice, _ := models.Float(c.PostForm("market_price"))
	price, _ := models.Float(c.PostForm("price"))
	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr")
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	goodsColorArr := c.PostFormArray("goods_color")
	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	isHot, _ := models.Int(c.PostForm("is_hot"))
	isBest, _ := models.Int(c.PostForm("is_best"))
	isNew, _ := models.Int(c.PostForm("is_new"))
	goodsTypeId, _ := models.Int(c.PostForm("goods_type_id"))
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))
	addTime := int(models.GetUnix())

	goodsColorStr := strings.Join(goodsColorArr, ",")

	goodsImg, _ := models.UploadImg(c, "goods_img")
	if len(goodsImg) > 0 {
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(goodsImg)
				wg.Done()
			}()
		}
	}

	goods := models.Goods{
		Title:         title,
		SubTitle:      subTitle,
		GoodsSn:       goodsSn,
		CateId:        cateId,
		ClickCount:    100,
		GoodsNumber:   goodsNumber,
		MarketPrice:   marketPrice,
		Price:         price,
		RelationGoods: relationGoods,
		GoodsAttr:     goodsAttr,
		GoodsVersion:  goodsVersion,
		GoodsGift:     goodsGift,
		GoodsFitting:  goodsFitting,
		GoodsKeywords: goodsKeywords,
		GoodsDesc:     goodsDesc,
		GoodsContent:  goodsContent,
		IsDelete:      0,
		IsHot:         isHot,
		IsBest:        isBest,
		IsNew:         isNew,
		GoodsTypeId:   goodsTypeId,
		Sort:          sort,
		Status:        status,
		AddTime:       addTime,
		GoodsColor:    goodsColorStr,
		GoodsImg:      goodsImg,
		MerchantId:    merchantId,
	}
	err := models.DB.Create(&goods).Error
	if err != nil {
		con.Error(c, "增加失败", "/merchant/goods/add")
		return
	}

	// 增加图库信息
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()

	// 增加规格包装
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, attributeIdErr := models.Int(attrIdList[i])
			if attributeIdErr == nil {
				goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
				models.DB.Find(&goodsTypeAttributeObj)
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
				goodsAttrObj.AttributeValue = attrValueList[i]
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	con.Success(c, "增加数据成功", "/merchant/goods")
}

func (con GoodsController) Edit(c *gin.Context) {
	merchantId := getMerchantId(c)

	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入参数错误", "/merchant/goods")
		return
	}

	goods := models.Goods{Id: id}
	models.DB.Where("merchant_id=?", merchantId).Find(&goods)
	if goods.Id == 0 {
		con.Error(c, "商品不存在或不属于当前商户", "/merchant/goods")
		return
	}

	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)

	goodsColorSlice := strings.Split(goods.GoodsColor, ",")
	goodsColorMap := make(map[string]string)
	for _, v := range goodsColorSlice {
		goodsColorMap[v] = v
	}
	goodsColorList := []models.GoodsColor{}
	models.DB.Find(&goodsColorList)
	for i := 0; i < len(goodsColorList); i++ {
		if _, ok := goodsColorMap[models.String(goodsColorList[i].Id)]; ok {
			goodsColorList[i].Checked = true
		}
	}

	goodsImageList := []models.GoodsImage{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsImageList)

	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	goodsAttr := []models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Find(&goodsAttr)
	goodsAttrStr := ""
	for _, v := range goodsAttr {
		if v.AttributeType == 1 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: </span> <input type="hidden" name="attr_id_list" value="%v" />   <input type="text" name="attr_value_list" value="%v" /></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else if v.AttributeType == 2 {
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: 　</span><input type="hidden" name="attr_id_list" value="%v" />  <textarea cols="50" rows="3" name="attr_value_list">%v</textarea></li>`, v.AttributeTitle, v.AttributeId, v.AttributeValue)
		} else {
			goodsTypeArttribute := models.GoodsTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&goodsTypeArttribute)
			attrValueSlice := strings.Split(goodsTypeArttribute.AttrValue, "\n")
			goodsAttrStr += fmt.Sprintf(`<li><span>%v: 　</span>  <input type="hidden" name="attr_id_list" value="%v" /> `, v.AttributeTitle, v.AttributeId)
			goodsAttrStr += fmt.Sprintf(`<select name="attr_value_list">`)
			for i := 0; i < len(attrValueSlice); i++ {
				if attrValueSlice[i] == v.AttributeValue {
					goodsAttrStr += fmt.Sprintf(`<option value="%v" selected >%v</option>`, attrValueSlice[i], attrValueSlice[i])
				} else {
					goodsAttrStr += fmt.Sprintf(`<option value="%v">%v</option>`, attrValueSlice[i], attrValueSlice[i])
				}
			}
			goodsAttrStr += fmt.Sprintf(`</select></li>`)
		}
	}

	c.HTML(http.StatusOK, "merchant/goods/edit.html", gin.H{
		"goods":          goods,
		"goodsCateList":  goodsCateList,
		"goodsColorList": goodsColorList,
		"goodsTypeList":  goodsTypeList,
		"goodsAttrStr":   goodsAttrStr,
		"goodsImageList": goodsImageList,
		"prevPage":       c.Request.Referer(),
	})
}

func (con GoodsController) DoEdit(c *gin.Context) {
	merchantId := getMerchantId(c)

	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入参数错误", "/merchant/goods")
		return
	}

	prevPage := c.PostForm("prevPage")

	// 验证商品归属
	goods := models.Goods{Id: id}
	models.DB.Where("merchant_id=?", merchantId).Find(&goods)
	if goods.Id == 0 {
		con.Error(c, "商品不存在或不属于当前商户", "/merchant/goods")
		return
	}

	title := c.PostForm("title")
	subTitle := c.PostForm("sub_title")
	goodsSn := c.PostForm("goods_sn")
	cateId, _ := models.Int(c.PostForm("cate_id"))
	goodsNumber, _ := models.Int(c.PostForm("goods_number"))
	marketPrice, _ := models.Float(c.PostForm("market_price"))
	price, _ := models.Float(c.PostForm("price"))
	relationGoods := c.PostForm("relation_goods")
	goodsAttr := c.PostForm("goods_attr")
	goodsVersion := c.PostForm("goods_version")
	goodsGift := c.PostForm("goods_gift")
	goodsFitting := c.PostForm("goods_fitting")
	goodsColorArr := c.PostFormArray("goods_color")
	goodsKeywords := c.PostForm("goods_keywords")
	goodsDesc := c.PostForm("goods_desc")
	goodsContent := c.PostForm("goods_content")
	isHot, _ := models.Int(c.PostForm("is_hot"))
	isBest, _ := models.Int(c.PostForm("is_best"))
	isNew, _ := models.Int(c.PostForm("is_new"))
	goodsTypeId, _ := models.Int(c.PostForm("goods_type_id"))
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	goodsColorStr := strings.Join(goodsColorArr, ",")

	goods.Title = title
	goods.SubTitle = subTitle
	goods.GoodsSn = goodsSn
	goods.CateId = cateId
	goods.GoodsNumber = goodsNumber
	goods.MarketPrice = marketPrice
	goods.Price = price
	goods.RelationGoods = relationGoods
	goods.GoodsAttr = goodsAttr
	goods.GoodsVersion = goodsVersion
	goods.GoodsGift = goodsGift
	goods.GoodsFitting = goodsFitting
	goods.GoodsKeywords = goodsKeywords
	goods.GoodsDesc = goodsDesc
	goods.GoodsContent = goodsContent
	goods.IsHot = isHot
	goods.IsBest = isBest
	goods.IsNew = isNew
	goods.GoodsTypeId = goodsTypeId
	goods.Sort = sort
	goods.Status = status
	goods.GoodsColor = goodsColorStr

	goodsImg, err2 := models.UploadImg(c, "goods_img")
	if err2 == nil && len(goodsImg) > 0 {
		goods.GoodsImg = goodsImg
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(goodsImg)
				wg.Done()
			}()
		}
	}

	err3 := models.DB.Save(&goods).Error
	if err3 != nil {
		con.Error(c, "修改失败", "/merchant/goods/edit?id="+models.String(id))
		return
	}

	// 修改图库
	wg.Add(1)
	go func() {
		goodsImageList := c.PostFormArray("goods_image_list")
		for _, v := range goodsImageList {
			goodsImgObj := models.GoodsImage{}
			goodsImgObj.GoodsId = goods.Id
			goodsImgObj.ImgUrl = v
			goodsImgObj.Sort = 10
			goodsImgObj.Status = 1
			goodsImgObj.AddTime = int(models.GetUnix())
			models.DB.Create(&goodsImgObj)
		}
		wg.Done()
	}()

	// 修改规格包装
	goodsAttrObj := models.GoodsAttr{}
	models.DB.Where("goods_id=?", goods.Id).Delete(&goodsAttrObj)
	wg.Add(1)
	go func() {
		attrIdList := c.PostFormArray("attr_id_list")
		attrValueList := c.PostFormArray("attr_value_list")
		for i := 0; i < len(attrIdList); i++ {
			goodsTypeAttributeId, attributeIdErr := models.Int(attrIdList[i])
			if attributeIdErr == nil {
				goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
				models.DB.Find(&goodsTypeAttributeObj)
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
				goodsAttrObj.AttributeValue = attrValueList[i]
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}
		}
		wg.Done()
	}()
	wg.Wait()

	if len(prevPage) > 0 {
		con.Success(c, "修改数据成功", prevPage)
	} else {
		con.Success(c, "修改数据成功", "/merchant/goods")
	}
}

func (con GoodsController) Delete(c *gin.Context) {
	merchantId := getMerchantId(c)
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/merchant/goods")
		return
	}

	goods := models.Goods{Id: id}
	models.DB.Where("merchant_id=?", merchantId).Find(&goods)
	if goods.Id == 0 {
		con.Error(c, "商品不存在或不属于当前商户", "/merchant/goods")
		return
	}

	goods.IsDelete = 1
	goods.Status = 0
	models.DB.Save(&goods)

	prevPage := c.Request.Referer()
	if len(prevPage) > 0 {
		con.Success(c, "删除数据成功", prevPage)
	} else {
		con.Success(c, "删除数据成功", "/merchant/goods")
	}
}

// 富文本编辑器上传图片
func (con GoodsController) EditorImageUpload(c *gin.Context) {
	imgDir, err := models.UploadImg(c, "file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
	} else {
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(imgDir)
				wg.Done()
			}()
			c.JSON(http.StatusOK, gin.H{
				"link": "/" + imgDir,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"link": models.GetSettingFromColumn("OssDomain") + imgDir,
			})
		}
	}
}

// 图库上传图片
func (con GoodsController) GoodsImageUpload(c *gin.Context) {
	imgDir, err := models.UploadImg(c, "file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"link": "",
		})
	} else {
		if models.GetOssStatus() != 1 {
			wg.Add(1)
			go func() {
				models.ResizeGoodsImage(imgDir)
				wg.Done()
			}()
		}
		c.JSON(http.StatusOK, gin.H{
			"link": imgDir,
		})
	}
}

func (con GoodsController) ChangeGoodsImageColor(c *gin.Context) {
	goodsImageId, err1 := models.Int(c.Query("goods_image_id"))
	colorId, err2 := models.Int(c.Query("color_id"))
	goodsImage := models.GoodsImage{Id: goodsImageId}
	models.DB.Find(&goodsImage)
	goodsImage.ColorId = colorId
	err3 := models.DB.Save(&goodsImage).Error
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  "更新失败",
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result":  "更新成功",
			"success": true,
		})
	}
}

func (con GoodsController) RemoveGoodsImage(c *gin.Context) {
	goodsImageId, err1 := models.Int(c.Query("goods_image_id"))
	goodsImage := models.GoodsImage{Id: goodsImageId}
	err2 := models.DB.Delete(&goodsImage).Error
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"result":  "删除失败",
			"success": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result":  "删除成功",
			"success": true,
		})
	}
}
