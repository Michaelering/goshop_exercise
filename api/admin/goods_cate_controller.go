package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoodsCateController struct{}

func (con GoodsCateController) Index(c *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)
	common.Success(c, goodsCateList)
}

func (con GoodsCateController) Tree(c *gin.Context) {
	goodsCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Preload("GoodsCateItems").Find(&goodsCateList)
	common.Success(c, goodsCateList)
}

func (con GoodsCateController) Create(c *gin.Context) {
	title := strings.Trim(c.PostForm("title"), " ")
	cateImg := c.PostForm("cate_img")
	link := c.PostForm("link")
	template := c.PostForm("template")
	pid, _ := models.Int(c.PostForm("pid"))
	subTitle := c.PostForm("sub_title")
	keywords := c.PostForm("keywords")
	description := c.PostForm("description")
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	if title == "" {
		common.BadRequest(c, "分类标题不能为空")
		return
	}

	goodsCate := models.GoodsCate{
		Title:       title,
		CateImg:     cateImg,
		Link:        link,
		Template:    template,
		Pid:         pid,
		SubTitle:    subTitle,
		Keywords:    keywords,
		Description: description,
		Sort:        sort,
		Status:      status,
		AddTime:     int(models.GetUnix()),
	}
	err := models.DB.Create(&goodsCate).Error
	if err != nil {
		common.Error(c, 500, "增加分类失败")
		return
	}
	common.Success(c, goodsCate)
}

func (con GoodsCateController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	if goodsCate.Id == 0 {
		common.BadRequest(c, "分类不存在")
		return
	}

	goodsCate.Title = strings.Trim(c.PostForm("title"), " ")
	goodsCate.CateImg = c.PostForm("cate_img")
	goodsCate.Link = c.PostForm("link")
	goodsCate.Template = c.PostForm("template")
	goodsCate.Pid, _ = models.Int(c.PostForm("pid"))
	goodsCate.SubTitle = c.PostForm("sub_title")
	goodsCate.Keywords = c.PostForm("keywords")
	goodsCate.Description = c.PostForm("description")
	goodsCate.Sort, _ = models.Int(c.PostForm("sort"))
	goodsCate.Status, _ = models.Int(c.PostForm("status"))

	err = models.DB.Save(&goodsCate).Error
	if err != nil {
		common.Error(c, 500, "修改分类失败")
		return
	}
	common.Success(c, goodsCate)
}

func (con GoodsCateController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	// 检查是否有子分类
	var childCount int64
	models.DB.Where("pid=?", id).Table("goods_cate").Count(&childCount)
	if childCount > 0 {
		common.BadRequest(c, "该分类下有子分类，请先删除子分类")
		return
	}

	models.DB.Delete(&models.GoodsCate{Id: id})
	common.Success(c, nil)
}

func (con GoodsCateController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	goodsCate := models.GoodsCate{Id: id}
	models.DB.Find(&goodsCate)
	if goodsCate.Id == 0 {
		common.BadRequest(c, "分类不存在")
		return
	}

	// 获取所有顶级分类供下拉选择
	topCateList := []models.GoodsCate{}
	models.DB.Where("pid=0").Find(&topCateList)

	common.Success(c, gin.H{
		"goodsCate":    goodsCate,
		"topCateList": topCateList,
	})
}
