package controllers

import (
	"HongXunServer/middleware"
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"log"
)

type SessionController struct {
	Ctx     iris.Context
	Service services.UserService
}

func (c *SessionController) Post() {
	var authentication models.Authentication
	err := c.Ctx.ReadJSON(&authentication)
	if err != nil {
		log.Println(err)
	}
	_, err = c.Ctx.JSON(c.Service.Verify(&authentication))
	if err != nil {
		log.Println(err)
	}
}

func (c *SessionController) Get() {
	token := c.Ctx.GetHeader("Authorization")
	log.Println(token)
	var claims models.UserClaims
	err := middleware.J.VerifyTokenString(c.Ctx, token[6:], &claims)
	if err != nil {
		log.Println(err)
	}
	log.Println("user from:", claims.UserId)
	_, err = c.Ctx.JSON(c.Service.AutoLogin(claims.UserId))
	if err != nil {
		log.Println(err)
	}
}
