package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/destination"
	"ajebackend/model/transaction"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func TransactionRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	destinationRepository := destination.NewRepository(db)
	destinationService := destination.NewService(destinationRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	transactionHandler := handler.NewTransactionHandler(transactionService, historyService, validate, logService, destinationService, userIupopkService)

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

	transactionRouting.Post("/create/dn/:iupopk_id", transactionHandler.CreateTransactionDN)
	transactionRouting.Post("/create/ln/:iupopk_id", transactionHandler.CreateTransactionLN)

	transactionRouting.Put("/update/dn/:id/:iupopk_id", transactionHandler.UpdateTransactionDN)
	transactionRouting.Put("/update/ln/:id/:iupopk_id", transactionHandler.UpdateTransactionLN)

	transactionRouting.Get("/list/:transaction_type/:iupopk_id", transactionHandler.ListData)
	transactionRouting.Get("/detail/:transaction_type/:id/:iupopk_id", transactionHandler.DetailTransaction)
	transactionRouting.Delete("/delete/:transaction_type/:id/:iupopk_id", transactionHandler.DeleteTransaction)
	transactionRouting.Put("/update/document/:transaction_type/:id/:type/:iupopk_id", transactionHandler.UpdateDocumentTransaction)
}
