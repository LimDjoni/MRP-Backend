package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/unit"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func UnitRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	UnitRepository := unit.NewRepository(db)
	UnitService := unit.NewService(UnitRepository)

	unitsHandler := handler.NewUnitHandler(userService, UnitService, validate)

	unitsRouting := app.Group("/unit")

	unitsRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	unitsRouting.Post("/create", unitsHandler.CreateUnit)

	unitsRouting.Get("/list", unitsHandler.GetUnit)
	unitsRouting.Get("/list/pagination", unitsHandler.GetListUnit)

	unitsRouting.Get("/detail/:id", unitsHandler.GetUnitById)

	unitsRouting.Put("/update/:id", unitsHandler.UpdateUnit)

	unitsRouting.Delete("/delete/:id", unitsHandler.DeleteUnit)

}
