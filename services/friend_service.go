package services

import (
	"HongXunServer/models"
	"HongXunServer/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

const (
	addSuccessCode = 0
	addSuccessMsg  = "添加成功"
	addSelfCode    = 1
	addSelfMsg     = "不能添加自己为好友"
	addNoFoundCode   = 2
	addNoFoundMsg    = "好友不存在"
)

type FriendService interface {
	AddFriend(ownerId, friendId primitive.ObjectID) models.Response
	LoadFriend(ownerId primitive.ObjectID) models.Response
}

type friendService struct {
	friendRepository repositories.FriendRepository
	userRepository   repositories.UserRepository
}

func (s *friendService) AddFriend(ownerId, friendId primitive.ObjectID) models.Response {
	log.Println("AddFriend:", friendId)

	if ownerId == friendId {
		return models.Response{
			ErrCode: addSelfCode,
			ErrMsg:  addSelfMsg,
			Data:    nil,
		}
	}

	_, err := s.userRepository.FindById(friendId)
	if err != nil {
		return models.Response{
			ErrCode: addNoFoundCode,
			ErrMsg: addNoFoundMsg,
			Data: nil,
		}
	}

	_ = s.friendRepository.Save(&models.Friend{
		Id:         primitive.NewObjectID(),
		OwnerId:    ownerId,
		FriendId:   friendId,
		CreateTime: time.Now(),
	})

	_ = s.friendRepository.Save(&models.Friend{
		Id:         primitive.NewObjectID(),
		OwnerId:    friendId,
		FriendId:   ownerId,
		CreateTime: time.Now(),
	})

	return models.Response{
		ErrCode: addSuccessCode,
		ErrMsg:  addSuccessMsg,
		Data:    nil,
	}
}

func (s *friendService) LoadFriend(ownerId primitive.ObjectID) models.Response {
	friends, _ := s.friendRepository.FindAllByOwnerId(ownerId)
	var users []*models.UserMini
	for _, friend := range friends {
		user, _ := s.userRepository.FindById(friend.FriendId)
		users = append(users, &models.UserMini{
			Nickname: user.Nickname,
			Email: user.Email,
			Icon: user.Icon,
			Signature: user.Signature,
		})
	}

	return models.Response{
		ErrCode: 0,
		ErrMsg:  "成功",
		Data:    users,
	}
}

func NewFriendService(rf repositories.FriendRepository, ru repositories.UserRepository) FriendService {
	log.Println("NewFriendService")
	return &friendService{friendRepository: rf, userRepository: ru}
}
