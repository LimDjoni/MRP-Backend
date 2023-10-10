package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"

	"ajebackend/model/haulingsynchronize"
	"ajebackend/model/logs"
	"ajebackend/model/transactionshauling"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func HaulingTransactionRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	haulingSynchronizeRepository := haulingsynchronize.NewRepository(db)
	haulingSynchronizeService := haulingsynchronize.NewService(haulingSynchronizeRepository)

	transactionsHaulingRepository := transactionshauling.NewRepository(db)
	transactionsHaulingService := transactionshauling.NewService(transactionsHaulingRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	haulingTransactionRoutingHandler := handler.NewHaulingHandler(haulingSynchronizeService, transactionsHaulingService, userIupopkService, logService)

	haulingRouting := app.Group("/hauling")

	haulingRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	haulingRouting.Get("/list/stock/:iupopk_id", haulingTransactionRoutingHandler.ListStockRom)
	haulingRouting.Get("/list/transaction/:iupopk_id", haulingTransactionRoutingHandler.ListTransactionHauling)
	haulingRouting.Get("/detail/stock/:id/:iupopk_id", haulingTransactionRoutingHandler.DetailStockRom)
	haulingRouting.Get("/detail/transaction/:id/:iupopk_id", haulingTransactionRoutingHandler.DetailTransactionHauling)
}
