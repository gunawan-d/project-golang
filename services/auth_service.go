package services

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{}

// CreateToken akan membuat token JWT berdasarkan payload
var SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

func (repo *AuthService) CreateToken(payload map[string]interface{}) (string, error) {
    // Pastikan SECRET_KEY ada
    if len(os.Getenv("JWT_SECRET")) == 0 {
        log.Fatalf("JWT_SECRET is not set. Current value: %s", os.Getenv("JWT_SECRET"))
    }

    // Buat claims untuk token
    claims := jwt.MapClaims{
        "exp":    time.Now().Add(time.Hour).Unix(), // Token berlaku selama 1 jam
        "email":  payload["email"],                // Ambil email dari parameter payload
        "name":   payload["name"],
        "roleID": payload["roleID"],
    }

    // Tambahkan payload ke klaim token
    for key, value := range payload {
        claims[key] = value
    }

    // Buat token JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(SECRET_KEY)
}
