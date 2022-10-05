package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/company"
	"ajebackend/model/logs"
	"ajebackend/model/trader"
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func TraderRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	traderRepository := trader.NewRepository(db)
	traderService := trader.NewService(traderRepository)

	companyRepository := company.NewRepository(db)
	companyService := company.NewService(companyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	traderHandler := handler.NewTraderHandler(userService, traderService, companyService, logService, validate)

	traderRouting := app.Group("/trader")

	traderRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err": err.Error(),
			})
		},
	}))

	traderRouting.Get("/trader", traderHandler.ListTrader)
}
