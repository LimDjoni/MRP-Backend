package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/alatberat"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func AlatBeratRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	AlatBeratRepository := alatberat.NewRepository(db)
	AlatBeratService := alatberat.NewService(AlatBeratRepository)

	alatBeratHandler := handler.NewAlatBeratHandler(userService, AlatBeratService, validate)

	alatBeratRouting := app.Group("/alatberat")

	alatBeratRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	alatBeratRouting.Post("/create/alatBerat", alatBeratHandler.CreateAlatBerat)

	alatBeratRouting.Get("/list/alatBerat", alatBeratHandler.GetAlatBerat)
	alatBeratRouting.Get("/list/pagination", alatBeratHandler.GetListAlatBerat)

	alatBeratRouting.Get("/detail/alatBerat/:id", alatBeratHandler.GetAlatBeratById)

	alatBeratRouting.Get("/consumption/:brandId/:heavyEquipmentId/:seriesId", alatBeratHandler.GetConsumption)

	alatBeratRouting.Put("/update/alatBerat/:id", alatBeratHandler.UpdateAlatBerat)

	alatBeratRouting.Delete("/delete/alatBerat/:id", alatBeratHandler.DeleteAlatBerat)

}
