package dao

import "gorm.io/gorm"

/**
 * @Description
 * @Date 2025/5/25 10:44
 **/
func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
