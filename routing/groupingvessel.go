package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/destination"
	"ajebackend/model/groupingvesseldn"
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

	groupingVesselLnHandler := handler.NewGroupingVesselLnHandler(transactionService, userService, historyService, validate, logService, groupingVesselLnService)

	groupingVesselDnRepository := groupingvesseldn.NewRepository(db)
	groupingVesselDnService := groupingvesseldn.NewService(groupingVesselDnRepository)

	destinationRepository := destination.NewRepository(db)
	destinationService := destination.NewService(destinationRepository)

	groupingVesselDnHandler := handler.NewGroupingVesselDnHandler(transactionService, userService, historyService, validate, logService, groupingVesselDnService, destinationService)

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

	// LN
	groupingVesselRouting.Post("/create/ln", groupingVesselLnHandler.CreateGroupingVesselLn)
	groupingVesselRouting.Get("/detail/ln/:id", groupingVesselLnHandler.GetDetailGroupingVesselLn)
	groupingVesselRouting.Put("/update/ln/:id", groupingVesselLnHandler.EditGroupingVesselLn)
	groupingVesselRouting.Put("/update/document/ln/:id/:type", groupingVesselLnHandler.UploadDocumentGroupingVesselLn)
	groupingVesselRouting.Delete("/delete/ln/:id", groupingVesselLnHandler.DeleteGroupingVesselLn)
	groupingVesselRouting.Get("/list/ln", groupingVesselLnHandler.ListGroupingVesselLn)

	// DN
	groupingVesselRouting.Get("/list/dn", groupingVesselDnHandler.ListGroupingVesselDn)
	groupingVesselRouting.Post("/create/dn", groupingVesselDnHandler.CreateGroupingVesselDn)
	groupingVesselRouting.Put("/update/dn/:id", groupingVesselDnHandler.EditGroupingVesselDn)
	groupingVesselRouting.Put("/update/document/dn/:id/:type", groupingVesselDnHandler.UploadDocumentGroupingVesselDn)
	groupingVesselRouting.Delete("/delete/dn/:id", groupingVesselDnHandler.DeleteGroupingVesselDn)
	groupingVesselRouting.Get("/detail/dn/:id", groupingVesselDnHandler.GetDetailGroupingVesselDn)
}
