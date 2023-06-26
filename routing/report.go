package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/masterreport"
	"ajebackend/model/transactionrequestreport"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func ReportRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	masterReportRepository := masterreport.NewRepository(db)
	masterReportService := masterreport.NewService(masterReportRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	allMasterRepository := allmaster.NewRepository(db)
	allMasterService := allmaster.NewService(allMasterRepository)

	transactionRequestReportRepository := transactionrequestreport.NewRepository(db)
	transactionRequestReportService := transactionrequestreport.NewService(transactionRequestReportRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	reportHandler := handler.NewMasterReportHandler(masterReportService, userIupopkService, validate, allMasterService, transactionRequestReportService, logService)

	reportRouting := app.Group("/report")

	reportRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	reportRouting.Get("/recap/:year/:iupopk_id", reportHandler.RecapDmo)
	reportRouting.Get("/realization/:year/:iupopk_id", reportHandler.RealizationReport)
	reportRouting.Get("/detail/:year/:iupopk_id", reportHandler.SaleDetailReport)

	reportRouting.Get("/download/recap/:year/:iupopk_id", reportHandler.DownloadRecapDmo)
	reportRouting.Get("/download/realization/:year/:iupopk_id", reportHandler.DownloadRealizationReport)
	reportRouting.Get("/download/detail/:year/:iupopk_id", reportHandler.DownloadSaleDetailReport)

	reportRouting.Post("/transactionrequest/preview/:iupopk_id", reportHandler.PreviewTransactionReport)
	reportRouting.Get("/transactionrequest/detail/:id/:iupopk_id", reportHandler.DetailTransactionReport)
	reportRouting.Get("/transactionrequest/list/:iupopk_id", reportHandler.ListTransactionReport)

	reportRouting.Post("/transactionrequest/create/:iupopk_id", reportHandler.CreateTransactionRequestReport)
	reportRouting.Put("/transactionrequest/create/document/:id/:iupopk_id", reportHandler.UpdateJobTransactionRequestReport)

	reportRouting.Delete("/transactionrequest/delete/:iupopk_id", reportHandler.DeleteTransactionReport)
	reportRouting.Delete("/transactionrequest/delete/:id/:iupopk_id", reportHandler.DeleteTransactionReportById)
}
