package datasource

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client, _ = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
var DB *mongo.Database

func init()  {
	DB = client.Database("HongXun")
}

func ConnectToDatabase() {
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping to database server `%s`: %s\n", "mongodb://localhost:27017", err.Error())
	}
}

func DisconnectToDatabase()  {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Println(err)
	}
}