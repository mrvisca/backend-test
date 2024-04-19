package routes

import (
	"backend-test/config"
	"backend-test/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

func GetProfile(biodata models.Biodata) models.Profile {
	return models.Profile{
		ID:        biodata.UserId,
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

func CheckToken(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "Login Sukses!",
		"message": "Login Sukses! Token Valid",
	})
}

func RedirectHandler(c *gin.Context) {
	provider := c.Param("provider")

	// In this case we use a map to store our secrets, but you can use dotenv or your framework configuration
	// for example, in revel you could use revel.Config.StringDefault(provider + "_clientID", "") etc.
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GH"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GH"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_G"),
			"clientSecret": os.Getenv("CLIENT_SECRET_G"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL)
}

func CallbackHandler(c *gin.Context) {
	// Ambil parameter kueri untuk status dan kode
	state := c.Query("state")
	code := c.Query("code")
	provider := c.Query("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, user)
	var newToken = createToken(&newUser)

	c.JSON(200, gin.H{
		"data":  newUser,
		"token": newToken,
		// "gh_token": token,
		"message": "Berhasil Login Aplikasi! Pastikan anda mengingat kode pengguna dan token dalam request selanjutnya",
		"link":    "http://127.0.0.1:8080/api/v1/profile/detail/{kode}",
	})
}

func getOrRegisterUser(provider string, user *structs.User) models.User {
	var userData models.User
	var count int64

	config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)
	config.DB.Model(&models.User{}).Count(&count)

	nomor := uint(count) + 1

	if userData.ID == 0 {
		newUser := models.User{
			Username: user.Username,
			Fullname: user.FullName,
			Email:    user.Email,
			SocialId: user.ID,
			Provider: provider,
			Avatar:   user.Avatar,
			Role:     true,
			Kode:     "BE_00" + strconv.Itoa(int(uint(nomor))),
		}
		config.DB.Create(&newUser)

		newBio := models.Biodata{
			UserId: newUser.ID,
		}
		config.DB.Create(&newBio)

		return newUser
	} else {
		return userData
	}
}

func createToken(user *models.User) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.Role,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
		"iat":       time.Now().Unix(),
	})

	// Sign and get the complete encode token as a string using the secret
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println(err)
	}

	return tokenString
}

func ProfileUser(c *gin.Context) {
	user_id := int(c.MustGet("jwt_user_id").(float64))
	var user models.Biodata

	// Eager loading (mengakses data pada beberapa tabel yang memiliki relasi untuk ditampilakan bersamaan tanpa harus query 1 per satu)
	config.DB.Where("user_id = ?", user_id).Preload("User", "id = ?", user_id).Find(&user) // Office dari user struct Offices untuk penghubung

	res_data := GetProfile(user)

	c.JSON(200, gin.H{
		"status": "Berhasil akses data profil pengguna",
		"data":   res_data,
	})
}
