package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/coareportln"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func CoaReportLnRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {

	coaReportLnRepository := coareportln.NewRepository(db)
	coaReportLnService := coareportln.NewService(coaReportLnRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	coaReportLnHandler := handler.NewCoaReportLnHandler(coaReportLnService, logService, userIupopkService, historyService, notificationUserService, validate)

	coaReportLnRouting := app.Group("/coareportln")

	coaReportLnRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	coaReportLnRouting.Get("/list/:iupopk_id", coaReportLnHandler.ListCoaReportLn)
	coaReportLnRouting.Post("/preview/:iupopk_id", coaReportLnHandler.ListCoaReportLnTransaction)
	coaReportLnRouting.Post("/create/:iupopk_id", coaReportLnHandler.CreateCoaReportLn)
	coaReportLnRouting.Delete("/delete/:id/:iupopk_id", coaReportLnHandler.DeleteCoaReportLn)
	coaReportLnRouting.Get("/detail/:id/:iupopk_id", coaReportLnHandler.DetailCoaReportLn)
	coaReportLnRouting.Post("/create/excel/:id/:iupopk_id", coaReportLnHandler.RequestCreateExcelCoaReportLn)
	// update job document
	coaReportLnRouting.Put("/update/document/:id/:iupopk_id", coaReportLnHandler.UpdateDocumentCoaReportLn)
}
