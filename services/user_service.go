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
	registerSuccessCode = 0
	registerSuccessMsg  = "注册成功"
	registerExistCode   = 1
	registerExistMsg    = "用户已存在"
	registerErrorCode   = 2
	registerErrorMsg    = "未知错误"
	verifySuccessCode   = 0
	verifySuccessMsg    = "验证成功"
	verifyExistCode     = 1
	verifyExistMsg      = "用户不存在"
	verifyErrorCode     = 2
	verifyErrorMsg      = "验证失败"
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
		Keys:    bsonx.Doc{{"nickname", bsonx.Int32(1)}},
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
			ErrCode: registerExistCode,
			ErrMsg:  registerExistMsg,
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
			ErrCode: registerErrorCode,
			ErrMsg:  registerErrorMsg,
			Data:    nil,
		}
	}
	return models.Response{
		ErrCode: registerSuccessCode,
		ErrMsg:  registerSuccessMsg,
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
				UserId: user.Id,
			}
			accessToken, _ := j.Token(claims)
			user.Token = accessToken
			return models.Response{
				ErrCode: verifySuccessCode,
				ErrMsg:  verifySuccessMsg,
				Data:    user,
			}
		} else {
			log.Println("密码错误")
			return models.Response{
				ErrCode: verifyErrorCode,
				ErrMsg:  verifyErrorMsg,
			}
		}
	} else {
		log.Println("用户不存在")
		return models.Response{
			ErrCode: verifyExistCode,
			ErrMsg:  verifyExistMsg,
		}
	}
}
