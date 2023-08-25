package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func Init() (*gorm.DB, error) {
	// 连接数据库
	dsn := "root:1111@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("mysql connection have a error:", err)
		return nil, err
	}
	DB.AutoMigrate(&Favorite{})
	return DB, nil
}
