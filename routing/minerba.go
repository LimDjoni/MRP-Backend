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

	minerbaHandler := handler.NewMinerbaHandler(transactionService, userService, historyService, logService, minerbaService, notificationUserService, validate)

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
	minerbaRouting.Get("/list/dn", minerbaHandler.ListMinerba)
	minerbaRouting.Get("/list/transaction/dn", minerbaHandler.ListDataDNWithoutMinerba)
	minerbaRouting.Get("/detail/dn/:id", minerbaHandler.DetailMinerba)
	minerbaRouting.Post("/create/dn", minerbaHandler.CreateMinerba)
	minerbaRouting.Delete("/delete/dn/:id", minerbaHandler.DeleteMinerba)
	minerbaRouting.Put("/update/dn/:id", minerbaHandler.UpdateMinerba)
	minerbaRouting.Put("/update/document/dn/:id", minerbaHandler.UpdateDocumentMinerba)
	minerbaRouting.Post("/create/excel/dn/:id", minerbaHandler.RequestCreateExcelMinerba)
	minerbaRouting.Post("/check/dn", minerbaHandler.CheckValidPeriodMinerba)

	// LN
	minerbaRouting.Get("/list/ln", minerbaLnHandler.ListMinerbaLn)
	minerbaRouting.Get("/list/transaction/ln", minerbaLnHandler.ListDataLNWithoutMinerba)
	minerbaRouting.Get("/detail/ln/:id", minerbaLnHandler.DetailMinerbaLn)
	minerbaRouting.Post("/create/ln", minerbaLnHandler.CreateMinerbaLn)
	minerbaRouting.Delete("/delete/ln/:id", minerbaLnHandler.DeleteMinerbaLn)
	minerbaRouting.Put("/update/ln/:id", minerbaLnHandler.UpdateMinerbaLn)
	minerbaRouting.Put("/update/document/ln/:id", minerbaLnHandler.UpdateDocumentMinerbaLn)
	minerbaRouting.Post("/create/excel/ln/:id", minerbaLnHandler.RequestCreateExcelMinerbaLn)
	minerbaRouting.Post("/check/ln", minerbaLnHandler.CheckValidPeriodMinerbaLn)
}
