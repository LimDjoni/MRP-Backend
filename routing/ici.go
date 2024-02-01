package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/ici"
	"ajebackend/model/logs"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func IciRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	IciRepository := ici.NewRepository(db)
	IciService := ici.NewService(IciRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	logsRepository := logs.NewRepository(db)
	logsService := logs.NewService(logsRepository)

	iciHandler := handler.NewIciHandler(IciService, logsService, validate, userIupopkService)

	iciRouting := app.Group("/ici")

	iciRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	iciRouting.Post("/create/:iupopk_id", iciHandler.CreateIci)

}
