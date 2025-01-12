package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type CustomClaims struct {
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.StandardClaims
}

// jwt tocken generation
func GenerateToken(email string, id string) (string, error) {
	godotenv.Load()

	//strcture of the token
	claims := &CustomClaims{
		Email: email,
		Id: id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    "192.168.0.155",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
