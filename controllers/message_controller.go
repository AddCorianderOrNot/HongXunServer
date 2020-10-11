package controllers

import (
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type MessageController struct {
	Ctx     iris.Context
	Service services.MessageService
}

func (c *MessageController) Get() {
	var claims models.UserClaims
	err := jwt.ReadClaims(c.Ctx, &claims)
	userFrom, _ := primitive.ObjectIDFromHex(c.Ctx.URLParam("user_from"))
	m, err := c.Service.LoadMessage(claims.UserId, userFrom)
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
	var claims models.UserClaims
	var message models.Message
	err := c.Ctx.ReadJSON(&message)
	if err != nil {
		log.Println(err)
	}
	err = jwt.ReadClaims(c.Ctx, &claims)
	if err != nil {
		log.Println(err)
	}
	log.Println("user from:", claims.UserId)
	message.UserFrom = claims.UserId
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
