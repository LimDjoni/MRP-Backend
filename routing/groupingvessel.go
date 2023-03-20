package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/destination"
	"ajebackend/model/transaction"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func GroupingVesselLnRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	transactionRepository := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	groupingVesselLnRepository := groupingvesselln.NewRepository(db)
	groupingVesselLnService := groupingvesselln.NewService(groupingVesselLnRepository)

	groupingVesselDnRepository := groupingvesseldn.NewRepository(db)
	groupingVesselDnService := groupingvesseldn.NewService(groupingVesselDnRepository)

	destinationRepository := destination.NewRepository(db)
	destinationService := destination.NewService(destinationRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	groupingVesselDnHandler := handler.NewGroupingVesselDnHandler(transactionService, historyService, validate, logService, groupingVesselDnService, destinationService, userIupopkService)

	groupingVesselLnHandler := handler.NewGroupingVesselLnHandler(transactionService, historyService, validate, logService, groupingVesselLnService, userIupopkService)

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

	// DN
	groupingVesselRouting.Get("/list/dn/:iupopk_id", groupingVesselDnHandler.ListGroupingVesselDn)
	groupingVesselRouting.Post("/create/dn/:iupopk_id", groupingVesselDnHandler.CreateGroupingVesselDn)
	groupingVesselRouting.Put("/update/dn/:id/:iupopk_id", groupingVesselDnHandler.EditGroupingVesselDn)
	groupingVesselRouting.Put("/update/document/dn/:id/:type/:iupopk_id", groupingVesselDnHandler.UploadDocumentGroupingVesselDn)
	groupingVesselRouting.Delete("/delete/dn/:id/:iupopk_id", groupingVesselDnHandler.DeleteGroupingVesselDn)
	groupingVesselRouting.Get("/detail/dn/:id/:iupopk_id", groupingVesselDnHandler.GetDetailGroupingVesselDn)
	groupingVesselRouting.Get("/list/dn/transaction/:iupopk_id", groupingVesselDnHandler.ListDnWithoutGroup)

	// LN
	groupingVesselRouting.Post("/create/ln/:iupopk_id", groupingVesselLnHandler.CreateGroupingVesselLn)
	groupingVesselRouting.Get("/detail/ln/:id/:iupopk_id", groupingVesselLnHandler.GetDetailGroupingVesselLn)
	groupingVesselRouting.Put("/update/ln/:id/:iupopk_id", groupingVesselLnHandler.EditGroupingVesselLn)
	groupingVesselRouting.Put("/update/document/ln/:id/:type/:iupopk_id", groupingVesselLnHandler.UploadDocumentGroupingVesselLn)
	groupingVesselRouting.Delete("/delete/ln/:id/:iupopk_id", groupingVesselLnHandler.DeleteGroupingVesselLn)
	groupingVesselRouting.Get("/list/ln/:iupopk_id", groupingVesselLnHandler.ListGroupingVesselLn)
	groupingVesselRouting.Get("/list/ln/transaction/:iupopk_id", groupingVesselLnHandler.ListLnWithoutGroup)
}
