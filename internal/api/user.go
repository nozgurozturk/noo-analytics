package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"github.com/nozgurozturk/noo-analytics/pkg/user"
	"net/http"
)

func UserRouter(app fiber.Router, service user.Service) {
	app.Post("/", createUser(service))
	app.Post("/find", findUserByEmail(service))
}

func createUser(s user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request entities.User
		if err := c.BodyParser(&request); err != nil {
			jsonErr := errors.BadRequest("Invalid user body")
			c.Status(jsonErr.Status).JSON(jsonErr)
			return nil
		}

		exist, _ := s.FindByEmail(request.Email)

		if exist != nil {
			existErr := errors.AlreadyExist("This email is taken by another user")
			c.Status(existErr.Status).JSON(existErr)
			return nil
		}

		u, err := s.Create(&request)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}
		c.Status(http.StatusCreated).JSON(u)
		return nil
	}
}

func findUserByEmail(s user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&request); err != nil {
			jsonErr := errors.BadRequest("Invalid Email Object")
			c.Status(jsonErr.Status).JSON(jsonErr)
			return nil
		}
		found, err := s.FindByEmail(request.Email)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}
		if found == nil {
			notFoundErr := errors.NotFound("User is not found")
			c.Status(notFoundErr.Status).JSON(notFoundErr)
			return nil
		}
		c.Status(http.StatusOK).JSON(found)
		return nil
	}
}
