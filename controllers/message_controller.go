package controllers

import (
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"log"
)

type MessageController struct {
	Ctx     iris.Context
	Service services.MessageService
}

func (c *MessageController) Get() {
	m, err := c.Service.LoadMessage(c.Ctx.URLParam("user_to"), c.Ctx.URLParam("user_from"))
	if err != nil {
		c.Ctx.NotFound()
		return
	}
	_, err = c.Ctx.JSON(m)

	if err != nil {
		log.Println(err)
	}
}

func (c *MessageController) Post() {
	var message models.Message
	err := c.Ctx.ReadJSON(&message)
	if err != nil {
		log.Println(err)
		c.Ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem())
		return
	}
	err = c.Service.SaveMessage(message)
	if err != nil {
		log.Println(err)
		c.Ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem())
		return
	}
	c.Ctx.JSON(iris.Map{"is_success": true})
}
