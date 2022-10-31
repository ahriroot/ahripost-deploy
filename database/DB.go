package database

import (
	model_v1 "ahripost_deploy/models/v1"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dsn := "host=127.0.0.1 user=postgres password=Aa12345. dbname=ahripost port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db.AutoMigrate(
		&model_v1.User{},
		&model_v1.Token{},
		&model_v1.Project{},
		&model_v1.Item{},
	)

	DB = db
}
