package admin

import (
	"fmt"
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type FocusController struct{}

func (con FocusController) Index(c *gin.Context) {
	focusList := []models.Focus{}
	models.DB.Find(&focusList)
	common.Success(c, focusList)
}

func (con FocusController) Create(c *gin.Context) {
	title := c.PostForm("title")
	focusType, _ := models.Int(c.PostForm("focus_type"))
	link := c.PostForm("link")
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	focusImgSrc, err4 := models.UploadImg(c, "focus_img")
	if err4 != nil {
		fmt.Println(err4)
	}

	focus := models.Focus{
		Title:     title,
		FocusType: focusType,
		FocusImg:  focusImgSrc,
		Link:      link,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}
	err5 := models.DB.Create(&focus).Error
	if err5 != nil {
		common.Error(c, 500, "增加轮播图失败")
		return
	}
	common.Success(c, focus)
}

func (con FocusController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	focus := models.Focus{Id: id}
	models.DB.Find(&focus)
	if focus.Id == 0 {
		common.BadRequest(c, "轮播图不存在")
		return
	}

	focus.Title = c.PostForm("title")
	focus.FocusType, _ = models.Int(c.PostForm("focus_type"))
	focus.Link = c.PostForm("link")
	focus.Sort, _ = models.Int(c.PostForm("sort"))
	focus.Status, _ = models.Int(c.PostForm("status"))

	focusImg, _ := models.UploadImg(c, "focus_img")
	if focusImg != "" {
		focus.FocusImg = focusImg
	}

	err = models.DB.Save(&focus).Error
	if err != nil {
		common.Error(c, 500, "修改轮播图失败")
		return
	}
	common.Success(c, focus)
}

func (con FocusController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.Focus{Id: id})
	common.Success(c, nil)
}

func (con FocusController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	focus := models.Focus{Id: id}
	models.DB.Find(&focus)
	if focus.Id == 0 {
		common.BadRequest(c, "轮播图不存在")
		return
	}
	common.Success(c, focus)
}
