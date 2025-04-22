package dao

/**
 * @Description
 * @Date 2025/4/22 14:42
 **/
import (
	"gorm.io/gorm"
)

type UserDao interface {
	FindByAP(username string) (string, bool)
	InsertUser(username string, password string) error
}

type UserGromDao struct {
	db *gorm.DB
}

func (r *UserGromDao) InsertUser(username string, password string) error {
	// 插入用户到数据库
	// 这里先不请求数据库
	return nil
}

func NewUserGromDao(db *gorm.DB) UserDao {
	return &UserGromDao{db: db}
}

// FindByAP 根据用户名和密码查找用户
func (r *UserGromDao) FindByAP(username string) (string, bool) {
	var user struct {
		Password string
		UserName string
	}
	//if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
	//	return "", false
	//}
	user.UserName = username

	// 这里先不请求数据库
	user.Password = "testpassword"

	return user.Password, true
}
