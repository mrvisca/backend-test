package config

import (
	"backend-test/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open("mysql", "root:@(localhost)/backendtest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Koneksi ke database gagal bosque!")
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Biodata{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	DB.AutoMigrate(&models.Submission{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	DB.Model(&models.User{}).Related(&models.Biodata{}).Related(&models.Submission{})
}
