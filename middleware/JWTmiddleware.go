package middleware

import (
	"fmt"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

func IsAuth() gin.HandlerFunc {
	return checkJWT()
}

func checkJWT() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
	
				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
		
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				context.Set("jwt_user_id", claims["user_id"])
				context.Set("jwt_isAdmin", claims["user_role"])
			} else {
				context.JSON(422, gin.H{"message": "Invalid Token", "error": err})
				context.Abort()
				return
			}
		} else {
			context.JSON(422, gin.H{"message": "Authorization token not provided"})
			context.Abort()
			return
		}
	}
}