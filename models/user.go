package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string
	Fullname string
	Email    string `gorm:"unique;not null"`
	SocialId string
	Provider string
	Avatar   string
	Role     bool `gorm:"default:0"`
	Telpon   string
	Ttl      string
	Kode     string
	Mbti     string
}
