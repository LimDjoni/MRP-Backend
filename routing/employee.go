package routing

import (
	"mrpbackend/handler"
	"mrpbackend/helper"
	"mrpbackend/model/employee"
	"mrpbackend/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"gorm.io/gorm"
)

func EmployeeRouting(db *gorm.DB, app fiber.Router, validate *validator.Validate) {
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	employeeRepository := employee.NewRepository(db)
	employeeService := employee.NewService(employeeRepository)

	employeeHandler := handler.NewEmployeeHandler(userService, employeeService, validate)

	employeeRouting := app.Group("/employee") // /api

	employeeRouting.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(helper.GetEnvWithKey("JWT_SECRET_KEY")),
		SigningMethod: jwtware.HS256,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
				"err":   err.Error(),
			})
		},
	}))

	employeeRouting.Post("/create", employeeHandler.CreateEmployee)
	employeeRouting.Get("/list/export/:empCode", employeeHandler.GetEmployee)
	employeeRouting.Get("/list/pagination", employeeHandler.GetListEmployee)
	employeeRouting.Get("/detail/:id", employeeHandler.GetEmployeeById)
	employeeRouting.Get("/department/:departmentId", employeeHandler.GetEmployeeByDepartmentId)
	employeeRouting.Get("/detail/employees/department/:departmentid", employeeHandler.GetEmployeeByDepartmentId)
	employeeRouting.Put("/update/:id", employeeHandler.UpdateEmployee)
	employeeRouting.Delete("/delete/:id", employeeHandler.DeleteEmployee)

	employeeRouting.Get("/list/dashboard/:codeEmpId", employeeHandler.ListDashboard)
	employeeRouting.Get("/list/dashboardTurnOver/:codeEmpId", employeeHandler.ListDashboardTurnover)
	employeeRouting.Get("/list/dashboardKontrak/:codeEmpId", employeeHandler.ListDashboardKontrak)

}
