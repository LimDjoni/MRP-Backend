package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/fuelratio"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func FuelRatioRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	FuelRatioRepository := fuelratio.NewRepository(db)
	FuelRatioService := fuelratio.NewService(FuelRatioRepository)

	fuelratiosHandler := handler.NewFuelRatioHandler(userService, FuelRatioService, validate)

	fuelratiosRouting := app.Group("/fuelratio")

	fuelratiosRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	fuelratiosRouting.Post("/create", fuelratiosHandler.CreateFuelRatio)

	fuelratiosRouting.Get("/list", fuelratiosHandler.GetFuelRatio)
	fuelratiosRouting.Get("/list/pagination", fuelratiosHandler.GetListFuelRatio)
	fuelratiosRouting.Get("/list/summary", fuelratiosHandler.GetFindFuelRatioExport)
	fuelratiosRouting.Get("/list/summary/pagination", fuelratiosHandler.GetListRangkuman)

	fuelratiosRouting.Get("/detail/:id", fuelratiosHandler.GetFuelRatioById)

	fuelratiosRouting.Put("/update/:id", fuelratiosHandler.UpdateFuelRatio)

	fuelratiosRouting.Delete("/delete/:id", fuelratiosHandler.DeleteFuelRatio)

}
