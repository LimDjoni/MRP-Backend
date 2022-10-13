package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/logs"
	"ajebackend/model/notification"
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func NotificationRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	notificationRepository := notification.NewRepository(db)
	notificationService := notification.NewService(notificationRepository)

	logsRepository := logs.NewRepository(db)
	logsService := logs.NewService(logsRepository)

	notificationHandler := handler.NewNotificationHandler(userService, notificationService, logsService, validate)

	notificationRouting := app.Group("/notification")

	notificationRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err": err.Error(),
			})
		},
	}))

	notificationRouting.Post("/create", notificationHandler.CreateNotification)
	notificationRouting.Put("/update/read", notificationHandler.UpdateNotification)
	notificationRouting.Delete("/delete", notificationHandler.DeleteNotification)
	notificationRouting.Get("/list", notificationHandler.GetNotification)
}
