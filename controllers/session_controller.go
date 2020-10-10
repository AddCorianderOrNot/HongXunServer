package controllers

import (
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"log"
	"time"
)

type SessionController struct {
	Ctx     iris.Context
	Service services.UserService
}

type UserClaims struct {
	jwt.Claims
	email string
}

func (c *SessionController) Post() {
	var authentication models.Authentication
	err := c.Ctx.ReadJSON(&authentication)
	if err != nil {
		log.Println(err)
	}
	isVerified := c.Service.Verify(authentication)
	if isVerified {
		j := jwt.HMAC(15*time.Minute, "secret", "itsa16bytesecret")
		claims := UserClaims{
			Claims: j.Expiry(jwt.Claims{
				Issuer:   "an-issuer",
				Audience: jwt.Audience{"an-audience"},
			}),
			email: authentication.Email,
		}
		accessToken, _ := j.Token(claims)
		_, err = c.Ctx.JSON(iris.Map{"is_success": true, "token": accessToken})
	} else {
		_, err = c.Ctx.JSON(iris.Map{"is_success": false, "token": nil})
	}
}
