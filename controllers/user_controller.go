package controllers

import (
	"HongXunServer/middleware"
	"HongXunServer/models"
	"HongXunServer/services"
	"HongXunServer/utils"
	"github.com/kataras/iris/v12"
	"log"
	"path/filepath"
	"strconv"
	"time"
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
	_, err = c.Ctx.JSON(c.Service.Register(&user))
}

func (c *UserController) Get() {
	_, err := c.Ctx.JSON(c.Service.SearchUser(c.Ctx.URLParam("nickname")))
	if err != nil {
		log.Println(err)
	}
}

func (c *UserController) Patch() {
	token := c.Ctx.GetHeader("Authorization")
	log.Println(token)
	var claims models.UserClaims
	err := middleware.J.VerifyTokenString(c.Ctx, token[6:], &claims)
	if err != nil {
		log.Println(err)
		_, _ = c.Ctx.JSON(models.Response{
			ErrCode: 3,
			ErrMsg:  "没有权限",
			Data:    nil,
		})
	}
	log.Println("user from:", claims.UserId)
	key := c.Ctx.FormValue("key")
	log.Println("update:", key)
	if key == "nickname" || key == "signature" {
		value := c.Ctx.FormValue("value")
		_, err = c.Ctx.JSON(c.Service.UpdateUser(claims.UserId, key, value))
		if err != nil {
			log.Println(err)
		}
	} else if key == "image" {
		f, fh, err := c.Ctx.FormFile("value")
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		uploadDate := strconv.FormatInt(time.Now().Unix(), 10)
		saveName := utils.Md5(uploadDate + fh.Filename)
		_, err = c.Ctx.SaveFormFile(fh, filepath.Join("./uploads", saveName))
		if err != nil {
			log.Println(err)
		}
		_, err = c.Ctx.JSON(c.Service.UpdateUser(claims.UserId, "icon", "http://www.brotherye.site:8080/image/" + saveName))
	} else {
		_, _ = c.Ctx.JSON(models.Response{
			ErrCode: 2,
			ErrMsg:  "请不要更新不支持的字段",
			Data:    nil,
		})
	}
}