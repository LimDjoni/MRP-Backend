package main

import (
	"mrpbackend/helper"
	"mrpbackend/model/adjuststock"
	"mrpbackend/model/alatberat"
	"mrpbackend/model/employee"
	"mrpbackend/model/fuelin"
	"mrpbackend/model/fuelratio"
	"mrpbackend/model/master/apd"
	"mrpbackend/model/master/bpjskesehatan"
	"mrpbackend/model/master/bpjsketenagakerjaan"
	"mrpbackend/model/master/brand"
	"mrpbackend/model/master/department"
	"mrpbackend/model/master/departmentform"
	"mrpbackend/model/master/doh"
	"mrpbackend/model/master/form"
	heavyequiment "mrpbackend/model/master/heavyequipment"
	"mrpbackend/model/master/history"
	"mrpbackend/model/master/jabatan"
	"mrpbackend/model/master/kartukeluarga"
	"mrpbackend/model/master/ktp"
	"mrpbackend/model/master/laporan"
	"mrpbackend/model/master/mcu"
	"mrpbackend/model/master/npwp"
	"mrpbackend/model/master/pendidikan"
	"mrpbackend/model/master/position"
	"mrpbackend/model/master/role"
	"mrpbackend/model/master/roleform"
	"mrpbackend/model/master/series"
	"mrpbackend/model/master/sertifikat"
	"mrpbackend/model/master/userrole"
	"mrpbackend/model/unit"
	"mrpbackend/model/user"
	"mrpbackend/model/userposition"
	routing2 "mrpbackend/routing"

	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var port string
	port = "8080"

	dbUrl := helper.GetEnvWithKey("DATABASE_URL")
	dbUrlStg := helper.GetEnvWithKey("DATABASE_URL_STAGING")

	var dsn string

	if len(dbUrl) > 0 {
		dsn = dbUrl
	} else {
		dsn = dbUrlStg
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		createDB(dsn)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect database after creation: " + err.Error())
		}
	}

	if db != nil {
		// Auto Migrate All Table
		errMigrate := db.AutoMigrate(
			&position.Position{},
			&user.User{},
			&department.Department{},
			&role.Role{},
			&kartukeluarga.KartuKeluarga{},
			&ktp.KTP{},
			&pendidikan.Pendidikan{},
			&jabatan.Jabatan{},
			&sertifikat.Sertifikat{},
			&mcu.MCU{},
			&laporan.Laporan{},
			&apd.APD{},
			&npwp.NPWP{},
			&bpjskesehatan.BPJSKesehatan{},
			&bpjsketenagakerjaan.BPJSKetenagakerjaan{},
			&history.History{},
			&form.Form{},
			&userrole.UserRole{},
			&departmentform.DepartmentForm{},
			&roleform.RoleForm{},
			&brand.Brand{},
			&heavyequiment.HeavyEquipment{},
			&series.Series{},
			&alatberat.AlatBerat{},
			&unit.Unit{},
			&employee.Employee{},
			&doh.DOH{},
			&fuelratio.FuelRatio{},
			&userposition.UserPosition{},
			&fuelin.FuelIn{},
			&adjuststock.AdjustStock{},
		)
		fmt.Println(errMigrate)
	}

	var validate = validator.New()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, OPTIONS, PUT, DELETE",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization",
		MaxAge:           2592000,
	}))

	apiV1 := app.Group("/api/v1") // /api

	Setup(db, validate, apiV1)

	app.Listen(":" + port)

}

func createDB(dsn string) {
	// Base DSN use for if there is no database (only for creating new database)
	baseDsn := helper.GetEnvWithKey("BASE_DATABASE_URL_STAGING")
	db, err := gorm.Open(postgres.Open(baseDsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to base DB:", err)
		return
	} else {
		dbName := helper.GetEnvWithKey("DATABASE_NAME")
		dbExec := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		db = db.Exec(dbExec)

		if db.Error != nil {
			fmt.Println(db.Error)
			errAssumingExist := fmt.Sprintf("Unable to create DB %s, attempting to connect assuming it exists...", dbName)
			fmt.Println(errAssumingExist)
		}
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect db")
		return
	}
}

func Setup(db *gorm.DB, validate *validator.Validate, route fiber.Router) {
	routing2.UserRouting(db, route, validate)
	routing2.MasterRouting(db, route, validate)
	routing2.AlatBeratRouting(db, route, validate)
	routing2.UnitRouting(db, route, validate)
	routing2.FuelRatioRouting(db, route, validate)
	routing2.EmployeeRouting(db, route, validate)
	routing2.StockFuelRouting(db, route, validate)
	routing2.FuelInRouting(db, route, validate)
	routing2.AdjustStockRouting(db, route, validate)
}
