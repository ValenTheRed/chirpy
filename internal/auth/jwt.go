package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var signingMethod = jwt.SigningMethodHS256

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(signingMethod, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
		jwt.WithValidMethods([]string{signingMethod.Alg()}),
	)
	if err != nil {
		return uuid.Nil, err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	// authorization format
	// 	Bearer TOKEN_STRING
	tokenString, found := strings.CutPrefix(authorization, "Bearer ")
	if !found {
		return "", errors.New("Error: authorization token format is incorrect")
	}
	return tokenString, nil
}
