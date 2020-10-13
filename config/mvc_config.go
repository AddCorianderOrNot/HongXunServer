package config

import (
	"HongXunServer/controllers"
	"HongXunServer/datasource"
	"HongXunServer/middleware"
	"HongXunServer/repositories"
	"HongXunServer/services"
	"github.com/kataras/iris/v12/mvc"
)

var (
	userCollection     = datasource.DB.Collection("users")
	userRepository     = repositories.NewUserRepository(userCollection)
	userService        = services.NewUserService(userRepository)
	friendCollection   = datasource.DB.Collection("friends")
	friendRepository   = repositories.NewFriendRepository(friendCollection)
	friendService      = services.NewFriendService(friendRepository, userRepository)
	messagesCollection = datasource.DB.Collection("messages")
	messageService     = services.NewMessageService(messagesCollection)
)

func ChatConfigure(app *mvc.Application) {
	app.Router.Use(middleware.Verify)
	app.Handle(new(controllers.ChatController))
}

func MessageConfigure(app *mvc.Application) {
	app.Router.Use(middleware.Verify)
	app.Register(messageService)
	app.Handle(new(controllers.MessageController))
}

func FriendConfigure(app *mvc.Application) {
	app.Router.Use(middleware.Verify)
	app.Register(friendService)
	app.Handle(new(controllers.FriendController))
}

func SessionConfigure(app *mvc.Application) {
	app.Register(userService)
	app.Handle(new(controllers.SessionController))
}

func UserConfigure(app *mvc.Application) {
	app.Register(userService)
	app.Handle(new(controllers.UserController))
}
