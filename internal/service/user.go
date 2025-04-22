package service

import (
	"errors"
	"log"
	"wsprotGame/internal/repository"
	proto2 "wsprotGame/proto/gen"
)

type UserService interface {
	HandleLogin(req proto2.LoginRequest) error
	HandleRegister(req proto2.RegisterRequest) error
}

type userService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}
func (s *userService) HandleRegister(req proto2.RegisterRequest) error {
	log.Printf("Received LoginRequest: Username = %s, Password = %s", req.Username, req.Password)
	// 正常需要对密码进行加密处理
	err := s.UserRepo.CreateUser(req.Username, req.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) HandleLogin(req proto2.LoginRequest) error {
	// 从持久化层获取用户信息
	log.Printf("Received LoginRequest: Username = %s, Password = %s", req.Username, req.Password)
	password, exists := s.UserRepo.FindByAP(req.Username)
	if !exists {
		return errors.New("用户名不存在")
	}
	// 这里可能对密码进行加密验证
	if password != req.Password {
		return errors.New("密码错误")
	}

	// 登录成功，返回 nil
	return nil
}
