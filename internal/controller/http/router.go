package http

import (
	fiberswaggerdocsV1 "github/smile-ko/go-template/api/openapi/v1"
	"github/smile-ko/go-template/config"
	"github/smile-ko/go-template/internal/controller/http/middleware"
	v1 "github/smile-ko/go-template/internal/controller/http/v1"
	"github/smile-ko/go-template/pkg/logger"
	"github/smile-ko/go-template/pkg/postgres"
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/swagger"
)

func NewRouter(app *fiber.App, cfg *config.Config, pg *postgres.Postgres, l logger.ILogger) {
	if cfg.App.EnvName == "dev" {
		// cors
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "*",
			AllowCredentials: false,
		}))
	}

	// option
	app.Use(helmet.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Swagger
	if cfg.Swagger.Enabled {
		setupSwagger(app, cfg)
	}

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New(cfg.App.Name)
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// logger
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	apiV1 := app.Group("/api/v1")
	{
		v1.NewV1Router(apiV1, pg, l)
	}

	// handler error
	app.Use(middleware.NotFound())
}

func setupSwagger(app *fiber.App, cfg *config.Config) {
	// Swagger V1
	fiberswaggerdocsV1.SwaggerInfov1.Title = cfg.App.Name + " v1"
	fiberswaggerdocsV1.SwaggerInfov1.Version = "v1"
	fiberswaggerdocsV1.SwaggerInfov1.BasePath = "/api/v1"
	app.Get("/swagger/v1/*", swagger.New(swagger.Config{
		InstanceName: "v1",
	}))
}
