package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type NavController struct{}

func (con NavController) Index(c *gin.Context) {
	page, _ := models.Int(c.DefaultQuery("page", "1"))
	pageSize, _ := models.Int(c.DefaultQuery("pageSize", "10"))

	navList := []models.Nav{}
	var count int64
	models.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&navList)
	models.DB.Table("nav").Count(&count)

	common.List(c, navList, count, page, pageSize)
}

func (con NavController) Create(c *gin.Context) {
	title := c.PostForm("title")
	link := c.PostForm("link")
	position, _ := models.Int(c.PostForm("position"))
	isOpennew, _ := models.Int(c.PostForm("is_opennew"))
	relation := c.PostForm("relation")
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	if title == "" {
		common.BadRequest(c, "标题不能为空")
		return
	}

	nav := models.Nav{
		Title:     title,
		Link:      link,
		Position:  position,
		IsOpennew: isOpennew,
		Relation:  relation,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}
	err := models.DB.Create(&nav).Error
	if err != nil {
		common.Error(c, 500, "增加导航失败")
		return
	}
	common.Success(c, nav)
}

func (con NavController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	nav := models.Nav{Id: id}
	models.DB.Find(&nav)
	if nav.Id == 0 {
		common.BadRequest(c, "导航不存在")
		return
	}

	nav.Title = c.PostForm("title")
	nav.Link = c.PostForm("link")
	nav.Position, _ = models.Int(c.PostForm("position"))
	nav.IsOpennew, _ = models.Int(c.PostForm("is_opennew"))
	nav.Relation = c.PostForm("relation")
	nav.Sort, _ = models.Int(c.PostForm("sort"))
	nav.Status, _ = models.Int(c.PostForm("status"))

	err = models.DB.Save(&nav).Error
	if err != nil {
		common.Error(c, 500, "修改导航失败")
		return
	}
	common.Success(c, nav)
}

func (con NavController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.Nav{Id: id})
	common.Success(c, nil)
}

func (con NavController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	nav := models.Nav{Id: id}
	models.DB.Find(&nav)
	if nav.Id == 0 {
		common.BadRequest(c, "导航不存在")
		return
	}
	common.Success(c, nav)
}
