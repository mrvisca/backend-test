package main

import (
	"backend-test/config"
	"backend-test/middleware"
	"backend-test/routes"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	config.InitDB()
	defer config.DB.Close()
	gotenv.Load()

	router := gin.Default()

	v1 := router.Group("api/v1/")
	{
		v1.GET("/", routes.HomeFungsi)

		// Route Autentitakasi
		v1.GET("/auth/:provider", routes.RedirectHandler)
		v1.GET("/auth/:provider/callback", routes.CallbackHandler)

		// Step 1
		bio := v1.Group("/profile/")
		{
			bio.GET("detail/:kode", middleware.IsAuth(), routes.DetailProfile)
			bio.POST("/update/", middleware.IsAuth(), routes.PostProfile)
			bio.GET("check-passcode/:passcode", middleware.IsAuth(), routes.PascodeCheck)
		}
	}

	router.Run()
}
