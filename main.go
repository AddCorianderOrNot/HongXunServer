package main

import (
	"HongXunServer/controllers"
	"HongXunServer/services"
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/kataras/iris/v12/mvc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())
	db := client.Database("hungxun")
	var (
		messagesCollection = db.Collection("messages")
		messageService     = services.NewMessageService(messagesCollection)
		userCollection     = db.Collection("user")
		userService        = services.NewUserService(userCollection)
	)

	message := mvc.New(app.Party("/message"))
	message.Register(messageService)
	message.Router.Use(jwt.HMAC(15*time.Minute, "secret", "itsa16bytesecret").Verify)
	message.Handle(new(controllers.MessageController))

	user := mvc.New(app.Party("/user"))
	user.Register(userService)
	user.Handle(new(controllers.UserController))

	session := mvc.New(app.Party("/session"))
	session.Register(userService)
	session.Handle(new(controllers.SessionController))

	app.Run(
		// Start the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}
