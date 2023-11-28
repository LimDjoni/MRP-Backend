package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/jettybalance"
	"ajebackend/model/logs"
	"ajebackend/model/notificationuser"
	"ajebackend/model/pitloss"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func JettyBalanceRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {

	pitLossRepository := pitloss.NewRepository(db)
	pitLossService := pitloss.NewService(pitLossRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	jettyBalanceRepository := jettybalance.NewRepository(db)
	jettyBalanceService := jettybalance.NewService(jettyBalanceRepository)

	jettyBalanceHandler := handler.NewJettyBalanceHandler(pitLossService, historyService, validate, logService, notificationUserService, userIupopkService, jettyBalanceService)

	jettyBalanceRouting := app.Group("/jettybalance")

	jettyBalanceRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	jettyBalanceRouting.Get("/detail/:id/:iupopk_id", jettyBalanceHandler.DetailJettyBalance)
	jettyBalanceRouting.Get("/list/:iupopk_id", jettyBalanceHandler.ListJettyBalance)

	jettyBalanceRouting.Post("/create/:iupopk_id", jettyBalanceHandler.CreateJettyBalance)
	jettyBalanceRouting.Put("/update/:id/:iupopk_id", jettyBalanceHandler.UpdateJettyBalance)

	jettyBalanceRouting.Delete("/delete/:id/:iupopk_id", jettyBalanceHandler.DeleteJettyBalance)
}
