package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/masterreport"
	"ajebackend/model/notificationuser"
	"ajebackend/model/royaltyrecon"
	"ajebackend/model/royaltyreport"
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

	royaltyReconRepository := royaltyrecon.NewRepository(db)
	royaltyReconService := royaltyrecon.NewService(royaltyReconRepository)

	royaltyReportRepository := royaltyreport.NewRepository(db)
	royaltyReportService := royaltyreport.NewService(royaltyReportRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	reportHandler := handler.NewMasterReportHandler(masterReportService, userIupopkService, validate, allMasterService, transactionRequestReportService, logService, royaltyReconService, royaltyReportService, historyService, notificationUserService)

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

	// Royalty Recon
	reportRouting.Post("/royalty/recon/preview/:iupopk_id", reportHandler.PreviewRoyaltyRecon)
	reportRouting.Post("/royalty/recon/create/:iupopk_id", reportHandler.CreateRoyaltyRecon)
	reportRouting.Delete("/royalty/recon/delete/:id/:iupopk_id", reportHandler.DeleteRoyaltyRecon)
	reportRouting.Get("/royalty/recon/detail/:id/:iupopk_id", reportHandler.DetailRoyaltyRecon)
	reportRouting.Post("/royalty/recon/create/excel/:id/:iupopk_id", reportHandler.RequestCreateExcelRoyaltyRecon)
	reportRouting.Get("/royalty/recon/list/:iupopk_id", reportHandler.ListRoyaltyRecon)
	reportRouting.Put("/royalty/recon/update/document/:id/:iupopk_id", reportHandler.UpdateDocumentRoyaltyRecon)

	// Royalty Report
	reportRouting.Post("/royalty/report/preview/:iupopk_id", reportHandler.PreviewRoyaltyReport)
	reportRouting.Post("/royalty/report/create/:iupopk_id", reportHandler.CreateRoyaltyReport)
	reportRouting.Delete("/royalty/report/delete/:id/:iupopk_id", reportHandler.DeleteRoyaltyReport)
	reportRouting.Get("/royalty/report/detail/:id/:iupopk_id", reportHandler.DetailRoyaltyReport)
	reportRouting.Post("/royalty/report/create/excel/:id/:iupopk_id", reportHandler.RequestCreateExcelRoyaltyReport)
	reportRouting.Get("/royalty/report/list/:iupopk_id", reportHandler.ListRoyaltyReport)
	reportRouting.Put("/royalty/report/update/document/:id/:iupopk_id", reportHandler.UpdateDocumentRoyaltyReport)
}
