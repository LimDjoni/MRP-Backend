package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/trader"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func DmoRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	dmoRepository := dmo.NewRepository(db)
	dmoService := dmo.NewService(dmoRepository)

	traderRepository := trader.NewRepository(db)
	traderService := trader.NewService(traderRepository)

	traderDmoRepository := traderdmo.NewRepository(db)
	traderDmoService := traderdmo.NewService(traderDmoRepository)

	dmoHandler := handler.NewDmoHandler(transactionService, userService, historyService, logService, dmoService, traderService, traderDmoService, validate)

	dmoRouting := app.Group("/dmo")

	dmoRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err": err.Error(),
			})
		},
	}))

	dmoRouting.Post("/create", dmoHandler.CreateDmo)
	dmoRouting.Get("/list", dmoHandler.ListDmo)
	dmoRouting.Get("/list/transaction", dmoHandler.ListDataDNWithoutDmo)
	dmoRouting.Get("/detail/:id", dmoHandler.DetailDmo)
	dmoRouting.Delete("/delete/:id", dmoHandler.DeleteDmo)
	dmoRouting.Put("/update/document/:id", dmoHandler.UpdateDocumentDmo)
	dmoRouting.Put("/update/document/downloaded/:id/:type", dmoHandler.UpdateIsDownloadedDocumentDmo)
	dmoRouting.Put("/update/document/signed/:id/:type", dmoHandler.UpdateIsSignedDocumentDmo)
}
