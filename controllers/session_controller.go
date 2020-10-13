package controllers

import (
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
	response := c.Service.Verify(&authentication)
	_, err = c.Ctx.JSON(response)
}
