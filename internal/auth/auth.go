package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {

	bytes := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {

	pswdBytes := []byte(password)
	hashBytes := []byte(hash)

	if err := bcrypt.CompareHashAndPassword(hashBytes, pswdBytes); err != nil {
		return err

	}

	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{

		Issuer:    string(TokenTypeAccess),                             // кто выдал
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),                // когда выдан
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)), // когда инстекает
		Subject:   userID.String(),                                     // для кого

	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Fatalf("Couldn't create jws by error: %v", err)
	}
	return tokenString, nil
}
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	str := headers.Get("Authorization")
	if str == "" {
		return "", fmt.Errorf("autorization token value is empty")
	}
	tokenStr := strings.TrimSpace(strings.TrimPrefix(str, "Bearer"))
	if tokenStr == "" {
		return "", fmt.Errorf("autorization token value is empty")
	}
	return tokenStr, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
