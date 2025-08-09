package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/fuelin"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func FuelInRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	FuelInRepository := fuelin.NewRepository(db)
	FuelInService := fuelin.NewService(FuelInRepository)

	fuelInHandler := handler.NewFuelInHandler(userService, FuelInService, validate)

	fuelInRouting := app.Group("/fuelin")

	fuelInRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	fuelInRouting.Post("/create", fuelInHandler.CreateFuelIn)

	fuelInRouting.Get("/list", fuelInHandler.GetFuelIn)
	fuelInRouting.Get("/list/pagination", fuelInHandler.GetListFuelIn)

	fuelInRouting.Get("/detail/:id", fuelInHandler.GetFuelInById)

	fuelInRouting.Put("/update/:id", fuelInHandler.UpdateFuelIn)

	fuelInRouting.Delete("/delete/:id", fuelInHandler.DeleteFuelIn)

}
