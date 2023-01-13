package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/destination"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/transaction"
	"ajebackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func TransactionRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	destinationRepository := destination.NewRepository(db)
	destinationService := destination.NewService(destinationRepository)

	transactionHandler := handler.NewTransactionHandler(transactionService, userService, historyService, validate, logService, destinationService)

	transactionRouting := app.Group("/transaction")

	transactionRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	transactionRouting.Post("/create/dn", transactionHandler.CreateTransactionDN)
	transactionRouting.Get("/list/dn", transactionHandler.ListDataDN)
	transactionRouting.Get("/detail/dn/:id", transactionHandler.DetailTransactionDN)
	transactionRouting.Delete("/delete/dn/:id", transactionHandler.DeleteTransactionDN)
	transactionRouting.Put("/update/dn/:id", transactionHandler.UpdateTransactionDN)
	transactionRouting.Put("/update/document/dn/:id/:type", transactionHandler.UpdateDocumentTransactionDN)
}
