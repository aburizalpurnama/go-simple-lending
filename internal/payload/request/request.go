package request

import (
	"github.com/go-playground/validator/v10"
)

type (
	CreateAccount struct {
		Name  string `validate:"required" json:"name"`
		Limit int    `validate:"required" json:"limit"`
	}

	CreateLoan struct {
		Limit int `validate:"required" json:"amount"`
	}

	CreatePayment struct {
		Amount int `validate:"required" json:"amount"`
	}
)

func (c *CreateAccount) Validate(validate *validator.Validate) (err error) {
	return validate.Struct(c)
}

func (c *CreateLoan) Validate(validate *validator.Validate) (err error) {
	return validate.Struct(c)
}

func (c *CreatePayment) Validate(validate *validator.Validate) (err error) {
	return validate.Struct(c)
}
