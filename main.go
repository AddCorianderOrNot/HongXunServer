package main

import (
	"HongXunServer/config"
	"HongXunServer/datasource"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug")

	datasource.ConnectToDatabase()
	defer datasource.DisconnectToDatabase()

	mvc.Configure(app.Party("/user"), config.UserConfigure)
	mvc.Configure(app.Party("/session"), config.SessionConfigure)
	mvc.Configure(app.Party("/friend"), config.FriendConfigure)
	mvc.Configure(app.Party("/message"), config.MessageConfigure)
	mvc.Configure(app.Party("/chat"), config.ChatConfigure)

	_ = app.Run(
		// Start the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}


