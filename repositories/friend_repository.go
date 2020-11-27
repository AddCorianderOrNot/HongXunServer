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

type FriendRepository interface {
	FindAllByOwnerId(ownerId primitive.ObjectID) ([]*models.Friend, error)
	UpdateTime(ownerId, friendId primitive.ObjectID, time int64) (*models.Friend, error)
	FindByBothId(ownerId, friendId primitive.ObjectID) (*models.Friend, error)
	Save(friend *models.Friend) error
}

type friendRepository struct {
	c *mongo.Collection
}

func (r *friendRepository) FindAllByOwnerId(ownerId primitive.ObjectID) ([]*models.Friend, error) {
	var friends []*models.Friend
	findOptions := options.Find().SetLimit(10)
	cur, err := r.c.Find(context.TODO(), bson.D{
		{"owner_id", ownerId},
	}, findOptions)
	if err != nil {
		log.Println(err)
	}
	for cur.Next(context.TODO()) {
		var elem models.Friend
		_ = cur.Decode(&elem)
		friends = append(friends, &elem)
	}
	err = cur.Close(context.TODO())
	return friends, err
}

func (r *friendRepository) UpdateTime(ownerId, friendId primitive.ObjectID, time int64) (*models.Friend, error) {
	_, err := r.c.UpdateOne(context.TODO(), bson.D{{"owner_id", ownerId}, {"friend_id", friendId}}, bson.M{"$set": bson.M{"read_time": time}})
	friend, _ := r.FindByBothId(ownerId, friendId)
	return friend, err
}

func (r *friendRepository) Save(friend *models.Friend) error {
	_, err := r.c.InsertOne(context.TODO(), friend)
	return err
}

func (r *friendRepository) FindByBothId(ownerId, friendId primitive.ObjectID) (*models.Friend, error) {
	log.Println("FindByBothId:", ownerId, friendId)
	var friend models.Friend
	err := r.c.FindOne(context.TODO(), bson.D{{"owner_id", ownerId}, {"friend_id", friendId}}).Decode(&friend)
	return &friend, err
}



func NewFriendRepository(collection *mongo.Collection) FriendRepository {
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("friendIndex").
		SetUnique(false).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bsonx.Doc{{"owner_id", bsonx.Int32(1)}, {"friend_id", bsonx.Int32(1)}},
		Options: indexOpt,
	})
	if err != nil {
		log.Println(err)
	}
	return &friendRepository{c: collection}
}

