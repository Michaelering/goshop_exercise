package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoodsTypeAttrController struct{}

func (con GoodsTypeAttrController) Index(c *gin.Context) {
	cateId, _ := models.Int(c.Query("cateId"))
	goodsTypeAttrList := []models.GoodsTypeAttribute{}

	if cateId > 0 {
		models.DB.Where("cate_id=?", cateId).Find(&goodsTypeAttrList)
	} else {
		models.DB.Find(&goodsTypeAttrList)
	}
	common.Success(c, goodsTypeAttrList)
}

func (con GoodsTypeAttrController) ByType(c *gin.Context) {
	cateId, err := models.Int(c.Param("cateId"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsTypeAttrList := []models.GoodsTypeAttribute{}
	models.DB.Where("cate_id=?", cateId).Find(&goodsTypeAttrList)
	common.Success(c, goodsTypeAttrList)
}

func (con GoodsTypeAttrController) Create(c *gin.Context) {
	cateId, _ := models.Int(c.PostForm("cate_id"))
	title := strings.Trim(c.PostForm("title"), " ")
	attrType, _ := models.Int(c.PostForm("attr_type"))
	attrValue := c.PostForm("attr_value")
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	if title == "" || cateId == 0 {
		common.BadRequest(c, "参数不完整")
		return
	}

	attr := models.GoodsTypeAttribute{
		CateId:    cateId,
		Title:     title,
		AttrType:  attrType,
		AttrValue: attrValue,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}
	err := models.DB.Create(&attr).Error
	if err != nil {
		common.Error(c, 500, "增加属性失败")
		return
	}
	common.Success(c, attr)
}

func (con GoodsTypeAttrController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	attr := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&attr)
	if attr.Id == 0 {
		common.BadRequest(c, "属性不存在")
		return
	}

	attr.CateId, _ = models.Int(c.PostForm("cate_id"))
	attr.Title = strings.Trim(c.PostForm("title"), " ")
	attr.AttrType, _ = models.Int(c.PostForm("attr_type"))
	attr.AttrValue = c.PostForm("attr_value")
	attr.Sort, _ = models.Int(c.PostForm("sort"))
	attr.Status, _ = models.Int(c.PostForm("status"))

	err = models.DB.Save(&attr).Error
	if err != nil {
		common.Error(c, 500, "修改属性失败")
		return
	}
	common.Success(c, attr)
}

func (con GoodsTypeAttrController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.GoodsTypeAttribute{Id: id})
	common.Success(c, nil)
}

func (con GoodsTypeAttrController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	attr := models.GoodsTypeAttribute{Id: id}
	models.DB.Find(&attr)
	if attr.Id == 0 {
		common.BadRequest(c, "属性不存在")
		return
	}

	// 获取所有商品类型供下拉选择
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)

	common.Success(c, gin.H{
		"attr":          attr,
		"goodsTypeList": goodsTypeList,
	})
}
