package services

import (
	"HongXunServer/models"
	"context"
	"github.com/kataras/iris/v12/middleware/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

const (
	RegisterSuccessCode = 0
	RegisterSuccessMsg  = "注册成功"
	RegisterExistCode   = 1
	RegisterExistMsg    = "用户已存在"
	RegisterErrorCode   = 2
	RegisterErrorMsg    = "未知错误"
	VerifySuccessCode   = 0
	VerifySuccessMsg    = "验证成功"
	VerifyExistCode     = 1
	VerifyExistMsg      = "用户不存在"
	VerifyErrorCode     = 2
	VerifyErrorMsg      = "验证失败"
)

type UserService interface {
	Register(user models.User) models.Response
	Verify(authentication models.Authentication) models.Response
	isExist(email string) (bool, models.User)
}

type userService struct {
	C *mongo.Collection
}

func NewUserService(collection *mongo.Collection) UserService {
	log.Println("NewUserService")
	indexOpt := new(options.IndexOptions)
	indexOpt.SetName("userIndex").
		SetUnique(false).
		SetBackground(true).
		SetSparse(true)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bsonx.Doc{{
			Key:   "email",
			Value: bsonx.String("text"),
		}},
		Options: indexOpt,
	})
	if err != nil {
		log.Println(err)
	}

	return &userService{C: collection}
}

func (s *userService) isExist(email string) (bool, models.User) {
	log.Println("Find:", email)
	var user models.User
	err := s.C.FindOne(context.TODO(), bson.D{{"email", email}}).Decode(&user)
	log.Println(user, err)
	if err == nil {
		return true, user
	} else {
		return false, user
	}
}

func (s *userService) Register(user models.User) models.Response {
	log.Println("Register")
	exist, _ := s.isExist(user.Email)
	log.Println(exist)
	if exist {
		return models.Response{
			ErrCode: RegisterExistCode,
			ErrMsg:  RegisterExistMsg,
			Data:    nil,
		}
	}
	if user.Id.IsZero() {
		user.Id = primitive.NewObjectID()
	}
	log.Println("Insert:", user)
	_, err := s.C.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println(err)
		return models.Response{
			ErrCode: RegisterErrorCode,
			ErrMsg:  RegisterErrorMsg,
			Data:    nil,
		}
	}
	return models.Response{
		ErrCode: RegisterSuccessCode,
		ErrMsg:  RegisterSuccessMsg,
		Data:    nil,
	}
}

func (s *userService) Verify(authentication models.Authentication) models.Response {
	exist, user := s.isExist(authentication.Email)
	log.Println(exist, user)
	if exist {
		if authentication.Password == user.Password {
			log.Println("密码正确")
			j := jwt.HMAC(15*time.Minute, "secret", "itsa16bytesecret")
			claims := models.UserClaims{
				Claims: j.Expiry(jwt.Claims{
					Issuer:   "an-issuer",
					Audience: jwt.Audience{"an-audience"},
				}),
				Email: authentication.Email,
			}
			accessToken, _ := j.Token(claims)
			user.Token = accessToken
			return models.Response{
				ErrCode: VerifySuccessCode,
				ErrMsg:  VerifySuccessMsg,
				Data:    user,
			}
		} else {
			log.Println("密码错误")
			return models.Response{
				ErrCode: VerifyErrorCode,
				ErrMsg:  VerifyErrorMsg,
			}
		}
	} else {
		log.Println("用户不存在")
		return models.Response{
			ErrCode: VerifyExistCode,
			ErrMsg:  VerifyExistMsg,
		}
	}
}
