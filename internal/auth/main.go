package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	ISSUER = "mini-youtube"
)

func HashPassword(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}

	return string(data), nil
}

func VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateJWTToken(userID uuid.UUID, expiresIn time.Duration, secretKey string) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    ISSUER,
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(secretKey))
}

func VerifyJWT(jwtSecret, token string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != ISSUER {
		return uuid.Nil, fmt.Errorf("issuer don't match")
	}

	rawId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(rawId)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
