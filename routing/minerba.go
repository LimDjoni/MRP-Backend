package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/notification"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func MinerbaRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	minerbaRepository := minerba.NewRepository(db)
	minerbaService := minerba.NewService(minerbaRepository)

	notificationRepository := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepository)

	minerbaHandler := handler.NewMinerbaHandler(transactionService, userService, historyService, logService, minerbaService, notificationService, validate)

	minerbaRouting := app.Group("/minerba")

	minerbaRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err": err.Error(),
			})
		},
	}))

	minerbaRouting.Get("/list", minerbaHandler.ListMinerba)
	minerbaRouting.Get("/list/transaction", minerbaHandler.ListDataDNWithoutMinerba)
	minerbaRouting.Get("/detail/:id", minerbaHandler.DetailMinerba)
	minerbaRouting.Post("/create", minerbaHandler.CreateMinerba)
	minerbaRouting.Delete("/delete/:id", minerbaHandler.DeleteMinerba)
	minerbaRouting.Put("update/:id", minerbaHandler.UpdateMinerba)
	minerbaRouting.Put("/update/document/:id", minerbaHandler.UpdateDocumentMinerba)
	minerbaRouting.Post("/create/excel/:id", minerbaHandler.RequestCreateExcelMinerba)
}
