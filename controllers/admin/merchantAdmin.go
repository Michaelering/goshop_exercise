package admin

import (
	"ginshop58/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type MerchantAdminController struct {
	BaseController
}

func (con MerchantAdminController) Index(c *gin.Context) {
	merchantList := []models.Merchant{}
	models.DB.Find(&merchantList)
	c.HTML(http.StatusOK, "admin/merchant/index.html", gin.H{
		"merchantList": merchantList,
	})
}

func (con MerchantAdminController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/merchant/add.html", gin.H{})
}

func (con MerchantAdminController) DoAdd(c *gin.Context) {
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	shopName := strings.Trim(c.PostForm("shop_name"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	if len(username) < 2 || len(password) < 6 {
		con.Error(c, "用户名或者密码的长度不合法", "/admin/merchant/add")
		return
	}
	if shopName == "" {
		con.Error(c, "店铺名称不能为空", "/admin/merchant/add")
		return
	}

	merchantList := []models.Merchant{}
	models.DB.Where("username=?", username).Find(&merchantList)
	if len(merchantList) > 0 {
		con.Error(c, "此商户已存在", "/admin/merchant/add")
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
		con.Error(c, "增加商户失败", "/admin/merchant/add")
		return
	}
	con.Success(c, "增加商户成功", "/admin/merchant")
}

func (con MerchantAdminController) Edit(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/merchant")
		return
	}
	merchant := models.Merchant{Id: id}
	models.DB.Find(&merchant)
	c.HTML(http.StatusOK, "admin/merchant/edit.html", gin.H{
		"merchant": merchant,
	})
}

func (con MerchantAdminController) DoEdit(c *gin.Context) {
	id, err1 := models.Int(c.PostForm("id"))
	if err1 != nil {
		con.Error(c, "传入数据错误", "/admin/merchant")
		return
	}
	username := strings.Trim(c.PostForm("username"), " ")
	password := strings.Trim(c.PostForm("password"), " ")
	shopName := strings.Trim(c.PostForm("shop_name"), " ")
	email := strings.Trim(c.PostForm("email"), " ")
	mobile := strings.Trim(c.PostForm("mobile"), " ")

	if shopName == "" {
		con.Error(c, "店铺名称不能为空", "/admin/merchant/edit?id="+models.String(id))
		return
	}

	merchant := models.Merchant{Id: id}
	models.DB.Find(&merchant)
	merchant.Username = username
	merchant.ShopName = shopName
	merchant.Email = email
	merchant.Mobile = mobile

	if password != "" {
		if len(password) < 6 {
			con.Error(c, "密码的长度不合法 密码长度不能小于6位", "/admin/merchant/edit?id="+models.String(id))
			return
		}
		merchant.Password = models.Md5(password)
	}
	err2 := models.DB.Save(&merchant).Error
	if err2 != nil {
		con.Error(c, "修改数据失败", "/admin/merchant/edit?id="+models.String(id))
		return
	}
	con.Success(c, "修改数据成功", "/admin/merchant")
}

func (con MerchantAdminController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/merchant")
	} else {
		merchant := models.Merchant{Id: id}
		models.DB.Delete(&merchant)
		con.Success(c, "删除数据成功", "/admin/merchant")
	}
}
