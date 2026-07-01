package admin

import (
	"ginshop58/api/common"
	"ginshop58/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type MerchantMgmtController struct{}

func (con MerchantMgmtController) Index(c *gin.Context) {
	merchantList := []models.Merchant{}
	models.DB.Find(&merchantList)
	common.Success(c, merchantList)
}

func (con MerchantMgmtController) Create(c *gin.Context) {
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	shopName := strings.Trim(c.PostForm("shop_name"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	if len(username) < 2 || len(password) < 6 {
		common.BadRequest(c, "用户名或者密码的长度不合法")
		return
	}
	if shopName == "" {
		common.BadRequest(c, "店铺名称不能为空")
		return
	}

	var count int64
	models.DB.Where("username=?", username).Table("merchant").Count(&count)
	if count > 0 {
		common.BadRequest(c, "此商户已存在")
		return
	}

	merchant := models.Merchant{
		Username: username,
		Password: models.Md5(password),
		ShopName: shopName,
		Email:    email,
		Mobile:   mobile,
		Status:   1,
		AddTime:  int(models.GetUnix()),
	}
	err := models.DB.Create(&merchant).Error
	if err != nil {
		common.Error(c, 500, "增加商户失败")
		return
	}
	merchant.Password = ""
	common.Success(c, merchant)
}

func (con MerchantMgmtController) Update(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	merchant := models.Merchant{Id: id}
	models.DB.Find(&merchant)
	if merchant.Id == 0 {
		common.BadRequest(c, "商户不存在")
		return
	}

	merchant.Username = strings.Trim(c.PostForm("username"), " ")
	merchant.ShopName = strings.Trim(c.PostForm("shop_name"), " ")
	merchant.Email = c.PostForm("email")
	merchant.Mobile = c.PostForm("mobile")
	merchant.Status, _ = models.Int(c.PostForm("status"))

	password := strings.Trim(c.PostForm("password"), " ")
	if password != "" {
		if len(password) < 6 {
			common.BadRequest(c, "密码长度不能小于6位")
			return
		}
		merchant.Password = models.Md5(password)
	}

	err = models.DB.Save(&merchant).Error
	if err != nil {
		common.Error(c, 500, "修改商户失败")
		return
	}
	merchant.Password = ""
	common.Success(c, merchant)
}

func (con MerchantMgmtController) Delete(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}
	models.DB.Delete(&models.Merchant{Id: id})
	common.Success(c, nil)
}

func (con MerchantMgmtController) Get(c *gin.Context) {
	id, err := models.Int(c.Param("id"))
	if err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	merchant := models.Merchant{Id: id}
	models.DB.Find(&merchant)
	if merchant.Id == 0 {
		common.BadRequest(c, "商户不存在")
		return
	}
	merchant.Password = ""
	common.Success(c, merchant)
}
