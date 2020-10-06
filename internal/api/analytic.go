package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	analytics "github.com/nozgurozturk/noo-analytics/pkg/analytic"
	"net/http"
)

func AnalyticsRouter(app fiber.Router, service analytics.Service) {
	app.Post("/action/date", findActionsByDate(service))
}

func findActionsByDate(s analytics.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request entities.AnalyticsActionRequest
		if err := c.BodyParser(&request); err != nil {
			jsonErr := errors.BadRequest("Invalid Trace Object")
			c.Status(jsonErr.Status).JSON(jsonErr)
			return err
		}

		response, err := s.FindActionsByDate(&request)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}
		c.Status(http.StatusOK).JSON(response)
		return nil

	}
}
