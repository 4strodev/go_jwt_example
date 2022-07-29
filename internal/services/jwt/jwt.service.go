package jwt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/4strodev/jwt/pkg/utils/slices"
	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	// tokens secrets
	refreshTokenSecret []byte
	accessTokenSecret  []byte

	// refresh token list
	refreshTokens []string
)

func init() {
	refreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
	accessTokenSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	if string(refreshTokenSecret) == "" {
		log.Fatal("Env variables don't loaded")
	}
}

func GenerateRefreshToken(payload fiber.Map) ([]byte, error) {
	var err error
	var signedToken []byte

	token := jwt.New()

	// Setting toekn payload
	for key, value := range payload {
		token.Set(key, value)
	}

	signedToken, err = jwt.Sign(token, jwt.WithKey(jwa.HS256, refreshTokenSecret))
	refreshTokens = append(refreshTokens, string(signedToken))

	return signedToken, err
}

func GenerateAccessToken(payload fiber.Map) ([]byte, error) {
	var err error
	var signedToken []byte

	token := jwt.New()
	expirationDate := time.Now().Add(time.Second * 10)

	// Setting token payload
	token.Set(jwt.ExpirationKey, expirationDate)
	for key, value := range payload {
		token.Set(key, value)
	}

	signedToken, err = jwt.Sign(token, jwt.WithKey(jwa.HS256, accessTokenSecret))

	if err != nil {
		return []byte{}, err
	}

	return signedToken, err
}

func VerifyAccessToken(token []byte) error {
	var err error

	_, err = jwt.Parse(token, jwt.WithKey(jwa.HS256, accessTokenSecret))
	if err != nil {
		if !jwt.IsValidationError(err) {
			return fmt.Errorf("Error verifying access token")
		}

		if errors.Is(err, jwt.ErrTokenExpired()) {
			return fmt.Errorf("Token expired")
		}
	}
	return err
}

func VerifyRefreshtoken(token []byte) error {
	var err error

	_, err = jwt.Parse(token, jwt.WithKey(jwa.HS256, refreshTokenSecret))
	if err != nil {
		if jwt.IsValidationError(err) {
			return fmt.Errorf("Invalid token")
		}

		log.Println(err)
		return fmt.Errorf("Error validating token")
	}

	if !slices.Contains(refreshTokens, string(token)) {
		return fmt.Errorf("Refresh token invalid")
	}

	return err
}

func RevokeRefreshToken(refreshToken []byte) {
	refreshTokens = slices.Filter(refreshTokens, func(token string) bool {
		return token == string(refreshToken)
	})
}
