package services

import (
	"HongXunServer/middleware"
	"HongXunServer/models"
	"HongXunServer/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
	Register(user *models.User) models.Response
	Verify(authentication *models.Authentication) models.Response
	SearchUser(nickname string) models.Response
	AutoLogin(id primitive.ObjectID) models.Response
	UpdateUser(id primitive.ObjectID, key, value string) models.Response
	isExist(email string) (bool, *models.User)
}

type userService struct {
	r repositories.UserRepository
}

func NewUserService(r repositories.UserRepository) UserService {
	log.Println("NewUserService")
	return &userService{r: r}
}

func (s *userService) UpdateUser(id primitive.ObjectID, key, value string) models.Response {
	user, err := s.r.Update(id, key, value)
	if err != nil {
		log.Println(err)
		return models.Response{
			ErrCode: 1,
			ErrMsg: "用户信息更新失败",
			Data: nil,
		}
	}
	user.Token = middleware.GenerateToken(user.Nickname, user.Id, user.Email)
	return models.Response{
		ErrCode: 0,
		ErrMsg: "更新成功",
		Data: user,
	}
}

func (s *userService) isExist(email string) (bool, *models.User) {
	log.Println("Find:", email)
	user, err := s.r.FindByEmail(email)
	log.Println(user, err)
	if err == nil {
		return true, user
	} else {
		return false, user
	}
}

func (s *userService) AutoLogin(id primitive.ObjectID) models.Response  {
	user, err := s.r.FindById(id)
	if err != nil {
		log.Println(err)
		return models.Response{
			ErrCode: 1,
			ErrMsg: "自动登录失败",
			Data: nil,
		}
	}
	user.Token = middleware.GenerateToken(user.Nickname, user.Id, user.Email)
	return models.Response{
		ErrCode: 0,
		ErrMsg: "自动登录成功",
		Data: user,
	}
}

func (s *userService) SearchUser(nickname string) models.Response {
	log.Println("Find:", nickname)
	users, err := s.r.FindByNickname(nickname)
	if err != nil {
		log.Println(err)
		return models.Response{
			ErrCode: 1,
			ErrMsg:  "未知错误",
			Data:    nil,
		}
	}
	if len(users) == 0 {
		return models.Response{
			ErrCode: 2,
			ErrMsg:  "用户不存在",
			Data:    nil,
		}
	} else {
		return models.Response{
			ErrCode: 0,
			ErrMsg:  "查询成功",
			Data:    users,
		}
	}

}

func (s *userService) Register(user *models.User) models.Response {
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
	log.Println("Save:", user)
	err := s.r.Save(user)
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

func (s *userService) Verify(authentication *models.Authentication) models.Response {
	exist, user := s.isExist(authentication.Email)
	log.Println(exist, user)
	if exist {
		if authentication.Password == user.Password {
			user.Token = middleware.GenerateToken(user.Nickname, user.Id, user.Email)
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
