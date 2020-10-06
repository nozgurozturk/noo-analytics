package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"github.com/nozgurozturk/noo-analytics/pkg/trace"
	"net/http"
)

func TraceRooter(app fiber.Router, service trace.Service) {
	app.Post("/", createTrace(service))
}

func createTrace(s trace.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request entities.TraceDTO

		if err := c.BodyParser(&request); err != nil {
			jsonErr := errors.BadRequest("Invalid Trace Object")
			c.Status(jsonErr.Status).JSON(jsonErr)
			return err
		}

		_, err := s.Insert(&request)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}
		c.Status(http.StatusOK)
		return nil
	}
}
