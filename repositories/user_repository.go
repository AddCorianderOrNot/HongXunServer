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
	Save(user *models.User) error
	Update(id primitive.ObjectID, key, value string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindById(id primitive.ObjectID) (*models.User, error)
	FindByNickname(nickname string) ([]*models.UserMini, error)
	findAllBy(key string, value interface{}) ([]*models.UserMini, error)
	findOneBy(key string, value interface{}) (*models.User, error)
}

type userRepository struct {
	c *mongo.Collection
}

func (r *userRepository) Update(id primitive.ObjectID, key, value string) (*models.User, error) {
	log.Println(key, value)
	_, err := r.c.UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.M{"$set": bson.M{key: value}})
	user, _ := r.FindById(id)
	return user, err
}

func (r *userRepository) Save(user *models.User) error {
	_, err := r.c.InsertOne(context.TODO(), user)
	return err
}

func (r *userRepository) FindByNickname(nickname string) ([]*models.UserMini, error) {
	return r.findAllBy("nickname", nickname)
}

func (r *userRepository) FindById(id primitive.ObjectID) (*models.User, error) {
	return r.findOneBy("_id", id)
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	log.Println("FindByEmail:", email)
	return r.findOneBy("email", email)
}

func (r *userRepository) findOneBy(key string, value interface{}) (*models.User, error) {
	var user models.User
	log.Println("findOneBy", key, value)
	err := r.c.FindOne(context.TODO(), bson.D{{key, value}}).Decode(&user)
	return &user, err
}

func (r *userRepository) findAllBy(key string, value interface{}) ([]*models.UserMini, error) {
	var users []*models.UserMini
	cur, err := r.c.Find(context.TODO(), bson.D{{key, bson.M{"$regex": value, "$options": "$i"}}})
	if err != nil {
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var elem models.UserMini
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
