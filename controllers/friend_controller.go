package controllers

import (
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"log"
)

type FriendController struct {
	Ctx     iris.Context
	Service services.FriendService
}

func (c *FriendController) Get() {
	var claims models.UserClaims
	err := jwt.ReadClaims(c.Ctx, &claims)
	if err != nil {
		log.Println(err)
	}
	_, err = c.Ctx.JSON(c.Service.LoadFriend(claims.UserId))
	if err != nil {
		log.Println(err)
	}
}

func (c *FriendController) Post() {
	var claims models.UserClaims
	err := jwt.ReadClaims(c.Ctx, &claims)
	if err != nil {
		log.Println(err)
	}
	log.Println("user from:", claims.UserId)

	var friend models.Friend
	err = c.Ctx.ReadJSON(&friend)
	if err != nil {
		log.Println(err)
	}

	_, err = c.Ctx.JSON(c.Service.AddFriend(claims.UserId, friend.FriendId))
	if err != nil {
		log.Println(err)
	}
}
