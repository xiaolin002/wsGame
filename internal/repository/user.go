package repository

import "wsprotGame/internal/repository/dao"

/**
 * @Description
 * @Date 2025/4/22 14:44
 **/

type UserRepository interface {
	FindByAP(username string) (string, bool)
	CreateUser(username string, password string) error
}

type UserCacheRepository struct {
	dao dao.UserDao
}

func (r *UserCacheRepository) CreateUser(username string, password string) error {
	err := r.dao.InsertUser(username, password)
	if err != nil {
		return err

	}
	return nil
}

func NewUserCacheRepository(dao dao.UserDao) UserRepository {
	return &UserCacheRepository{dao: dao}
}

func (r *UserCacheRepository) FindByAP(username string) (string, bool) {
	return r.dao.FindByAP(username)
}
