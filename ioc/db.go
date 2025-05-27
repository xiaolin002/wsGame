package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"wsprotGame/internal/repository/dao"
)

func InitDB() *gorm.DB {
	dsn := "root:123456@tcp(47.93.78.200:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库驱动错误")
	}
	err = dao.InitTables(db)
	if err != nil {
		panic("数据库初始化错误")
	}
	return db
}
