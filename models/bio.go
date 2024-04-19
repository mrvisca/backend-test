package models

import (
	"github.com/jinzhu/gorm"
)

type Biodata struct {
	gorm.Model
	User      User `gorm:"foreignkey:UserId"`
	UserId    uint
	Motivasi  string
	Kekuatan  string
	Kelemahan string
	Kemampuan string
}

type Profile struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	Telpon    string `json:"telpon"`
	Ttl       string `json:"ttl"`
	Kode      string `json:"kode"`
	Mbti      string `json:"mbti"`
	Motivasi  string `json:"motivasi"`
	Kekuatan  string `json:"kekuatan"`
	Kelemahan string `json:"kelemahan"`
	Kemampuan string `json:"kemampuan"`
}
