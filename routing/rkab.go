package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/notificationuser"
	"ajebackend/model/rkab"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func RkabRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {

	rkabRepository := rkab.NewRepository(db)
	rkabService := rkab.NewService(rkabRepository)

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

	rkabHandler := handler.NewRkabHandler(rkabService, logService, userIupopkService, historyService, notificationUserService, validate, allMasterService)

	rkabRouting := app.Group("/rkab")

	rkabRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	rkabRouting.Get("/list/:iupopk_id", rkabHandler.ListRkab)
	rkabRouting.Post("/create/:iupopk_id", rkabHandler.CreateRkab)
	rkabRouting.Delete("delete/:id/:iupopk_id", rkabHandler.DeleteRkab)
}
