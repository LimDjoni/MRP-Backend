package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/company"
	"ajebackend/model/master/trader"
	"ajebackend/model/notificationuser"
	"ajebackend/model/reportdmo"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func DmoRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

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

	companyRepository := company.NewRepository(db)
	companyService := company.NewService(companyRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	dmoVesselRepository := dmovessel.NewRepository(db)
	dmoVesselService := dmovessel.NewService(dmoVesselRepository)

	reportDmoRepository := reportdmo.NewRepository(db)
	reportDmoService := reportdmo.NewService(reportDmoRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	dmoHandler := handler.NewDmoHandler(transactionService, historyService, logService, dmoService, traderService, traderDmoService, notificationUserService, companyService, validate, dmoVesselService, userIupopkService)

	reportDmoHandler := handler.NewReportDmoHandler(transactionService, historyService, logService, notificationUserService, validate, reportDmoService, userIupopkService)

	dmoRouting := app.Group("/dmo")

	dmoRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	dmoRouting.Post("/create/:iupopk_id", dmoHandler.CreateDmo)
	dmoRouting.Get("/list/:iupopk_id", dmoHandler.ListDmo)
	dmoRouting.Get("/list/transaction/:iupopk_id", dmoHandler.ListTransactionForDmo)
	dmoRouting.Get("/detail/:id/:iupopk_id", dmoHandler.DetailDmo)
	dmoRouting.Delete("/delete/:id/:iupopk_id", dmoHandler.DeleteDmo)
	dmoRouting.Put("/update/document/:id/:iupopk_id", dmoHandler.UpdateDocumentDmo)
	dmoRouting.Put("/update/document/downloaded/:id/:type/:iupopk_id", dmoHandler.UpdateIsDownloadedDocumentDmo)
	dmoRouting.Put("/update/document/signed/:id/:type/:iupopk_id", dmoHandler.UpdateTrueIsSignedDmoDocument)
	dmoRouting.Put("/update/document/not_signed/:id/:type/:iupopk_id", dmoHandler.UpdateFalseIsSignedDmoDocument)
	dmoRouting.Put("/update/:id/:iupopk_id", dmoHandler.UpdateDmo)

	// Report
	dmoRouting.Post("/create/report/:iupopk_id", reportDmoHandler.CreateReportDmo)
	dmoRouting.Put("/update/report/document/:id/:iupopk_id", reportDmoHandler.UpdateDocumentReportDmo)
	dmoRouting.Put("/update/report/:id/:iupopk_id", reportDmoHandler.UpdateReportDmo)
	dmoRouting.Delete("/delete/report/:id/:iupopk_id", reportDmoHandler.DeleteReportDmo)
	dmoRouting.Post("/create/excel/:id/:iupopk_id", reportDmoHandler.RequestCreateExcelReportDmo)
	dmoRouting.Get("/list/report/transaction/:iupopk_id", reportDmoHandler.GetListForReport)
	dmoRouting.Post("/validate/report/:iupopk_id", reportDmoHandler.CheckValidPeriodReportDmo)
	dmoRouting.Get("/detail/report/:id/:iupopk_id", reportDmoHandler.DetailReportDmo)
	dmoRouting.Get("/list/report/:iupopk_id", reportDmoHandler.ListReportDmo)
}
