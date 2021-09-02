package api

import (
	"errors"
	"fmt"
	"leekbox/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenBody struct {
	UserInfo model.User `json:"userinfo"`
	jwt.StandardClaims
}

var secret = []byte("leekbox")

func GenToken(data interface{}) (string, error) {
	const expiresDuration = time.Hour * 24
	conf := TokenBody{
		UserInfo: data.(model.User),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresDuration).Unix(),
			Issuer:    "leekbox",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, conf)
	return token.SignedString(secret)
}

func ParseToken(token string) (*TokenBody, error) {
	tokenBody, err := jwt.ParseWithClaims(token, &TokenBody{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if temp, ok := tokenBody.Claims.(*TokenBody); ok && tokenBody.Valid {
		return temp, nil
	}
	return nil, errors.New("invaild token")
}

func AttachToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Set("authed", false)
		} else {
			temp, err := ParseToken(authHeader)
			if err != nil {
				c.Set("authed", false)
			} else {
				c.Set("authed", true)
				c.Set("userInfo", temp.UserInfo)
				fmt.Println("authed")
			}
		}
		c.Next()
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authed, exist := c.Get("authed"); authed.(bool) && exist {
			c.Next()
			return
		}
		fmt.Println("needAuth")
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Abort()
			c.JSON(http.StatusForbidden, model.Return(40100, nil, "token为空"))
			return
		}
		temp, err := ParseToken(authHeader)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusForbidden, model.Return(40300, nil, err.Error()))
			return
		}
		c.Set("userInfo", temp.UserInfo)
		c.Next()
	}
}
