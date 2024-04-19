package models

import "github.com/jinzhu/gorm"

type Submission struct {
	gorm.Model
	UserId      uint
	Passcode    string
	PorfilePass bool
	DatePass    string
	IsPass      bool
}
