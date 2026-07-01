package models

type Setting struct {
	Id              int    `form:"id" json:"id"`
	SiteTitle       string `form:"site_title" json:"site_title"`
	SiteLogo        string `json:"site_logo"`
	SiteKeywords    string `form:"site_keywords" json:"site_keywords"`
	SiteDescription string `form:"site_description" json:"site_description"`
	NoPicture       string `json:"no_picture"`
	SiteIcp         string `form:"site_icp" json:"site_icp"`
	SiteTel         string `form:"site_tel" json:"site_tel"`
	SearchKeywords  string `form:"search_keywords" json:"search_keywords"`
	TongjiCode      string `form:"tongji_code" json:"tongji_code"`
	Appid           string `form:"appid" json:"appid"`
	AppSecret       string `form:"app_secret" json:"app_secret"`
	EndPoint        string `form:"end_point" json:"end_point"`
	BucketName      string `form:"bucket_name" json:"bucket_name"`
	OssStatus       int    `form:"oss_status" json:"oss_status"`
	OssDomain       string `form:"oss_domain" json:"oss_domain"`
	ThumbnailSize   string `form:"thumbnail_size" json:"thumbnail_size"`
}

func (Setting) TableName() string {
	return "setting"
}
