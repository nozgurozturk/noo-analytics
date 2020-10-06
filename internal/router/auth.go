package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/noo-analytics/internal/config"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"github.com/nozgurozturk/noo-analytics/pkg/auth"
	"net/http"
)

func AuthMiddleware(service auth.Service, config *config.ServerConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		at := ctx.Cookies("at")
		if at == "" {
			ctx.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Access token is not found"))
			return nil
		}

		_, errAt := service.ValidateToken(at, "Access", config)
		if errAt != nil {
			ctx.Status(errAt.Status).JSON(errAt)
			return nil
		}

		rt := ctx.Cookies("rt")
		if rt == "" {
			ctx.Status(http.StatusUnauthorized).JSON(errors.Unauthorized("Refresh token is not found"))
			return nil
		}

		_, errRt := service.ValidateToken(rt, "Refresh", config)
		if errRt != nil {

			ctx.Status(errRt.Status).JSON(errRt)
			return nil
		}
		return ctx.Next()
	}
}
