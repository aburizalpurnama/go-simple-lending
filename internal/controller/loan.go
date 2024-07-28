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
	"github.com/aburizalpurnama/go-simple-lending/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	Loan interface {
		Create(c *fiber.Ctx) error
		GetListByAccount(c *fiber.Ctx) error
	}

	loanImpl struct {
		db          *gorm.DB
		loanUsecase usecase.Loan
		loanRepo    repository.Loan
		validate    *validator.Validate
	}
)

func NewLoan(db *gorm.DB, validate *validator.Validate, loanUsecase usecase.Loan, loanRepo repository.Loan) *loanImpl {
	return &loanImpl{
		db:          db,
		validate:    validate,
		loanUsecase: loanUsecase,
		loanRepo:    loanRepo,
	}
}

func (l *loanImpl) Create(c *fiber.Ctx) error {
	param := c.Params("id")
	accountId, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: "invalid account id",
			Data:        nil,
		})
	}

	var req request.CreateLoan
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: "invalid payload",
			Data:        nil,
		})
	}

	if err := req.Validate(l.validate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: err.Error(),
			Data:        nil,
		})
	}

	loan, err := l.loanUsecase.Create(c.Context(), accountId, req)
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
		Data:        loan,
	})
}

func (a *loanImpl) GetListByAccount(c *fiber.Ctx) error {
	param := c.Params("id")
	accountId, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: "invalid account id",
			Data:        nil,
		})
	}

	var loans []model.Loan
	err = a.db.Transaction(func(tx *gorm.DB) error {
		var err error
		loans, err = a.loanRepo.GetListByAccountId(c.Context(), tx, accountId)
		if err != nil {
			return err
		}

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
		Data:        loans,
	})
}
