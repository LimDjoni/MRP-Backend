package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/notificationuser"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func MinerbaRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	minerbaRepository := minerba.NewRepository(db)
	minerbaService := minerba.NewService(minerbaRepository)

	minerbaLnRepository := minerbaln.NewRepository(db)
	minerbaLnService := minerbaln.NewService(minerbaLnRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	minerbaHandler := handler.NewMinerbaHandler(transactionService, userService, historyService, logService, minerbaService, notificationUserService, validate, userIupopkService)

	minerbaLnHandler := handler.NewMinerbaLnHandler(transactionService, userService, historyService, logService, minerbaLnService, notificationUserService, validate)

	minerbaRouting := app.Group("/minerba")

	minerbaRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	// DN
	minerbaRouting.Get("/list/dn/:iupopk_id", minerbaHandler.ListMinerba)
	minerbaRouting.Get("/list/transaction/dn/:iupopk_id", minerbaHandler.ListDataDNWithoutMinerba)
	minerbaRouting.Get("/detail/dn/:id/:iupopk_id", minerbaHandler.DetailMinerba)
	minerbaRouting.Post("/create/dn/:iupopk_id", minerbaHandler.CreateMinerba)
	minerbaRouting.Delete("/delete/dn/:id/:iupopk_id", minerbaHandler.DeleteMinerba)
	minerbaRouting.Put("/update/dn/:id/:iupopk_id", minerbaHandler.UpdateMinerba)
	minerbaRouting.Put("/update/document/dn/:id/:iupopk_id", minerbaHandler.UpdateDocumentMinerba)
	minerbaRouting.Post("/create/excel/dn/:id/:iupopk_id", minerbaHandler.RequestCreateExcelMinerba)
	minerbaRouting.Post("/check/dn/:iupopk_id", minerbaHandler.CheckValidPeriodMinerba)

	// LN
	minerbaRouting.Get("/list/ln/:iupopk_id", minerbaLnHandler.ListMinerbaLn)
	minerbaRouting.Get("/list/transaction/ln/:iupopk_id", minerbaLnHandler.ListDataLNWithoutMinerba)
	minerbaRouting.Get("/detail/ln/:id/:iupopk_id", minerbaLnHandler.DetailMinerbaLn)
	minerbaRouting.Post("/create/ln/:iupopk_id", minerbaLnHandler.CreateMinerbaLn)
	minerbaRouting.Delete("/delete/ln/:id/:iupopk_id", minerbaLnHandler.DeleteMinerbaLn)
	minerbaRouting.Put("/update/ln/:id/:iupopk_id", minerbaLnHandler.UpdateMinerbaLn)
	minerbaRouting.Put("/update/document/ln/:id/:iupopk_id", minerbaLnHandler.UpdateDocumentMinerbaLn)
	minerbaRouting.Post("/create/excel/ln/:id/:iupopk_id", minerbaLnHandler.RequestCreateExcelMinerbaLn)
	minerbaRouting.Post("/check/ln/:iupopk_id", minerbaLnHandler.CheckValidPeriodMinerbaLn)
}
