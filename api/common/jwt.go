package common

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/ini.v1"
)

// JWT 自定义 claims
type AdminClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	RoleId   int    `json:"role_id"`
	IsSuper  int    `json:"is_super"`
	Type     string `json:"type"` // "admin"
	jwt.RegisteredClaims
}

type MerchantClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	ShopName string `json:"shop_name"`
	Type     string `json:"type"` // "merchant"
	jwt.RegisteredClaims
}

// 从配置文件加载 JWT 密钥
func getJwtSecret() []byte {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}
	secret := config.Section("jwt").Key("secret").String()
	if secret == "" {
		secret = "goshop-jwt-secret-key"
	}
	return []byte(secret)
}

func getJwtExpireHours() int {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}
	hours, _ := config.Section("jwt").Key("expire_hours").Int()
	if hours <= 0 {
		hours = 24
	}
	return hours
}

// 生成 Admin JWT
func GenerateAdminToken(userId int, username string, roleId int, isSuper int) (string, error) {
	claims := AdminClaims{
		UserId:   userId,
		Username: username,
		RoleId:   roleId,
		IsSuper:  isSuper,
		Type:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(getJwtExpireHours()) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtSecret())
}

// 生成 Merchant JWT
func GenerateMerchantToken(userId int, username string, shopName string) (string, error) {
	claims := MerchantClaims{
		UserId:   userId,
		Username: username,
		ShopName: shopName,
		Type:     "merchant",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(getJwtExpireHours()) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtSecret())
}

// 解析 Admin Token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		if claims.Type != "admin" {
			return nil, errors.New("token 类型错误")
		}
		return claims, nil
	}
	return nil, errors.New("无效的 token")
}

// 解析 Merchant Token
func ParseMerchantToken(tokenString string) (*MerchantClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MerchantClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MerchantClaims); ok && token.Valid {
		if claims.Type != "merchant" {
			return nil, errors.New("token 类型错误")
		}
		return claims, nil
	}
	return nil, errors.New("无效的 token")
}
