package main

import (
	"ajebackend/helper"
	routing2 "ajebackend/routing"

	// Hauling

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
		// errMigrate := db.AutoMigrate(
		// 	&dmo.Dmo{},
		// 	&dmovessel.DmoVessel{},
		// 	&history.History{},
		// 	&logs.Logs{},
		// 	&minerba.Minerba{},
		// 	&trader.Trader{},
		// 	&traderdmo.TraderDmo{},
		// 	&transaction.Transaction{},
		// 	&user.User{},
		// 	&company.Company{},
		// 	&notification.Notification{},
		// 	&notificationuser.NotificationUser{},
		// 	&production.Production{},
		// 	&groupingvesselln.GroupingVesselLn{},
		// 	&minerbaln.MinerbaLn{},
		// 	&insw.Insw{},
		// 	&destination.Destination{},
		// 	&barge.Barge{},
		// 	&country.Country{},
		// 	&currency.Currency{},
		// 	&documenttype.DocumentType{},
		// 	&industrytype.IndustryType{},
		// 	&insurancecompany.InsuranceCompany{},
		// 	&iupopk.Iupopk{},
		// 	&navycompany.NavyCompany{},
		// 	&navyship.NavyShip{},
		// 	&pabeanoffice.PabeanOffice{},
		// 	&ports.Port{},
		// 	&portinsw.PortInsw{},
		// 	&portlocation.PortLocation{},
		// 	&salessystem.SalesSystem{},
		// 	&surveyor.Surveyor{},
		// 	&unit.Unit{},
		// 	&useriupopk.UserIupopk{},
		// 	&vessel.Vessel{},
		// 	&tugboat.Tugboat{},
		// 	&counter.Counter{},
		// 	&electricassignment.ElectricAssignment{},
		// 	&electricassignmentenduser.ElectricAssignmentEndUser{},
		// 	&cafassignment.CafAssignment{},
		// 	&cafassignmentenduser.CafAssignmentEndUser{},
		// 	&rkab.Rkab{},
		// 	&coareport.CoaReport{},
		// 	&coareportln.CoaReportLn{},
		// 	&transactionrequestreport.TransactionRequestReport{},
		// 	&categoryindustrytype.CategoryIndustryType{},
		// 	&royaltyrecon.RoyaltyRecon{},
		// 	&royaltyreport.RoyaltyReport{},
		// 	// Hauling section
		// 	&contractor.Contractor{},
		// 	&isp.Isp{},
		// 	&iupopk.Iupopk{},
		// 	&jetty.Jetty{},
		// 	&pit.Pit{},
		// 	&site.Site{},
		// 	&truck.Truck{},

		// 	&transactionispjetty.TransactionIspJetty{},
		// 	&transactionjetty.TransactionJetty{},
		// 	&transactiontoisp.TransactionToIsp{},
		// 	&transactiontojetty.TransactionToJetty{},

		// 	&haulingsynchronize.HaulingSynchronize{},
		// )

		// seeding.UpdateTransactionsRoyalty(db)
		// seeding.SeedingTraderAndCompanyData(db)
		// seeding.SeedingDestination(db)
		// seeding.UpdateNaming(db)
		// seeding.UpdateTransactionsQuantity(db)
		// seedingmaster.SeedingBarge(db)
		// seedingmaster.SeedingCountry(db)
		// seedingmaster.SeedingCurrency(db)
		// seedingmaster.SeedingDocumentType(db)
		// seedingmaster.SeedingIndustryType(db)
		// seedingmaster.SeedingInsuranceCompany(db)
		// seedingmaster.SeedingIupopk(db)
		// seedingmaster.SeedingPabeanOffice(db)
		// seedingmaster.SeedingPortInsw(db)
		// seedingmaster.SeedingPortsAndLocation(db)
		// seedingmaster.SeedingSalesSystem(db)
		// seedingmaster.SeedingSurveyor(db)
		// seedingmaster.SeedingTugboat(db)
		// seedingmaster.SeedingUnit(db)
		// seedingmaster.SeedingVessel(db)
		// seedingmaster.SeedingCounter(db)
		// seedingmaster.SeedingCategoryIndustryType(db)
		// seeding.UpdateIupopk(db)
		// fmt.Println(errMigrate)
	}

	var validate = validator.New()

	errSalesSystem := validate.RegisterValidation("SalesSystemValidation", validatorfunc.CheckEnum)

	if errSalesSystem != nil {
		fmt.Println(errSalesSystem.Error())
	}

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

	errLongMonth := validate.RegisterValidation("ShortMonth", validatorfunc.CheckEnum)

	if errLongMonth != nil {
		fmt.Println(errLongMonth.Error())
	}

	errCategory := validate.RegisterValidation("CategoryValidation", validatorfunc.CheckEnum)

	if errCategory != nil {
		fmt.Println(errCategory.Error())
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
	routing2.MasterRouting(db, route, validate)
	routing2.InswRouting(db, route, validate)
	routing2.CoaReportRouting(db, route, validate)
	routing2.CoaReportLnRouting(db, route, validate)
	routing2.RkabRouting(db, route, validate)
	routing2.ElectricAssignmentRouting(db, route, validate)
	routing2.CafAssignmentRouting(db, route, validate)
	routing2.HaulingSynchronizeRouting(db, route, validate)
	routing2.HaulingTransactionRouting(db, route, validate)
}
