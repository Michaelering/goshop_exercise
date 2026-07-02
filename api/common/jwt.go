package common

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/ini.v1"
)

var (
	jwtSecret     []byte
	jwtExpireHours int
	jwtInitOnce   sync.Once
)

// initJwtConfig 一次性加载 JWT 配置到内存，避免每次请求读磁盘
func initJwtConfig() {
	jwtInitOnce.Do(func() {
		config, iniErr := ini.Load("./conf/app.ini")
		if iniErr != nil {
			fmt.Printf("Fail to read file: %v", iniErr)
			os.Exit(1)
		}
		secret := config.Section("jwt").Key("secret").String()
		if secret == "" {
			secret = "goshop-jwt-secret-key"
		}
		jwtSecret = []byte(secret)

		hours, _ := config.Section("jwt").Key("expire_hours").Int()
		if hours <= 0 {
			hours = 24
		}
		jwtExpireHours = hours
	})
}

// AdminClaims 管理员 JWT claims
type AdminClaims struct {
	UserId    int    `json:"user_id"`
	Username  string `json:"username"`
	RoleId    int    `json:"role_id"`
	RoleTitle string `json:"role_title"` // "超级管理员" | "管理员" | 自定义角色名
	Type      string `json:"type"`       // "admin"
	jwt.RegisteredClaims
}

// MerchantClaims 商户 JWT claims
type MerchantClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	ShopName string `json:"shop_name"`
	Type     string `json:"type"` // "merchant"
	jwt.RegisteredClaims
}

// GenerateAdminToken 生成管理员 JWT
func GenerateAdminToken(userId int, username string, roleId int, roleTitle string) (string, error) {
	initJwtConfig()
	claims := AdminClaims{
		UserId:    userId,
		Username:  username,
		RoleId:    roleId,
		RoleTitle: roleTitle,
		Type:      "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateMerchantToken 生成商户 JWT
func GenerateMerchantToken(userId int, username string, shopName string) (string, error) {
	initJwtConfig()
	claims := MerchantClaims{
		UserId:   userId,
		Username: username,
		ShopName: shopName,
		Type:     "merchant",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseAdminToken 解析管理员 Token
func ParseAdminToken(tokenString string) (*AdminClaims, error) {
	initJwtConfig()
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
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

// ParseMerchantToken 解析商户 Token
func ParseMerchantToken(tokenString string) (*MerchantClaims, error) {
	initJwtConfig()
	token, err := jwt.ParseWithClaims(tokenString, &MerchantClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
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
