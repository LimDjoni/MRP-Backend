package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/adjuststock"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func AdjustStockRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	AdjustStockRepository := adjuststock.NewRepository(db)
	AdjustStockService := adjuststock.NewService(AdjustStockRepository)

	adjustStockHandler := handler.NewAdjustStockHandler(userService, AdjustStockService, validate)

	adjustStockRouting := app.Group("/adjuststock")

	adjustStockRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	adjustStockRouting.Post("/create", adjustStockHandler.CreateAdjustStock)

	adjustStockRouting.Get("/list", adjustStockHandler.GetAdjustStock)
	adjustStockRouting.Get("/list/pagination", adjustStockHandler.GetListAdjustStock)

	adjustStockRouting.Get("/detail/:id", adjustStockHandler.GetAdjustStockById)

	adjustStockRouting.Put("/update/:id", adjustStockHandler.UpdateAdjustStock)

	adjustStockRouting.Delete("/delete/:id", adjustStockHandler.DeleteAdjustStock)

}
