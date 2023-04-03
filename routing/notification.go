package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func NotificationRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	notificationRepository := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepository)

	notificationUserRepository := notificationuser.NewRepository(db)
	notificationUserService := notificationuser.NewService(notificationUserRepository)

	logsRepository := logs.NewRepository(db)
	logsService := logs.NewService(logsRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	notificationHandler := handler.NewNotificationHandler(notificationService, notificationUserService, logsService, validate, userIupopkService)

	notificationRouting := app.Group("/notification")

	notificationRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	notificationRouting.Post("/create/:iupopk_id", notificationHandler.CreateNotification)
	notificationRouting.Put("/update/:iupopk_id", notificationHandler.UpdateNotification)
	notificationRouting.Delete("/delete/:iupopk_id", notificationHandler.DeleteNotification)
	notificationRouting.Get("/list/:iupopk_id", notificationHandler.GetNotification)
}
