package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/stockfuel"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func StockFuelRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	StockFuelRepository := stockfuel.NewRepository(db)
	StockFuelService := stockfuel.NewService(StockFuelRepository)

	stockfuelsHandler := handler.NewStockFuelHandler(userService, StockFuelService, validate)

	stockfuelsRouting := app.Group("/stockfuel")

	stockfuelsRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	stockfuelsRouting.Get("/list/summary", stockfuelsHandler.ListStockFuel)

}
