package repository

import (
	"context"
	"database/sql"
	"wsprotGame/internal/domain"
	"wsprotGame/internal/repository/dao"
)

/**
 * @Description
 * @Date 2025/4/22 14:44
 **/

type UserRepository interface {
	FindByAP(ctx context.Context, username string) (domain.User, bool)
	CreateUser(ctx context.Context, user domain.User) error
}

type UserCacheRepository struct {
	dao dao.UserDao
}

func (r *UserCacheRepository) CreateUser(ctx context.Context, user domain.User) error {
	err := r.dao.InsertUser(ctx, r.toDaoUser(user))
	if err != nil {
		return err

	}
	return nil
}

func NewUserCacheRepository(dao dao.UserDao) UserRepository {
	return &UserCacheRepository{dao: dao}
}

func (r *UserCacheRepository) FindByAP(ctx context.Context, username string) (domain.User, bool) {
	u, b := r.dao.FindByAP(ctx, username)
	if b == false {
		return r.toDomainUser(u), false
	}

	return r.toDomainUser(u), true
}

func (r *UserCacheRepository) toDaoUser(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		// valid 取值为true不为空 取值为false为空
		Password: u.Password,
		//Birthday: u.Birthday.UnixMilli(),
		NickName: u.NickName,
		AboutMe:  u.AboutMe,
	}
}
func (r *UserCacheRepository) toDomainUser(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Uid:      u.Uid,
		Phone:    u.Phone.String,
		NickName: u.NickName,
		Password: u.Password,
	}
}
