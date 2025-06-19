package http

import (
	"github/smile-ko/go-template/config"
	fiberswaggerdocsV1 "github/smile-ko/go-template/docs/v1"
	fiberswaggerdocsV2 "github/smile-ko/go-template/docs/v2"
	"github/smile-ko/go-template/internal/interfaces/http/middleware"
	v1 "github/smile-ko/go-template/internal/interfaces/http/v1"
	v2 "github/smile-ko/go-template/internal/interfaces/http/v2"
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func NewRouter(app *fiber.App, cfg *config.Config, pg *postgres.Postgres, l logger.Interface) {
	// Options
	if cfg.App.EnvName == "dev" {
		// cors
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowCredentials: false,
		}))

		// logger dev
		app.Use(fiberlogger.New(fiberlogger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}))
	} else {
		// logger production
		app.Use(middleware.Logger(l))
		app.Use(middleware.Recovery(l))
	}

	// Swagger docs
	if cfg.Swagger.Enabled {
		setupSwagger(app, cfg)
	}

	apiV1 := app.Group("/api/v1")
	{
		v1.RegisterRoutes(apiV1, pg, l)
	}

	apiV2 := app.Group("/api/v2")
	{
		v2.RegisterRoutes(apiV2, pg, l)
	}
}

func setupSwagger(app *fiber.App, cfg *config.Config) {
	// Swagger V1
	fiberswaggerdocsV1.SwaggerInfov1.Title = cfg.App.Name + " v1"
	fiberswaggerdocsV1.SwaggerInfov1.Version = "v1"
	fiberswaggerdocsV1.SwaggerInfov1.BasePath = "/api/v1"
	app.Get("/v1/swagger/*", swagger.New(swagger.Config{
		InstanceName: "v1",
	}))

	// Swagger V2
	fiberswaggerdocsV2.SwaggerInfov2.Title = cfg.App.Name + " v2"
	fiberswaggerdocsV2.SwaggerInfov2.Version = "v2"
	fiberswaggerdocsV2.SwaggerInfov2.BasePath = "/api/v2"
	app.Get("/v2/swagger/*", swagger.New(swagger.Config{
		InstanceName: "v2",
	}))
}
