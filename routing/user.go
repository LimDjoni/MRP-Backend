package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/user"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func UserRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	userHandler := handler.NewUserHandler(userService, userIupopkService, validate)

	userRouting := app.Group("/user") // /api

	userRouting.Post("/register", userHandler.RegisterUser)
	userRouting.Post("/login", userHandler.LoginUser)

	userValidateRouting := app.Group("/validate")

	userValidateRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	userValidateRouting.Get("/", userHandler.Validate)

	userIupopkRouting := app.Group("/user")

	userIupopkRouting.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			helper.GetEnvWithKey("USERNAME_BASIC"): helper.GetEnvWithKey("PASSWORD_BASIC"),
		},
	}))

	userIupopkRouting.Post("/create/iupopk/:user_id/:iupopk_id", userHandler.CreateUserIupopk)
	userIupopkRouting.Delete("/delete/iupopk/:user_id/:iupopk_id", userHandler.DeleteUserIupopk)

	userIupopkRouting.Put("/reset/password", userHandler.ResetPassword)

	userUpdateRouting := app.Group("/password")

	userUpdateRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	userUpdateRouting.Put("/update", userHandler.ChangePassword)
}
