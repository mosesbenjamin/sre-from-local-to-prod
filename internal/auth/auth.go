package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	// TokenAccess -
	TokenTypeAccess TokenType = "student-access"
	// TokenRefresh
	TokenTypeRefresh TokenType = "student-refresh"
)

// ErrNoAuthHeaderIncluded
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// HashPassword
func HashPassword(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CheckPasswordHash
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT
func MakeJWT(studentID, tokenSecret string, expiresIn time.Duration, tokenType TokenType) (string, error) {
	signingKey := []byte(tokenSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(tokenType),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   studentID,
	})
	return token.SignedString(signingKey)
}

// RefreshToken
func RefreshToken(tokenstring, tokenSecret string) (string, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenstring,
		&claimsStruct,
		func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	studentID, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeRefresh) {
		return "", errors.New("invalid issuer")
	}

	newToken, err := MakeJWT(
		studentID,
		tokenSecret,
		time.Hour,
		TokenTypeAccess,
	)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

// ValidateJWT
func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimStruct,
		func(t *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	studentID, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeAccess) {
		return "", errors.New("invalid issuer")
	}
	return studentID, nil
}

// GetBearerToken
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}

// GetAPIKey
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
