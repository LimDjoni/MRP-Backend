package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/coareport"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func CoaReportRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {

	coaReportRepository := coareport.NewRepository(db)
	coaReportService := coareport.NewService(coaReportRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	coaReportHandler := handler.NewCoaReportHandler(coaReportService, logService, userIupopkService, historyService, notificationUserService, validate)

	coaReportRouting := app.Group("/coareport")

	coaReportRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	coaReportRouting.Get("/list/:iupopk_id", coaReportHandler.ListCoaReport)
	coaReportRouting.Post("/preview/:iupopk_id", coaReportHandler.ListCoaReportTransaction)
	coaReportRouting.Post("/create/:iupopk_id", coaReportHandler.CreateCoaReport)
	coaReportRouting.Delete("/delete/:id/:iupopk_id", coaReportHandler.DeleteCoaReport)
	coaReportRouting.Get("/detail/:id/:iupopk_id", coaReportHandler.DetailCoaReport)
	coaReportRouting.Post("/create/excel/:id/:iupopk_id", coaReportHandler.RequestCreateExcelCoaReport)
	// update job document
	coaReportRouting.Put("/update/document/:id/:iupopk_id", coaReportHandler.UpdateDocumentCoaReport)
}
