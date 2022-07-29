package server

import (
	jwtMiddlewares "github.com/4strodev/jwt/internal/middlewares/jwt"
	jwtServices "github.com/4strodev/jwt/internal/services/jwt"
	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

func init() {
	App = fiber.New()

	App.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	App.Post("/login", func(c *fiber.Ctx) error {
		payload := fiber.Map{
			"user": "astro",
			"role": "admin",
		}

		accessToken, err := jwtServices.GenerateAccessToken(payload)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"err": "Error creating access token",
			})
		}

		refreshToken, err := jwtServices.GenerateRefreshToken(payload)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"err": "Error creating refresh token",
			})
		}

		refreshTokenCookie := new(fiber.Cookie)
		refreshTokenCookie.Name = "refresh_token"
		refreshTokenCookie.Value = string(refreshToken)
		refreshTokenCookie.HTTPOnly = true

		accessTokenCookie := new(fiber.Cookie)
		accessTokenCookie.Name = "access_token"
		accessTokenCookie.Value = string(accessToken)
		accessTokenCookie.HTTPOnly = true

		c.Cookie(accessTokenCookie)
		c.Cookie(refreshTokenCookie)

		return c.JSON(fiber.Map{
			"msg": "logged successfully",
		})
	})

	App.Post("/token", func(c *fiber.Ctx) error {
		payload := fiber.Map{
			"user": "astro",
			"role": "admin",
		}

		refreshToken := c.Cookies("refresh_token", "")
		if refreshToken == "" {
			return c.Status(403).JSON(fiber.Map{
				"err": "Token not provided",
			})
		}

		// Verify refresh token
		err := jwtServices.VerifyRefreshtoken([]byte(refreshToken))
		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"err": err.Error(),
			})
		}

		accessToken, err := jwtServices.GenerateAccessToken(payload)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"err": "Error creating access token",
			})
		}

		tokenCookie := new(fiber.Cookie)
		tokenCookie.Name = "access_token"
		tokenCookie.Value = string(accessToken)
		tokenCookie.HTTPOnly = true

		c.Cookie(tokenCookie)

		return c.JSON(fiber.Map{
			"msg": "token refreshed successfully",
		})
	})

	App.Delete("/logout", func(c *fiber.Ctx) error {
		refreshToken := c.Cookies("refresh_token", "")
		jwtServices.RevokeRefreshToken([]byte(refreshToken))

		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
	})

	App.Get("/resource", jwtMiddlewares.AuthenticateJwt, func(c *fiber.Ctx) error {
		return c.SendString("A basic resource with authentication")
	})
}
