package routing

import (
"ajebackend/handler"
"ajebackend/helper"
"ajebackend/model/logs"
"ajebackend/model/transaction"
"ajebackend/model/user"
"github.com/go-playground/validator/v10"
"github.com/gofiber/fiber/v2"
jwtware "github.com/gofiber/jwt/v3"
"gorm.io/gorm"
)

func ReportRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	reportHandler := handler.NewReportHandler(transactionService, userService,  validate, logService)

	reportRouting := app.Group("/report")

	reportRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err": err.Error(),
			})
		},
	}))

	reportRouting.Post("/", reportHandler.Report)
}
