package admin

import (
	"ginshop58/models"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// 商品创建/更新的通用参数
type GoodsCreateParams struct {
	MerchantId  int
	Title       string
	SubTitle    string
	GoodsSn     string
	CateId      int
	GoodsNumber int
	MarketPrice float64
	Price       float64
	RelationGoods string
	GoodsAttr   string
	GoodsVersion string
	GoodsGift   string
	GoodsFitting string
	GoodsColor  []string
	GoodsKeywords string
	GoodsDesc   string
	GoodsContent string
	IsHot       int
	IsBest      int
	IsNew       int
	GoodsTypeId int
	Sort        int
	Status      int
	GoodsImg    string
	// 图片列表 和 属性列表（在 controller 中处理 db 写入）
	GoodsImageList []string
	AttrIdList     []string
	AttrValueList  []string
}

// 创建商品
func CreateGoods(params *GoodsCreateParams) (*models.Goods, error) {
	goodsColorStr := strings.Join(params.GoodsColor, ",")

	goods := models.Goods{
		Title:         params.Title,
		SubTitle:      params.SubTitle,
		GoodsSn:       params.GoodsSn,
		CateId:        params.CateId,
		ClickCount:    100,
		GoodsNumber:   params.GoodsNumber,
		MarketPrice:   params.MarketPrice,
		Price:         params.Price,
		RelationGoods: params.RelationGoods,
		GoodsAttr:     params.GoodsAttr,
		GoodsVersion:  params.GoodsVersion,
		GoodsGift:     params.GoodsGift,
		GoodsFitting:  params.GoodsFitting,
		GoodsKeywords: params.GoodsKeywords,
		GoodsDesc:     params.GoodsDesc,
		GoodsContent:  params.GoodsContent,
		IsDelete:      0,
		IsHot:         params.IsHot,
		IsBest:        params.IsBest,
		IsNew:         params.IsNew,
		GoodsTypeId:   params.GoodsTypeId,
		Sort:          params.Sort,
		Status:        params.Status,
		AddTime:       int(models.GetUnix()),
		GoodsColor:    goodsColorStr,
		GoodsImg:      params.GoodsImg,
		MerchantId:    params.MerchantId,
	}

	err := models.DB.Create(&goods).Error
	if err != nil {
		return nil, err
	}

	// 增加图库
	for _, v := range params.GoodsImageList {
		goodsImgObj := models.GoodsImage{}
		goodsImgObj.GoodsId = goods.Id
		goodsImgObj.ImgUrl = v
		goodsImgObj.Sort = 10
		goodsImgObj.Status = 1
		goodsImgObj.AddTime = int(models.GetUnix())
		models.DB.Create(&goodsImgObj)
	}

	// 增加规格包装
	for i := 0; i < len(params.AttrIdList); i++ {
		goodsTypeAttributeId, attributeIdErr := models.Int(params.AttrIdList[i])
		if attributeIdErr == nil && i < len(params.AttrValueList) {
			goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
			models.DB.Find(&goodsTypeAttributeObj)
			if goodsTypeAttributeObj.Id > 0 {
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
				goodsAttrObj.AttributeValue = params.AttrValueList[i]
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}
		}
	}

	return &goods, nil
}

// 更新商品
func UpdateGoods(id int, params *GoodsCreateParams) error {
	goods := models.Goods{Id: id}
	models.DB.Find(&goods)
	if goods.Id == 0 {
		return nil // 商品不存在
	}

	goodsColorStr := strings.Join(params.GoodsColor, ",")
	goods.Title = params.Title
	goods.SubTitle = params.SubTitle
	goods.GoodsSn = params.GoodsSn
	goods.CateId = params.CateId
	goods.GoodsNumber = params.GoodsNumber
	goods.MarketPrice = params.MarketPrice
	goods.Price = params.Price
	goods.RelationGoods = params.RelationGoods
	goods.GoodsAttr = params.GoodsAttr
	goods.GoodsVersion = params.GoodsVersion
	goods.GoodsGift = params.GoodsGift
	goods.GoodsFitting = params.GoodsFitting
	goods.GoodsKeywords = params.GoodsKeywords
	goods.GoodsDesc = params.GoodsDesc
	goods.GoodsContent = params.GoodsContent
	goods.IsHot = params.IsHot
	goods.IsBest = params.IsBest
	goods.IsNew = params.IsNew
	goods.GoodsTypeId = params.GoodsTypeId
	goods.Sort = params.Sort
	goods.Status = params.Status
	goods.GoodsColor = goodsColorStr
	if params.GoodsImg != "" {
		goods.GoodsImg = params.GoodsImg
	}

	err := models.DB.Save(&goods).Error
	if err != nil {
		return err
	}

	// 重新处理图库
	for _, v := range params.GoodsImageList {
		goodsImgObj := models.GoodsImage{}
		goodsImgObj.GoodsId = goods.Id
		goodsImgObj.ImgUrl = v
		goodsImgObj.Sort = 10
		goodsImgObj.Status = 1
		goodsImgObj.AddTime = int(models.GetUnix())
		models.DB.Create(&goodsImgObj)
	}

	// 删除旧的规格包装，重新添加
	models.DB.Where("goods_id=?", goods.Id).Delete(&models.GoodsAttr{})
	for i := 0; i < len(params.AttrIdList); i++ {
		goodsTypeAttributeId, attributeIdErr := models.Int(params.AttrIdList[i])
		if attributeIdErr == nil && i < len(params.AttrValueList) {
			goodsTypeAttributeObj := models.GoodsTypeAttribute{Id: goodsTypeAttributeId}
			models.DB.Find(&goodsTypeAttributeObj)
			if goodsTypeAttributeObj.Id > 0 {
				goodsAttrObj := models.GoodsAttr{}
				goodsAttrObj.GoodsId = goods.Id
				goodsAttrObj.AttributeTitle = goodsTypeAttributeObj.Title
				goodsAttrObj.AttributeType = goodsTypeAttributeObj.AttrType
				goodsAttrObj.AttributeId = goodsTypeAttributeObj.Id
				goodsAttrObj.AttributeCateId = goodsTypeAttributeObj.CateId
				goodsAttrObj.AttributeValue = params.AttrValueList[i]
				goodsAttrObj.Status = 1
				goodsAttrObj.Sort = 10
				goodsAttrObj.AddTime = int(models.GetUnix())
				models.DB.Create(&goodsAttrObj)
			}
		}
	}

	return nil
}

// 构建商品属性HTML（用于编辑页面）
func BuildGoodsAttrHtml(goodsId int) string {
	goodsAttr := []models.GoodsAttr{}
	models.DB.Where("goods_id=?", goodsId).Find(&goodsAttr)
	goodsAttrStr := ""
	for _, v := range goodsAttr {
		if v.AttributeType == 1 {
			goodsAttrStr += `<li><span>` + v.AttributeTitle + `: </span> <input type="hidden" name="attr_id_list" value="` + models.String(v.AttributeId) + `" />   <input type="text" name="attr_value_list" value="` + v.AttributeValue + `" /></li>`
		} else if v.AttributeType == 2 {
			goodsAttrStr += `<li><span>` + v.AttributeTitle + `: 　</span><input type="hidden" name="attr_id_list" value="` + models.String(v.AttributeId) + `" />  <textarea cols="50" rows="3" name="attr_value_list">` + v.AttributeValue + `</textarea></li>`
		} else {
			goodsTypeArttribute := models.GoodsTypeAttribute{Id: v.AttributeId}
			models.DB.Find(&goodsTypeArttribute)
			attrValueSlice := strings.Split(goodsTypeArttribute.AttrValue, "\n")
			goodsAttrStr += `<li><span>` + v.AttributeTitle + `: 　</span>  <input type="hidden" name="attr_id_list" value="` + models.String(v.AttributeId) + `" /> `
			goodsAttrStr += `<select name="attr_value_list">`
			for _, attrVal := range attrValueSlice {
				selected := ""
				if attrVal == v.AttributeValue {
					selected = ` selected`
				}
				goodsAttrStr += `<option value="` + attrVal + `"` + selected + `>` + attrVal + `</option>`
			}
			goodsAttrStr += `</select></li>`
		}
	}
	return goodsAttrStr
}
