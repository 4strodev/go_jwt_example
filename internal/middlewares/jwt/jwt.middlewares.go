package jwt

import (
	jwtServices "github.com/4strodev/jwt/internal/services/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthenticateJwt(c *fiber.Ctx) error {
	token := c.Cookies("access_token", "")

	if token == "" {
		return c.Status(400).JSON(fiber.Map{
			"err": "Token not provided",
		})
	}

	err := jwtServices.VerifyAccessToken([]byte(token))
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Next()
}
