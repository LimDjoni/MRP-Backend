package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/counter"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/master/destination"
	"ajebackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func MasterRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	destinationRepository := destination.NewRepository(db)
	destinationService := destination.NewService(destinationRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	allMasterRepository := allmaster.NewRepository(db)
	allMasterService := allmaster.NewService(allMasterRepository)

	counterRepository := counter.NewRepository(db)
	counterService := counter.NewService(counterRepository)

	masterHandler := handler.NewMasterHandler(destinationService, userService, allMasterService, counterService, validate)

	masterRouting := app.Group("/master")

	masterRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	masterRouting.Get("/global", masterHandler.GetListMaster)

	masterRouting.Post("/create/iupopk", masterHandler.CreateIupopk)
	masterRouting.Put("/update/counter", masterHandler.UpdateCounter)
}
