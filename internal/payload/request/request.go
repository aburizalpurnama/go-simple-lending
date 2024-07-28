package request

import (
	"github.com/go-playground/validator/v10"
)

type (
	CreateAccount struct {
		Name  string `validate:"required" json:"name"`
		Limit int    `validate:"required,gt=0" json:"limit"`
	}

	CreateLoan struct {
		Amount int `validate:"required,gt=0" json:"amount"`
		Tenor  int `validate:"required,oneof=1 3 6 12 24" json:"tenor"`
	}

	CreatePayment struct {
		Amount int `validate:"required,gt=0" json:"amount"`
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
