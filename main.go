package main

import (
	"ajebackend/helper"
	"ajebackend/model/company"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbatransaction"
	"ajebackend/model/notification"
	"ajebackend/model/notificationuser"
	"ajebackend/model/trader"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	routing2 "ajebackend/routing"
	"ajebackend/seeding"
	"ajebackend/validatorfunc"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
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
		)

		db.Migrator().RenameColumn(&transaction.Transaction{}, "ship_name", "tugboat_name")

		seeding.UpdateTransactionsRoyalty(db)
		seeding.SeedingTraderData(db)

		errDropTable := db.Migrator().DropTable(
			&minerbatransaction.MinerbaTransaction{},
		)

		if errDropTable != nil {
			fmt.Println(errDropTable)
		}
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

	file, err := os.OpenFile("./logging.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, OPTIONS, PUT, DELETE",
		AllowCredentials: true,
		AllowHeaders: "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization",
		MaxAge:           2592000,
	}), logger.New(logger.Config{
		Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n query params : ${queryParams}\n body: ${body}\n response body: ${resBody}\n\n",
		TimeFormat: "02-Jan-2006 03:04:05 PM",
		TimeZone:   "Asia/Jakarta",
		Output: file,
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
		db = db.Exec("CREATE DATABASE deli_aje_development;")

		if db.Error != nil {
			fmt.Println(db.Error)
			fmt.Println("Unable to create DB deli_aje_development, attempting to connect assuming it exists...")
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
}
