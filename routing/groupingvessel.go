package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/transaction"
	"ajebackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func GroupingVesselLnRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	groupingVesselLnRepository := groupingvesselln.NewRepository(db)
	groupingVesselLnService := groupingvesselln.NewService(groupingVesselLnRepository)

	groupingVesselHandler := handler.NewGroupingVesselHandler(transactionService, userService, historyService, validate, logService, groupingVesselLnService)

	groupingVesselRouting := app.Group("/groupingvessel")

	groupingVesselRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	groupingVesselRouting.Post("/create/ln", groupingVesselHandler.CreateGroupingVesselLn)
	groupingVesselRouting.Get("/detail/ln/:id", groupingVesselHandler.GetDetailGroupingVesselLn)
	groupingVesselRouting.Put("/update/ln/:id", groupingVesselHandler.EditGroupingVesselLn)
	groupingVesselRouting.Put("/update/document/ln/:id/:type", groupingVesselHandler.UploadDocumentGroupingVesselLn)
	groupingVesselRouting.Delete("/delete/ln/:id", groupingVesselHandler.DeleteGroupingVesselLn)
	groupingVesselRouting.Get("/list/ln", groupingVesselHandler.ListGroupingVesselLn)
}
