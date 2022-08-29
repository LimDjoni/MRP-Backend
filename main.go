package main

import (
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/dmobarge"
	"ajebackend/model/dmovessel"
	"ajebackend/model/history"
	"ajebackend/model/logs"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbatransaction"
	"ajebackend/model/trader"
	"ajebackend/model/traderdmo"
	"ajebackend/model/transaction"
	"ajebackend/model/user"
	routing2 "ajebackend/routing"
	"ajebackend/validatorfunc"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {

	LoadEnv()
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
			&dmobarge.DmoBarge{},
			&dmovessel.DmoVessel{},
			&history.History{},
			&logs.Logs{},
			&minerba.Minerba{},
			&minerbatransaction.MinerbaTransaction{},
			&trader.Trader{},
			&traderdmo.TraderDmo{},
			&transaction.Transaction{},
			&user.User{},
		)

		fmt.Println(errMigrate)
	}

	var validate = validator.New()

	// Make Validation for Gender
	errDate := validate.RegisterValidation("DateValidation", validatorfunc.CheckDateString)
	if errDate != nil {
		fmt.Println(errDate.Error())
		fmt.Println("error validate date")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://cdf2-103-121-18-7.ap.ngrok.io",
		AllowMethods:     "GET, POST, OPTIONS, PUT, DELETE",
		AllowCredentials: true,
		AllowHeaders: "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin, Authorization",
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

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func Setup(db *gorm.DB, validate *validator.Validate, route fiber.Router) {
	routing2.TransactionRouting(db, route, validate)
	routing2.UserRouting(db, route, validate)

}
