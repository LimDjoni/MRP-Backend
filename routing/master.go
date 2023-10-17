package routing

import (
	"ajebackend/handler"
	"ajebackend/helper"
	"ajebackend/model/counter"
	"ajebackend/model/master/allmaster"
	"ajebackend/model/master/destination"
	"ajebackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"gorm.io/gorm"
)

func MasterRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	destinationRepository := destination.NewRepository(db)
	destinationService := destination.NewService(destinationRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	allMasterRepository := allmaster.NewRepository(db)
	allMasterService := allmaster.NewService(allMasterRepository)

	counterRepository := counter.NewRepository(db)
	counterService := counter.NewService(counterRepository)

	masterHandler := handler.NewMasterHandler(destinationService, userService, allMasterService, counterService, validate)

	masterRouting := app.Group("/master")

	masterRouting.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			helper.GetEnvWithKey("USERNAME_BASIC"): helper.GetEnvWithKey("PASSWORD_BASIC"),
		},
	}))

	masterRouting.Get("/global/:iupopk_id", masterHandler.GetListMaster)
	masterRouting.Put("/update/counter", masterHandler.UpdateCounter)

	masterRouting.Get("/list/trader", masterHandler.ListTrader)
	masterRouting.Get("/list/company", masterHandler.ListCompany)

	masterRouting.Post("/create/iupopk", masterHandler.CreateIupopk)
	masterRouting.Post("/create/barge", masterHandler.CreateBarge)
	masterRouting.Post("/create/tugboat", masterHandler.CreateTugboat)
	masterRouting.Post("/create/vessel", masterHandler.CreateVessel)
	masterRouting.Post("/create/portlocation", masterHandler.CreatePortLocation)
	masterRouting.Post("/create/port", masterHandler.CreatePort)
	masterRouting.Post("/create/company", masterHandler.CreateCompany)
	masterRouting.Post("/create/trader", masterHandler.CreateTrader)
	masterRouting.Post("/create/industrytype", masterHandler.CreateIndustryType)

	masterRouting.Put("/update/barge/:id", masterHandler.UpdateBarge)
	masterRouting.Put("/update/tugboat/:id", masterHandler.UpdateTugboat)
	masterRouting.Put("/update/vessel/:id", masterHandler.UpdateVessel)
	masterRouting.Put("/update/portlocation/:id", masterHandler.UpdatePortLocation)
	masterRouting.Put("/update/port/:id", masterHandler.UpdatePort)
	masterRouting.Put("/update/company/:id", masterHandler.UpdateCompany)
	masterRouting.Put("/update/trader/:id", masterHandler.UpdateTrader)
	masterRouting.Put("/update/industrytype/:id", masterHandler.UpdateIndustryType)

	masterRouting.Delete("/delete/barge/:id", masterHandler.DeleteBarge)
	masterRouting.Delete("/delete/tugboat/:id", masterHandler.DeleteTugboat)
	masterRouting.Delete("/delete/vessel/:id", masterHandler.DeleteVessel)
	masterRouting.Delete("/delete/portlocation/:id", masterHandler.DeletePortLocation)
	masterRouting.Delete("/delete/port/:id", masterHandler.DeletePort)
	masterRouting.Delete("/delete/company/:id", masterHandler.DeleteCompany)
	masterRouting.Delete("/delete/trader/:id", masterHandler.DeleteTrader)
	masterRouting.Delete("/delete/industrytype/:id", masterHandler.DeleteIndustryType)

}
