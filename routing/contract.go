package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/contract"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func ContractRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	contractRepository := contract.NewRepository(db)
	contractService := contract.NewService(contractRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	allMasterRepository := allmaster.NewRepository(db)
	allMasterService := allmaster.NewService(allMasterRepository)

	contractHandler := handler.NewContractHandler(contractService, historyService, logService, validate, userIupopkService, allMasterService)

	contractRouting := app.Group("/contract")

	contractRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	contractRouting.Post("/create/:iupopk_id", contractHandler.CreateContract)

	contractRouting.Put("/update/:id/:iupopk_id", contractHandler.UpdateContract)

	contractRouting.Get("/list/:iupopk_id", contractHandler.ListContract)
	contractRouting.Get("/detail/:id/:iupopk_id", contractHandler.DetailContract)
}
