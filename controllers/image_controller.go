package controllers

import (
	"github.com/kataras/iris/v12"
)

type ImageController struct {
	Ctx     iris.Context
}

func (c *ImageController) GetBy(imageName string) {
	src := "./uploads/" + imageName
	c.Ctx.SendFile(src, imageName)
}