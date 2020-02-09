package config

import (
	"github.com/wawandx/rest-api-gin/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() {
	var err error

	DB, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:8889)/learn-gin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&models.User{})
  DB.AutoMigrate(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	
	DB.Model(&models.User{}).Related(&models.Article{})
}