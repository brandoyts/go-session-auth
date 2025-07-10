package main

import (
	"context"
	"net/http"

	"github.com/brandoyts/go-session-auth/internal/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func authSession(redisClient *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// check cookie named "session_token"
		cookie := c.Cookies(shared.COOKIE_SESSION_TOKEN)
		if cookie == "" {
			return c.SendStatus(http.StatusUnauthorized)
		}

		cache := redisClient.Get(context.Background(), string(cookie))
		if cache.Err() != nil {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.Next()
	}
}
