package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoodsTypeController struct{}

func (con GoodsTypeController) Index(c *gin.Context) {
	goodsTypeList := []models.GoodsType{}
	models.DB.Find(&goodsTypeList)
	common.Success(c, goodsTypeList)
}

func (con GoodsTypeController) Create(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), " ")
	description := c.PostForm("description")

	if title == "" {
		common.BadRequest(c, "类型名称不能为空")
		return
	}

	goodsType := models.GoodsType{
		Title:       title,
		Description: description,
		Status:      1,
		AddTime:     int(models.GetUnix()),
	}
	err := models.DB.Create(&goodsType).Error
	if err != nil {
		common.Error(c, 500, "增加类型失败")
		return
	}
	common.Success(c, goodsType)
}

func (con GoodsTypeController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsType := models.GoodsType{Id: id}
	models.DB.Find(&goodsType)
	if goodsType.Id == 0 {
		common.BadRequest(c, "类型不存在")
		return
	}

	goodsType.Title = strings.Trim(c.PostForm("title"), " ")
	goodsType.Description = c.PostForm("description")
	goodsType.Status, _ = models.Int(c.PostForm("status"))

	err = models.DB.Save(&goodsType).Error
	if err != nil {
		common.Error(c, 500, "修改类型失败")
		return
	}
	common.Success(c, goodsType)
}

func (con GoodsTypeController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.GoodsType{Id: id})
	common.Success(c, nil)
}

func (con GoodsTypeController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsType := models.GoodsType{Id: id}
	models.DB.Find(&goodsType)
	if goodsType.Id == 0 {
		common.BadRequest(c, "类型不存在")
		return
	}
	common.Success(c, goodsType)
}
