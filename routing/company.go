package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/logs"
	"ajebackend/model/master/company"
	"ajebackend/model/master/trader"
	"ajebackend/model/useriupopk"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func CompanyRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	traderRepository := trader.NewRepository(db)
	traderService := trader.NewService(traderRepository)

	companyRepository := company.NewRepository(db)
	companyService := company.NewService(companyRepository)

	logRepository := logs.NewRepository(db)
	logService := logs.NewService(logRepository)

	userIupopkRepository := useriupopk.NewRepository(db)
	userIupopkService := useriupopk.NewService(userIupopkRepository)

	companyHandler := handler.NewCompanyHandler(companyService, traderService, logService, validate, userIupopkService)

	companyRouting := app.Group("/company")

	companyRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	companyRouting.Get("/list", companyHandler.ListCompany)
	companyRouting.Get("/detail/:id", companyHandler.DetailCompany)
}
