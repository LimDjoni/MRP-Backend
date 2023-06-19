package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/masterreport"
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

	reportHandler := handler.NewMasterReportHandler(masterReportService, userIupopkService, validate, allMasterService)

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

	reportRouting.Get("/alltransaction/:type/:iupopk_id", reportHandler.GetTransactionReport)
}
