package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"

	"github.com/gin-gonic/gin"
)

type SettingController struct{}

func (con SettingController) Index(c *gin.Context) {
	setting := models.Setting{}
	models.DB.First(&setting)
	common.Success(c, setting)
}

func (con SettingController) Update(c *gin.Context) {
	setting := models.Setting{}
	models.DB.First(&setting)

	setting.SiteTitle = c.PostForm("site_title")
	setting.SiteKeywords = c.PostForm("site_keywords")
	setting.SiteDescription = c.PostForm("site_description")
	setting.SiteIcp = c.PostForm("site_icp")
	setting.SiteTel = c.PostForm("site_tel")
	setting.SearchKeywords = c.PostForm("search_keywords")
	setting.TongjiCode = c.PostForm("tongji_code")
	setting.Appid = c.PostForm("appid")
	setting.AppSecret = c.PostForm("app_secret")
	setting.EndPoint = c.PostForm("end_point")
	setting.BucketName = c.PostForm("bucket_name")
	setting.OssStatus, _ = models.Int(c.PostForm("oss_status"))
	setting.OssDomain = c.PostForm("oss_domain")
	setting.ThumbnailSize = c.PostForm("thumbnail_size")

	// Logo 和默认图上传
	logo, _ := models.UploadImg(c, "site_logo")
	if logo != "" {
		setting.SiteLogo = logo
	}
	noPic, _ := models.UploadImg(c, "no_picture")
	if noPic != "" {
		setting.NoPicture = noPic
	}

	err := models.DB.Save(&setting).Error
	if err != nil {
		common.Error(c, 500, "修改设置失败")
		return
	}
	common.Success(c, setting)
}
