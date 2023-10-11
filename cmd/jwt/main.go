package main

import (
	"fmt"

	"github.com/gerich/s3test/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	cfg := config.New()
	for _, u := range cfg.HTTP.AllowedUsers {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user"] = u
		tokenString, err := token.SignedString([]byte(cfg.HTTP.JWTSecret))
		if err != nil {
			panic(err)
		}
		fmt.Println(u, " | ", tokenString)
	}
}
