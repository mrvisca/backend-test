package routes

import "github.com/gin-gonic/gin"

func HomeFungsi(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "Akses endpoint berhasil!",
		"message": "Hai, selamat datang...! ini test untuk kandidat backend, semoga kamu terpilih ya :)!",
		"link":    "http://127.0.0.1:8080/api/v1/auth/github.com",
	})
}
