package controllers

import (
	"HongXunServer/models"
	"HongXunServer/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"log"
	"strconv"
)

type MessageController struct {
	Ctx     iris.Context
	Service services.MessageService
}

func (c *MessageController) GetMessage() {
	m, err := c.Service.LoadMessage()
	if err != nil {
		c.Ctx.NotFound()
		return
	}
	_, err = c.Ctx.JSON(m)

	if err != nil {
		log.Println(err)
	}
}

func (c *MessageController) PostMessage() mvc.Result {
	var (
		id, _         = strconv.Atoi(c.Ctx.FormValue("id"))
		createTime, _ = strconv.Atoi(c.Ctx.FormValue("createTime"))
		userFrom, _   = strconv.Atoi(c.Ctx.FormValue("userFrom"))
		userTo, _     = strconv.Atoi(c.Ctx.FormValue("userTo"))
		content       = c.Ctx.FormValue("content")
		isViewed, _   = strconv.ParseBool(c.Ctx.FormValue("isViewed"))
	)

	err := c.Service.SaveMessage(models.Message{
		Id:         id,
		CreateTime: createTime,
		UserFrom:   userFrom,
		UserTo:     userTo,
		Content:    content,
		IsViewed:   isViewed,
	})
	return mvc.Response{
		// if not nil then this error will be shown instead.
		Err: err,
		// redirect to /user/me.
		Path: "/chat/message",
		// When redirecting from POST to GET request you -should- use this HTTP status code,
		// however there're some (complicated) alternatives if you
		// search online or even the HTTP RFC.
		// Status "See Other" RFC 7231, however iris can automatically fix that
		// but it's good to know you can set a custom code;
		// Code: 303,
	}
}
