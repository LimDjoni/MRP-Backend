package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/logs"
	"ajebackend/model/master/company"
	"ajebackend/model/master/trader"
	"ajebackend/model/traderdmo"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func TraderRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	traderRepository := trader.NewRepository(db)
	traderService := trader.NewService(traderRepository)

	companyRepository := company.NewRepository(db)
	companyService := company.NewService(companyRepository)

	traderDmoRepository := traderdmo.NewRepository(db)
	traderDmoService := traderdmo.NewService(traderDmoRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	traderHandler := handler.NewTraderHandler(traderService, companyService, traderDmoService, logService, validate, userIupopkService)

	traderRouting := app.Group("/trader")

	traderRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	traderRouting.Get("/list", traderHandler.ListTrader)
	traderRouting.Post("/create", traderHandler.CreateTrader)
	traderRouting.Put("/update/:id", traderHandler.UpdateTrader)
	traderRouting.Delete("/delete/:id", traderHandler.DeleteTrader)
	traderRouting.Get("/detail/:id", traderHandler.DetailTrader)
}
