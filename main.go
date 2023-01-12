package main

import (
	"ajebackend/helper"
	"ajebackend/model/company"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/history"
	"ajebackend/model/insw"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/production"
	"ajebackend/model/trader"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	"ajebackend/model/vessel"
	routing2 "ajebackend/routing"
	"ajebackend/seeding"
	"ajebackend/validatorfunc"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var port string
	if len(os.Getenv("PORT")) < 2 {
		port = "8080"
	} else {
		port = helper.GetEnvWithKey("PORT")
	}

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
	}

	if db != nil {
		// Auto Migrate All Table
		errMigrate := db.AutoMigrate(
			&dmo.Dmo{},
			&dmovessel.DmoVessel{},
			&history.History{},
			&logs.Logs{},
			&minerba.Minerba{},
			&trader.Trader{},
			&traderdmo.TraderDmo{},
			&transaction.Transaction{},
			&user.User{},
			&company.Company{},
			&notification.Notification{},
			&notificationuser.NotificationUser{},
			&production.Production{},
			&vessel.Vessel{},
			&groupingvesselln.GroupingVesselLn{},
			&minerbaln.MinerbaLn{},
			&insw.Insw{},
		)

		seeding.UpdateTransactionsRoyalty(db)
		seeding.SeedingTraderAndCompanyData(db)

		fmt.Println(errMigrate)
	}

	var validate = validator.New()

	// Make Validation for Date
	errDate := validate.RegisterValidation("DateValidation", validatorfunc.CheckDateString)
	if errDate != nil {
		fmt.Println(errDate.Error())
	}

	// Make Validation for Period
	errPeriod := validate.RegisterValidation("PeriodValidation", validatorfunc.ValidationPeriod)

	if errPeriod != nil {
		fmt.Println(errPeriod.Error())
	}

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
		fmt.Println("Success connect base db")
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
	routing2.TransactionRouting(db, route, validate)
	routing2.UserRouting(db, route, validate)
	routing2.MinerbaRouting(db, route, validate)
	routing2.DmoRouting(db, route, validate)
	routing2.TraderRouting(db, route, validate)
	routing2.CompanyRouting(db, route, validate)
	routing2.NotificationRouting(db, route, validate)
	routing2.ProductionRouting(db, route, validate)
	routing2.ReportRouting(db, route, validate)
	routing2.GroupingVesselLnRouting(db, route, validate)
}
