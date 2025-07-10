package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/brandoyts/go-session-auth/infrastructure/mongodb"
	"github.com/brandoyts/go-session-auth/infrastructure/redisClient"
	"github.com/brandoyts/go-session-auth/internal/auth"
	"github.com/brandoyts/go-session-auth/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type services struct {
	auth *auth.Service
	user *user.Service
}

type handlers struct {
	user *user.FiberHandler
	auth *auth.FiberHandler
}

type dependency struct {
	db          *mongo.Database
	redisClient *redis.Client
	services    *services
	handlers    *handlers
}

func injectDependencies() dependency {
	redisClient, err := redisClient.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
	})
	if err != nil {
		log.Fatal(err)
	}

	mongoURL := fmt.Sprintf("mongodb://%v", os.Getenv("MONGO_HOST"))
	db, err := mongodb.NewMongoDb(os.Getenv("MONGO_DATABASE_NAME"), mongoURL, options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// user dependency
	userRepository := user.NewMongoRepository(db)
	userService := user.NewService(userRepository)

	// auth dependency
	sessionRepository := auth.NewRedisRepository(redisClient)
	authService := auth.NewService(sessionRepository, userRepository)

	// handlers
	userHandler := user.NewFiberHandler(*userService)
	authHandler := auth.NewFiberHandler(*authService)

	return dependency{
		db:          db,
		redisClient: redisClient,
		services: &services{
			auth: authService,
			user: userService,
		},
		handlers: &handlers{
			user: userHandler,
			auth: authHandler,
		},
	}
}

func main() {
	deps := injectDependencies()

	// disconnect database gracefully
	defer func() {
		if err := deps.db.Client().Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// close redis connection gracefully
	defer deps.redisClient.Close()

	app := fiber.New()

	// register default middlewares
	app.Use(logger.New(), recover.New())

	apiRouter := app.Group("/api/v1")

	// health check router
	apiRouter.Get("/health-check", authSession(deps.redisClient), func(c *fiber.Ctx) error {
		return c.SendString("healthy")
	})

	// user router
	userRouter := apiRouter.Group("/users").Use(authSession(deps.redisClient))
	userRouter.Post("/create", deps.handlers.user.CreateUser)
	userRouter.Get(":id", deps.handlers.user.GetUserById)

	// auth router
	authRouter := apiRouter.Group("/auth")
	authRouter.Post("/login", deps.handlers.auth.Login)
	authRouter.Post("/logout", authSession(deps.redisClient), deps.handlers.auth.Logout)

	app.Listen(":8000")
}
