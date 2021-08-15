package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type TokenExpiredError struct {
}

func (*TokenExpiredError) Error() string {
	return "token expired"
}

// secret key being used to sign tokens
var (
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt token and assign a username to it's claims and return it
func GenerateToken(username string, id int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the username in it's claims
func ParseToken(tokenStr string) (string, int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if token == nil || err != nil {
		return "", -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err = nil
		if claims.VerifyExpiresAt(time.Now().Unix(), false) {
			err = &TokenExpiredError{}
		}

		username := claims["username"].(string)
		id := int(claims["id"].(float64))
		return username, id, err
	} else {
		return "", -1, err
	}
}
