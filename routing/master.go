package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/trader"
	"ajebackend/model/user"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func MasterRouting(db *gorm.DB, app fiber.Router) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	traderRepository := trader.NewRepository(db)
	traderService := trader.NewService(traderRepository)

	masterHandler := handler.NewMasterHandler(userService, traderService)

	masterRouting := app.Group("/master")

	masterRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err": err.Error(),
			})
		},
	}))

	masterRouting.Get("/trader", masterHandler.ListTrader)
}
