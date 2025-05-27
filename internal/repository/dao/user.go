package dao

/**
 * @Description
 * @Date 2025/4/22 14:42
 **/
import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
	"wsprotGame/pkg/uidGenerate"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao interface {
	FindByAP(ctx context.Context, username string) (User, bool)
	InsertUser(ctx context.Context, u User) error
}

type UserGromDao struct {
	db *gorm.DB
	rd redis.Cmdable
}

func (r *UserGromDao) InsertUser(ctx context.Context, u User) error {
	// 注意数据库中唯一索引冲突的问题
	now := time.Now().Unix()
	u.Ctime = now
	u.Utime = now
	// 创建 UIDGenerator 实例
	generator, err := uidGenerate.NewUIDGenerator(123, 45, r.rd)
	if err != nil {
		panic(err)
	}

	// 生成 UID
	uid, err := generator.Generate(context.Background())
	if err != nil {
		panic(err)
	}
	u.Uid = uid
	err = r.db.WithContext(ctx).Create(&u).Error
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); ok {
			const duplicateErr uint16 = 1062
			if me.Number == duplicateErr {
				// 用户冲突  邮箱冲突
				return ErrDuplicateEmail
			}
		}
	}
	return err
}

func NewUserGromDao(db *gorm.DB, rd redis.Cmdable) UserDao {

	dao := &UserGromDao{
		db: db,
		rd: rd,
	}
	return dao
}

// FindByAP 根据用户名和密码查找用户
func (r *UserGromDao) FindByAP(ctx context.Context, username string) (User, bool) {
	var u User
	err := r.db.WithContext(ctx).Where("nick_name= ?", username).First(&u).Error
	if err != nil {
		return u, false
	}
	return u, true
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 唯一索引
	Uid      uint64         `gorm:"unique"`
	NickName string         `gorm:"type=varchar(128)"`
	Phone    sql.NullString `gorm:"unique"`

	// NullString  代表可以为空
	Email    sql.NullString `gorm:"unique"`
	Password string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
	// NullString  代表可以为空

	//Birthday int64
	AboutMe string `gorm:"type=varchar(4096)"`
}
