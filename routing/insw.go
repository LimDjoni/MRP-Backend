package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/insw"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/transaction"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func InswRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	groupingVesselLnRepository := groupingvesselln.NewRepository(db)
	groupingVesselLnService := groupingvesselln.NewService(groupingVesselLnRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	inswRepository := insw.NewRepository(db)
	inswService := insw.NewService(inswRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	inswHandler := handler.NewInswHandler(transactionService, historyService, validate, logService, groupingVesselLnService, notificationUserService, inswService, userIupopkService)

	inswRouting := app.Group("/insw")

	inswRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	inswRouting.Get("/list", inswHandler.ListInsw)
	inswRouting.Post("/create", inswHandler.CreateInsw)
	inswRouting.Post("/preview", inswHandler.ListGroupingVesselLnWithPeriod)
	inswRouting.Get("/detail/:id", inswHandler.DetailInsw)
	inswRouting.Delete("/delete/:id", inswHandler.DeleteInsw)
	inswRouting.Put("/update/document/:id", inswHandler.UpdateDocumentInsw)
	inswRouting.Post("/create/excel/:id", inswHandler.RequestCreateExcelInsw)
}
