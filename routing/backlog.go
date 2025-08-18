package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/backlog"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func BackLogRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	BackLogRepository := backlog.NewRepository(db)
	BackLogService := backlog.NewService(BackLogRepository)

	backlogsHandler := handler.NewBackLogHandler(userService, BackLogService, validate)

	backlogsRouting := app.Group("/backlog")

	backlogsRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	backlogsRouting.Post("/create", backlogsHandler.CreateBackLog)

	backlogsRouting.Get("/list", backlogsHandler.GetBackLog)
	backlogsRouting.Get("/list/pagination", backlogsHandler.GetListBackLog)
	backlogsRouting.Get("/dashboard", backlogsHandler.GetListDashboardBackLog)

	backlogsRouting.Get("/detail/:id", backlogsHandler.GetBackLogById)

	backlogsRouting.Put("/update/:id", backlogsHandler.UpdateBackLog)

	backlogsRouting.Delete("/delete/:id", backlogsHandler.DeleteBackLog)

}
