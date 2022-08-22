package main

import (
	"ajebackend/model/dmo"
	"ajebackend/model/dmotongkang"
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
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	LoadEnv()

	dbUrl := ""
	dbUrlStg := "host=localhost user=postgres password=postgres dbname=deli_aje_development port=5432 sslmode=disable TimeZone=Asia/Shanghai"

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
			&dmotongkang.DmoTongkang{},
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
	apiV1 := app.Group("/api/v1") // /api

	routing2.TransactionRouting(db, apiV1, validate)
	routing2.UserRouting(db, apiV1, validate)

	log.Fatal(app.Listen(":3000"))
}

func createDB(dsn string) {

	// Base DSN use for if there is no database (only for creating new database)

	baseDsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
