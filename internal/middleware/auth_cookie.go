package middleware

import (
	"hacko-app/pkg/jwthandler"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("accessToken")

	if cookie == "" {
		log.Warn().Msg("middleware::AuthMiddleware - Unauthorized [Token not found in cookie]")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Token expired",
			"success": false,
		})
	}

	claims, err := jwthandler.ParseTokenString(cookie)
	if err != nil {
		if jwthandler.IsTokenExpired(err) {
			log.Warn().Msg("middleware::AuthMiddleware - Token expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: Token expired",
				"success": false,
			})
		}

		log.Error().Err(err).Msg("middleware::AuthMiddleware - Invalid token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid token",
			"success": false,
		})
	}

	c.Locals("user_id", claims.UserId)
	c.Locals("role", claims.Role)

	return c.Next()
}
