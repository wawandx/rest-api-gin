package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Articles []Article
	Username string
	FullName string 
	Email string `gorm: "unique_index"`
	SocialId string
	Provider string
	Avatar string
	Role bool `gorm: "default:0"`
}