package user

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type FiberHandler struct {
	userService Service
}

func NewFiberHandler(userService Service) *FiberHandler {
	return &FiberHandler{userService: userService}
}

func (fh *FiberHandler) CreateUser(c *fiber.Ctx) error {
	requestBody := new(User)

	err := c.BodyParser(requestBody)
	if err != nil {
		return err
	}

	result, err := fh.userService.Create(context.Background(), *requestBody)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

func (fh *FiberHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := fh.userService.FindUserById(context.Background(), id)
	if err != nil {
		return err
	}

	return c.JSON(result)
}
