package controllers

import (
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"log"
)

type UserController struct {
	Ctx     iris.Context
	Service services.UserService
}

func (c *UserController) Post() {
	var user models.User
	err := c.Ctx.ReadJSON(&user)
	if err != nil {
		log.Println(err)
	}
	_, err = c.Ctx.JSON(c.Service.Register(user))
}
