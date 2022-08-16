package routing

import (
	"ajebackend/handler"
	"ajebackend/model/history"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	jwtware "github.com/gofiber/jwt/v3"
)

func TransactionRouting(db *gorm.DB, app fiber.Router) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	transactionHandler := handler.NewTransactionHandler(transactionService, userService, historyService)

	transactionRouting := app.Group("/transaction") // /api

	// Reference to edit the error - https://www.youtube.com/watch?v=ejEizICXm9w
	transactionRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte("aFhF234aiI"),
		SigningMethod: jwtware.HS256,
	}))

	transactionRouting.Post("/create/dn", transactionHandler.CreateTransactionDN)
	transactionRouting.Get("/list/dn", transactionHandler.ListDataDN)
	transactionRouting.Get("/detail/dn/:id", transactionHandler.DetailTransactionDN)
	transactionRouting.Delete("/delete/dn/:id", transactionHandler.DeleteTransaction)
}
