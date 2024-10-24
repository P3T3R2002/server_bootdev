package auth

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password),  bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))	
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	} else if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	} else {
		ID, err := token.Claims.GetSubject()
		if err != nil {
			return uuid.Nil, err
		}
		return uuid.Parse(ID)
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	return strings.TrimPrefix(auth, "Bearer "), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	return strings.TrimPrefix(auth, "ApiKey "), nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}