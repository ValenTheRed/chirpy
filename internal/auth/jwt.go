package auth

import (
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
		return uuid.UUID([16]byte{}), err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID([16]byte{}), err
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID([16]byte{}), err
	}
	return userID, nil
}
