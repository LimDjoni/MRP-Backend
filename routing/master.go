package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/allmaster"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func MasterRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	allMasterRepository := allmaster.NewRepository(db)
	allMasterService := allmaster.NewService(allMasterRepository)

	masterHandler := handler.NewMasterHandler(userService, allMasterService, validate)

	masterRouting := app.Group("/master")

	masterRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	masterRouting.Post("/create/brand", masterHandler.CreateBrand)
	masterRouting.Post("/create/heavyequipment", masterHandler.CreateHeavyEquipment)
	masterRouting.Post("/create/series", masterHandler.CreateSeries)
	masterRouting.Post("/create/kartukeluarga", masterHandler.CreateKartuKeluarga)
	masterRouting.Post("/create/ktp", masterHandler.CreateKTP)
	masterRouting.Post("/create/pendidikan", masterHandler.CreatePendidikan)
	masterRouting.Post("/create/doh", masterHandler.CreateDOH)
	masterRouting.Post("/create/jabatan", masterHandler.CreateJabatan)
	masterRouting.Post("/create/sertifikat", masterHandler.CreateSertifikat)
	masterRouting.Post("/create/mcu", masterHandler.CreateMCU)
	masterRouting.Post("/create/history", masterHandler.CreateHistory)

	masterRouting.Get("/list/brand", masterHandler.GetBrand)
	masterRouting.Get("/list/heavyequipment", masterHandler.GetHeavyEquipment)
	masterRouting.Get("/list/series", masterHandler.GetSeries)
	masterRouting.Get("/list/department", masterHandler.GetDepartment)
	masterRouting.Get("/list/role", masterHandler.GetRole)
	masterRouting.Get("/list/position", masterHandler.GetPosition)
	masterRouting.Get("/list/expireddoh", masterHandler.GetDohKontrak)

	masterRouting.Get("/detail/brand/:id", masterHandler.GetBrandById)
	masterRouting.Get("/detail/heavyequipment/:id", masterHandler.GetHeavyEquipmentById)
	masterRouting.Get("/detail/heavyequipment/brand/:brandId", masterHandler.GetHeavyEquipmentByBrandId)
	masterRouting.Get("/detail/series/:id", masterHandler.GetSeriesById)
	masterRouting.Get("/detail/series/heavyequipment/:brandId/:heavyequipmentId", masterHandler.GetSeriesByBrandAndEquipmentdID)

	masterRouting.Put("/update/doh/:id", masterHandler.UpdateDOH)
	masterRouting.Delete("/delete/doh/:id", masterHandler.DeleteDOH)

	masterRouting.Put("/update/jabatan/:id", masterHandler.UpdateJabatan)
	masterRouting.Delete("/delete/jabatan/:id", masterHandler.DeleteJabatan)

	masterRouting.Put("/update/sertifikat/:id", masterHandler.UpdateSertifikat)
	masterRouting.Delete("/delete/sertifikat/:id", masterHandler.DeleteSertifikat)

	masterRouting.Put("/update/mcu/:id", masterHandler.UpdateMCU)
	masterRouting.Delete("/delete/mcu/:id", masterHandler.DeleteMCU)
	masterRouting.Get("/list/mcuberkala", masterHandler.GetMCUBerkala)

	masterRouting.Put("/update/history/:id", masterHandler.UpdateHistory)
	masterRouting.Delete("/delete/history/:id", masterHandler.DeleteHistory)

	masterRouting.Get("/sidebar/:userId", masterHandler.GenerateSideBar)

	userRoleRouting := app.Group("/userRole") // /api

	userRoleRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	userRoleRouting.Post("/create", masterHandler.CreateUserRole)
	userRoleRouting.Get("/list", masterHandler.GetUserRole)
	userRoleRouting.Get("/detail/:id", masterHandler.GetUserRoleById)

}
