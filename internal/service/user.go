package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"wsprotGame/internal/domain"
	"wsprotGame/internal/repository"
	"wsprotGame/server/connection"
)

type UserService interface {
	HandleLogin(ctx context.Context, conn *connection.ConnInfo, user domain.User) error
	HandleRegister(ctx context.Context, user domain.User) error
}

type userService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{UserRepo: userRepo}
}
func (s *userService) HandleRegister(ctx context.Context, user domain.User) error {
	//对密码进行加密处理
	pwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(pwd)
	return s.UserRepo.CreateUser(ctx, user)

}

func (s *userService) HandleLogin(ctx context.Context, conn *connection.ConnInfo, user domain.User) error {
	// 从持久化层获取用户信息
	dao, exists := s.UserRepo.FindByAP(ctx, user.NickName)
	if !exists {
		return errors.New("用户名不存在")
	}
	err := bcrypt.CompareHashAndPassword([]byte(dao.Password), []byte(user.Password))
	if err != nil {
		return errors.New("密码错误")
	}
	conn.SetUid(dao.Uid)
	// 登录成功，返回 nil
	return nil
}
