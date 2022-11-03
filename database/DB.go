package database

import (
	model_v1 "ahripost_deploy/models/v1"
	"ahripost_deploy/tools"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var cfg = tools.Cfg

	var err error
	var db *gorm.DB

	if cfg.DBType == "postgres" {
		dsn := "host=" + cfg.Postgres.Host + " user=" + cfg.Postgres.User + " password=" + cfg.Postgres.Pass + " dbname=" + cfg.Postgres.Name + " port=" + cfg.Postgres.Port + " sslmode=disable TimeZone=" + cfg.TimeZone
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		cfg.Sqlite.Path = "/data/ahripost.db" // friendly for build docker image and set volume
		db, err = gorm.Open(sqlite.Open(cfg.Sqlite.Path), &gorm.Config{})
	}

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&model_v1.User{},
		&model_v1.Token{},
		&model_v1.Project{},
		&model_v1.Item{},
		&model_v1.Member{},
	)

	var count int64
	db.Model(&model_v1.User{}).Count(&count)
	if count == 0 {
		db.Create(&model_v1.User{
			Username: "admin",
			Password: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
			Token:    "",
		})
	}

	DB = db
}
