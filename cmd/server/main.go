package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nozgurozturk/noo-analytics/internal/config"
	"github.com/nozgurozturk/noo-analytics/internal/router"
	"github.com/nozgurozturk/noo-analytics/internal/storage"
)

func main() {
	cnf := config.Config()

	mongoDB, err := storage.MongoConnect(cnf.Mongo)
	if err != nil {
		return
	}
	redis, err := storage.RedisConnect(cnf.Redis)
	if err != nil {
		return
	}
	app := fiber.New()
	app.Use(cors.New())

	router.New(&router.Router{
		App:     app,
		MongoDB: mongoDB,
		Redis:   redis,
		Config:  cnf.Server,
	})

	_ = app.Listen(cnf.Server.Port)
}
