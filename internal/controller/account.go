package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	"github.com/aburizalpurnama/go-simple-lending/internal/payload/request"
	"github.com/aburizalpurnama/go-simple-lending/internal/payload/response"
	"github.com/aburizalpurnama/go-simple-lending/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	Account interface {
		Create(c *fiber.Ctx) error
	}

	accountImpl struct {
		db          *gorm.DB
		accountRepo repository.Account
		loanRepo    repository.Loan
		validate    *validator.Validate
	}
)

func NewAccount(db *gorm.DB, validate *validator.Validate, accountRepo repository.Account, loanRepo repository.Loan) *accountImpl {
	return &accountImpl{
		db:          db,
		validate:    validate,
		accountRepo: accountRepo,
		loanRepo:    loanRepo,
	}
}

func (a *accountImpl) Create(c *fiber.Ctx) error {
	var req request.CreateAccount
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: "invalid payload",
			Data:        nil,
		})
	}

	if err := req.Validate(a.validate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: err.Error(),
			Data:        nil,
		})
	}

	account := model.Account{
		Name:  req.Name,
		Limit: req.Limit,
	}

	err := a.db.Transaction(func(tx *gorm.DB) error {
		id, err := a.accountRepo.Create(c.Context(), tx, account)
		if err != nil {
			return err
		}

		account.Id = id

		return nil
	})
	if err != nil {
		var custErr *custerror.Error
		if errors.As(err, &custErr) {
			return c.Status(custErr.HttpStatusCode).JSON(response.Base{
				Description: err.Error(),
				Data:        nil,
			})
		} else {
			fmt.Printf("%+v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(response.Base{
				Description: err.Error(),
				Data:        nil,
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(response.Base{
		Description: "success",
		Data:        account,
	})
}

func (a *accountImpl) GetDetail(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: "invalid id",
			Data:        nil,
		})
	}

	var resp response.GetAccount
	err = a.db.Transaction(func(tx *gorm.DB) error {
		account, err := a.accountRepo.GetById(c.Context(), tx, id)
		if err != nil {
			return err
		}

		loans, err := a.loanRepo.GetListByAccountId(c.Context(), tx, id)
		if err != nil {
			return err
		}

		var totalLoanAmount, totalPaidAmount int
		for _, loan := range loans {
			totalLoanAmount += loan.Amount
			totalPaidAmount += loan.PaidAmount
		}

		osAmount := totalLoanAmount - totalPaidAmount

		resp.Account = account
		resp.AvailableLimit = account.Limit - osAmount
		resp.TotalLoanAmount = totalLoanAmount
		resp.TotalPaidAmount = totalPaidAmount
		resp.TotalOutstandingAmount = osAmount

		return nil
	})
	if err != nil {
		var custErr *custerror.Error
		if errors.As(err, &custErr) {
			return c.Status(custErr.HttpStatusCode).JSON(response.Base{
				Description: err.Error(),
				Data:        nil,
			})
		} else {
			log.Printf("%+v\n", err)
			return c.Status(http.StatusInternalServerError).JSON(response.Base{
				Description: err.Error(),
				Data:        nil,
			})
		}
	}

	return c.Status(http.StatusOK).JSON(response.Base{
		Description: "success",
		Data:        resp,
	})
}
