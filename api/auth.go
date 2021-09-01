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
	const expiresDuration = time.Hour * 1
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

func AuthMiddleWare() func(*gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			resp := model.Resp{
				Code:    40100,
				Data:    nil,
				Message: "token为空",
			}
			c.Abort()
			c.JSON(http.StatusForbidden, resp)
			return
		}
		temp, err := ParseToken(authHeader)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusForbidden,
				model.Resp{
					Code:    40300,
					Data:    nil,
					Message: err.Error(),
				})
			return
		}
		c.Set("tokenBody", temp)
		c.Set("userInfo", temp.UserInfo)
		c.Next()
	}
}
