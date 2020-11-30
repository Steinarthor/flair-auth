package main

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email string
	jwt.StandardClaims
}

// HashAndSalt generates a hashed password
func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Println(err.Error())
	}
	return string(hash)
}

// ComparePassword compares a hashed password with plain text
func ComparePassword(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// GenerateJWT returns a JWT token on success and error on failure.
func GenerateJWT(email string) (error, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	JwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	jwtSignedKeyBytes := []byte(JwtSigningKey)

	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtSignedKeyBytes)

	if err != nil {
		return err, ""
	}

	return nil, ss
}
