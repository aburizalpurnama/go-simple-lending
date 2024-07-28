package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aburizalpurnama/go-simple-lending/internal/controller"
	"github.com/aburizalpurnama/go-simple-lending/internal/repository"
	"github.com/aburizalpurnama/go-simple-lending/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	// wire up all dependencies
	dbconn, err := initDB()
	if err != nil {
		panic(err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	accountRepo := repository.NewAccount()
	loanRepo := repository.NewLoan()
	instRepo := repository.NewInstallment()

	loanUsecase := usecase.NewLoan(dbconn, accountRepo, loanRepo, instRepo)

	accountCtrl := controller.NewAccount(dbconn, validate, accountRepo, loanRepo)
	loanCtrl := controller.NewLoan(dbconn, validate, loanUsecase, loanRepo)
	instCtrl := controller.NewInstallment(dbconn, validate, instRepo)

	app := fiber.New(fiber.Config{
		AppName: "simple-lending",
	})

	app.Get("/health-check", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	account := app.Group("/accounts")
	account.Post("/", accountCtrl.Create)
	account.Get("/:id", accountCtrl.GetDetail)

	loan := account.Group("/:id/loans")
	loan.Post("/", loanCtrl.Create)
	loan.Get("/", loanCtrl.GetListByAccount)

	account.Get("/:id/installments", instCtrl.GetListByAccount)

	if err := app.Listen(":8089"); err != nil {
		panic(err)
	}
}

func initDB() (*gorm.DB, error) {
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
		},
	)

	var err error
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=admin password=secret dbname=go-simple-lending host=localhost port=5432 sslmode=disable TimeZone=Asia/Jakarta",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
