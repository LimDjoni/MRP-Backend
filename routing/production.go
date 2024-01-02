package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/production"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func ProductionRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	productionRepository := production.NewRepository(db)
	productionService := production.NewService(productionRepository)

	logsRepository := logs.NewRepository(db)
	logsService := logs.NewService(logsRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	productionHandler := handler.NewProductionHandler(historyService, productionService, logsService, validate, userIupopkService)

	productionRouting := app.Group("/production")

	productionRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	productionRouting.Post("/create/:iupopk_id", productionHandler.CreateProduction)
	productionRouting.Put("/update/:id/:iupopk_id", productionHandler.UpdateProduction)
	productionRouting.Delete("/delete/:id/:iupopk_id", productionHandler.DeleteProduction)
	productionRouting.Get("/list/:iupopk_id", productionHandler.ListProduction)
	productionRouting.Get("/detail/:id/:iupopk_id", productionHandler.DetailProduction)
	productionRouting.Get("/summary/:iupopk_id", productionHandler.SummaryProduction)
}
