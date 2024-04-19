package routes

import (
	"backend-test/config"
	"backend-test/models"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func DProfile(biodata models.Biodata) models.Profile {
	return models.Profile{
		ID:        biodata.User.ID,
		Username:  biodata.User.Username,
		Fullname:  biodata.User.Fullname,
		Email:     biodata.User.Email,
		Telpon:    biodata.User.Telpon,
		Ttl:       biodata.User.Ttl,
		Kode:      biodata.User.Kode,
		Mbti:      biodata.User.Mbti,
		Motivasi:  biodata.Motivasi,
		Kekuatan:  biodata.Kekuatan,
		Kelemahan: biodata.Kelemahan,
		Kemampuan: biodata.Kemampuan,
	}
}

func DetailProfile(c *gin.Context) {
	kode := c.Param("kode")
	userid := uint(c.MustGet("jwt_user_id").(float64))
	var item models.Biodata

	if config.DB.Preload("User", "kode = ?", kode).First(&item, "user_id = ?", userid).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  "Akses Elor!",
			"message": "Elor, Data Pengguna tidak ditemukan",
		})
		c.Abort()
		return
	}

	res_data := DProfile(item)

	c.JSON(200, gin.H{
		"status": "Berhasil akses endpoint!, silahkan update profile pada link yang sudah disediakan",
		"data":   res_data,
		"link":   "http://127.0.0.1:8080/api/v1/profile/update/",
	})
}

func PostProfile(c *gin.Context) {
	var olditem models.Biodata
	var pengguna models.User
	userid := uint(c.MustGet("jwt_user_id").(float64))
	if config.DB.First(&olditem, "user_id = ?", userid).RecordNotFound() {
		c.JSON(400, gin.H{
			"status":  "Elor akses data!",
			"message": "Data kamu tidak ditemukan, silahkan ulangi step dari awal ya :) !",
		})
		c.Abort()
		return
	}

	// Generate a random 6-digit integer
	randomNum := rand.Intn(1000000) // Generates a random integer between 0 and 999999

	config.DB.Model(&pengguna).Where("id = ?", userid).Updates(models.User{
		Telpon: c.PostForm("telpon"),
		Ttl:    c.PostForm("ttl"),
		Mbti:   c.PostForm("mbti"),
	})

	config.DB.Model(&olditem).Where("user_id = ?", userid).Updates(models.Biodata{
		Motivasi:  c.PostForm("motivasi"),
		Kekuatan:  c.PostForm("kekuatan"),
		Kelemahan: c.PostForm("kelemahan"),
		Kemampuan: c.PostForm("kemampuan"),
	})

	item := models.Submission{
		UserId:      userid,
		Passcode:    fmt.Sprint(randomNum),
		PorfilePass: true,
		IsPass:      false,
	}

	config.DB.Create(&item)

	c.JSON(200, gin.H{
		"status":   "Berhasil update profile kamu, langkah terakhir reedem passcode kamu ke link yang telah di tentukan",
		"passcode": randomNum,
		"link":     "http://127.0.0.1:8080/api/v1/profile/check-passcode/{kode passcode}",
	})
}

func PascodeCheck(c *gin.Context) {
	passcode := c.Param("passcode")
	userid := uint(c.MustGet("jwt_user_id").(float64))
	var sub models.Submission

	if config.DB.Find(&sub, "user_id = ? AND passcode= ?", userid, passcode).RecordNotFound() {
		c.JSON(400, gin.H{
			"status":  "Gagal reedem passcode!",
			"message": "Passcode yang anda masukan tidak valid, silahkan masukan passcode yang valid!",
		})
		c.Abort()
		return
	}

	tanggal := time.Now()

	config.DB.Model(&sub).Where("passcode = ? AND user_id = ?", passcode, userid).Updates(models.Submission{
		DatePass: tanggal.Format("2006-01-02 15:04:05"),
		IsPass:   true,
	})

	c.JSON(200, gin.H{
		"status":  "Reedem Passcode Berhasil!",
		"message": "Selamat ya, kamu lolos pada test kali ini, silahkan kasih tahu admin kalau test mu sudah selesai!",
	})
}
