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
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"gorm.io/gorm"
)

func HaulingSynchronizeRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	haulingSynchronizeRepository := haulingsynchronize.NewRepository(db)
	haulingSynchronizeService := haulingsynchronize.NewService(haulingSynchronizeRepository)

	transactionsHaulingRepository := transactionshauling.NewRepository(db)
	transactionsHaulingService := transactionshauling.NewService(transactionsHaulingRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	haulingSynchronizeRoutingHandler := handler.NewHaulingHandler(haulingSynchronizeService, transactionsHaulingService, userIupopkService, logService)

	syncRouting := app.Group("/synchronize")

	syncRouting.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			helper.GetEnvWithKey("USERNAME_BASIC"): helper.GetEnvWithKey("PASSWORD_BASIC"),
		},
	}))

	syncRouting.Post("/isp", haulingSynchronizeRoutingHandler.SyncHaulingDataIsp)
	syncRouting.Post("/jetty", haulingSynchronizeRoutingHandler.SyncHaulingDataJetty)
}
