package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/aburizalpurnama/go-simple-lending/internal/custerror"
	"github.com/aburizalpurnama/go-simple-lending/internal/model"
	"github.com/aburizalpurnama/go-simple-lending/internal/payload/response"
	"github.com/aburizalpurnama/go-simple-lending/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type (
	Installment interface {
		GetListByAccount(c *fiber.Ctx) error
	}

	installmentImpl struct {
		db              *gorm.DB
		installmentRepo repository.Installment
		validate        *validator.Validate
	}
)

func NewInstallment(db *gorm.DB, validate *validator.Validate, installmentRepo repository.Installment) *installmentImpl {
	return &installmentImpl{
		db:              db,
		validate:        validate,
		installmentRepo: installmentRepo,
	}
}

func (a *installmentImpl) GetListByAccount(c *fiber.Ctx) error {
	param := c.Params("id")
	accountId, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.Base{
			Description: "invalid account id",
			Data:        nil,
		})
	}

	activeOnlyStr := c.Query("active_only", "true")
	activeOnly, err := strconv.ParseBool(activeOnlyStr)
	if err != nil {
		activeOnly = true
	}

	var installments []model.Installment
	err = a.db.Transaction(func(tx *gorm.DB) error {
		var err error
		if activeOnly {
			installments, err = a.installmentRepo.GetListAciveByAccountId(c.Context(), tx, accountId)
			if err != nil {
				return err
			}
		} else {
			installments, err = a.installmentRepo.GetListByAccountId(c.Context(), tx, accountId)
			if err != nil {
				return err
			}
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
		Data:        installments,
	})
}
