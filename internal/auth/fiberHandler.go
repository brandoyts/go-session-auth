package auth

import (
	"context"
	"net/http"

	"github.com/brandoyts/go-session-auth/internal/shared"
	"github.com/brandoyts/go-session-auth/internal/user"
	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	authService Service
}

func NewFiberHandler(authService Service) *FiberHandler {
	return &FiberHandler{authService: authService}
}

func (fh *FiberHandler) Login(c *fiber.Ctx) error {
	requestBody := new(user.User)

	err := c.BodyParser(requestBody)
	if err != nil {
		return err
	}

	result, err := fh.authService.Login(context.Background(), *requestBody)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     shared.COOKIE_SESSION_TOKEN,
		Value:    result.Key,
		Expires:  result.TTL,
		HTTPOnly: true,
	})

	return c.SendStatus(http.StatusOK)
}

func (fh *FiberHandler) Logout(c *fiber.Ctx) error {
	sessionToken := c.Cookies(shared.COOKIE_SESSION_TOKEN)

	err := fh.authService.Logout(context.Background(), sessionToken)
	if err != nil {
		return err
	}
	c.ClearCookie(shared.COOKIE_SESSION_TOKEN)

	return c.SendStatus(http.StatusOK)
}
