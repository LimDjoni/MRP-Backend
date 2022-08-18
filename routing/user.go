package routing

import (
	"ajebackend/handler"
	"ajebackend/model/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRouting(db *gorm.DB, app fiber.Router) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	app.Post("/register", userHandler.RegisterUser)
	app.Post("/login", userHandler.LoginUser)
}
