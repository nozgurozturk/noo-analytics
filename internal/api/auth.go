package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/noo-analytics/entities"
	"github.com/nozgurozturk/noo-analytics/internal/config"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"github.com/nozgurozturk/noo-analytics/pkg/auth"
	"github.com/nozgurozturk/noo-analytics/pkg/user"
	"net/http"
	"time"
)

func AuthRouter(app fiber.Router, userService user.Service, authService auth.Service, config *config.ServerConfig) {
	app.Post("/login", login(userService, authService, config))
	app.Post("/signup", signUp(userService, authService, config))
}

func login(userService user.Service, authService auth.Service, config *config.ServerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var l entities.Login

		if err := c.BodyParser(&l); err != nil {
			parseError := errors.BadRequest("Invalid login body")
			c.Status(parseError.Status).JSON(parseError)
			return nil
		}

		if l.Email == "" {
			emailErr := errors.BadRequest("Email is required")
			c.Status(emailErr.Status).JSON(emailErr)
			return nil
		}
		if l.Password == "" {
			passErr := errors.BadRequest("Password is required")
			c.Status(passErr.Status).JSON(passErr)
			return nil
		}

		currentUser, err := userService.FindByEmail(l.Email)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}

		if passErr := user.VerifyPassword(currentUser.Password, l.Password); passErr != nil {
			c.Status(http.StatusForbidden).JSON(passErr.Error())
			return nil
		}

		tokens, err := authService.CreateToken(currentUser, config)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}

		accessToken := new(fiber.Cookie)
		accessToken.Name = "at"
		accessToken.Value = tokens.AccessToken.Token
		accessToken.Expires = time.Unix(tokens.AccessToken.Expires, 0)
		accessToken.HTTPOnly = true
		c.Cookie(accessToken)

		refreshToken := new(fiber.Cookie)
		refreshToken.Name = "rt"
		refreshToken.Value = tokens.RefreshToken.Token
		refreshToken.Expires = time.Unix(tokens.RefreshToken.Expires, 0)
		refreshToken.HTTPOnly = true
		c.Cookie(refreshToken)

		c.Status(http.StatusOK).JSON(&entities.UserResponse{
			Email: currentUser.Email,
			Name:  currentUser.Name,
		})
		return nil
	}
}

func signUp(userService user.Service, authService auth.Service, config *config.ServerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := new(entities.User)

		if err := c.BodyParser(&u); err != nil {
			parseError := errors.BadRequest("Invalid user body")
			c.Status(parseError.Status).JSON(parseError)
			return nil
		}

		exist, _ := userService.FindByEmail(u.Email)

		if exist != nil {
			existErr := errors.AlreadyExist("This email is taken by another user")
			c.Status(existErr.Status).JSON(existErr)
			return nil
		}

		tokens, err := authService.CreateToken(u, config)

		if err != nil {
			c.Status(err.Status).JSON(err)
			return nil
		}

		createdUser, err := userService.Create(u)

		if err != nil {
			if redisErr := authService.DeleteToken(tokens.AccessToken); redisErr != nil {
				c.Status(redisErr.Status).JSON(redisErr)
				return nil
			}
			c.Status(err.Status).JSON(err)
			return nil
		}

		accessToken := new(fiber.Cookie)
		accessToken.Name = "at"
		accessToken.Value = tokens.AccessToken.Token
		accessToken.Expires = time.Unix(tokens.AccessToken.Expires, 0)
		accessToken.HTTPOnly = true
		c.Cookie(accessToken)

		refreshToken := new(fiber.Cookie)
		refreshToken.Name = "rt"
		refreshToken.Value = tokens.RefreshToken.Token
		refreshToken.Expires = time.Unix(tokens.RefreshToken.Expires, 0)
		refreshToken.HTTPOnly = true
		c.Cookie(refreshToken)

		c.Status(http.StatusOK).JSON(&entities.UserResponse{
			Email: createdUser.Email,
			Name:  createdUser.Name,
		})
		return nil
	}
}
