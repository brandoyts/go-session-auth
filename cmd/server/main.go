package main

import (
	"context"
	"log"

	"github.com/brandoyts/go-session-auth/infrastructure/mongodb"
	"github.com/brandoyts/go-session-auth/infrastructure/redisClient"
	"github.com/brandoyts/go-session-auth/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	redis, err := redisClient.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()

	db, err := mongodb.NewMongoDb("go-session-auth", "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.Client().Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	userRepository := user.NewMongoRepository(db)
	userService := user.NewService(userRepository)

	app := fiber.New()
	app.Use(logger.New())

	apiRouter := app.Group("/api/v1")

	apiRouter.Get("/health-check", func(c *fiber.Ctx) error {
		return c.SendString("healthy")
	})

	// user router
	userRouter := apiRouter.Group("/users")
	userRouter.Post("/create", func(c *fiber.Ctx) error {
		requestBody := new(user.User)

		err := c.BodyParser(requestBody)
		if err != nil {
			return err
		}

		result, err := userService.Create(context.Background(), *requestBody)
		if err != nil {
			return err
		}

		return c.JSON(result)

	})
	userRouter.Get(":id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := userService.FindUserById(context.Background(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)

	})

	app.Listen(":8000")
}
