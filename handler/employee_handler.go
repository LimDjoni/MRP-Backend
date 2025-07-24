package handler

import (
	"mrpbackend/model/employee"
	"mrpbackend/model/user"
	"mrpbackend/validatorfunc"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type employeesHandler struct {
	userService     user.Service
	employeeService employee.Service
	v               *validator.Validate
}

func NewEmployeeHandler(userService user.Service, employeesService employee.Service, v *validator.Validate) *employeesHandler {
	return &employeesHandler{
		userService,
		employeesService,
		v,
	}
}

// Master Data

func (h *employeesHandler) CreateEmployee(c *fiber.Ctx) error {
	employeesInput := new(employee.RegisterEmployeeInput)

	// Binds the request body to the Person struct
	if err := c.BodyParser(employeesInput); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	errors := h.v.Struct(*employeesInput)

	if errors != nil {
		dataErrors := validatorfunc.ValidateStruct(errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": dataErrors,
		})
	}

	createdEmployee, createdEmployeeErr := h.employeeService.CreateEmployee(*employeesInput)

	if createdEmployeeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": createdEmployeeErr.Error(),
		})
	}

	return c.Status(201).JSON(createdEmployee)
}

func (h *employeesHandler) GetEmployee(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	empCode := c.Params("empCode")

	empCodeId, err := strconv.Atoi(empCode)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	brands, brandsErr := h.employeeService.FindEmployee(uint(empCodeId))

	if brandsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": brandsErr.Error(),
		})
	}

	return c.Status(200).JSON(brands)
}

func (h *employeesHandler) GetEmployeeById(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	id := c.Params("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	heavyEquipments, heavyEquipmentsErr := h.employeeService.FindEmployeeById(uint(idInt))

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *employeesHandler) GetEmployeeByDepartmentId(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	departmentId := c.Params("departmentId")

	departmentIdInt, err := strconv.Atoi(departmentId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	employees, employeesErr := h.employeeService.FindEmployeeByDepartmentId(uint(departmentIdInt))

	if employeesErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": employeesErr.Error(),
		})
	}

	return c.Status(200).JSON(employees)
}

func (h *employeesHandler) GetListEmployee(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	page := c.Query("page")

	pageNumber, err := strconv.Atoi(page)

	if err != nil && page != "" {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if page == "" {
		pageNumber = 1
	}

	var filterEmployee employee.SortFilterEmployee

	filterEmployee.NomorKaryawan = c.Query("nomor_karyawan")
	filterEmployee.DepartmentId = c.Query("department_id")
	filterEmployee.Firstname = c.Query("firstname")
	filterEmployee.HireBy = c.Query("hire_by")
	filterEmployee.Agama = c.Query("agama")
	filterEmployee.Level = c.Query("level")
	filterEmployee.Gender = c.Query("gender")
	filterEmployee.KategoriLokalNonLokal = c.Query("kategori_lokal_non_lokal")
	filterEmployee.KategoriTriwulan = c.Query("kategori_triwulan")
	filterEmployee.Status = c.Query("status")
	filterEmployee.Kontrak = c.Query("kontrak")
	filterEmployee.RoleId = c.Query("role_id")
	filterEmployee.PositionId = c.Query("position_id")
	filterEmployee.CodeEmp = c.Query("code_emp")
	filterEmployee.Field = c.Query("field")
	filterEmployee.Sort = c.Query("sort")

	listEmployee, listEmployeeErr := h.employeeService.GetListEmployee(pageNumber, filterEmployee)

	if listEmployeeErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": listEmployeeErr.Error(),
		})
	}

	return c.Status(200).JSON(listEmployee)
}

func (h *employeesHandler) UpdateEmployee(c *fiber.Ctx) error {
	//Get User
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	responseUnauthorized := map[string]interface{}{
		"error": "unauthorized",
	}
	// Check User Login
	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(401).JSON(responseUnauthorized)
	}

	//Get Input
	inputUpdateEmployee := new(employee.UpdateEmployeeInput)
	if err := c.BodyParser(inputUpdateEmployee); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//Get ID
	id := c.Params("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	updateEmployee, err := h.employeeService.UpdateEmployee(*inputUpdateEmployee, idInt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(201).JSON(updateEmployee)
}

func (h *employeesHandler) DeleteEmployee(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid ID parameter",
			"error":   err.Error(),
		})
	}

	employee, err := h.employeeService.FindEmployeeById(uint(id))
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to find employee",
			"error":   err.Error(),
		})
	}

	// Optional: Use employee.ID for extra safety
	if _, err := h.employeeService.DeleteEmployee(employee.ID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "record not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"message": "failed to delete employee",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success delete employee",
	})
}

func (h *employeesHandler) ListDashboard(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	codeEmpId := c.Params("codeEmpId")

	codeEmpIdInt, err := strconv.Atoi(codeEmpId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	var filterDashboard employee.SortFilterDashboardEmployee

	filterDashboard.PT = c.Query("pt")
	filterDashboard.DepartmentId = c.Query("department_id")

	heavyEquipments, heavyEquipmentsErr := h.employeeService.ListDashboard(uint(codeEmpIdInt), filterDashboard)

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *employeesHandler) ListDashboardTurnover(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	codeEmpId := c.Params("codeEmpId")

	codeEmpIdInt, err := strconv.Atoi(codeEmpId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	var filterDashboard employee.SortFilterDashboardEmployeeTurnOver

	filterDashboard.PT = c.Query("pt")
	filterDashboard.Year = c.Query("year")

	heavyEquipments, heavyEquipmentsErr := h.employeeService.ListDashboardTurnover(uint(codeEmpIdInt), filterDashboard)

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}

func (h *employeesHandler) ListDashboardKontrak(c *fiber.Ctx) error {
	// Safe check and type assertion
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - token not found",
		})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil || reflect.TypeOf(claims["id"]).Kind() != reflect.Float64 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - invalid claims",
		})
	}

	// Check user existence
	userID := uint(claims["id"].(float64))
	checkUser, err := h.userService.FindUser(userID)
	if err != nil || !checkUser.IsActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized - user not found or inactive",
		})
	}

	codeEmpId := c.Params("codeEmpId")

	codeEmpIdInt, err := strconv.Atoi(codeEmpId)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "record not found",
		})
	}

	var filterDashboard employee.SortFilterDashboardEmployeeTurnOver

	filterDashboard.PT = c.Query("pt")
	filterDashboard.Year = c.Query("year")

	heavyEquipments, heavyEquipmentsErr := h.employeeService.ListDashboardKontrak(uint(codeEmpIdInt), filterDashboard)

	if heavyEquipmentsErr != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": heavyEquipmentsErr.Error(),
		})
	}

	return c.Status(200).JSON(heavyEquipments)
}
