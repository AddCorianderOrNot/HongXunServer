package repositories

import (
	"HongXunServer/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindById(id primitive.ObjectID) (*models.User, error)
	Save(user *models.User) error
	FindAll() ([]*models.User, error)
}

type userRepository struct {
	c *mongo.Collection
}

func (r *userRepository) FindById(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.c.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	return &user, err
}

func (r *userRepository) Save(user *models.User) error {
	_, err := r.c.InsertOne(context.TODO(), user)
	return err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.c.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	return &user, err
}

func (r *userRepository) FindAll() ([]*models.User, error) {
	var users []*models.User
	cur, err := r.c.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var elem models.User
		_ = cur.Decode(&elem)
		users = append(users, &elem)
	}
	err = cur.Close(context.TODO())
	return users, err
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("userIndex").
		SetUnique(false).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bsonx.Doc{{"nickname", bsonx.Int32(1)}},
		Options: indexOpt,
	})
	if err != nil {
		log.Println(err)
	}
	return &userRepository{c: collection}
}
