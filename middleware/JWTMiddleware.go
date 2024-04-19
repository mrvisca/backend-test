package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWT_SECRET = os.Getenv("JWT_SECRET")

func IsAuth() gin.HandlerFunc {
	return checkJWT(true)
}

func checkJWT(middlewareAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Token string di dapat dari header postman
		authHeader := c.Request.Header.Get("Authorization")
		// Mengambil token dari "Bearer <token>"
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			// Parse takes the token string and a function for looking up the key. The latter is especially
			// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
			// head of the token to identify which key to use, but the parsed token (head and claims) is provided
			// to the callback, providing flexibility.
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userRole := bool(claims["user_role"].(bool))
				c.Set("jwt_user_id", claims["user_id"])
				c.Set("jwt_isAdmin", claims["user_role"])

				if middlewareAdmin && !userRole {
					c.JSON(403, gin.H{
						"status":  "Elor",
						"message": "Akses hanya untuk yang sudah melakukan register / login",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(422, gin.H{
					"status":  "Elor Token",
					"message": "Token tidak valid, pastikan token valid!",
					"error":   err,
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(422, gin.H{
				"status":  "Elor login",
				"message": "Autorisasi diperlukan untuk akses endpoint ini",
			})
			c.Abort()
			return
		}
	}
}
