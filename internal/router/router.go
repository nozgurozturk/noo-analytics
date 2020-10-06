package router

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/noo-analytics/internal/api"
	"github.com/nozgurozturk/noo-analytics/internal/config"
	analytics "github.com/nozgurozturk/noo-analytics/pkg/analytic"
	"github.com/nozgurozturk/noo-analytics/pkg/auth"
	"github.com/nozgurozturk/noo-analytics/pkg/trace"
	"github.com/nozgurozturk/noo-analytics/pkg/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Router struct {
	App     *fiber.App
	MongoDB *mongo.Database
	Redis   *redis.Client
	Config  *config.ServerConfig
}

func New(r *Router) *Router {
	router := &Router{
		App:     r.App,
		MongoDB: r.MongoDB,
		Redis:   r.Redis,
		Config:  r.Config,
	}
	router.initializeRouters()
	return router
}

func (r *Router) initializeRouters() {

	traceCollection := r.MongoDB.Collection("traces")
	userCollection := r.MongoDB.Collection("users")

	// User Router
	authRepository := auth.NewRepository(r.Redis)
	authService := auth.NewService(authRepository)
	userRepository := user.NewRepository(userCollection)
	userService := user.NewService(userRepository)
	userRoute := r.App.Group("/user", AuthMiddleware(authService, r.Config))
	authRoute := r.App.Group("/auth")
	api.UserRouter(userRoute, userService)
	api.AuthRouter(authRoute, userService, authService, r.Config)


	// Trace Router
	traceRepository := trace.NewRepository(traceCollection)
	traceService := trace.NewService(traceRepository)
	traceRoute := r.App.Group("/trace")
	api.TraceRooter(traceRoute, traceService)

	// Analytics Router
	analyticsRepository := analytics.NewRepository(traceCollection)
	analyticsService := analytics.NewService(analyticsRepository)
	analyticsRoute := r.App.Group("/analytics", AuthMiddleware(authService, r.Config))
	api.AnalyticsRouter(analyticsRoute, analyticsService)


}
