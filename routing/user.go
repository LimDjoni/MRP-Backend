package routing

import (
	"ajebackend/handler"
	"ajebackend/model/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService, validate)

	userRouting := app.Group("/user") // /api

	userRouting.Post("/register", userHandler.RegisterUser)
	userRouting.Post("/login", userHandler.LoginUser)
}
