package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/notificationuser"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func CafAssignmentRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {

	cafAssignmentEndUserRepository := cafassignmentenduser.NewRepository(db)
	cafAssignmentEndUserService := cafassignmentenduser.NewService(cafAssignmentEndUserRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	historyRepository := history.NewRepository(db)
	historyService := history.NewService(historyRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	allMasterRepository := allmaster.NewRepository(db)
	allMasterService := allmaster.NewService(allMasterRepository)

	cafAssignmentHandler := handler.NewCafAssignmentHandler(cafAssignmentEndUserService, logService, userIupopkService, historyService, notificationUserService, validate, allMasterService)

	cafAssignmentRouting := app.Group("/cafassignment")

	cafAssignmentRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	cafAssignmentRouting.Post("/create/:iupopk_id", cafAssignmentHandler.CreateCafAssignment)
	cafAssignmentRouting.Get("/detail/:id/:iupopk_id", cafAssignmentHandler.DetailCafAssignment)
	cafAssignmentRouting.Put("/update/:id/:iupopk_id", cafAssignmentHandler.UpdateCafAssignment)
	cafAssignmentRouting.Delete("/delete/:id/:iupopk_id", cafAssignmentHandler.DeleteCafAssignment)
}
