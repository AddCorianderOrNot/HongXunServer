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

	var email models.FriendEmail
	err = c.Ctx.ReadJSON(&email)
	if err != nil {
		log.Println(err)
	}
	log.Println("friendEmail:", email.Email)
	_, err = c.Ctx.JSON(c.Service.AddFriendByEmail(claims.UserId, email.Email))
	if err != nil {
		log.Println(err)
	}
}

func (c *FriendController) GetTime() {
	var claims models.UserClaims
	err := jwt.ReadClaims(c.Ctx, &claims)
	if err != nil {
		log.Println(err)
	}

	var email models.FriendEmail
	err = c.Ctx.ReadQuery(&email)
	if err != nil {
		log.Println(err)
	}
	log.Println("GetTime", claims.UserId, email.Email)
	_, err = c.Ctx.JSON(c.Service.GetReadTime(claims.UserId, email.Email))
	if err != nil {
		log.Println(err)
	}
}

func (c *FriendController) PostTime() {
	var claims models.UserClaims
	err := jwt.ReadClaims(c.Ctx, &claims)
	if err != nil {
		log.Println(err)
	}
	var friend models.FriendTime
	err = c.Ctx.ReadJSON(&friend)
	if err != nil {
		log.Println(err)
	}
	log.Println("PostTime", claims.UserId, friend.Email)
	_, err = c.Ctx.JSON(c.Service.UpdateReadTime(claims.UserId, friend.Email, friend.ReadTime))
	if err != nil {
		log.Println(err)
	}
}